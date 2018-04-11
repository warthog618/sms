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
}
