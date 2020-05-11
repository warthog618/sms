// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

package gsm7

const cr byte = 0x0d

// Pack7Bit packs an array of septets into an 8bit array as per the packing
// rules defined in 3GPP TS 23.038 Section 6.1.2.1
//
// The padBits is the number of bits of pad to place at the beginning of the
// packed array, as the packed septets may not start on an octet boundary.
//
// Packed arrays containing 8n or 8n-1 digits both return 8n septets. The
// caller must be aware of the number of expected digits in order to
// distinguish between a 0 septet ending the sequence in the 8n case, and 0
// padding in the 8n-1 case.
func Pack7Bit(u []byte, fillBits int) []byte {
	if len(u) == 0 {
		return append(u[:0:0], u...)
	}
	p := make([]byte, 0, (len(u)*7+7+fillBits)/8)
	var r, s byte
	rbits := uint(fillBits)
	for _, s = range u {
		if rbits == 0 {
			// no residual bits so not enough for a full octet
			r = s
			rbits = 7
			continue
		}
		r = (r | s<<rbits) & 0xff
		p = append(p, r)
		r = s >> (8 - rbits)
		rbits--
	}
	if rbits != 0 {
		p = append(p, r)
	}
	return p
}

// Unpack7Bit unpacks septets, packed into an 8bit array as per the packing
// rules defined in 3GPP TS 23.038 Section 6.1.2.1, into an array of septets.
//
// The fillBits is the number of bits of pad at the beginning of the src, as
// the packed septets may not start on an octet boundary.
func Unpack7Bit(p []byte, fillBits int) []byte {
	if len(p) == 0 {
		return append(p[:0:0], p...)
	}
	u := make([]byte, 0, (len(p)*8+6+fillBits)/7)
	var r byte
	var rbits uint
	if fillBits != 0 {
		rbits = uint(7 - fillBits)
	}
	for _, o := range p {
		r = (r | o<<rbits) & 0x7f
		u = append(u, r)
		if rbits == 6 {
			// only needed 1 bit from p, so there is a complete septet left...
			u = append(u, o>>1)
			rbits = 0
			r = 0
		} else {
			// each octet provides one extra residual bit
			rbits++
			r = o >> (8 - rbits)
		}
	}
	if fillBits > 0 {
		u = u[1:]
	}
	return u
}

// Pack7BitUSSD packs an array of septets into an 8bit array as per the packing
// rules defined in 3GPP TS 23.038 Section 6.1.2.3
//
// The padBits is the number of bits of pad to place at the beginning of the
// packed array, as the packed septets may not start on an octet boundary.
//
// A filler CR is added to the final octet if there are 7 bits unused (to
// distinguish from the 0x00 septet), or if the last septet is CR and ends on
// an octet boundary (so it wont be considered filler).
func Pack7BitUSSD(u []byte, fillBits int) []byte {
	b := Pack7Bit(u, fillBits)
	if len(b) == 0 {
		return append(b[:0:0], b...)
	}
	last := len(b) - 1
	if b[last]&^0x1 == 0 && u[len(u)-1] != 0 {
		b[last] = b[last] | (cr << 1)
	} else if len(u)&0x7 == 0 && u[len(u)-1] == cr {
		b = append(b, cr)
	}
	return b
}

// Unpack7BitUSSD unpacks septets, packed into an 8bit array, as per the
// packing rules defined in 3GPP TS 23.038 Section 6.1.2.3, into an array of
// septets.
//
// The fillBits is the number of bits of pad at the beginning of the src, as
// the packed septets may not start on an octet boundary.
//
// Any trailing CR is assumed to be filler if it ends on an octet boundary, or
// if it starts on an octet boundary and the previous character is also CR.
func Unpack7BitUSSD(p []byte, fillBits int) []byte {
	u := Unpack7Bit(p, fillBits)
	// remove any trailing filler
	if len(p) > 1 && ((p[len(p)-1]>>1 == cr) || (p[len(p)-1] == cr && p[len(p)-2]>>1 == cr)) {
		u = u[:len(u)-1]
	}
	return u
}
