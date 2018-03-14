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

type marshalStatusReportTestPattern struct {
	name string
	in   StatusReport
	out  []byte
	err  error
}

var marshalStatusReportTestPatterns = []marshalStatusReportTestPattern{
	{"minimal",
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3, udhiMask: 0x04},
			mr:       0x42,
			ra:       Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 6*3600))},
			st: 0xab,
		},
		[]byte{0x3, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05,
			0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab},
		nil},
	{"full",
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 7, pid: 0x89, dcs: 0x04, udhiMask: 0x04, ud: []byte("report")},
			mr:       0x42, pi: 0x07,
			ra: Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 6*3600))},
			st: 0xab,
		},
		[]byte{0x7, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05,
			0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x7, 0x89, 0x04, 0x06,
			0x72, 0x65, 0x70, 0x6f, 0x72, 0x74},
		nil},
	{"bad ra",
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3},
			mr:       0x42,
			ra:       Address{Addr: "63d1", TOA: 0x91},
		},
		nil,
		EncodeError("ra.addr", semioctet.ErrInvalidDigit('d'))},
	{"bad scts",
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3, udhiMask: 0x04},
			mr:       0x42,
			ra:       Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 24*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 6*3600))},
			st: 0xab},
		nil,
		EncodeError("scts", bcd.ErrInvalidInteger(96))},
	{"bad dt",
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3, udhiMask: 0x04},
			mr:       0x42,
			ra:       Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 24*3600))},
			st: 0xab},
		nil,
		EncodeError("dt", bcd.ErrInvalidInteger(96))},
	{"bad ud",
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 1, udhiMask: 0x02, dcs: 0x80, ud: []byte("report")},
			pi:       0x06},
		nil,
		EncodeError("ud.alphabet", ErrInvalid)},
}

func TestStatusReportMarshalBinary(t *testing.T) {
	for _, p := range marshalStatusReportTestPatterns {
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

type unmarshalStatusReportTestPattern struct {
	name string
	in   []byte
	out  StatusReport
	err  error
}

var unmarshalStatusReportTestPatterns = []unmarshalStatusReportTestPattern{
	{"minimal",
		[]byte{0x3, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05,
			0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3},
			mr:       0x42,
			ra:       Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 6*3600))},
			st: 0xab,
		},
		nil},
	{"pid",
		[]byte{0x3, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05,
			0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x01, 0x89},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3, pid: 0x89},
			mr:       0x42, pi: 0x01,
			ra: Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 6*3600))},
			st: 0xab,
		},
		nil},
	{"dcs",
		[]byte{0x3, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05,
			0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x02, 0x04},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3, dcs: 0x04},
			mr:       0x42, pi: 0x02,
			ra: Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 6*3600))},
			st: 0xab,
		},
		nil},
	{"ud",
		[]byte{0x3, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05,
			0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x06, 0x04, 0x06,
			0x72, 0x65, 0x70, 0x6f, 0x72, 0x74},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3, udhiMask: 0x04, dcs: 0x04, ud: []byte("report")},
			mr:       0x42, pi: 0x06,
			ra: Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 6*3600))},
			st: 0xab,
		},
		nil},
	{"underflow fo", []byte{}, StatusReport{}, DecodeError("firstOctet", 0, ErrUnderflow)},
	{"underflow mr", []byte{0x03},
		StatusReport{BaseTPDU: BaseTPDU{firstOctet: 0x03}},
		DecodeError("mr", 1, ErrUnderflow)},
	{"underflow ra", []byte{0x03, 0x42},
		StatusReport{BaseTPDU: BaseTPDU{firstOctet: 3}, mr: 0x42},
		DecodeError("ra.addr", 2, ErrUnderflow)},
	{"underflow scts", []byte{0x03, 0x42, 0x04, 0x91, 0x36, 0x19},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3},
			mr:       0x42,
			ra:       Address{Addr: "6391", TOA: 0x91},
		},
		DecodeError("scts", 6, ErrUnderflow)},
	{"underflow dt", []byte{0x03, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71,
		0x32, 0x20, 0x05, 0x23},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3},
			mr:       0x42,
			ra:       Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
		},
		DecodeError("dt", 13, ErrUnderflow)},
	{"bad scts",
		[]byte{0x03, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0xf1, 0x32, 0x20, 0x05, 0x23},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3},
			mr:       0x42,
			ra:       Address{Addr: "6391", TOA: 0x91},
		},
		DecodeError("scts", 6, bcd.ErrInvalidOctet(0xf1))},
	{"bad dt",
		[]byte{0x3, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05,
			0x23, 0x51, 0x40, 0xc1, 0x32, 0x20, 0x05, 0x42},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3},
			mr:       0x42,
			ra:       Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
		},
		DecodeError("dt", 13, bcd.ErrInvalidOctet(0xc1))},
	{"underflow st",
		[]byte{0x3, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05,
			0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3},
			mr:       0x42,
			ra:       Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 6*3600))},
		},
		DecodeError("st", 20, ErrUnderflow)},
	{"underflow pid",
		[]byte{0x3, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05,
			0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x01},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3},
			mr:       0x42, pi: 0x01,
			ra: Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 6*3600))},
			st: 0xab},
		DecodeError("pid", 22, ErrUnderflow)},
	{"underflow dcs",
		[]byte{0x3, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05,
			0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x02},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3},
			mr:       0x42, pi: 0x02,
			ra: Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 6*3600))},
			st: 0xab},
		DecodeError("dcs", 22, ErrUnderflow)},
	{"underflow ud",
		[]byte{0x3, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05,
			0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x04},
		StatusReport{
			BaseTPDU: BaseTPDU{firstOctet: 3, udhiMask: 0x04},
			mr:       0x42, pi: 0x04,
			ra: Address{Addr: "6391", TOA: 0x91},
			scts: Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
				time.FixedZone("SCTS", 8*3600))},
			dt: Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
				time.FixedZone("SCTS", 6*3600))},
			st: 0xab},
		DecodeError("ud.udl", 22, ErrUnderflow)},
}

func TestStatusReportUnmarshalBinary(t *testing.T) {
	for _, p := range unmarshalStatusReportTestPatterns {
		f := func(t *testing.T) {
			d := StatusReport{}
			err := d.UnmarshalBinary(p.in)
			if err != p.err {
				t.Errorf("error decoding '%v': %v", p.in, err)
			}
			assert.Equal(t, p.out, d)
		}
		t.Run(p.name, f)
	}
}

func TestRegisterStatusReportDecoder(t *testing.T) {
	dec := Decoder{map[byte]ConcreteDecoder{}}
	err := RegisterStatusReportDecoder(&dec)
	if err != nil {
		t.Errorf("registration should not fail")
	}
	k := byte(MtCommand) | (byte(MT) << 2)
	if cd, ok := dec.d[k]; !ok {
		t.Errorf("not registered with the correct key")
	} else {
		testDecodeStatusReport(t, cd)
	}
	err = RegisterStatusReportDecoder(&dec)
	if err == nil {
		t.Errorf("repeated registration should fail")
	}
}

func testDecodeStatusReport(t *testing.T, cd ConcreteDecoder) {
	b, derr := cd([]byte{})
	expected := DecodeError("firstOctet", 0, ErrUnderflow)
	if derr != expected {
		t.Errorf("returned unexpected error, expected %v, got %v\n", expected, derr)
	}
	if b != nil {
		t.Errorf("returned unexpected tpdu, expected nil, got %v\n", b)
	}
	b, derr = cd([]byte{0x3, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71,
		0x32, 0x20, 0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab})
	if derr != nil {
		t.Errorf("returned unexpected error %v\n", derr)
	}
	if b != nil {
		_, ok := b.(*StatusReport)
		if !ok {
			t.Error("returned unexpected tpdu type")
		}
	}
}
