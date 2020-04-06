// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package sms_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/warthog618/sms"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestNewCollector(t *testing.T) {
	c := sms.NewCollector()
	assert.NotNil(t, c)
}

func TestCollectorClose(t *testing.T) {
	c := sms.NewCollector(sms.WithReassemblyTimeout(time.Minute, nil))
	require.NotNil(t, c)
	c.Close() // when open
	c.Close() // when closed
	c = sms.NewCollector(sms.WithReassemblyTimeout(time.Minute, nil))
	require.NotNil(t, c)
	d := tpdu.TPDU{}
	d.OA = tpdu.Address{Addr: "1234", TOA: 0x91}
	d.SetUDH(tpdu.UserDataHeader{tpdu.InformationElement{ID: 0, Data: []byte{5, 2, 1}}})
	c.Collect(d)
	c.Close() // with pipe active
}

func TestCollectorCollect(t *testing.T) {
	patterns := []struct {
		name string
		in   tpdu.TPDU
		out  []*tpdu.TPDU
		err  error
	}{
		// The patterns are Collected sequentially, so the return value depends
		// on the preceding set of PDUs, not just the individual in. The
		// resulting tests must be run as a complete set.
		{
			"deliver single segment",
			tpdu.TPDU{OA: tpdu.Address{Addr: "1234", TOA: 0x91}},
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				},
			},
			nil,
		},
		// 1 segment (shouldn't be seen in practice, but test in case)
		{
			"deliver one segment",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{1, 1, 1}},
				},
			},
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{1, 1, 1}},
					},
				},
			},
			nil,
		},
		{
			"deliver two a",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 1}},
				},
			},
			nil,
			nil,
		},
		{
			"deliver two b",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 2}},
				},
			},
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 1}},
					},
				},
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 2}},
					},
				},
			},
			nil,
		},
		{
			"deliver three a",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{3, 3, 1}},
				},
			},
			nil,
			nil,
		},
		{
			"deliver three b",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{3, 3, 2}},
				},
			},
			nil,
			nil,
		},
		{
			"deliver three c",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{3, 3, 3}},
				},
			},
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 3, 1}},
					},
				},
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 3, 2}},
					},
				},
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{3, 3, 3}},
					},
				},
			},
			nil,
		},
		{
			"jumbled a",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{4, 3, 1}},
				},
			},
			nil,
			nil,
		},
		{
			"jumbled c",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{4, 3, 3}},
				},
			},
			nil,
			nil,
		},
		{
			"duplicate",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{4, 3, 3}},
				},
			},
			nil,
			sms.ErrDuplicateSegment,
		},
		{
			"jumbled b",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{4, 3, 2}},
				},
			},
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{4, 3, 1}},
					},
				},
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{4, 3, 2}},
					},
				},
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{4, 3, 3}},
					},
				},
			},
			nil,
		},
		{
			"concurrent one a",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{5, 2, 1}},
				},
			},
			nil,
			nil,
		},
		{
			"concurrent two a",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{6, 2, 1}},
				},
			},
			nil,
			nil,
		},
		{
			"concurrent one b",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{5, 2, 2}},
				},
			},
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{5, 2, 1}},
					},
				},
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{5, 2, 2}},
					},
				},
			},
			nil,
		},
		{
			"concurrent two b",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{6, 2, 2}},
				},
			},
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{6, 2, 1}},
					},
				},
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{6, 2, 2}},
					},
				},
			},
			nil,
		},
		{
			"deliver 16bit concat a",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 8, Data: []byte{4, 4, 2, 1}},
				},
			},
			nil,
			nil,
		},
		{
			"deliver 16bit concat b",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 8, Data: []byte{4, 4, 2, 2}},
				},
			},
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 8, Data: []byte{4, 4, 2, 1}},
					},
				},
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 8, Data: []byte{4, 4, 2, 2}},
					},
				},
			},
			nil,
		},
		{
			"submit two a",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				DA:         tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 1}},
				},
			},
			nil,
			nil,
		},
		{
			"submit two b",
			tpdu.TPDU{
				Direction:  tpdu.MO,
				FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
				DA:         tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 2}},
				},
			},
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					Direction:  tpdu.MO,
					FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
					DA:         tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 1}},
					},
				},
				&tpdu.TPDU{
					Direction:  tpdu.MO,
					FirstOctet: tpdu.FirstOctet(tpdu.MtSubmit),
					DA:         tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 2}},
					},
				},
			},
			nil,
		},
		{
			"zero seqno",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{7, 2, 0}},
				},
			},
			nil,
			sms.ErrReassemblyInconsistency,
		},
		{
			"large seqno",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{8, 2, 3}},
				},
			},
			nil,
			sms.ErrReassemblyInconsistency,
		},
		{
			"deliverreport concat",
			tpdu.TPDU{
				Direction: tpdu.MO,
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{8, 2, 1}},
				},
			},
			nil,
			tpdu.ErrUnsupportedSmsType(tpdu.SmsDeliverReport),
		},
	}
	var ae []*tpdu.TPDU
	exph := func(pp []*tpdu.TPDU) {
		ae = pp
	}
	c := sms.NewCollector(sms.WithReassemblyTimeout(time.Minute, exph))
	require.NotNil(t, c)
	for _, p := range patterns {
		f := func(t *testing.T) {
			ae = nil
			out, err := c.Collect(p.in)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, out)
			assert.Nil(t, ae)
		}
		t.Run(p.name, f)
	}
	c.Close()
	patterns = []struct {
		name string
		in   tpdu.TPDU
		out  []*tpdu.TPDU
		err  error
	}{
		{
			"closed single segment",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
			},
			nil,
			sms.ErrClosed,
		},
		{
			"closed concat",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{4, 3, 1}},
				},
			},
			nil,
			sms.ErrClosed,
		},
	}
	for _, p := range patterns {
		ae = nil
		out, err := c.Collect(p.in)
		assert.Equal(t, p.err, err, p.name)
		assert.Equal(t, p.out, out, p.name)
		assert.Nil(t, ae, p.name)
	}
}

func TestCollectorReasemmblyTimeout(t *testing.T) {
	patterns := []struct {
		name string
		in   []tpdu.TPDU
	}{
		{
			"one",
			[]tpdu.TPDU{
				tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}},
					},
				},
			},
		},
		{
			"two",
			[]tpdu.TPDU{
				tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 2}},
					},
				},
				tpdu.TPDU{
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{1, 3, 1}},
					},
				},
			},
		},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			var done = make(chan struct{})
			var aechan = make(chan []*tpdu.TPDU)
			exph := func(tt []*tpdu.TPDU) {
				close(done)
				aechan <- tt
			}
			c := sms.NewCollector(sms.WithReassemblyTimeout(time.Millisecond, exph))
			pexp := make([]*tpdu.TPDU, len(p.in)+1)
			for i, s := range p.in {
				_, seqno, _, _ := s.ConcatInfo()
				pexp[seqno-1] = &p.in[i]
				m, err := c.Collect(s)
				assert.Nil(t, err)
				assert.Nil(t, m)
			}
			select {
			case <-done:
			case <-time.After(5 * time.Millisecond):
				t.Fatalf("didn't expire")
			}
			pipes := c.Pipes()
			assert.Zero(t, len(pipes))
			texp := <-aechan
			assert.Equal(t, pexp, texp)
		}
		t.Run(p.name, f)
	}
}

func TestCollectorPipes(t *testing.T) {
	c := sms.NewCollector()
	patterns := []struct {
		name string
		in   tpdu.TPDU
		m    []*tpdu.TPDU
		out  map[string][]*tpdu.TPDU
	}{
		{
			"deliver one a",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}},
				},
			},
			nil,
			map[string][]*tpdu.TPDU{
				"0:91:1234:1:2": []*tpdu.TPDU{
					&tpdu.TPDU{
						OA: tpdu.Address{Addr: "1234", TOA: 0x91},
						UDH: tpdu.UserDataHeader{
							tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}},
						},
					},
					nil,
				},
			},
		},
		{
			"deliver two b",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 2}},
				},
			},
			nil,
			map[string][]*tpdu.TPDU{
				"0:91:1234:1:2": []*tpdu.TPDU{
					&tpdu.TPDU{
						OA: tpdu.Address{Addr: "1234", TOA: 0x91},
						UDH: tpdu.UserDataHeader{
							tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}},
						},
					},
					nil,
				},
				"0:91:1234:2:2": []*tpdu.TPDU{
					nil,
					&tpdu.TPDU{
						OA: tpdu.Address{Addr: "1234", TOA: 0x91},
						UDH: tpdu.UserDataHeader{
							tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 2}},
						},
					},
				},
			},
		},
		{
			"deliver one b",
			tpdu.TPDU{
				OA: tpdu.Address{Addr: "1234", TOA: 0x91},
				UDH: tpdu.UserDataHeader{
					tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}},
				},
			},
			[]*tpdu.TPDU{
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 1}},
					},
				},
				&tpdu.TPDU{
					OA: tpdu.Address{Addr: "1234", TOA: 0x91},
					UDH: tpdu.UserDataHeader{
						tpdu.InformationElement{ID: 0, Data: []byte{1, 2, 2}},
					},
				},
			},
			map[string][]*tpdu.TPDU{
				"0:91:1234:2:2": []*tpdu.TPDU{
					nil,
					&tpdu.TPDU{
						OA: tpdu.Address{Addr: "1234", TOA: 0x91},
						UDH: tpdu.UserDataHeader{
							tpdu.InformationElement{ID: 0, Data: []byte{2, 2, 2}},
						},
					},
				},
			},
		},
	}
	for _, p := range patterns {

		m, err := c.Collect(p.in)
		assert.Nil(t, err, p.name)
		assert.Equal(t, p.m, m, p.name)
		out := c.Pipes()
		assert.Equal(t, p.out, out, p.name)
	}
}
