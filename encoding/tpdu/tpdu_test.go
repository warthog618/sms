// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package tpdu_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warthog618/sms/encoding/bcd"
	"github.com/warthog618/sms/encoding/semioctet"
	"github.com/warthog618/sms/encoding/tpdu"
)

type BadOption struct {
	err error
}

func (o BadOption) ApplyTPDUOption(t *tpdu.TPDU) error {
	return o.err
}

func TestNew(t *testing.T) {
	s, err := tpdu.New()
	require.Nil(t, err)
	addr := tpdu.NewAddress()
	assert.Equal(t, addr, s.OA)
	assert.Equal(t, addr, s.DA)
	assert.Equal(t, addr, s.RA)
	assert.Equal(t, tpdu.SmsDeliver, s.SmsType())

	s, err = tpdu.New(tpdu.SmsCommand)
	require.Nil(t, err)
	assert.Equal(t, tpdu.SmsCommand, s.SmsType())

	inerr := errors.New("failed TPDU option")
	s, err = tpdu.New(BadOption{inerr})
	assert.Equal(t, inerr, err)
	assert.Nil(t, s)
}

func TestNewDeliver(t *testing.T) {
	s, err := tpdu.NewDeliver()
	require.Nil(t, err)
	addr := tpdu.NewAddress()
	assert.Equal(t, addr, s.OA)
	assert.Equal(t, addr, s.DA)
	assert.Equal(t, addr, s.RA)
	assert.Equal(t, tpdu.SmsDeliver, s.SmsType())

	inerr := errors.New("failed TPDU option")
	s, err = tpdu.NewDeliver(BadOption{inerr})
	assert.Equal(t, inerr, err)
	assert.Nil(t, s)
}

func TestNewSubmit(t *testing.T) {
	s, err := tpdu.NewSubmit()
	require.Nil(t, err)
	addr := tpdu.NewAddress()
	assert.Equal(t, addr, s.OA)
	assert.Equal(t, addr, s.DA)
	assert.Equal(t, addr, s.RA)
	assert.Equal(t, tpdu.SmsSubmit, s.SmsType())

	inerr := errors.New("failed TPDU option")
	s, err = tpdu.NewSubmit(BadOption{inerr})
	assert.Equal(t, inerr, err)
	assert.Nil(t, s)
}

func TestAlphabet(t *testing.T) {
	patterns := []dcsAlphabetPattern{
		{0x00, tpdu.Alpha7Bit, nil},
		{0x04, tpdu.Alpha8Bit, nil},
		{0x08, tpdu.AlphaUCS2, nil},
		{0x0c, tpdu.Alpha7Bit, nil},
		{0x80, tpdu.Alpha7Bit, tpdu.ErrInvalid},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			d := tpdu.TPDU{}
			d.DCS = tpdu.DCS(p.in)
			c, err := d.Alphabet()
			require.Equal(t, p.err, err, p.in)
			assert.Equal(t, p.out, c)
		}
		t.Run(fmt.Sprintf("%02x", p.in), f)
	}
}

func TestConcat(t *testing.T) {
	patterns := []concatTestPattern{
		{
			"empty",
			tpdu.UserDataHeader{},
			0,
			0,
			0,
			false,
		},
		{
			"empty data",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{}},
			},
			0,
			0,
			0,
			false,
		},
		{
			"nil data",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: nil},
			},
			0,
			0,
			0,
			false,
		},
		{
			"concat8",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}},
			},
			3,
			2,
			1,
			true,
		},
		{
			"id 1",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 1, Data: []byte{3, 2, 1}},
			},
			0,
			0,
			0,
			false,
		},
		{
			"concat16",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 8, Data: []byte{4, 3, 2, 1}},
			},
			1027,
			2,
			1,
			true,
		},
		{
			"short concat8",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{2, 1}},
			},
			0,
			0,
			0,
			false,
		},
		{
			"short concat16",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 8, Data: []byte{3, 2, 1}},
			},
			0,
			0,
			0,
			false,
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			s, err := tpdu.New(tpdu.WithUDH(p.udh))
			require.Nil(t, err)
			require.NotNil(t, s)
			segments, seqno, mref, ok := s.ConcatInfo()
			assert.Equal(t, p.ok, ok)
			assert.Equal(t, p.segments, segments)
			assert.Equal(t, p.seqno, seqno)
			assert.Equal(t, p.mref, mref)
		}
		t.Run(p.name, f)
	}
}

func TestMarshalBinary(t *testing.T) {
	patterns := []struct {
		name string
		in   tpdu.TPDU
		out  []byte
		err  error
	}{
		{
			"unsupported SMS type",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x03,
			},
			nil,
			tpdu.ErrUnsupportedSmsType(7),
		},
		{
			"SmsCommand full",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x02,
				PID:        0xab,
				UD:         []byte("a command"),
				MR:         0x42,
				CT:         0x89,
				MN:         0x34,
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			[]byte{
				0x02, 0x42, 0xab, 0x89, 0x34, 0x04, 0x91, 0x36, 0x19, 0x09, 0x61,
				0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64,
			},
			nil,
		},
		{
			"SmsCommand bad da",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x02,
				PID:        0xab,
				UD:         []byte("a command"),
				MR:         0x42,
				CT:         0x89,
				MN:         0x34,
				DA:         tpdu.Address{Addr: "d391", TOA: 0x91},
			},
			nil,
			tpdu.EncodeError("SmsCommand.da.addr", semioctet.ErrInvalidDigit('d')),
		},
		{
			"SmsDeliver haha",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 4,
				UD:         []byte("Hahahaha"),
				OA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600))},
			},
			[]byte{
				0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71, 0x32,
				0x20, 0x05, 0x23, 0x08, 0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3,
			},
			nil,
		},
		{
			"SmsDeliver bad oa",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 4,
				UD:         []byte("Hahahaha"),
				OA:         tpdu.Address{Addr: "d391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			nil,
			tpdu.EncodeError("SmsDeliver.oa.addr", semioctet.ErrInvalidDigit('d')),
		},
		{
			"SmsDeliver bad scts",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 4,
				UD:         []byte("Hahahaha"),
				OA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 24*3600)),
				},
			},
			nil,
			tpdu.EncodeError("SmsDeliver.scts", bcd.ErrInvalidInteger(96)),
		},
		{
			"SmsDeliver bad ud",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 4,
				DCS:        0x80,
				UD:         []byte("Hahahaha"),
				OA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			nil,
			tpdu.EncodeError("SmsDeliver.ud.alphabet", tpdu.ErrInvalid),
		},
		{
			"SmsDeliverReport minimal",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0,
				FCS:        0x12,
			},
			[]byte{0x00, 0x12, 0x00},
			nil,
		},
		{
			"SmsDeliverReport pid",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0,
				PID:        0xab,
				DCS:        0x04,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x01,
			},
			[]byte{
				0x00, 0x12, 0x01, 0xab,
			},
			nil,
		},
		{
			"SmsDeliverReport dcs",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0,
				PID:        0xab,
				DCS:        0x04,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x02,
			},
			[]byte{
				0x00, 0x12, 0x02, 0x04,
			},
			nil,
		},
		{
			"SmsDeliverReport ud",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0,
				PID:        0xab,
				DCS:        0x04,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x04,
			},
			[]byte{
				0x00, 0x12, 0x04, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
			},
			nil,
		},
		{
			"SmsDeliverReport full",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0,
				PID:        0xab,
				DCS:        0x04,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x07,
			},
			[]byte{
				0x00, 0x12, 0x07, 0xab, 0x04, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72,
				0x74,
			},
			nil,
		},
		{
			"SmsDeliverReport bad ud",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0,
				DCS:        0x80,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x06,
			},
			nil,
			tpdu.EncodeError("SmsDeliverReport.ud.alphabet", tpdu.ErrInvalid),
		},
		{
			"SmsStatusReport minimal",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
				ST: 0xab,
			},
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab,
			},
			nil,
		},
		{
			"SmsStatusReport full",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x06,
				PID:        0x89,
				DCS:        0x04,
				UD:         []byte("report"),
				MR:         0x42,
				PI:         0x07,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
				ST: 0xab,
			},
			[]byte{
				0x06, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x07,
				0x89, 0x04, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
			},
			nil,
		},
		{
			"SmsStatusReport pidless",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x06,
				PID:        0x89, // not in PI so wont be encoded
				DCS:        0x04,
				UD:         []byte("report"),
				MR:         0x42,
				PI:         0x06, // no PI set
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
				ST: 0xab,
			},
			[]byte{
				0x06, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x06,
				0x04, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
			},
			nil,
		},
		{
			"SmsStatusReport dcsless",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x06,
				PID:        0x89,
				DCS:        0x04, // not in PI so wont be encoded
				UD:         []byte("report"),
				MR:         0x42,
				PI:         0x05, // no DCS set
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
				ST: 0xab,
			},
			[]byte{
				0x06, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x05,
				0x89, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
			},
			nil,
		},
		{
			"SmsStatusReport bad ra",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				RA:         tpdu.Address{Addr: "63d1", TOA: 0x91},
			},
			nil,
			tpdu.EncodeError("SmsStatusReport.ra.addr", semioctet.ErrInvalidDigit('d')),
		},
		{
			"SmsStatusReport bad scts",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 24*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
				ST: 0xab,
			},
			nil,
			tpdu.EncodeError("SmsStatusReport.scts", bcd.ErrInvalidInteger(96)),
		},
		{
			"SmsStatusReport bad dt",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 24*3600)),
				},
				ST: 0xab,
			},
			nil,
			tpdu.EncodeError("SmsStatusReport.dt", bcd.ErrInvalidInteger(96)),
		},
		{
			"SmsStatusReport bad ud",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				DCS:        0x80,
				UD:         []byte("report"),
				PI:         0x06,
			},
			nil,
			tpdu.EncodeError("SmsStatusReport.ud.alphabet", tpdu.ErrInvalid),
		},
		{
			"SmsSubmit haha",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
				UD:         []byte("Hahahaha"),
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			[]byte{
				0x01, 0x00, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x08, 0xC8, 0x30,
				0x3A, 0x8C, 0x0E, 0xA3, 0xC3,
			},
			nil,
		},
		{
			"SmsSubmit vp",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
				UD:         []byte("Hahahaha"),
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				VP: tpdu.ValidityPeriod{
					Format:   tpdu.VpfRelative,
					Duration: time.Duration(6000000000000),
				},
			},
			[]byte{
				0x01, 0x00, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x13, 0x08, 0xC8,
				0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3,
			},
			nil,
		},
		{
			"SmsSubmit bad da",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
				UD:         []byte("Hahahaha"),
				DA:         tpdu.Address{Addr: "d391", TOA: 0x91},
			},
			nil,
			tpdu.EncodeError("SmsSubmit.da.addr", semioctet.ErrInvalidDigit('d')),
		},
		{
			"SmsSubmit bad vp",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
				UD:         []byte("Hahahaha"),
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				VP:         tpdu.ValidityPeriod{Format: 6},
			},
			nil,
			tpdu.EncodeError("SmsSubmit.vp.vpf", tpdu.ErrInvalid),
		},
		{
			"SmsSubmit bad ud",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
				DCS:        0x80,
				UD:         []byte("Hahahaha"),
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			nil,
			tpdu.EncodeError("SmsSubmit.ud.alphabet", tpdu.ErrInvalid),
		},
		{
			"SmsSubmitReport minimal",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				FCS:        0x12,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			[]byte{0x01, 0x12, 0x00, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23},
			nil,
		},
		{
			"SmsSubmitReport pi",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				PID:        0xab,
				DCS:        0x04,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x01,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			[]byte{
				0x01, 0x12, 0x01, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0xab,
			},
			nil,
		},
		{
			"SmsSubmitReport dcs",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				PID:        0xab,
				DCS:        0x04,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x02,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			[]byte{
				0x01, 0x12, 0x02, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0x04,
			},
			nil,
		},
		{
			"SmsSubmitReport ud",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				PID:        0xab,
				DCS:        0x04,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x04,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			[]byte{
				0x01, 0x12, 0x04, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0x06,
				0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
			},
			nil,
		},
		{
			"SmsSubmitReport full",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				PID:        0xab,
				DCS:        0x04,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x07,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			[]byte{
				0x01, 0x12, 0x07, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0xab,
				0x04, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
			},
			nil,
		},
		{
			"SmsSubmitReport bad scts",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				DCS:        0x80,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x07,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 24*3600))},
			},
			nil,
			tpdu.EncodeError("SmsSubmitReport.scts", bcd.ErrInvalidInteger(96)),
		},
		{
			"SmsSubmitReport bad ud",
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				DCS:        0x80,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x06,
			},
			nil,
			tpdu.EncodeError("SmsSubmitReport.ud.alphabet", tpdu.ErrInvalid),
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

func TestMTI(t *testing.T) {
	b := tpdu.TPDU{}
	m := b.MTI()
	assert.Equal(t, tpdu.MtDeliver, m)
	for _, p := range []tpdu.FirstOctet{0x00, 0xab, 0x00, 0xff} {
		b.FirstOctet = p
		m = b.MTI()
		assert.Equal(t, tpdu.MessageType(p&0x3), m)
	}
}

func TestSetDCS(t *testing.T) {
	b := tpdu.TPDU{}
	assert.False(t, b.PI.DCS())
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.SetDCS(p)
		assert.Equal(t, p, byte(b.DCS))
		assert.True(t, b.PI.DCS())
	}
}

func TestSegment(t *testing.T) {
	patterns := []struct {
		name    string
		in      tpdu.TPDU
		msg     []byte
		options []tpdu.SegmentationOption
		out     []tpdu.TPDU
	}{
		{
			"nil msg",
			tpdu.TPDU{},
			nil,
			nil,
			nil,
		},
		{
			"empty msg",
			tpdu.TPDU{},
			[]byte{},
			nil,
			nil,
		},
		{
			"single segment",
			tpdu.TPDU{},
			[]byte("hello"),
			nil,
			[]tpdu.TPDU{
				tpdu.TPDU{
					UD: []byte("hello"),
				},
			},
		},
		{
			"two segment 7bit",
			tpdu.TPDU{},
			[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think"),
			nil,
			[]tpdu.TPDU{
				tpdu.TPDU{
					FirstOctet: tpdu.FoUDHI,
					PI:         tpdu.PiUDL,
					UDH: []tpdu.InformationElement{
						tpdu.InformationElement{
							ID:   0,
							Data: []byte{1, 2, 1},
						},
					},
					UD: []byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you mi"),
				},
				tpdu.TPDU{
					FirstOctet: tpdu.FoUDHI,
					PI:         tpdu.PiUDL,
					UDH: []tpdu.InformationElement{
						tpdu.InformationElement{
							ID:   0,
							Data: []byte{1, 2, 2},
						},
					},
					UD: []byte("ght think"),
				},
			},
		},
		{
			"three segment 7bit withMR",
			tpdu.TPDU{},
			[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think, but wait, then we also need a really really long message to trigger a three segment concatenation which requires even more characters than I care to count"),
			[]tpdu.SegmentationOption{
				tpdu.WithMR(&counter{10}),
				tpdu.WithConcatRef(&counter{6}),
			},
			[]tpdu.TPDU{
				tpdu.TPDU{
					FirstOctet: tpdu.FoUDHI,
					PI:         tpdu.PiUDL,
					MR:         11,
					UDH: []tpdu.InformationElement{
						tpdu.InformationElement{
							ID:   0,
							Data: []byte{7, 3, 1},
						},
					},
					UD: []byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you mi"),
				},
				tpdu.TPDU{
					FirstOctet: tpdu.FoUDHI,
					PI:         tpdu.PiUDL,
					MR:         12,
					UDH: []tpdu.InformationElement{
						tpdu.InformationElement{
							ID:   0,
							Data: []byte{7, 3, 2},
						},
					},
					UD: []byte("ght think, but wait, then we also need a really really long message to trigger a three segment concatenation which requires even more characters than I c"),
				},
				tpdu.TPDU{
					FirstOctet: tpdu.FoUDHI,
					PI:         tpdu.PiUDL,
					MR:         13,
					UDH: []tpdu.InformationElement{
						tpdu.InformationElement{
							ID:   0,
							Data: []byte{7, 3, 3},
						},
					},
					UD: []byte("are to count"),
				},
			},
		},
		{
			"three segment 7bit with16BitConcatRef",
			tpdu.TPDU{},
			[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think, but wait, then we also need a really really long message to trigger a three segment concatenation which requires even more characters than I care to count"),
			[]tpdu.SegmentationOption{
				tpdu.WithMR(&counter{20}),
				tpdu.WithConcatRef(&counter{0x507}),
				tpdu.With16BitConcatRef,
			},
			[]tpdu.TPDU{
				tpdu.TPDU{
					FirstOctet: tpdu.FoUDHI,
					PI:         tpdu.PiUDL,
					MR:         21,
					UDH: []tpdu.InformationElement{
						tpdu.InformationElement{
							ID:   8,
							Data: []byte{5, 8, 3, 1},
						},
					},
					UD: []byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you m"),
				},
				tpdu.TPDU{
					FirstOctet: tpdu.FoUDHI,
					PI:         tpdu.PiUDL,
					MR:         22,
					UDH: []tpdu.InformationElement{
						tpdu.InformationElement{
							ID:   8,
							Data: []byte{5, 8, 3, 2},
						},
					},
					UD: []byte("ight think, but wait, then we also need a really really long message to trigger a three segment concatenation which requires even more characters than I"),
				},
				tpdu.TPDU{
					FirstOctet: tpdu.FoUDHI,
					PI:         tpdu.PiUDL,
					MR:         23,
					UDH: []tpdu.InformationElement{
						tpdu.InformationElement{
							ID:   8,
							Data: []byte{5, 8, 3, 3},
						},
					},
					UD: []byte(" care to count"),
				},
			},
		},
		{
			"8bit",
			tpdu.TPDU{
				DCS: tpdu.Dcs8BitData,
			},
			[]byte("hello"),
			nil,
			[]tpdu.TPDU{
				tpdu.TPDU{
					DCS: tpdu.Dcs8BitData,
					UD:  []byte("hello"),
				},
			},
		},
		{
			"ucs2",
			tpdu.TPDU{
				DCS: tpdu.DcsUCS2Data,
			},
			[]byte("hello"),
			nil,
			[]tpdu.TPDU{
				tpdu.TPDU{
					DCS: tpdu.DcsUCS2Data,
					UD:  []byte("hello"),
				},
			},
		},
		{
			"single segment withMR",
			tpdu.TPDU{},
			[]byte("hello"),
			[]tpdu.SegmentationOption{
				tpdu.WithMR(&counter{42}),
			},
			[]tpdu.TPDU{
				tpdu.TPDU{
					MR: 43,
					UD: []byte("hello"),
				},
			},
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out := p.in.Segment(p.msg, p.options...)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}

func TestSetPID(t *testing.T) {
	b := tpdu.TPDU{}
	assert.Zero(t, b.PI)
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.SetPID(p)
		assert.Equal(t, p, b.PID)
		assert.True(t, b.PI.PID())
	}
}

func TestSetSmsType(t *testing.T) {
	patterns := []struct {
		in  tpdu.SmsType
		err error
	}{
		{tpdu.SmsDeliver, nil},
		{tpdu.SmsSubmit, nil},
		{7, tpdu.ErrInvalid},
	}
	for _, p := range patterns {
		s, err := tpdu.New()
		require.Nil(t, err)
		err = s.SetSmsType(p.in)
		assert.Equal(t, p.err, err)
		if err == nil {
			assert.Equal(t, p.in, s.SmsType())
		} else {
			assert.Equal(t, tpdu.SmsType(0), s.SmsType())
		}
	}
}

func TestSetUD(t *testing.T) {
	b := tpdu.TPDU{}
	assert.Zero(t, b.PI)
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.SetUD([]byte{p})
		assert.Equal(t, []byte{p}, []byte(b.UD))
		assert.True(t, b.PI.UDL())
	}
	// Reset
	b.SetUD(nil)
	assert.Nil(t, b.UD)
	assert.False(t, b.PI.UDL())
}

func TestSetUDH(t *testing.T) {
	// also tests tpdu.TPDU.UDH
	b := tpdu.TPDU{}
	udh := b.UDH
	assert.Zero(t, len(udh))
	for _, p := range []tpdu.UserDataHeader{
		nil,
		{
			tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
		},
		{
			tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
			tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
		},
		nil,
	} {
		b.SetUDH(p)
		udh = b.UDH
		assert.Equal(t, len(p) != 0, b.UDHI())
		assert.Equal(t, udh, p)
	}
}

func TestSetValidityPeriod(t *testing.T) {
	// also tests Submit.VP
	s := tpdu.TPDU{}
	vp := s.VP
	assert.Equal(t, tpdu.VpfNotPresent, vp.Format)
	pvp := tpdu.ValidityPeriod{}
	pvp.SetRelative(time.Duration(100000000))
	for _, p := range []struct {
		vp tpdu.ValidityPeriod
		fo tpdu.FirstOctet
	}{
		{tpdu.ValidityPeriod{}, 0x00},
		{pvp, 0x10},
		{tpdu.ValidityPeriod{}, 0x00},
	} {
		s.SetVP(p.vp)
		vp = s.VP
		assert.Equal(t, p.fo, s.FirstOctet)
		assert.Equal(t, vp, p.vp)
	}
}

func TestSmsType(t *testing.T) {
	patterns := []struct {
		dirn tpdu.Direction
		mt   tpdu.MessageType
		st   tpdu.SmsType
	}{
		{tpdu.MT, tpdu.MtDeliver, tpdu.SmsDeliver},
		{tpdu.MT, tpdu.MtSubmit, tpdu.SmsSubmitReport},
		{tpdu.MT, tpdu.MtCommand, tpdu.SmsStatusReport},
		{tpdu.MO, tpdu.MtDeliver, tpdu.SmsDeliverReport},
		{tpdu.MO, tpdu.MtSubmit, tpdu.SmsSubmit},
		{tpdu.MO, tpdu.MtCommand, tpdu.SmsCommand},
	}
	s := tpdu.TPDU{}
	for _, p := range patterns {
		s.Direction = p.dirn
		s.FirstOctet = tpdu.FirstOctet(p.mt)
		assert.Equal(t, p.st, s.SmsType())
	}
}

func TestSmsTypeString(t *testing.T) {
	patterns := []struct {
		st  tpdu.SmsType
		out string
	}{
		{0, "SmsDeliver"},
		{1, "SmsDeliverReport"},
		{2, "SmsSubmitReport"},
		{3, "SmsSubmit"},
		{4, "SmsStatusReport"},
		{5, "SmsCommand"},
		{6, "Unknown"},
	}
	for _, p := range patterns {
		assert.Equal(t, p.out, p.st.String())
	}
}

func TestUDBlockSize(t *testing.T) {
	patterns := []struct {
		name string
		pdu  tpdu.TPDU
		out  int
	}{
		{
			"command 7bit",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtCommand),
			},
			166,
		},
		{
			"command 8bit",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtCommand),
				DCS:        0xf4,
			},
			146,
		},
		{
			"deliver 7bit",
			tpdu.TPDU{},
			160,
		},
		{
			"deliver UDH 7bit",
			tpdu.TPDU{
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
				},
			},
			153,
		},
		{
			"deliver 8bit",
			tpdu.TPDU{
				DCS: 0xf4,
			},
			140,
		},
		{
			"deliver UDH 8bit",
			tpdu.TPDU{
				DCS: 0xf4,
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
				},
			},
			134,
		},
		{
			"deliver UCS2",
			tpdu.TPDU{
				DCS: 0xe0,
			},
			140,
		},
		{
			"deliverreport RP-ACK 7bit",
			tpdu.TPDU{
				Direction: tpdu.MO,
			},
			181,
		},
		{
			"deliverreport RP-ERROR 7bit",
			tpdu.TPDU{
				Direction: tpdu.MO,
				FCS:       0x90,
			},
			180,
		},
		{
			"deliverreport RP-ACK 8bit",
			tpdu.TPDU{
				Direction: tpdu.MO,
				DCS:       0xf4,
			},
			159,
		},
		{
			"deliverreport RP-ERROR 8bit",
			tpdu.TPDU{
				Direction: tpdu.MO,
				DCS:       0xf4,
				FCS:       0x90,
			},
			158,
		},
		{
			"statusreport 7bit",
			tpdu.TPDU{
				FirstOctet: tpdu.FirstOctet(tpdu.MtCommand),
			},
			149,
		},
		{
			"statusreport 8bit",
			tpdu.TPDU{
				FirstOctet: tpdu.FirstOctet(tpdu.MtCommand),
				DCS:        0xf4,
			},
			131,
		},
		{
			"statusreport UCS2",
			tpdu.TPDU{
				FirstOctet: tpdu.FirstOctet(tpdu.MtCommand),
				DCS:        0xe0,
			},
			130,
		},
		{
			"submit 7bit",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
			},
			160,
		},
		{
			"submit 8bit",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				DCS:        0xf4,
			},
			140,
		},
		{
			"submit UCS2",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				DCS:        0xe0,
			},
			140,
		},
		{
			"submitreport RP-ACK 7bit",
			tpdu.TPDU{
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
			},
			173,
		},
		{
			"submitreport RP-ERROR 7bit",
			tpdu.TPDU{
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				FCS:        0x90,
			},
			172,
		},
		{
			"submitreport RP-ACK 8bit",
			tpdu.TPDU{
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				DCS:        0xf4,
			},
			152,
		},
		{
			"submitreport RP-ERROR 8bit",
			tpdu.TPDU{
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				DCS:        0xf4,
				FCS:        0x90,
			},
			151,
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			assert.Equal(t, p.out, p.pdu.UDBlockSize())
		}
		t.Run(p.name, f)
	}
}

func TestUDHI(t *testing.T) {
	// also tests tpdu.TPDU.SetUDH
	b := tpdu.TPDU{}
	assert.False(t, b.UDHI())
	for _, p := range []tpdu.UserDataHeader{
		nil,
		{
			tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
		},
		{
			tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
			tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
		},
		nil,
	} {
		b.SetUDH(p)
		assert.Equal(t, len(p) != 0, b.UDHI())
	}
}

func TestUnmarshalBinary(t *testing.T) {
	patterns := []struct {
		name string
		in   []byte
		dirn tpdu.Direction
		out  tpdu.TPDU
		err  error
	}{
		{
			"underflow fo",
			[]byte{},
			tpdu.MT,
			tpdu.TPDU{},
			tpdu.DecodeError("tpdu.firstOctet", 0, tpdu.ErrUnderflow),
		},
		{
			"unsupported SMS type",
			[]byte{0x03},
			tpdu.MT,
			tpdu.TPDU{},
			tpdu.DecodeError("tpdu.firstOctet", 0, tpdu.ErrUnsupportedSmsType(6)),
		},
		{
			"SmsCommand",
			[]byte{
				0x02, 0x42, 0xab, 0x89, 0x34, 0x04, 0x91, 0x36, 0x19, 0x09, 0x61,
				0x20, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x02,
				DCS:        0x04,
				PID:        0xab,
				UD:         []byte("a command"),
				MR:         0x42,
				CT:         0x89,
				MN:         0x34,
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			nil,
		},
		{
			"SmsCommand underflow mr",
			[]byte{0x02},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x02,
			},
			tpdu.DecodeError("SmsCommand.mr", 1, tpdu.ErrUnderflow),
		},
		{
			"SmsCommand underflow pid",
			[]byte{0x02, 0x42},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x02,
				MR:         0x42,
			},
			tpdu.DecodeError("SmsCommand.pid", 2, tpdu.ErrUnderflow),
		},
		{
			"SmsCommand underflow ct",
			[]byte{0x02, 0x42, 0xab},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x02,
				PID:        0xab,
				MR:         0x42,
			},
			tpdu.DecodeError("SmsCommand.ct", 3, tpdu.ErrUnderflow),
		},
		{
			"SmsCommand underflow mn",
			[]byte{0x02, 0x42, 0xab, 0x89},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x02,
				PID:        0xab,
				MR:         0x42,
				CT:         0x89,
			},
			tpdu.DecodeError("SmsCommand.mn", 4, tpdu.ErrUnderflow),
		},
		{
			"SmsCommand underflow da",
			[]byte{0x02, 0x42, 0xab, 0x89, 0x34, 0x04, 0x91, 0x36, 0xF9, 0x09},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x02,
				PID:        0xab,
				MR:         0x42,
				CT:         0x89,
				MN:         0x34,
			},
			tpdu.DecodeError("SmsCommand.da.addr", 7, tpdu.ErrUnderflow),
		},
		{
			"SmsCommand underflow ud",
			[]byte{0x02, 0x42, 0xab, 0x89, 0x34, 0x04, 0x91, 0x36, 0x19},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x02,
				PID:        0xab,
				DCS:        0x04,
				MR:         0x42,
				CT:         0x89,
				MN:         0x34,
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsCommand.ud.udl", 9, tpdu.ErrUnderflow),
		},
		{
			"SmsDeliver haha",
			[]byte{
				0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71, 0x32,
				0x20, 0x05, 0x23, 0x08, 0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x04,
				UD:         []byte("Hahahaha"),
				OA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			nil,
		},
		{
			"SmsDeliver underflow oa",
			[]byte{0x04, 0x04, 0x91, 0x36, 0xF9, 0x00, 0x00},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x04,
			},
			tpdu.DecodeError("SmsDeliver.oa.addr", 3, tpdu.ErrUnderflow),
		},
		{
			"SmsDeliver underflow pid",
			[]byte{0x04, 0x04, 0x91, 0x36, 0x19},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x04,
				OA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsDeliver.pid", 5, tpdu.ErrUnderflow),
		},
		{
			"SmsDeliver underflow dcs",
			[]byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x04,
				OA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsDeliver.dcs", 6, tpdu.ErrUnderflow),
		},
		{
			"SmsDeliver underflow scts",
			[]byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x04,
				OA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsDeliver.scts", 7, tpdu.ErrUnderflow),
		},
		{
			"SmsDeliver bad scts",
			[]byte{
				0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0xf1, 0x32,
				0x20, 0x05, 0x23,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x04,
				OA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsDeliver.scts", 7, bcd.ErrInvalidOctet(0xf1)),
		},
		{
			"SmsDeliver underflow ud",
			[]byte{
				0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71, 0x32,
				0x20, 0x05, 0x23, 0x08,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x04,
				OA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			tpdu.DecodeError("SmsDeliver.ud.sm", 15, tpdu.ErrUnderflow),
		},
		{
			"SmsDeliverReport minimal",
			[]byte{0x00, 0x12, 0x00},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x00,
				FCS:        0x12,
			},
			nil,
		},
		{
			"SmsDeliverReport pid",
			[]byte{0x00, 0x12, 0x01, 0xab},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x00,
				PID:        0xab,
				FCS:        0x12,
				PI:         0x01,
			},
			nil,
		},
		{
			"SmsDeliverReport dcs",
			[]byte{0x00, 0x12, 0x02, 0x04},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x00,
				DCS:        0x04,
				FCS:        0x12,
				PI:         0x02,
			},
			nil,
		},
		{
			"SmsDeliverReport ud",
			[]byte{
				0x00, 0x12, 0x06, 0x04, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
			},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x00,
				DCS:        0x04,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x06,
			},
			nil,
		},
		{
			"SmsDeliverReport underflow fcs",
			[]byte{0x00},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x00,
			},
			tpdu.DecodeError("SmsDeliverReport.fcs", 1, tpdu.ErrUnderflow),
		},
		{
			"SmsDeliverReport underflow pi",
			[]byte{0x00, 0x12},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x00,
				FCS:        0x12,
			},
			tpdu.DecodeError("SmsDeliverReport.pi", 2, tpdu.ErrUnderflow),
		},
		{
			"SmsDeliverReport underflow pid",
			[]byte{0x00, 0x12, 0x01},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x00,
				FCS:        0x12,
				PI:         0x01,
			},
			tpdu.DecodeError("SmsDeliverReport.pid", 3, tpdu.ErrUnderflow),
		},
		{
			"SmsDeliverReport underflow dcs",
			[]byte{0x00, 0x12, 0x02},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x00,
				FCS:        0x12,
				PI:         0x02,
			},
			tpdu.DecodeError("SmsDeliverReport.dcs", 3, tpdu.ErrUnderflow),
		},
		{
			"SmsDeliverReport underflow ud",
			[]byte{0x00, 0x12, 0x04},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x00,
				FCS:        0x12,
				PI:         0x04,
			},
			tpdu.DecodeError("SmsDeliverReport.ud.udl", 3, tpdu.ErrUnderflow),
		},
		{
			"SmsStatusReport minimal",
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
				ST: 0xab,
			},
			nil,
		},
		{
			"SmsStatusReport pid",
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x01,
				0x89,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				PID:        0x89,
				MR:         0x42,
				PI:         0x01,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
				ST: 0xab,
			},
			nil,
		},
		{
			"SmsStatusReport dcs",
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x02,
				0x04,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				DCS:        0x04,
				MR:         0x42,
				PI:         0x02,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
				ST: 0xab,
			},
			nil,
		},
		{
			"SmsStatusReport ud",
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x06,
				0x04, 0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				DCS:        0x04,
				UD:         []byte("report"),
				MR:         0x42,
				PI:         0x06,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
				ST: 0xab,
			},
			nil,
		},
		{
			"SmsStatusReport underflow mr",
			[]byte{0x02},
			tpdu.MT,
			tpdu.TPDU{
				FirstOctet: 0x02,
			},
			tpdu.DecodeError("SmsStatusReport.mr", 1, tpdu.ErrUnderflow),
		},
		{
			"SmsStatusReport underflow ra",
			[]byte{0x02, 0x42},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
			},
			tpdu.DecodeError("SmsStatusReport.ra.addr", 2, tpdu.ErrUnderflow),
		},
		{
			"SmsStatusReport underflow scts",
			[]byte{0x02, 0x42, 0x04, 0x91, 0x36, 0x19},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsStatusReport.scts", 6, tpdu.ErrUnderflow),
		},
		{
			"SmsStatusReport underflow dt",
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			tpdu.DecodeError("SmsStatusReport.dt", 13, tpdu.ErrUnderflow),
		},
		{
			"SmsStatusReport bad scts",
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0xf1, 0x32, 0x20,
				0x05, 0x23,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsStatusReport.scts", 6, bcd.ErrInvalidOctet(0xf1)),
		},
		{
			"SmsStatusReport bad dt",
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0xc1, 0x32, 0x20, 0x05, 0x42,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			tpdu.DecodeError("SmsStatusReport.dt", 13, bcd.ErrInvalidOctet(0xc1)),
		},
		{
			"SmsStatusReport underflow st",
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
			},
			tpdu.DecodeError("SmsStatusReport.st", 20, tpdu.ErrUnderflow),
		},
		{
			"SmsStatusReport underflow pid",
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x01,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				PI:         0x01,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
					time.FixedZone("SCTS", 8*3600))},
				DT: tpdu.Timestamp{Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
					time.FixedZone("SCTS", 6*3600))},
				ST: 0xab,
			},
			tpdu.DecodeError("SmsStatusReport.pid", 22, tpdu.ErrUnderflow),
		},
		{
			"SmsStatusReport underflow dcs",
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x02,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				PI:         0x02,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
				ST: 0xab,
			},
			tpdu.DecodeError("SmsStatusReport.dcs", 22, tpdu.ErrUnderflow),
		},
		{
			"SmsStatusReport underflow ud",
			[]byte{
				0x02, 0x42, 0x04, 0x91, 0x36, 0x19, 0x51, 0x50, 0x71, 0x32, 0x20,
				0x05, 0x23, 0x51, 0x40, 0x81, 0x32, 0x20, 0x05, 0x42, 0xab, 0x04,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x02,
				MR:         0x42,
				PI:         0x04,
				RA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
				DT: tpdu.Timestamp{
					Time: time.Date(2015, time.April, 18, 23, 02, 50, 0,
						time.FixedZone("SCTS", 6*3600)),
				},
				ST: 0xab,
			},
			tpdu.DecodeError("SmsStatusReport.ud.udl", 22, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmit haha",
			[]byte{
				0x01, 0x23, 0x04, 0x91, 0x36, 0x19, 0x34, 0x00, 0x08, 0xC8, 0x30,
				0x3A, 0x8C, 0x0E, 0xA3, 0xC3,
			},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
				PID:        0x34,
				UD:         []byte("Hahahaha"),
				MR:         0x23,
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			nil},
		{
			"SmsSubmit vp",
			[]byte{
				0x11, 0x23, 0x04, 0x91, 0x36, 0x19, 0x34, 0x00, 0x45, 0x08, 0xC8,
				0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3,
			},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x11,
				PID:        0x34,
				UD:         []byte("Hahahaha"),
				MR:         0x23,
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
				VP: tpdu.ValidityPeriod{
					Format:   tpdu.VpfRelative,
					Duration: time.Duration(60 * 350 * 1000000000),
				},
			},
			nil,
		},
		{
			"SmsSubmit underflow mr",
			[]byte{0x01},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
			},
			tpdu.DecodeError("SmsSubmit.mr", 1, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmit underflow da",
			[]byte{0x01, 0x23},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
				MR:         0x23,
			},
			tpdu.DecodeError("SmsSubmit.da.addr", 2, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmit bad da",
			[]byte{0x01, 0x23, 0x04, 0x91, 0x36, 0xF9, 0x00},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
				MR:         0x23,
			},
			tpdu.DecodeError("SmsSubmit.da.addr", 4, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmit underflow pid",
			[]byte{0x01, 0x23, 0x04, 0x91, 0x36, 0x19},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
				MR:         0x23,
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsSubmit.pid", 6, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmit underflow dcs",
			[]byte{0x01, 0x23, 0x04, 0x91, 0x36, 0x19, 0x00},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
				MR:         0x23,
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsSubmit.dcs", 7, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmit underflow vp",
			[]byte{0x11, 0x23, 0x04, 0x91, 0x36, 0x19, 0x34, 0x00},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x11,
				PID:        0x34,
				MR:         0x23,
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsSubmit.vp", 8, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmit bad vp",
			[]byte{
				0x19, 0x23, 0x04, 0x91, 0x36, 0x19, 0x34, 0x00, 0x45, 0x08, 0xC8,
				0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3,
			},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x19,
				PID:        0x34,
				MR:         0x23,
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsSubmit.vp", 8, bcd.ErrInvalidOctet(0xc8)),
		},
		{
			"SmsSubmit underflow ud",
			[]byte{
				0x01, 0x23, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71,
				0x32, 0x20, 0x05, 0x23, 0x08,
			},
			tpdu.MO,
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: 0x01,
				MR:         0x23,
				DA:         tpdu.Address{Addr: "6391", TOA: 0x91},
			},
			tpdu.DecodeError("SmsSubmit.ud.sm", 9, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmitReport minimal",
			[]byte{0x01, 0x12, 0x00, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				FCS:        0x12,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			nil,
		},
		{
			"SmsSubmitReport pid",
			[]byte{0x01, 0x12, 0x01, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0xab},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				PID:        0xab,
				FCS:        0x12,
				PI:         tpdu.PiPID,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			nil,
		},
		{
			"SmsSubmitReport dcs",
			[]byte{0x01, 0x12, 0x02, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0x04},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				DCS:        0x04,
				FCS:        0x12,
				PI:         tpdu.PiDCS,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			nil,
		},
		{
			"SmsSubmitReport ud",
			[]byte{
				0x01, 0x12, 0x06, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0x04,
				0x06, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74,
			},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				DCS:        0x04,
				UD:         []byte("report"),
				FCS:        0x12,
				PI:         0x06,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			nil,
		},
		{
			"SmsSubmitReport underflow fcs",
			[]byte{0x01},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
			},
			tpdu.DecodeError("SmsSubmitReport.fcs", 1, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmitReport underflow pi",
			[]byte{0x01, 0x12},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				FCS:        0x12,
			},
			tpdu.DecodeError("SmsSubmitReport.pi", 2, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmitReport underflow scts",
			[]byte{0x01, 0x12, 0x00},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				FCS:        0x12,
			},
			tpdu.DecodeError("SmsSubmitReport.scts", 3, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmitReport bad scts",
			[]byte{0x01, 0x12, 0x00, 0x51, 0x50, 0xf1, 0x32, 0x20, 0x05, 0x23},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				FCS:        0x12,
			},
			tpdu.DecodeError("SmsSubmitReport.scts", 3, bcd.ErrInvalidOctet(0xf1)),
		},
		{
			"SmsSubmitReport underflow pid",
			[]byte{0x01, 0x12, 0x01, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				FCS:        0x12,
				PI:         0x01,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600))},
			},
			tpdu.DecodeError("SmsSubmitReport.pid", 10, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmitReport underflow dcs",
			[]byte{0x01, 0x12, 0x02, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				FCS:        0x12,
				PI:         0x02,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600))},
			},
			tpdu.DecodeError("SmsSubmitReport.dcs", 10, tpdu.ErrUnderflow),
		},
		{
			"SmsSubmitReport underflow ud",
			[]byte{0x01, 0x12, 0x06, 0x51, 0x50, 0x71, 0x32, 0x20, 0x05, 0x23, 0x04},
			tpdu.MT,
			tpdu.TPDU{
				Direction:  tpdu.MT,
				FirstOctet: 0x01,
				DCS:        0x04,
				FCS:        0x12,
				PI:         0x06,
				SCTS: tpdu.Timestamp{
					Time: time.Date(2015, time.May, 17, 23, 02, 50, 0,
						time.FixedZone("SCTS", 8*3600)),
				},
			},
			tpdu.DecodeError("SmsSubmitReport.ud.udl", 11, tpdu.ErrUnderflow),
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			d := tpdu.TPDU{Direction: p.dirn}
			err := d.UnmarshalBinary(p.in)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, d)
		}
		t.Run(p.name, f)
	}
}

// counter is an implementation of the tpdu.Counter interface.
//
// It is never used in a multi-threaded setting and so is not MT safe.
type counter struct {
	c int
}

// Count increments and returns the counter.
func (c *counter) Count() int {
	c.c++
	return c.c
}
