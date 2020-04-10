// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package sms_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms"
	"github.com/warthog618/sms/encoding/gsm7/charset"
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/encoding/ucs2"
)

func TestConcatenate(t *testing.T) {
	patterns := []struct {
		name    string
		in      []*tpdu.TPDU
		options []sms.ConcatOption
		out     []byte
		err     error
	}{
		{
			"two segment 7bit",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UD: []byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you mi"),
				},
				&tpdu.TPDU{
					UD: []byte("ght think"),
				},
			},
			nil,
			twoSegmentMsg,
			nil,
		},
		{
			"two segment 8bit",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					DCS: tpdu.Dcs8BitData,
					UD:  []byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you mi"),
				},
				&tpdu.TPDU{
					DCS: tpdu.Dcs8BitData,
					UD:  []byte("ght think"),
				},
			},
			nil,
			twoSegmentMsg,
			nil,
		},
		{
			"single segment 8bit",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					DCS: tpdu.Dcs8BitData,
					UD:  []byte("this is not a very long message"),
				},
			},
			nil,
			[]byte("this is not a very long message"),
			nil,
		},
		{
			"single segment 7bit",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UD: []byte("hello \x03"),
				},
			},
			nil,
			[]byte("hello ¥"),
			nil,
		},
		{
			"single segment 7bit urdu",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 25, Data: []byte{0x0d}},
					},
					UD: []byte("hello \x03"),
				},
			},
			[]sms.ConcatOption{sms.WithCharset(charset.Urdu)},
			[]byte("hello ٻ"),
			nil,
		},
		{
			"single segment 7bit locking urdu",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 25, Data: []byte{0x0d}},
					},
					UD: []byte("hello \x03"),
				},
			},
			[]sms.ConcatOption{sms.WithLockingCharset(charset.Urdu)},
			[]byte("hello ٻ"),
			nil,
		},
		{
			"single segment 7bit shift urdu",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 24, Data: []byte{0x0d}},
					},
					UD: []byte("hello \x1b\x2b"),
				},
			},
			[]sms.ConcatOption{sms.WithShiftCharset(charset.Urdu)},
			[]byte("hello ؏"),
			nil,
		},
		{
			"single segment 7bit urdu unsupported",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 25, Data: []byte{0x0d}},
					},
					UD: []byte("hello \x03"),
				},
			},
			[]sms.ConcatOption{sms.WithCharset(charset.Kannada)},
			[]byte("hello ¥"),
			nil,
		},
		{
			"ucs2 split surrogate",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					DCS: tpdu.DcsUCS2Data,
					UD:  []byte{0xd8, 0x3d, 0xde, 0x01, 0xd8, 0x3d},
				},
				&tpdu.TPDU{
					DCS: tpdu.DcsUCS2Data,
					UD:  []byte{0xde, 0x01, 0xd8, 0x3d, 0xde, 0x01},
				},
			},
			nil,
			[]byte("😁😁😁"),
			nil,
		},
		{
			"ucs2 dangling surrogate",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					DCS: tpdu.DcsUCS2Data,
					UD:  []byte{0xd8, 0x3d, 0xde, 0x01, 0xd8, 0x3d},
				},
			},
			nil,
			nil,
			ucs2.ErrDanglingSurrogate([]byte{0xd8, 0x3d}),
		},
		{
			"ucs2 odd length",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					DCS: tpdu.DcsUCS2Data,
					UD:  []byte{0xd8, 0x3d, 0xde, 0x01, 0xd8},
				},
			},
			nil,
			nil,
			ucs2.ErrInvalidLength,
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out, err := sms.Concatenate(p.in, p.options...)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}

func TestIsCompleteMessage(t *testing.T) {
	patterns := []struct {
		name string
		in   []*tpdu.TPDU
		out  bool
	}{
		{
			"nil",
			nil,
			false,
		},
		{
			"empty",
			[]*tpdu.TPDU{},
			false,
		},
		{
			"single segment",
			[]*tpdu.TPDU{
				&tpdu.TPDU{},
			},
			true,
		},
		{
			"segment count mismatch",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}},
					},
				},
			},
			false,
		},
		{
			"segments mismatch",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}},
					},
				},
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 3, 2}},
					},
				},
			},
			false,
		},
		{
			"concatRef mismatch",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}},
					},
				},
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{4, 2, 2}},
					},
				},
			},
			false,
		},
		{
			"misordered segments",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 2}},
					},
				},
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}},
					},
				},
			},
			false,
		},
		{
			"missing concat",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}},
					},
				},
				&tpdu.TPDU{},
			},
			false,
		},
		{
			"no concat",
			[]*tpdu.TPDU{
				&tpdu.TPDU{},
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 3, Data: []byte{3, 2, 2}},
					},
				},
			},
			false,
		},
		{
			"two segments",
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}},
					},
				},
				&tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 2}},
					},
				},
			},
			true,
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out := sms.IsCompleteMessage(p.in)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}

func TestDecode(t *testing.T) {
	tz1 := time.FixedZone("SCTS", 3600)
	tz8 := time.FixedZone("SCTS", 28800)
	patterns := []struct {
		name    string
		in      []byte
		options []sms.DecodeOption
		out     *tpdu.TPDU
		err     error
	}{
		{
			"empty",
			nil,
			nil,
			nil,
			tpdu.NewDecodeError("tpdu.firstOctet", 0, tpdu.ErrUnderflow),
		},
		{
			"deliver single segment",
			[]byte{
				0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71,
				0x32, 0x20, 0x05, 0x23, 0x08, 0xC8, 0x30, 0x3A, 0x8C, 0x0E,
				0xA3, 0xC3,
			},
			nil,
			&tpdu.TPDU{
				FirstOctet: 0x04,
				OA:         tpdu.Address{TOA: 0x91, Addr: "6391"},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, 5, 17, 23, 2, 50, 0, tz8),
				},
				UD: []byte("Hahahaha"),
			},
			nil,
		},
		{
			"submit single segment",
			[]byte{
				0x31, 0x07, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0xA9, 0x08,
				0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3,
			},
			[]sms.DecodeOption{sms.AsMO},
			&tpdu.TPDU{
				Direction:  1,
				FirstOctet: 0x31,
				MR:         0x07,
				DA:         tpdu.Address{TOA: 0x91, Addr: "6391"},
				VP: tpdu.ValidityPeriod{
					Format:   tpdu.VpfRelative,
					Duration: time.Hour * 72,
				},
				UD: []byte("Hahahaha"),
			},
			nil,
		},
		{
			"submitreport",
			[]byte{
				0x06, 0x18, 0x04, 0x91, 0x36, 0x19, 0x11, 0x10, 0x11, 0x71,
				0x95, 0x51, 0x40, 0x11, 0x10, 0x11, 0x71, 0x95, 0x71, 0x40,
				0x00,
			},
			nil,
			&tpdu.TPDU{
				FirstOctet: 0x06,
				MR:         0x18,
				RA:         tpdu.Address{TOA: 0x91, Addr: "6391"},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2011, 1, 11, 17, 59, 15, 0, tz1),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2011, 1, 11, 17, 59, 17, 0, tz1),
				},
			},
			nil,
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out, err := sms.Decode(p.in, p.options...)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}
