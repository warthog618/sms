// SPDX-License-Identifier: MIT
//
// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.

package tpdu

// FirstOctet is the first byte of a SMS-TPDU.
type FirstOctet byte

// LP returns true if the TP-LP flag is set.
func (f FirstOctet) LP() bool {
	return f&FoLP != 0
}

// MMS returns true if the TP-MMS flag is set.
func (f FirstOctet) MMS() bool {
	return f&FoMMS != 0
}

// MTI returns the message type field.
func (f FirstOctet) MTI() MessageType {
	return MessageType(f & FoMTIMask)
}

// RD returns true if the TP-RD flag is set.
func (f FirstOctet) RD() bool {
	return f&FoRD != 0
}

// RP returns true if the TP-RP flag is set.
func (f FirstOctet) RP() bool {
	return f&FoRP != 0
}

// SRI returns true if the TP-SRI flag is set.
func (f FirstOctet) SRI() bool {
	return f&FoSRI != 0
}

// SRR returns true if the TP-SRR flag is set.
func (f FirstOctet) SRR() bool {
	return f&FoSRR != 0
}

// SRQ returns true if the TP-SRQ flag is set.
func (f FirstOctet) SRQ() bool {
	return f&FoSRQ != 0
}

// VPF returns the TP-VPF field.
func (f FirstOctet) VPF() ValidityPeriodFormat {
	return ValidityPeriodFormat((f & FoVPFMask) >> FoVPFShift)
}

// WithMTI returns a FirstOctet with the TP-MTI field set.
func (f FirstOctet) WithMTI(mti MessageType) FirstOctet {
	f &^= FoMTIMask
	f |= FirstOctet(mti << FoMTIShift)
	return f
}

// WithVPF returns a FirstOctet with the TP-VPF field set.
func (f FirstOctet) WithVPF(vpf ValidityPeriodFormat) FirstOctet {
	f &^= FoVPFMask
	f |= FirstOctet(vpf << FoVPFShift)
	return f
}

// UDHI returns true if the TP-UDHI flag is set.
func (f FirstOctet) UDHI() bool {
	return f&FoUDHI != 0
}

const (
	// FirstOctet bit fields

	// FoMTIMask masks the bit for the TP-MTI field
	FoMTIMask = 0x3

	// FoMTIShift defines the shift required to move the MTI field to/from bit 0
	FoMTIShift = 0

	// FoMMS defines the TP-MMS More Messages to Send bit
	//
	// Only applies to SMS-DELIVER and SMS-STATUS-REPORT
	FoMMS = 0x4

	// FoRD defines the TP-RD Reject Duplicates bit
	//
	// Only applies to SMS-SUBMIT
	FoRD = 0x4

	// FoLP defines the TP-LP Loop Prevention bit
	//
	// Only applies to SMS-DELIVER and SMS-STATUS-REPORT
	FoLP = 0x8

	// FoVPFMask masks the bit for the TP-VPF field
	//
	// Only applies to SMS-SUBMIT
	FoVPFMask = 0x18

	// FoVPFShift defines the shift required to move the VPF field to/from bit 0
	FoVPFShift = 3

	// FoSRI defines the TP-SRI bit
	//
	// Only applies to SMS-DELIVER
	FoSRI = 0x20

	// FoSRR defines the TP-SRR bit
	//
	// Only applies to the SMS-SUBMIT and SMS-COMMAND
	FoSRR = 0x20

	// FoSRQ defines the TP-SRQ bit
	//
	// Only applies to the SMS-STATUS-REPORT
	FoSRQ = 0x20

	// FoUDHI defines the TP-UDHI bit
	FoUDHI = 0x40

	// FoRP defines the TP-RP bit
	//
	// Only applies to the SMS-SUBMIT and SMS-DELIVER
	FoRP = 0x80
)
