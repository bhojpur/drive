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

import "github.com/bhojpur/drive/pkg/mountinfo"

// MakeShared ensures a mounted filesystem has the SHARED mount option enabled.
// See the supported options in flags.go for further reference.
func MakeShared(mountPoint string) error {
	return ensureMountedAs(mountPoint, SHARED)
}

// MakeRShared ensures a mounted filesystem has the RSHARED mount option enabled.
// See the supported options in flags.go for further reference.
func MakeRShared(mountPoint string) error {
	return ensureMountedAs(mountPoint, RSHARED)
}

// MakePrivate ensures a mounted filesystem has the PRIVATE mount option enabled.
// See the supported options in flags.go for further reference.
func MakePrivate(mountPoint string) error {
	return ensureMountedAs(mountPoint, PRIVATE)
}

// MakeRPrivate ensures a mounted filesystem has the RPRIVATE mount option
// enabled. See the supported options in flags.go for further reference.
func MakeRPrivate(mountPoint string) error {
	return ensureMountedAs(mountPoint, RPRIVATE)
}

// MakeSlave ensures a mounted filesystem has the SLAVE mount option enabled.
// See the supported options in flags.go for further reference.
func MakeSlave(mountPoint string) error {
	return ensureMountedAs(mountPoint, SLAVE)
}

// MakeRSlave ensures a mounted filesystem has the RSLAVE mount option enabled.
// See the supported options in flags.go for further reference.
func MakeRSlave(mountPoint string) error {
	return ensureMountedAs(mountPoint, RSLAVE)
}

// MakeUnbindable ensures a mounted filesystem has the UNBINDABLE mount option
// enabled. See the supported options in flags.go for further reference.
func MakeUnbindable(mountPoint string) error {
	return ensureMountedAs(mountPoint, UNBINDABLE)
}

// MakeRUnbindable ensures a mounted filesystem has the RUNBINDABLE mount
// option enabled. See the supported options in flags.go for further reference.
func MakeRUnbindable(mountPoint string) error {
	return ensureMountedAs(mountPoint, RUNBINDABLE)
}

// MakeMount ensures that the file or directory given is a mount point,
// bind mounting it to itself it case it is not.
func MakeMount(mnt string) error {
	mounted, err := mountinfo.Mounted(mnt)
	if err != nil {
		return err
	}
	if mounted {
		return nil
	}

	return mount(mnt, mnt, "none", uintptr(BIND), "")
}

func ensureMountedAs(mnt string, flags int) error {
	if err := MakeMount(mnt); err != nil {
		return err
	}

	return mount("", mnt, "none", uintptr(flags), "")
}
