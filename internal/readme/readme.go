// SPDX-License-Identifier: MIT
//
// Copyright © 2020 Kent Gibson <warthog618@gmail.com>.

package main

import (
	"github.com/warthog618/sms"
	"github.com/warthog618/sms/encoding/gsm7/charset"
	"github.com/warthog618/sms/encoding/tpdu"
)

func main() {

}

func encode() {
	tpdus, _ := sms.Encode([]byte("hello world"))
	for _, p := range tpdus {
		b, _ := p.MarshalBinary()
		sendPDU(b) // send binary TPDU...
	}
}

func encoder() {
	msgChan := make(chan []byte)
	e := sms.NewEncoder()
	for {
		msg := <-msgChan
		tpdus, _ := e.Encode(msg)
		for _, p := range tpdus {
			b, _ := p.MarshalBinary()
			sendPDU(b) // send binary TPDU...
		}
	}
}

func unmarshal() tpdu.TPDU {
	bintpdu := []byte{}

	pdu, _ := sms.Unmarshal(bintpdu)

	return *pdu
}

func decodeOne() []byte {
	pdu := &tpdu.TPDU{}

	msg, _ := sms.Decode([]*tpdu.TPDU{pdu})

	return msg
}

func decodeMany() []byte {
	tpdus := []*tpdu.TPDU{}

	msg, _ := sms.Decode(tpdus)

	return msg
}

func collect() {
	pduChan := make(chan []byte)
	c := sms.NewCollector()
	for {
		bintpdu := <-pduChan
		pdu, _ := sms.Unmarshal(bintpdu)
		tpdus, _ := c.Collect(*pdu)
		if len(tpdus) > 0 {
			msg, _ := sms.Decode(tpdus)
			// handle msg...
			handleMsg(msg)
		}
	}

}

func to() []tpdu.TPDU {
	tpdus, _ := sms.Encode([]byte("hello"), sms.To("12345"))

	return tpdus
}

func urdu() []tpdu.TPDU {

	tpdus, _ := sms.Encode([]byte("hello ٻ"), sms.WithCharset(charset.Urdu))

	return tpdus
}

func deliver() []tpdu.TPDU {

	tpdus, _ := sms.Encode([]byte("hello"), sms.AsDeliver, sms.From("12345"))

	return tpdus
}

func unmarshalOpt() tpdu.TPDU {
	bintpdu := []byte("")

	tpdu, _ := sms.Unmarshal(bintpdu, sms.AsMO)

	return *tpdu
}

func handleMsg([]byte) {

}

func sendPDU([]byte) {
}
