// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package pdumode_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/semioctet"
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/ms/pdumode"
)

type testPattern struct {
	name string
	pdu  string
	smsc *pdumode.SMSCAddress
	tpdu []byte
	err  error
}

func TestDecode(t *testing.T) {
	decodePatterns := []testPattern{
		{
			"empty",
			"",
			nil,
			nil,
			tpdu.NewDecodeError("length", 0, tpdu.ErrUnderflow),
		},
		{
			"valid",
			"0791361907002039010203040506070809",
			&pdumode.SMSCAddress{
				tpdu.Address{Addr: "639170000293", TOA: 0x91},
			},
			[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09},
			nil,
		},
	}
	for _, p := range decodePatterns {
		f := func(t *testing.T) {
			b, err := hex.DecodeString(p.pdu)
			if err != nil {
				t.Fatalf("error converting in: %v", err)
			}
			smsc, tpdu, err := pdumode.Decode(b)
			assert.Equal(t, p.smsc, smsc)
			assert.Equal(t, p.tpdu, tpdu)
			assert.Equal(t, p.err, err)
		}
		t.Run(p.name, f)
	}
}

func TestDecodeString(t *testing.T) {
	decodePatterns := []testPattern{
		{
			"nothex",
			"nothex",
			nil,
			nil,
			hex.InvalidByteError('n'),
		},
		{
			"empty",
			"",
			nil,
			nil,
			tpdu.NewDecodeError("length", 0, tpdu.ErrUnderflow),
		},
		{
			"valid", "0791361907002039010203040506070809",
			&pdumode.SMSCAddress{
				tpdu.Address{Addr: "639170000293", TOA: 0x91},
			},
			[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09},
			nil,
		},
	}
	for _, p := range decodePatterns {
		f := func(t *testing.T) {
			smsc, tpdu, err := pdumode.DecodeString(p.pdu)
			assert.Equal(t, p.smsc, smsc)
			assert.Equal(t, p.tpdu, tpdu)
			assert.Equal(t, p.err, err)
		}
		t.Run(p.name, f)
	}
}

func TestEncode(t *testing.T) {
	patterns := []testPattern{
		{
			"empty",
			"00",
			&pdumode.SMSCAddress{},
			nil,
			nil,
		},
		{
			"valid", "0791361907002039010203040506070809",
			&pdumode.SMSCAddress{
				tpdu.Address{Addr: "639170000293", TOA: 0x91},
			},
			[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09},
			nil,
		},
		{
			"invalid addr", "",
			&pdumode.SMSCAddress{
				tpdu.Address{Addr: "banana"},
			},
			nil,
			tpdu.EncodeError("addr", semioctet.ErrInvalidDigit(0x6e)),
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			pdu, err := pdumode.Encode(*p.smsc, p.tpdu)
			s := hex.EncodeToString(pdu)
			assert.Equal(t, p.pdu, s)
			assert.Equal(t, p.err, err)
		}
		t.Run(p.name, f)
	}
}

func TestEncodeToString(t *testing.T) {
	patterns := []testPattern{
		{
			"empty",
			"00",
			&pdumode.SMSCAddress{},
			nil,
			nil,
		},
		{
			"valid",
			"0791361907002039010203040506070809",
			&pdumode.SMSCAddress{
				tpdu.Address{Addr: "639170000293", TOA: 0x91},
			},
			[]byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09},
			nil,
		},
		{
			"invalid addr",
			"",
			&pdumode.SMSCAddress{
				tpdu.Address{Addr: "banana"},
			},
			nil,
			tpdu.EncodeError("addr", semioctet.ErrInvalidDigit(0x6e)),
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			pdu, err := pdumode.EncodeToString(*p.smsc, p.tpdu)
			assert.Equal(t, p.pdu, pdu)
			assert.Equal(t, p.err, err)
		}
		t.Run(p.name, f)
	}
}
