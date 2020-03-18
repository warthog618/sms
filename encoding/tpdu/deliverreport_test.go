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
}

func TestDeliverReportSetDCS(t *testing.T) {
	b := tpdu.DeliverReport{}
	pi := b.PI
	if pi != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.SetDCS(p)
		d := byte(b.DCS)
		if d != p {
			t.Errorf("expected dcs %d, got %d", p, d)
		}
		pi = b.PI
		if pi&0x02 == 0x00 {
			t.Errorf("expected pi 0x02, got 0x%02x", pi)
		}
	}
}

func TestDeliverReportSetPID(t *testing.T) {
	b := tpdu.DeliverReport{}
	pi := b.PI
	if pi != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.SetPID(p)
		d := b.PID
		if d != p {
			t.Errorf("expected pid %d, got %d", p, d)
		}
		pi = b.PI
		if pi&0x01 == 0x00 {
			t.Errorf("expected pi 0x01, got 0x%02x", pi)
		}
	}
}

func TestDeliverReportSetUD(t *testing.T) {
	b := tpdu.DeliverReport{}
	pi := b.PI
	if pi != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		b.SetUD([]byte{p})
		d := b.UD
		if !bytes.Equal(d, []byte{p}) {
			t.Errorf("expected ud %d, got %d", p, d)
		}
		pi = b.PI
		if pi&0x04 == 0x00 {
			t.Errorf("expected pi 0x01, got 0x%02x", pi)
		}
	}
}

func TestDeliverReportSetUDH(t *testing.T) {
	b := tpdu.DeliverReport{}
	pi := b.PI
	if pi != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []tpdu.UserDataHeader{{
		tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}}}} {
		b.SetUDH(p)
		d := b.UDH
		if !assert.Equal(t, d, p) {
			t.Errorf("expected udh %d, got %d", p, d)
		}
		pi = b.PI
		if pi&0x04 == 0x00 {
			t.Errorf("expected pi 0x01, got 0x%02x", pi)
		}
	}
}
