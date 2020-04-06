// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package charset

var (
	defaultDecoder    Decoder
	defaultExtDecoder = Decoder{
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
	defaultEncoder    Encoder
	defaultExtEncoder Encoder
	defaultRunes      = []rune(
		"@£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞ\x1bÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?" +
			"¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà")
)

func generateDefaultEncoder() Encoder {
	return generateEncoderFromRunes(defaultRunes)
}

func generateDefaultDecoder() Decoder {
	return generateDecoderFromRunes(defaultRunes)
}

func generateDefaultExtEncoder() Encoder {
	return generateEncoder(defaultExtDecoder)
}
