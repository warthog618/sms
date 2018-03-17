// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package gsm7

import (
	"fmt"

	"github.com/warthog618/sms/encoding/gsm7/charset"
	"github.com/warthog618/sms/encoding/gsm7/charset/basic"
)

const (
	esc byte = 0x1b
	sp  byte = 0x20
)

// Decoder converts from GSM7 to UTF-8 using a particular character set.
type Decoder struct {
	set    charset.Decoder
	ext    charset.Decoder
	strict bool
}

// Encoder converts from UTF-8 to GSM7 using a particular character set.
type Encoder struct {
	set charset.Encoder
	ext charset.Encoder
}

// NewDecoder returns a new GSM7 decoder which uses the basic (default) character set.
func NewDecoder() Decoder {
	return Decoder{basic.NewDecoder(), basic.NewExtDecoder(), false}
}

// NewEncoder returns a new GSM7 encoder which uses the basic (default) character set.
func NewEncoder() Encoder {
	return Encoder{basic.NewEncoder(), basic.NewExtEncoder()}
}

// Decode converts the src from unpacked GSM7 to UTF-8.
func (d *Decoder) Decode(src []byte) ([]byte, error) {
	dst := make([]byte, 0, len(src))
	escaped := false
	for _, g := range src {
		if escaped { // must be first to deal with double escapes
			escaped = false
			if g == esc {
				dst = append(dst, sp)
				continue
			}
			if m, ok := d.ext[g]; ok {
				dst = append(dst, []byte(string(m))...)
				continue
			}
			if d.strict {
				return nil, ErrInvalidSeptet(g)
			}
		} else if g == esc { // then regular escapes
			escaped = true
			continue
		}
		if m, ok := d.set[g]; ok {
			dst = append(dst, []byte(string(m))...)
			continue
		}
		if d.strict {
			return nil, ErrInvalidSeptet(g)
		}
		dst = append(dst, sp)
	}
	// handle dangling escape
	if escaped {
		dst = append(dst, sp)
	}
	return dst, nil
}

// WithCharset replaces the character set map used by the Decoder.
func (d Decoder) WithCharset(set charset.Decoder) Decoder {
	d.set = set
	return d
}

// WithExtCharset replaces the extension character set map used by the Decoder.
func (d Decoder) WithExtCharset(ext charset.Decoder) Decoder {
	d.ext = ext
	return d
}

// Strict makes the Decoder return an error if an unknown character is detected
// when looking up a septet in the character set (not the extension set).
func (d Decoder) Strict() Decoder {
	d.strict = true
	return d
}

// Encode converts the src from UTF-8 to GSM7 and writes the result to dst.
// The return value includes the encoded GSM7 bytes, and any error that
// occured during encoding.
func (e *Encoder) Encode(src []byte) ([]byte, error) {
	dst := make([]byte, 0, len(src))
	for _, u := range string(src) {
		g, ok := e.set[u]
		if ok {
			dst = append(dst, g)
			continue
		}
		g, ok = e.ext[u]
		if ok {
			dst = append(dst, esc, g)
			continue
		}
		return nil, ErrInvalidUTF8(u)
	}
	return dst, nil
}

// WithCharset replaces the character set map used by the Encoder.
func (e Encoder) WithCharset(set charset.Encoder) Encoder {
	e.set = set
	return e
}

// WithExtCharset replaces the extension character set map used by the Encoder.
func (e Encoder) WithExtCharset(ext charset.Encoder) Encoder {
	e.ext = ext
	return e
}

type ErrInvalidSeptet byte

func (e ErrInvalidSeptet) Error() string {
	return fmt.Sprintf("gsm7: invalid septet 0x%02x", int(e))
}

type ErrInvalidUTF8 rune

func (e ErrInvalidUTF8) Error() string {
	return fmt.Sprintf("gsm7: invalid utf8 '%c' (%U)", rune(e), int(e))
}
