package tencent_test

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
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"

	tencent "github.com/bhojpur/drive/pkg/provider/tencent"
	"github.com/bhojpur/drive/tests"
)

func TestClient_Get(t *testing.T) {

}

var client *tencent.Client

func init() {
	client = tencent.New(&tencent.Config{
		AppID:     "1252882253",
		AccessID:  "AKIDToxukQWBG8nGXcBN8i662nOo12sc5Wjl",
		AccessKey: "40jNrBf5mLiuuiU8HH7lDTXP5at00sbA",
		Bucket:    "tets-1252882253",
		Region:    "ap-shanghai",
		ACL:       "public-read", // private，public-read-write，public-read；默认值：private
		//Endpoint:  config.Public.Endpoint,
	})
}

func TestClient_Put(t *testing.T) {
	f, err := ioutil.ReadFile("/home/shashi.rai/Downloads/2.png")
	if err != nil {
		t.Error(err)
		return
	}

	client.Put("test.png", bytes.NewReader(f))
}

func TestClient_Put2(t *testing.T) {
	tests.TestAll(client, t)
}

func TestClient_Delete(t *testing.T) {
	fmt.Println(client.Delete("test.png"))
}
