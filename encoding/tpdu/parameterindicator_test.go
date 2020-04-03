// SPDX-License-Identifier: MIT
//
// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.

package tpdu_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestParameterIndicatorPID(t *testing.T) {
	assert.False(t, tpdu.PI(0).PID())
	assert.True(t, tpdu.PI(tpdu.PiPID).PID())
}

func TestParameterIndicatorDCS(t *testing.T) {
	assert.False(t, tpdu.PI(0).DCS())
	assert.True(t, tpdu.PI(tpdu.PiDCS).DCS())
}

func TestParameterIndicatorUDL(t *testing.T) {
	assert.False(t, tpdu.PI(0).UDL())
	assert.True(t, tpdu.PI(tpdu.PiUDL).UDL())
}

func TestParameterIndicatorString(t *testing.T) {
	patterns := []struct {
		in  tpdu.PI
		out string
	}{
		{0, ""},
		{tpdu.PiPID, "PID"},
		{tpdu.PiDCS, "DCS"},
		{tpdu.PiUDL, "UDL"},
		{tpdu.PiPID | tpdu.PiDCS, "PID|DCS"},
		{tpdu.PiPID | tpdu.PiUDL, "PID|UDL"},
		{tpdu.PiDCS | tpdu.PiUDL, "DCS|UDL"},
		{tpdu.PiPID | tpdu.PiDCS | tpdu.PiUDL, "PID|DCS|UDL"},
	}
	for _, p := range patterns {
		assert.Equal(t, p.out, p.in.String())
	}
}
