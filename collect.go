// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package sms

import (
	"fmt"
	"sync"
	"time"

	"github.com/warthog618/sms/encoding/tpdu"
)

// Collector contains reassembly pipes that buffer concatenated TPDUs until a
// full set is available to be concatenated.
type Collector struct {
	sync.Mutex    // covers pipes and closing closed
	pipes         map[string]*pipe
	closed        bool
	duration      time.Duration
	expiryHandler func([]*tpdu.TPDU)
}

// CollectorOption alters the behaviour of a Collector.
type CollectorOption interface {
	ApplyCollectorOption(*Collector)
}

type reassemblyTimeoutOption struct {
	d  time.Duration
	eh func([]*tpdu.TPDU)
}

func (o reassemblyTimeoutOption) ApplyCollectorOption(c *Collector) {
	c.duration = o.d
	c.expiryHandler = o.eh
}

// WithReassemblyTimeout limits the time allowed for a collection of TPDUs to
// be collected.
//
// If the timer expires before the collection is complete then the collected
// TPDUs are passed to the expiryHandler. The expiry handler can be nil in
// which case the collected TPDUs are simply discarded.
//
// A zero duration disables the timeout.
func WithReassemblyTimeout(d time.Duration, eh func([]*tpdu.TPDU)) CollectorOption {
	return reassemblyTimeoutOption{d, eh}
}

// NewCollector creates a Collector.
func NewCollector(options ...CollectorOption) *Collector {
	c := Collector{
		pipes: make(map[string]*pipe),
	}
	for _, o := range options {
		o.ApplyCollectorOption(&c)
	}
	return &c
}

// Close shuts down the Collector and all active pipes.
func (c *Collector) Close() {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return
	}
	c.closed = true
	for _, p := range c.pipes {
		if p.cleanup != nil {
			p.cleanup.Stop()
		}
	}
}

// Pipes returns a snapshot of the reassembly pipes.
//
// This is intended for diagnostics.
func (c *Collector) Pipes() map[string][]*tpdu.TPDU {
	c.Lock()
	m := map[string][]*tpdu.TPDU{}
	for k, v := range c.pipes {
		m[k] = v.segments
	}
	c.Unlock()
	return m
}

// Collect adds a TPDU to the collection.
//
// If all the components of a concatenated TPDU are available then they are
// returned.
func (c *Collector) Collect(pdu tpdu.TPDU) (d []*tpdu.TPDU, err error) {
	c.Lock()
	defer c.Unlock()
	if c.closed {
		return nil, ErrClosed
	}
	segments, seqno, concatRef, ok := pdu.ConcatInfo()
	if !ok || segments < 2 {
		// short circuit single segment - no need for a pipe
		return []*tpdu.TPDU{&pdu}, nil
	}
	if seqno < 1 || seqno > segments {
		return nil, ErrReassemblyInconsistency
	}
	key, err := pduKey(pdu, segments, concatRef)
	p, ok := c.pipes[key]
	if ok {
		if p.segments[seqno-1] != nil {
			return nil, ErrDuplicateSegment
		}
		if p.cleanup != nil && !p.cleanup.Stop() {
			// timer has fired, but cleanup hasn't been performed yet - so need
			// a new pipe
			ok = false
		}
	}
	if !ok {
		p = &pipe{nil, make([]*tpdu.TPDU, segments), 0}
		c.pipes[key] = p
	}
	p.segments[seqno-1] = &pdu
	p.frags++
	if p.frags == segments {
		delete(c.pipes, key)
		return p.segments, nil
	}
	if c.duration != 0 {
		p.cleanup = time.AfterFunc(c.duration, func() {
			c.Lock()
			m := c.pipes[key]
			if m == p {
				delete(c.pipes, key)
			}
			c.Unlock()
			if c.expiryHandler != nil {
				c.expiryHandler(p.segments)
			}
		})
	}
	return nil, err
}

func pduKey(pdu tpdu.TPDU, segments, concatRef int) (string, error) {
	st := pdu.SmsType()
	var key string
	switch st {
	case tpdu.SmsSubmit:
		key = fmt.Sprintf("%d:%02x:%s:%d:%d",
			st,
			pdu.DA.TOA,
			pdu.DA.Addr,
			concatRef,
			segments)
	case tpdu.SmsDeliver:
		key = fmt.Sprintf("%d:%02x:%s:%d:%d",
			st,
			pdu.OA.TOA,
			pdu.OA.Addr,
			concatRef,
			segments)
	default:
		return "", tpdu.ErrUnsupportedSmsType(st)
	}
	return key, nil
}

// pipe is a buffer that contains the individual TPDUs in a concatenation set
// until the complete set is available or the reassembly times out.
type pipe struct {
	cleanup  *time.Timer
	segments []*tpdu.TPDU
	frags    int
}
