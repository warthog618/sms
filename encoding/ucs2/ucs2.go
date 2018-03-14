package ucs2

import (
	"encoding/binary"
	"errors"
	"fmt"
	"unicode/utf16"
)

// Decode converts an array of UCS2 characters into an array of runes.
// As the UCS2 characters are packed into a byte array, the length of the
// byte array provided must be even.
func Decode(src []byte) ([]rune, error) {
	if len(src)&0x01 == 0x01 {
		return nil, ErrInvalidLength
	}
	l := len(src) / 2
	dst := make([]rune, l)
	ri := 0
	for i := 0; i < l; i++ {
		dst[i] = rune(binary.BigEndian.Uint16(src[ri:]))
		ri += 2
	}
	return dst, nil
}

// Encode converts an array of UCS2 runes into an array of bytes, where pairs of
// bytes (in Big Endian) represent a UCS2 character.
func Encode(src []rune) ([]byte, error) {
	dst := make([]byte, len(src)*2)
	wi := 0
	for _, r := range src {
		if r1, r2 := utf16.EncodeRune(r); r1 != '\ufffd' || r2 != '\ufffd' {
			return dst[:wi], ErrInvalidRune(r)
		}
		binary.BigEndian.PutUint16(dst[wi:], uint16(r))
		wi += 2
	}
	return dst[:wi], nil
}

// ErrInvalidRune indicates a rune cannot be converted to UCS2.
type ErrInvalidRune rune

func (e ErrInvalidRune) Error() string {
	return fmt.Sprintf("ucs2: invalid rune: 0x%04x", uint32(e))
}

var (
	// ErrInvalidLength indicates the binary provided has an invalid (odd) length.
	ErrInvalidLength = errors.New("ucs2: length must be even")
)
