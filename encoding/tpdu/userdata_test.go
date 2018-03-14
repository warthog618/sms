// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"bytes"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
)

type udhMarshalPattern struct {
	name string
	in   tpdu.UserDataHeader
	out  []byte
	err  error
}

func TestUserDataHeaderMarshalBinary(t *testing.T) {
	patterns := []udhMarshalPattern{
		{"empty", tpdu.UserDataHeader{}, nil, nil},
		{"one", tpdu.UserDataHeader{
			tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}}},
			[]byte{5, 1, 3, 1, 2, 3}, nil},
		{"three", tpdu.UserDataHeader{
			tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
			tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
			tpdu.InformationElement{ID: 2, Data: []byte{1, 2, 3}}},
			[]byte{15, 1, 3, 1, 2, 3, 1, 3, 5, 6, 7, 2, 3, 1, 2, 3}, nil},
		// error is always nil so no error cases to test.
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			b, err := p.in.MarshalBinary()
			if err != p.err {
				t.Fatalf("error marshalling %v: %v", p.in, err)
			}
			if !bytes.Equal(b, p.out) {
				t.Errorf("failed to marshal %v: expected %v, got %v", p.in, p.out, b)
			}
		}
		t.Run(p.name, f)
	}
}

type udhUnmarshalPattern struct {
	name string
	in   []byte
	out  tpdu.UserDataHeader
	n    int
	err  error
}

func TestUserDataHeaderUnmarshalBinary(t *testing.T) {
	patterns := []udhUnmarshalPattern{
		{"one", []byte{5, 1, 3, 1, 2, 3},
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}}},
			6, nil},
		{"three", []byte{15, 1, 3, 1, 2, 3, 1, 3, 5, 6, 7, 2, 3, 1, 2, 3},
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
				tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
				tpdu.InformationElement{ID: 2, Data: []byte{1, 2, 3}}},
			16, nil},
		{"short udhl", nil, tpdu.UserDataHeader{}, 0, tpdu.DecodeError("udhl", 0, tpdu.ErrUnderflow)},
		{"short udh", []byte{5, 1, 3, 1, 2}, tpdu.UserDataHeader{}, 1, tpdu.DecodeError("ie", 1, tpdu.ErrUnderflow)},
		{"short ie", []byte{1, 1}, tpdu.UserDataHeader{}, 1, tpdu.DecodeError("ie", 1, tpdu.ErrUnderflow)},
		{"short ied", []byte{3, 1, 3, 1, 2}, tpdu.UserDataHeader{}, 3, tpdu.DecodeError("ied", 3, tpdu.ErrUnderflow)},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			a := tpdu.UserDataHeader{}
			n, err := a.UnmarshalBinary(p.in)
			if errors.Cause(err) != p.err {
				t.Fatalf("error unmarshalling %v: %v", p.in, err)
			}
			if n != p.n {
				t.Errorf("unmarshal %v read incorrect number of characters, expected %d, read %d", p.in, p.n, n)
			}
			if !assert.Equal(t, a, p.out) {
				t.Errorf("failed to unmarshal %v: expected %v, got %v", p.in, p.out, a)
			}
		}
		t.Run(p.name, f)
	}
}

func TestUserDataHeaderIE(t *testing.T) {
	u := tpdu.UserDataHeader{
		tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
		tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
		tpdu.InformationElement{ID: 2, Data: []byte{1, 2, 3}},
	}
	_, ok := u.IE(0)
	if ok {
		t.Errorf("found non existent IE")
	}
	i, ok := u.IE(1)
	if !ok {
		t.Fatalf("failed to find IE")
	}
	if i.ID != 1 {
		t.Errorf("returned wrong IE, expected ID %d, got %d", 1, i.ID)
	}
	if !bytes.Equal(i.Data, u[1].Data) {
		t.Errorf("returned wrong IE, expected Data %v, got %v", u[1].Data, i.Data)
	}
}

func TestUserDataHeaderIEs(t *testing.T) {
	u := tpdu.UserDataHeader{
		tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
		tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
		tpdu.InformationElement{ID: 2, Data: []byte{1, 2, 3}},
	}
	i := u.IEs(0)
	if i != nil {
		t.Errorf("found non existent IE")
	}
	i = u.IEs(2)
	if len(i) != 1 {
		t.Errorf("returned wrong numner of IEs, expected %d, got %d", 1, len(i))
	}
	if i[0].ID != 2 {
		t.Errorf("returned wrong IE, expected ID %d, got %d", 2, i[0].ID)
	}
	if !bytes.Equal(i[0].Data, u[2].Data) {
		t.Errorf("returned wrong IE, expected Data %v, got %v", u[2].Data, i[0].Data)
	}
	i = u.IEs(1)
	if len(i) != 2 {
		t.Errorf("returned wrong numner of IEs, expected %d, got %d", 2, len(i))
	}
	if i[0].ID != 1 {
		t.Errorf("returned wrong IE, expected ID %d, got %d", 1, i[0].ID)
	}
	if !bytes.Equal(i[0].Data, u[0].Data) {
		t.Errorf("returned wrong IE, expected Data %v, got %v", u[0].Data, i[0].Data)
	}
	if i[1].ID != 1 {
		t.Errorf("returned wrong IE, expected ID %d, got %d", 1, i[1].ID)
	}
	if !bytes.Equal(i[1].Data, u[1].Data) {
		t.Errorf("returned wrong IE, expected Data %v, got %v", u[1].Data, i[1].Data)
	}
}
