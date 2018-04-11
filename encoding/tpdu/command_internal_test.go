// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/semioctet"
)

type marshalCommandTestPattern struct {
	name string
	in   Command
	out  []byte
	err  error
}

var marshalCommandTestPatterns = []marshalCommandTestPattern{
	{"command",
		Command{
			TPDU: TPDU{FirstOctet: 2,
				PID: 0xab, UD: []byte("a command")},
			MR: 0x42, CT: 0x89, MN: 0x34,
			DA: Address{Addr: "6391", TOA: 0x91},
		},
		[]byte{0x02, 0x42, 0xab, 0x89, 0x34, 0x04, 0x91, 0x36, 0x19, 0x09, 0x61,
			0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64},
		nil},
	{"bad da",
		Command{
			TPDU: TPDU{FirstOctet: 2,
				PID: 0xab, UD: []byte("a command")},
			MR: 0x42, CT: 0x89, MN: 0x34,
			DA: Address{Addr: "d391", TOA: 0x91},
		},
		nil,
		EncodeError("da.addr", semioctet.ErrInvalidDigit('d'))},
}

func TestCommandMarshalBinary(t *testing.T) {
	for _, p := range marshalCommandTestPatterns {
		f := func(t *testing.T) {
			b, err := p.in.MarshalBinary()
			if err != p.err {
				t.Errorf("error encoding '%v': %v", p.in, err)
			}
			assert.Equal(t, p.out, b)
		}
		t.Run(p.name, f)
	}
}

type unmarshalCommandTestPattern struct {
	name string
	in   []byte
	out  Command
	err  error
}

var unmarshalCommandTestPatterns = []unmarshalCommandTestPattern{
	{"command", []byte{0x02, 0x42, 0xab, 0x89, 0x34, 0x04, 0x91, 0x36, 0x19, 0x09, 0x61,
		0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64},
		Command{
			TPDU: TPDU{FirstOctet: 2, DCS: 0x04,
				PID: 0xab, UD: []byte("a command")},
			MR: 0x42, CT: 0x89, MN: 0x34,
			DA: Address{Addr: "6391", TOA: 0x91},
		},
		nil},
	{"underflow fo", []byte{}, Command{}, DecodeError("firstOctet", 0, ErrUnderflow)},
	{"underflow mr", []byte{0x02},
		Command{TPDU: TPDU{FirstOctet: 2}},
		DecodeError("mr", 1, ErrUnderflow)},
	{"underflow pid", []byte{0x02, 0x42},
		Command{
			TPDU: TPDU{FirstOctet: 2},
			MR:   0x42},
		DecodeError("pid", 2, ErrUnderflow)},
	{"underflow ct", []byte{0x02, 0x42, 0xab},
		Command{
			TPDU: TPDU{FirstOctet: 2, PID: 0xab},
			MR:   0x42},
		DecodeError("ct", 3, ErrUnderflow)},
	{"underflow mn", []byte{0x02, 0x42, 0xab, 0x89},
		Command{
			TPDU: TPDU{FirstOctet: 2, PID: 0xab},
			MR:   0x42, CT: 0x89},
		DecodeError("mn", 4, ErrUnderflow)},
	{"underflow da", []byte{0x02, 0x42, 0xab, 0x89, 0x34, 0x04, 0x91, 0x36, 0xF9, 0x09},
		Command{
			TPDU: TPDU{FirstOctet: 2, PID: 0xab},
			MR:   0x42, CT: 0x89, MN: 0x34},
		DecodeError("da.addr", 7, ErrUnderflow)},
	{"underflow ud", []byte{0x02, 0x42, 0xab, 0x89, 0x34, 0x04, 0x91, 0x36, 0x19},
		Command{
			TPDU: TPDU{FirstOctet: 2, PID: 0xab, DCS: 0x04},
			MR:   0x42, CT: 0x89, MN: 0x34,
			DA: Address{Addr: "6391", TOA: 0x91},
		},
		DecodeError("ud.udl", 9, ErrUnderflow)},
}

func TestCommandUnmarshalBinary(t *testing.T) {
	for _, p := range unmarshalCommandTestPatterns {
		f := func(t *testing.T) {
			d := Command{}
			err := d.UnmarshalBinary(p.in)
			if err != p.err {
				t.Errorf("error decoding '%v': %v", p.in, err)
			}
			assert.Equal(t, p.out, d)
		}
		t.Run(p.name, f)
	}
}

func TestRegisterCommandDecoder(t *testing.T) {
	dec := Decoder{map[byte]ConcreteDecoder{}}
	err := RegisterCommandDecoder(&dec)
	if err != nil {
		t.Errorf("registration should not fail")
	}
	k := byte(MtCommand) | (byte(MO) << 2)
	if cd, ok := dec.d[k]; !ok {
		t.Errorf("not registered with the correct key")
	} else {
		testDecodeCommand(t, cd)
	}
	err = RegisterCommandDecoder(&dec)
	if err == nil {
		t.Errorf("repeated registration should fail")
	}
}

func testDecodeCommand(t *testing.T, cd ConcreteDecoder) {
	b, derr := cd([]byte{})
	expected := DecodeError("firstOctet", 0, ErrUnderflow)
	if derr != expected {
		t.Errorf("returned unexpected error, expected %v, got %v\n", expected, derr)
	}
	if b != nil {
		t.Errorf("returned unexpected tpdu, expected nil, got %v\n", b)
	}
	b, derr = cd([]byte{0x02, 0x42, 0xab, 0x89, 0x34, 0x04, 0x91, 0x36, 0x19, 0x09, 0x61,
		0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64})
	if derr != nil {
		t.Errorf("returned unexpected error: %v\n", derr)
	}
	if b != nil {
		deli, ok := b.(*Command)
		if !ok {
			t.Error("returned unexpected tpdu type")
		}
		if string(deli.UD) != "a command" {
			t.Error("returned unexpected tpdu user data")
		}
	}
}
