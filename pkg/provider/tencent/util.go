package tencent

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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"
)

func sha(s string) string {
	sha := sha1.New()
	sha.Write([]byte(s))
	b := sha.Sum(nil)

	return hex.EncodeToString(b)
}

func hmacSha(k, s string) string {
	enc := hmac.New(sha1.New, []byte(k))
	enc.Write([]byte(s))
	b := enc.Sum(nil)

	return hex.EncodeToString(b)
}

func getSignTime() string {
	now := time.Now()
	expired := now.Add(time.Second * 1800)
	return fmt.Sprintf("%d;%d", now.Unix(), expired.Unix())
}

func getSignature(k string, req *http.Request, signTime string) string {
	httpString := fmt.Sprintf("%s\n%s\n%s\n%s\n", strings.ToLower(req.Method),
		req.URL.Path, getParamsStr(req.URL.RawQuery), getHeadStr(req.Header))

	httpString = sha(httpString)
	signKey := hmacSha(k, signTime)
	signStr := fmt.Sprintf("sha1\n%s\n%s\n", signTime, httpString)

	return hmacSha(signKey, signStr)
}
func getHeadKeys(headers http.Header) string {
	if headers == nil || len(headers) == 0 {
		return ""
	}

	tmp := []string{}
	for k := range headers {
		tmp = append(tmp, strings.ToLower(k))
	}
	sort.Strings(tmp)

	return strings.Join(tmp, ";")
}
func getParamsKeys(p string) string {
	if p == "" {
		return ""
	}
	uv, err := url.ParseQuery(p)
	if err != nil {
		return ""
	}
	tmp := []string{}
	for k := range uv {
		tmp = append(tmp, strings.ToLower(k))
	}
	sort.Strings(tmp)

	return strings.Join(tmp, ";")
}
func getHeadStr(headers http.Header) string {
	if headers == nil || len(headers) == 0 {
		return ""
	}

	tmp := []string{}
	for k, v := range headers {
		str := fmt.Sprintf("%s=%s", strings.ToLower(k), escape(v[0]))
		tmp = append(tmp, str)
	}
	sort.Strings(tmp)

	return strings.Join(tmp, "&")
}

func getParamsStr(p string) string {
	if p == "" {
		return ""
	}
	uv, err := url.ParseQuery(p)
	if err != nil {
		return ""
	}
	tmp := []string{}
	for k, v := range uv {
		str := fmt.Sprintf("%s=%s", strings.ToLower(k), escape(v[0]))
		tmp = append(tmp, str)
	}
	sort.Strings(tmp)

	return strings.Join(tmp, "&")
}

func escape(str string) string {
	return strings.Replace(url.QueryEscape(str), "+", "%20", -1)
}
