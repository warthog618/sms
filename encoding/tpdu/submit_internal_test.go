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

type marshalSubmitTestPattern struct {
	name string
	in   Submit
	out  []byte
	err  error
}

var marshalSubmitTestPatterns = []marshalSubmitTestPattern{
	{"haha",
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 1, ud: []byte("Hahahaha")},
			da:       Address{Addr: "6391", TOA: 0x91},
		},
		[]byte{0x01, 0x00, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x08, 0xC8, 0x30,
			0x3A, 0x8C, 0x0E, 0xA3, 0xC3},
		nil},
	{"vp",
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 1, ud: []byte("Hahahaha")},
			da:       Address{Addr: "6391", TOA: 0x91},
			vp:       ValidityPeriod{VpfRelative, Timestamp{}, time.Duration(6000000000000), 0},
		},
		[]byte{0x01, 0x00, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x13, 0x08, 0xC8,
			0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3},
		nil},
	{"bad da",
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 4, ud: []byte("Hahahaha")},
			da:       Address{Addr: "d391", TOA: 0x91},
		},
		nil,
		EncodeError("da.addr", semioctet.ErrInvalidDigit('d'))},
	{"bad vp",
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 1, ud: []byte("Hahahaha")},
			da:       Address{Addr: "6391", TOA: 0x91},
			vp:       ValidityPeriod{6, Timestamp{}, 0, 0},
		},
		nil,
		EncodeError("vp.vpf", ErrInvalid)},
	{"bad ud",
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 4, dcs: 0x80, ud: []byte("Hahahaha")},
			da:       Address{Addr: "6391", TOA: 0x91},
		},
		nil,
		EncodeError("ud.alphabet", ErrInvalid)},
}

func TestSubmitMarshalBinary(t *testing.T) {
	for _, p := range marshalSubmitTestPatterns {
		f := func(t *testing.T) {
			b, err := p.in.MarshalBinary()
			if err != p.err {
				t.Errorf("error encoding '%v': %v %v", p.in, err, p.err)
			}
			assert.Equal(t, p.out, b)
		}
		t.Run(p.name, f)
	}
}

type unmarshalSubmitTestPattern struct {
	name string
	in   []byte
	out  Submit
	err  error
}

var unmarshalSubmitTestPatterns = []unmarshalSubmitTestPattern{
	{"haha", []byte{0x01, 0x23, 0x04, 0x91, 0x36, 0x19, 0x34, 0x00, 0x08, 0xC8,
		0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3},
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 1, pid: 0x34, ud: []byte("Hahahaha")},
			mr:       0x23,
			da:       Address{Addr: "6391", TOA: 0x91},
		},
		nil},
	{"vp", []byte{0x11, 0x23, 0x04, 0x91, 0x36, 0x19, 0x34, 0x00, 0x45, 0x08, 0xC8,
		0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3},
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 17, pid: 0x34, ud: []byte("Hahahaha")},
			mr:       0x23,
			da:       Address{Addr: "6391", TOA: 0x91},
			vp:       ValidityPeriod{Format: VpfRelative, Duration: time.Duration(60 * 350 * 1000000000)},
		},
		nil},
	{"underflow fo", []byte{}, Submit{}, DecodeError("firstOctet", 0, ErrUnderflow)},
	{"underflow mr", []byte{0x01},
		Submit{BaseTPDU: BaseTPDU{firstOctet: 1}},
		DecodeError("mr", 1, ErrUnderflow)},
	{"underflow da",
		[]byte{0x01, 0x23},
		Submit{BaseTPDU: BaseTPDU{firstOctet: 1},
			mr: 0x23},
		DecodeError("da.addr", 2, ErrUnderflow)},
	{"bad da",
		[]byte{0x01, 0x23, 0x04, 0x91, 0x36, 0xF9, 0x00},
		Submit{BaseTPDU: BaseTPDU{firstOctet: 1},
			mr: 0x23},
		DecodeError("da.addr", 4, ErrUnderflow)},
	{"underflow pid",
		[]byte{0x01, 0x23, 0x04, 0x91, 0x36, 0x19},
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 1},
			mr:       0x23,
			da:       Address{Addr: "6391", TOA: 0x91}},
		DecodeError("pid", 6, ErrUnderflow)},
	{"underflow dcs",
		[]byte{0x01, 0x23, 0x04, 0x91, 0x36, 0x19, 0x00},
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 1},
			mr:       0x23,
			da:       Address{Addr: "6391", TOA: 0x91}},
		DecodeError("dcs", 7, ErrUnderflow)},
	{"underflow vp",
		[]byte{0x11, 0x23, 0x04, 0x91, 0x36, 0x19, 0x34, 0x00},
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 17, pid: 0x34},
			mr:       0x23,
			da:       Address{Addr: "6391", TOA: 0x91}},
		DecodeError("vp", 8, ErrUnderflow)},
	{"bad vp", []byte{0x19, 0x23, 0x04, 0x91, 0x36, 0x19, 0x34, 0x00, 0x45, 0x08,
		0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3},
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 0x19, pid: 0x34},
			mr:       0x23,
			da:       Address{Addr: "6391", TOA: 0x91},
		},
		DecodeError("vp", 8, bcd.ErrInvalidOctet(0xc8))},
	{"underflow ud",
		[]byte{0x01, 0x23, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50,
			0x71, 0x32, 0x20, 0x05, 0x23, 0x08},
		Submit{
			BaseTPDU: BaseTPDU{firstOctet: 1},
			mr:       0x23,
			da:       Address{Addr: "6391", TOA: 0x91},
		},
		DecodeError("ud.sm", 9, ErrUnderflow)},
}

func TestSubmitUnmarshalBinary(t *testing.T) {
	for _, p := range unmarshalSubmitTestPatterns {
		f := func(t *testing.T) {
			s := Submit{}
			err := s.UnmarshalBinary(p.in)
			if err != p.err {
				t.Errorf("error decoding '%v': %v", p.in, err)
			}
			assert.Equal(t, p.out, s)
		}
		t.Run(p.name, f)
	}
}

func TestRegisterSubmitDecoder(t *testing.T) {
	dec := Decoder{map[byte]ConcreteDecoder{}}
	err := RegisterSubmitDecoder(&dec)
	if err != nil {
		t.Errorf("registration should not fail")
	}
	k := byte(MtSubmit) | (byte(MO) << 2)
	if cd, ok := dec.d[k]; !ok {
		t.Errorf("not registered with the correct key")
	} else {
		testDecodeSubmit(t, cd)
	}
	err = RegisterSubmitDecoder(&dec)
	if err == nil {
		t.Errorf("repeated registration should fail")
	}
}

func testDecodeSubmit(t *testing.T, cd ConcreteDecoder) {
	b, derr := cd([]byte{})
	expected := DecodeError("firstOctet", 0, ErrUnderflow)
	if derr != expected {
		t.Errorf("returned unexpected error, expected %v, got %v\n", expected, derr)
	}
	if b != nil {
		t.Errorf("returned unexpected tpdu, expected nil, got %v\n", b)
	}
	b, derr = cd([]byte{0x01, 0x00, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x08, 0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3})
	if derr != nil {
		t.Errorf("returned unexpected error %v\n", derr)
	}
	if b != nil {
		s, ok := b.(*Submit)
		if !ok {
			t.Error("returned unexpected tpdu type")
		}
		if string(s.UD()) != "Hahahaha" {
			t.Error("returned unexpected tpdu user data")
		}
	}
}
