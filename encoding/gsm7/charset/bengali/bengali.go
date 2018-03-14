// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package bengali

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
		0x00: '\u0981',
		0x01: '\u0982',
		0x02: '\u0983',
		0x03: '\u0985',
		0x04: '\u0986',
		0x05: '\u0987',
		0x06: '\u0988',
		0x07: '\u0989',
		0x08: '\u098a',
		0x09: '\u098b',
		0x0a: '\n',
		0x0b: '\u098c',
		0x0d: '\r',
		0x0f: '\u098f',
		0x10: '\u0990',
		0x13: '\u0993',
		0x14: '\u0994',
		0x15: '\u0995',
		0x16: '\u0996',
		0x17: '\u0997',
		0x18: '\u0998',
		0x19: '\u0999',
		0x1a: '\u099a',
		0x1b: 0x1b,
		0x1c: '\u099b',
		0x1d: '\u099c',
		0x1e: '\u099d',
		0x1f: '\u099e',
		0x20: 0x20,
		0x21: '!',
		0x22: '\u099f',
		0x23: '\u09a0',
		0x24: '\u09a1',
		0x25: '\u09a2',
		0x26: '\u09a3',
		0x27: '\u09a4',
		0x28: ')',
		0x29: '(',
		0x2a: '\u09a5',
		0x2b: '\u09a6',
		0x2c: ',',
		0x2d: '\u09a7',
		0x2e: '.',
		0x2f: '\u09a8',
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
		0x3d: '\u09aa',
		0x3e: '\u09ab',
		0x3f: '?',
		0x40: '\u09ac',
		0x41: '\u09ad',
		0x42: '\u09ae',
		0x43: '\u09af',
		0x44: '\u09b0',
		0x46: '\u09b2',
		0x4a: '\u09b6',
		0x4b: '\u09b7',
		0x4c: '\u09b8',
		0x4d: '\u09b9',
		0x4e: '\u09bc',
		0x4f: '\u09bd',
		0x50: '\u09be',
		0x51: '\u09bf',
		0x52: '\u09c0',
		0x53: '\u09c1',
		0x54: '\u09c2',
		0x55: '\u09c3',
		0x56: '\u09c4',
		0x59: '\u09c7',
		0x5a: '\u09c8',
		0x5d: '\u09cb',
		0x5e: '\u09cc',
		0x5f: '\u09cd',
		0x60: '\u09ce',
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
		0x7b: '\u09d7',
		0x7c: '\u09dc',
		0x7d: '\u09dd',
		0x7e: '\u09f0',
		0x7f: '\u09f1',
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
		0x19: '\u09e6',
		0x1a: '\u09e7',
		0x1b: 0x1b,
		0x1c: '\u09e8',
		0x1d: '\u09e9',
		0x1e: '\u09ea',
		0x1f: '\u09eb',
		0x20: '\u09ec',
		0x21: '\u09ed',
		0x22: '\u09ee',
		0x23: '\u09ef',
		0x24: '\u09df',
		0x25: '\u09e0',
		0x26: '\u09e1',
		0x27: '\u09e2',
		0x28: '{',
		0x29: '}',
		0x2a: '\u09e3',
		0x2b: '\u09f2',
		0x2c: '\u09f3',
		0x2d: '\u09f4',
		0x2e: '\u09f5',
		0x2f: '\\',
		0x30: '\u09f6',
		0x31: '\u09f7',
		0x32: '\u09f8',
		0x33: '\u09f9',
		0x34: '\u09fa',
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
