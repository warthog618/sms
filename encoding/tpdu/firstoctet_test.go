// SPDX-License-Identifier: MIT
//
// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.

package tpdu_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestFirstOctetLP(t *testing.T) {
	assert.False(t, tpdu.FirstOctet(0).LP())
	assert.True(t, tpdu.FirstOctet(tpdu.FoLP).LP())
}

func TestFirstOctetMMS(t *testing.T) {
	assert.False(t, tpdu.FirstOctet(0).MMS())
	assert.True(t, tpdu.FirstOctet(tpdu.FoMMS).MMS())
}

func TestFirstMTI(t *testing.T) {
	patterns := []struct {
		inout tpdu.MessageType
	}{
		{tpdu.MtCommand},
		{tpdu.MtDeliver},
		{tpdu.MtSubmit},
		{tpdu.MtReserved},
	}
	for _, p := range patterns {
		assert.Equal(t, p.inout, tpdu.FirstOctet(p.inout).MTI())
	}
}
func TestFirstOctetRD(t *testing.T) {
	assert.False(t, tpdu.FirstOctet(0).RD())
	assert.True(t, tpdu.FirstOctet(tpdu.FoRD).RD())
}

func TestFirstOctetRP(t *testing.T) {
	assert.False(t, tpdu.FirstOctet(0).RP())
	assert.True(t, tpdu.FirstOctet(tpdu.FoRP).RP())
}

func TestFirstOctetSRI(t *testing.T) {
	assert.False(t, tpdu.FirstOctet(0).SRI())
	assert.True(t, tpdu.FirstOctet(tpdu.FoSRI).SRI())
}

func TestFirstOctetSRR(t *testing.T) {
	assert.False(t, tpdu.FirstOctet(0).SRR())
	assert.True(t, tpdu.FirstOctet(tpdu.FoSRR).SRR())
}

func TestFirstOctetSRQ(t *testing.T) {
	assert.False(t, tpdu.FirstOctet(0).SRQ())
	assert.True(t, tpdu.FirstOctet(tpdu.FoSRQ).SRQ())
}

func TestFirstOctetUDHI(t *testing.T) {
	assert.False(t, tpdu.FirstOctet(0).UDHI())
	assert.True(t, tpdu.FirstOctet(tpdu.FoUDHI).UDHI())
}

func TestFirstOctetVPF(t *testing.T) {
	patterns := []struct {
		in  tpdu.FirstOctet
		out tpdu.ValidityPeriodFormat
	}{
		{0, tpdu.VpfNotPresent},
		{0x10, tpdu.VpfRelative},
		{0x18, tpdu.VpfAbsolute},
		{0x08, tpdu.VpfEnhanced},
	}
	for _, p := range patterns {
		fo := tpdu.FirstOctet(p.in)
		assert.Equal(t, p.out, fo.VPF())
	}
}

func TestFirstWithMTI(t *testing.T) {
	patterns := []struct {
		inout tpdu.MessageType
	}{
		{tpdu.MtCommand},
		{tpdu.MtDeliver},
		{tpdu.MtSubmit},
		{tpdu.MtReserved},
	}
	for _, p := range patterns {
		fo := tpdu.FirstOctet(0).WithMTI(p.inout)
		assert.Equal(t, p.inout, fo.MTI())
	}
}

func TestFirstOctetWithVPF(t *testing.T) {
	patterns := []struct {
		inout tpdu.ValidityPeriodFormat
	}{
		{tpdu.VpfNotPresent},
		{tpdu.VpfRelative},
		{tpdu.VpfAbsolute},
		{tpdu.VpfEnhanced},
	}
	for _, p := range patterns {
		assert.Equal(t, p.inout, tpdu.FirstOctet(0).WithVPF(p.inout).VPF())
	}
}
