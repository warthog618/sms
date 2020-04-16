// SPDX-License-Identifier: MIT
//
// Copyright © 2019 Kent Gibson <warthog618@gmail.com>.
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/warthog618/sms"
	"github.com/warthog618/sms/encoding/pdumode"
)

func main() {
	var pm bool
	flag.BoolVar(&pm, "p", false, "PDU is prefixed with SCA (PDU mode)")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		flag.Usage()
		os.Exit(1)
	}
	c := sms.NewCollector()
	defer c.Close()
	for _, a := range flag.Args() {
		b, err := hex.DecodeString(a)
		if err != nil {
			log.Fatal(err)
		}
		tb := b
		if pm {
			pdu, err := pdumode.UnmarshalBinary(b)
			if err != nil {
				log.Fatal(err)
			}
			tb = pdu.TPDU
		}
		t, err := sms.Unmarshal(tb)
		if err != nil {
			log.Printf("unmarshal error: %v", err)
			continue
		}
		pdus, err := c.Collect(*t)
		if err != nil {
			log.Printf("collect error: %v", err)
		}
		if pdus == nil {
			continue
		}
		msg, err := sms.Decode(pdus)
		if err != nil {
			log.Printf("decode error: %v", err)
		}
		if msg != nil {
			fmt.Printf("%s: %s\n", pdus[0].OA.Number(), msg)
		}
	}
	// report active collect pipes
	pipes := c.Pipes()
	for k, v := range pipes {
		fmt.Println("incomplete reassembly: ", k)
		fmt.Println(v)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "smsdeliver decodes and displays the message from one or more SMS Deliver TPDUs.\n"+
		"Usage: smsdeliver [-p] <pdu> [pdu...]\n")
	flag.PrintDefaults()
}
