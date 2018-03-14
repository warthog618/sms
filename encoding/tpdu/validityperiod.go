// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"time"

	"github.com/warthog618/sms/encoding/bcd"
)

// ValidityPeriod represents the validity period as defined in 3GPP TS 34.040 Section 9.2.3.12.
type ValidityPeriod struct {
	Format   ValidityPeriodFormat
	Time     Timestamp     // for VpfAbsolute
	Duration time.Duration // for VpfRelative and VpfEnhanced
	Efi      byte          // enhanced functionality indicator - first octet of enhanced format
}

// SetAbsolute seth the validity period to an absolute time.
func (v *ValidityPeriod) SetAbsolute(t Timestamp) {
	v.Format = VpfAbsolute
	v.Duration = 0
	v.Time = t
	v.Efi = 0
}

// SetRelative sets the validity period to a relative time.
func (v *ValidityPeriod) SetRelative(d time.Duration) {
	v.Format = VpfRelative
	v.Duration = d
	v.Time = Timestamp{}
	v.Efi = 0
}

// SetEnhanced sets the validity period to an enhnaced format as determined
// from the functionality identifier (efi).
func (v *ValidityPeriod) SetEnhanced(d time.Duration, efi byte) {
	v.Format = VpfEnhanced
	v.Duration = d
	v.Time = Timestamp{}
	v.Efi = efi
}

// MarshalBinary marshals a ValidityPeriod.
func (v *ValidityPeriod) MarshalBinary() ([]byte, error) {
	switch v.Format {
	case VpfAbsolute:
		return v.Time.MarshalBinary()
	case VpfEnhanced:
		evpf := EnhancedValidityPeriodFormat(v.Efi & 0x7)
		if evpf > EvpfRelativeHHMMSS {
			return nil, EncodeError("fi", ErrInvalid)
		}
		dst := make([]byte, 7)
		dst[0] = v.Efi
		switch evpf {
		case EvpfRelative:
			dst[1] = durationToRelative(v.Duration)
		case EvpfRelativeSeconds:
			secs := v.Duration / time.Second
			if secs > 255 {
				secs = 255
			}
			dst[1] = byte(secs)
		case EvpfRelativeHHMMSS:
			f := []int{int(v.Duration.Hours()) % 100, int(v.Duration.Minutes()) % 60, int(v.Duration.Seconds()) % 60}
			for i, tf := range f {
				t, err := bcd.Encode(tf)
				// this should never trip, as the encoded values should always be valid, but just in case...
				if err != nil {
					return nil, EncodeError("enhanced", err)
				}
				dst[i+1] = t
			}
		}
		return dst, nil
	case VpfRelative:
		t := durationToRelative(v.Duration)
		return []byte{t}, nil
	case VpfNotPresent:
		return nil, nil
	}
	return nil, EncodeError("vpf", ErrInvalid)
}

// UnmarshalBinary unmarshals a ValidityPeriod stored in the given format.
// Returns the number of bytes read from the src, and any error detected
// during the unmarshalling.
func (v *ValidityPeriod) UnmarshalBinary(src []byte, vpf ValidityPeriodFormat) (int, error) {
	v.Format = VpfNotPresent
	switch vpf {
	case VpfAbsolute:
		t := Timestamp{}
		err := t.UnmarshalBinary(src)
		if err == nil {
			v.Time = t
			v.Format = vpf
		}
		return 7, err
	case VpfEnhanced:
		if len(src) < 7 {
			return 0, ErrUnderflow
		}
		efi := src[0]
		evpf := EnhancedValidityPeriodFormat(efi & 0x7)
		used := 0
		d := time.Duration(0)
		switch evpf {
		case EvpfNotPresent:
		case EvpfRelative:
			d = relativeToDuration(src[1])
			used = 1
		case EvpfRelativeSeconds:
			d = time.Second * time.Duration(src[1])
			used = 1
		case EvpfRelativeHHMMSS:
			i := make([]int, 3)
			var err error
			for idx := 0; idx < 3; idx++ {
				i[idx], err = bcd.Decode(src[idx+1])
				if err != nil {
					return 4, DecodeError("enhanced", 1, err)
				}
			}
			d = time.Duration(i[0])*time.Hour + time.Duration(i[1])*time.Minute + time.Duration(i[2])*time.Second
			used = 3
		default:
			return 7, DecodeError("enhanced", 0, ErrInvalid)
		}
		for i := used + 1; i < 7; i++ {
			if src[i] != 0 {
				return used + 1, DecodeError("enhanced", i, ErrNonZero)
			}
		}
		v.Efi = efi
		v.Duration = d
		v.Format = vpf
		return 7, nil
	case VpfRelative:
		if len(src) < 1 {
			return 0, ErrUnderflow
		}
		v.Duration = relativeToDuration(src[0])
		v.Format = vpf
		return 1, nil
	case VpfNotPresent:
		return 0, nil
	}
	return 0, DecodeError("vpf", 0, ErrInvalid)
}

// ValidityPeriodFormat identifies the format of the ValidityPeriod when encoded to binary.
type ValidityPeriodFormat byte

const (
	// VpfNotPresent indicates no VP is present.
	VpfNotPresent ValidityPeriodFormat = iota
	// VpfEnhanced indicates the VP is stored in enhanced format as per 3GPP TS 23.038 Section 9.2.3.12.3.
	VpfEnhanced
	// VpfRelative indicates the VP is stored in relative format as per 3GPP TS 23.038 Section 9.2.3.12.1.
	VpfRelative
	// VpfAbsolute indicates the VP is stored in absolute format as per 3GPP TS 23.038 Section 9.2.3.12.2.
	// The absolute format is the same format as the SCTS.
	VpfAbsolute
)

// EnhancedValidityPeriodFormat identifies the subformat of the ValidityPeriod
// when encoded to binary in enhanced format, as per 3GPP TS 23.038 Section 9.2.3.12.3
type EnhancedValidityPeriodFormat byte

const (
	// EvpfNotPresent indicates no VP is present.
	EvpfNotPresent EnhancedValidityPeriodFormat = iota
	// EvpfRelative indicates the VP is stored in relative format as per 3GPP TS 23.038 Section 9.2.3.12.1.
	EvpfRelative
	// EvpfRelativeSeconds indicates the VP is stored in relative format as an
	// integer number of seconds, from 0 to 255.
	EvpfRelativeSeconds
	// EvpfRelativeHHMMSS indicates the VP is stored in relative format as a period of
	// hours, minutes and seconds in semioctet format as per SCTS time.
	EvpfRelativeHHMMSS
	// All other values currently reserved.
)

func durationToRelative(d time.Duration) byte {
	switch {
	case d < time.Hour*12:
		t := byte(d / (time.Minute * 5))
		if t > 1 {
			t--
		}
		return t
	case d < time.Hour*24:
		return 119 + byte(d/(time.Minute*30))
	case d < time.Hour*24*30:
		return 166 + byte(d/(time.Hour*24))
	case d < time.Hour*24*7*63:
		return 192 + byte(d/(time.Hour*24*7))
	default:
		return 255
	}
}

func relativeToDuration(t byte) time.Duration {
	switch {
	case t < 144:
		return time.Minute * 5 * time.Duration(t+1)
	case t < 168:
		return time.Minute * 30 * time.Duration(t-119)
	case t < 197:
		return time.Hour * 24 * time.Duration(t-166)
	default:
		return time.Hour * 24 * 7 * time.Duration(t-192)
	}
}
