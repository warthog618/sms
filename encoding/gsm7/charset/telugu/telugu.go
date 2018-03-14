// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package telugu

// Decoder provides a mapping from GSM7 byte to UTF8 rune.
type Decoder map[byte]rune

// Encoder provides a mapping from UTF8 rune to GSM7 byte.
type Encoder map[rune]byte

// NewDecoder returns the mapping table from GSM7 to UTF8.
func NewDecoder() Decoder {
	return dset
}

// NewExtDecoder returns the extension mapping table from GSM7 to UTF8.
func NewExtDecoder() Decoder {
	return dext
}

// NewEncoder returns the mapping table from UTF8 to GSM7.
func NewEncoder() Encoder {
	return eset
}

// NewExtEncoder returns the extention mapping table from UTF8 to GSM7.
func NewExtEncoder() Encoder {
	return eext
}

var (
	dset = Decoder{
		0x00: '\u0c01',
		0x01: '\u0c02',
		0x02: '\u0c03',
		0x03: '\u0c05',
		0x04: '\u0c06',
		0x05: '\u0c07',
		0x06: '\u0c08',
		0x07: '\u0c09',
		0x08: '\u0c0a',
		0x09: '\u0c0b',
		0x0a: '\n',
		0x0b: '\u0c0c',
		0x0d: '\r',
		0x0e: '\u0c0e',
		0x0f: '\u0c0f',
		0x10: '\u0c10',
		0x12: '\u0c12',
		0x13: '\u0c13',
		0x14: '\u0c14',
		0x15: '\u0c15',
		0x16: '\u0c16',
		0x17: '\u0c17',
		0x18: '\u0c18',
		0x19: '\u0c19',
		0x1a: '\u0c1a',
		0x1b: 0x1b,
		0x1c: '\u0c1b',
		0x1d: '\u0c1c',
		0x1e: '\u0c1d',
		0x1f: '\u0c1e',
		0x20: 0x20,
		0x21: '!',
		0x22: '\u0c1f',
		0x23: '\u0c20',
		0x24: '\u0c21',
		0x25: '\u0c22',
		0x26: '\u0c23',
		0x27: '\u0c24',
		0x28: ')',
		0x29: '(',
		0x2a: '\u0c25',
		0x2b: '\u0c26',
		0x2c: ',',
		0x2d: '\u0c27',
		0x2e: '.',
		0x2f: '\u0c28',
		0x30: '0',
		0x31: '1',
		0x32: '2',
		0x33: '3',
		0x34: '4',
		0x35: '5',
		0x36: '6',
		0x37: '7',
		0x38: '8',
		0x39: '9',
		0x3a: ':',
		0x3b: ';',
		0x3d: '\u0c2a',
		0x3e: '\u0c2b',
		0x3f: '?',
		0x40: '\u0c2c',
		0x41: '\u0c2d',
		0x42: '\u0c2e',
		0x43: '\u0c2f',
		0x44: '\u0c30',
		0x45: '\u0c31',
		0x46: '\u0c32',
		0x47: '\u0c33',
		0x49: '\u0c35',
		0x4a: '\u0c36',
		0x4b: '\u0c37',
		0x4c: '\u0c38',
		0x4d: '\u0c39',
		0x4f: '\u0c3d',
		0x50: '\u0c3e',
		0x51: '\u0c3f',
		0x52: '\u0c40',
		0x53: '\u0c41',
		0x54: '\u0c42',
		0x55: '\u0c43',
		0x56: '\u0c44',
		0x58: '\u0c46',
		0x59: '\u0c47',
		0x5a: '\u0c48',
		0x5c: '\u0c4a',
		0x5d: '\u0c4b',
		0x5e: '\u0c4c',
		0x5f: '\u0c4d',
		0x60: '\u0c55',
		0x61: 'a',
		0x62: 'b',
		0x63: 'c',
		0x64: 'd',
		0x65: 'e',
		0x66: 'f',
		0x67: 'g',
		0x68: 'h',
		0x69: 'i',
		0x6a: 'j',
		0x6b: 'k',
		0x6c: 'l',
		0x6d: 'm',
		0x6e: 'n',
		0x6f: 'o',
		0x70: 'p',
		0x71: 'q',
		0x72: 'r',
		0x73: 's',
		0x74: 't',
		0x75: 'u',
		0x76: 'v',
		0x77: 'w',
		0x78: 'x',
		0x79: 'y',
		0x7a: 'z',
		0x7b: '\u0c56',
		0x7c: '\u0c60',
		0x7d: '\u0c61',
		0x7e: '\u0c62',
		0x7f: '\u0c63',
	}
	dext = Decoder{
		0x00: '@',
		0x01: '£',
		0x02: '$',
		0x03: '¥',
		0x04: '¿',
		0x05: '"',
		0x06: '¤',
		0x07: '%',
		0x08: '&',
		0x09: '\'',
		0x0a: '\f',
		0x0b: '*',
		0x0c: '+',
		0x0d: '\r',
		0x0e: '-',
		0x0f: '/',
		0x10: '<',
		0x11: '=',
		0x12: '>',
		0x13: '¡',
		0x14: '^',
		0x15: '¡',
		0x16: '_',
		0x17: '#',
		0x18: '*',
		0x1b: 0x1b,
		0x1c: '\u0c66',
		0x1d: '\u0c67',
		0x1e: '\u0c68',
		0x1f: '\u0c69',
		0x20: '\u0c6a',
		0x21: '\u0c6b',
		0x22: '\u0c6c',
		0x23: '\u0c6d',
		0x24: '\u0c6e',
		0x25: '\u0c6f',
		0x26: '\u0c58',
		0x27: '\u0c59',
		0x28: '{',
		0x29: '}',
		0x2a: '\u0c78',
		0x2b: '\u0c79',
		0x2c: '\u0c7a',
		0x2d: '\u0c7b',
		0x2e: '\u0c7c',
		0x2f: '\\',
		0x30: '\u0c7d',
		0x31: '\u0c7e',
		0x32: '\u0c7f',
		0x3c: '[',
		0x3d: '~',
		0x3e: ']',
		0x40: '|',
		0x41: 'A',
		0x42: 'B',
		0x43: 'C',
		0x44: 'D',
		0x45: 'E',
		0x46: 'F',
		0x47: 'G',
		0x48: 'H',
		0x49: 'I',
		0x4a: 'J',
		0x4b: 'K',
		0x4c: 'L',
		0x4d: 'M',
		0x4e: 'N',
		0x4f: 'O',
		0x50: 'P',
		0x51: 'Q',
		0x52: 'R',
		0x53: 'S',
		0x54: 'T',
		0x55: 'U',
		0x56: 'V',
		0x57: 'W',
		0x58: 'X',
		0x59: 'Y',
		0x5a: 'Z',
		0x65: '€',
	}
	eset Encoder
	eext Encoder
)

func init() {
	eset = make(Encoder, len(dset))
	for k, v := range dset {
		eset[v] = k
	}
	eext = make(Encoder, len(dext))
	for k, v := range dext {
		if ko, ok := eext[v]; !ok || ko > k {
			eext[v] = k
		}
	}
}
