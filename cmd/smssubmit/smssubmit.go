// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/warthog618/sms/encoding/gsm7/charset"
	"github.com/warthog618/sms/encoding/tpdu"
	"github.com/warthog618/sms/ms/message"
	"github.com/warthog618/sms/ms/sar"
)

func main() {
	var number, msg string
	var language uint
	flag.StringVar(&number, "number", "", "Destination number in international format")
	flag.StringVar(&msg, "message", "", "The message to encode")
	flag.UintVar(&language, "language", 0, "The NLI of a character set to use in addition to the default")
	flag.Usage = usage
	flag.Parse()
	if number == "" || msg == "" {
		flag.Usage()
		os.Exit(1)
	}

	e, err := tpdu.NewUDEncoder()
	if err != nil {
		log.Println(err)
		return
	}
	if language != 0 {
		e.AddLockingCharset(charset.NationalLanguageIdentifier(language))
		e.AddShiftCharset(charset.NationalLanguageIdentifier(language))
	}
	sg := sar.NewSegmenter()
	t := message.NewEncoder(e, sg)
	pdus, err := t.Encode(number, msg)
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
