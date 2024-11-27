// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package charset

var (
	portugueseDecoder    Decoder
	portugueseExtDecoder = Decoder{
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
		0x1f: 'Ê',
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
	portugueseEncoder    Encoder
	portugueseExtEncoder Encoder
	portugueseRunes      = []rune(
		"@£$¥êéúíóç\nÔô\rÁáΔ_ªÇÀ∞^\\€Ó|\x1bÂâÊÉ !\"#º%&'()*+,-./0123456789:;<=>?" +
			"ÍABCDEFGHIJKLMNOPQRSTUVWXYZÃÕÚÜ§~abcdefghijklmnopqrstuvwxyzãõ`üà")
)

func generatePortugueseEncoder() Encoder {
	return generateEncoderFromRunes(portugueseRunes)
}

func generatePortugueseDecoder() Decoder {
	return generateDecoderFromRunes(portugueseRunes)
}

func generatePortugueseExtEncoder() Encoder {
	return generateEncoder(portugueseExtDecoder)
}
