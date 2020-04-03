// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

// Package tpdu provides the TPDU type and conversions to and from its binary
// form.
package tpdu

import (
	"bytes"
	"encoding/binary"

	"github.com/warthog618/sms/encoding/gsm7"
)

// TPDU represents all SMS TPDUs.
type TPDU struct {
	// Direction indicates whether the TPDU is mobile originated (MO) or
	// terminated (MT).
	Direction Direction

	// FirstOctet is the first octet of all TPDUs.
	FirstOctet FirstOctet

	// OA contains the TP-OA Originating Address field.
	//
	// Only applies to SMS-DELIVER
	OA Address

	// FCS contains the TP-FCS Failure Cause field.
	//
	// Only applies to SMS-DELIVER-REPORT and MS-SUBMIT-REPORT
	FCS byte

	// MR contains the TP-MP Message Reference field.
	//
	// Only applies to SMS-COMMAND, SMS-SUBMIT and SMS-STATUS-REPORT
	MR byte

	// CT contains the TP-CT Command Type field.
	//
	// Only applies to SMS-COMMAND
	CT byte

	// MN contains the TP-MN Message Number field.
	//
	// Only applies to SMS-COMMAND
	MN byte

	// DA contains the TP-DA Destination Address field.
	//
	// Only applies to SMS-COMMAND and SMS-SUBMIT
	DA Address

	// RA contains the TP-RA Recipient Address field.
	//
	// Only applies to SMS-STATUS-REPORT
	RA Address

	// PI contains the TP-PI Parameter Indicator field.
	//
	//  Only applies to SMS-DELIVER-REPORT and SMS-SUBMIT-REPORT
	PI PI

	// SCTS contains the TP-SCTS Service Center Time Stamp field.
	//
	// The SCTS timestamp indicates the time the SMS was sent.
	// The time is the originator's local time, the timezone of which may
	// differ from the receiver's.
	//
	// Only applies to SMS-DELIVER, SMS-SUBMIT-REPORT and SMS-STATUS-REPORT
	SCTS Timestamp

	// DT contains the TP-DT Discharge Time field.
	//
	// Only applies to SMS-STATUS-REPORT
	DT Timestamp

	// ST contains the TP-ST Status field.
	//
	// Only applies to SMS-STATUS-REPORT
	ST byte

	// PID contains the TP-PID field.
	PID byte

	// DCS contains the TP-DCS Data Coding Scheme field.
	DCS DCS

	// VP contains the TP-VP Validity Period field.
	//
	//  Only applies to SMS-SUBMIT
	VP ValidityPeriod

	// UDH contains the TP-UDH User Data Header field.
	UDH UserDataHeader

	// UD contains the short message from the User Data.
	//
	// It does not include the User Data Header, which is provided separately
	// in the UDH.
	// The interpretation of UD depends on the Alphabet:
	// For Alpha7Bit, UD is an array of GSM7 septets, each septet stored in the
	// lower 7 bits of a byte.
	//  These have NOT been converted to the corresponding UTF8.
	//  Use the gsm7 package to convert to UTF8.
	// For AlphaUCS2, UD is an array of UCS2 characters packed into a byte
	// array in Big Endian.
	//  These have NOT been converted to the corresponding UTF8.
	//  Use the usc2 package to convert to UTF8.
	// For Alpha8Bit, UD contains the raw octets.
	UD UserData
}

// New creates a new TPDU
func New(options ...Option) (*TPDU, error) {
	t := TPDU{OA: NewAddress(), DA: NewAddress(), RA: NewAddress()}
	for _, option := range options {
		err := option.ApplyTPDUOption(&t)
		if err != nil {
			return nil, err
		}
	}
	return &t, nil
}

// NewDeliver creates a new TPDU of type SmsDeliver.
func NewDeliver(options ...Option) (*TPDU, error) {
	options = append([]Option{SmsDeliver}, options...)
	return New(options...)
}

// NewSubmit creates a new TPDU of type SmsSubmit.
func NewSubmit(options ...Option) (*TPDU, error) {
	options = append([]Option{SmsSubmit}, options...)
	return New(options...)
}

// Alphabet returns the alphabet field from the DCS of the SMS TPDU.
func (t *TPDU) Alphabet() (Alphabet, error) {
	return t.DCS.Alphabet()
}

// ConcatInfo extracts the segmentation info contained in the provided User
// Data Header.
func (t *TPDU) ConcatInfo() (segments, seqno, mref int, ok bool) {
	return t.UDH.ConcatInfo()
}

// MTI returns the MessageType from the first octet of the SMS TPDU.
func (t *TPDU) MTI() MessageType {
	return t.FirstOctet.MTI()
}

// Counter provides a reference couunter that is incremented every time Count
// is called.
type Counter interface {
	Count() int
}

type segmentationConfig struct {
	// concat IE factory
	ief func(msgCount int, segCount int, segment int) InformationElement

	// concat ref generator
	cr Counter

	// MR generator
	mr Counter
}

// SegmentationOption provides an option to modify the behaviour of segmentation.
type SegmentationOption func(*segmentationConfig)

// Segment returns the set of SMS TPDUs required to transmit the message.
//
// The TPDU acts as the template for the generated TPDUs and provides all the
// fields in the resulting TPDUs, other than the UD, which is populated using
// the message.  For multi-part messages, the UDH provided in the TPDU is
// extended with a concatenation IE. The TPDU UDH must not contain a
// concatenation IE (ID 0 or 8) or the resulting TPDUs will be non-conformant.
func (t TPDU) Segment(msg []byte, options ...SegmentationOption) []TPDU {
	if len(msg) == 0 {
		return nil
	}
	cfg := segmentationConfig{newInfoElement, nil, nil}
	for _, o := range options {
		o(&cfg)
	}
	bs := t.UDBlockSize()
	if len(msg) <= bs {
		// single segment
		t.UD = msg
		if cfg.mr != nil {
			t.MR = byte(cfg.mr.Count())
		}
		return []TPDU{t}
	}
	// add contcat IE and recalc bs
	t.SetUDH(append(t.UDH, cfg.ief(0, 0, 0)))
	bs = t.UDBlockSize()
	t.UDH = t.UDH[:len(t.UDH)-1]
	alpha, _ := t.Alphabet()
	chunks := chunk(msg, alpha, bs)
	count := len(chunks)
	pdus := make([]TPDU, count)
	concatRef := 1
	if cfg.cr != nil {
		concatRef = cfg.cr.Count()
	}
	for i := 0; i < count; i++ {
		pdus[i] = t
		if cfg.mr != nil {
			pdus[i].MR = byte(cfg.mr.Count())
		}
		udh := append(t.UDH[:0:0], t.UDH...)
		udh = append(udh, cfg.ief(concatRef, count, i+1))
		pdus[i].SetUDH(udh)
		pdus[i].UD = chunks[i]
	}
	return pdus
}

// With16BitConcatRef specifies the usage of concat IEs with 16 bit reference
// numbers (ID=8).
//
// By default 8bit reference numbers are used.
var With16BitConcatRef = func(so *segmentationConfig) {
	so.ief = newInfoElement16bit
}

// WithMR provides an MR generator to provide the TP-MR field for TPDUs.
//
// By default the MR is copied from the template TPDU.
func WithMR(mr Counter) SegmentationOption {
	return func(so *segmentationConfig) {
		so.mr = mr
	}
}

// WithConcatRef provides a generator to provide the reference for concatenation IEs.
//
// By default the field is set to 1, which is only suitable for one-off messages.
func WithConcatRef(cr Counter) SegmentationOption {
	return func(so *segmentationConfig) {
		so.cr = cr
	}
}

// SetDCS sets the dcs field and the corresponding bit of the PI.
func (t *TPDU) SetDCS(dcs byte) {
	t.PI |= PiDCS
	t.DCS = DCS(dcs)
}

// SetPID sets the TPDU pid field and the corresponding bit of the PI.
func (t *TPDU) SetPID(pid byte) {
	t.PI |= PiPID
	t.PID = pid
}

// SetVP sets the validity period and the corresponding VPF bits
// in the firstOctet.
func (t *TPDU) SetVP(vp ValidityPeriod) {
	t.FirstOctet &^= FoVPFMask
	t.FirstOctet |= (FirstOctet(vp.Format<<FoVPFShift) & FoVPFMask)
	t.VP = vp
}

// SetUD sets the TPDU ud field and the corresponding bit of the PI.
func (t *TPDU) SetUD(ud UserData) {
	t.UD = ud
	if ud == nil {
		t.PI &^= PiUDL
	} else {
		t.PI |= PiUDL
	}
}

// SetUDH sets the User Data Header of the TPDU and the TP-UDHI flag.
func (t *TPDU) SetUDH(udh UserDataHeader) {
	t.UDH = udh
	if udh == nil {
		t.FirstOctet &^= FoUDHI
	} else {
		t.PI |= PiUDL
		t.FirstOctet |= FoUDHI
	}
}

// SmsType returns the type of SMS-TPDU this TPDU represents.
func (t *TPDU) SmsType() SmsType {
	return smsType(t.FirstOctet.MTI(), t.Direction)
}

// SetSmsType returns the type of SMS-TPDU this TPDU represents.
func (t *TPDU) SetSmsType(st SmsType) error {
	if st < 0 || st > SmsCommand {
		return ErrInvalid
	}
	t.Direction = st.Direction()
	t.FirstOctet = t.FirstOctet.WithMTI(st.MTI())
	return nil
}

// UDBlockSize returns the maximum size of a block of UserData that can fit in
// this TPDU.
//
// The interpretation of the size depends on the encoding - for 7bit encoding
// it is the number of septets. For all other encodings it is the number of
// octets.
func (t *TPDU) UDBlockSize() int {
	var bs int
	switch t.SmsType() {
	case SmsSubmit, SmsDeliver:
		bs = 140
	case SmsCommand:
		bs = 146 // conservative
		// precise answer depends on variable length fields...
	case SmsSubmitReport:
		if t.FCS == 0 {
			bs = 152 // for RP-ACK
		} else {
			bs = 151 // for RP-ERROR
		}
	case SmsDeliverReport:
		if t.FCS == 0 {
			bs = 159 // for RP-ACK
		} else {
			bs = 158 // for RP-ERROR
		}
	case SmsStatusReport:
		bs = 131 // conservative
		// precise answer depends on variable length fields...
	}
	alpha, _ := t.Alphabet()
	udhl := t.UDHL()
	if alpha == Alpha7Bit {
		// work in septets
		bs = (bs * 8) / 7
		if udhl == 0 {
			return bs
		}
		// remove septets used by UDH, including UDHL and fill bits
		bs -= ((udhl+1)*8 + 6) / 7
		return bs
	}
	if udhl > 0 {
		bs -= (udhl + 1)
	}
	if alpha == AlphaUCS2 {
		bs &^= 0x1
	}
	return bs
}

// UDHI returns the User Data Header Indicator bit from the SMS TPDU first
// octet.
//
// This is generally the same as testing the length of the udh - unless the dcs
// has been intentionally overwritten to create an inconsistency.
func (t *TPDU) UDHI() bool {
	return t.FirstOctet.UDHI()
}

// UDHL returns the encoded length of the UDH, not including the UDHL itself.
func (t *TPDU) UDHL() int {
	return t.UDH.UDHL()
}

// MarshalBinary marshals a SMS TPDU into the corresponding byte array.
func (t *TPDU) MarshalBinary() (dst []byte, err error) {
	st := smsType(t.FirstOctet.MTI(), t.Direction)
	switch st {
	case SmsDeliver:
		dst, err = t.marshalDeliver()
	case SmsDeliverReport:
		dst, err = t.marshalDeliverReport()
	case SmsSubmitReport:
		dst, err = t.marshalSubmitReport()
	case SmsSubmit:
		dst, err = t.marshalSubmit()
	case SmsStatusReport:
		dst, err = t.marshalStatusReport()
	case SmsCommand:
		dst, err = t.marshalCommand()
	default:
		return nil, ErrUnsupportedSmsType(st)
	}
	if err != nil {
		err = EncodeError(st.String(), err)
	}
	return
}

func (t *TPDU) marshalCommand() ([]byte, error) {
	da, err := t.DA.MarshalBinary()
	if err != nil {
		return nil, EncodeError("da", err)
	}
	cdl := len(t.UD)
	l := 6 + len(da) + cdl
	b := make([]byte, 0, l)
	b = append(b, byte(t.FirstOctet), t.MR, t.PID, t.CT, t.MN)
	b = append(b, da...)
	b = append(b, byte(cdl))
	b = append(b, t.UD...)
	return b, nil
}

func (t *TPDU) marshalDeliver() ([]byte, error) {
	oa, err := t.OA.MarshalBinary()
	if err != nil {
		return nil, EncodeError("oa", err)
	}
	scts, err := t.SCTS.MarshalBinary()
	if err != nil {
		return nil, EncodeError("scts", err)
	}
	ud, err := t.encodeUserData()
	if err != nil {
		return nil, EncodeError("ud", err)
	}
	l := 3 + len(oa) + len(scts) + len(ud)
	b := make([]byte, 0, l)
	b = append(b, byte(t.FirstOctet))
	b = append(b, oa...)
	b = append(b, t.PID, byte(t.DCS))
	b = append(b, scts...)
	b = append(b, ud...)
	return b, nil
}

func (t *TPDU) marshalDeliverReport() ([]byte, error) {
	ud := []byte{}
	if t.PI.UDL() {
		var err error
		ud, err = t.encodeUserData()
		if err != nil {
			return nil, EncodeError("ud", err)
		}
	}
	l := 5 + len(ud) // assume FCS, PID and DCS
	b := make([]byte, 0, l)
	b = append(b, byte(t.FirstOctet))
	if t.FCS != 0 {
		b = append(b, t.FCS)
	}
	b = append(b, byte(t.PI))
	if t.PI.PID() {
		b = append(b, t.PID)
	}
	if t.PI.DCS() {
		b = append(b, byte(t.DCS))
	}
	b = append(b, ud...)
	return b, nil
}

func (t *TPDU) marshalStatusReport() ([]byte, error) {
	ra, err := t.RA.MarshalBinary()
	if err != nil {
		return nil, EncodeError("ra", err)
	}
	scts, err := t.SCTS.MarshalBinary()
	if err != nil {
		return nil, EncodeError("scts", err)
	}
	dt, err := t.DT.MarshalBinary()
	if err != nil {
		return nil, EncodeError("dt", err)
	}
	var ud []byte
	if t.PI.UDL() {
		ud, err = t.encodeUserData()
		if err != nil {
			return nil, EncodeError("ud", err)
		}
	}
	l := 6 + len(ra) + len(scts) + len(dt) + len(ud) // assume PID and DCS
	b := make([]byte, 0, l)
	b = append(b, byte(t.FirstOctet), t.MR)
	b = append(b, ra...)
	b = append(b, scts...)
	b = append(b, dt...)
	b = append(b, t.ST)
	if t.PI == 0x00 {
		return b, nil
	}
	b = append(b, byte(t.PI))
	if t.PI.PID() {
		b = append(b, t.PID)
	}
	if t.PI.DCS() {
		b = append(b, byte(t.DCS))
	}
	b = append(b, ud...)
	return b, nil
}

func (t *TPDU) marshalSubmit() ([]byte, error) {
	da, err := t.DA.MarshalBinary()
	if err != nil {
		return nil, EncodeError("da", err)
	}
	ud, err := t.encodeUserData()
	if err != nil {
		return nil, EncodeError("ud", err)
	}
	var vp []byte
	if t.VP.Format != VpfNotPresent {
		vp, err = t.VP.MarshalBinary()
		if err != nil {
			return nil, EncodeError("vp", err)
		}
	}
	l := 4 + len(da) + len(ud) + len(vp)
	b := make([]byte, 0, l)
	b = append(b, byte(t.FirstOctet), t.MR)
	b = append(b, da...)
	b = append(b, t.PID, byte(t.DCS))
	b = append(b, vp...)
	b = append(b, ud...)
	return b, nil
}

func (t *TPDU) marshalSubmitReport() ([]byte, error) {
	scts, err := t.SCTS.MarshalBinary()
	if err != nil {
		return nil, EncodeError("scts", err)
	}
	var ud []byte
	if t.PI.UDL() {
		ud, err = t.encodeUserData()
		if err != nil {
			return nil, EncodeError("ud", err)
		}
	}
	l := 5 + len(scts) + len(ud) // assume PID and DCS
	b := make([]byte, 0, l)
	b = append(b, byte(t.FirstOctet))
	if t.FCS != 0 {
		b = append(b, t.FCS)
	}
	b = append(b, byte(t.PI))
	b = append(b, scts...)
	if t.PI.PID() {
		b = append(b, t.PID)
	}
	if t.PI.DCS() {
		b = append(b, byte(t.DCS))
	}
	b = append(b, ud...)
	return b, nil
}

// UnmarshalBinary unmarshals a SMS TPDU from the corresponding byte array.
//
// In the case of error the TPDU will be partially unmarshalled, up to the
// point that the decoding error was detected.
func (t *TPDU) UnmarshalBinary(src []byte) (err error) {
	if len(src) < 1 {
		return DecodeError("tpdu.firstOctet", 0, ErrUnderflow)
	}
	t.FirstOctet = FirstOctet(src[0])
	st := smsType(t.FirstOctet.MTI(), t.Direction)
	switch st {
	case SmsDeliver:
		err = t.unmarshalDeliver(src[1:])
	case SmsDeliverReport:
		err = t.unmarshalDeliverReport(src[1:])
	case SmsSubmitReport:
		err = t.unmarshalSubmitReport(src[1:])
	case SmsSubmit:
		err = t.unmarshalSubmit(src[1:])
	case SmsStatusReport:
		err = t.unmarshalStatusReport(src[1:])
	case SmsCommand:
		err = t.unmarshalCommand(src[1:])
	default:
		t.FirstOctet = 0
		return DecodeError("tpdu.firstOctet", 0, ErrUnsupportedSmsType(st))
	}
	if err != nil {
		return DecodeError(st.String(), 1, err)
	}
	return nil
}

func (t *TPDU) unmarshalCommand(src []byte) (err error) {
	b := bytes.NewBuffer(src)
	t.MR, err = b.ReadByte()
	if err != nil {
		return DecodeError("mr", len(src)-b.Len(), err)
	}
	t.PID, err = b.ReadByte()
	if err != nil {
		return DecodeError("pid", len(src)-b.Len(), err)
	}
	t.CT, err = b.ReadByte()
	if err != nil {
		return DecodeError("ct", len(src)-b.Len(), err)
	}
	t.MN, err = b.ReadByte()
	if err != nil {
		return DecodeError("mn", len(src)-b.Len(), err)
	}
	n, err := t.DA.UnmarshalBinary(b.Bytes())
	if err != nil {
		return DecodeError("da", len(src)-b.Len(), err)
	}
	b.Next(n)
	t.DCS = Dcs8BitData // force TPDU to interpret UD as 8bit, if not set already
	err = t.decodeUserData(b.Bytes())
	if err != nil {
		return DecodeError("ud", len(src)-b.Len(), err)
	}
	return nil
}

func (t *TPDU) unmarshalDeliver(src []byte) (err error) {
	n, err := t.OA.UnmarshalBinary(src)
	if err != nil {
		return DecodeError("oa", 0, err)
	}
	b := bytes.NewBuffer(src[n:])
	t.PID, err = b.ReadByte()
	if err != nil {
		return DecodeError("pid", len(src)-b.Len(), err)
	}
	dcs, err := b.ReadByte()
	if err != nil {
		return DecodeError("dcs", len(src)-b.Len(), err)
	}
	t.DCS = DCS(dcs)
	err = t.SCTS.UnmarshalBinary(b.Bytes())
	if err != nil {
		return DecodeError("scts", len(src)-b.Len(), err)
	}
	b.Next(7)
	err = t.decodeUserData(b.Bytes())
	if err != nil {
		return DecodeError("ud", len(src)-b.Len(), err)
	}
	return nil
}

func (t *TPDU) unmarshalDeliverReport(src []byte) error {
	ri := 0
	if len(src) <= ri {
		return DecodeError("fcs", ri, ErrUnderflow)
	}
	t.FCS = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("pi", ri, ErrUnderflow)
	}
	t.PI = PI(src[ri])
	ri++
	if t.PI.PID() {
		if len(src) <= ri {
			return DecodeError("pid", ri, ErrUnderflow)
		}
		t.PID = src[ri]
		ri++
	}
	if t.PI.DCS() {
		if len(src) <= ri {
			return DecodeError("dcs", ri, ErrUnderflow)
		}
		t.DCS = DCS(src[ri])
		ri++
	}
	if t.PI.UDL() {
		err := t.decodeUserData(src[ri:])
		if err != nil {
			return DecodeError("ud", ri, err)
		}
	}
	return nil
}

func (t *TPDU) unmarshalStatusReport(src []byte) error {
	ri := 0
	if len(src) <= ri {
		return DecodeError("mr", ri, ErrUnderflow)
	}
	t.MR = src[ri]
	ri++
	n, err := t.RA.UnmarshalBinary(src[ri:])
	if err != nil {
		return DecodeError("ra", ri, err)
	}
	ri += n
	if len(src) < ri+7 {
		return DecodeError("scts", ri, ErrUnderflow)
	}
	err = t.SCTS.UnmarshalBinary(src[ri : ri+7])
	if err != nil {
		return DecodeError("scts", ri, err)
	}
	ri += 7
	if len(src) < ri+7 {
		return DecodeError("dt", ri, ErrUnderflow)
	}
	err = t.DT.UnmarshalBinary(src[ri : ri+7])
	if err != nil {
		return DecodeError("dt", ri, err)
	}
	ri += 7
	if len(src) <= ri {
		return DecodeError("st", ri, ErrUnderflow)
	}
	t.ST = src[ri]
	ri++
	if len(src) > ri {
		return t.unmarshalSROptionals(ri, src)
	}
	return nil
}

// unmarshal the optional fields at the end of the StatusReport TPDU.
func (t *TPDU) unmarshalSROptionals(ri int, src []byte) error {
	t.PI = PI(src[ri])
	ri++
	if t.PI.PID() {
		if len(src) <= ri {
			return DecodeError("pid", ri, ErrUnderflow)
		}
		t.PID = src[ri]
		ri++
	}
	if t.PI.DCS() {
		if len(src) <= ri {
			return DecodeError("dcs", ri, ErrUnderflow)
		}
		t.DCS = DCS(src[ri])
		ri++
	}
	if t.PI.UDL() {
		err := t.decodeUserData(src[ri:])
		if err != nil {
			return DecodeError("ud", ri, err)
		}
	}
	return nil
}

func (t *TPDU) unmarshalSubmit(src []byte) error {
	if len(src) < 1 {
		return DecodeError("mr", 0, ErrUnderflow)
	}
	t.MR = src[0]
	ri := 1
	n, err := t.DA.UnmarshalBinary(src[ri:])
	if err != nil {
		return DecodeError("da", ri, err)
	}
	ri += n
	if len(src) <= ri {
		return DecodeError("pid", ri, ErrUnderflow)
	}
	t.PID = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("dcs", ri, ErrUnderflow)
	}
	t.DCS = DCS(src[ri])
	ri++
	n, err = t.VP.UnmarshalBinary(src[ri:], t.FirstOctet.VPF())
	if err != nil {
		return DecodeError("vp", ri, err)
	}
	ri += n
	err = t.decodeUserData(src[ri:])
	if err != nil {
		return DecodeError("ud", ri, err)
	}
	return nil
}

func (t *TPDU) unmarshalSubmitReport(src []byte) error {
	ri := 0
	if len(src) < 1 {
		return DecodeError("fcs", ri, ErrUnderflow)
	}
	t.FCS = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("pi", ri, ErrUnderflow)
	}
	t.PI = PI(src[ri])
	ri++
	if len(src) < ri+7 {
		return DecodeError("scts", ri, ErrUnderflow)
	}
	err := t.SCTS.UnmarshalBinary(src[ri : ri+7])
	if err != nil {
		return DecodeError("scts", ri, err)
	}
	ri += 7
	if t.PI.PID() {
		if len(src) <= ri {
			return DecodeError("pid", ri, ErrUnderflow)
		}
		t.PID = src[ri]
		ri++
	}
	if t.PI.DCS() {
		if len(src) <= ri {
			return DecodeError("dcs", ri, ErrUnderflow)
		}
		t.DCS = DCS(src[ri])
		ri++
	}
	if t.PI.UDL() {
		err := t.decodeUserData(src[ri:])
		if err != nil {
			return DecodeError("ud", ri, err)
		}
	}
	return nil
}

// decodeUserData unmarshals the User Data field from the binary src.
func (t *TPDU) decodeUserData(src []byte) error {
	if len(src) < 1 {
		return DecodeError("udl", 0, ErrUnderflow)
	}
	udl := int(src[0])
	if udl == 0 {
		return nil
	}
	var udh UserDataHeader
	sml7 := 0
	ri := 1
	alphabet, err := t.Alphabet()
	if err != nil {
		return DecodeError("alphabet", ri, err)
	}
	if alphabet == Alpha7Bit {
		sml7 = udl
		// length is septets - convert to octets
		udl = (sml7*7 + 7) / 8
	}
	if len(src) < ri+udl {
		return DecodeError("sm", ri, ErrUnderflow)
	}
	if len(src) > ri+udl {
		return DecodeError("ud", ri, ErrOverlength)
	}
	var udhl int // Note that in this context udhl includes itself.
	udhi := t.UDHI()
	if udhi {
		udh = make(UserDataHeader, 0)
		l, err := udh.UnmarshalBinary(src[ri:])
		if err != nil {
			return DecodeError("udh", ri, err)
		}
		udhl = l
		ri += udhl
	}
	if ri == len(src) {
		t.UDH = udh
		return nil
	}
	switch alphabet {
	case Alpha7Bit:
		sm, err := decode7Bit(sml7, udhl, src[ri:])
		if err != nil {
			return DecodeError("sm", ri, err)
		}
		t.UD = sm
	case AlphaUCS2:
		if len(src[ri:])&0x01 == 0x01 {
			return DecodeError("sm", ri, ErrOddUCS2Length)
		}
		fallthrough
	case Alpha8Bit:
		t.UD = append([]byte(nil), src[ri:]...)
	}
	t.UDH = udh
	return nil
}

// decode7Bit decodes the GSM7 encoded binary src into a byte array.
//
// sml is the number of septets expected, and udhl is the number of octets in
// the UDH, including the UDHL field.
func decode7Bit(sml, udhl int, src []byte) ([]byte, error) {
	var fillBits int
	if udhl > 0 {
		if dangling := udhl % 7; dangling != 0 {
			fillBits = 7 - dangling
		}
		sml = sml - (udhl*8+fillBits)/7
	}
	sm := gsm7.Unpack7Bit(src, fillBits)
	// this is a double check on the math and should never trip...
	if len(sm) < sml {
		return nil, ErrUnderflow
	}
	if len(sm) > sml {
		if len(sm) > sml+1 || sm[sml] != 0 {
			return nil, ErrOverlength
		}
		// drop trailing 0 septet
		sm = sm[:sml]
	}
	return sm, nil
}

// encodeUserData marshals the User Data into binary.
//
// The User Data Header is also encoded if present.
// If Alphabet is GSM7 then the User Data is assumed to be unpacked GSM7
// septets and is packed prior to encoding.
// For other alphabet values the User Data is encoded as is.
// No checks of encoded size are performed here as that depends on concrete
// TPDU type, and that can check the length of the returned b.
func (t *TPDU) encodeUserData() (b []byte, err error) {
	udh, err := t.UDH.MarshalBinary()
	if err != nil {
		// never trips as UDH marshalling never fails...
		return nil, EncodeError("udh", err)
	}
	ud := t.UD
	alphabet, err := t.Alphabet()
	if err != nil {
		return nil, EncodeError("alphabet", err)
	}
	udl := len(t.UD) // assume octets
	switch alphabet {
	case Alpha7Bit:
		fillBits := 0
		if dangling := len(udh) % 7; dangling != 0 {
			fillBits = 7 - dangling
		}
		ud = gsm7.Pack7Bit(t.UD, fillBits)
		// udl is in septets so convert
		if udl > 0 {
			udl = udl + (len(udh)*8+fillBits)/7
		} else {
			udl = (len(udh) * 8) / 7
		}
	case AlphaUCS2:
		if udl&0x01 == 0x01 {
			return nil, EncodeError("sm", ErrOddUCS2Length)
		}
		fallthrough
	case Alpha8Bit:
		// udl is in octets
		udl = udl + len(udh)
	}
	b = make([]byte, 0, 1+len(udh)+len(ud))
	b = append(b, byte(udl))
	b = append(b, udh...)
	b = append(b, ud...)
	return b, nil
}

// MaxUDL is the maximum number of octets that can be encoded into the UD.
// Note that for 7bit encoding this can result in up to 160 septets.
const MaxUDL = 140

// MessageType identifies the type of TPDU encoded in a binary stream, as
// defined in 3GPP TS 23.040 Section 9.2.3.1.
// Note that the direction of the TPDU must also be known to determine how to
// interpret the TPDU.
type MessageType int

const (
	// MtDeliver identifies the message as a SMS-Deliver or SMS-Deliver-Report
	// TPDU.
	MtDeliver MessageType = iota

	// MtSubmit identifies the message as a SMS-Submit or SMS-Submit-Report
	// TPDU.
	MtSubmit

	// MtCommand identifies the message as a SMS-Command or SMS-Status-Report
	// TPDU.
	MtCommand

	// MtReserved identifies the message as an unknown type of SMS TPDU.
	MtReserved
)

// ApplyTPDUOption sets the TPDU MTI.
func (mti MessageType) ApplyTPDUOption(t *TPDU) error {
	t.FirstOctet = t.FirstOctet.WithMTI(mti)
	return nil
}

// Direction indicates the direction that the SMS TPDU is carried.
type Direction int

const (
	// MT indicates that the SMS TPDU is intended to be received by the MS.
	MT Direction = iota

	// MO indicates that the SMS TPDU is intended to be sent by the MS.
	MO
)

// ApplyTPDUOption sets the direction of the TPDU.
func (d Direction) ApplyTPDUOption(t *TPDU) error {
	t.Direction = d
	return nil
}

// SmsType indicatges the type of SMS TPDU type represented by the TPDU.
type SmsType int

const (
	// SmsDeliver indiates the TPDU represents a SMS-DELIVER
	SmsDeliver SmsType = iota

	// SmsDeliverReport indiates the TPDU represents a SMS-DELIVER-REPORT
	SmsDeliverReport

	// SmsSubmitReport indiates the TPDU represents a SMS-SUBMIT-REPORT
	SmsSubmitReport

	// SmsSubmit indiates the TPDU represents a SMS-SUBMIT
	SmsSubmit

	// SmsStatusReport indiates the TPDU represents a SMS-STATUS-REPORT
	SmsStatusReport

	// SmsCommand indiates the TPDU represents a SMS-COMMAND
	SmsCommand
)

// MTI returns the MessageType corresponding to the SmsType.
func (st SmsType) MTI() MessageType {
	return MessageType(st >> 1)
}

// Direction returns the direction corresponding to the SmsType.
func (st SmsType) Direction() Direction {
	return Direction(st & 0x01)
}

func (st SmsType) String() string {
	switch st {
	case SmsDeliver:
		return "SmsDeliver"
	case SmsDeliverReport:
		return "SmsDeliverReport"
	case SmsSubmitReport:
		return "SmsSubmitReport"
	case SmsSubmit:
		return "SmsSubmit"
	case SmsStatusReport:
		return "SmsStatusReport"
	case SmsCommand:
		return "SmsCommand"
	default:
		return "Unknown"
	}
}

// ApplyTPDUOption sets the TPDU direction and MTI to match the SmsType.
func (st SmsType) ApplyTPDUOption(t *TPDU) error {
	return t.SetSmsType(st)
}

func smsType(mt MessageType, dir Direction) SmsType {
	return SmsType(byte(mt<<1) | byte(dir))
}

func newInfoElement(msgCount, segCount, segment int) InformationElement {
	ie := InformationElement{}
	ie.ID = 0
	ie.Data = []byte{byte(msgCount), byte(segCount), byte(segment)}
	return ie
}

func newInfoElement16bit(msgCount, segCount, segment int) InformationElement {
	ie := InformationElement{}
	ie.ID = 8
	ie.Data = []byte{0, 0, byte(segCount), byte(segment)}
	binary.BigEndian.PutUint16(ie.Data, uint16(msgCount))
	return ie
}

const (
	esc byte = 0x1b
)

// chunk splits a message into chunks that are not larger than bs.
func chunk(msg []byte, alpha Alphabet, bs int) [][]byte {
	switch alpha {
	default: // default to 7Bit
		return chunk7Bit(msg, bs)
	case AlphaUCS2:
		return chunkUCS2(msg, bs)
	case Alpha8Bit:
		return chunk8Bit(msg, bs)
	}
}

// chunk7Bit splits a GSM7 message into chunks that are not larger than bs.
//
// Escaped characters are not split across blocks, so the resulting blocks may
// be one septet shorter than bs.
func chunk7Bit(msg []byte, bs int) [][]byte {
	if len(msg) == 0 {
		return nil
	}
	count := 1 + len(msg)/bs
	chunks := make([][]byte, 0, count)
	bstart := 0
	bend := bs
	for bend < len(msg) {
		// don't split escapes
		if msg[bend-1] == esc && msg[bend-2] != esc {
			bend--
		}
		chunks = append(chunks, msg[bstart:bend])
		bstart = bend
		bend = bstart + bs
	}
	chunks = append(chunks, msg[bstart:])
	return chunks
}

// chunk8Bit splits a raw 8bit message into chunks that are bs, except for the
// last segment which contains any residual bytes.
func chunk8Bit(msg []byte, bs int) [][]byte {
	if len(msg) == 0 {
		return nil
	}
	count := 1 + len(msg)/bs
	chunks := make([][]byte, 0, count)
	bstart := 0
	bend := bs
	for bend < len(msg) {
		chunks = append(chunks, msg[bstart:bend])
		bstart = bend
		bend = bstart + bs
	}
	chunks = append(chunks, msg[bstart:])
	return chunks
}

const (
	surrHighStart = 0xd800
	surrLowStart  = 0xdc00
)

// chunkUCS2 splits a UCS2/UTF-16 message into chunks that are not larger than bs.
//
// bs should be even, but if odd is reduced by one.
// To allow for reassemblers that cannot handle split surrogate pairs, they are
// not split during chunking, so the resulting blocks may be slightly smaller
// than bs whenever a surrogate pair would span a block boundary.
// While the msg should have even length for UCS2, the chunker does not enforce
// this, and if an odd length message is presented then the final chunk will
// have an odd length.
func chunkUCS2(msg []byte, bs int) [][]byte {
	if len(msg) == 0 {
		return nil
	}
	bs &^= 0x1
	// rough count of blocks - may be off due to not splitting surrogates, but
	// not worth working out the precise count in advance.
	count := 1 + len(msg)/bs
	chunks := make([][]byte, 0, count)
	bstart := 0
	bend := bstart + bs
	for bend < len(msg) {
		// check last uint16 is a high surrogate, if so then leave for later
		r := binary.BigEndian.Uint16(msg[bend-2 : bend])
		if surrHighStart <= r && r < surrLowStart {
			bend = bend - 2
		}
		chunks = append(chunks, msg[bstart:bend])
		bstart = bend
		bend = bstart + bs
	}
	chunks = append(chunks, msg[bstart:])
	return chunks
}
