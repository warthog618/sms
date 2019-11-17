// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// Command represents an SMS Command TPDU as defined in 3GPP TS 23.040 Section
// 9.2.2.4.
type Command struct {
	TPDU
	MR byte
	CT byte
	MN byte
	DA Address
	// TPDU UD, including the UDH, is taken to be the CD
}

// NewCommand creates a Command TPDU and initialises non-zero fields.
func NewCommand() *Command {
	// Command doesn't use dcs, but TPDU does to determine UD alphabet, so set
	// it to 8bit.
	return &Command{TPDU: TPDU{FirstOctet: byte(MtCommand), DCS: 0x04}}
}

// MarshalBinary marshals an SMS-Command-Report TPDU.
func (c *Command) MarshalBinary() ([]byte, error) {
	b := []byte{c.FirstOctet, c.MR, c.PID, c.CT, c.MN}
	da, err := c.DA.MarshalBinary()
	if err != nil {
		return nil, EncodeError("da", err)
	}
	b = append(b, da...)
	cdl := len(c.TPDU.UD)
	b = append(b, byte(cdl))
	b = append(b, c.UD...)
	return b, nil
}

// UnmarshalBinary unmarshals an SMS-Command-Report TPDU.
func (c *Command) UnmarshalBinary(src []byte) error {
	if len(src) < 1 {
		return DecodeError("firstOctet", 0, ErrUnderflow)
	}
	c.FirstOctet = src[0]
	ri := 1
	if len(src) <= ri {
		return DecodeError("mr", ri, ErrUnderflow)
	}
	c.MR = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("pid", ri, ErrUnderflow)
	}
	c.PID = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("ct", ri, ErrUnderflow)
	}
	c.CT = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("mn", ri, ErrUnderflow)
	}
	c.MN = src[ri]
	ri++
	n, err := c.DA.UnmarshalBinary(src[ri:])
	if err != nil {
		return DecodeError("da", ri, err)
	}
	ri += n
	c.DCS = 0x04 // force TPDU to interpret UD as 8bit, if nt set already
	err = c.decodeUserData(src[ri:])
	if err != nil {
		return DecodeError("ud", ri, err)
	}
	return nil
}

func decodeCommand(src []byte) (interface{}, error) {
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
