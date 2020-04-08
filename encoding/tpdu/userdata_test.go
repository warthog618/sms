// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warthog618/sms/encoding/gsm7/charset"
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/encoding/ucs2"
)

func TestUserDataHeaderMarshalBinary(t *testing.T) {
	patterns := []struct {
		name string
		in   tpdu.UserDataHeader
		out  []byte
		err  error
	}{
		{"empty",
			tpdu.UserDataHeader{},
			nil,
			nil,
		},
		{"one",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
			},
			[]byte{5, 1, 3, 1, 2, 3},
			nil,
		},
		{"three",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
				tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
				tpdu.InformationElement{ID: 2, Data: []byte{1, 2, 3}},
			},
			[]byte{15, 1, 3, 1, 2, 3, 1, 3, 5, 6, 7, 2, 3, 1, 2, 3}, nil},
		// error is always nil so no error cases to test.
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			b, err := p.in.MarshalBinary()
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, b)
		}
		t.Run(p.name, f)
	}
}

func TestUserDataHeaderUnmarshalBinary(t *testing.T) {
	patterns := []struct {
		name string
		in   []byte
		out  tpdu.UserDataHeader
		n    int
		err  error
	}{
		{"one",
			[]byte{5, 1, 3, 1, 2, 3},
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
			},
			6,
			nil,
		},
		{"three",
			[]byte{15, 1, 3, 1, 2, 3, 1, 3, 5, 6, 7, 2, 3, 1, 2, 3},
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
				tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
				tpdu.InformationElement{ID: 2, Data: []byte{1, 2, 3}},
			},
			16,
			nil,
		},
		{"short udhl",
			nil,
			tpdu.UserDataHeader{},
			0,
			tpdu.NewDecodeError("udhl", 0, tpdu.ErrUnderflow),
		},
		{"short udh",
			[]byte{5, 1, 3, 1, 2},
			tpdu.UserDataHeader{},
			1,
			tpdu.NewDecodeError("ie", 1, tpdu.ErrUnderflow),
		},
		{"short ie",
			[]byte{1, 1},
			tpdu.UserDataHeader{},
			1,
			tpdu.NewDecodeError("ie", 1, tpdu.ErrUnderflow),
		},
		{"short ied",
			[]byte{3, 1, 3, 1, 2},
			tpdu.UserDataHeader{},
			3,
			tpdu.NewDecodeError("ied", 3, tpdu.ErrUnderflow),
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			a := tpdu.UserDataHeader{}
			n, err := a.UnmarshalBinary(p.in)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.n, n)
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
	assert.False(t, ok)
	i, ok := u.IE(1)
	assert.True(t, ok)
	assert.Equal(t, uint8(1), i.ID)
	assert.Equal(t, u[1].Data, i.Data)
}

func TestUserDataHeaderIEs(t *testing.T) {
	u := tpdu.UserDataHeader{
		tpdu.InformationElement{ID: 1, Data: []byte{1, 2, 3}},
		tpdu.InformationElement{ID: 1, Data: []byte{5, 6, 7}},
		tpdu.InformationElement{ID: 2, Data: []byte{1, 2, 3}},
	}
	i := u.IEs(0)
	assert.Nil(t, i)
	i = u.IEs(2)
	assert.Equal(t, 1, len(i))
	assert.Equal(t, uint8(2), i[0].ID)
	assert.Equal(t, u[2].Data, i[0].Data)
	i = u.IEs(1)
	assert.Equal(t, 2, len(i))
	assert.Equal(t, uint8(1), i[0].ID)
	assert.Equal(t, u[0].Data, i[0].Data)
	assert.Equal(t, uint8(1), i[1].ID)
	assert.Equal(t, u[1].Data, i[1].Data)
}

type concatTestPattern struct {
	name     string
	udh      tpdu.UserDataHeader
	mref     int
	segments int
	seqno    int
	ok       bool
}

func TestConcatInfo(t *testing.T) {
	patterns := []concatTestPattern{
		{"empty",
			tpdu.UserDataHeader{},
			0,
			0,
			0,
			false,
		},
		{"empty data",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{}},
			},
			0,
			0,
			0,
			false,
		},
		{"nil data",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: nil},
			},
			0,
			0,
			0,
			false,
		},
		{"concat8",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}},
			},
			3,
			2,
			1,
			true,
		},
		{"id 1",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 1, Data: []byte{3, 2, 1}},
			},
			0,
			0,
			0,
			false,
		},
		{"concat16",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 8, Data: []byte{4, 3, 2, 1}},
			},
			1027,
			2,
			1,
			true,
		},
		{"short concat8",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{2, 1}},
			},
			0,
			0,
			0,
			false,
		},
		{"short concat16",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 8, Data: []byte{3, 2, 1}},
			},
			0,
			0,
			0,
			false,
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			segments, seqno, mref, ok := p.udh.ConcatInfo()
			assert.Equal(t, p.ok, ok)
			assert.Equal(t, p.segments, segments)
			assert.Equal(t, p.seqno, seqno)
			assert.Equal(t, p.mref, mref)
		}
		t.Run(p.name, f)
	}
}

func TestConcatInfo8(t *testing.T) {
	patterns := []concatTestPattern{
		{"empty",
			tpdu.UserDataHeader{},
			0,
			0,
			0,
			false,
		},
		{"empty data",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{}},
			},
			0,
			0,
			0,
			false,
		},
		{"nil data",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: nil},
			},
			0,
			0,
			0,
			false,
		},
		{"concat8",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}},
			},
			3,
			2,
			1,
			true,
		},
		{"id 1",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 1, Data: []byte{3, 2, 1}},
			},
			0,
			0,
			0,
			false,
		},
		{"concat16",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 8, Data: []byte{4, 3, 2, 1}},
			},
			0,
			0,
			0,
			false,
		},
		{"short concat8",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{2, 1}},
			},
			0,
			0,
			0,
			false,
		},
		{"short concat16",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 8, Data: []byte{3, 2, 1}},
			},
			0,
			0,
			0,
			false,
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			segments, seqno, mref, ok := p.udh.ConcatInfo8()
			assert.Equal(t, p.ok, ok)
			assert.Equal(t, p.segments, segments)
			assert.Equal(t, p.seqno, seqno)
			assert.Equal(t, p.mref, mref)
		}
		t.Run(p.name, f)
	}
}

func TestConcatInfo16(t *testing.T) {
	patterns := []concatTestPattern{
		{"empty",
			tpdu.UserDataHeader{},
			0,
			0,
			0,
			false,
		},
		{"empty data",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{}},
			},
			0,
			0,
			0,
			false,
		},
		{"nil data",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: nil},
			},
			0,
			0,
			0,
			false,
		},
		{"concat8",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{3, 2, 1}},
			},
			0,
			0,
			0,
			false,
		},
		{"id 1",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 1, Data: []byte{3, 2, 1}},
			},
			0,
			0,
			0,
			false,
		},
		{"concat16",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 8, Data: []byte{4, 3, 2, 1}},
			},
			1027,
			2,
			1,
			true,
		},
		{"short concat8",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 0, Data: []byte{2, 1}},
			},
			0,
			0,
			0,
			false,
		},
		{"short concat16",
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 8, Data: []byte{3, 2, 1}},
			},
			0,
			0,
			0,
			false,
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			segments, seqno, mref, ok := p.udh.ConcatInfo16()
			assert.Equal(t, p.ok, ok)
			assert.Equal(t, p.segments, segments)
			assert.Equal(t, p.seqno, seqno)
			assert.Equal(t, p.mref, mref)
		}
		t.Run(p.name, f)

	}
}

func TestDecodeUserData(t *testing.T) {
	// Also tests NewUDDecoder, AddLockingCharset and AddShiftCharset
	patterns := []struct {
		name    string
		ud      tpdu.UserData
		udh     tpdu.UserDataHeader
		alpha   tpdu.Alphabet
		options []tpdu.UDDecodeOption
		msg     []byte
		err     error
	}{
		{"empty",
			nil,
			nil,
			0,
			nil,
			nil,
			nil},
		{"message 7bit",
			[]byte("message\x10"),
			nil,
			tpdu.Alpha7Bit,
			nil,
			[]byte("messageΔ"),
			nil,
		},
		{"message reserved",
			[]byte("message\x10"),
			nil,
			tpdu.AlphaReserved,
			nil,
			[]byte("messageΔ"),
			nil,
		},
		{"message 7bit esc",
			[]byte("message\x1b"),
			nil,
			tpdu.Alpha7Bit,
			nil,
			[]byte("message "),
			nil,
		},
		{"message 7bit locking",
			[]byte("\x01\x02\x03"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 25, Data: []byte{byte(charset.Kannada)}},
			},
			tpdu.Alpha7Bit,
			[]tpdu.UDDecodeOption{tpdu.WithLockingCharset(charset.Kannada)},
			[]byte("\u0c82\u0c83\u0c85"),
			nil,
		},
		{"message 7bit shift", []byte("\x1b\x1e\x1b\x1f\x1b\x20"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 24, Data: []byte{byte(charset.Kannada)}},
			},
			tpdu.Alpha7Bit,
			[]tpdu.UDDecodeOption{tpdu.WithShiftCharset(charset.Kannada)},
			[]byte("\u0ce8\u0ce9\u0cea"),
			nil,
		},
		{"message 8bit",
			[]byte("message\x1b"),
			nil,
			tpdu.Alpha8Bit,
			nil,
			[]byte("message\x1b"),
			nil,
		},
		{"euro",
			[]byte("\x1be"),
			nil,
			tpdu.Alpha7Bit,
			nil,
			[]byte("€"),
			nil,
		},
		{"grin",
			[]byte{0xd8, 0x3d, 0xde, 0x01},
			nil,
			tpdu.AlphaUCS2,
			nil,
			[]byte("😁"),
			nil,
		},
		// repeat the GSM7 Kannada tests without charset to force decoding to
		// fallback to default
		{"message 7bit locking defaulted",
			[]byte("\x01\x02\x03"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 25, Data: []byte{byte(charset.Kannada)}},
			},
			tpdu.Alpha7Bit,
			nil,
			[]byte("£$¥"),
			nil,
		},
		{"message 7bit shift defaulted",
			[]byte("\x1b\x1e\x1b\x1f\x1b\x20"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 24, Data: []byte{byte(charset.Kannada)}},
			},
			tpdu.Alpha7Bit,
			nil,
			[]byte("ßÉ "),
			nil,
		},
		// error tests
		{"dangling surrogate",
			[]byte{0xd8, 0x3d, 0xde, 0x01, 0xd8, 0x3d},
			nil,
			tpdu.AlphaUCS2,
			nil,
			[]byte("😁"), ucs2.ErrDanglingSurrogate([]byte{0xd8, 0x3d}),
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			msg, err := tpdu.DecodeUserData(p.ud, p.udh, p.alpha, p.options...)
			require.Equal(t, p.err, err)
			assert.Equal(t, p.msg, msg)
		}
		t.Run(p.name, f)
	}
}

func TestDecodeUserDataAllCharsets(t *testing.T) {
	// Tests decode with AddAllCharsets
	patterns := []struct {
		name    string
		ud      tpdu.UserData
		udh     tpdu.UserDataHeader
		msg     []byte
		options []tpdu.UDDecodeOption
		err     error
	}{
		{"empty",
			nil,
			nil,
			nil,
			nil,
			nil,
		},
		{"message 7bit",
			[]byte("message\x10"),
			nil,
			[]byte("messageΔ"),
			nil,
			nil,
		},
		{"message reserved",
			[]byte("message\x10"),
			nil,
			[]byte("messageΔ"),
			nil,
			nil,
		},
		{"message 7bit esc",
			[]byte("message\x1b"),
			nil,
			[]byte("message "),
			nil,
			nil,
		},
		{"message 7bit locking all cs",
			[]byte("\x01\x02\x03"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 25, Data: []byte{byte(charset.Kannada)}},
			},
			[]byte("\u0c82\u0c83\u0c85"),
			[]tpdu.UDDecodeOption{
				tpdu.WithAllCharsets,
			},
			nil,
		},
		{"message 7bit locking kannada",
			[]byte("\x01\x02\x03"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 25, Data: []byte{byte(charset.Kannada)}},
			},
			[]byte("\u0c82\u0c83\u0c85"),
			[]tpdu.UDDecodeOption{
				tpdu.WithCharset(charset.Kannada),
			},
			nil,
		},
		{"message 7bit shift all cs",
			[]byte("\x1b\x1e\x1b\x1f\x1b\x20"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 24, Data: []byte{byte(charset.Kannada)}},
			},
			[]byte("\u0ce8\u0ce9\u0cea"),
			[]tpdu.UDDecodeOption{
				tpdu.WithAllCharsets,
			},
			nil,
		},
		{"message 7bit shift kannada",
			[]byte("\x1b\x1e\x1b\x1f\x1b\x20"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 24, Data: []byte{byte(charset.Kannada)}},
			},
			[]byte("\u0ce8\u0ce9\u0cea"),
			[]tpdu.UDDecodeOption{
				tpdu.WithCharset(charset.Kannada),
			},
			nil,
		},
		{"euro",
			[]byte("\x1be"),
			nil,
			[]byte("€"),
			nil,
			nil,
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			msg, err := tpdu.DecodeUserData(p.ud, p.udh, tpdu.Alpha7Bit, p.options...)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.msg, msg)
		}
		t.Run(p.name, f)
	}
}

func TestEncodeUserData(t *testing.T) {
	// Also tests WithAllCharsets, WithCharset, WithLockingCharset and
	// WithShiftCharset.
	patterns := []struct {
		name    string
		ud      tpdu.UserData
		udh     tpdu.UserDataHeader
		alpha   tpdu.Alphabet
		options []tpdu.UDEncodeOption
		msg     []byte
	}{
		{"empty",
			nil,
			nil,
			0,
			nil,
			nil,
		},
		{"message 7bit",
			[]byte("message\x10"),
			nil,
			tpdu.Alpha7Bit,
			nil,
			[]byte("messageΔ"),
		},
		{"message 7bit all cs",
			[]byte("message\x10"),
			nil,
			tpdu.Alpha7Bit,
			[]tpdu.UDEncodeOption{
				tpdu.WithAllCharsets,
			},
			[]byte("messageΔ"),
		},
		{"message 7bit locking all cs",
			[]byte("\x01\x02\x03"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 25, Data: []byte{byte(charset.Kannada)}},
			},
			tpdu.Alpha7Bit,
			[]tpdu.UDEncodeOption{
				tpdu.WithAllCharsets,
			},
			[]byte("\u0c82\u0c83\u0c85"),
		},
		{"message 7bit shift all cs",
			[]byte("\x1b\x1e\x1b\x1f\x1b\x20"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 24, Data: []byte{byte(charset.Kannada)}},
			},
			tpdu.Alpha7Bit,
			[]tpdu.UDEncodeOption{
				tpdu.WithAllCharsets,
			},
			[]byte("\u0ce8\u0ce9\u0cea"),
		},
		{"message 7bit kannada",
			[]byte("\x01\x02\x03"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 25, Data: []byte{byte(charset.Kannada)}},
			},
			tpdu.Alpha7Bit,
			[]tpdu.UDEncodeOption{
				tpdu.WithCharset(charset.Kannada),
			},
			[]byte("\u0c82\u0c83\u0c85"),
		},
		{"message 7bit locking kannada",
			[]byte("\x01\x02\x03"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 25, Data: []byte{byte(charset.Kannada)}},
			},
			tpdu.Alpha7Bit,
			[]tpdu.UDEncodeOption{
				tpdu.WithLockingCharset(charset.Kannada),
			},
			[]byte("\u0c82\u0c83\u0c85"),
		},
		{"message 7bit shift kannada",
			[]byte("\x1b\x1e\x1b\x1f\x1b\x20"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 24, Data: []byte{byte(charset.Kannada)}},
			},
			tpdu.Alpha7Bit,
			[]tpdu.UDEncodeOption{
				tpdu.WithShiftCharset(charset.Kannada),
			},
			[]byte("\u0ce8\u0ce9\u0cea"),
		},
		{"message 7bit locking and shift urdu",
			[]byte("hello \x07\x1b\x2a"),
			tpdu.UserDataHeader{
				tpdu.InformationElement{ID: 25, Data: []byte{byte(charset.Urdu)}},
				tpdu.InformationElement{ID: 24, Data: []byte{byte(charset.Urdu)}},
			},
			tpdu.Alpha7Bit,
			[]tpdu.UDEncodeOption{
				tpdu.WithLockingCharset(charset.Urdu),
				tpdu.WithShiftCharset(charset.Urdu),
			},
			[]byte("hello ت؎"),
		},
		{"euro",
			[]byte("\x1be"),
			nil,
			tpdu.Alpha7Bit,
			nil,
			[]byte("€"),
		},
		{"grin",
			[]byte{0xd8, 0x3d, 0xde, 0x01},
			nil,
			tpdu.AlphaUCS2,
			nil,
			[]byte("😁"),
		},
		{"grin all cs",
			[]byte{0xd8, 0x3d, 0xde, 0x01},
			nil,
			tpdu.AlphaUCS2,
			[]tpdu.UDEncodeOption{
				tpdu.WithAllCharsets,
			},
			[]byte("😁"),
		},
		// repeat the GSM7 Kannada tests without charset to force encoding to UCS2
		{"message ucs2 locking",
			[]byte{0x0c, 0x82, 0x0c, 0x83, 0x0c, 0x85},
			nil,
			tpdu.AlphaUCS2,
			nil,
			[]byte("\u0c82\u0c83\u0c85"),
		},
		{"message ucs2 shift",
			[]byte{0x0c, 0xe8, 0x0c, 0xe9, 0x0c, 0xea},
			nil,
			tpdu.AlphaUCS2,
			nil,
			[]byte("\u0ce8\u0ce9\u0cea"),
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			ud, udh, alpha := tpdu.EncodeUserData(p.msg, p.options...)
			assert.Equal(t, p.ud, ud)
			assert.Equal(t, p.udh, udh)
			assert.Equal(t, p.alpha, alpha)
		}
		t.Run(p.name, f)
	}
}
