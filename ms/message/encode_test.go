// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package message_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/ms/message"
	"github.com/warthog618/sms/ms/sar"
)

func TestNewEncoder(t *testing.T) {
	ude, _ := tpdu.NewUDEncoder()
	s := sar.NewSegmenter()
	e := message.NewEncoder(ude, s)
	if e == nil {
		t.Fatalf("failed to create Encoder")
	}
}

type encodeOutPattern struct {
	da  tpdu.Address
	dcs int
	udh tpdu.UserDataHeader
	ud  tpdu.UserData
}

func TestEncode(t *testing.T) {
	patterns := []struct {
		name   string
		number string
		msg    string
		out    []encodeOutPattern
		err    error
	}{
		{"empty", "", "", nil, nil},
		{"single segment", "1234", "hello",
			[]encodeOutPattern{{tpdu.Address{Addr: "1234", TOA: 0x91}, 0, nil, []byte("hello")}}, nil},
		{"plus number", "+1234", "hello",
			[]encodeOutPattern{{tpdu.Address{Addr: "1234", TOA: 0x91}, 0, nil, []byte("hello")}}, nil},
		{"two segment 7bit", "1234", "this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think",
			[]encodeOutPattern{
				{tpdu.Address{Addr: "1234", TOA: 0x91}, 0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you mi")},
				{tpdu.Address{Addr: "1234", TOA: 0x91}, 0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}}},
					[]byte("ght think")}},
			nil,
		},
	}
	ude, _ := tpdu.NewUDEncoder()
	s := sar.NewSegmenter()
	e := message.NewEncoder(ude, s)
	for _, p := range patterns {
		f := func(t *testing.T) {
			out, err := e.Encode(p.number, p.msg)
			if err != p.err {
				t.Errorf("encode returned unexpected error %v", err)
			}
			expected := make([]tpdu.Submit, len(p.out))
			if len(p.out) == 0 {
				expected = nil
			}
			for i, o := range p.out {
				expected[i].SetFirstOctet(1)
				expected[i].SetDA(o.da)
				expected[i].SetDCS(tpdu.DCS(o.dcs))
				expected[i].SetUDH(o.udh)
				expected[i].SetUD(o.ud)
			}
			assert.Equal(t, expected, out)
		}
		t.Run(p.name, f)
	}
}

type MockUDEncoder struct{}

func (m MockUDEncoder) Encode(msg string) (tpdu.UserData, tpdu.UserDataHeader, tpdu.Alphabet, error) {
	return nil, nil, 0, fmt.Errorf("mock encode failed for '%s'", msg)
}

func TestEncodeError(t *testing.T) {
	ude := MockUDEncoder{}
	s := sar.NewSegmenter()
	e := message.NewEncoder(ude, s)
	out, err := e.Encode("1234", "hello")
	if err.Error() != "mock encode failed for 'hello'" {
		t.Errorf("encode returned unexpected error %v", err)
	}
	if out != nil {
		t.Errorf("encode returned unexpected result %v", out)
	}
}

func TestEncodeWithTemplate(t *testing.T) {
	patterns := []struct {
		name   string
		number string
		msg    string
		out    []encodeOutPattern
		err    error
	}{
		{"empty", "", "", nil, nil},
		{"single segment", "1234", "hello",
			[]encodeOutPattern{{tpdu.Address{Addr: "1234", TOA: 0x91}, 0,
				tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}},
				[]byte("hello")}},
			nil},
		{"plus number", "+1234", "hello",
			[]encodeOutPattern{{tpdu.Address{Addr: "1234", TOA: 0x91}, 0,
				tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}},
				[]byte("hello")}},
			nil},
		{"two segment 7bit", "1234", "this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think",
			[]encodeOutPattern{
				{tpdu.Address{Addr: "1234", TOA: 0x91}, 0,
					tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}},
						tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than ")},
				{tpdu.Address{Addr: "1234", TOA: 0x91}, 0,
					tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}},
						tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}}},
					[]byte("you might think")}},
			nil,
		},
	}
	ude, _ := tpdu.NewUDEncoder()
	s := sar.NewSegmenter()
	e := message.NewEncoder(ude, s)
	tmpl := tpdu.NewSubmit()
	tmpl.SetDCS(0xe3) // doesn't support alphabet
	tmpl.SetUDH(tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}})
	e.SetT(tmpl)
	for _, p := range patterns {
		f := func(t *testing.T) {
			out, err := e.Encode(p.number, p.msg)
			if err != p.err {
				t.Errorf("encode returned unexpected error %v", err)
			}
			expected := make([]tpdu.Submit, len(p.out))
			if len(p.out) == 0 {
				expected = nil
			}
			for i, o := range p.out {
				expected[i].SetFirstOctet(65)
				expected[i].SetDA(o.da)
				expected[i].SetDCS(tpdu.DCS(o.dcs))
				expected[i].SetUDH(o.udh)
				expected[i].SetUD(o.ud)
			}
			assert.Equal(t, expected, out)
		}
		t.Run(p.name, f)
	}
}

func TestEncode8Bit(t *testing.T) {
	patterns := []struct {
		name   string
		number string
		msg    string
		out    []encodeOutPattern
		err    error
	}{
		{"empty", "", "", nil, nil},
		{"single segment", "1234", "hello",
			[]encodeOutPattern{{tpdu.Address{Addr: "1234", TOA: 0x91}, 4, nil, []byte("hello")}}, nil},
		{"plus number", "+1234", "hello",
			[]encodeOutPattern{{tpdu.Address{Addr: "1234", TOA: 0x91}, 4, nil, []byte("hello")}}, nil},
		{"two segment 7bit", "1234", "this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think",
			[]encodeOutPattern{
				{tpdu.Address{Addr: "1234", TOA: 0x91}, 4, tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters ")},
				{tpdu.Address{Addr: "1234", TOA: 0x91}, 4, tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}}},
					[]byte("is more than you might think")}},
			nil,
		},
	}
	ude, _ := tpdu.NewUDEncoder()
	s := sar.NewSegmenter()
	e := message.NewEncoder(ude, s)
	for _, p := range patterns {
		f := func(t *testing.T) {
			out, err := e.Encode8Bit(p.number, []byte(p.msg))
			if err != p.err {
				t.Errorf("encode returned unexpected error %v", err)
			}
			expected := make([]tpdu.Submit, len(p.out))
			if len(p.out) == 0 {
				expected = nil
			}
			for i, o := range p.out {
				expected[i].SetFirstOctet(1)
				expected[i].SetDA(o.da)
				expected[i].SetDCS(tpdu.DCS(o.dcs))
				expected[i].SetUDH(o.udh)
				expected[i].SetUD(o.ud)
			}
			assert.Equal(t, expected, out)
		}
		t.Run(p.name, f)
	}
}

func TestEncode8BitWithTemplate(t *testing.T) {
	patterns := []struct {
		name   string
		number string
		msg    string
		out    []encodeOutPattern
		err    error
	}{
		{"empty", "", "", nil, nil},
		{"single segment", "1234", "hello",
			[]encodeOutPattern{{tpdu.Address{Addr: "1234", TOA: 0x91}, 4,
				tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}},
				[]byte("hello")}},
			nil},
		{"plus number", "+1234", "hello",
			[]encodeOutPattern{{tpdu.Address{Addr: "1234", TOA: 0x91}, 4,
				tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}},
				[]byte("hello")}},
			nil},
		{"two segment 7bit", "1234", "this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think",
			[]encodeOutPattern{
				{tpdu.Address{Addr: "1234", TOA: 0x91}, 4,
					tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}},
						tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 charac")},
				{tpdu.Address{Addr: "1234", TOA: 0x91}, 4,
					tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}},
						tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}}},
					[]byte("ters is more than you might think")}},
			nil,
		},
	}
	ude, _ := tpdu.NewUDEncoder()
	s := sar.NewSegmenter()
	e := message.NewEncoder(ude, s)
	tmpl := tpdu.NewSubmit()
	tmpl.SetDCS(0xe3) // doesn't support alphabet
	tmpl.SetUDH(tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}})
	e.SetT(tmpl)
	for _, p := range patterns {
		f := func(t *testing.T) {
			out, err := e.Encode8Bit(p.number, []byte(p.msg))
			if err != p.err {
				t.Errorf("encode returned unexpected error %v", err)
			}
			expected := make([]tpdu.Submit, len(p.out))
			if len(p.out) == 0 {
				expected = nil
			}
			for i, o := range p.out {
				expected[i].SetFirstOctet(65)
				expected[i].SetDA(o.da)
				expected[i].SetDCS(tpdu.DCS(o.dcs))
				expected[i].SetUDH(o.udh)
				expected[i].SetUD(o.ud)
			}
			assert.Equal(t, expected, out)
		}
		t.Run(p.name, f)
	}
}
