// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package tpdu

// Deliver represents a SMS-Deliver PDU as defined in 3GPP TS 23.038 Section 9.2.2.1.
type Deliver struct {
	TPDU
	OA Address
	// The SCTS timestamp indicates the time the SMS was sent.
	// The time is the originator's local time, the timezone of which may
	// differ from the receiver's.
	SCTS Timestamp
}

// NewDeliver creates a Deliver TPDU and initialises non-zero fields.
func NewDeliver() *Deliver {
	return &Deliver{
		TPDU: TPDU{FirstOctet: byte(MtDeliver)},
		OA:   Address{TOA: 0x80},
	}
}

// MaxUDL returns the maximum number of octets that can be encoded into the UD.
// Note that for 7bit encoding this can result in up to 160 septets.
func (d *Deliver) MaxUDL() int {
	return 140
}

// MarshalBinary marshals a SMS-Deliver PDU into the corresponding byte array.
func (d *Deliver) MarshalBinary() ([]byte, error) {
	b := []byte{d.FirstOctet}
	oa, err := d.OA.MarshalBinary()
	if err != nil {
		return nil, EncodeError("oa", err)
	}
	b = append(b, oa...)
	b = append(b, d.PID, byte(d.DCS))
	scts, err := d.SCTS.MarshalBinary()
	if err != nil {
		return nil, EncodeError("scts", err)
	}
	b = append(b, scts...)
	ud, err := d.encodeUserData()
	if err != nil {
		return nil, EncodeError("ud", err)
	}
	b = append(b, ud...)
	return b, nil
}

// UnmarshalBinary unmarshals a SMS-Deliver PDU from the corresponding byte array.
// In the case of error the Deliver will be partially unmarshalled, up to the
// point that the decoding error was detected.
func (d *Deliver) UnmarshalBinary(src []byte) error {
	if len(src) < 1 {
		return DecodeError("firstOctet", 0, ErrUnderflow)
	}
	d.FirstOctet = src[0]
	ri := 1
	n, err := d.OA.UnmarshalBinary(src[ri:])
	if err != nil {
		return DecodeError("oa", ri, err)
	}
	ri += n
	if len(src) <= ri {
		return DecodeError("pid", ri, ErrUnderflow)
	}
	d.PID = src[ri]
	ri++
	if len(src) <= ri {
		return DecodeError("dcs", ri, ErrUnderflow)
	}
	d.DCS = DCS(src[ri])
	ri++
	if len(src) < ri+7 {
		return DecodeError("scts", ri, ErrUnderflow)
	}
	err = d.SCTS.UnmarshalBinary(src[ri : ri+7])
	if err != nil {
		return DecodeError("scts", ri, err)
	}
	ri += 7
	err = d.decodeUserData(src[ri:])
	if err != nil {
		return DecodeError("ud", ri, err)
	}
	return nil
}

func decodeDeliver(src []byte) (interface{}, error) {
	d := NewDeliver()
	if err := d.UnmarshalBinary(src); err != nil {
		return nil, err
	}
	return d, nil
}

// RegisterDeliverDecoder registers a decoder for the Deliver TPDU.
func RegisterDeliverDecoder(d *Decoder) error {
	return d.RegisterDecoder(MtDeliver, MT, decodeDeliver)
}

// RegisterReservedDecoder registers a decoder for the Deliver TPDU for the
// Reserved message type.
func RegisterReservedDecoder(d *Decoder) error {
	return d.RegisterDecoder(MtReserved, MT, decodeDeliver)
}
