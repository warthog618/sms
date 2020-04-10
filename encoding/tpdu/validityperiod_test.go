// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package tpdu_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/bcd"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestVPEnhancedFormat(t *testing.T) {
	patterns := []struct {
		in  byte
		out tpdu.EnhancedValidityPeriodFormat
	}{
		{0x00, 0x00},
		{0x07, 0x07},
		{0x0f, 0x07},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out := tpdu.EnhancedFormat(p.in)
			assert.Equal(t, p.out, out)
		}
		t.Run(fmt.Sprintf("%02x", p.in), f)
	}
}

func TestVPSetAbsolute(t *testing.T) {
	v := tpdu.ValidityPeriod{}
	soon := tpdu.Timestamp{Time: time.Now().Add(300 * time.Second)}
	v.SetAbsolute(soon)
	if v.Format != tpdu.VpfAbsolute {
		t.Errorf("format is %v, expected %v", v.Format, tpdu.VpfAbsolute)
	}
	if v.Time != soon {
		t.Errorf("time is %v, expected %v", v.Time, soon)
	}
}

func TestVPSetRelative(t *testing.T) {
	v := tpdu.ValidityPeriod{}
	v.SetRelative(123 * time.Second)
	if v.Format != tpdu.VpfRelative {
		t.Errorf("format is %v, expected %v", v.Format, tpdu.VpfRelative)
	}
	if v.Duration != 123*time.Second {
		t.Errorf("duration is %v, expected %v", v.Duration, 123*time.Second)
	}
}

func TestVPSetEnhanced(t *testing.T) {
	v := tpdu.ValidityPeriod{}
	seconds := 123 * time.Second
	efi := 4
	v.SetEnhanced(seconds, byte(efi))
	if v.Format != tpdu.VpfEnhanced {
		t.Errorf("format is %v, expected %v", v.Format, tpdu.VpfRelative)
	}
	if v.Duration != seconds {
		t.Errorf("duration is %v, expected %v", v.Duration, seconds)
	}
	if v.EFI != byte(efi) {
		t.Errorf("efi is %x, expected %x", v.EFI, efi)
	}
}

func TestVPMarshalBinary(t *testing.T) {
	patterns := []struct {
		name string
		in   tpdu.ValidityPeriod
		out  []byte
		err  error
	}{
		{"notpresent", tpdu.ValidityPeriod{}, nil, nil},
		{"absolute",
			tpdu.ValidityPeriod{
				Format: tpdu.VpfAbsolute,
				Time: tpdu.Timestamp{
					Time: time.Date(2017, time.August, 31, 11, 21, 54, 0,
						time.FixedZone("SCTS", 8*3600))},
			},
			[]byte{0x71, 0x80, 0x13, 0x11, 0x12, 0x45, 0x23},
			nil},
		{"relativeMinutes",
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfRelative,
				Duration: 11 * time.Hour},
			[]byte{0x83},
			nil},
		{"relativeHours",
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfRelative,
				Duration: 23 * time.Hour},
			[]byte{0xa5},
			nil},
		{"relativeDays",
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfRelative,
				Duration: 29 * 24 * time.Hour},
			[]byte{0xc3},
			nil},
		{"relativeWeeks",
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfRelative,
				Duration: 62 * 7 * 24 * time.Hour},
			[]byte{0xfe},
			nil},
		{"relativeMax",
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfRelative,
				Duration: 63 * 7 * 24 * time.Hour},
			[]byte{0xff},
			nil},
		{"enhancedNotPresent",
			tpdu.ValidityPeriod{
				Format: tpdu.VpfEnhanced,
				EFI:    byte(tpdu.EvpfNotPresent)},
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			nil},
		{"enhancedRelative5m",
			tpdu.ValidityPeriod{
				Format: tpdu.VpfEnhanced,
				EFI:    byte(tpdu.EvpfRelative)},
			[]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			nil},
		{"enhancedRelative10m",
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfEnhanced,
				EFI:      byte(tpdu.EvpfRelative),
				Duration: 10 * time.Minute},
			[]byte{0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00},
			nil},
		{"enhancedRelativeSeconds",
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfEnhanced,
				EFI:      byte(tpdu.EvpfRelativeSeconds),
				Duration: time.Hour},
			[]byte{0x02, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00},
			nil},
		{"enhancedHHMMSS",
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfEnhanced,
				EFI:      byte(tpdu.EvpfRelativeHHMMSS),
				Duration: 3*time.Hour + 12*time.Minute + 45*time.Second,
			},
			[]byte{0x03, 0x30, 0x21, 0x54, 0x00, 0x00, 0x00},
			nil},
		{"invalid enhanced",
			tpdu.ValidityPeriod{Format: tpdu.VpfEnhanced, EFI: 0xff},
			nil,
			tpdu.EncodeError("fi", tpdu.ErrInvalid)},
	}

	for _, p := range patterns {
		f := func(t *testing.T) {
			b, err := p.in.MarshalBinary()
			if err != p.err {
				t.Errorf("error encoding '%v': %v", p.in, err)
			}
			assert.Equal(t, p.out, b)
		}
		t.Run(p.name, f)
	}
}

func TestVPUnmarshalBinary(t *testing.T) {
	patterns := []struct {
		name      string
		in        []byte
		vpf       tpdu.ValidityPeriodFormat
		readCount int
		out       tpdu.ValidityPeriod
		err       error
	}{
		{"notpresent", nil, tpdu.VpfNotPresent, 0, tpdu.ValidityPeriod{}, nil},
		{"absolute",
			[]byte{0x71, 0x80, 0x13, 0x11, 0x12, 0x45, 0x23},
			tpdu.VpfAbsolute, 7,
			tpdu.ValidityPeriod{
				Format: tpdu.VpfAbsolute,
				Time: tpdu.Timestamp{Time: time.Date(2017, time.August, 31, 11, 21, 54, 0,
					time.FixedZone("SCTS", 8*3600))}},
			nil},
		{"relativeMinutes",
			[]byte{0x83},
			tpdu.VpfRelative, 1,
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfRelative,
				Duration: 11 * time.Hour},
			nil},
		{"relativeHours",
			[]byte{0xa5},
			tpdu.VpfRelative, 1,
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfRelative,
				Duration: 23 * time.Hour},
			nil},
		{"relativeDays",
			[]byte{0xc3},
			tpdu.VpfRelative, 1,
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfRelative,
				Duration: 29 * 24 * time.Hour},
			nil},
		{"relativeWeeks",
			[]byte{0xfe},
			tpdu.VpfRelative, 1,
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfRelative,
				Duration: 62 * 7 * 24 * time.Hour},
			nil},
		{"relativeMax",
			[]byte{0xff},
			tpdu.VpfRelative, 1,
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfRelative,
				Duration: 63 * 7 * 24 * time.Hour},
			nil},
		{"enhancedNotPresent",
			[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			tpdu.VpfEnhanced, 7,
			tpdu.ValidityPeriod{
				Format: tpdu.VpfEnhanced,
				EFI:    byte(tpdu.EvpfNotPresent)},
			nil},
		{"enhancedRelative5m",
			[]byte{0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			tpdu.VpfEnhanced, 7,
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfEnhanced,
				EFI:      byte(tpdu.EvpfRelative),
				Duration: 5 * time.Minute},
			nil},
		{"enhancedRelative10m",
			[]byte{0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00},
			tpdu.VpfEnhanced, 7,
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfEnhanced,
				EFI:      byte(tpdu.EvpfRelative),
				Duration: 10 * time.Minute},
			nil},
		{"enhancedRelativeSeconds",
			[]byte{0x02, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00},
			tpdu.VpfEnhanced, 7,
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfEnhanced,
				EFI:      byte(tpdu.EvpfRelativeSeconds),
				Duration: 255 * time.Second},
			nil},
		{"enhancedHHMMSS",
			[]byte{0x03, 0x30, 0x21, 0x54, 0x00, 0x00, 0x00},
			tpdu.VpfEnhanced, 7,
			tpdu.ValidityPeriod{
				Format:   tpdu.VpfEnhanced,
				EFI:      byte(tpdu.EvpfRelativeHHMMSS),
				Duration: 3*time.Hour + 12*time.Minute + 45*time.Second,
			},
			nil},
		{"underflow relative",
			nil,
			tpdu.VpfRelative, 0,
			tpdu.ValidityPeriod{},
			tpdu.ErrUnderflow},
		{"underflow enhanced",
			[]byte{0x00, 0x01},
			tpdu.VpfEnhanced, 0,
			tpdu.ValidityPeriod{},
			tpdu.ErrUnderflow},
		{"invalid vpf",
			nil,
			0x4, 0,
			tpdu.ValidityPeriod{},
			tpdu.NewDecodeError("vpf", 0, tpdu.ErrInvalid)},
		{"invalid evpf",
			[]byte{0x07, 0x01, 0x2d, 0x54, 0x00, 0x00, 0x00},
			tpdu.VpfEnhanced, 7,
			tpdu.ValidityPeriod{},
			tpdu.NewDecodeError("enhanced", 0, tpdu.ErrInvalid)},
		{"invalid enhancedHHMMSS",
			[]byte{0x03, 0x30, 0x2d, 0x54, 0x00, 0x00, 0x00},
			tpdu.VpfEnhanced, 4,
			tpdu.ValidityPeriod{},
			tpdu.NewDecodeError("enhanced", 1, bcd.ErrInvalidOctet(0x2d))},
		{"nonzero pad enhancedNotPresent",
			[]byte{0x00, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00},
			tpdu.VpfEnhanced, 1,
			tpdu.ValidityPeriod{},
			tpdu.NewDecodeError("enhanced", 1, tpdu.ErrNonZero)},
		{"nonzero pad enhancedRelative",
			[]byte{0x01, 0x00, 0x01, 0x00, 0x00, 0x00, 0x00},
			tpdu.VpfEnhanced, 2,
			tpdu.ValidityPeriod{},
			tpdu.NewDecodeError("enhanced", 2, tpdu.ErrNonZero)},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			s := tpdu.ValidityPeriod{}
			n, err := s.UnmarshalBinary(p.in, p.vpf)
			if err != p.err {
				t.Errorf("error decoding '%v': %v", p.in, err)
			}
			if n != p.readCount {
				t.Errorf("expected to read %d characters, read %d", p.readCount, n)
			}
			assert.Equal(t, p.out, s)
		}
		t.Run(p.name, f)
	}
}

func TestVPFString(t *testing.T) {
	patterns := []struct {
		vpf tpdu.ValidityPeriodFormat
		out string
	}{
		{0, "Not Present"},
		{1, "Enhanced"},
		{2, "Relative"},
		{3, "Absolute"},
		{4, "Unknown"},
	}
	for _, p := range patterns {
		assert.Equal(t, p.out, p.vpf.String())
	}
}

func TestEVPFString(t *testing.T) {
	patterns := []struct {
		evpf tpdu.EnhancedValidityPeriodFormat
		out  string
	}{
		{0, "Not Present"},
		{1, "Relative"},
		{2, "Relative Seconds"},
		{3, "Relative HHMMSS"},
		{4, "Unknown"},
	}
	for _, p := range patterns {
		assert.Equal(t, p.out, p.evpf.String())
	}
}
