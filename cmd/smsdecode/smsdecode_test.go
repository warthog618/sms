// SPDX-License-Identifier: MIT
//
// Copyright © 2019 Kent Gibson <warthog618@gmail.com>.package main

package main

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestDecode(t *testing.T) {
	patterns := []struct {
		name string
		pdu  string
		pm   bool
		mo   bool
		out  *tpdu.TPDU
		err  error
	}{
		{
			"invalid hex",
			"bad hex",
			false,
			false,
			nil,
			hex.InvalidByteError(' '),
		},
		{
			"decode fail",
			"010005912143f500000bc8329bfd06dddf7236",
			false,
			true,
			nil,
			tpdu.DecodeError{
				Field:  "SmsSubmit.ud.sm",
				Offset: 10,
				Err:    tpdu.ErrUnderflow,
			},
		},
		{
			"pdumode decode fail",
			"07911614220991",
			true,
			true,
			nil,
			tpdu.DecodeError{
				Field:  "addr",
				Offset: 2,
				Err:    tpdu.ErrUnderflow,
			},
		},
		{
			"submit",
			"010005912143f500000bc8329bfd06dddf723619",
			false,
			true,
			&tpdu.TPDU{
				Direction:  1,
				FirstOctet: 1,
				DA:         tpdu.Address{TOA: 0x91, Addr: "12345"},
				UD: []byte{
					0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64,
				},
			},
			nil,
		},
		{
			"submit pdumode",
			"07911614220991f1010005912143f500000bc8329bfd06dddf723619",
			true,
			true,
			&tpdu.TPDU{
				Direction:  1,
				FirstOctet: 1,
				DA:         tpdu.Address{TOA: 0x91, Addr: "12345"},
				UD: []byte{
					0x48, 0x65, 0x6c, 0x6c, 0x6f, 0x20, 0x77, 0x6f, 0x72, 0x6c, 0x64,
				},
			},
			nil,
		},
	}

	for _, p := range patterns {
		f := func(t *testing.T) {
			out, _, err := decode(p.pdu, p.pm, p.mo)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}
