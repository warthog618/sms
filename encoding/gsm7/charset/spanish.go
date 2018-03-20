// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package charset

var (
	spanishExtDecoder = Decoder{
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
	spanishExtEncoder Encoder
)

func generateSpanishExtEncoder() Encoder {
	return generateEncoder(spanishExtDecoder)
}
