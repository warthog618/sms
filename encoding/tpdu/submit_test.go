// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestNewSubmit(t *testing.T) {
	d := tpdu.NewSubmit()
	if d.MTI() != tpdu.MtSubmit {
		t.Errorf("didn't set MTI - expected %v, got %v", tpdu.MtDeliver, d.MTI())
	}
	if d.UDHI() {
		t.Errorf("UDHI initially set to true")
	}
}

func TestSubmitMaxUDL(t *testing.T) {
	s := tpdu.NewSubmit()
	if s.MaxUDL() != 140 {
		t.Errorf("bad maxUDL expected %d, got %d", 140, s.MaxUDL())
	}
}

func TestSubmitSetValidityPeriod(t *testing.T) {
	// also tests Submit.VP
	s := tpdu.Submit{}
	vp := s.VP
	if vp.Format != tpdu.VpfNotPresent {
		t.Errorf("initial vp should be nil")
	}
	pvp := tpdu.ValidityPeriod{}
	pvp.SetRelative(time.Duration(100000000))
	for _, p := range []struct {
		vp tpdu.ValidityPeriod
		fo byte
	}{{tpdu.ValidityPeriod{}, 0x00},
		{pvp, 0x08},
		{tpdu.ValidityPeriod{}, 0x00}} {
		s.SetVP(p.vp)
		vp = s.VP
		if p.fo != s.FirstOctet {
			t.Errorf("expected firstOctet 0x%02x, got 0x%02x", p.fo, s.FirstOctet)
		}
		if !assert.Equal(t, vp, p.vp) {
			t.Errorf("expected vp %v, got %v", p.vp, vp)
		}
	}
}
