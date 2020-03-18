// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sar_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/ms/sar"
)

func TestNewSegmenter(t *testing.T) {
	s := sar.NewSegmenter()
	if s == nil {
		t.Fatalf("failed to create Segmenter")
	}
}

type segmentInPattern struct {
	msg []byte
	dcs byte
	udh tpdu.UserDataHeader
}

type segmentOutPattern struct {
	dcs byte
	udh tpdu.UserDataHeader
	ud  tpdu.UserData
}

func TestSegment(t *testing.T) {
	patterns := []struct {
		name string
		in   segmentInPattern
		out  []segmentOutPattern
	}{
		{"empty",
			segmentInPattern{nil, 0, nil},
			nil},
		{"single segment",
			segmentInPattern{[]byte("hello"), 0, nil},
			[]segmentOutPattern{{0, nil, []byte("hello")}}},
		{"two segment 7bit",
			segmentInPattern{[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think"), 0, nil},
			[]segmentOutPattern{
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you mi")},
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}}},
					[]byte("ght think")}},
		},
		{"three segment 7bit",
			segmentInPattern{[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think, but wait, then we also need a really really long message to trigger a three segment concatenation which requires even more characters than I care to count"), 0, nil},
			[]segmentOutPattern{
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{2, 3, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you mi")},
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{2, 3, 2}}},
					[]byte("ght think, but wait, then we also need a really really long message to trigger a three segment concatenation which requires even more characters than I c")},
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{2, 3, 3}}},
					[]byte("are to count")},
			},
		},
		{"8bit",
			segmentInPattern{[]byte("hello"), byte(tpdu.Alpha8Bit << 2), nil},
			[]segmentOutPattern{{4, nil, []byte("hello")}}},
		{"ucs2",
			segmentInPattern{[]byte("hello!"), byte(tpdu.AlphaUCS2 << 2), nil},
			[]segmentOutPattern{{8, nil, []byte("hello!")}}},
		{"reserved",
			segmentInPattern{[]byte("hello"), byte(tpdu.AlphaReserved << 2), nil},
			[]segmentOutPattern{{12, nil, []byte("hello")}}},
		{"7bit udh",
			segmentInPattern{[]byte("hello"), 0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}}},
			[]segmentOutPattern{{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}}, []byte("hello")}}},
		{"8bit udh",
			segmentInPattern{[]byte("hello"), byte(tpdu.Alpha8Bit << 2), tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}}},
			[]segmentOutPattern{{4, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}}, []byte("hello")}}},
		{"ucs udh",
			segmentInPattern{[]byte("hello"), byte(tpdu.AlphaUCS2 << 2), tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}}},
			[]segmentOutPattern{{8, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}}, []byte("hello")}}},
		{"two segment 7bit udh",
			segmentInPattern{[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think"),
				0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}}},
			[]segmentOutPattern{
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}, tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than ")},
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}, tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 2}}},
					[]byte("you might think")}},
		},
		{"two segment 8bit udh",
			segmentInPattern{[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think"),
				byte(tpdu.Alpha8Bit << 2), tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}}},
			[]segmentOutPattern{
				{4, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}, tpdu.InformationElement{ID: 0, Data: []byte{4, 2, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 charac")},
				{4, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}, tpdu.InformationElement{ID: 0, Data: []byte{4, 2, 2}}},
					[]byte("ters is more than you might think")}},
		},
		{"two segment ucs2 udh",
			segmentInPattern{[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think"),
				byte(tpdu.AlphaUCS2 << 2), tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}}},
			[]segmentOutPattern{
				{8, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}, tpdu.InformationElement{ID: 0, Data: []byte{5, 2, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 chara")},
				{8, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}, tpdu.InformationElement{ID: 0, Data: []byte{5, 2, 2}}},
					[]byte("cters is more than you might think")}},
		},
	}
	s := sar.NewSegmenter()
	for _, p := range patterns {
		f := func(t *testing.T) {
			tmpl := tpdu.Submit{}
			tmpl.DCS = tpdu.DCS(p.in.dcs)
			tmpl.SetUDH(p.in.udh)
			out := s.Segment(p.in.msg, &tmpl)
			expected := make([]tpdu.Submit, len(p.out))
			if len(p.out) == 0 {
				expected = nil
			}
			for i, o := range p.out {
				expected[i].DCS = tpdu.DCS(o.dcs)
				expected[i].SetUDH(o.udh)
				expected[i].UD = o.ud
			}
			assert.Equal(t, expected, out)
		}
		t.Run(p.name, f)
	}
}

func TestWith16BitMR(t *testing.T) {
	patterns := []struct {
		name string
		in   segmentInPattern
		out  []segmentOutPattern
	}{
		{"empty",
			segmentInPattern{nil, 0, nil},
			nil},
		{"single segment",
			segmentInPattern{[]byte("hello"), 0, nil},
			[]segmentOutPattern{{0, nil, []byte("hello")}}},
		{"two segment 7bit",
			segmentInPattern{[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think"), 0, nil},
			[]segmentOutPattern{
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 8, Data: []byte{0, 1, 2, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you mi")},
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 8, Data: []byte{0, 1, 2, 2}}},
					[]byte("ght think")}},
		},
		{"three segment 7bit",
			segmentInPattern{[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think, but wait, then we also need a really really long message to trigger a three segment concatenation which requires even more characters than I care to count"), 0, nil},
			[]segmentOutPattern{
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 8, Data: []byte{0, 2, 3, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you mi")},
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 8, Data: []byte{0, 2, 3, 2}}},
					[]byte("ght think, but wait, then we also need a really really long message to trigger a three segment concatenation which requires even more characters than I c")},
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 8, Data: []byte{0, 2, 3, 3}}},
					[]byte("are to count")},
			},
		},
		{"two segment 7bit udh",
			segmentInPattern{[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think"),
				0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}}},
			[]segmentOutPattern{
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}, tpdu.InformationElement{ID: 8, Data: []byte{0, 3, 2, 1}}},
					[]byte("this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than ")},
				{0, tpdu.UserDataHeader{tpdu.InformationElement{ID: 3, Data: []byte{1, 2, 3}}, tpdu.InformationElement{ID: 8, Data: []byte{0, 3, 2, 2}}},
					[]byte("you might think")}},
		},
	}
	s := sar.NewSegmenter(sar.With16BitMR)
	for _, p := range patterns {
		f := func(t *testing.T) {
			tmpl := tpdu.Submit{}
			tmpl.DCS = tpdu.DCS(p.in.dcs)
			tmpl.SetUDH(p.in.udh)
			out := s.Segment(p.in.msg, &tmpl)
			expected := make([]tpdu.Submit, len(p.out))
			if len(p.out) == 0 {
				expected = nil
			}
			for i, o := range p.out {
				expected[i].DCS = tpdu.DCS(o.dcs)
				expected[i].SetUDH(o.udh)
				expected[i].UD = o.ud
			}
			assert.Equal(t, expected, out)
		}
		t.Run(p.name, f)
	}
}
