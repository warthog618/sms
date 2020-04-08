// SPDX-License-Identifier: MIT
//
// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.

package tpdu

import "strings"

// PI is the parameter indicator bitfield.
type PI byte

// PID returns true if a PID field is present in the TPDU.
func (p PI) PID() bool {
	return p&PiPID != 0
}

// DCS returns true if a DCS field is present in the TPDU.
func (p PI) DCS() bool {
	return p&PiDCS != 0
}

// UDL returns true if a UDL, and hence a UD, field is present in the TPDU.
func (p PI) UDL() bool {
	return p&PiUDL != 0
}

func (p PI) String() string {
	if p == 0 {
		return "0"
	}
	elems := []string{}
	if p.PID() {
		elems = append(elems, "PID")
	}
	if p.DCS() {
		elems = append(elems, "DCS")
	}
	if p.UDL() {
		elems = append(elems, "UDL")
	}
	return strings.Join(elems, "|")
}

const (
	// PI bit fields

	// PiPID indicates a TP-PID field is present in the TPDU
	PiPID = 1 << iota

	// PiDCS indicates a TP-DCS field is present in the TPDU
	PiDCS

	// PiUDL indicates a TP-UDL field is present in the TPDU
	PiUDL
)
