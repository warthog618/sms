// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestNewDeliverReport(t *testing.T) {
	d := tpdu.NewDeliverReport()
	if d.MTI() != tpdu.MtDeliver {
		t.Errorf("didn't set MTI - expected %v, got %v", tpdu.MtDeliver, d.MTI())
	}
	if d.UDHI() {
		t.Errorf("UDHI initially set to true")
	}
	d.SetFirstOctet(0x04)
	if !d.UDHI() {
		t.Errorf("UDHI can't be set - wrong udhiMask?")
	}
}

func TestDeliverReportSetFCS(t *testing.T) {
	// also tests DeliverReport.FCS
	d := tpdu.DeliverReport{}
	m := d.FCS()
	if m != 0 {
		t.Errorf("initial fcs should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		d.SetFCS(p)
		m = d.FCS()
		if m != p {
			t.Errorf("expected fcs %d, got %d", p, m)
		}
	}
}

func TestDeliverReportSetPI(t *testing.T) {
	// also tests DeliverReport.PI
	d := tpdu.DeliverReport{}
	m := d.PI()
	if m != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		d.SetPI(p)
		m = d.PI()
		if m != p {
			t.Errorf("expected pi %d, got %d", p, m)
		}
	}
}

func TestDeliverReportSetDCS(t *testing.T) {
	b := tpdu.DeliverReport{}
	pi := b.PI()
	if pi != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []tpdu.DCS{0x00, 0xab, 0x00, 0xff} {
		b.SetDCS(p)
		d := b.DCS()
		if d != p {
			t.Errorf("expected dcs %d, got %d", p, d)
		}
		pi = b.PI()
		if pi&0x02 == 0x00 {
			t.Errorf("expected pi 0x02, got 0x%02x", pi)
		}
	}
}

func TestDeliverReportSetPID(t *testing.T) {
	b := tpdu.DeliverReport{}
	pi := b.PI()
	if pi != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.SetPID(p)
		d := b.PID()
		if d != p {
			t.Errorf("expected pid %d, got %d", p, d)
		}
		pi = b.PI()
		if pi&0x01 == 0x00 {
			t.Errorf("expected pi 0x01, got 0x%02x", pi)
		}
	}
}

func TestDeliverReportSetUD(t *testing.T) {
	b := tpdu.DeliverReport{}
	pi := b.PI()
	if pi != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.SetUD([]byte{p})
		d := b.UD()
		if !bytes.Equal(d, []byte{p}) {
			t.Errorf("expected ud %d, got %d", p, d)
		}
		pi = b.PI()
		if pi&0x04 == 0x00 {
			t.Errorf("expected pi 0x01, got 0x%02x", pi)
		}
	}
}

func TestDeliverReportSetUDH(t *testing.T) {
	b := tpdu.DeliverReport{}
	pi := b.PI()
	if pi != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []tpdu.UserDataHeader{{
		tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}}}} {
		b.SetUDH(p)
		d := b.UDH()
		if !assert.Equal(t, d, p) {
			t.Errorf("expected udh %d, got %d", p, d)
		}
		pi = b.PI()
		if pi&0x04 == 0x00 {
			t.Errorf("expected pi 0x01, got 0x%02x", pi)
		}
	}
}
