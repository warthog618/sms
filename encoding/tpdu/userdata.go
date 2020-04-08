// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

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

func (ie InformationElement) marshalledLen() int {
	return 2 + len(ie.Data)
}

// UDHL returns the encoded length of the UDH, not including the UDHL itself.
func (udh UserDataHeader) UDHL() int {
	udhl := 0
	for _, ie := range udh {
		udhl += ie.marshalledLen()
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
		return 0, NewDecodeError("udhl", 0, ErrUnderflow)
	}
	udhl := int(src[0])
	udhl++ // so it includes itself
	ri := 1
	if len(src) < udhl {
		return ri, NewDecodeError("ie", ri, ErrUnderflow)
	}
	ies := []InformationElement(nil)
	for ri < udhl {
		if udhl < ri+2 {
			return ri, NewDecodeError("ie", ri, ErrUnderflow)
		}
		var ie InformationElement
		ie.ID = src[ri]
		ri++
		iedl := int(src[ri])
		ri++
		if len(src) < ri+iedl {
			return ri, NewDecodeError("ied", ri, ErrUnderflow)
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
//
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

type udDecodeConfig struct {
	locking map[int]bool
	shift   map[int]bool
}

// UDDecodeOption provides behavioural modifiers for DecodeUserData,
// specifically the character sets available to decode GSM7.
type UDDecodeOption interface {
	applyDecodeOption(udDecodeConfig) udDecodeConfig
}

// DecodeUserData converts TPDU UD into the corresponding UTF8 message.
//
// The UD is expected to be unpacked, as stored in TPDU UD. If the UD is GSM7
// encoded then it is translated to UTF8 with the default character set, or
// with the character set specified in the UDH, assuming the corresponding
// language has been registered with the UDDecoder. If the UDH specifies a
// character set that has not been registered then the translation will fall
// back to the default character set.
func DecodeUserData(ud UserData, udh UserDataHeader, alpha Alphabet, options ...UDDecodeOption) ([]byte, error) {
	switch alpha {
	case AlphaUCS2:
		m, err := ucs2.Decode(ud)
		return []byte(string(m)), err
	case Alpha8Bit:
		return ud, nil
	case Alpha7Bit:
		fallthrough
	default:
		cfg := udDecodeConfig{locking: map[int]bool{}, shift: map[int]bool{}}
		for _, option := range options {
			cfg = option.applyDecodeOption(cfg)
		}
		options := []gsm7.DecoderOption{}
		if ie, ok := udh.IE(lockingIEI); ok {
			if len(ie.Data) >= 1 {
				nli := int(ie.Data[0])
				if _, ok := cfg.locking[nli]; ok {
					options = append(options, gsm7.WithCharset(nli))
				}
			}
		}
		if ie, ok := udh.IE(shiftIEI); ok {
			if len(ie.Data) >= 1 {
				nli := int(ie.Data[0])
				if _, ok := cfg.shift[nli]; ok {
					options = append(options, gsm7.WithExtCharset(nli))
				}
			}
		}
		return gsm7.Decode(ud, options...)
	}
}

type udEncodeConfig struct {
	locking []int // locking charsets in order
	shift   []int // shift charsets in order
}

// UDEncodeOption provides behavioural modifiers for EncodeUserData,
// specifically the locking and shift character sets available, in addition to
// the default character set.
type UDEncodeOption interface {
	applyEncodeOption(udEncodeConfig) udEncodeConfig
}

// CharsetOption adds the locking and shift character sets available for
// encoding and decoding.
//
// These are in addition to the default character set.
type CharsetOption struct {
	nli []int
}

func (o CharsetOption) applyDecodeOption(d udDecodeConfig) udDecodeConfig {
	for _, n := range o.nli {
		d.locking[n] = true
		d.shift[n] = true
	}
	return d
}

func (o CharsetOption) applyEncodeOption(e udEncodeConfig) udEncodeConfig {
	e.locking = append(e.locking, o.nli...)
	e.shift = append(e.shift, o.nli...)
	return e
}

// LockingCharsetOption adds to the locking character sets available for
// encoding and decoding.
//
// These are in addition to the default character set.
type LockingCharsetOption struct {
	nli []int
}

func (o LockingCharsetOption) applyDecodeOption(d udDecodeConfig) udDecodeConfig {
	for _, n := range o.nli {
		d.locking[n] = true
	}
	return d
}

func (o LockingCharsetOption) applyEncodeOption(e udEncodeConfig) udEncodeConfig {
	e.locking = append(e.locking, o.nli...)
	return e
}

// ShiftCharsetOption adds the shift character sets available for encoding
// and decoding.
//
// These are in addition to the default character set.
type ShiftCharsetOption struct {
	nli []int
}

func (o ShiftCharsetOption) applyDecodeOption(d udDecodeConfig) udDecodeConfig {
	for _, n := range o.nli {
		d.shift[n] = true
	}
	return d
}

func (o ShiftCharsetOption) applyEncodeOption(e udEncodeConfig) udEncodeConfig {
	e.shift = append(e.shift, o.nli...)
	return e
}

// AllCharsetsOption specifies that all character sets are available for
// encoding and decoding.
type AllCharsetsOption struct{}

func (o AllCharsetsOption) applyDecodeOption(d udDecodeConfig) udDecodeConfig {
	for nli := charset.Start; nli < charset.End; nli++ {
		d.locking[nli] = true
		d.shift[nli] = true
	}
	return d
}

func (o AllCharsetsOption) applyEncodeOption(e udEncodeConfig) udEncodeConfig {
	e.locking = make([]int, charset.Size)
	for nli := charset.Start; nli < charset.End; nli++ {
		e.locking[nli-1] = nli
	}
	e.shift = e.locking
	return e
}

// WithAllCharsets makes all possible character sets available to encode or
// decode.
//
// This is equivalent to calling WithCharset with all possible
// NationalLanguageIdentifiers, in increasing order.
var WithAllCharsets = AllCharsetsOption{}

// WithCharset sets the set of character sets available to encode or decode.
//
// These are in addition to the default character set.
func WithCharset(nli ...int) CharsetOption {
	return CharsetOption{nli}
}

// WithLockingCharset sets the set of locking character sets available to
// encode or decode.
//
// These are in addition to the default character set.
func WithLockingCharset(nli ...int) LockingCharsetOption {
	return LockingCharsetOption{nli}
}

// WithShiftCharset sets the set of shift character sets available to
// encode or decode.
//
// These are in addition to the default character set.
func WithShiftCharset(nli ...int) ShiftCharsetOption {
	return ShiftCharsetOption{nli}
}

const (
	shiftIEI   byte = 24
	lockingIEI byte = 25
)

// EncodeUserData converts a UTF8 message into corresponding TPDU User Data.
//
// Note that the UD size is not limited to the size available in a single TPDU,
// and so may need to be segmented into several concatenated messages. Encode
// attempts to pick the most compact alphabet for the given message. It assumes
// GSM7 is the most compact, and, if the default character set is insufficient,
// tries combinations of supported language character sets, in the order they
// were added to the UDEncoder.
//
// This is not optimal as it performs language selection on the whole message,
// rather than determining the best for each segment in turn. (which is totally
// allowed as stated in 3GPP TS 23.040 9.2.3.24.15 + 16), but this may be a
// safer approach - to allow for the decoder being non-compliant, and the
// benefit of per-segment language encoding is minimal. In most cases there is
// no benefit at all.
//
// Failing GSM7 conversion it falls back to UCS2/UTF16.
func EncodeUserData(msg []byte, options ...UDEncodeOption) (UserData, UserDataHeader, Alphabet) {
	enc, err := gsm7.Encode([]byte(msg)) // default charset
	if err == nil {
		return enc, nil, Alpha7Bit
	}
	cfg := udEncodeConfig{}
	for _, option := range options {
		cfg = option.applyEncodeOption(cfg)
	}
	// try locking tables with default shift
	for _, nli := range cfg.locking {
		enc, err = gsm7.Encode(msg, gsm7.WithCharset(nli))
		if err == nil {
			return enc, UserDataHeader{
					InformationElement{ID: lockingIEI, Data: []byte{byte(nli)}},
				},
				Alpha7Bit
		}
	}
	// try default with language shift tables
	for _, nli := range cfg.shift {
		enc, err = gsm7.Encode(msg, gsm7.WithExtCharset(nli))
		if err == nil {
			return enc, UserDataHeader{
					InformationElement{ID: shiftIEI, Data: []byte{byte(nli)}},
				},
				Alpha7Bit
		}
	}
	// try combination of locking AND shift for same charset
	for _, nli := range cfg.locking {
		for _, snli := range cfg.shift {
			if nli != snli {
				continue
			}
			enc, err = gsm7.Encode(msg, gsm7.WithCharset(nli), gsm7.WithExtCharset(nli))
			if err == nil {
				return enc, UserDataHeader{
						InformationElement{ID: lockingIEI, Data: []byte{byte(nli)}},
						InformationElement{ID: shiftIEI, Data: []byte{byte(nli)}},
					},
					Alpha7Bit
			}
		}
	}
	// could also try other combos of locking AND shift, but unlikely to help??...

	// fallback to ucs-2
	enc = ucs2.Encode([]rune(string(msg)))
	return enc, nil, AlphaUCS2
}
