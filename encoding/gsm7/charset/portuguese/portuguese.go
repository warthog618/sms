// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package portuguese

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
		0x05: 'ê',
		0x09: 'ç',
		0x0a: '\f',
		0x0b: 'Ô',
		0x0c: 'ô',
		0x0d: '\n',
		0x0e: 'Á',
		0x0f: 'á',
		0x12: 'Φ',
		0x13: 'Γ',
		0x14: '^',
		0x15: 'Ω',
		0x16: 'Π',
		0x17: 'Ψ',
		0x18: 'Σ',
		0x19: 'Θ',
		0x28: '{',
		0x29: '}',
		0x2f: '\\',
		0x3c: '[',
		0x3d: '~',
		0x3e: ']',
		0x40: '|',
		0x41: 'À',
		0x49: 'Í',
		0x4f: 'Ó',
		0x55: 'Ú',
		0x5b: 'Ã',
		0x5c: 'Õ',
		0x61: 'Â',
		0x65: '€',
		0x69: 'í',
		0x6f: 'ó',
		0x75: 'ú',
		0x7b: 'ã',
		0x7c: 'õ',
		0x7f: 'â',
	}
	eset charset.Encoder
	eext charset.Encoder
)

func init() {
	// the decoder mapping table, in string form.
	b := []rune(
		"@£$¥êéúíóç\nÔô\rÁáΔ_ªÇÀ∞^\\€Ó|\x1bÂâÊÉ !\"#º%&'()*+,-./0123456789:;<=>?" +
			"ÍABCDEFGHIJKLMNOPQRSTUVWXYZÃÕÚÜ§~abcdefghijklmnopqrstuvwxyzãõ`üà")
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
