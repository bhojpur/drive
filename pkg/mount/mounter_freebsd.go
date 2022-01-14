//go:build freebsd && cgo
// +build freebsd,cgo

package mount

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

/*
#include <errno.h>
#include <stdlib.h>
#include <string.h>
#include <sys/_iovec.h>
#include <sys/mount.h>
#include <sys/param.h>
*/
import "C"

import (
	"strings"
	"syscall"
	"unsafe"
)

func allocateIOVecs(options []string) []C.struct_iovec {
	out := make([]C.struct_iovec, len(options))
	for i, option := range options {
		out[i].iov_base = unsafe.Pointer(C.CString(option))
		out[i].iov_len = C.size_t(len(option) + 1)
	}
	return out
}

func mount(device, target, mType string, flag uintptr, data string) error {
	isNullFS := false

	xs := strings.Split(data, ",")
	for _, x := range xs {
		if x == "bind" {
			isNullFS = true
		}
	}

	options := []string{"fspath", target}
	if isNullFS {
		options = append(options, "fstype", "nullfs", "target", device)
	} else {
		options = append(options, "fstype", mType, "from", device)
	}
	rawOptions := allocateIOVecs(options)
	for _, rawOption := range rawOptions {
		defer C.free(rawOption.iov_base)
	}

	if errno := C.nmount(&rawOptions[0], C.uint(len(options)), C.int(flag)); errno != 0 {
		return &mountError{
			op:     "mount",
			source: device,
			target: target,
			flags:  flag,
			err:    syscall.Errno(errno),
		}
	}
	return nil
}
