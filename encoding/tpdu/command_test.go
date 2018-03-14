// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"testing"

	"github.com/warthog618/sms/encoding/tpdu"
)

func TestNewCommand(t *testing.T) {
	c := tpdu.NewCommand()
	if c.MTI() != tpdu.MtCommand {
		t.Errorf("didn't set MTI - expected %v, got %v", tpdu.MtCommand, c.MTI())
	}
	if c.UDHI() {
		t.Errorf("UDHI initially set to true")
	}
	c.SetFirstOctet(0x04)
	if !c.UDHI() {
		t.Errorf("UDHI can't be set - wrong udhiMask?")
	}
}

func TestCommandSetMR(t *testing.T) {
	// also tests Command.MR
	c := tpdu.Command{}
	m := c.MR()
	if m != 0 {
		t.Errorf("initial mr should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		c.SetMR(p)
		m = c.MR()
		if m != p {
			t.Errorf("expected mr %d, got %d", p, m)
		}
	}
}

func TestCommandSetCT(t *testing.T) {
	// also tests Command.CT
	c := tpdu.Command{}
	m := c.CT()
	if m != 0 {
		t.Errorf("initial ct should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		c.SetCT(p)
		m = c.CT()
		if m != p {
			t.Errorf("expected ct %d, got %d", p, m)
		}
	}
}

func TestCommandSetMN(t *testing.T) {
	// also tests Command.MN
	c := tpdu.Command{}
	m := c.MN()
	if m != 0 {
		t.Errorf("initial mn should be 0")
	}
	for _, p := range []byte{0x00, 0xab, 0x00, 0xff} {
		c.SetMN(p)
		m = c.MN()
		if m != p {
			t.Errorf("expected mn %d, got %d", p, m)
		}
	}
}

func TestCommandSetDA(t *testing.T) {
	// also tests Command.DA
	c := tpdu.Command{}
	da := c.DA()
	nda := tpdu.Address{}
	if da != nda {
		t.Errorf("initial da should be 0")
	}
	p := tpdu.Address{Addr: "61409865629", TOA: 0x91}
	c.SetDA(p)
	da = c.DA()
	if da != p {
		t.Errorf("expected da %v, got %v", p, da)
	}
}
