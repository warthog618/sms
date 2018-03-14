// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package punjabi

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
		0x00: '\u0a01',
		0x01: '\u0a02',
		0x02: '\u0a03',
		0x03: '\u0a05',
		0x04: '\u0a06',
		0x05: '\u0a07',
		0x06: '\u0a08',
		0x07: '\u0a09',
		0x08: '\u0a0a',
		0x0a: '\n',
		0x0d: '\r',
		0x0f: '\u0a0f',
		0x10: '\u0a10',
		0x13: '\u0a13',
		0x14: '\u0a14',
		0x15: '\u0a15',
		0x16: '\u0a16',
		0x17: '\u0a17',
		0x18: '\u0a18',
		0x19: '\u0a19',
		0x1a: '\u0a1a',
		0x1b: 0x1b,
		0x1c: '\u0a1b',
		0x1d: '\u0a1c',
		0x1e: '\u0a1d',
		0x1f: '\u0a1e',
		0x20: 0x20,
		0x21: '!',
		0x22: '\u0a1f',
		0x23: '\u0a20',
		0x24: '\u0a21',
		0x25: '\u0a22',
		0x26: '\u0a23',
		0x27: '\u0a24',
		0x28: ')',
		0x29: '(',
		0x2a: '\u0a25',
		0x2b: '\u0a26',
		0x2c: ',',
		0x2d: '\u0a27',
		0x2e: '.',
		0x2f: '\u0a28',
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
		0x3d: '\u0a2a',
		0x3e: '\u0a2b',
		0x3f: '?',
		0x40: '\u0a2c',
		0x41: '\u0a2d',
		0x42: '\u0a2e',
		0x43: '\u0a2f',
		0x44: '\u0a30',
		0x46: '\u0a32',
		0x47: '\u0a33',
		0x49: '\u0a35',
		0x4a: '\u0a36',
		0x4c: '\u0a38',
		0x4d: '\u0a39',
		0x4e: '\u0a3c',
		0x50: '\u0a3e',
		0x51: '\u0a3f',
		0x52: '\u0a40',
		0x53: '\u0a41',
		0x54: '\u0a42',
		0x59: '\u0a47',
		0x5a: '\u0a48',
		0x5d: '\u0a4b',
		0x5e: '\u0a4c',
		0x5f: '\u0a4d',
		0x60: '\u0a51',
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
		0x7b: '\u0a70',
		0x7c: '\u0a71',
		0x7d: '\u0a72',
		0x7e: '\u0a73',
		0x7f: '\u0a74',
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
		0x19: '\u0964',
		0x1a: '\u0965',
		0x1b: 0x1b,
		0x1c: '\u0a66',
		0x1d: '\u0a67',
		0x1e: '\u0a68',
		0x1f: '\u0a69',
		0x20: '\u0a6a',
		0x21: '\u0a6b',
		0x22: '\u0a6c',
		0x23: '\u0a6d',
		0x24: '\u0a6e',
		0x25: '\u0a6f',
		0x26: '\u0a59',
		0x27: '\u0a5a',
		0x28: '{',
		0x29: '}',
		0x2a: '\u0a5b',
		0x2b: '\u0a5c',
		0x2c: '\u0a5e',
		0x2d: '\u0a75',
		0x2f: '\\',
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
