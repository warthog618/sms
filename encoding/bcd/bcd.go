// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package bcd

import (
	"fmt"
)

// Decode decodes a BCD encoded octet into the equivalent integer.
// The lowest nibble is taken as being the most significant.
func Decode(bcd byte) (int, error) {
	msn := bcd & 0x0f
	lsn := bcd >> 4
	if msn > 9 || lsn > 9 {
		return 0, ErrInvalidOctet(bcd)
	}
	return int(msn*10 + lsn), nil
}

// DecodeSigned decodes a BCD encoded octet where bit 3 of the msn indicates
// the sign of the encoded integer.
func DecodeSigned(bcd byte) (int, error) {
	msn := bcd & 0x07
	lsn := bcd >> 4
	if lsn > 9 {
		return 0, ErrInvalidOctet(bcd)
	}
	retval := int(msn*10 + lsn)
	if bcd&0x08 != 0 {
		retval = -retval
	}
	return retval, nil
}

// Encode converts an integer in the range 0..99 into two BCD digits.
// The return value is the two BCD digits encoded into a byte in big endian,
// and any error detected during conversion.
func Encode(u int) (byte, error) {
	if u < 0 || u > 99 {
		return 0, ErrInvalidInteger(u)
	}
	msn := u % 10
	lsn := u / 10
	b := (msn << 4) | lsn
	return byte(b), nil
}

// EncodeSigned converts an integer in the range -79..79 into two BCD digits.
// The return value is the two BCD digits encoded into a byte, with the most
// significant digit stored in the lowest nibble, and any error detected
// during conversion.  If the integer is negative then bit 3 of the byte is
// set to 1 .
func EncodeSigned(s int) (byte, error) {
	if s < -79 || s > 79 {
		return 0, ErrInvalidInteger(s)
	}
	b := 0
	if s < 0 {
		b = 0x08
		s = -s
	}
	msn := s % 10
	lsn := s / 10
	b = b | (msn << 4) | lsn
	return byte(b), nil
}

// ErrInvalidOctet indicates that at least one of the nibbles in the BCD octet
// is invalid, i.e. greater than 9.
// For DecodeSigned only the upper (least significant) nibble can be invalid.
type ErrInvalidOctet byte

func (e ErrInvalidOctet) Error() string {
	return fmt.Sprintf("bcd: invalid octet: 0x%02x", byte(e))
}

// ErrInvalidInteger indicates that the integer is outside the range that can
// be encoded.
type ErrInvalidInteger int

func (e ErrInvalidInteger) Error() string {
	return fmt.Sprintf("bcd: invalid integer: %d", int(e))
}
