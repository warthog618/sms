// SPDX-License-Identifier: MIT
//
// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.

package tpdu

// Option applies a construction option to a TPDU.
type Option interface {
	ApplyTPDUOption(*TPDU) error
}

// DAOption specifies the DA for the TPDU.
type DAOption struct {
	addr Address
}

// ApplyTPDUOption applies the DA to the TPDU.
func (o DAOption) ApplyTPDUOption(t *TPDU) error {
	t.DA = o.addr
	return nil
}

// WithDA creates a DAOption to apply to a TPDU.
func WithDA(addr Address) DAOption {
	return DAOption{addr}
}

// OAOption specifies the OA for the TPDU.
type OAOption struct {
	addr Address
}

// ApplyTPDUOption applies the OA to the TPDU.
func (o OAOption) ApplyTPDUOption(t *TPDU) error {
	t.OA = o.addr
	return nil
}

// WithOA creates a OAOption to apply to a TPDU.
func WithOA(addr Address) OAOption {
	return OAOption{addr}
}

// UDHOption specifies the UDH for the TPDU.
type UDHOption struct {
	udh UserDataHeader
}

// ApplyTPDUOption applies the UDH to the TPDU.
func (o UDHOption) ApplyTPDUOption(t *TPDU) error {
	t.UDH = o.udh
	return nil
}

// WithUDH creates a UDHOption to apply to a TPDU.
func WithUDH(udh UserDataHeader) UDHOption {
	return UDHOption{udh}
}
