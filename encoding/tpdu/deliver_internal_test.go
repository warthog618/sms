// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/bcd"
	"github.com/warthog618/sms/encoding/semioctet"
)

type marshalDeliverTestPattern struct {
	name string
	in   Deliver
	out  []byte
	err  error
}

var marshalDeliverTestPatterns = []marshalDeliverTestPattern{
	{"haha",
		Deliver{
			TPDU: TPDU{FirstOctet: 4, UD: []byte("Hahahaha")},
			OA:   Address{Addr: "6391", TOA: 0x91},
			SCTS: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))}},
		[]byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71, 0x32, 0x20,
			0x05, 0x23, 0x08, 0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3},
		nil},
	{"bad oa",
		Deliver{
			TPDU: TPDU{FirstOctet: 4, UD: []byte("Hahahaha")},
			OA:   Address{Addr: "d391", TOA: 0x91},
			SCTS: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))}},
		nil,
		EncodeError("oa.addr", semioctet.ErrInvalidDigit('d'))},
	{"bad scts",
		Deliver{
			TPDU: TPDU{FirstOctet: 4, UD: []byte("Hahahaha")},
			OA:   Address{Addr: "6391", TOA: 0x91},
			SCTS: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 24*3600))}},
		nil,
		EncodeError("scts", bcd.ErrInvalidInteger(96))},
	{"bad ud",
		Deliver{
			TPDU: TPDU{FirstOctet: 4, DCS: 0x80, UD: []byte("Hahahaha")},
			OA:   Address{Addr: "6391", TOA: 0x91},
			SCTS: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))}},
		nil,
		EncodeError("ud.alphabet", ErrInvalid)},
}

func TestDeliverMarshalBinary(t *testing.T) {
	for _, p := range marshalDeliverTestPatterns {
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

type unmarshalDeliverTestPattern struct {
	name string
	in   []byte
	out  Deliver
	err  error
}

var unmarshalDeliverTestPatterns = []unmarshalDeliverTestPattern{
	{"haha", []byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71, 0x32,
		0x20, 0x05, 0x23, 0x08, 0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3},
		Deliver{
			TPDU: TPDU{FirstOctet: 4, UD: []byte("Hahahaha")},
			OA:   Address{Addr: "6391", TOA: 0x91},
			SCTS: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))}},
		nil},
	{"underflow fo", []byte{}, Deliver{}, DecodeError("firstOctet", 0, ErrUnderflow)},
	{"underflow oa", []byte{0x04, 0x04, 0x91, 0x36, 0xF9, 0x00, 0x00},
		Deliver{TPDU: TPDU{FirstOctet: 4}},
		DecodeError("oa.addr", 3, ErrUnderflow)},
	{"underflow pid", []byte{0x04, 0x04, 0x91, 0x36, 0x19},
		Deliver{
			TPDU: TPDU{FirstOctet: 4},
			OA:   Address{Addr: "6391", TOA: 0x91}},
		DecodeError("pid", 5, ErrUnderflow)},
	{"underflow dcs", []byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00},
		Deliver{
			TPDU: TPDU{FirstOctet: 4},
			OA:   Address{Addr: "6391", TOA: 0x91}},
		DecodeError("dcs", 6, ErrUnderflow)},
	{"underflow scts", []byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51},
		Deliver{
			TPDU: TPDU{FirstOctet: 4},
			OA:   Address{Addr: "6391", TOA: 0x91}},
		DecodeError("scts", 7, ErrUnderflow)},
	{"bad scts", []byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0xf1,
		0x32, 0x20, 0x05, 0x23},
		Deliver{
			TPDU: TPDU{FirstOctet: 4},
			OA:   Address{Addr: "6391", TOA: 0x91}},
		DecodeError("scts", 7, bcd.ErrInvalidOctet(0xf1))},
	{"underflow ud", []byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50,
		0x71, 0x32, 0x20, 0x05, 0x23, 0x08},
		Deliver{
			TPDU: TPDU{FirstOctet: 4},
			OA:   Address{Addr: "6391", TOA: 0x91},
			SCTS: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))}},
		DecodeError("ud.sm", 15, ErrUnderflow)},
}

func TestDeliverUnmarshalBinary(t *testing.T) {
	for _, p := range unmarshalDeliverTestPatterns {
		f := func(t *testing.T) {
			d := Deliver{}
			err := d.UnmarshalBinary(p.in)
			if err != p.err {
				t.Errorf("error decoding '%v': %v", p.in, err)
			}
			assert.Equal(t, p.out, d)
		}
		t.Run(p.name, f)
	}
}

func TestRegisterDeliverDecoder(t *testing.T) {
	dec := Decoder{map[byte]ConcreteDecoder{}}
	err := RegisterDeliverDecoder(&dec)
	if err != nil {
		t.Errorf("registration should not fail")
	}
	k := byte(MtDeliver) | (byte(MT) << 2)
	if cd, ok := dec.d[k]; !ok {
		t.Errorf("not registered with the correct key")
	} else {
		testDecodeDeliver(t, cd)
	}
	err = RegisterDeliverDecoder(&dec)
	if err == nil {
		t.Errorf("repeated registration should fail")
	}
}

func TestRegisterReservedDecoder(t *testing.T) {
	dec := Decoder{map[byte]ConcreteDecoder{}}
	err := RegisterReservedDecoder(&dec)
	if err != nil {
		t.Errorf("registration should not fail")
	}
	k := byte(MtReserved) | (byte(MT) << 2)
	if cd, ok := dec.d[k]; !ok {
		t.Errorf("not registered with the correct key")
	} else {
		testDecodeDeliver(t, cd)
	}
}

func testDecodeDeliver(t *testing.T, cd ConcreteDecoder) {
	b, derr := cd([]byte{})
	expected := DecodeError("firstOctet", 0, ErrUnderflow)
	if derr != expected {
		t.Errorf("returned unexpected error, expected %v, got %v\n", expected, derr)
	}
	if b != nil {
		t.Errorf("returned unexpected tpdu, expected nil, got %v\n", b)
	}
	b, derr = cd([]byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71,
		0x32, 0x20, 0x05, 0x23, 0x08, 0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3})
	if derr != nil {
		t.Errorf("returned unexpected error %v\n", derr)
	}
	if b != nil {
		deli, ok := b.(*Deliver)
		if !ok {
			t.Error("returned unexpected tpdu type")
		}
		if string(deli.UD) != "Hahahaha" {
			t.Error("returned unexpected tpdu user data")
		}
	}
}
