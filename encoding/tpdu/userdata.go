// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"encoding/binary"

	"github.com/warthog618/sms/encoding/gsm7"
	"github.com/warthog618/sms/encoding/gsm7/charset"
	"github.com/warthog618/sms/encoding/ucs2"
)

// UserData represents the User Data field as defined in 3GPP TS 23.040 Section
// 9.2.3.24.
//
// The UserData is comprised of an optional User Data Header and a short
// message field.
type UserData []byte

// UserDataHeader represents the header section of the User Data as defined in
// 3GPP TS 23.040 Section 9.2.3.24.
type UserDataHeader []InformationElement

// InformationElement represents one of the information elements contained in
// the User Data Header.
type InformationElement struct {
	ID   byte
	Data []byte
}

// UDHL returns the encoded length of the UDH, not including the UDHL itself.
func (udh UserDataHeader) UDHL() int {
	udhl := 0
	for _, ie := range udh {
		udhl += (2 + len(ie.Data))
	}
	return udhl
}

// MarshalBinary marshals the User Data Header, including the UDHL, into
// binary.
func (udh UserDataHeader) MarshalBinary() ([]byte, error) {
	if len(udh) == 0 {
		return nil, nil
	}
	udhl := udh.UDHL()
	b := make([]byte, 0, udhl+1)
	b = append(b, byte(udhl))
	for _, ie := range udh {
		b = append(b, ie.ID, byte(len(ie.Data)))
		b = append(b, ie.Data...)
	}
	return b, nil
}

// UnmarshalBinary reads the InformationElements from the binary User Data
// Header.
//
// The src contains the complete UDH, including the UDHL and all IEs.
// The function returns the number of bytes read from src, and any error
// detected while unmarshalling.
func (udh *UserDataHeader) UnmarshalBinary(src []byte) (int, error) {
	if len(src) < 1 {
		return 0, DecodeError("udhl", 0, ErrUnderflow)
	}
	udhl := int(src[0])
	udhl++ // so it includes itself
	ri := 1
	if len(src) < udhl {
		return ri, DecodeError("ie", ri, ErrUnderflow)
	}
	ies := []InformationElement(nil)
	for ri < udhl {
		if udhl < ri+2 {
			return ri, DecodeError("ie", ri, ErrUnderflow)
		}
		var ie InformationElement
		ie.ID = src[ri]
		ri++
		iedl := int(src[ri])
		ri++
		if len(src) < ri+iedl {
			return ri, DecodeError("ied", ri, ErrUnderflow)
		}
		ie.Data = append([]byte(nil), src[ri:ri+iedl]...)
		ri += iedl
		ies = append(ies, ie)
	}
	*udh = ies
	return udhl, nil
}

// IE returns the last instance of the IE with the given id in the UDH.
//
// If no such IE is found then the function returns false.
func (udh UserDataHeader) IE(id byte) (InformationElement, bool) {
	for i := len(udh) - 1; i >= 0; i-- {
		if udh[i].ID == id {
			return udh[i], true
		}
	}
	return InformationElement{}, false
}

// IEs returns all instances of the IEs with the given id in the UDH.
func (udh UserDataHeader) IEs(id byte) []InformationElement {
	ies := []InformationElement(nil)
	for _, ie := range udh {
		if ie.ID == id {
			ies = append(ies, ie)
		}
	}
	return ies
}

// ConcatInfo extracts the segmentation info contained in the provided User
// Data Header.
//
// If the UDH contains no segmentation information then ok is false and zero
// values are returned.
// The returned values do not distinguish between 8bit and 16bit message
// reference numbers.
func (udh UserDataHeader) ConcatInfo() (segments, seqno, mref int, ok bool) {
	if len(udh) == 0 {
		// single segment - most likely case
		return
	}
	if segments, seqno, mref, ok = udh.ConcatInfo8(); ok {
		return
	}
	return udh.ConcatInfo16()
}

// ConcatInfo8 extracts the segmentation info contained in the provided User
// Data Header, for the 8bit message reference case.
//
// If the UDH contains no segmentation information then ok is false and zero
// values are returned.
func (udh UserDataHeader) ConcatInfo8() (segments, seqno, mref int, ok bool) {
	if c, k := udh.IE(0x00); k && len(c.Data) == 3 {
		ok = true
		mref = int(c.Data[0])
		segments = int(c.Data[1])
		seqno = int(c.Data[2])
	}
	return
}

// ConcatInfo16 extracts the segmentation info contained in the provided User
// Data Header, for the 16bit message reference case.
// If the UDH contains no segmentation information then ok is false and zero
// values are returned.
func (udh UserDataHeader) ConcatInfo16() (segments, seqno, mref int, ok bool) {
	if c, k := udh.IE(0x08); k && len(c.Data) == 4 {
		ok = true
		mref = int(binary.BigEndian.Uint16(c.Data[0:2]))
		segments = int(c.Data[2])
		seqno = int(c.Data[3])
	}
	return
}

// UDDecoder converts TPDU UD to UTF8.
// By default the translator only supports the default character set.
// Additional character sets can be added using the AddLockingCharset and
// AddShiftCharset methods.
type UDDecoder struct {
	locking map[int]bool
	shift   map[int]bool
}

// UDDecoderOption is a construction option for UDDecoders.
type UDDecoderOption interface {
	applyDecoderOption(*UDDecoder)
}

// NewUDDecoder creates a new UD Decoder.
func NewUDDecoder(options ...UDDecoderOption) *UDDecoder {
	d := UDDecoder{}
	for _, option := range options {
		option.applyDecoderOption(&d)
	}
	return &d
}

// Decode converts TPDU UD into the corresponding UTF8 message.
// The UD is expected to be unpacked, as stored in TPDU UD.
// If the UD is GSM7 encoded then it is translated to UTF8 with the default
// character set, or with the character set specified in the UDH, assuming the
// corresponding language has been registered with the UDDecoder.
// If the UDH specifies a character set that has not been registered then the
// translation will fall back to the default character set.
func (d *UDDecoder) Decode(ud UserData, udh UserDataHeader, alpha Alphabet) ([]byte, error) {
	switch alpha {
	case AlphaUCS2:
		m, err := ucs2.Decode(ud)
		return []byte(string(m)), err
	case Alpha8Bit:
		return ud, nil
	case Alpha7Bit:
		fallthrough
	default:
		options := []gsm7.DecoderOption{}
		if ie, ok := udh.IE(lockingIEI); ok {
			if len(ie.Data) >= 1 {
				nli := int(ie.Data[0])
				if _, ok := d.locking[nli]; ok {
					options = append(options, gsm7.WithCharset(nli))
				}
			}
		}
		if ie, ok := udh.IE(shiftIEI); ok {
			if len(ie.Data) >= 1 {
				nli := int(ie.Data[0])
				if _, ok := d.shift[nli]; ok {
					options = append(options, gsm7.WithExtCharset(nli))
				}
			}
		}
		return gsm7.Decode(ud, options...)
	}
}

// UDEncoder converts TPDU UD into the corresponding binary UD.
// By default the translator only supports the default character set.
// Additional character sets can be added using the AddLockingCharset and
// AddShiftCharset methods.
type UDEncoder struct {
	l []int // locking charsets in order
	s []int // shift charsets in order
}

// UDEncoderOption is a construction option for UDEncoders.
type UDEncoderOption interface {
	applyEncoderOption(*UDEncoder)
}

// NewUDEncoder creates a new UDEncoder.
func NewUDEncoder(options ...UDEncoderOption) *UDEncoder {
	e := UDEncoder{}
	for _, option := range options {
		option.applyEncoderOption(&e)
	}
	return &e
}

// CharsetOption specifies the locking and shift character sets available for
// encoding and decoding.
//
// These are in addition to the default character set.
type CharsetOption struct {
	nli []int
}

func (o CharsetOption) applyDecoderOption(d *UDDecoder) {
	if d.locking == nil {
		d.locking = make(map[int]bool)
	}
	if d.shift == nil {
		d.shift = make(map[int]bool)
	}
	for _, nli := range o.nli {
		d.locking[nli] = true
		d.shift[nli] = true
	}
}

func (o CharsetOption) applyEncoderOption(e *UDEncoder) {
	e.l = append(e.l, o.nli...)
	e.s = append(e.s, o.nli...)
}

// LockingCharsetOption specifies the locking character sets available for
// encoding and decoding.
//
// These are in addition to the default character set.
type LockingCharsetOption struct {
	nli []int
}

func (o LockingCharsetOption) applyDecoderOption(d *UDDecoder) {
	if d.locking == nil {
		d.locking = make(map[int]bool)
	}
	for _, nli := range o.nli {
		d.locking[nli] = true
	}
}

func (o LockingCharsetOption) applyEncoderOption(e *UDEncoder) {
	e.l = append(e.l, o.nli...)
}

// ShiftCharsetOption specifies the shift character sets available for encoding
// and decoding.
//
// These are in addition to the default character set.
type ShiftCharsetOption struct {
	nli []int
}

func (o ShiftCharsetOption) applyDecoderOption(d *UDDecoder) {
	if d.shift == nil {
		d.shift = make(map[int]bool)
	}
	for _, nli := range o.nli {
		d.shift[nli] = true
	}
}

func (o ShiftCharsetOption) applyEncoderOption(e *UDEncoder) {
	e.s = append(e.s, o.nli...)
}

// AllCharsetsOption specifies that all character sets are available for
// encoding and decoding.
type AllCharsetsOption struct{}

func (o AllCharsetsOption) applyEncoderOption(e *UDEncoder) {
	e.l = make([]int, charset.Size)
	for nli := charset.Start; nli < charset.End; nli++ {
		e.l[nli-1] = nli
	}
	e.s = e.l
}

func (o AllCharsetsOption) applyDecoderOption(d *UDDecoder) {
	if d.locking == nil {
		d.locking = make(map[int]bool)
	}
	if d.shift == nil {
		d.shift = make(map[int]bool)
	}
	for nli := charset.Start; nli < charset.End; nli++ {
		d.locking[nli] = true
		d.shift[nli] = true
	}
}

// WithAllCharsets makes all possible character sets available to Encode.
// This is equivalent to calling AddLockingCharset and AddShiftCharset for all
// possible NationalLanguageIdentifiers, in increasing order.
var WithAllCharsets = AllCharsetsOption{}

// WithCharset adds character sets to the sets available to Encode.
func WithCharset(nli ...int) CharsetOption {
	return CharsetOption{nli}
}

// WithLockingCharset adds a set of locking character sets to the sets available to Encode.
func WithLockingCharset(nli ...int) LockingCharsetOption {
	return LockingCharsetOption{nli}
}

// WithShiftCharset adds a set of shift character sets to the sets available to Encode.
func WithShiftCharset(nli ...int) ShiftCharsetOption {
	return ShiftCharsetOption{nli}
}

const (
	shiftIEI   byte = 24
	lockingIEI byte = 25
)

// Encode converts a UTF8 message into corresponding TPDU User Data.
// Note that the UD size is not limited to the szie available in a single
// TPDU, and so may need to be segmented into several concatenated messages.
// Encode attempts to pick the most compact alphabet for the given message.
// It assumes GSM7 is the most compact, and, if the default character set is
// insufficient, tries combinations of supported language character sets, in
// the order they were added to the UDEncoder.
// It is not optimal as it performs language selection on the whole message,
// rather than determining the best for each segment in turn. (which is totally
// allowed as stated in 3GPP TS 23.040 9.2.3.24.15 + 16)
// But this may be a safer approach - to allow for the decoder being
// non-compliant, and the benefit of per-segment language encoding is minimal.
// In most cases there is no benefit at all.
//
// Failing GSM7 conversion it falls back to UCS2/UTF16.
func (e *UDEncoder) Encode(msg string) (UserData, UserDataHeader, Alphabet, error) {
	enc, err := gsm7.Encode([]byte(msg)) // default charset
	if err == nil {
		return enc, nil, Alpha7Bit, nil
	}
	// try locking tables with default shift
	for _, nli := range e.l {
		enc, err = gsm7.Encode([]byte(msg), gsm7.WithCharset(nli))
		if err == nil {
			return enc, UserDataHeader{
					InformationElement{ID: lockingIEI, Data: []byte{byte(nli)}}},
				Alpha7Bit, nil
		}
	}
	// try default with language shift tables
	for _, nli := range e.s {
		enc, err = gsm7.Encode([]byte(msg), gsm7.WithExtCharset(nli))
		if err == nil {
			return enc, UserDataHeader{
					InformationElement{ID: shiftIEI, Data: []byte{byte(nli)}}},
				Alpha7Bit, nil
		}
	}
	// could also try combos of locking AND shift, but unlikely to help...

	// fallback to ucs-2
	enc = ucs2.Encode([]rune(msg))
	return enc, nil, AlphaUCS2, nil
}
