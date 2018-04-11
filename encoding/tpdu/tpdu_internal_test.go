// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTPDUdecodeUserData(t *testing.T) {
	patterns := []struct {
		name   string
		inPDU  TPDU
		inSrc  []byte
		outUD  UserData
		outUDH UserDataHeader
		err    error
	}{
		{"nil", TPDU{}, nil, nil, nil, DecodeError("udl", 0, ErrUnderflow)},
		{"empty", TPDU{}, []byte{0}, nil, nil, nil},
		{"7bit", TPDU{},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			[]byte("message"),
			nil, nil},
		{"sm overlength 7bit", TPDU{},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0xf1},
			nil, nil, DecodeError("sm", 1, ErrOverlength)},
		{"sm underflow", TPDU{},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97},
			nil, nil, DecodeError("sm", 1, ErrUnderflow)},
		{"7bit udh", TPDU{FirstOctet: 0x40},
			[]byte{0x0e, 5, 0, 3, 1, 2, 3, 0xda, 0xe5, 0xf9, 0x3c, 0x7c, 0x2e, 0x03},
			[]byte("message"),
			UserDataHeader([]InformationElement{{ID: 0, Data: []byte{1, 2, 3}}}),
			nil},
		{"8bit", TPDU{DCS: 0xf4},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			[]byte{0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			nil, nil},
		{"ucs2", TPDU{DCS: 0xe0},
			[]byte{0x0e, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65},
			[]byte{0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65},
			nil, nil},
		{"odd ucs2", TPDU{DCS: 0xe0},
			[]byte{0x0d, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00},
			nil, nil, DecodeError("sm", 1, ErrOverlength)},
		{"udh only", TPDU{FirstOctet: 0x40},
			[]byte{6, 5, 1, 3, 1, 2, 3},
			nil,
			UserDataHeader([]InformationElement{{ID: 1, Data: []byte{1, 2, 3}}}),
			nil},
		{"bad dcs", TPDU{FirstOctet: 0x40, DCS: 0xaa},
			[]byte{6, 4, 1, 3, 1, 2}, nil, nil, DecodeError("alphabet", 1, ErrInvalid)},
		{"overlength", TPDU{FirstOctet: 0x40},
			[]byte{6, 5, 1, 3, 1, 2, 3, 4}, nil, nil, DecodeError("ud", 1, ErrOverlength)},
		{"short udh", TPDU{FirstOctet: 0x40},
			[]byte{5, 5, 1, 3, 1, 2}, nil, nil, DecodeError("udh.ie", 2, ErrUnderflow)},
		{"bad udh", TPDU{FirstOctet: 0x40},
			[]byte{5, 4, 1, 3, 1, 2}, nil, nil, DecodeError("udh.ied", 4, ErrUnderflow)},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			err := p.inPDU.decodeUserData(p.inSrc)
			if err != p.err {
				t.Fatalf("error decoding '%v': %v", p.inSrc, err)
			}
			assert.Equal(t, p.outUDH, p.inPDU.UDH)
			assert.Equal(t, p.outUD, p.inPDU.UD)
		}
		t.Run(p.name, f)
	}
}

func TestTPDUencodeUserData(t *testing.T) {
	patterns := []struct {
		name string
		in   TPDU
		out  []byte
		err  error
	}{
		{"empty 7bit", TPDU{}, []byte{0}, nil},
		{"empty 8bit", TPDU{DCS: 0xf4}, []byte{0}, nil},
		{"empty ucs2", TPDU{DCS: 0xe0}, []byte{0}, nil},
		{"7bit", TPDU{UD: []byte("message")},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			nil},
		{"7bit udh", TPDU{FirstOctet: 0x40,
			UDH: UserDataHeader([]InformationElement{{ID: 0, Data: []byte{1, 2, 3}}}),
			UD:  []byte("message")},
			[]byte{0x0e, 5, 0, 3, 1, 2, 3, 0xda, 0xe5, 0xf9, 0x3c, 0x7c, 0x2e, 0x03},
			nil},
		{"8bit udh", TPDU{FirstOctet: 0x40, DCS: 0xf4,
			UDH: UserDataHeader([]InformationElement{{ID: 0, Data: []byte{1, 2, 3}}}),
			UD:  []byte("message")},
			[]byte{0x0d, 5, 0, 3, 1, 2, 3, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65},
			nil},
		{"ucs2 udh", TPDU{FirstOctet: 0x40, DCS: 0xe0,
			UDH: UserDataHeader([]InformationElement{{ID: 0, Data: []byte{1, 2, 3}}}),
			UD:  []byte{0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65}},
			[]byte{0x14, 5, 0, 3, 1, 2, 3, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65},
			nil},
		{"8bit", TPDU{DCS: 0xf4,
			UD: []byte{0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01}},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			nil},
		{"ucs2", TPDU{DCS: 0xe0,
			UD: []byte{0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65}},
			[]byte{0x0e, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00, 0x65},
			nil},
		{"odd ucs2", TPDU{DCS: 0xe0,
			UD: []byte{0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00}},
			nil,
			EncodeError("sm", ErrOddUCS2Length)},
		{"udh only", TPDU{FirstOctet: 0x40,
			UDH: UserDataHeader([]InformationElement{{ID: 1, Data: []byte{1, 2, 3}}})},
			[]byte{6, 5, 1, 3, 1, 2, 3},
			nil},
		{"unknown alphabet", TPDU{DCS: 0x80},
			nil,
			EncodeError("alphabet", ErrInvalid)},
	}
	for _, p := range patterns {
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
