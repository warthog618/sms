// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

import (
	"time"

	"github.com/warthog618/sms/encoding/bcd"
)

// Timestamp represents a SCTS timestamp, as defined in 3GPP TS 23.040 Section 9.2.3.11.
type Timestamp struct {
	time.Time
}

// MarshalBinary encodes the SCTS timestamp into binary.
func (t *Timestamp) MarshalBinary() (dst []byte, err error) {
	dst = make([]byte, 7)
	y := t.Year() % 100
	f := []int{y, int(t.Month()), t.Day(), t.Hour(), t.Minute(), t.Second()}
	for i, v := range f {
		dst[i], err = bcd.Encode(v)
		// this should never trip, assuming the time methods return values in
		// the expected ranges...
		if err != nil {
			return nil, err
		}
	}
	_, tz := t.Zone()
	dst[6], err = bcd.EncodeSigned(tz / (15 * 60))
	if err != nil {
		return nil, err
	}
	return dst, nil
}

// UnmarshalBinary decodes the SCTS timestamp.
func (t *Timestamp) UnmarshalBinary(src []byte) error {
	if len(src) < 7 {
		return ErrUnderflow
	}
	i := make([]int, 6)
	var err error
	for idx := 0; idx < 6; idx++ {
		i[idx], err = bcd.Decode(src[idx])
		if err != nil {
			return err
		}
	}
	tz, err := bcd.DecodeSigned(src[6])
	if err != nil {
		return err
	}
	loc := time.UTC
	if tz != 0 {
		tzoffset := tz * 15 * 60 // seconds east of UTC
		loc = time.FixedZone("SCTS", tzoffset)
	}
	year := i[0]
	if year < 70 {
		year += 2000
	} else {
		year += 1900
	}
	mon := time.Month(i[1])
	t.Time = time.Date(year, mon, i[2], i[3], i[4], i[5], 0, loc)
	return nil
}
