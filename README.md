# sms

A Go library for encoding and decoding SMSs.

[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/warthog618/sms)
[![Build Status](https://travis-ci.org/warthog618/sms.svg)](https://travis-ci.org/warthog618/sms)
[![Coverage Status](https://coveralls.io/repos/github/warthog618/sms/badge.svg?branch=master)](https://coveralls.io/github/warthog618/sms?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/warthog618/sms)](https://goreportcard.com/report/github.com/warthog618/sms)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://github.com/warthog618/sms/blob/master/LICENSE)

sms is a Go library for encoding and decoding SMS TPDUs, as specified in 3GPP TS 23.040 and 3GPP TS 23.038.

The initial impetus was to provide functionality to send and receive SMSs via a
GSM modem, but the library may be generally useful anywhere encoding and
decoding SMS TPDUs or their fields is required.

## Features

Supports the following functionality:

- Creation of SMS TPDUs from UTF-8 strings, including emoji's 😁
- Segmentation of long messages into several concatenated SMS TPDUs
- Automatic selection of alphabet and language when encoding
- Decoding of SMS TPDUs into UTF-8 strings
- Reassembly of concatenated SMS TPDUs into a long message
- Support for all GSM character sets
- Encoding and decoding SMS TPDUs in PDU mode for exchange with GSM modems

## Usage

```go
import "github.com/warthog618/sms"
```

In the following usage examples the error handling is omitted for brevity.

### Encode

Creating the TPDUs to contain a message is referred to as encoding.

A one-off message can be encoded using *sms.Encode*:

```go
tpdus, _ := sms.Encode([]byte("hello world"))
for _, p := range tpdus {
    b, _ := p.MarshalBinary()
    // send binary TPDU...
}
```

Sending multiple messages requires maintaining multiple counter fields and
encoding them in the TPDU.  This is performed by an *sms.Encoder*:

```go
e := sms.NewEncoder()
for {
    msg := <- msgChan
    tpdus, _ := e.Encode(msg)
    for _, p := range tpdus {
        b, _ := p.MarshalBinary()
        // send binary TPDU...
    }
}
```

### Unmarshal

Reassembling received TPDUs into a complete message is a multi-step process.
The first step is to unmarhsal the binary SMS TPDU into a TPDU object using
*sms.Unmarshal*:

```go
pdu, _ := sms.Unmarshal(bintpdu)
```

### Decode

A single segment TPDU can be decoded using *sms.Decode*:

```go
msg, _ := sms.Decode([]*tpdu.TPDUs{pdu})
```

For concatenated messages, the set of TPDUs containing a message is reassembled
into a complete message using *sms.Decode*:

```go
msg, _ := sms.Decode(tpdus)
```

### Collect

The segments of concatenated messages must be collected before they can be
decoded.  The Collector collects received segments and returns the complete set
once the final segment is received.

```go
c := sms.NewCollector()
for {
    bintpdu := <- pduChan
    pdu, _ := sms.Unmarshal(bintpdu)
    tpdus, _ := c.Collect(pdu)
    if len(tpdus) > 0 {
        msg, _ := sms.Decode(tpdus)
        // handle msg...
    }
}
```

### Options

The core API is aimed at the most common use cases, those performed to the
mobile station. e.g. By default, *sms.Encode* creates an SMS-SUBMIT TPDU and
only uses the default character set.  By default, *sms.Decode* uses all
character sets.  By default, *sms.Unmarshal* assumes the TPDU is mobile
terminating.

The behaviour of the core API functions can be altered for other use cases
using optional parameters.

e.g. to specify the destination number for a SMS-SUBMIT message:

```go
tpdus, _ := sms.Encode("hello",sms.To("12345"))
```

or to encode a message using a particular character set, if necessary:

```go
tpdus, _ := sms.Encode("hello ٻ",sms.WithCharset(charset.Urdu))
```

or to specify the encoding of a SMS-DELIVER message:

```go
tpdus, _ := sms.Encode("hello",sms.AsDeliver,sms.From("12345"))
```

or to unmarshal a TPDU from the mobile station:

```go
pdu, _ := sms.Unmarshal(bintpdu,sms.AsMO)
```

The full set of supplied options:

Option | Category | Description
---|---|---
*WithReassemblyTimeout(duration,handler)*|Collect|Limit the time allowed to wait for the TPDUs of a complete reassembly
*WithTemplateOption(tpdu)*|Encode|Use the provided TPDU as the template for encoded TPDUs.
*To(number)*|Encode|Set the DA of the encoded TPDU to the number provided
*From(number)*|Encode|Set the OA of the encoded TPDU to the number provided
*WithAllCharsets*|Decode,Encode|Make all GSM7 character sets available
*WithDefaultCharset*|Decode,Encode|Make only the default character set available
*WithCharset(nli...)*|Decode,Encode|Make the specified character set(s) available
*WithLockingCharset(nli...)*|Decode,Encode|Make the specified character set(s) available as a locking character set
*WithShiftCharset(nli...)*|Decode,Encode|Make the specified character set(s) available as a shift character set
*AsSubmit*|Encode|Encode the TPDU as a SMS-SUBMIT (default)
*AsDeliver*|Encode|Encode the TPDU as a SMS-DELIVER
*As8Bit*|Encode|Force the encoding of user data as 8-bit
*AsUCS2*|Encode|Force the encoding of user data as UCS-2
*AsMO*|Unmarshal|Treat the TPDU as originating from the mobile station
*AsMT*|Unmarshal|Treat the TPDU as terminating at the mobile station (default)

## Tools

The [cmd](cmd) directory contains basic commands tools to exercise, debug and
demonstrate the core functionality of the library, including:

- encoding messages into SMS-SUBMIT TPDUs [(smssubmit)](cmd/smssubmit/smssubmit.go)
- decoding SMS-DELIVER TPDUs into messages [(smsdeliver)](cmd/smsdeliver/smsdeliver.go)
- decoding arbitrary TPDUs [(smsdecode)](cmd/smsdecode/smsdecode.go)
- displaying supported character sets [(charsets)](cmd/charsets/charsets.go).

The following examples demonstrate the example commands, and their code provides some examples of using the library.

### Submit Encoding

Creating an SMS to send:

```shell
$ smssubmit -number 12345 -message "Hello world"
Submit TPDU:
010105912143f500000bc8329bfd06dddf723619
```

Long messages are split into a concatenated message spanning several TPDUs:

```shell
smssubmit -number 12345 -message "this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think 😁"
Submit TPDU 1:
410105912143f500088c050003010301007400680069007300200069007300200061002000760065007200790020006c006f006e00670020006d0065007300730061006700650020007400680061007400200064006f006500730020006e006f0074002000660069007400200069006e00200061002000730069006e0067006c006500200053004d00530020006d0065007300730061
Submit TPDU 2:
410205912143f500088c05000301030200670065002c0020006100740020006c0065006100730074002000690074002000770069006c006c002000690066002000490020006b00650065007000200061006400640069006e00670020006d006f0072006500200074006f0020006900740020006100730020003100360030002000630068006100720061006300740065007200730020
Submit TPDU 3:
410305912143f5000844050003010303006900730020006d006f007200650020007400680061006e00200079006f00750020006d00690067006800740020007400680069006e006b0020d83dde01
```

### Deliver Decoding

Decoding an SMS received from a GSM modem in PDU mode:

```shell
$ smsdeliver -p 07911614220991F1040B911605935713F200008140806113912304D7F79B0E
+61503975312: Woot
```

### Concatenated Message Decoding

Concatenating and displaying a message split into multiple SMS-DELIVER TPDUs.

```shell
$ smsdeliver 400B911605935713F20008814080611373238C050003C00301007400680069007300200069007300200061002000760065007200790020006C006F006E00670020006D0065007300730061006700650020007400680061007400200064006F006500730020006E006F0074002000660069007400200069006E00200061002000730069006E0067006C006500200053004D00530020006D0065007300730061 400B911605935713F20008814080611373238C050003C0030200670065002C0020006100740020006C0065006100730074002000690074002000770069006C006C002000690066002000490020006B00650065007000200061006400640069006E00670020006D006F0072006500200074006F0020006900740020006100730020003100360030002000630068006100720061006300740065007200730020 440B911605935713F200088140806113832344050003C00303006900730020006D006F007200650020007400680061006E00200079006F00750020006D00690067006800740020007400680069006E006B0020D83DDE01
+61503975312: this is a very long message that does not fit in a single SMS message, at least it will if I keep adding more to it as 160 characters is more than you might think 😁
```

### General TPDU Decoding

Decoding the Submit TPDU created above:

```shell
$ smsdecode -o 010105912143f500000bc8329bfd06dddf723619
TPDU: SMS-SUBMIT
TP-MTI: 0x01 Submit
TP-RD: false
TP-VPF: 0x00 Not Present
TP-RP: false
TP-UDHI: false
TP-SRR: false
TP-MR: 1
TP-DA: +12345
TP-PID: 0x00
TP-DCS: 0x00 7bit
TP-VP: Not Present
TP-UD: 00000000  48 65 6c 6c 6f 20 77 6f  72 6c 64                 |Hello world|
```

Decoding the Deliver TPDU above:

```shell
$ smsdecode -p 07911614220991F1040B911605935713F200008140806113912304D7F79B0E
SMSC: +61412290191
TPDU: SMS-DELIVER
TP-MTI: 0x00 Deliver
TP-MMS: true
TP-LP: false
TP-RP: false
TP-UDHI: false
TP-SRI: false
TP-OA: +61503975312
TP-PID: 0x00
TP-DCS: 0x00 7bit
TP-SCTS: 2018-04-08 16:31:19 +0800
TP-UD: 00000000  57 6f 6f 74                                       |Woot|
```

Decoding the first Deliver TPDU of the concatenated message above:

```shell
$ smsdecode 400B911605935713F20008814080611373238C050003C00301007400680069007300200069007300200061002000760065007200790020006C006F006E00670020006D0065007300730061006700650020007400680061007400200064006F006500730020006E006F0074002000660069007400200069006E00200061002000730069006E0067006C006500200053004D00530020006D0065007300730061
TPDU: SMS-DELIVER
TP-MTI: 0x00 Deliver
TP-MMS: false
TP-LP: false
TP-RP: false
TP-UDHI: true
TP-SRI: false
TP-OA: +61503975312
TP-PID: 0x00
TP-DCS: 0x08 UCS-2
TP-SCTS: 2018-04-08 16:31:37 +0800
TP-UDH: ID: 0  Data: [192 3 1]
TP-UD: 00000000  00 74 00 68 00 69 00 73  00 20 00 69 00 73 00 20  |.t.h.i.s. .i.s. |
       00000010  00 61 00 20 00 76 00 65  00 72 00 79 00 20 00 6c  |.a. .v.e.r.y. .l|
       00000020  00 6f 00 6e 00 67 00 20  00 6d 00 65 00 73 00 73  |.o.n.g. .m.e.s.s|
       00000030  00 61 00 67 00 65 00 20  00 74 00 68 00 61 00 74  |.a.g.e. .t.h.a.t|
       00000040  00 20 00 64 00 6f 00 65  00 73 00 20 00 6e 00 6f  |. .d.o.e.s. .n.o|
       00000050  00 74 00 20 00 66 00 69  00 74 00 20 00 69 00 6e  |.t. .f.i.t. .i.n|
       00000060  00 20 00 61 00 20 00 73  00 69 00 6e 00 67 00 6c  |. .a. .s.i.n.g.l|
       00000070  00 65 00 20 00 53 00 4d  00 53 00 20 00 6d 00 65  |.e. .S.M.S. .m.e|
       00000080  00 73 00 73 00 61                                 |.s.s.a|
```

Decoding the second Submit TPDU in the concatenated message above:

```shell
smsdecode -o 410205912143f500088c05000301030200670065002c0020006100740020006c0065006100730074002000690074002000770069006c006c002000690066002000490020006b00650065007000200061006400640069006e00670020006d006f0072006500200074006f0020006900740020006100730020003100360030002000630068006100720061006300740065007200730020
TPDU: SMS-SUBMIT
TP-MTI: 0x01 Submit
TP-RD: false
TP-VPF: 0x00 Not Present
TP-RP: false
TP-UDHI: true
TP-SRR: false
TP-MR: 2
TP-DA: +12345
TP-PID: 0x00
TP-DCS: 0x08 UCS-2
TP-VP: Not Present
TP-UDH: ID: 0  Data: [1 3 2]
TP-UD: 00000000  00 67 00 65 00 2c 00 20  00 61 00 74 00 20 00 6c  |.g.e.,. .a.t. .l|
       00000010  00 65 00 61 00 73 00 74  00 20 00 69 00 74 00 20  |.e.a.s.t. .i.t. |
       00000020  00 77 00 69 00 6c 00 6c  00 20 00 69 00 66 00 20  |.w.i.l.l. .i.f. |
       00000030  00 49 00 20 00 6b 00 65  00 65 00 70 00 20 00 61  |.I. .k.e.e.p. .a|
       00000040  00 64 00 64 00 69 00 6e  00 67 00 20 00 6d 00 6f  |.d.d.i.n.g. .m.o|
       00000050  00 72 00 65 00 20 00 74  00 6f 00 20 00 69 00 74  |.r.e. .t.o. .i.t|
       00000060  00 20 00 61 00 73 00 20  00 31 00 36 00 30 00 20  |. .a.s. .1.6.0. |
       00000070  00 63 00 68 00 61 00 72  00 61 00 63 00 74 00 65  |.c.h.a.r.a.c.t.e|
       00000080  00 72 00 73 00 20                                 |.r.s. |
```

## Subpackages

The [tpdu](encoding/tpdu) package [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/warthog618/sms/encoding/tpdu) provides the core TPDU types and conversions to and from their binary form.

The [pdumode](ms/pdumode) package [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/warthog618/sms/ms/pdumode) provides encoding and decoding of PDUs exchanged with GSM modems in PDU mode.

A number of packages provide functionality to encode and decode TPDU fields:

The [bcd](encoding/bcd) package [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/warthog618/sms/encoding/bcd) provides conversions to and from BCD format.

The [gsm7](encoding/gsm7) package [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/warthog618/sms/encoding/gsm7) provides conversions to and from 7bit packed user data.

The [charset](encoding/gsm7/charset) package [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/warthog618/sms/encoding/gsm7/charset) provides the character sets used to encode user data in GSM 7bit format as specified in 3GPP TS 23.038.

The [semioctet](encoding/semioctet) package [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/warthog618/sms/encoding/semioctet) provides conversions to and from semioctet format.

The [ucs2](encoding/ucs2) package [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=for-the-badge)](https://pkg.go.dev/github.com/warthog618/sms/encoding/ucs2) provides conversions between UCS-2 and UTF-8.
