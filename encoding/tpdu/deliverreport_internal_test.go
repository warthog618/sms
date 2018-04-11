// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type marshalDeliverReportTestPattern struct {
	name string
	in   DeliverReport
	out  []byte
	err  error
}

var marshalDeliverReportTestPatterns = []marshalDeliverReportTestPattern{
	{"minimal",
		DeliverReport{
			TPDU: TPDU{FirstOctet: 0},
			FCS:  0x12,
		},
		[]byte{0x00, 0x12, 0x00},
		nil},
	{"full",
		DeliverReport{
			TPDU: TPDU{FirstOctet: 0, PID: 0xab, DCS: 0x04, UD: []byte("report")},
			FCS:  0x12, PI: 0x07,
		},
		[]byte{0x00, 0x12, 0x07, 0xab, 0x04, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74},
		nil},
	{"bad ud",
		DeliverReport{
			TPDU: TPDU{FirstOctet: 0, DCS: 0x80, UD: []byte("report")},
			FCS:  0x12, PI: 0x06},
		nil,
		EncodeError("ud.alphabet", ErrInvalid)},
}

func TestDeliverReportMarshalBinary(t *testing.T) {
	for _, p := range marshalDeliverReportTestPatterns {
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

type unmarshalDeliverReportTestPattern struct {
	name string
	in   []byte
	out  DeliverReport
	err  error
}

var unmarshalDeliverReportTestPatterns = []unmarshalDeliverReportTestPattern{
	{"minimal", []byte{0x00, 0x12, 0x00},
		DeliverReport{
			TPDU: TPDU{FirstOctet: 0},
			FCS:  0x12,
		},
		nil},
	{"pid", []byte{0x00, 0x12, 0x01, 0xab},
		DeliverReport{
			TPDU: TPDU{FirstOctet: 0, PID: 0xab},
			FCS:  0x12, PI: 0x01,
		},
		nil},
	{"dcs", []byte{0x00, 0x12, 0x02, 0x04},
		DeliverReport{
			TPDU: TPDU{FirstOctet: 0, DCS: 0x04},
			FCS:  0x12, PI: 02,
		},
		nil},
	{"ud", []byte{0x00, 0x12, 0x06, 0x04, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74},
		DeliverReport{
			TPDU: TPDU{FirstOctet: 0, DCS: 0x04, UD: []byte("report")},
			FCS:  0x12, PI: 0x06,
		},
		nil},
	{"underflow fo", []byte{}, DeliverReport{}, DecodeError("firstOctet", 0, ErrUnderflow)},
	{"underflow fcs", []byte{0x00},
		DeliverReport{TPDU: TPDU{FirstOctet: 0x00}},
		DecodeError("fcs", 1, ErrUnderflow)},
	{"underflow pi", []byte{0x00, 0x12},
		DeliverReport{TPDU: TPDU{FirstOctet: 0x00},
			FCS: 0x12},
		DecodeError("pi", 2, ErrUnderflow)},
	{"underflow pid", []byte{0x00, 0x12, 0x01},
		DeliverReport{TPDU: TPDU{FirstOctet: 0x00},
			FCS: 0x12, PI: 0x01},
		DecodeError("pid", 3, ErrUnderflow)},
	{"underflow dcs", []byte{0x00, 0x12, 0x02},
		DeliverReport{TPDU: TPDU{FirstOctet: 0x00},
			FCS: 0x12, PI: 0x02},
		DecodeError("dcs", 3, ErrUnderflow)},
	{"underflow ud", []byte{0x00, 0x12, 0x04},
		DeliverReport{TPDU: TPDU{FirstOctet: 0x00},
			FCS: 0x12, PI: 0x04},
		DecodeError("ud.udl", 3, ErrUnderflow)},
}

func TestDeliverReportUnmarshalBinary(t *testing.T) {
	for _, p := range unmarshalDeliverReportTestPatterns {
		f := func(t *testing.T) {
			d := DeliverReport{}
			err := d.UnmarshalBinary(p.in)
			if err != p.err {
				t.Errorf("error decoding '%v': %v", p.in, err)
			}
			assert.Equal(t, p.out, d)
		}
		t.Run(p.name, f)
	}
}

func TestRegisterDeliverReportDecoder(t *testing.T) {
	dec := Decoder{map[byte]ConcreteDecoder{}}
	err := RegisterDeliverReportDecoder(&dec)
	if err != nil {
		t.Errorf("registration should not fail")
	}
	k := byte(MtDeliver) | (byte(MO) << 2)
	if cd, ok := dec.d[k]; !ok {
		t.Errorf("not registered with the correct key")
	} else {
		testDecodeDeliverReport(t, cd)
	}
	err = RegisterDeliverReportDecoder(&dec)
	if err == nil {
		t.Errorf("repeated registration should fail")
	}
}

func testDecodeDeliverReport(t *testing.T, cd ConcreteDecoder) {
	b, derr := cd([]byte{})
	expected := DecodeError("firstOctet", 0, ErrUnderflow)
	if derr != expected {
		t.Errorf("returned unexpected error, expected %v, got %v\n", expected, derr)
	}
	if b != nil {
		t.Errorf("returned unexpected tpdu, expected nil, got %v\n", b)
	}
	b, derr = cd([]byte{0x00, 0x12, 0x00})
	if derr != nil {
		t.Errorf("returned unexpected error %v\n", derr)
	}
	if b != nil {
		_, ok := b.(*DeliverReport)
		if !ok {
			t.Error("returned unexpected tpdu type")
		}
	}
}
