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

import (
	"golang.org/x/sys/unix"
)

const (
	// ptypes is the set propagation types.
	ptypes = unix.MS_SHARED | unix.MS_PRIVATE | unix.MS_SLAVE | unix.MS_UNBINDABLE

	// pflags is the full set valid flags for a change propagation call.
	pflags = ptypes | unix.MS_REC | unix.MS_SILENT

	// broflags is the combination of bind and read only
	broflags = unix.MS_BIND | unix.MS_RDONLY
)

// isremount returns true if either device name or flags identify a remount request, false otherwise.
func isremount(device string, flags uintptr) bool {
	switch {
	// We treat device "" and "none" as a remount request to provide compatibility with
	// requests that don't explicitly set MS_REMOUNT such as those manipulating bind mounts.
	case flags&unix.MS_REMOUNT != 0, device == "", device == "none":
		return true
	default:
		return false
	}
}

func mount(device, target, mType string, flags uintptr, data string) error {
	oflags := flags &^ ptypes
	if !isremount(device, flags) || data != "" {
		// Initial call applying all non-propagation flags for mount
		// or remount with changed data
		if err := unix.Mount(device, target, mType, oflags, data); err != nil {
			return &mountError{
				op:     "mount",
				source: device,
				target: target,
				flags:  oflags,
				data:   data,
				err:    err,
			}
		}
	}

	if flags&ptypes != 0 {
		// Change the propagation type.
		if err := unix.Mount("", target, "", flags&pflags, ""); err != nil {
			return &mountError{
				op:     "remount",
				target: target,
				flags:  flags & pflags,
				err:    err,
			}
		}
	}

	if oflags&broflags == broflags {
		// Remount the bind to apply read only.
		if err := unix.Mount("", target, "", oflags|unix.MS_REMOUNT, ""); err != nil {
			return &mountError{
				op:     "remount-ro",
				target: target,
				flags:  oflags | unix.MS_REMOUNT,
				err:    err,
			}
		}
	}

	return nil
}
