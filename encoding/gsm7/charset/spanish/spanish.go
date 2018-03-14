// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package spanish

// Decoder provides a mapping from GSM7 byte to UTF8 rune.
type Decoder map[byte]rune

// Encoder provides a mapping from UTF8 rune to GSM7 byte.
type Encoder map[rune]byte

// NewExtDecoder returns the extension mapping table from GSM7 to UTF8.
func NewExtDecoder() Decoder {
	return dext
}

// NewExtEncoder returns the extention mapping table from UTF8 to GSM7.
func NewExtEncoder() Encoder {
	return eext
}

var (
	dext = Decoder{
		0x09: 'ç',
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
		0x41: 'Á',
		0x49: 'Í',
		0x4f: 'Ó',
		0x55: 'Ú',
		0x61: 'á',
		0x65: '€',
		0x69: 'í',
		0x6f: 'ó',
		0x75: 'ú',
	}
	eext Encoder
)

func init() {
	eext = make(Encoder, len(dext))
	for k, v := range dext {
		eext[v] = k
	}
}
