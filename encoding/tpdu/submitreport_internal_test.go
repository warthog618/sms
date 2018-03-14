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
)

type marshalSubmitReportTestPattern struct {
	name string
	in   SubmitReport
	out  []byte
	err  error
}

var marshalSubmitReportTestPatterns = []marshalSubmitReportTestPattern{
	{"minimal",
		SubmitReport{
			BaseTPDU: BaseTPDU{firstOctet: 1, udhiMask: 0x04},
			fcs:      0x12,
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
		},
		[]byte{0x01, 0x12, 0x00, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23},
		nil},
	{"full",
		SubmitReport{
			BaseTPDU: BaseTPDU{firstOctet: 1, udhiMask: 0x04, pid: 0xab, dcs: 0x04, ud: []byte("report")},
			fcs:      0x12, pi: 0x07,
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
		},
		[]byte{0x1, 0x12, 0x07, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0xab,
			0x04, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74},
		nil},
	{"bad scts",
		SubmitReport{
			BaseTPDU: BaseTPDU{firstOctet: 1, udhiMask: 0x04, dcs: 0x80, ud: []byte("report")},
			fcs:      0x12, pi: 0x07,
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 24*3600))}},
		nil,
		EncodeError("scts", bcd.ErrInvalidInteger(96))},
	{"bad ud",
		SubmitReport{
			BaseTPDU: BaseTPDU{firstOctet: 1, udhiMask: 0x04, dcs: 0x80, ud: []byte("report")},
			fcs:      0x12, pi: 0x06},
		nil,
		EncodeError("ud.alphabet", ErrInvalid)},
}

func TestSubmitReportMarshalBinary(t *testing.T) {
	for _, p := range marshalSubmitReportTestPatterns {
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

type unmarshalSubmitReportTestPattern struct {
	name string
	in   []byte
	out  SubmitReport
	err  error
}

var unmarshalSubmitReportTestPatterns = []unmarshalSubmitReportTestPattern{
	{"minimal", []byte{0x01, 0x12, 0x00, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23},
		SubmitReport{
			BaseTPDU: BaseTPDU{firstOctet: 1, udhiMask: 0x04},
			fcs:      0x12,
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
		},
		nil},
	{"pid", []byte{0x01, 0x12, 0x01, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0xab},
		SubmitReport{
			BaseTPDU: BaseTPDU{firstOctet: 1, udhiMask: 0x04, pid: 0xab},
			fcs:      0x12, pi: 0x01,
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
		},
		nil},
	{"dcs", []byte{0x00, 0x12, 0x02, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0x04},
		SubmitReport{
			BaseTPDU: BaseTPDU{firstOctet: 0, udhiMask: 0x04, dcs: 0x04},
			fcs:      0x12, pi: 02,
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
		},
		nil},
	{"ud", []byte{0x00, 0x12, 0x06, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0x04, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74},
		SubmitReport{
			BaseTPDU: BaseTPDU{firstOctet: 0, udhiMask: 0x04, dcs: 0x04, ud: []byte("report")},
			fcs:      0x12, pi: 0x06,
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
		},
		nil},
	{"underflow fo", []byte{}, SubmitReport{}, DecodeError("firstOctet", 0, ErrUnderflow)},
	{"underflow fcs", []byte{0x00},
		SubmitReport{BaseTPDU: BaseTPDU{firstOctet: 0x00}},
		DecodeError("fcs", 1, ErrUnderflow)},
	{"underflow pi", []byte{0x00, 0x12},
		SubmitReport{BaseTPDU: BaseTPDU{firstOctet: 0x00},
			fcs: 0x12},
		DecodeError("pi", 2, ErrUnderflow)},
	{"underflow scts", []byte{0x01, 0x12, 0x00},
		SubmitReport{
			BaseTPDU: BaseTPDU{firstOctet: 1},
			fcs:      0x12,
		},
		DecodeError("scts", 3, ErrUnderflow)},
	{"bad scts", []byte{0x01, 0x12, 0x00, 0x51, 0x50, 0xf1,
		0x32, 0x20, 0x05, 0x23},
		SubmitReport{
			BaseTPDU: BaseTPDU{firstOctet: 1},
			fcs:      0x12,
		},
		DecodeError("scts", 3, bcd.ErrInvalidOctet(0xf1))},

	{"underflow pid", []byte{0x00, 0x12, 0x01, 0x51, 0x50, 0x71,
		0x32, 0x20, 0x05, 0x23},
		SubmitReport{BaseTPDU: BaseTPDU{firstOctet: 0x00},
			fcs: 0x12, pi: 0x01,
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))}},
		DecodeError("pid", 10, ErrUnderflow)},
	{"underflow dcs", []byte{0x00, 0x12, 0x02, 0x51, 0x50, 0x71,
		0x32, 0x20, 0x05, 0x23},
		SubmitReport{BaseTPDU: BaseTPDU{firstOctet: 0x00},
			fcs: 0x12, pi: 0x02,
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))}},
		DecodeError("dcs", 10, ErrUnderflow)},
	{"underflow ud", []byte{0x00, 0x12, 0x06, 0x51, 0x50, 0x71,
		0x32, 0x20, 0x05, 0x23, 0x04},
		SubmitReport{BaseTPDU: BaseTPDU{firstOctet: 0x00, udhiMask: 0x04, dcs: 0x04},
			fcs: 0x12, pi: 0x06,
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
		},
		DecodeError("ud.udl", 11, ErrUnderflow)},
}

func TestSubmitReportUnmarshalBinary(t *testing.T) {
	for _, p := range unmarshalSubmitReportTestPatterns {
		f := func(t *testing.T) {
			d := SubmitReport{}
			err := d.UnmarshalBinary(p.in)
			if err != p.err {
				t.Errorf("error decoding '%v': %v", p.in, err)
			}
			assert.Equal(t, p.out, d)
		}
		t.Run(p.name, f)
	}
}

func TestRegisterSubmitReportDecoder(t *testing.T) {
	dec := Decoder{map[byte]ConcreteDecoder{}}
	err := RegisterSubmitReportDecoder(&dec)
	if err != nil {
		t.Errorf("registration should not fail")
	}
	k := byte(MtSubmit) | (byte(MT) << 2)
	if cd, ok := dec.d[k]; !ok {
		t.Errorf("not registered with the correct key")
	} else {
		testDecodeSubmitReport(t, cd)
	}
	err = RegisterSubmitReportDecoder(&dec)
	if err == nil {
		t.Errorf("repeated registration should fail")
	}
}

func testDecodeSubmitReport(t *testing.T, cd ConcreteDecoder) {
	b, derr := cd([]byte{})
	expected := DecodeError("firstOctet", 0, ErrUnderflow)
	if derr != expected {
		t.Errorf("returned unexpected error, expected %v, got %v\n", expected, derr)
	}
	if b != nil {
		t.Errorf("returned unexpected tpdu, expected nil, got %v\n", b)
	}
	b, derr = cd([]byte{0x01, 0x12, 0x00, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23})
	if derr != nil {
		t.Errorf("returned unexpected error %v\n", derr)
	}
	if b != nil {
		_, ok := b.(*SubmitReport)
		if !ok {
			t.Error("returned unexpected tpdu type")
		}
	}
}
