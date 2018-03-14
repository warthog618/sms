// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"fmt"
	"testing"

	"github.com/pkg/errors"

	"github.com/warthog618/sms/encoding/tpdu"
)

func TestNewDecoder(t *testing.T) {
	errs := []error{nil, errors.New("something is broken")}
	for _, e := range errs {
		f := func(t *testing.T) {
			optCalled := false
			opt := func(d *tpdu.Decoder) error {
				optCalled = true
				return e
			}
			d, err := tpdu.NewDecoder(opt)
			if err != e {
				t.Errorf("error constructing decoder: expected %v, got %v", e, err)
			}
			if err == nil && d == nil {
				t.Errorf("failed to create decoder")
			}
			if optCalled == false {
				t.Errorf("option not called")
			}
		}
		t.Run(fmt.Sprintf("%s", e), f)
	}
}

func TestNewDecoderMO(t *testing.T) {
	d, err := tpdu.NewDecoderMO()
	if err != nil {
		t.Errorf("error crearing decoder: %v", err)
	}
	if d == nil {
		t.Error("failed to create decoder")
	}
	// !!! should check it has the correct set of decoders by performing a decode of
	// an example of each type...
}

func TestNewDecoderMT(t *testing.T) {
	d, err := tpdu.NewDecoderMT()
	if err != nil {
		t.Errorf("error crearing decoder: %v", err)
	}
	if d == nil {
		t.Error("failed to create decoder")
	}
	// !!! should check it has the correct set of decoders by performing a decode of
	// an example of each type...
}

func TestDecoderRegisterDecoder(t *testing.T) {
	d, err := tpdu.NewDecoder()
	if err != nil {
		t.Fatalf("error creating decoder: %v", err)
	}
	_, err = d.Decode([]byte{byte(tpdu.MtDeliver)}, tpdu.MT)
	if err != tpdu.DecodeError("firstOctet", 0, tpdu.ErrUnsupportedMTI(0)) {
		t.Fatal(err)
	}
	e := errors.New("called registered decoder")
	f := func(src []byte) (tpdu.TPDU, error) {
		return nil, e
	}
	err = d.RegisterDecoder(tpdu.MtDeliver, tpdu.MT, f)
	if err != nil {
		t.Fatalf("error registering decoder: %v", err)
	}
	err = d.RegisterDecoder(tpdu.MtDeliver, tpdu.MT, f)
	if err == nil {
		t.Errorf("error registering decoder: %v", err)
	}
	_, err = d.Decode([]byte{byte(tpdu.MtDeliver)}, tpdu.MT)
	if err != e {
		t.Errorf("failed to call registered decoder")
	}
}

func TestDecoderDecode(t *testing.T) {
	d, err := tpdu.NewDecoder()
	if err != nil {
		t.Fatalf("error creating decoder: %v", err)
	}
	_, err = d.Decode([]byte{byte(tpdu.MtDeliver)}, tpdu.MT)
	if err != tpdu.DecodeError("firstOctet", 0, tpdu.ErrUnsupportedMTI(0)) {
		t.Fatal(err)
	}
	e := errors.New("called registered decoder")
	f := func(src []byte) (tpdu.TPDU, error) {
		return nil, e
	}
	err = d.RegisterDecoder(tpdu.MtDeliver, tpdu.MT, f)
	if err != nil {
		t.Fatalf("error registering decoder: %v", err)
	}
	p, err := d.Decode(nil, tpdu.MT)
	if err == nil || errors.Cause(err) != tpdu.DecodeError("firstOctet", 0, tpdu.ErrUnderflow) {
		t.Errorf("unexpected error on decode: got %v", err)
	}
	if p != nil {
		t.Errorf("unexpected TPDU returned on decode: got %v", p)
	}
	s := tpdu.Submit{}
	f = func(src []byte) (tpdu.TPDU, error) {
		return &s, nil
	}
	err = d.RegisterDecoder(tpdu.MtSubmit, tpdu.MO, f)
	if err != nil {
		t.Fatalf("error registering decoder: %v", err)
	}
	p, err = d.Decode([]byte{byte(tpdu.MtSubmit)}, tpdu.MO)
	if err != nil {
		t.Errorf("unexpected error on decode: got %v", err)
	}
	if p != &s {
		t.Errorf("unexpected TPDU returned on decode: got %v", p)
	}
}
