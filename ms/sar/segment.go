// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sar

import (
	"encoding/binary"
	"sync"

	"github.com/warthog618/sms/encoding/tpdu"
)

// Segmenter segments a large outgoing message into the set of Submit TPDUs
// required to contain it.
type Segmenter struct {
	mutex    sync.Mutex // covers msgCount and wide
	msgCount int
	wide     bool
}

// NewSegmenter creates a Segmenter.
func NewSegmenter() *Segmenter {
	return new(Segmenter)
}

// SetWide sets the Segmenter wide flag.
// If true then the segmenter will use 16bit message references instead of 8bit.
func (s *Segmenter) SetWide(w bool) {
	s.mutex.Lock()
	s.wide = w
	s.mutex.Unlock()
}

// Segment returns the set of SMS-Submit TPDUs required to transmit the message
// using the given alphabet.
// A template for the SMS-Submit TPDUs is passed in, and provides all the
// fields in the resulting TPDUs, other than the UD, which is populated using
// the message.  For multi-part messages, the UDH provided in the template is
// extended with a concatenation IE.
// The template UDH must not contain a concatenation IE (ID 0) or the resulting
// TPDUs will be non-conformant.
func (s *Segmenter) Segment(msg []byte, t *tpdu.Submit) []tpdu.Submit {
	if len(msg) == 0 || t == nil {
		return nil
	}
	alpha, _ := t.Alphabet()
	udhl := t.UDH.UDHL()
	bs := maxSML(t.MaxUDL(), udhl, alpha)
	if len(msg) <= bs {
		// single segment
		pdus := make([]tpdu.Submit, 1)
		pdus[0] = *t
		pdus[0].UD = msg
		return pdus
	}
	// allow for concat entry in UDH
	bs = maxSML(t.MaxUDL(), udhl+5, alpha)
	// any point checking for bs==0?
	var chunks [][]byte
	switch alpha {
	default: // default to 7Bit
		chunks = chunk7Bit(msg, bs)
	case tpdu.AlphaUCS2:
		chunks = chunkUCS2(msg, bs)
	case tpdu.Alpha8Bit:
		chunks = chunk8Bit(msg, bs)
	}
	count := len(chunks)
	pdus := make([]tpdu.Submit, count)
	s.mutex.Lock()
	s.msgCount++
	msgCount := s.msgCount
	wide := s.wide
	s.mutex.Unlock()
	for i := 0; i < count; i++ {
		sg := &pdus[i]
		*sg = *t
		ie := tpdu.InformationElement{}
		if wide {
			ie.ID = 8
			ie.Data = []byte{0, 0, byte(count), byte(i + 1)}
			binary.BigEndian.PutUint16(ie.Data, uint16(msgCount))
		} else {
			ie.ID = 0
			ie.Data = []byte{byte(msgCount), byte(count), byte(i + 1)}
		}
		sg.SetUDH(append(t.UDH, ie))
		sg.UD = chunks[i]
	}
	return pdus
}

const (
	esc byte = 0x1b
)

// chunk7Bit splits a GSM7 message into chunks that are not larger than bs.
// Escaped characters are not split across blocks, so the resulting blocks may
// be one septet shorter than bs.
func chunk7Bit(msg []byte, bs int) [][]byte {
	if len(msg) == 0 {
		return nil
	}
	count := 1 + len(msg)/bs
	chunks := make([][]byte, 0, count)
	bstart := 0
	bend := bs
	for bend < len(msg) {
		// don't split escapes
		if msg[bend-1] == esc && msg[bend-2] != esc {
			bend--
		}
		chunks = append(chunks, msg[bstart:bend])
		bstart = bend
		bend = bstart + bs
	}
	chunks = append(chunks, msg[bstart:])
	return chunks
}

// chunk8Bit splits a raw 8bit message into chunks that are bs, except for the
// last segment which contains any residual bytes.
func chunk8Bit(msg []byte, bs int) [][]byte {
	if len(msg) == 0 {
		return nil
	}
	count := 1 + len(msg)/bs
	chunks := make([][]byte, 0, count)
	bstart := 0
	bend := bs
	for bend < len(msg) {
		chunks = append(chunks, msg[bstart:bend])
		bstart = bend
		bend = bstart + bs
	}
	chunks = append(chunks, msg[bstart:])
	return chunks
}

const (
	surrHighStart = 0xd800
	surrLowStart  = 0xdc00
)

// chunkUCS2 splits a UCS2/UTF-16 message into chunks that are not larger than bs.
// bs should be even, but if odd is reduced by one.
// To allow for reassemblers that cannot handle split surrogate pairs, they are
// not split during chunking, so the resulting blocks may be slightly smaller
// than bs whenever a surrogate pair would span a block boundary.
// While the msg should have even length for UCS2, the chunker does not enforce
// this, and if an odd length message is presented then the final chunk will
// have an odd length.
func chunkUCS2(msg []byte, bs int) [][]byte {
	if len(msg) == 0 {
		return nil
	}
	bs = bs &^ 0x1
	// rough count of blocks - may be off due to not splitting surrogates, but
	// not worth working out the precise count in advance.
	count := 1 + len(msg)/bs
	chunks := make([][]byte, 0, count)
	bstart := 0
	bend := bstart + bs
	for bend < len(msg) {
		// check last uint16 is a high surrogate, if so then leave for later
		r := binary.BigEndian.Uint16(msg[bend-2 : bend])
		if surrHighStart <= r && r < surrLowStart {
			bend = bend - 2
		}
		chunks = append(chunks, msg[bstart:bend])
		bstart = bend
		bend = bstart + bs
	}
	chunks = append(chunks, msg[bstart:])
	return chunks
}

// maxSML returns the block size for the SM in concatentated SMSs.
// For 8bit and UCS-2 it returns the number of bytes.
// For 7bit it returns the number of septets, though ,as the 7bit is unpacked
// at this stage, it also corresponds to the number of bytes.
func maxSML(maxUDL, udhl int, alpha tpdu.Alphabet) int {
	bs := maxUDL
	if alpha == tpdu.Alpha7Bit {
		// work in septets
		bs = (bs * 8) / 7
		if udhl == 0 {
			return bs
		}
		// remove septets used by UDH, including UDHL and fill bits
		bs = bs - ((udhl+1)*8+6)/7
		return bs
	}
	if udhl > 0 {
		bs = bs - udhl - 1
	}
	if alpha == tpdu.AlphaUCS2 {
		bs = bs &^ 0x1
	}
	return bs
}
