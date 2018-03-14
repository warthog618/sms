// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type decodeUDTestPattern struct {
	name   string
	inPDU  BaseTPDU
	inSrc  []byte
	outUD  UserData
	outUDH UserDataHeader
	err    error
}

var decodeUDTestPatterns = []decodeUDTestPattern{
	{"nil", BaseTPDU{}, nil, nil, nil, DecodeError("udl", 0, ErrUnderflow)},
	{"empty", BaseTPDU{}, []byte{0}, nil, nil, nil},
	{"7bit", BaseTPDU{},
		[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
		[]byte("message"),
		nil, nil},
	{"sm overlength 7bit", BaseTPDU{},
		[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0xf1},
		nil, nil, DecodeError("sm", 1, ErrOverlength)},
	{"sm underflow", BaseTPDU{},
		[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97},
		nil, nil, DecodeError("sm", 1, ErrUnderflow)},
	{"7bit udh", BaseTPDU{udhiMask: 0x20, firstOctet: 0x20},
		[]byte{0x0e, 5, 0, 3, 1, 2, 3, 0xda, 0xe5, 0xf9, 0x3c, 0x7c, 0x2e, 0x03},
		[]byte("message"),
		UserDataHeader([]InformationElement{InformationElement{ID: 0, Data: []byte{1, 2, 3}}}),
		nil},
	{"8bit", BaseTPDU{dcs: 0xf4},
		[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
		[]byte{0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
		nil, nil},
	{"ucs2", BaseTPDU{dcs: 0xe0},
		[]byte{0x0e, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65},
		[]byte{0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65},
		nil, nil},
	{"odd ucs2", BaseTPDU{dcs: 0xe0},
		[]byte{0x0d, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00},
		nil, nil, DecodeError("sm", 1, ErrOverlength)},
	{"udh only", BaseTPDU{udhiMask: 0x20, firstOctet: 0x20},
		[]byte{6, 5, 1, 3, 1, 2, 3},
		nil,
		UserDataHeader([]InformationElement{InformationElement{ID: 1, Data: []byte{1, 2, 3}}}),
		nil},
	{"bad dcs", BaseTPDU{udhiMask: 0x20, firstOctet: 0x20, dcs: 0xaa},
		[]byte{6, 4, 1, 3, 1, 2}, nil, nil, DecodeError("alphabet", 1, ErrInvalid)},
	{"overlength", BaseTPDU{udhiMask: 0x20, firstOctet: 0x20},
		[]byte{6, 5, 1, 3, 1, 2, 3, 4}, nil, nil, DecodeError("udh", 1, ErrOverlength)},
	{"short udh", BaseTPDU{udhiMask: 0x20, firstOctet: 0x20},
		[]byte{5, 5, 1, 3, 1, 2}, nil, nil, DecodeError("udh.ie", 2, ErrUnderflow)},
	{"bad udh", BaseTPDU{udhiMask: 0x20, firstOctet: 0x20},
		[]byte{5, 4, 1, 3, 1, 2}, nil, nil, DecodeError("udh.ied", 4, ErrUnderflow)},
}

func TestBaseTPDUdecodeUserData(t *testing.T) {
	for _, p := range decodeUDTestPatterns {
		f := func(t *testing.T) {
			err := p.inPDU.decodeUserData(p.inSrc)
			if err != p.err {
				t.Fatalf("error decoding '%v': %v", p.inSrc, err)
			}
			assert.Equal(t, p.outUDH, p.inPDU.udh)
			assert.Equal(t, p.outUD, p.inPDU.ud)
		}
		t.Run(p.name, f)
	}
}

type encodeUDTestPattern struct {
	name string
	in   BaseTPDU
	out  []byte
	err  error
}

var encodeUDTestPatterns = []encodeUDTestPattern{
	{"empty 7bit", BaseTPDU{}, []byte{0}, nil},
	{"empty 8bit", BaseTPDU{dcs: 0xf4}, []byte{0}, nil},
	{"empty ucs2", BaseTPDU{dcs: 0xe0}, []byte{0}, nil},
	{"7bit", BaseTPDU{ud: []byte("message")},
		[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
		nil},
	{"7bit udh", BaseTPDU{udhiMask: 0x20, firstOctet: 0x20,
		udh: UserDataHeader([]InformationElement{InformationElement{ID: 0, Data: []byte{1, 2, 3}}}),
		ud:  []byte("message")},
		[]byte{0x0e, 5, 0, 3, 1, 2, 3, 0xda, 0xe5, 0xf9, 0x3c, 0x7c, 0x2e, 0x03},
		nil},
	{"8bit udh", BaseTPDU{udhiMask: 0x20, firstOctet: 0x20, dcs: 0xf4,
		udh: UserDataHeader([]InformationElement{InformationElement{ID: 0, Data: []byte{1, 2, 3}}}),
		ud:  []byte("message")},
		[]byte{0x0d, 5, 0, 3, 1, 2, 3, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
		nil},
	{"ucs2 udh", BaseTPDU{udhiMask: 0x20, firstOctet: 0x20, dcs: 0xe0,
		udh: UserDataHeader([]InformationElement{InformationElement{ID: 0, Data: []byte{1, 2, 3}}}),
		ud:  []byte{0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65}},
		[]byte{0x14, 5, 0, 3, 1, 2, 3, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65},
		nil},
	{"8bit", BaseTPDU{dcs: 0xf4,
		ud: []byte{0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01}},
		[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
		nil},
	{"ucs2", BaseTPDU{dcs: 0xe0,
		ud: []byte{0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65}},
		[]byte{0x0e, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65},
		nil},
	{"odd ucs2", BaseTPDU{dcs: 0xe0,
		ud: []byte{0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00}},
		nil,
		EncodeError("sm", ErrOddUCS2Length)},
	{"udh only", BaseTPDU{udhiMask: 0x20, firstOctet: 0x20,
		udh: UserDataHeader([]InformationElement{InformationElement{ID: 1, Data: []byte{1, 2, 3}}})},
		[]byte{6, 5, 1, 3, 1, 2, 3},
		nil},
	{"unknown alphabet", BaseTPDU{dcs: 0x80},
		nil,
		EncodeError("alphabet", ErrInvalid)},
}

func TestBaseTPDUencodeUserData(t *testing.T) {
	for _, p := range encodeUDTestPatterns {
		f := func(t *testing.T) {
			d, err := p.in.encodeUserData()
			if err != p.err {
				t.Fatalf("error encoding '%v': %v", p.in, err)
			}
			assert.Equal(t, p.out, d)
		}
		t.Run(p.name, f)
	}
}
