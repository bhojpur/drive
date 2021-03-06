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
)

// GetMounts retrieves a list of mounts for the current running process,
// with an optional filter applied (use nil for no filter).
func GetMounts(f FilterFunc) ([]*Info, error) {
	return parseMountTable(f)
}

// Mounted determines if a specified path is a mount point. In case of any
// error, false (and an error) is returned.
//
// The non-existent path returns an error. If a caller is not interested
// in this particular error, it should handle it separately using e.g.
// errors.Is(err, os.ErrNotExist).
func Mounted(path string) (bool, error) {
	// root is always mounted
	if path == string(os.PathSeparator) {
		return true, nil
	}
	return mounted(path)
}

// Info reveals information about a particular mounted filesystem. This
// struct is populated from the content in the /proc/<pid>/mountinfo file.
type Info struct {
	// ID is a unique identifier of the mount (may be reused after umount).
	ID int

	// Parent is the ID of the parent mount (or of self for the root
	// of this mount namespace's mount tree).
	Parent int

	// Major and Minor are the major and the minor components of the Dev
	// field of unix.Stat_t structure returned by unix.*Stat calls for
	// files on this filesystem.
	Major, Minor int

	// Root is the pathname of the directory in the filesystem which forms
	// the root of this mount.
	Root string

	// Mountpoint is the pathname of the mount point relative to the
	// process's root directory.
	Mountpoint string

	// Options is a comma-separated list of mount options.
	Options string

	// Optional are zero or more fields of the form "tag[:value]",
	// separated by a space.  Currently, the possible optional fields are
	// "shared", "master", "propagate_from", and "unbindable". For more
	// information, see mount_namespaces(7) Linux man page.
	Optional string

	// FSType is the filesystem type in the form "type[.subtype]".
	FSType string

	// Source is filesystem-specific information, or "none".
	Source string

	// VFSOptions is a comma-separated list of superblock options.
	VFSOptions string
}
