// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package tpdu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDecodeUserData(t *testing.T) {
	patterns := []struct {
		name   string
		inPDU  TPDU
		inSrc  []byte
		outUD  UserData
		outUDH UserDataHeader
		err    error
	}{
		{"nil",
			TPDU{},
			nil,
			nil,
			nil,
			NewDecodeError("udl", 0, ErrUnderflow),
		},
		{"empty",
			TPDU{},
			[]byte{0},
			nil,
			nil,
			nil,
		},
		{"7bit",
			TPDU{},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			[]byte("message"),
			nil,
			nil,
		},
		{"sm overlength 7bit",
			TPDU{},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0xf1},
			nil,
			nil,
			NewDecodeError("sm", 1, ErrOverlength),
		},
		{"sm underflow",
			TPDU{},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97},
			nil,
			nil,
			NewDecodeError("sm", 1, ErrUnderflow),
		},
		{"7bit udh",
			TPDU{FirstOctet: 0x40},
			[]byte{
				0x0e, 0x05, 0x00, 0x03, 0x01, 0x02, 0x03, 0xda, 0xe5, 0xf9,
				0x3c, 0x7c, 0x2e, 0x03,
			},
			[]byte("message"),
			UserDataHeader([]InformationElement{{ID: 0, Data: []byte{1, 2, 3}}}),
			nil,
		},
		{"8bit",
			TPDU{DCS: 0xf4},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			[]byte{0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			nil,
			nil,
		},
		{"ucs2",
			TPDU{DCS: 0xe0},
			[]byte{
				0x0e, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00,
				0x61, 0x00, 0x67, 0x00, 0x65,
			},
			[]byte{
				0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61,
				0x00, 0x67, 0x00, 0x65,
			},
			nil,
			nil,
		},
		{"odd ucs2",
			TPDU{DCS: 0xe0},
			[]byte{
				0x0d, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00,
				0x61, 0x00, 0x67, 0x00,
			},
			nil,
			nil,
			NewDecodeError("sm", 1, ErrOddUCS2Length),
		},
		{"udh only",
			TPDU{FirstOctet: 0x40},
			[]byte{0x06, 0x05, 0x01, 0x03, 0x01, 0x02, 0x03},
			nil,
			UserDataHeader([]InformationElement{{ID: 1, Data: []byte{1, 2, 3}}}),
			nil,
		},
		{"bad dcs",
			TPDU{
				FirstOctet: 0x40,
				DCS:        0xaa,
			},
			[]byte{0x06, 0x04, 0x01, 0x03, 0x01, 0x02},
			nil,
			nil,
			NewDecodeError("alphabet", 1, ErrInvalid),
		},
		{"overlength",
			TPDU{FirstOctet: 0x40},
			[]byte{0x06, 0x05, 0x01, 0x03, 0x01, 0x02, 0x03, 0x04},
			nil,
			nil,
			NewDecodeError("ud", 1, ErrOverlength),
		},
		{"short udh",
			TPDU{FirstOctet: 0x40},
			[]byte{0x05, 0x05, 0x01, 0x03, 0x01, 0x02},
			nil,
			nil,
			NewDecodeError("udh.ie", 2, ErrUnderflow),
		},
		{"bad udh",
			TPDU{FirstOctet: 0x40},
			[]byte{0x05, 0x04, 0x01, 0x03, 0x01, 0x02},
			nil,
			nil,
			NewDecodeError("udh.ied", 4, ErrUnderflow),
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			err := p.inPDU.decodeUserData(p.inSrc)
			require.Equal(t, p.err, err)
			assert.Equal(t, p.outUDH, p.inPDU.UDH)
			assert.Equal(t, p.outUD, p.inPDU.UD)
		}
		t.Run(p.name, f)
	}
}

func TestEncodeUserData(t *testing.T) {
	patterns := []struct {
		name string
		in   TPDU
		out  []byte
		err  error
	}{
		{"empty 7bit",
			TPDU{},
			[]byte{0},
			nil,
		},
		{"empty 8bit",
			TPDU{DCS: 0xf4},
			[]byte{0},
			nil,
		},
		{"empty ucs2",
			TPDU{DCS: 0xe0},
			[]byte{0},
			nil,
		},
		{"7bit",
			TPDU{UD: []byte("message")},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			nil,
		},
		{"7bit udh",
			TPDU{
				FirstOctet: 0x40,
				UDH: UserDataHeader(
					[]InformationElement{{ID: 0, Data: []byte{1, 2, 3}}}),
				UD: []byte("message"),
			},
			[]byte{
				0x0e, 0x05, 0x00, 0x03, 0x01, 0x02, 0x03, 0xda, 0xe5, 0xf9,
				0x3c, 0x7c, 0x2e, 0x03,
			},
			nil,
		},
		{"8bit udh",
			TPDU{
				FirstOctet: 0x40,
				DCS:        0xf4,
				UDH: UserDataHeader(
					[]InformationElement{{ID: 0, Data: []byte{1, 2, 3}}}),
				UD: []byte("message"),
			},
			[]byte{
				0x0d, 0x05, 0x00, 0x03, 0x01, 0x02, 0x03, 0x6d, 0x65, 0x73,
				0x73, 0x61, 0x67, 0x65,
			},
			nil,
		},
		{"ucs2 udh",
			TPDU{
				FirstOctet: 0x40,
				DCS:        0xe0,
				UDH: UserDataHeader(
					[]InformationElement{{ID: 0, Data: []byte{1, 2, 3}}}),
				UD: []byte{
					0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61,
					0x00, 0x67, 0x00, 0x65,
				},
			},
			[]byte{
				0x14, 0x05, 0x00, 0x03, 0x01, 0x02, 0x03, 0x00, 0x6D, 0x00,
				0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61, 0x00, 0x67, 0x00,
				0x65,
			},
			nil,
		},
		{"8bit",
			TPDU{
				DCS: 0xf4,
				UD:  []byte{0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01}},
			[]byte{0x07, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0x01},
			nil,
		},
		{"ucs2",
			TPDU{
				DCS: 0xe0,
				UD: []byte{
					0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61,
					0x00, 0x67, 0x00, 0x65,
				},
			},
			[]byte{
				0x0e, 0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00,
				0x61, 0x00, 0x67, 0x00, 0x65,
			},
			nil,
		},
		{"odd ucs2",
			TPDU{
				DCS: 0xe0,
				UD: []byte{
					0x00, 0x6D, 0x00, 0x65, 0x00, 0x73, 0x00, 0x73, 0x00, 0x61,
					0x00, 0x67, 0x00,
				},
			},
			nil,
			EncodeError("sm", ErrOddUCS2Length),
		},
		{"udh only",
			TPDU{
				FirstOctet: 0x40,
				UDH: UserDataHeader(
					[]InformationElement{{ID: 1, Data: []byte{1, 2, 3}}})},
			[]byte{0x06, 0x05, 0x01, 0x03, 0x01, 0x02, 0x03},
			nil,
		},
		{"unknown alphabet",
			TPDU{DCS: 0x80},
			nil,
			EncodeError("alphabet", ErrInvalid),
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			d, err := p.in.encodeUserData()
			require.Equal(t, p.err, err)
			assert.Equal(t, p.out, d)
		}
		t.Run(p.name, f)
	}
}

func TestChunk(t *testing.T) {
	patterns := []struct {
		name  string
		in    []byte
		alpha Alphabet
		bs    int
		out   [][]byte
	}{
		{
			"7bit",
			[]byte{1, 2, 0x1b, 4, 5, 6, 7, 8},
			Alpha7Bit,
			3,
			[][]byte{
				{0x01, 0x02},
				{0x1b, 0x04, 0x05},
				{0x06, 0x07, 0x08},
			},
		},
		{
			"8bit",
			[]byte{1, 2, 0x1b, 4, 5, 6, 7, 8},
			Alpha8Bit,
			3,
			[][]byte{
				{0x01, 0x02, 0x1b},
				{0x04, 0x05, 0x06},
				{0x07, 0x08},
			},
		},
		{
			"ucs2",
			[]byte{1, 2, 0xd8, 4, 5, 6, 7, 8, 9, 10},
			AlphaUCS2,
			4,
			[][]byte{
				{0x01, 0x02},
				{0xd8, 0x04, 0x05, 0x06},
				{0x07, 0x08, 0x09, 0x0a},
			},
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out := chunk(p.in, p.alpha, p.bs)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}

func TestChunk7Bit(t *testing.T) {
	patterns := []struct {
		name string
		in   []byte
		bs   int
		out  [][]byte
	}{
		{
			"nil",
			nil,
			2,
			nil,
		},
		{
			"empty",
			[]byte{},
			2,
			nil,
		},
		{
			"integral",
			[]byte{1, 2, 3, 4},
			2,
			[][]byte{
				{1, 2},
				{3, 4},
			},
		},
		{
			"residual",
			[]byte{1, 2, 3, 4},
			3,
			[][]byte{
				{1, 2, 3},
				{4},
			},
		},
		{
			"three",
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			3,
			[][]byte{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8},
			},
		},
		{
			"escaped",
			[]byte{1, 2, 0x1b, 4, 5, 6, 7, 8},
			3,
			[][]byte{
				{1, 2},
				{0x1b, 4, 5},
				{6, 7, 8},
			},
		},
		{
			"double escaped",
			[]byte{1, 0x1b, 0x1b, 4, 5, 6, 7, 8},
			3,
			[][]byte{
				{1, 0x1b, 0x1b},
				{4, 5, 6},
				{7, 8},
			},
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out := chunk7Bit(p.in, p.bs)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}

func TestChunk8Bit(t *testing.T) {
	patterns := []struct {
		name string
		in   []byte
		bs   int
		out  [][]byte
	}{
		{
			"nil",
			nil,
			2,
			nil,
		},
		{
			"empty",
			[]byte{},
			2,
			nil,
		},
		{
			"integral",
			[]byte{1, 2, 3, 4},
			2,
			[][]byte{
				{1, 2},
				{3, 4},
			},
		},
		{
			"residual",
			[]byte{1, 2, 3, 4},
			3,
			[][]byte{
				{1, 2, 3},
				{4},
			},
		},
		{
			"three",
			[]byte{1, 2, 3, 4, 5, 6, 7, 8},
			3,
			[][]byte{
				{1, 2, 3},
				{4, 5, 6},
				{7, 8},
			},
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out := chunk8Bit(p.in, p.bs)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}

func TestChunkUCS2(t *testing.T) {
	patterns := []struct {
		name string
		in   []byte
		bs   int
		out  [][]byte
	}{
		{
			"nil",
			nil,
			2,
			nil,
		},
		{
			"empty",
			[]byte{},
			2,
			nil,
		},
		{
			"integral",
			[]byte{1, 2, 3, 4},
			2,
			[][]byte{
				{1, 2},
				{3, 4},
			},
		},
		{
			"odd bs",
			[]byte{1, 2, 3, 4},
			3,
			[][]byte{
				{1, 2},
				{3, 4},
			},
		},
		{
			"three",
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			4,
			[][]byte{
				{1, 2, 3, 4},
				{5, 6, 7, 8},
				{9, 10},
			},
		},
		{
			"surrogate",
			[]byte{1, 2, 0xd8, 4, 5, 6, 7, 8, 9, 10},
			4,
			[][]byte{
				{1, 2},
				{0xd8, 4, 5, 6},
				{7, 8, 9, 10},
			},
		},
		{
			"odd msg",
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9},
			4,
			[][]byte{
				{1, 2, 3, 4},
				{5, 6, 7, 8},
				{9},
			},
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out := chunkUCS2(p.in, p.bs)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}
