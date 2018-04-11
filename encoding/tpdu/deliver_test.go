// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"testing"

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
}

func TestDeliverMaxUDL(t *testing.T) {
	d := tpdu.NewDeliver()
	if d.MaxUDL() != 140 {
		t.Errorf("bad maxUDL expected %d, got %d", 140, d.MaxUDL())
	}
}
