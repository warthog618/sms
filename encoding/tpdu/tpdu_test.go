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

func TestBaseTPDUAlphabet(t *testing.T) {
	patterns := []dcsAlphabetPattern{
		{0x00, tpdu.Alpha7Bit, nil},
		{0x04, tpdu.Alpha8Bit, nil},
		{0x08, tpdu.AlphaUCS2, nil},
		{0x0c, tpdu.Alpha7Bit, nil},
		{0x80, tpdu.Alpha7Bit, tpdu.ErrInvalid},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			d := tpdu.BaseTPDU{}
			d.SetDCS(tpdu.DCS(p.in))
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

func TestBaseTPDUMTI(t *testing.T) {
	b := tpdu.BaseTPDU{}
	m := b.MTI()
	if m != tpdu.MtDeliver {
		t.Errorf("initial MTI should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.SetFirstOctet(p)
		m = b.MTI()
		if m != tpdu.MessageType(p&0x3) {
			t.Errorf("expected MTI %v, got %v", p, m)
		}
	}
}

func TestBaseTPDUSetDCS(t *testing.T) {
	// also tests BaseTPDU.DCS
	b := tpdu.BaseTPDU{}
	d := b.DCS()
	if d != 0 {
		t.Errorf("initial dcs should be 0")
	}
	for _, p := range []tpdu.DCS{0x00, 0xab, 0x00, 0xff} {
		b.SetDCS(p)
		d = b.DCS()
		if d != p {
			t.Errorf("expected dcs %d, got %d", p, d)
		}
	}
}

func TestBaseTPDUSetFirstOctet(t *testing.T) {
	// also tests BaseTPDU.FirstOctet
	b := tpdu.BaseTPDU{}
	f := b.FirstOctet()
	if f != 0 {
		t.Errorf("initial firstOctet should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.SetFirstOctet(p)
		f = b.FirstOctet()
		if f != p {
			t.Errorf("expected firstOctet %d, got %d", p, f)
		}
	}
}

func TestBaseTPDUSetPID(t *testing.T) {
	// also tests BaseTPDU.PID
	b := tpdu.BaseTPDU{}
	i := b.PID()
	if i != 0 {
		t.Errorf("initial pid should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.SetPID(p)
		i = b.PID()
		if i != p {
			t.Errorf("expected pid %d, got %d", p, i)
		}
	}
}

func TestBaseTPDUSetUD(t *testing.T) {
	// also tests BaseTPDU.UD
	b := tpdu.BaseTPDU{}
	ud := b.UD()
	if len(ud) != 0 {
		t.Errorf("initial ud should be empty")
	}
	for _, p := range []tpdu.UserData{
		nil,
		{5, 6, 7},
		nil,
	} {
		b.SetUD(p)
		ud = b.UD()
		assert.Equal(t, ud, p)
	}
}

func TestBaseTPDUSetUDH(t *testing.T) {
	// also tests BaseTPDU.UDH
	b := tpdu.BaseTPDU{}
	udh := b.UDH()
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
		udh = b.UDH()
		assert.Equal(t, udh, p)
	}
}

func TestBaseTPDUUDHI(t *testing.T) {
	// also tests BaseTPDU.SetUDH
	b := tpdu.BaseTPDU{}
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
