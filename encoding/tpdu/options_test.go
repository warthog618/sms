// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestWithDA(t *testing.T) {
	addr := tpdu.NewAddress(tpdu.FromNumber("12345"))
	s, err := tpdu.New(tpdu.WithDA(addr))
	require.Nil(t, err)
	assert.Equal(t, addr, s.DA)
}

func TestWithOA(t *testing.T) {
	addr := tpdu.NewAddress(tpdu.FromNumber("12345"))
	s, err := tpdu.New(tpdu.WithOA(addr))
	require.Nil(t, err)
	assert.Equal(t, addr, s.OA)
}

func TestWithUDH(t *testing.T) {
	udh := tpdu.UserDataHeader{
		tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}},
	}
	s, err := tpdu.New(tpdu.WithUDH(udh))
	require.Nil(t, err)
	assert.Equal(t, udh, s.UDH)
}

func TestWithMTI(t *testing.T) {
	s, err := tpdu.New(tpdu.MtSubmit)
	require.Nil(t, err)
	assert.Equal(t, tpdu.MtSubmit, s.MTI())
}

func TestWithDirection(t *testing.T) {
	s, err := tpdu.New(tpdu.MO)
	require.Nil(t, err)
	assert.Equal(t, tpdu.MO, s.Direction)
}
