// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package gsm7_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/gsm7"
)

type testPattern struct {
	name    string
	padBits int
	p       []byte
	u       []byte
}

var (
	testPatterns = []testPattern{
		{
			"nil",
			0,
			nil,
			nil,
		},
		{
			"empty",
			0,
			[]byte{},
			[]byte{},
		},
		{
			"empty fill",
			1,
			[]byte{},
			[]byte{},
		},
		{
			"cr",
			0,
			[]byte{13},
			[]byte("\r"),
		},
		{
			"one",
			0,
			[]byte{49},
			[]byte("1"),
		},
		{
			"two",
			0,
			[]byte{48, 25},
			[]byte("02"),
		},
		{
			"message",
			0,
			[]byte{0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			[]byte("message\x00"),
		},
		{
			"seven",
			0,
			[]byte{240, 48, 157, 94, 150, 187, 1},
			[]byte("pattern\x00"), // packed is zero filled
		},
		{
			"eight",
			0,
			[]byte{240, 48, 157, 94, 150, 187, 127},
			[]byte("pattern?"),
		},
		{
			"nine",
			0,
			[]byte{240, 48, 157, 94, 150, 187, 67, 33},
			[]byte("pattern!!"),
		},
		{
			"long",
			0,
			[]byte{
				97, 144, 189, 44, 207, 131, 216, 111, 247, 25, 68, 47, 207,
				233, 32, 120, 152, 78, 47, 203, 221,
			},
			[]byte("a very long test pattern"),
		},
		{
			"filler",
			0,
			[]byte{230, 52, 155, 93, 150, 255, 0},
			[]byte("filler?\x00"),
		},
		{
			"fill1",
			1,
			[]byte{0xfe},
			[]byte{0x7f},
		},
		{
			"fill2",
			2,
			[]byte{0xfc, 1},
			[]byte{0x7f, 0x00},
		},
		{
			"fill3",
			3,
			[]byte{0xf8, 3},
			[]byte{0x7f},
		},
		{
			"fill4",
			4,
			[]byte{0xf0, 7},
			[]byte{0x7f},
		},
		{
			"fill5",
			5,
			[]byte{0xe0, 0xf},
			[]byte{0x7f},
		},
		{
			"fill6",
			6,
			[]byte{0xc0, 0x1f},
			[]byte{0x7f},
		},
	}
	ussdTestPatterns = []testPattern{
		{
			"nil",
			0,
			nil,
			nil,
		},
		{
			"empty",
			0,
			[]byte{},
			[]byte{},
		},
		{
			"cr",
			0,
			[]byte{13},
			[]byte("\r"),
		},
		{
			"one",
			0,
			[]byte{49},
			[]byte("1"),
		},
		{
			"two",
			0,
			[]byte{48, 25},
			[]byte("02"),
		},
		{
			"message0",
			0,
			[]byte{0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			[]byte("message\x00"),
		},
		{
			"message",
			0,
			[]byte{0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x1b},
			[]byte("message"),
		},
		{
			"seven",
			0,
			[]byte{240, 48, 157, 94, 150, 187, 27}, // packed is filled with CR
			[]byte("pattern"),
		},
		{
			"eight",
			0,
			[]byte{240, 48, 157, 94, 150, 187, 127},
			[]byte("pattern?"),
		},
		{
			"nine",
			0,
			[]byte{240, 48, 157, 94, 150, 187, 67, 33},
			[]byte("pattern!!"),
		},
		{
			"long",
			0,
			[]byte{
				97, 144, 189, 44, 207, 131, 216, 111, 247, 25, 68, 47, 207, 233,
				32, 120, 152, 78, 47, 203, 221,
			},
			[]byte("a very long test pattern"),
		},
		{
			"octet cr",
			0,
			[]byte{240, 48, 157, 94, 150, 187, 27, 13},
			[]byte("pattern\r"),
		},
		{
			"filler",
			0,
			[]byte{230, 52, 155, 93, 150, 255, 26},
			[]byte("filler?"),
		},
		{
			"null filler",
			0,
			[]byte{230, 52, 155, 93, 150, 255, 0},
			[]byte("filler?\x00"),
		},
		{
			"fill1",
			1,
			[]byte{0xfe},
			[]byte{0x7f},
		},
		{
			"fill2",
			2,
			[]byte{0xfc, 27}, // packed is filled with CR
			[]byte{0x7f},
		},
		{
			"fill3",
			3,
			[]byte{0xf8, 3},
			[]byte{0x7f},
		},
		{
			"fill4",
			4,
			[]byte{0xf0, 7},
			[]byte{0x7f},
		},
		{
			"fill5",
			5,
			[]byte{0xe0, 0xf},
			[]byte{0x7f},
		},
		{
			"fill6",
			6,
			[]byte{0xc0, 0x1f},
			[]byte{0x7f},
		},
	}
)

func TestUnpack7Bit(t *testing.T) {
	for _, p := range testPatterns {
		f := func(t *testing.T) {
			u := gsm7.Unpack7Bit(p.p, p.padBits)
			assert.Equal(t, p.u, u)
		}
		t.Run(p.name, f)
	}
}

func TestPack7Bit(t *testing.T) {
	for _, p := range testPatterns {
		f := func(t *testing.T) {
			d := gsm7.Pack7Bit(p.u, p.padBits)
			assert.Equal(t, p.p, d)
		}
		t.Run(p.name, f)
	}
}

func TestUnpack7BitUSSD(t *testing.T) {
	for _, p := range ussdTestPatterns {
		f := func(t *testing.T) {
			u := gsm7.Unpack7BitUSSD(p.p, p.padBits)
			assert.Equal(t, p.u, u)
		}
		t.Run(p.name, f)
	}
}

func TestPack7BitUSSD(t *testing.T) {
	for _, p := range ussdTestPatterns {
		f := func(t *testing.T) {
			d := gsm7.Pack7BitUSSD(p.u, p.padBits)
			assert.Equal(t, p.p, d)
		}
		t.Run(p.name, f)
	}
}
