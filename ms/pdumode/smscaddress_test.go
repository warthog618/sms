// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package pdumode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/semioctet"
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/ms/pdumode"
)

func TestSMSCAddressMarshalBinary(t *testing.T) {
	patterns := []struct {
		name string
		in   pdumode.SMSCAddress
		out  []byte
		err  error
	}{
		{
			"empty",
			pdumode.SMSCAddress{},
			[]byte{0},
			nil,
		},
		{
			"number",
			pdumode.SMSCAddress{
				tpdu.Address{Addr: "61409865629", TOA: 0x91},
			},
			[]byte{7, 0x91, 0x16, 0x04, 0x89, 0x56, 0x26, 0xf9},
			nil,
		},
		{
			"number alphabet",
			pdumode.SMSCAddress{
				tpdu.Address{Addr: "0123456789*#abc", TOA: 0x91},
			},
			[]byte{9, 0x91, 0x10, 0x32, 0x54, 0x76, 0x98, 0xba, 0xdc, 0xfe},
			nil,
		},
		{
			"alpha",
			pdumode.SMSCAddress{
				tpdu.Address{Addr: "messages", TOA: 0xd1},
			},
			nil,
			tpdu.EncodeError("addr", semioctet.ErrInvalidDigit(0x6d)),
		},
		{
			"invalid number",
			pdumode.SMSCAddress{
				tpdu.Address{Addr: "6140f98656", TOA: 0x91},
			},
			nil,
			tpdu.EncodeError("addr", semioctet.ErrInvalidDigit('f')),
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

func TestSMSCAddressUnmarshalBinary(t *testing.T) {
	patterns := []struct {
		name string
		in   []byte
		out  pdumode.SMSCAddress
		n    int
		err  error
	}{
		{
			"nil",
			nil,
			pdumode.SMSCAddress{},
			0,
			tpdu.NewDecodeError("length", 0, tpdu.ErrUnderflow),
		},
		{
			"only toa",
			[]byte{1, 0},
			pdumode.SMSCAddress{},
			2,
			nil,
		},
		{
			"number",
			[]byte{7, 0x91, 0x16, 0x04, 0x89, 0x56, 0x26, 0xf9},
			pdumode.SMSCAddress{
				tpdu.Address{Addr: "61409865629", TOA: 0x91},
			},
			8,
			nil,
		},
		{
			"number alphabet",
			[]byte{9, 0x91, 0x10, 0x32, 0x54, 0x76, 0x98, 0xba, 0xdc, 0xfe},
			pdumode.SMSCAddress{
				tpdu.Address{Addr: "0123456789*#abc", TOA: 0x91},
			},
			10,
			nil,
		},
		{
			"alpha",
			[]byte{8, 0xd1, 0xED, 0xF2, 0x7C, 0x1E, 0x3E, 0x97, 0xE7},
			pdumode.SMSCAddress{
				tpdu.Address{Addr: "bc2a7c1c3797c", TOA: 0xd1},
			},
			9,
			nil,
		},
		{
			"zero length",
			[]byte{0},
			pdumode.SMSCAddress{},
			1,
			nil,
		},
		{
			"short number",
			[]byte{11, 0x91, 0x16, 0x04, 0x89, 0x56},
			pdumode.SMSCAddress{},
			6,
			tpdu.NewDecodeError("addr", 2, tpdu.ErrUnderflow),
		},
		{
			"short number pad",
			[]byte{12, 0x91, 0x16, 0x04, 0x89, 0x56, 0x97, 0xf7},
			pdumode.SMSCAddress{},
			8,
			tpdu.NewDecodeError("addr", 2, tpdu.ErrUnderflow),
		},
		{
			"underflow alpha",
			[]byte{5, 0xd1, 0xCF, 0xE5, 0x39},
			pdumode.SMSCAddress{},
			5,
			tpdu.NewDecodeError("addr", 2, tpdu.ErrUnderflow),
		},
		{
			"underflow toa",
			[]byte{1}, pdumode.SMSCAddress{},
			1,
			tpdu.NewDecodeError("toa", 1, tpdu.ErrUnderflow),
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			a := pdumode.SMSCAddress{}
			n, err := a.UnmarshalBinary(p.in)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.n, n)
			assert.Equal(t, p.out, a)
		}
		t.Run(p.name, f)
	}
}
