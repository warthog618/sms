// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warthog618/sms/encoding/tpdu"
)

type dcsAlphabetPattern struct {
	in  byte
	out tpdu.Alphabet
	err error
}

func TestDCSAlphabet(t *testing.T) {
	// might as well test them all..
	patterns := []dcsAlphabetPattern{}
	for i := 0; i < 8; i++ {
		m := byte(i << 4)
		patterns = append(patterns,
			dcsAlphabetPattern{0x00 | m, tpdu.Alpha7Bit, nil},
			dcsAlphabetPattern{0x04 | m, tpdu.Alpha8Bit, nil},
			dcsAlphabetPattern{0x08 | m, tpdu.AlphaUCS2, nil},
			dcsAlphabetPattern{0x0c | m, tpdu.Alpha7Bit, nil},
		)
	}
	for i := 0x80; i < 0xc0; i++ {
		patterns = append(patterns,
			dcsAlphabetPattern{byte(i), tpdu.Alpha7Bit, tpdu.ErrInvalid},
		)
	}
	for i := 0xc0; i < 0xe0; i++ {
		patterns = append(patterns,
			dcsAlphabetPattern{byte(i), tpdu.Alpha7Bit, nil},
		)
	}
	for i := 0xe0; i < 0xf0; i++ {
		patterns = append(patterns,
			dcsAlphabetPattern{byte(i), tpdu.AlphaUCS2, nil},
		)
	}
	for i := 0xf0; i <= 0xff; i++ {
		a := tpdu.Alpha7Bit
		if i&0x04 == 0x04 {
			a = tpdu.Alpha8Bit
		}
		patterns = append(patterns,
			dcsAlphabetPattern{byte(i), a, nil},
		)
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			d := tpdu.DCS(p.in)
			a, err := d.Alphabet()
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, a)
		}
		t.Run(fmt.Sprintf("%08b", p.in), f)
	}
}

func TestApplyTPDUOption(t *testing.T) {
	s, err := tpdu.New(tpdu.DCS(0x34))
	require.Nil(t, err)
	assert.Equal(t, tpdu.DCS(0x34), s.DCS)
}

type dcsWithAlphabetIn struct {
	in byte
	a  tpdu.Alphabet
}

type dcsWithAlphabetOut struct {
	out tpdu.DCS
	err error
}

func TestDCSWithAlphabet(t *testing.T) {
	// might as well test them all..
	patterns := make(map[dcsWithAlphabetIn]dcsWithAlphabetOut) // the good ones
	for i := 0x00; i < 0x80; i++ {
		for a := 0; a < 4; a++ {
			patterns[dcsWithAlphabetIn{byte(i), tpdu.Alphabet(a)}] =
				dcsWithAlphabetOut{tpdu.DCS(i&^0x0c | a<<2), nil}
		}
	}
	for i := 0xc0; i < 0xe0; i++ {
		patterns[dcsWithAlphabetIn{byte(i), tpdu.Alpha7Bit}] =
			dcsWithAlphabetOut{tpdu.DCS(i), nil}
	}
	for i := 0xe0; i < 0xf0; i++ {
		patterns[dcsWithAlphabetIn{byte(i), tpdu.AlphaUCS2}] =
			dcsWithAlphabetOut{tpdu.DCS(i), nil}
	}
	for i := 0xf0; i <= 0xff; i++ {
		for a := 0; a < 2; a++ {
			patterns[dcsWithAlphabetIn{byte(i), tpdu.Alphabet(a)}] =
				dcsWithAlphabetOut{tpdu.DCS(i&^0x0c | a<<2), nil}
		}
	}
	for i := 0x00; i <= 0xff; i++ {
		for a := 0; a < 4; a++ {
			p, ok := patterns[dcsWithAlphabetIn{byte(i), tpdu.Alphabet(a)}]
			if !ok {
				p = dcsWithAlphabetOut{tpdu.DCS(i), tpdu.ErrInvalid} // the bad ones
			}
			f := func(t *testing.T) {
				d := tpdu.DCS(i)
				dcs, err := d.WithAlphabet(tpdu.Alphabet(a))
				assert.Equal(t, p.err, err)
				assert.Equal(t, p.out, dcs)
			}
			t.Run(fmt.Sprintf("%08b_%d", i, a), f)
		}
	}
}

type dcsClassPattern struct {
	in  byte
	out tpdu.MessageClass
	err error
}

func TestDCSClass(t *testing.T) {
	// might as well test them all..
	patterns := []dcsClassPattern{}
	for i := 0; i < 4; i++ {
		m := byte(i << 5)
		patterns = append(patterns,
			dcsClassPattern{0x00 | m, tpdu.MClassUnknown, tpdu.ErrInvalid},
			dcsClassPattern{0x01 | m, tpdu.MClassUnknown, tpdu.ErrInvalid},
			dcsClassPattern{0x02 | m, tpdu.MClassUnknown, tpdu.ErrInvalid},
			dcsClassPattern{0x03 | m, tpdu.MClassUnknown, tpdu.ErrInvalid},
			dcsClassPattern{0x10 | m, tpdu.MClass0, nil},
			dcsClassPattern{0x11 | m, tpdu.MClass1, nil},
			dcsClassPattern{0x12 | m, tpdu.MClass2, nil},
			dcsClassPattern{0x13 | m, tpdu.MClass3, nil},
		)
	}
	for i := 0x80; i < 0xc0; i++ {
		patterns = append(patterns,
			dcsClassPattern{byte(i), tpdu.MClassUnknown, tpdu.ErrInvalid},
		)
	}
	for i := 0xc0; i < 0xe0; i++ {
		patterns = append(patterns,
			dcsClassPattern{byte(i), tpdu.MClassUnknown, nil},
		)
	}
	for i := 0xe0; i < 0xf0; i++ {
		patterns = append(patterns,
			dcsClassPattern{byte(i), tpdu.MClassUnknown, nil},
		)
	}
	for i := 0xf0; i <= 0xff; i++ {
		patterns = append(patterns,
			dcsClassPattern{byte(i), tpdu.MessageClass(i & 0x3), nil},
		)
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			d := tpdu.DCS(p.in)
			c, err := d.Class()
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, c)
		}
		t.Run(fmt.Sprintf("%08b", p.in), f)
	}
}

type dcsWithClassIn struct {
	in byte
	c  tpdu.MessageClass
}

type dcsWithClassOut struct {
	out tpdu.DCS
	err error
}

func TestDCSWithClass(t *testing.T) {
	// might as well test them all..
	patterns := make(map[dcsWithClassIn]dcsWithClassOut) // the good ones
	for i := 0x00; i < 0x80; i++ {
		for a := 0; a < 4; a++ {
			patterns[dcsWithClassIn{byte(i), tpdu.MessageClass(a)}] =
				dcsWithClassOut{tpdu.DCS(i&^0x03 | 0x10 | a), nil}
		}
	}
	for i := 0xf0; i <= 0xff; i++ {
		for a := 0; a < 4; a++ {
			patterns[dcsWithClassIn{byte(i), tpdu.MessageClass(a)}] =
				dcsWithClassOut{tpdu.DCS(i&^0x03 | a), nil}
		}
	}
	for i := 0x00; i <= 0xff; i++ {
		for a := 0; a < 4; a++ {
			p, ok := patterns[dcsWithClassIn{byte(i), tpdu.MessageClass(a)}]
			if !ok {
				p = dcsWithClassOut{tpdu.DCS(i), tpdu.ErrInvalid} // the bad ones
			}
			f := func(t *testing.T) {
				d := tpdu.DCS(i)
				dcs, err := d.WithClass(tpdu.MessageClass(a))
				assert.Equal(t, p.err, err)
				assert.Equal(t, p.out, dcs)
			}
			t.Run(fmt.Sprintf("%08b_%d", i, a), f)
		}
	}
}

func TestDCSCompressed(t *testing.T) {
	patterns := []struct {
		in  int
		out bool
	}{
		{0x00, false},
		{0x10, false},
		{0x20, true},
		{0x30, true},
		{0x40, false},
		{0x50, false},
		{0x60, true},
		{0x70, true},
		{0x80, false},
		{0x90, false},
		{0xa0, false},
		{0xb0, false},
		{0xc0, false},
		{0xd0, false},
		{0xe0, false},
		{0xf0, false},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			d := tpdu.DCS(p.in)
			c := d.Compressed()
			assert.Equal(t, p.out, c)
		}
		t.Run(fmt.Sprintf("%02x", p.in), f)
	}
}
