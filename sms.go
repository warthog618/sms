// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package sms provides encoders and decoders for SMS PDUs.
package sms

import (
	"log"

	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/ms/message"
)

// Encode generates the set of Submit TPDUs to send a message to a number.
//
// number is the destination phone number, which should be international.
// msg is the message to be sent.
// Returns the set of Submit TPDUs, or any error detected during encoding.
func Encode(number, msg string, options ...message.EncoderOption) ([]tpdu.Submit, error) {
	e := message.NewEncoder(options...)
	return e.Encode(number, msg)
}

// Encode8bit generates the set of Submit TPDUs to send a block of 8bit data to a number.
//
// number is the destination phone number, which should be international.
// d is the data to be sent.
// Returns the set of Submit TPDUs, or any error detected during encoding.
func Encode8bit(number string, d []byte, options ...message.EncoderOption) ([]tpdu.Submit, error) {
	e := message.NewEncoder(options...)
	return e.Encode8Bit(number, d)
}

func Decode(src []byte, drn tpdu.Direction) (interface{}, error) {
	td, err := tpdu.NewDecoder(
		// !!! fully loaded decoder should be the default...
		tpdu.RegisterCommandDecoder,
		tpdu.RegisterDeliverDecoder,
		tpdu.RegisterDeliverReportDecoder,
		tpdu.RegisterReservedDecoder,
		tpdu.RegisterSubmitDecoder,
		tpdu.RegisterSubmitReportDecoder,
		tpdu.RegisterStatusReportDecoder,
	)
	if err != nil {
		log.Fatal(err)
	}
	return td.Decode(src, drn)
}
