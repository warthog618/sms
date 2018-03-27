// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package sar

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/warthog618/sms/encoding/tpdu"
)

// Collector contains reassembly pipes that buffer concatenated TPDUs until
// a full set is available to be concatenated.
type Collector struct {
	sync.Mutex // covers pipes and closing closed
	pipes      map[string]*pipe
	duration   time.Duration
	closed     chan struct{}
	asyncError func(error)
}

// NewCollector creates a Collector.
// The asyncError function is called when a reassembly fails asynchronously.
// The asyncError function must be safe to be called from multiple goroutines.
func NewCollector(d time.Duration, asyncError func(error)) *Collector {
	return &Collector{
		sync.Mutex{},
		make(map[string]*pipe),
		d,
		make(chan struct{}),
		asyncError,
	}
}

// Close shuts down the Collector and all active pipes.
func (c *Collector) Close() {
	c.Lock()
	select {
	case <-c.closed:
	default:
		close(c.closed)
		for _, p := range c.pipes {
			p.cleanup.Stop()
		}
	}
	c.Unlock()
}

// Collect adds a TPDU to the collection.
// If all the components of a concatenated TPDU are available then they are returned.
func (c *Collector) Collect(pdu *tpdu.Deliver) (d []*tpdu.Deliver, err error) {
	segments, seqno, mref, ok := pdu.UDH().ConcatInfo()
	if !ok || segments < 2 {
		// short circuit single segment - no need for a pipe
		return []*tpdu.Deliver{pdu}, nil
	}
	if seqno < 1 || seqno > segments {
		return nil, ErrReassemblyInconsistency
	}
	oa := pdu.OA()
	key := fmt.Sprintf("%02x:%s:%d:%d", oa.TOA, oa.Addr, mref, segments)
	c.Lock()
	defer c.Unlock()
	select {
	case <-c.closed:
		return nil, ErrClosed
	default:
	}
	p, ok := c.pipes[key]
	if ok {
		if p.segments[seqno-1] != nil {
			return nil, ErrDuplicateSegment
		}
		if !p.cleanup.Stop() {
			// timer has fired, but cleanup hasn't been performed yet - so need a new pipe
			ok = false
		}
	}
	if !ok {
		p = &pipe{nil, make([]*tpdu.Deliver, segments), 0}
		c.pipes[key] = p
	}
	p.segments[seqno-1] = pdu
	p.frags++
	if p.frags == segments {
		delete(c.pipes, key)
		return p.segments, nil
	}
	p.cleanup = time.AfterFunc(c.duration, func() {
		c.Lock()
		m := c.pipes[key]
		if m == p {
			delete(c.pipes, key)
		}
		c.Unlock()
		c.asyncError(ErrExpired{p.segments})
	})
	return nil, err
}

// pipe is a buffer that contains the individual TPDUs in a concatenation set
// until the complete set is available or the reassembly times out.
type pipe struct {
	cleanup  *time.Timer
	segments []*tpdu.Deliver
	frags    int
}

// ErrExpired indicates that a reassembly has timed out.
// The segments of the aborted reassembly are returned in the error.
type ErrExpired struct {
	T []*tpdu.Deliver
}

func (e ErrExpired) Error() string {
	return fmt.Sprintf("sar: timed out reassembling %v", e.T)
}

var (
	// ErrClosed indicates that the collector has been closed and is no longer
	// accepting PDUs.
	ErrClosed = errors.New("closed")
	// ErrDuplicateSegment indicates a segment has arrived for a reassembly
	// that already has that segment.
	// The segments are duplicates in terms of their concatentation information.
	// They may differ in other fields, particularly UD, but those fields cannot
	// be used to determine which of the two may better fit the reassembly, so
	// the first is kept and the second discarded.
	ErrDuplicateSegment = errors.New("duplcate segment")
	// ErrReassemblyInconsistency indicates a segment has arrived for a reassembly
	// that has a seqno greater than the number of segments in the reassembly.
	ErrReassemblyInconsistency = errors.New("reassembly inconsistency")
)
