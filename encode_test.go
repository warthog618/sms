// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package sms_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms"
	"github.com/warthog618/sms/encoding/gsm7/charset"
	"github.com/warthog618/sms/encoding/tpdu"
)

var twoSegmentMsg = []byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think")

func TestNewEncoder(t *testing.T) {
	e := sms.NewEncoder()
	assert.NotNil(t, e)
}

var patterns = []struct {
	name    string
	msg     []byte
	options []sms.EncoderOption
	out     []tpdu.TPDU
	err     error
}{
	{
		"nil",
		nil,
		nil,
		nil,
		nil,
	},
	{
		"empty",
		[]byte{},
		nil,
		nil,
		nil,
	},
	{
		"single segment",
		[]byte("hello"),
		nil,
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				MR:         1,
				UD:         []byte("hello"),
			},
		},
		nil,
	},
	{
		"single segment grin",
		[]byte("hello 😁"),
		nil,
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				PI:         tpdu.PiDCS, // not relevant for Submit, but set as side-effect
				DCS:        tpdu.DcsUCS2Data,
				MR:         1,
				UD: []byte{
					0x00, 0x68, 0x00, 0x65, 0x00, 0x6c, 0x00, 0x6c, 0x00, 0x6f,
					0x00, 0x20, 0xd8, 0x3d, 0xde, 0x01,
				},
			},
		},
		nil,
	},
	{
		"three grins",
		[]byte("😁😁😁"),
		nil,
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				PI:         tpdu.PiDCS, // not relevant for Submit, but set as side-effect
				DCS:        tpdu.DcsUCS2Data,
				MR:         1,
				UD: []byte{
					0xd8, 0x3d, 0xde, 0x01, 0xd8, 0x3d, 0xde, 0x01, 0xd8, 0x3d,
					0xde, 0x01,
				},
			},
		},
		nil,
	},
	{
		"single segment unused urdu",
		[]byte("hello"),
		[]sms.EncoderOption{sms.WithCharset(charset.Urdu)},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				MR:         1,
				UD:         []byte("hello"),
			},
		},
		nil,
	},
	{
		"single segment with urdu",
		[]byte("hello ت"),
		[]sms.EncoderOption{sms.WithCharset(charset.Urdu)},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41, // UDHI | Submit
				MR:         1,
				PI:         tpdu.PiUDL,
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 25, Data: []byte{0x0d}},
				},
				UD: []byte("hello \x07"),
			},
		},
		nil,
	},
	{
		"single segment with locking urdu",
		[]byte("hello ت"),
		[]sms.EncoderOption{sms.WithLockingCharset(charset.Urdu)},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41, // UDHI | Submit
				MR:         1,
				PI:         tpdu.PiUDL,
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 25, Data: []byte{0x0d}},
				},
				UD: []byte("hello \x07"),
			},
		},
		nil,
	},
	{
		"single segment with shift urdu",
		[]byte("hello ؎"),
		[]sms.EncoderOption{sms.WithShiftCharset(charset.Urdu)},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41, // UDHI | Submit
				MR:         1,
				PI:         tpdu.PiUDL,
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 24, Data: []byte{0x0d}},
				},
				UD: []byte("hello \x1b\x2a"),
			},
		},
		nil,
	},
	{
		"single segment discovered urdu",
		[]byte("hello ت؎"),
		[]sms.EncoderOption{sms.WithAllCharsets},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41, // UDHI | Submit
				MR:         1,
				PI:         tpdu.PiUDL,
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 25, Data: []byte{0x0d}},
					tpdu.InformationElement{ID: 24, Data: []byte{0x0d}},
				},
				UD: []byte("hello \x07\x1b\x2a"),
			},
		},
		nil,
	},
	{
		"single segment UCS2",
		[]byte("hello!"), // this isn't UCS2, but demonstrates it is passed raw, not re-encoded.
		[]sms.EncoderOption{sms.AsUCS2},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				PI:         tpdu.PiDCS, // not relevant for Submit, but set as side-effect
				DCS:        tpdu.DcsUCS2Data,
				MR:         1,
				UD:         []byte("hello!"),
			},
		},
		nil,
	},
	{
		"dcs conflict",
		[]byte("hello 😁"),
		[]sms.EncoderOption{sms.WithTemplateOption(tpdu.DCS(0x80))},
		nil,
		sms.ErrDcsConflict,
	},
	{
		"deliver single segment",
		[]byte("hello"),
		[]sms.EncoderOption{sms.AsDeliver},
		[]tpdu.TPDU{
			tpdu.TPDU{
				MR: 1,
				UD: []byte("hello"),
			},
		},
		nil,
	},
	{
		"number",
		[]byte("hello"),
		[]sms.EncoderOption{sms.To("1234")},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				MR:         1,
				DA:         tpdu.Address{TOA: 0x91, Addr: "1234"},
				UD:         []byte("hello"),
			},
		},
		nil,
	},
	{
		"deliver number",
		[]byte("hello"),
		[]sms.EncoderOption{sms.AsDeliver, sms.From("1234")},
		[]tpdu.TPDU{
			tpdu.TPDU{
				MR: 1,
				OA: tpdu.Address{TOA: 0x91, Addr: "1234"},
				UD: []byte("hello"),
			},
		},
		nil,
	},
	{
		"plus number",
		[]byte("hello"),
		[]sms.EncoderOption{sms.To("+1234")},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				MR:         1,
				DA:         tpdu.Address{TOA: 0x91, Addr: "1234"},
				UD:         []byte("hello"),
			},
		},
		nil,
	},
	{
		"two segment 7bit",
		twoSegmentMsg,
		[]sms.EncoderOption{sms.To("1234")},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41, // Submit | UDHI
				MR:         1,
				PI:         tpdu.PiUDL, // not relevant for Submit, but set as side-effect
				DA:         tpdu.Address{TOA: 0x91, Addr: "1234"},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}},
				},
				UD: []byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you mi"),
			},
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41,
				MR:         2,
				PI:         tpdu.PiUDL, // not relevant for Submit, but set as side-effect
				DA:         tpdu.Address{TOA: 0x91, Addr: "1234"},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}},
				},
				UD: []byte("ght think"),
			},
		},
		nil,
	},
	{
		"two segment 7bit with UDH",
		twoSegmentMsg,
		[]sms.EncoderOption{
			sms.To("1234"),
			sms.WithTemplateOption(
				tpdu.WithUDH(
					tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}},
					},
				)),
		},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41, // Submit | UDHI
				MR:         1,
				PI:         tpdu.PiUDL, // not relevant for Submit, but set as side-effect
				DA:         tpdu.Address{TOA: 0x91, Addr: "1234"},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}},
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}},
				},
				UD: []byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than "),
			},
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41,
				MR:         2,
				PI:         tpdu.PiUDL, // not relevant for Submit, but set as side-effect
				DA:         tpdu.Address{TOA: 0x91, Addr: "1234"},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}},
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}},
				},
				UD: []byte("you might think"),
			},
		},
		nil,
	},
	{
		"two segment 7bit with template",
		twoSegmentMsg,
		[]sms.EncoderOption{
			sms.To("1234"), // overridden by template
			sms.WithTemplate(
				tpdu.TPDU{
					DCS: 0x34,
					DA:  tpdu.Address{TOA: 0x91, Addr: "4321"},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}},
					},
				}),
			sms.AsSubmit,
		},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41, // Submit | UDHI
				MR:         1,
				PI:         tpdu.PiUDL, // not relevant for Submit, but set as side-effect
				DCS:        0x34,
				DA:         tpdu.Address{TOA: 0x91, Addr: "4321"},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}},
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}},
				},
				UD: []byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 charac"),
			},
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41,
				MR:         2,
				PI:         tpdu.PiUDL, // not relevant for Submit, but set as side-effect
				DCS:        0x34,
				DA:         tpdu.Address{TOA: 0x91, Addr: "4321"},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}},
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}},
				},
				UD: []byte("ters is more than you might think"),
			},
		},
		nil,
	},
	{
		"two segment 8bit",
		twoSegmentMsg,
		[]sms.EncoderOption{sms.To("1234"), sms.As8Bit},
		[]tpdu.TPDU{
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41, // Submit | UDHI
				MR:         1,
				DCS:        tpdu.Dcs8BitData,
				PI:         tpdu.PiUDL | tpdu.PiDCS, // not relevant for Submit, but set as side-effect
				DA:         tpdu.Address{TOA: 0x91, Addr: "1234"},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}},
				},
				UD: []byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters "),
			},
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x41,
				MR:         2,
				DCS:        0x04,
				PI:         tpdu.PiUDL | tpdu.PiDCS, // not relevant for Submit, but set as side-effect
				DA:         tpdu.Address{TOA: 0x91, Addr: "1234"},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}},
				},
				UD: []byte("is more than you might think"),
			},
		},
		nil,
	},
	{
		"deliver two segment 7bit",
		twoSegmentMsg,
		[]sms.EncoderOption{sms.AsDeliver, sms.From("1234")},
		[]tpdu.TPDU{
			tpdu.TPDU{
				FirstOctet: 0x40, // UDHI
				MR:         1,
				PI:         tpdu.PiUDL, // not relevant for Deliver, but set as side-effect
				OA:         tpdu.Address{TOA: 0x91, Addr: "1234"},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}},
				},
				UD: []byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you mi"),
			},
			tpdu.TPDU{
				FirstOctet: 0x40,
				MR:         2,
				PI:         tpdu.PiUDL, // not relevant for Deliver, but set as side-effect
				OA:         tpdu.Address{TOA: 0x91, Addr: "1234"},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}},
				},
				UD: []byte("ght think"),
			},
		},
		nil,
	},
}

func TestEncode(t *testing.T) {
	for _, p := range patterns {
		f := func(t *testing.T) {
			out, err := sms.Encode(p.msg, p.options...)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}
func TestEncoderEncode(t *testing.T) {
	for _, p := range patterns {
		f := func(t *testing.T) {
			e := sms.NewEncoder(sms.AsSubmit)
			out, err := e.Encode(p.msg, p.options...)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}

func TestEncoderCounters(t *testing.T) {
	e := sms.NewEncoder()
	msgC, ok := e.MsgCount.(*sms.Counter)
	assert.True(t, ok)
	assert.Equal(t, 0, msgC.Read())
	concatC, ok := e.ConcatRef.(*sms.Counter)
	assert.True(t, ok)
	assert.Equal(t, 0, concatC.Read())

	p, err := e.Encode([]byte("blah"))
	assert.Nil(t, err)
	assert.Equal(t, 1, len(p))
	assert.True(t, ok)
	assert.Equal(t, 1, msgC.Read())
	assert.True(t, ok)
	assert.Equal(t, 0, concatC.Read())

	p, err = e.Encode(twoSegmentMsg)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(p))
	assert.True(t, ok)
	assert.Equal(t, 3, msgC.Read())
	assert.True(t, ok)
	assert.Equal(t, 1, concatC.Read())
}
