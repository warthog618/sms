// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package message_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/ms/message"
)

func TestNewReassembler(t *testing.T) {
	d := MockDecoder{}
	c := MockCollector{}
	r := message.NewReassembler(&d, &c)
	if r == nil {
		t.Fatalf("failed to create Reassembler")
	}
}

func TestReassemblerClose(t *testing.T) {
	closed := false
	c := MockCollector{CloseFunc: func() {
		closed = true
	}}
	d := MockDecoder{}
	r := message.NewReassembler(&d, &c)
	r.Close()
	if closed == false {
		t.Errorf("didn't call Collector.Close")
	}
}

func TestReassemble(t *testing.T) {
	patterns := []struct {
		name       string
		in         []byte
		segments   []*tpdu.Deliver
		collectErr error
		decoded    []byte
		decodeErr  error
		msg        *message.Message
		err        error
	}{
		{"empty", nil, nil, nil, nil, nil, nil, tpdu.DecodeError("firstOctet", 0, tpdu.ErrUnderflow)},
		{"single segment", []byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71, 0x32, 0x20,
			0x05, 0x23, 0x08, 0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3},
			[]*tpdu.Deliver{{}}, nil,
			[]byte("message"), nil,
			&message.Message{Msg: "message", Number: "", TPDUs: []*tpdu.Deliver{{}}}, nil},
		{"partial message", []byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71, 0x32, 0x20,
			0x05, 0x23, 0x08, 0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3},
			nil, nil,
			nil, nil,
			nil, nil},
		{"collect error", []byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71, 0x32, 0x20,
			0x05, 0x23, 0x08, 0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3},
			nil, errors.New("collect error"),
			nil, nil,
			nil, errors.New("collect error")},
		{"decode error", []byte{0x04, 0x04, 0x91, 0x36, 0x19, 0x00, 0x00, 0x51, 0x50, 0x71, 0x32, 0x20,
			0x05, 0x23, 0x08, 0xC8, 0x30, 0x3A, 0x8C, 0x0E, 0xA3, 0xC3},
			[]*tpdu.Deliver{{}}, nil,
			nil, errors.New("decode error"),
			nil, errors.New("decode error")},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			c := MockCollector{CollectFunc: func(pdu *tpdu.Deliver) (d []*tpdu.Deliver, err error) {
				return p.segments, p.collectErr
			}}
			d := MockDecoder{DecodeFunc: func(ud tpdu.UserData, udh tpdu.UserDataHeader, alpha tpdu.Alphabet) ([]byte, error) {
				return p.decoded, p.decodeErr
			}}
			r := message.NewReassembler(&d, &c)
			if r == nil {
				t.Fatalf("failed to create Reassembler")
			}
			m, err := r.Reassemble(p.in)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.msg, m)
		}
		t.Run(p.name, f)
	}
}

// TestConcatenate tests the dangling surrogate section of concatenate.
// Other functionality is covered implicitly by TestReassemble.
func TestConcatenate(t *testing.T) {
	// Two UCS-2 TPDUs to be concatenated.  The first ends with the first surrogate
	// of a surrogate pair.  In total there are three surrogate pairs - each one
	// an emoticon - 😁.
	u1 := tpdu.NewDeliver()
	dcs, _ := tpdu.DCS(0).WithAlphabet(tpdu.AlphaUCS2)
	oa := tpdu.Address{TOA: 0x91, Addr: "1234"}
	u1.SetDCS(dcs)
	u1.SetOA(oa)
	u1.SetUD([]byte{0xd8, 0x3d, 0xde, 0x01, 0xd8, 0x3d})
	u2 := tpdu.NewDeliver()
	u2.SetDCS(dcs)
	u2.SetOA(oa)
	u2.SetUD([]byte{0xde, 0x01, 0xd8, 0x3d, 0xde, 0x01})
	d, _ := tpdu.NewUDDecoder()
	c := message.NewConcatenator(d)
	m, err := c.Concatenate([]*tpdu.Deliver{u1, u2})
	assert.Equal(t, nil, err)
	expected := &message.Message{Msg: "😁😁😁", Number: "+1234", TPDUs: []*tpdu.Deliver{u1, u2}}
	assert.Equal(t, expected, m)
}

type MockCollector struct {
	CloseFunc   func()
	CollectFunc func(pdu *tpdu.Deliver) (d []*tpdu.Deliver, err error)
}

func (c *MockCollector) Collect(pdu *tpdu.Deliver) (d []*tpdu.Deliver, err error) {
	return c.CollectFunc(pdu)
}

func (c *MockCollector) Close() {
	c.CloseFunc()
}

type MockDecoder struct {
	DecodeFunc func(ud tpdu.UserData, udh tpdu.UserDataHeader, alpha tpdu.Alphabet) ([]byte, error)
}

func (d *MockDecoder) Decode(ud tpdu.UserData, udh tpdu.UserDataHeader, alpha tpdu.Alphabet) ([]byte, error) {
	return d.DecodeFunc(ud, udh, alpha)
}
