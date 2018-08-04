# sms

A Go library for encoding and decoding SMSs.

[![Build Status](https://travis-ci.org/warthog618/sms.svg)](https://travis-ci.org/warthog618/sms)
[![Coverage Status](https://coveralls.io/repos/github/warthog618/sms/badge.svg?branch=master)](https://coveralls.io/github/warthog618/sms?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/warthog618/sms)](https://goreportcard.com/report/github.com/warthog618/sms)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/warthog618/sms/blob/master/LICENSE)

sms is a Go library for encoding and decoding SMS TPDUs, as specified in 3GPP TS 23.040 and 3GPP TS 23.038.
The initial impetus was to provide functionality to send and receive SMSs via
a GSM modem, but the library may be generally useful anywhere encoding and decoding SMS TPDUs or their fields is required.

## Features

Supports the following functionality:

- Creation of SMS Submit TPDUs from UTF-8 strings
- Segmentation of large messages into several concatenated SMS Submit TPDUs
- Automatic selection of alphabet and language when encoding
- Decoding of SMS Deliver TPDUs into UTF-8 strings
- Reassembly of concatenated SMS Deliver TPDUs into a single large message
- Supports encoding and decoding all SMS TPDU types, not just Submit and Deliver
- Supports encoding and decoding SMS TPDUs to be sent and recevied via GSM modems in PDU mode

## Contained Packages

The [tpdu](encoding/tpdu) package [![GoDoc](https://godoc.org/github.com/warthog618/sms/encoding/tpdu?status.svg)](https://godoc.org/github.com/warthog618/sms/encoding/tpdu) provides the core TPDU types and conversions to and from their binary form.

Several packages build on top of tpdu to provide higher level functionality:

The [sar](ms/sar) package [![GoDoc](https://godoc.org/github.com/warthog618/sms/ms/sar?status.svg)](https://godoc.org/github.com/warthog618/sms/ms/sar) provides segmentation and reassembly of concatenated SMS TPDUs to implement large messages.

The [message](ms/message) package [![GoDoc](https://godoc.org/github.com/warthog618/sms/ms/message?status.svg)](https://godoc.org/github.com/warthog618/sms/ms/message) provides a layer above sar that allows simplfied encoding and decoding of messages with only the message and the destination number.

The [pdumode](ms/pdumode) package [![GoDoc](https://godoc.org/github.com/warthog618/sms/ms/pdumode?status.svg)](https://godoc.org/github.com/warthog618/sms/ms/pdumode) provides encoding and decoding of PDUs exchanged with GSM modems in PDU mode.

A number of packages provide functionality to encode and decode TPDU fields:

The [bcd](encoding/bcd) package [![GoDoc](https://godoc.org/github.com/warthog618/sms/encoding/bcd?status.svg)](https://godoc.org/github.com/warthog618/sms/encoding/bcd) provides conversions to and from BCD format.

The [gsm7](encoding/gsm7) package [![GoDoc](https://godoc.org/github.com/warthog618/sms/encoding/gsm7?status.svg)](https://godoc.org/github.com/warthog618/sms/encoding/gsm7) provides conversions to and from 7bit packed user data.

The [charset](encoding/gsm7/charset) package [![GoDoc](https://godoc.org/github.com/warthog618/sms/encoding/gsm7/charset?status.svg)](https://godoc.org/github.com/warthog618/sms/encoding/gsm7/charset) provides the character sets used to encode user data in GSM 7bit format as specified in 3GPP TS 23.038.

The [semioctet](encoding/semioctet) package [![GoDoc](https://godoc.org/github.com/warthog618/sms/encoding/semioctet?status.svg)](https://godoc.org/github.com/warthog618/sms/encoding/semioctet) provides conversions to and from semioctet format.

The [ucs2](encoding/ucs2) package [![GoDoc](https://godoc.org/github.com/warthog618/sms/encoding/ucs2?status.svg)](https://godoc.org/github.com/warthog618/sms/encoding/ucs2) provides conversions between UCS-2 and UTF-8.

## Examples

The [cmd](cmd) directory contains basic commands to exercise, debug and demonstrate the core functionality of the library, including:

- decoding arbitrary TPDUs [(smsdecode)](cmd/smsdecode/smsdecode.go)
- encoding messages into SMS Submit TPDUs [(smssubmit)](cmd/smssubmit/smssubmit.go)
- decoding SMS Deliver TPDUs into messages [(smsdeliver)](cmd/smsdeliver/smsdeliver.go)
- displaying supported character sets [(charsets)](cmd/charsets/charsets.go).

The following examples demonstrate the example commands, and their code provides some examples of using the library.

### Submit Encoding

Creating an SMS to send:

```shell
$ smssubmit -number 12345 -message "Hello world"
Submit TPDU:
010005912143f500000bc8329bfd06dddf723619
```

### Deliver Decoding

Decoding an SMS received from a GSM modem in PDU mode:

```shell
$ smsdeliver -p 07911614220991F1040B911605935713F200008140806113912304D7F79B0E
+61503975312: Woot
```

### Concatenated Message Decoding

Concatenating and displaying a message split into multiple Deliver TPDUs.

```shell
$ smsdeliver 400B911605935713F20008814080611373238C050003C00301007400680069007300200069007300200061002000760065007200790020006C006F006E00670020006D0065007300730061006700650020007400680061007400200064006F006500730020006E006F0074002000660069007400200069006E00200061002000730069006E0067006C006500200053004D00530020006D0065007300730061 400B911605935713F20008814080611373238C050003C0030200670065002C0020006100740020006C0065006100730074002000690074002000770069006C006C002000690066002000490020006B00650065007000200061006400640069006E00670020006D006F0072006500200074006F0020006900740020006100730020003100360030002000630068006100720061006300740065007200730020 440B911605935713F200088140806113832344050003C00303006900730020006D006F007200650020007400680061006E00200079006F00750020006D00690067006800740020007400680069006E006B0020D83DDE01
+61503975312: this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think 😁

```

### General TPDU Decoding

Decoding the Submit TPDU created above:

```shell
$ smsdecode -o 010005912143f500000bc8329bfd06dddf723619
(*tpdu.Submit)(0xc4200c0000)({
 TPDU: (tpdu.TPDU) {
  FirstOctet: (uint8) 1,
  PID: (uint8) 0,
  DCS: (uint8) 0,
  UDH: (tpdu.UserDataHeader) <nil>,
  UD: (tpdu.UserData) (len=11 cap=12) {
   00000000  48 65 6c 6c 6f 20 77 6f  72 6c 64                 |Hello world|
  }
 },
 MR: (uint8) 0,
 DA: (tpdu.Address) {
  TOA: (uint8) 145,
  Addr: (string) (len=5) "12345"
 },
 VP: (tpdu.ValidityPeriod) {
  Format: (tpdu.ValidityPeriodFormat) 0,
  Time: (tpdu.Timestamp) 0001-01-01 00:00:00 +0000 UTC,
  Duration: (time.Duration) 0s,
  EFI: (uint8) 0
 }
})
```