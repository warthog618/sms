// Copyright © 2018 Kent Gibson <warthog618@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package ms provides functionality specific to the Mobile Station.
// The contained packages add layers of functionality above tpdu:
// - sar provides segmentation and reassembly above tpdu
// - message provides conversion to abstract messages above sar
// - pdumode provides stuff...
package ms
