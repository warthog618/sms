// SPDX-License-Identifier: MIT
//
// Copyright © 2019 Kent Gibson <warthog618@gmail.com>.package main

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCount(t *testing.T) {
	patterns := []struct {
		name string
		msg  string
		nli  int
		out  Count
		err  error
	}{
		{
			"std",
			"content of the SMS",
			0,
			Count{"7BIT", 1, 18, 18, 160, 142},
			nil,
		},
		{
			"grin",
			"hello 😁",
			0,
			Count{"UCS-2", 1, 8, 8, 70, 62},
			nil,
		},
		{
			"urdu locking",
			"hi ت",
			13,
			Count{"7BIT", 1, 4, 4, 155, 151},
			nil,
		},
		{
			"urdu extended",
			"hi ؎",
			13,
			Count{"7BIT_EX", 1, 5, 5, 155, 150},
			nil,
		},
		{
			"urdu locking and extended",
			"hi ت؎",
			13,
			Count{"7BIT_EX", 1, 6, 6, 152, 146},
			nil,
		},
	}

	for _, p := range patterns {
		f := func(t *testing.T) {
			out, err := NewCount(p.msg, p.nli)
			assert.Equal(t, p.err, err)
			assert.Equal(t, p.out, out)
		}
		t.Run(p.name, f)
	}
}
