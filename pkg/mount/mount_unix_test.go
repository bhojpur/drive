//go:build !darwin && !windows
// +build !darwin,!windows

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
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/bhojpur/drive/pkg/mountinfo"
)

func TestMountOptionsParsing(t *testing.T) {
	options := "noatime,ro,noexec,size=10k"

	flag, data := parseOptions(options)

	if data != "size=10k" {
		t.Fatalf("Expected size=10 got %s", data)
	}

	expected := NOATIME | RDONLY | NOEXEC

	if flag != expected {
		t.Fatalf("Expected %d got %d", expected, flag)
	}
}

func TestMounted(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("root required")
	}

	tmp := path.Join(os.TempDir(), "mount-tests")
	if err := os.MkdirAll(tmp, 0o777); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	var (
		sourceDir  = path.Join(tmp, "source")
		targetDir  = path.Join(tmp, "target")
		sourcePath = path.Join(sourceDir, "file.txt")
		targetPath = path.Join(targetDir, "file.txt")
	)

	if err := os.Mkdir(sourceDir, 0o777); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(targetDir, 0o777); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(sourcePath, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(targetPath, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Mount(sourceDir, targetDir, "none", "bind,rw"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Unmount(targetDir); err != nil {
			t.Fatal(err)
		}
	}()

	mounted, err := mountinfo.Mounted(targetDir)
	if err != nil {
		t.Fatal(err)
	}
	if !mounted {
		t.Fatalf("Expected %s to be mounted", targetDir)
	}
	if _, err := os.Stat(targetDir); err != nil {
		t.Fatal(err)
	}
}

func TestMountTmpfsOptions(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("root required")
	}

	testCases := []struct {
		opts       string
		expected   string
		unexpected string
	}{
		{
			opts:       "exec",
			unexpected: "noexec",
		},
		{
			opts:       "noexec",
			expected:   "noexec",
			unexpected: "exec",
		},
	}

	target := path.Join(os.TempDir(), "mount-tmpfs-tests-"+t.Name())
	if err := os.MkdirAll(target, 0o777); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(target)

	for _, tc := range testCases {
		t.Run(tc.opts, func(t *testing.T) {
			if err := Mount("tmpfs", target, "tmpfs", tc.opts); err != nil {
				t.Fatal(err)
			}
			defer ensureUnmount(t, target)

			mounts, err := mountinfo.GetMounts(mountinfo.SingleEntryFilter(target))
			if err != nil {
				t.Fatal(err)
			}
			if len(mounts) != 1 {
				t.Fatal("Mount point ", target, " not found")
			}
			entry := mounts[0]
			opts := "," + entry.Options + ","
			if tc.expected != "" && !strings.Contains(opts, ","+tc.expected+",") {
				t.Fatal("Expected option ", tc.expected, " missing from ", entry.Options)
			}
			if tc.unexpected != "" && strings.Contains(opts, ","+tc.unexpected+",") {
				t.Fatal("Unexpected option ", tc.unexpected, " in ", entry.Options)
			}
		})
	}
}

func TestMountReadonly(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("root required")
	}

	tmp := path.Join(os.TempDir(), "mount-tests")
	if err := os.MkdirAll(tmp, 0o777); err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	var (
		sourceDir  = path.Join(tmp, "source")
		targetDir  = path.Join(tmp, "target")
		sourcePath = path.Join(sourceDir, "file.txt")
		targetPath = path.Join(targetDir, "file.txt")
	)

	if err := os.Mkdir(sourceDir, 0o777); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(targetDir, 0o777); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(sourcePath, []byte("hello"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := ioutil.WriteFile(targetPath, nil, 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Mount(sourceDir, targetDir, "none", "bind,ro"); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := Unmount(targetDir); err != nil {
			t.Fatal(err)
		}
	}()

	if err := ioutil.WriteFile(targetPath, []byte("hello"), 0o644); err == nil {
		t.Fatal("Should not be able to open a ro file as rw")
	}
}

func TestMergeTmpfsOptions(t *testing.T) {
	options := []string{"noatime", "ro", "size=10k", "defaults", "noexec", "atime", "defaults", "rw", "rprivate", "size=1024k", "slave", "exec"}
	expected := []string{"atime", "rw", "size=1024k", "slave", "exec"}
	merged, err := MergeTmpfsOptions(options)
	if err != nil {
		t.Fatal(err)
	}
	if len(expected) != len(merged) {
		t.Fatalf("Expected %s got %s", expected, merged)
	}
	for index := range merged {
		if merged[index] != expected[index] {
			t.Fatalf("Expected %s for the %dth option, got %s", expected, index, merged)
		}
	}

	options = []string{"noatime", "ro", "size=10k", "atime", "rw", "rprivate", "size=1024k", "slave", "size", "exec"}
	_, err = MergeTmpfsOptions(options)
	if err == nil {
		t.Fatal("Expected error got nil")
	}
}

func TestRecursiveUnmountTooGreedy(t *testing.T) {
	if os.Getuid() != 0 {
		t.Skip("root required")
	}

	tmp, err := ioutil.TempDir("", t.Name())
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	// Create a bunch of tmpfs mounts. Make sure "dir" itself is not
	// a mount point, or we'll hit the fast path in RecursiveUnmount.
	dirs := []string{"dir-other", "dir/subdir1", "dir/subdir1/subsub", "dir/subdir2/subsub"}
	for _, d := range dirs {
		dir := path.Join(tmp, d)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := Mount("tmpfs", dir, "tmpfs", ""); err != nil {
			t.Fatal(err)
		}
		//nolint:errcheck
		defer Unmount(dir)
	}
	// sanity check
	mounted, err := mountinfo.Mounted(path.Join(tmp, "dir-other"))
	if err != nil {
		t.Fatalf("[pre-check] error from mountinfo.mounted: %v", err)
	}
	if !mounted {
		t.Fatal("[pre-check] expected dir-other to be mounted, but it's not")
	}
	// Unmount dir, make sure dir-other is still mounted.
	if err := RecursiveUnmount(path.Join(tmp, "dir")); err != nil {
		t.Fatal(err)
	}
	mounted, err = mountinfo.Mounted(path.Join(tmp, "dir-other"))
	if err != nil {
		t.Fatalf("error from mountinfo.mounted: %v", err)
	}
	if !mounted {
		t.Fatal("expected dir-other to be mounted, but it's not")
	}
}
