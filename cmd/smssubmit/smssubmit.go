// SPDX-License-Identifier: MIT
//
// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.

// smssubmit provides an example of encoding a message into a set of SMS-SUBMIT
// TPDUs.
package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/warthog618/sms"
)

func main() {
	var number, msg string
	var nli int
	flag.StringVar(&number, "number", "", "Destination number in international format")
	flag.StringVar(&msg, "message", "", "The message to encode")
	flag.IntVar(&nli, "language", 0, "The NLI of a character set to use in addition to the default")
	flag.Usage = usage
	flag.Parse()
	if number == "" || msg == "" {
		flag.Usage()
		os.Exit(1)
	}

	options := []sms.EncoderOption{sms.AsSubmit, sms.To(number)}
	if nli != 0 {
		options = append(options, sms.WithCharset(nli))
	}
	pdus, err := sms.Encode([]byte(msg), options...)
	if err != nil {
		log.Println(err)
		return
	}
	if len(pdus) == 1 {
		b, _ := pdus[0].MarshalBinary()
		fmt.Printf("Submit TPDU:\n%s\n", hex.EncodeToString(b))
		return
	}
	for i, p := range pdus {
		b, _ := p.MarshalBinary()
		fmt.Printf("Submit TPDU %d:\n%s\n", i+1, hex.EncodeToString(b))
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "smssubmit encodes a message into a SMS Submit TPDU.\n"+
		"The message is encoded using the GSM7 default alphabet, or if necessary\n"+
		"an optionally specified character set, or failing those as UCS-2.\n"+
		"If the message is too long for a single PDU then it is split into several.\n\n"+
		"Usage: smssubmit -number <number> -message <message>\n")
	flag.PrintDefaults()
}
