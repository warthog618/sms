// SPDX-License-Identifier: MIT
//
// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.

package sms

import (
	"errors"
)

var (
	// ErrClosed indicates that the collector has been closed and is no longer
	// accepting PDUs.
	ErrClosed = errors.New("closed")
	// ErrDcsConflict indicates the required encoding for user data conflicts with the
	// encoding specified in the template TPDU DCS.
	ErrDcsConflict = errors.New("DCS conflict")
	// ErrDuplicateSegment indicates a segment has arrived for a reassembly
	// that already has that segment.
	// The segments are duplicates in terms of their concatentation information.
	// They may differ in other fields, particularly UD, but those fields
	// cannot be used to determine which of the two may better fit the
	// reassembly, so the first is kept and the second discarded.
	ErrDuplicateSegment = errors.New("duplicate segment")
	// ErrReassemblyInconsistency indicates a segment has arrived for a
	// reassembly that has a seqno greater than the number of segments in the
	// reassembly.
	ErrReassemblyInconsistency = errors.New("reassembly inconsistency")
)
