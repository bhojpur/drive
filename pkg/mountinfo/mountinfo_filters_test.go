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

import "testing"

func TestPrefixFilter(t *testing.T) {
	tests := []struct {
		prefix     string
		mountPoint string
		shouldSkip bool
	}{
		{prefix: "/a", mountPoint: "/a", shouldSkip: false},
		{prefix: "/a", mountPoint: "/a/b", shouldSkip: false},
		{prefix: "/a", mountPoint: "/aa", shouldSkip: true},
		{prefix: "/a", mountPoint: "/aa/b", shouldSkip: true},

		// invalid prefix: prefix path must be cleaned and have no trailing slash
		{prefix: "/a/", mountPoint: "/a", shouldSkip: true},
		{prefix: "/a/", mountPoint: "/a/b", shouldSkip: true},
	}
	for _, tc := range tests {
		filter := PrefixFilter(tc.prefix)
		skip, _ := filter(&Info{Mountpoint: tc.mountPoint})
		if skip != tc.shouldSkip {
			if tc.shouldSkip {
				t.Errorf("prefix %q: expected %q to be skipped", tc.prefix, tc.mountPoint)
			} else {
				t.Errorf("prefix %q: expected %q not to be skipped", tc.prefix, tc.mountPoint)
			}
		}
	}
}
