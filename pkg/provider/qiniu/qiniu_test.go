package qiniu_test

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
	"testing"

	cfgsvr "github.com/bhojpur/configure/pkg/markup"
	qiniu "github.com/bhojpur/drive/pkg/provider/qiniu"
	"github.com/bhojpur/drive/tests"
)

type Config struct {
	AccessID  string
	AccessKey string
	Region    string
	Bucket    string
	Endpoint  string
}

type AppConfig struct {
	Private Config
	Public  Config
}

var client *qiniu.Client
var privateClient *qiniu.Client

func init() {
	config := AppConfig{}
	cfgsvr.New(&cfgsvr.Config{ENVPrefix: "QINIU"}).Load(&config)
	if len(config.Private.AccessID) == 0 {
		return
	}

	client = qiniu.New(&qiniu.Config{
		AccessID:  config.Public.AccessID,
		AccessKey: config.Public.AccessKey,
		Region:    config.Public.Region,
		Bucket:    config.Public.Bucket,
		Endpoint:  config.Public.Endpoint,
	})
	privateClient = qiniu.New(&qiniu.Config{
		AccessID:   config.Private.AccessID,
		AccessKey:  config.Private.AccessKey,
		Region:     config.Private.Region,
		Bucket:     config.Private.Bucket,
		Endpoint:   config.Private.Endpoint,
		PrivateURL: true,
	})
}

func TestAll(t *testing.T) {
	if client == nil {
		t.Skip(`skip because of no config:


			`)
	}
	clis := []*qiniu.Client{client, privateClient}
	for _, cli := range clis {
		tests.TestAll(cli, t)
	}
}
