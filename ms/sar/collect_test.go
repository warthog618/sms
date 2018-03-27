// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sar_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/ms/sar"
)

func TestNewCollector(t *testing.T) {
	e := func(err error) {}
	c := sar.NewCollector(time.Minute, e)
	if c == nil {
		t.Fatalf("failed to create Collector")
	}
}

func TestCollectorClose(t *testing.T) {
	e := func(err error) {}
	c := sar.NewCollector(time.Minute, e)
	if c == nil {
		t.Fatalf("failed to create Collector")
	}
	c.Close() // when open
	c.Close() // when closed
	c = sar.NewCollector(time.Minute, e)
	if c == nil {
		t.Fatalf("failed to create Collector")
	}
	d := tpdu.Deliver{}
	d.SetOA(tpdu.Address{Addr: "1234", TOA: 0x91})
	d.SetUDH(tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{5, 2, 1}}})
	c.Collect(&d)
	c.Close() // with pipe active
}

type collectTestPattern struct {
	name string
	oa   tpdu.Address
	in   tpdu.UserDataHeader
	out  []*tpdu.UserDataHeader
	err  error
}

func TestCollectorCollect(t *testing.T) {
	patterns := []collectTestPattern{
		// the patterns are Collected sequentially, so the return value depends on
		// the preceding set of PDUs, not just the individual in.
		// The resulting tests must be run as a complete set, not individually, or some will fail.
		{"single segment", tpdu.Address{Addr: "1234", TOA: 0x91}, nil,
			[]*tpdu.UserDataHeader{nil}, nil},
		// 1 segment (shouldn't be seen in practice, but test in case)
		{"one segment", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 1, 1}}},
			[]*tpdu.UserDataHeader{{tpdu.InformationElement{ID: 0, Data: []byte{1, 1, 1}}}},
			nil},
		{"two a", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}},
			nil,
			nil},
		{"two b", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}}},
			[]*tpdu.UserDataHeader{
				{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}},
				{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}}},
			},
			nil},
		{"three a", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 1}}},
			nil,
			nil},
		{"three b", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 2}}},
			nil,
			nil},
		{"three c", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 3}}},
			[]*tpdu.UserDataHeader{
				{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 1}}},
				{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 2}}},
				{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 3}}},
			},
			nil},
		{"jumbled a", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 1}}},
			nil,
			nil},
		{"jumbled c", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 3}}},
			nil,
			nil},
		{"duplicate", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 3}}},
			nil,
			sar.ErrDuplicateSegment},
		{"jumbled b", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 2}}},
			[]*tpdu.UserDataHeader{
				{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 1}}},
				{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 2}}},
				{tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 3}}},
			},
			nil},
		{"concurrent one a", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}},
			nil,
			nil},
		{"concurrent two a", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 1}}},
			nil,
			nil},
		{"concurrent one b", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}}},
			[]*tpdu.UserDataHeader{
				{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}},
				{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}}},
			},
			nil},
		{"concurrent two b", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 2}}},
			[]*tpdu.UserDataHeader{
				{tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 1}}},
				{tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 2}}},
			},
			nil},
		{"zero seqno", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 0}}},
			nil,
			sar.ErrReassemblyInconsistency},
		{"large seqno", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 3}}},
			nil,
			sar.ErrReassemblyInconsistency},
	}
	closedPatterns := []collectTestPattern{
		{"closed single segment", tpdu.Address{Addr: "1234", TOA: 0x91}, nil,
			[]*tpdu.UserDataHeader{nil}, nil},
		{"closed concat", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}},
			nil,
			sar.ErrClosed},
	}
	var ae error
	asyncError := func(err error) {
		ae = err
	}
	c := sar.NewCollector(time.Minute, asyncError)
	if c == nil {
		t.Fatalf("failed to create Collector")
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			testCollect(t, c, p)
		}
		t.Run(p.name, f)
		if ae != nil {
			t.Errorf("%s expiry called unexpectedly: %v", p.name, ae)
			ae = nil
		}
	}
	c.Close()
	for _, p := range closedPatterns {
		f := func(t *testing.T) {
			testCollect(t, c, p)
		}
		t.Run(p.name, f)
		if ae != nil {
			t.Errorf("%s expiry called unexpectedly: %v", p.name, ae)
			ae = nil
		}
	}
}

func TestCollectorCollectExpiry(t *testing.T) {
	// not happy with expiry testing - would like more control to allow pipes
	// to have multiple segments before expiring....
	var aechan = make(chan error)
	asyncError := func(err error) {
		aechan <- err
	}
	c := sar.NewCollector(time.Millisecond, asyncError)
	expiredPatterns := []collectTestPattern{
		{"two a", tpdu.Address{Addr: "1234", TOA: 0x91},
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}},
			[]*tpdu.UserDataHeader{{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}}, nil},
			nil,
		},
	}
	for _, p := range expiredPatterns {
		f := func(t *testing.T) {
			in := tpdu.Deliver{}
			in.SetOA(p.oa)
			in.SetUDH(p.in)
			out, err := c.Collect(&in)
			if err != p.err {
				t.Errorf("collect returned unexpected error %v", err)
			}
			if out != nil {
				t.Errorf("collect returned unexpected error %v", err)
			}
			expected := make([]*tpdu.Deliver, len(p.out))
			if len(p.out) == 0 {
				expected = nil
			}
			for i, udh := range p.out {
				if udh != nil {
					r := tpdu.Deliver{}
					r.SetOA(p.oa)
					r.SetUDH(*udh)
					expected[i] = &r
				}
			}
			select {
			case d := <-aechan:
				if x, ok := d.(sar.ErrExpired); ok {
					assert.Equal(t, expected, x.T)
				} else {
					t.Errorf("collect returned unexpected async error %v", err)
				}
			case <-time.After(100 * time.Millisecond):
				t.Errorf("timeout waiting for expiry")
			}
		}
		t.Run(p.name, f)
	}
}

func testCollect(t *testing.T, c *sar.Collector, p collectTestPattern) {
	in := tpdu.Deliver{}
	in.SetOA(p.oa)
	in.SetUDH(p.in)
	out, err := c.Collect(&in)
	if err != p.err {
		t.Errorf("collect returned unexpected error %v", err)
	}
	expected := make([]*tpdu.Deliver, len(p.out))
	if len(p.out) == 0 {
		expected = nil
	}
	for i, udh := range p.out {
		r := tpdu.Deliver{}
		r.SetOA(p.oa)
		if udh != nil {
			r.SetUDH(*udh)
		}
		expected[i] = &r
	}
	assert.Equal(t, expected, out)
}

func TestErrExpired(t *testing.T) {
	d := tpdu.Deliver{}
	reassembly := make([]*tpdu.Deliver, 2)
	reassembly[0] = &d
	err := sar.ErrExpired{T: reassembly}
	expected := fmt.Sprintf("sar: timed out reassembling %v", reassembly)
	s := err.Error()
	if s != expected {
		t.Errorf("failed to stringify, expected '%s', got '%s'", expected, s)
	}
	assert.Equal(t, reassembly, err.T)
}
