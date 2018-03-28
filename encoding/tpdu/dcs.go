// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// DCS represents the SMS Data Coding Scheme field as defined in 3GPP TS 23.040 Section 4.
type DCS byte

// Alphabet defines the encoding of the SMS User Data, as defined in 3GPP TS 23.038 Section 4.
type Alphabet int

const (
	// Alpha7Bit indicates that the UD is encoded using GSM 7 bit encoding.
	// The character set used for the decoding is determined from the UDH.
	Alpha7Bit Alphabet = iota
	// Alpha8Bit indicates that the UD is encoded as raw 8bit data.
	Alpha8Bit
	// AlphaUCS2 indicates that the UD is encoded as UCS-2 (16bit) characters.
	AlphaUCS2
	// AlphaReserved indicates the alphabet is not defined.
	AlphaReserved
)

// Alphabet returns the alphabet used to encode the User Data according to the DCS.
// The DCS is assumed to be defined as per 3GPP TS 23.038 Section 4.
func (d DCS) Alphabet() (Alphabet, error) {
	alpha := Alpha7Bit
	switch {
	case d&0x80 == 0x00: // 0xxx
		alpha = Alphabet((d >> 2) & 0x3)
		if alpha == AlphaReserved {
			alpha = Alpha7Bit
		}
	case d&0xe0 == 0xc0: // 110x
	// is 7bit
	case d&0xf0 == 0xe0: // 1110
		alpha = AlphaUCS2
	case d&0xf0 == 0xf0: // 1111
		if d&0x04 == 0x04 {
			alpha = Alpha8Bit
		} // else 7bit
	default: // includes 10xx reserved coding groups
		return Alpha7Bit, ErrInvalid
	}
	return alpha, nil
}

// WithAlphabet sets the Alphabet bits of the DCS, given the state of the other
// bits.  An error is returned if the state is incompatible with setting the
// alphabet.
func (d DCS) WithAlphabet(a Alphabet) (DCS, error) {
	switch {
	case d&0x80 == 0x00: // 0xxx
		return d&^0x0c | (DCS(a) << 2), nil
	case d&0xe0 == 0xc0 && a == Alpha7Bit: // 110x is 7Bit
		return d, nil
	case d&0xf0 == 0xe0 && a == AlphaUCS2: // 1110 is UCS2
		return d, nil
	case d&0xf0 == 0xf0 && a <= Alpha8Bit:
		return d&^0x0c | (DCS(a) << 2), nil
	default: // includes 110x, 1110, and 10xx reserved coding groups
		return d, ErrInvalid
	}
}

// MessageClass indicates the
type MessageClass int

const (
	// MClass0 is a flash message which is not to be stored in memory.
	MClass0 MessageClass = iota
	// MClass1 is an ME specific message.
	MClass1
	// MClass2 is a SIM/USIM specific message.
	MClass2
	// MClass3 is a TE specific message.
	MClass3
	// MClassUnknown indicates no message class is set.
	MClassUnknown
)

// Class returns the MessageClass indicated by the DCS.
// The DCS is assumed to be defined as per 3GPP TS 23.038 Section 4.
func (d DCS) Class() (MessageClass, error) {
	switch {
	case d&0x90 == 0x10, d&0xf0 == 0xf0: // 0xx1 and 1111
		return MessageClass(d & 0x3), nil
	case d&0xe0 == 0xc0, d&0xf0 == 0xe0: // 110x and 1110
		return MClassUnknown, nil
	default: // includes 10xx reserved coding groups
		return MClassUnknown, ErrInvalid
	}
}

// WithClass sets the MessageClass bits of the DCS, given the state of the other
// bits.  An error is returned if the state is incompatible with setting the
// message class.
func (d DCS) WithClass(c MessageClass) (DCS, error) {
	switch {
	case d&0x80 == 0x00: // 0xxx
		return (d&^0x03 | 0x10 | DCS(c)), nil
	case d&0xf0 == 0xf0: // 1111
		return (d&^0x03 | DCS(c)), nil
	default: // includes 10xx reserved coding groups
		return d, ErrInvalid
	}
}

// Compressed indicates whether the text is compressed using the algorithm defined
// in 3GPP TS 23.024, as determined from the DCS.
// The DCS is assumed to be defined as per 3GPP TS 23.038 Section 4.
func (d DCS) Compressed() bool {
	// only true for 0x1xxxxx (binary)
	return (d&0xa0 == 0x20)
}
