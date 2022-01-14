//go:build openbsd && cgo
// +build openbsd,cgo

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
#include <sys/types.h>
#include <sys/mount.h>
*/
import "C"

import (
	"fmt"
	"syscall"
	"unsafe"
)

func createExportInfo(readOnly bool) C.struct_export_args {
	exportFlags := C.int(0)
	if readOnly {
		exportFlags = C.MNT_EXRDONLY
	}
	out := C.struct_export_args{
		ex_root:  0,
		ex_flags: exportFlags,
	}
	return out
}

func createUfsArgs(device string, readOnly bool) unsafe.Pointer {
	out := &C.struct_ufs_args{
		fspec:       C.CString(device),
		export_info: createExportInfo(readOnly),
	}
	return unsafe.Pointer(out)
}

func mount(device, target, mType string, flag uintptr, data string) error {
	readOnly := flag&RDONLY != 0

	var fsArgs unsafe.Pointer

	switch mType {
	case "ffs":
		fsArgs = createUfsArgs(device, readOnly)
	default:
		return &mountError{
			op:     "mount",
			source: device,
			target: target,
			flags:  flag,
			err:    fmt.Errorf("unsupported file system type: %s", mType),
		}
	}

	if errno := C.mount(C.CString(mType), C.CString(target), C.int(flag), fsArgs); errno != 0 {
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
