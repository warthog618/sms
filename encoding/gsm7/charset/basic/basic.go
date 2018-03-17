// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package basic

import "github.com/warthog618/sms/encoding/gsm7/charset"

// NewDecoder returns the mapping table from GSM7 to UTF8.
func NewDecoder() charset.Decoder {
	return dset
}

// NewExtDecoder returns the extension mapping table from GSM7 to UTF8.
func NewExtDecoder() charset.Decoder {
	return dext
}

// NewEncoder returns the mapping table from UTF8 to GSM7.
func NewEncoder() charset.Encoder {
	return eset
}

// NewExtEncoder returns the extention mapping table from UTF8 to GSM7.
func NewExtEncoder() charset.Encoder {
	return eext
}

var (
	dset charset.Decoder
	dext = charset.Decoder{
		0x0a: '\f',
		0x0d: '\n',
		0x14: '^',
		0x28: '{',
		0x29: '}',
		0x2f: '\\',
		0x3c: '[',
		0x3d: '~',
		0x3e: ']',
		0x40: '|',
		0x65: '€',
	}
	eset charset.Encoder
	eext charset.Encoder
)

func init() {
	// the decoder mapping table, in string form.
	b := []rune(
		"@£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞ\x1bÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?" +
			"¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà")
	dset = make(charset.Decoder, len(b))
	eset = make(charset.Encoder, len(b))
	for i, r := range b {
		dset[byte(i)] = r
		eset[r] = byte(i)
	}
	eext = make(charset.Encoder, len(dext))
	for k, v := range dext {
		eext[v] = k
	}
}
