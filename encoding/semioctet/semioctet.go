// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package semioctet

import (
	"errors"
	"fmt"
)

// Mapping from half-octet value to UTF-8 digit.
// The trailing 'F' is never returned but is provided to detect fill.
var decodeDigits = "0123456789*#abcF"

// Decode converts a semi-octet encoded field into the corresponding
// array of UTF-8 digits, as per 3GPP TS 23.040 Section 9.1.2.3
// Conversion is terminated by the length of the src or dst.
// The length of the dst indicates the maximum number of digits to return.
// The return values are the decoded field (a slice of dst),
// the number of bytes read from src, and any error detected during the
// conversion.
func Decode(dst, src []byte) ([]byte, int, error) {
	wi := 0
	ri := 0
	for wi < len(dst) && ri < len(src) {
		d := decodeDigits[src[ri]&0x0f]
		if d != 'F' {
			dst[wi] = d
			wi++
		}
		d = decodeDigits[src[ri]>>4]
		ri++
		if wi == len(dst) && d != 'F' {
			return nil, ri, ErrMissingFill
		}
		if d == 'F' {
			continue
		}
		dst[wi] = d
		wi++
	}
	return dst[:wi], ri, nil
}

var encodeDigits = map[byte]int{
	'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7,
	'8': 8, '9': 9, '*': 10, '#': 11, 'a': 12, 'b': 13, 'c': 14,
}

// Encode converts an array of UTF-8 digits into semi-octet format,
// as per 3GPP TS 23.040 Section 9.1.2.3
// The return values are the encoded field, and any error detected during the
// conversion.
func Encode(src []byte) ([]byte, error) {
	b := make([]byte, (len(src)+1)/2)
	hi := false
	p := 0
	wi := 0
	for _, d := range src {
		e, ok := encodeDigits[d]
		if !ok {
			return nil, ErrInvalidDigit(d)
		}
		if hi {
			p = p | int(e<<4)
			b[wi] = byte(p)
			wi++
			hi = false
		} else {
			p = e
			hi = true
		}
	}
	if hi == true {
		p = p | 0xf0
		b[wi] = byte(p)
	}
	return b, nil
}

// ErrInvalidDigit indicates that the digit can not be encoded into semioctet format.
type ErrInvalidDigit byte

func (e ErrInvalidDigit) Error() string {
	return fmt.Sprintf("semioctet: invalid digit: '%c' - 0x%x", byte(e), int(e))
}

var (
	// ErrMissingFill indicates the final src octet does not contain the expected
	// fill character 'F'.
	ErrMissingFill = errors.New("semioctet: last src octet missing expected fill")
)
