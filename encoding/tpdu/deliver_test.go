// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"testing"
	"time"

	"github.com/warthog618/sms/encoding/tpdu"
)

func TestNewDeliver(t *testing.T) {
	d := tpdu.NewDeliver()
	if d.MTI() != tpdu.MtDeliver {
		t.Errorf("didn't set MTI - expected %v, got %v", tpdu.MtDeliver, d.MTI())
	}
	if d.UDHI() {
		t.Errorf("UDHI initially set to true")
	}
	d.SetFirstOctet(0x20)
	if !d.UDHI() {
		t.Errorf("UDHI can't be set - wrong udhiMask?")
	}
}

func TestDeliverMaxUDL(t *testing.T) {
	d := tpdu.NewDeliver()
	if d.MaxUDL() != 140 {
		t.Errorf("bad maxUDL expected %d, got %d", 140, d.MaxUDL())
	}
}

func TestDeliverSetOA(t *testing.T) {
	// also tests Deliver.OA
	d := tpdu.Deliver{}
	oa := d.OA()
	noa := tpdu.Address{}
	if oa != noa {
		t.Errorf("initial oa should be 0")
	}
	p := tpdu.Address{Addr: "61409865629", TOA: 0x91}
	d.SetOA(p)
	oa = d.OA()
	if oa != p {
		t.Errorf("expected oa %v, got %v", p, oa)
	}
}

func TestDeliverSetSCTS(t *testing.T) {
	// also tests Deliver.SCTS
	d := tpdu.Deliver{}
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
