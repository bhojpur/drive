//go:build linux || (freebsd && cgo) || (openbsd && cgo) || (darwin && cgo)
// +build linux freebsd,cgo openbsd,cgo darwin,cgo

package mountinfo

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
	"fmt"
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

func mountedByStat(path string) (bool, error) {
	var st unix.Stat_t

	if err := unix.Lstat(path, &st); err != nil {
		return false, &os.PathError{Op: "stat", Path: path, Err: err}
	}
	dev := st.Dev
	parent := filepath.Dir(path)
	if err := unix.Lstat(parent, &st); err != nil {
		return false, &os.PathError{Op: "stat", Path: parent, Err: err}
	}
	if dev != st.Dev {
		// Device differs from that of parent,
		// so definitely a mount point.
		return true, nil
	}
	// NB: this does not detect bind mounts on Linux.
	return false, nil
}

func normalizePath(path string) (realPath string, err error) {
	if realPath, err = filepath.Abs(path); err != nil {
		return "", fmt.Errorf("unable to get absolute path for %q: %w", path, err)
	}
	if realPath, err = filepath.EvalSymlinks(realPath); err != nil {
		return "", fmt.Errorf("failed to canonicalise path for %q: %w", path, err)
	}
	if _, err := os.Stat(realPath); err != nil {
		return "", fmt.Errorf("failed to stat target of %q: %w", path, err)
	}
	return realPath, nil
}

func mountedByMountinfo(path string) (bool, error) {
	entries, err := GetMounts(SingleEntryFilter(path))
	if err != nil {
		return false, err
	}

	return len(entries) > 0, nil
}
