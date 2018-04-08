// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestNewStatusReport(t *testing.T) {
	d := tpdu.NewStatusReport()
	if d.MTI() != tpdu.MtCommand {
		t.Errorf("didn't set MTI - expected %v, got %v", tpdu.MtCommand, d.MTI())
	}
	if d.UDHI() {
		t.Errorf("UDHI initially set to true")
	}
}

func TestStatusReportSetPI(t *testing.T) {
	// also tests StatusReport.PI
	d := tpdu.StatusReport{}
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

func TestStatusReportSetDCS(t *testing.T) {
	b := tpdu.StatusReport{}
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

func TestStatusReportSetDT(t *testing.T) {
	// also tests StatusReport.DT
	d := tpdu.StatusReport{}
	scts := d.DT()
	nscts := tpdu.Timestamp{}
	if scts != nscts {
		t.Errorf("initial scts should be 0")
	}
	p := tpdu.Timestamp{Time: time.Date(2017, time.August, 31, 11, 21, 54, 0,
		time.FixedZone("SCTS", 8*3600))}
	d.SetDT(p)
	scts = d.DT()
	if scts != p {
		t.Errorf("expected scts %v, got %v", p, scts)
	}
}

func TestStatusReportSetMR(t *testing.T) {
	// also tests StatusReport.MR
	d := tpdu.StatusReport{}
	m := d.MR()
	if m != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		d.SetMR(p)
		m = d.MR()
		if m != p {
			t.Errorf("expected pi %d, got %d", p, m)
		}
	}
}

func TestStatusReportSetRA(t *testing.T) {
	// also tests StatusReport.RA
	d := tpdu.StatusReport{}
	ra := d.RA()
	nra := tpdu.Address{}
	if ra != nra {
		t.Errorf("initial ra should be 0")
	}
	p := tpdu.Address{Addr: "61409865629", TOA: 0x91}
	d.SetRA(p)
	ra = d.RA()
	if ra != p {
		t.Errorf("expected ra %v, got %v", p, ra)
	}
}

func TestStatusReportSetSCTS(t *testing.T) {
	// also tests StatusReport.SCTS
	d := tpdu.StatusReport{}
	scts := d.SCTS()
	nscts := tpdu.Timestamp{}
	if scts != nscts {
		t.Errorf("initial scts should be 0")
	}
	p := tpdu.Timestamp{Time: time.Date(2017, time.August, 31, 11, 21, 54, 0,
		time.FixedZone("SCTS", 8*3600))}
	d.SetSCTS(p)
	scts = d.SCTS()
	if scts != p {
		t.Errorf("expected scts %v, got %v", p, scts)
	}
}

func TestStatusReportSetPID(t *testing.T) {
	b := tpdu.StatusReport{}
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

func TestStatusReportSetST(t *testing.T) {
	// also tests StatusReport.ST
	d := tpdu.StatusReport{}
	m := d.ST()
	if m != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		d.SetST(p)
		m = d.ST()
		if m != p {
			t.Errorf("expected pi %d, got %d", p, m)
		}
	}
}

func TestStatusReportSetUD(t *testing.T) {
	b := tpdu.StatusReport{}
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

func TestStatusReportSetUDH(t *testing.T) {
	b := tpdu.StatusReport{}
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
			t.Errorf("expected pi 0x04, got 0x%02x", pi)
		}
	}
}
