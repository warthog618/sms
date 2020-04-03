// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package tpdu_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/gsm7"
	"github.com/warthog618/sms/encoding/semioctet"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestNewAddress(t *testing.T) {
	a := tpdu.NewAddress()
	assert.Equal(t, uint8(0x80), a.TOA)
	a = tpdu.NewAddress(tpdu.FromNumber("1234"))
	assert.Equal(t, uint8(0x91), a.TOA)
	assert.Equal(t, "1234", a.Addr)
	a = tpdu.NewAddress(tpdu.FromNumber("+4321"))
	assert.Equal(t, uint8(0x91), a.TOA)
	assert.Equal(t, "4321", a.Addr)
}

type addressMarshalPattern struct {
	name string
	in   tpdu.Address
	out  []byte
	err  error
}

func TestAddressMarshalBinary(t *testing.T) {
	patterns := []addressMarshalPattern{
		{"empty",
			tpdu.Address{},
			[]byte{0, 0},
			nil,
		},
		{"number",
			tpdu.Address{Addr: "61409865629", TOA: 0x91},
			[]byte{11, 0x91, 0x16, 0x04, 0x89, 0x56, 0x26, 0xf9},
			nil,
		},
		{"number alphabet",
			tpdu.Address{Addr: "0123456789*#abc", TOA: 0x91},
			[]byte{15, 0x91, 0x10, 0x32, 0x54, 0x76, 0x98, 0xba, 0xdc, 0xfe},
			nil,
		},
		{"alpha",
			tpdu.Address{Addr: "messages", TOA: 0xd1},
			[]byte{14, 0xd1, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0xE7},
			nil,
		},
		{"alpha odd",
			tpdu.Address{Addr: "message", TOA: 0xd1},
			[]byte{13, 0xd1, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			nil,
		},
		{"alpha Vodafone",
			tpdu.Address{Addr: "Vodafone", TOA: 0xd0},
			[]byte{14, 0xd0, 0xD6, 0x37, 0x39, 0x6C, 0x7E, 0xBB, 0xCB},
			nil,
		},
		{"invalid number",
			tpdu.Address{Addr: "6140f98656", TOA: 0x91},
			nil,
			tpdu.EncodeError("addr", semioctet.ErrInvalidDigit('f')),
		},
		// test characters only available in the extension table - which should
		// be unavailable.
		{"invalid alpha euro",
			tpdu.Address{Addr: "a euro €32", TOA: 0xd1},
			nil,
			tpdu.EncodeError("addr", gsm7.ErrInvalidUTF8('€')),
		},
		{"invalid alpha bar",
			tpdu.Address{Addr: "a bar | ", TOA: 0xd1},
			nil,
			tpdu.EncodeError("addr", gsm7.ErrInvalidUTF8('|')),
		},
		// test characters not available in the default character set at all.
		{"invalid alpha",
			tpdu.Address{Addr: "mes⌘sages", TOA: 0xd1},
			nil,
			tpdu.EncodeError("addr", gsm7.ErrInvalidUTF8('⌘')),
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			b, err := p.in.MarshalBinary()
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, b)
		}
		t.Run(p.name, f)
	}
}

type addressUnmarshalPattern struct {
	name string
	in   []byte
	out  tpdu.Address
	n    int
	err  error
}

func TestAddressUnmarshalBinary(t *testing.T) {
	patterns := []addressUnmarshalPattern{
		{"empty",
			[]byte{0, 0},
			tpdu.Address{},
			2,
			nil,
		},
		{"number",
			[]byte{11, 0x91, 0x16, 0x04, 0x89, 0x56, 0x26, 0xf9},
			tpdu.Address{Addr: "61409865629", TOA: 0x91},
			8,
			nil,
		},
		{"number alphabet",
			[]byte{15, 0x91, 0x10, 0x32, 0x54, 0x76, 0x98, 0xba, 0xdc, 0xfe},
			tpdu.Address{Addr: "0123456789*#abc", TOA: 0x91},
			10,
			nil},
		{"alpha",
			[]byte{14, 0xd1, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0xE7},
			tpdu.Address{Addr: "messages", TOA: 0xd1},
			9,
			nil,
		},
		{"alpha odd",
			[]byte{13, 0xd1, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			tpdu.Address{Addr: "message", TOA: 0xd1},
			9,
			nil,
		},
		{"alpha Vodafone",
			[]byte{14, 0xd0, 0xD6, 0x37, 0x39, 0x6C, 0x7E, 0xBB, 0xCB},
			tpdu.Address{Addr: "Vodafone", TOA: 0xd0},
			9,
			nil,
		},
		{"overlong number",
			[]byte{11, 0x91, 0x16, 0x04, 0x89, 0x56, 0x26, 0x09},
			tpdu.Address{},
			8,
			tpdu.DecodeError("addr", 2, semioctet.ErrMissingFill),
		},
		{"short binary",
			[]byte{0}, tpdu.Address{},
			0,
			tpdu.DecodeError("addr", 0, tpdu.ErrUnderflow),
		},
		{"short number",
			[]byte{11, 0x91, 0x16, 0x04, 0x89, 0x56},
			tpdu.Address{},
			6,
			tpdu.DecodeError("addr", 2, tpdu.ErrUnderflow),
		},
		{"short number pad",
			[]byte{12, 0x91, 0x16, 0x04, 0x89, 0x56, 0x97, 0xf7},
			tpdu.Address{},
			8,
			tpdu.DecodeError("addr", 2, tpdu.ErrUnderflow),
		},
		{"invalid alpha bar",
			[]byte{10, 0xd1, 0xED, 0xF2, 0x7C, 0x03, 0x9c, 0x87, 0xCF, 0xE5, 0x39},
			tpdu.Address{},
			2,
			tpdu.DecodeError("addr", 2, gsm7.ErrInvalidSeptet(0x40)),
		},
		{"underflow alpha",
			[]byte{10, 0xd1, 0xCF, 0xE5, 0x39},
			tpdu.Address{},
			5,
			tpdu.DecodeError("addr", 2, tpdu.ErrUnderflow),
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			a := tpdu.Address{}
			n, err := a.UnmarshalBinary(p.in)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.n, n)
			assert.Equal(t, p.out, a)
		}
		t.Run(p.name, f)
	}
}

func TestAddressNumber(t *testing.T) {
	a := tpdu.Address{Addr: "61409865629", TOA: 0}
	assert.Equal(t, "61409865629", a.Number())
	a.SetTypeOfNumber(tpdu.TonInternational)
	assert.Equal(t, "+61409865629", a.Number())
}

func TestAddressNumberingPlan(t *testing.T) {
	patterns := []tpdu.NumberingPlan{0, 2, 0xf, 0x12}
	for _, p := range patterns {
		f := func(t *testing.T) {
			a := tpdu.NewAddress()
			a.SetNumberingPlan((p))
			assert.Equal(t, p&0x0f, a.NumberingPlan())
		}
		t.Run(fmt.Sprintf("%02x", p), f)
	}
}

func TestAddressSetNumber(t *testing.T) {
	a := tpdu.NewAddress()
	assert.Equal(t, uint8(0x80), a.TOA)
	a.SetNumber("1234")
	assert.Equal(t, uint8(0x91), a.TOA)
	assert.Equal(t, "1234", a.Addr)
	assert.Equal(t, "+1234", a.Number())
	a.SetNumber("+4321")
	assert.Equal(t, uint8(0x91), a.TOA)
	assert.Equal(t, "4321", a.Addr)
	assert.Equal(t, "+4321", a.Number())
}

func TestAddressTypeOfNumber(t *testing.T) {
	patterns := []tpdu.TypeOfNumber{0, 2, 3, 5}
	for _, p := range patterns {
		f := func(t *testing.T) {
			a := tpdu.NewAddress()
			a.SetTypeOfNumber((p))
			ton := a.TypeOfNumber()
			assert.Equal(t, p&0x0f, ton)
		}
		t.Run(fmt.Sprintf("%02x", p), f)
	}
}
