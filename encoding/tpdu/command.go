// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// Command represents an SMS Command TPDU as defined in 3GPP TS 23.040 Section 9.2.2.4.
type Command struct {
	BaseTPDU
	mr byte
	ct byte
	mn byte
	da Address
	// BaseTPDU ud, including the udh, is taken to be the cd
}

// NewCommand creates a Command TPDU and initialises non-zero fields.
func NewCommand() *Command {
	// Command doesn't use dcs, but baseTPDU does to determine UD alphabet, so set it to 8bit.
	return &Command{BaseTPDU: BaseTPDU{firstOctet: byte(MtCommand), dcs: 0x04}}
}

// MR returns the Command mr.
func (c *Command) MR() byte {
	return c.mr
}

// CT returns the Command ct.
func (c *Command) CT() byte {
	return c.ct
}

// MN returns the Command mn.
func (c *Command) MN() byte {
	return c.mn
}

// DA returns the Command da.
func (c *Command) DA() Address {
	return c.da
}

// SetMR sets the Command mr field.
func (c *Command) SetMR(mr byte) {
	c.mr = mr
}

// SetCT sets the Command ct field.
func (c *Command) SetCT(ct byte) {
	c.ct = ct
}

// SetMN sets the Command mn field.
func (c *Command) SetMN(mn byte) {
	c.mn = mn
}

// SetDA sets the Command da field.
func (c *Command) SetDA(da Address) {
	c.da = da
}

// MarshalBinary marshals an SMS-Command-Report TPDU.
func (c *Command) MarshalBinary() ([]byte, error) {
	b := []byte{c.firstOctet, c.mr, c.pid, c.ct, c.mn}
	da, err := c.da.MarshalBinary()
	if err != nil {
		return nil, EncodeError("da", err)
	}
	b = append(b, da...)
	cdl := len(c.ud)
	b = append(b, byte(cdl))
	b = append(b, c.ud...)
	return b, nil
}

// UnmarshalBinary unmarshals an SMS-Command-Report TPDU.
func (c *Command) UnmarshalBinary(src []byte) error {
	if len(src) < 1 {
		return DecodeError("firstOctet", 0, ErrUnderflow)
	}
	c.firstOctet = src[0]
	ri := 1
	if len(src) <= ri {
		return DecodeError("mr", ri, ErrUnderflow)
	}
	c.mr = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("pid", ri, ErrUnderflow)
	}
	c.pid = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("ct", ri, ErrUnderflow)
	}
	c.ct = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("mn", ri, ErrUnderflow)
	}
	c.mn = src[ri]
	ri++
	n, err := c.da.UnmarshalBinary(src[ri:])
	if err != nil {
		return DecodeError("da", ri, err)
	}
	ri += n
	c.dcs = 0x04 // force BaseTPDU to interpret UD as 8bit, if nt set already
	err = c.decodeUserData(src[ri:])
	if err != nil {
		return DecodeError("ud", ri, err)
	}
	return nil
}

func decodeCommand(src []byte) (TPDU, error) {
	c := NewCommand()
	if err := c.UnmarshalBinary(src); err != nil {
		return nil, err
	}
	return c, nil
}

// RegisterCommandDecoder registers a decoder for the Command TPDU.
func RegisterCommandDecoder(d *Decoder) error {
	return d.RegisterDecoder(MtCommand, MO, decodeCommand)
}
