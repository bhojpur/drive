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

import "strings"

// FilterFunc is a type defining a callback function for GetMount(),
// used to filter out mountinfo entries we're not interested in,
// and/or stop further processing if we found what we wanted.
//
// It takes a pointer to the Info struct (fully populated with all available
// fields on the GOOS platform), and returns two booleans:
//
// skip: true if the entry should be skipped;
//
// stop: true if parsing should be stopped after the entry.
type FilterFunc func(*Info) (skip, stop bool)

// PrefixFilter discards all entries whose mount points do not start with, or
// are equal to the path specified in prefix. The prefix path must be absolute,
// have all symlinks resolved, and cleaned (i.e. no extra slashes or dots).
//
// PrefixFilter treats prefix as a path, not a partial prefix, which means that
// given "/foo", "/foo/bar" and "/foobar" entries, PrefixFilter("/foo") returns
// "/foo" and "/foo/bar", and discards "/foobar".
func PrefixFilter(prefix string) FilterFunc {
	return func(m *Info) (bool, bool) {
		skip := !strings.HasPrefix(m.Mountpoint+"/", prefix+"/")
		return skip, false
	}
}

// SingleEntryFilter looks for a specific entry.
func SingleEntryFilter(mp string) FilterFunc {
	return func(m *Info) (bool, bool) {
		if m.Mountpoint == mp {
			return false, true // don't skip, stop now
		}
		return true, false // skip, keep going
	}
}

// ParentsFilter returns all entries whose mount points
// can be parents of a path specified, discarding others.
//
// For example, given /var/lib/bhojpur/something, entries
// like /var/lib/bhojpur, /var and / are returned.
func ParentsFilter(path string) FilterFunc {
	return func(m *Info) (bool, bool) {
		skip := !strings.HasPrefix(path, m.Mountpoint)
		return skip, false
	}
}

// FSTypeFilter returns all entries that match provided fstype(s).
func FSTypeFilter(fstype ...string) FilterFunc {
	return func(m *Info) (bool, bool) {
		for _, t := range fstype {
			if m.FSType == t {
				return false, false // don't skip, keep going
			}
		}
		return true, false // skip, keep going
	}
}
