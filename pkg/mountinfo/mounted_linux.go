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
	"os"
	"path/filepath"

	"golang.org/x/sys/unix"
)

// mountedByOpenat2 is a method of detecting a mount that works for all kinds
// of mounts (incl. bind mounts), but requires a recent (v5.6+) linux kernel.
func mountedByOpenat2(path string) (bool, error) {
	dir, last := filepath.Split(path)

	dirfd, err := unix.Openat2(unix.AT_FDCWD, dir, &unix.OpenHow{
		Flags: unix.O_PATH | unix.O_CLOEXEC,
	})
	if err != nil {
		return false, &os.PathError{Op: "openat2", Path: dir, Err: err}
	}
	fd, err := unix.Openat2(dirfd, last, &unix.OpenHow{
		Flags:   unix.O_PATH | unix.O_CLOEXEC | unix.O_NOFOLLOW,
		Resolve: unix.RESOLVE_NO_XDEV,
	})
	_ = unix.Close(dirfd)
	switch err { //nolint:errorlint // unix errors are bare
	case nil: // definitely not a mount
		_ = unix.Close(fd)
		return false, nil
	case unix.EXDEV: // definitely a mount
		return true, nil
	}
	// not sure
	return false, &os.PathError{Op: "openat2", Path: path, Err: err}
}

func mounted(path string) (bool, error) {
	path, err := normalizePath(path)
	if err != nil {
		return false, err
	}
	// Try a fast path, using openat2() with RESOLVE_NO_XDEV.
	mounted, err := mountedByOpenat2(path)
	if err == nil {
		return mounted, nil
	}
	// Another fast path: compare st.st_dev fields.
	mounted, err = mountedByStat(path)
	// This does not work for bind mounts, so false negative
	// is possible, therefore only trust if return is true.
	if mounted && err == nil {
		return mounted, nil
	}

	// Fallback to parsing mountinfo
	return mountedByMountinfo(path)
}
