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
	"github.com/warthog618/sms/encoding/gsm7/charset"
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/encoding/ucs2"
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
			assert.Equal(t, a, p.out)
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

type concatTestPattern struct {
	udh      tpdu.UserDataHeader
	mref     int
	segments int
	seqno    int
	ok       bool
}

func TestConcatInfo(t *testing.T) {
	patterns := []concatTestPattern{
		{tpdu.UserDataHeader{}, 0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{}}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: nil}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}}},
			3, 2, 1, true},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 1, Data: []byte{3, 2, 1}}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 8, Data: []byte{4, 3, 2, 1}}},
			1027, 2, 1, true},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{2, 1}}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 8, Data: []byte{3, 2, 1}}},
			0, 0, 0, false},
	}
	for _, p := range patterns {
		segments, seqno, mref, ok := p.udh.ConcatInfo()
		if ok != p.ok {
			t.Errorf("%v expected ok %t, got %t", p.udh, p.ok, ok)
		}
		if segments != p.segments {
			t.Errorf("%v expected segments %d, got %d", p.udh, p.segments, segments)
		}
		if seqno != p.seqno {
			t.Errorf("%v expected seqno %d, got %d", p.udh, p.seqno, seqno)
		}
		if mref != p.mref {
			t.Errorf("%v expected mref %d, got %d", p.udh, p.mref, mref)
		}
	}
}

func TestConcatInfo8(t *testing.T) {
	patterns := []concatTestPattern{
		{tpdu.UserDataHeader{}, 0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{}}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: nil}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}}},
			3, 2, 1, true},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 1, Data: []byte{3, 2, 1}}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 8, Data: []byte{4, 3, 2, 1}}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{2, 1}}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 8, Data: []byte{3, 2, 1}}},
			0, 0, 0, false},
	}
	for _, p := range patterns {
		segments, seqno, mref, ok := p.udh.ConcatInfo8()
		if ok != p.ok {
			t.Errorf("%v expected ok %t, got %t", p.udh, p.ok, ok)
		}
		if segments != p.segments {
			t.Errorf("%v expected segments %d, got %d", p.udh, p.segments, segments)
		}
		if seqno != p.seqno {
			t.Errorf("%v expected seqno %d, got %d", p.udh, p.seqno, seqno)
		}
		if mref != p.mref {
			t.Errorf("%v expected mref %d, got %d", p.udh, p.mref, mref)
		}
	}
}

func TestConcatInfo16(t *testing.T) {
	patterns := []concatTestPattern{
		{tpdu.UserDataHeader{}, 0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{}}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: nil}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 1, Data: []byte{3, 2, 1}}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 8, Data: []byte{4, 3, 2, 1}}},
			1027, 2, 1, true},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{2, 1}}},
			0, 0, 0, false},
		{tpdu.UserDataHeader{tpdu.InformationElement{ID: 8, Data: []byte{3, 2, 1}}},
			0, 0, 0, false},
	}
	for _, p := range patterns {
		segments, seqno, mref, ok := p.udh.ConcatInfo16()
		if ok != p.ok {
			t.Errorf("%v expected ok %t, got %t", p.udh, p.ok, ok)
		}
		if segments != p.segments {
			t.Errorf("%v expected segments %d, got %d", p.udh, p.segments, segments)
		}
		if seqno != p.seqno {
			t.Errorf("%v expected seqno %d, got %d", p.udh, p.seqno, seqno)
		}
		if mref != p.mref {
			t.Errorf("%v expected mref %d, got %d", p.udh, p.mref, mref)
		}
	}
}

type udDecodeTestPattern struct {
	name    string
	ud      tpdu.UserData
	udh     tpdu.UserDataHeader
	alpha   tpdu.Alphabet
	locking charset.NationalLanguageIdentifier
	shift   charset.NationalLanguageIdentifier
	msg     []byte
	err     error
}

func TestUDDDecode(t *testing.T) {
	// Also tests NewUDDecoder, AddLockingCharset and AddShiftCharset
	patterns := []udDecodeTestPattern{
		{"empty", nil, nil, 0, 0, 0, nil, nil},
		{"message 7bit", []byte("message\x10"), nil, tpdu.Alpha7Bit, 0, 0, []byte("messageΔ"), nil},
		{"message reserved", []byte("message\x10"), nil, tpdu.AlphaReserved, 0, 0, []byte("messageΔ"), nil},
		{"message 7bit esc", []byte("message\x1b"), nil, tpdu.Alpha7Bit, 0, 0, []byte("message "), nil},
		{"message 7bit locking", []byte("\x01\x02\x03"),
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 25, Data: []byte{byte(charset.Kannada)}}},
			tpdu.Alpha7Bit, charset.Kannada, 0, []byte("\u0c82\u0c83\u0c85"), nil},
		{"message 7bit shift", []byte("\x1b\x1e\x1b\x1f\x1b\x20"),
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 24, Data: []byte{byte(charset.Kannada)}}},
			tpdu.Alpha7Bit, 0, charset.Kannada, []byte("\u0ce8\u0ce9\u0cea"), nil},
		{"message 8bit", []byte("message\x1b"), nil, tpdu.Alpha8Bit, 0, 0, []byte("message\x1b"), nil},
		{"euro", []byte("\x1be"), nil, tpdu.Alpha7Bit, 0, 0, []byte("€"), nil},
		{"grin", []byte{0xd8, 0x3d, 0xde, 0x01}, nil, tpdu.AlphaUCS2, 0, 0, []byte("😁"), nil},
		{"dangling surrogate", []byte{0xd8, 0x3d},
			nil, tpdu.AlphaUCS2, 0, 0, []byte{}, ucs2.ErrDanglingSurrogate(0xD83d)},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			d, err := tpdu.NewUDDecoder()
			if d == nil || err != nil {
				t.Fatal("failed to create decoder")
			}
			if p.locking != charset.Default {
				d.AddLockingCharset(p.locking)
			}
			if p.shift != charset.Default {
				d.AddShiftCharset(p.shift)
			}
			msg, err := d.Decode(p.ud, p.udh, p.alpha)
			if err != p.err {
				t.Fatalf("error decoding %v: %v", p.ud, err)
			}
			assert.Equal(t, p.msg, msg)
		}
		t.Run(p.name, f)
	}
}

func TestUDEEncode(t *testing.T) {
	// Also tests NewUDEncoder, AddLockingCharset and AddShiftCharset
	patterns := []udDecodeTestPattern{
		{"empty", nil, nil, 0, 0, 0, nil, nil},
		{"message 7bit", []byte("message\x10"),
			nil, tpdu.Alpha7Bit, 0, 0, []byte("messageΔ"), nil},
		{"message 7bit locking", []byte("\x01\x02\x03"),
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 25, Data: []byte{byte(charset.Kannada)}}},
			tpdu.Alpha7Bit, charset.Kannada, 0, []byte("\u0c82\u0c83\u0c85"), nil},
		{"message 7bit shift", []byte("\x1b\x1e\x1b\x1f\x1b\x20"),
			tpdu.UserDataHeader{tpdu.InformationElement{ID: 24, Data: []byte{byte(charset.Kannada)}}},
			tpdu.Alpha7Bit, 0, charset.Kannada, []byte("\u0ce8\u0ce9\u0cea"), nil},
		{"euro", []byte("\x1be"), nil, tpdu.Alpha7Bit, 0, 0, []byte("€"), nil},
		{"grin", []byte{0xd8, 0x3d, 0xde, 0x01}, nil, tpdu.AlphaUCS2, 0, 0, []byte("😁"), nil},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			e, err := tpdu.NewUDEncoder()
			if e == nil || err != nil {
				t.Fatal("failed to create encoder")
			}
			if p.locking != charset.Default {
				e.AddLockingCharset(p.locking)
			}
			if p.shift != charset.Default {
				e.AddShiftCharset(p.shift)
			}
			ud, udh, alpha, err := e.Encode(string(p.msg))
			if err != p.err {
				t.Fatalf("error encoding %v: %v", p.ud, err)
			}
			assert.Equal(t, p.ud, ud)
			assert.Equal(t, p.udh, udh)
			if p.alpha != alpha {
				t.Errorf("failed to encode %s: expected alphabet %v, got %v", p.msg, p.alpha, alpha)
			}
		}
		t.Run(p.name, f)
	}
}
