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

func TestNewSubmitReport(t *testing.T) {
	s := tpdu.NewSubmitReport()
	if s.MTI() != tpdu.MtSubmit {
		t.Errorf("didn't set MTI - expected %v, got %v", tpdu.MtSubmit, s.MTI())
	}
	if s.UDHI() {
		t.Errorf("UDHI initially set to true")
	}
}

func TestSubmitReportSetFCS(t *testing.T) {
	// also tests SubmitReport.FCS
	s := tpdu.SubmitReport{}
	m := s.FCS()
	if m != 0 {
		t.Errorf("initial fcs should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		s.SetFCS(p)
		m = s.FCS()
		if m != p {
			t.Errorf("expected fcs %d, got %d", p, m)
		}
	}
}

func TestSubmitReportSetPI(t *testing.T) {
	// also tests SubmitReport.PI
	s := tpdu.SubmitReport{}
	m := s.PI()
	if m != 0 {
		t.Errorf("initial pi should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		s.SetPI(p)
		m = s.PI()
		if m != p {
			t.Errorf("expected pi %d, got %d", p, m)
		}
	}
}

func TestSubmitReportSetDCS(t *testing.T) {
	b := tpdu.SubmitReport{}
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

func TestSubmitReportSetPID(t *testing.T) {
	b := tpdu.SubmitReport{}
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

func TestSubmitReportSetSCTS(t *testing.T) {
	// also tests SubmitReport.SCTS
	s := tpdu.SubmitReport{}
	scts := s.SCTS()
	nscts := tpdu.Timestamp{}
	if scts != nscts {
		t.Errorf("initial scts should be 0")
	}
	p := tpdu.Timestamp{Time: time.Date(2017, time.August, 31, 11, 21, 54, 0,
		time.FixedZone("SCTS", 8*3600))}
	s.SetSCTS(p)
	scts = s.SCTS()
	if scts != p {
		t.Errorf("expected scts %v, got %v", p, scts)
	}
}

func TestSubmitReportSetUD(t *testing.T) {
	b := tpdu.SubmitReport{}
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

func TestSubmitReportSetUDH(t *testing.T) {
	b := tpdu.SubmitReport{}
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
