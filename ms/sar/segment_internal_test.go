// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/warthog618/sms/encoding/tpdu"
)

func TestChunk7Bit(t *testing.T) {
	patterns := []struct {
		name string
		in   []byte
		bs   int
		out  [][]byte
	}{
		{"empty", nil, 2, nil},
		{"integral", []byte{1, 2, 3, 4}, 2, [][]byte{{1, 2}, {3, 4}}},
		{"residual", []byte{1, 2, 3, 4}, 3, [][]byte{{1, 2, 3}, {4}}},
		{"three", []byte{1, 2, 3, 4, 5, 6, 7, 8}, 3, [][]byte{{1, 2, 3}, {4, 5, 6}, {7, 8}}},
		{"escaped", []byte{1, 2, 0x1b, 4, 5, 6, 7, 8}, 3, [][]byte{{1, 2}, {0x1b, 4, 5}, {6, 7, 8}}},
		{"double escaped", []byte{1, 0x1b, 0x1b, 4, 5, 6, 7, 8}, 3, [][]byte{{1, 0x1b, 0x1b}, {4, 5, 6}, {7, 8}}},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out := chunk7Bit(p.in, p.bs)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}

func TestChunk8Bit(t *testing.T) {
	patterns := []struct {
		name string
		in   []byte
		bs   int
		out  [][]byte
	}{
		{"empty", nil, 2, nil},
		{"integral", []byte{1, 2, 3, 4}, 2, [][]byte{{1, 2}, {3, 4}}},
		{"residual", []byte{1, 2, 3, 4}, 3, [][]byte{{1, 2, 3}, {4}}},
		{"three", []byte{1, 2, 3, 4, 5, 6, 7, 8}, 3, [][]byte{{1, 2, 3}, {4, 5, 6}, {7, 8}}},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out := chunk8Bit(p.in, p.bs)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}

func TestChunkUCS2(t *testing.T) {
	patterns := []struct {
		name string
		in   []byte
		bs   int
		out  [][]byte
	}{
		{"empty", nil, 2, nil},
		{"integral", []byte{1, 2, 3, 4}, 2, [][]byte{{1, 2}, {3, 4}}},
		{"odd bs", []byte{1, 2, 3, 4}, 3, [][]byte{{1, 2}, {3, 4}}},
		{"three", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, 4, [][]byte{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10}}},
		{"surrogage", []byte{1, 2, 0xd8, 4, 5, 6, 7, 8, 9, 10}, 4, [][]byte{{1, 2}, {0xd8, 4, 5, 6}, {7, 8, 9, 10}}},
		{"odd msg", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9}, 4, [][]byte{{1, 2, 3, 4}, {5, 6, 7, 8}, {9}}},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out := chunkUCS2(p.in, p.bs)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}

func TestMaxSML(t *testing.T) {
	patterns := []struct {
		name   string
		maxUDL int
		udhl   int
		alpha  tpdu.Alphabet
		out    int
	}{
		{"empty", 0, 0, 0, 0},
		{"six 7bit", 6, 0, tpdu.Alpha7Bit, 6},
		{"seven 7bit", 7, 0, tpdu.Alpha7Bit, 8},
		{"eight 7bit", 8, 0, tpdu.Alpha7Bit, 9},
		{"six 7bit with udh", 6, 3, tpdu.Alpha7Bit, 1},
		{"seven 7bit with udh", 7, 3, tpdu.Alpha7Bit, 3},
		{"eight 7bit with udh", 8, 3, tpdu.Alpha7Bit, 4},
		{"odd 8bit", 13, 0, tpdu.Alpha8Bit, 13},
		{"even 8bit", 12, 0, tpdu.Alpha8Bit, 12},
		{"odd 8bit with udh", 13, 3, tpdu.Alpha8Bit, 9},
		{"even 8bit with udh", 12, 3, tpdu.Alpha8Bit, 8},
		{"odd ucs2", 13, 0, tpdu.AlphaUCS2, 12},
		{"even ucs2", 12, 0, tpdu.AlphaUCS2, 12},
		{"odd ucs2 with udh", 13, 3, tpdu.AlphaUCS2, 8},
		{"even ucs2 with udh", 12, 3, tpdu.AlphaUCS2, 8},
	}
	for _, p := range patterns {
		f := func(t *testing.T) {
			out := maxSML(p.maxUDL, p.udhl, p.alpha)
			if out != p.out {
				t.Errorf("failed: expected %d, got %d", p.out, out)
			}
		}
		t.Run(p.name, f)
	}
}
