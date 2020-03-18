// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestTPDUAlphabet(t *testing.T) {
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
			if err != p.err {
				t.Fatalf("error converting 0x%02x: %v", p.in, err)
			}
			if c != p.out {
				t.Errorf("expected result %v, got %v", p.out, c)
			}
		}
		t.Run(fmt.Sprintf("%02x", p.in), f)
	}

}

func TestTPDUMTI(t *testing.T) {
	b := tpdu.TPDU{}
	m := b.MTI()
	if m != tpdu.MtDeliver {
		t.Errorf("initial MTI should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.FirstOctet = p
		m = b.MTI()
		if m != tpdu.MessageType(p&0x3) {
			t.Errorf("expected MTI %v, got %v", p, m)
		}
	}
}

func TestTPDUSetUDH(t *testing.T) {
	// also tests tpdu.TPDU.UDH
	b := tpdu.TPDU{}
	udh := b.UDH
	if len(udh) != 0 {
		t.Errorf("initial udh should be empty")
	}
	for _, p := range []tpdu.UserDataHeader{
		nil,
		{tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}}},
		{tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
			tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
		},
		nil,
	} {
		b.SetUDH(p)
		udh = b.UDH
		assert.Equal(t, udh, p)
	}
}

func TestTPDUUDHI(t *testing.T) {
	// also tests tpdu.TPDU.SetUDH
	b := tpdu.TPDU{}
	udhi := b.UDHI()
	if udhi {
		t.Errorf("initial udhi should be false")
	}
	for _, p := range []tpdu.UserDataHeader{
		nil,
		{tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}}},
		{tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
			tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
		},
		nil,
	} {
		b.SetUDH(p)
		udhi = b.UDHI()
		if udhi != (len(p) != 0) {
			t.Errorf("for udh %v expected udhi %v, got %v", p, (len(p) != 0), udhi)
		}
	}
}
