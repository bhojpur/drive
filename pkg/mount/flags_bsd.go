//go:build freebsd || openbsd
// +build freebsd openbsd

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

import "golang.org/x/sys/unix"

const (
	// RDONLY will mount the filesystem as read-only.
	RDONLY = unix.MNT_RDONLY

	// NOSUID will not allow set-user-identifier or set-group-identifier bits to
	// take effect.
	NOSUID = unix.MNT_NOSUID

	// NOEXEC will not allow execution of any binaries on the mounted file system.
	NOEXEC = unix.MNT_NOEXEC

	// SYNCHRONOUS will allow any I/O to the file system to be done synchronously.
	SYNCHRONOUS = unix.MNT_SYNCHRONOUS

	// NOATIME will not update the file access time when reading from a file.
	NOATIME = unix.MNT_NOATIME
)

// These flags are unsupported.
const (
	BIND        = 0
	DIRSYNC     = 0
	MANDLOCK    = 0
	NODEV       = 0
	NODIRATIME  = 0
	UNBINDABLE  = 0
	RUNBINDABLE = 0
	PRIVATE     = 0
	RPRIVATE    = 0
	SHARED      = 0
	RSHARED     = 0
	SLAVE       = 0
	RSLAVE      = 0
	RBIND       = 0
	RELATIME    = 0
	REMOUNT     = 0
	STRICTATIME = 0
	mntDetach   = 0
)
