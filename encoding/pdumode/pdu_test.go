// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package pdumode_test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warthog618/sms/encoding/pdumode"
	"github.com/warthog618/sms/encoding/semioctet"
	"github.com/warthog618/sms/encoding/tpdu"
)

type testPattern struct {
	name string
	pdu  string
	smsc *pdumode.SMSCAddress
	tpdu []byte
	err  error
}

func TestUnmarshalBinary(t *testing.T) {
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
			pdu, err := pdumode.UnmarshalBinary(b)
			assert.Equal(t, p.err, err)
			if err == nil {
				require.NotNil(t, pdu)
				assert.Equal(t, *p.smsc, pdu.SMSC)
				assert.Equal(t, p.tpdu, pdu.TPDU)
			} else {
				assert.Nil(t, pdu)
			}
		}
		t.Run(p.name, f)
	}
}

func TestUnmarshalHexString(t *testing.T) {
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
			pdu, err := pdumode.UnmarshalHexString(p.pdu)
			assert.Equal(t, p.err, err)
			if err == nil {
				require.NotNil(t, pdu)
				assert.Equal(t, *p.smsc, pdu.SMSC)
				assert.Equal(t, p.tpdu, pdu.TPDU)
			} else {
				assert.Nil(t, pdu)
			}
		}
		t.Run(p.name, f)
	}
}

func TestMarshalBinary(t *testing.T) {
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
			pdu := pdumode.PDU{SMSC: *p.smsc, TPDU: p.tpdu}
			b, err := pdu.MarshalBinary()
			s := hex.EncodeToString(b)
			assert.Equal(t, p.pdu, s)
			assert.Equal(t, p.err, err)
		}
		t.Run(p.name, f)
	}
}

func TestMarshalHexString(t *testing.T) {
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
			pdu := pdumode.PDU{SMSC: *p.smsc, TPDU: p.tpdu}
			b, err := pdu.MarshalHexString()
			assert.Equal(t, p.pdu, b)
			assert.Equal(t, p.err, err)
		}
		t.Run(p.name, f)
	}
}
