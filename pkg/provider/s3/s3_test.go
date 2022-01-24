package s3_test

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
	"fmt"
	"testing"

	awss3 "github.com/aws/aws-sdk-go/service/s3"
	cfgsvr "github.com/bhojpur/configure/pkg/markup"
	"github.com/bhojpur/drive/pkg/provider/s3"
	"github.com/bhojpur/drive/tests"
)

type Config struct {
	AccessID  string `env:"BHOJPUR_AWS_ACCESS_KEY_ID"`
	AccessKey string `env:"BHOJPUR_AWS_SECRET_ACCESS_KEY"`
	Region    string `env:"BHOJPUR_AWS_REGION"`
	Bucket    string `env:"BHOJPUR_AWS_BUCKET"`
	Endpoint  string `env:"BHOJPUR_AWS_ENDPOINT"`
}

var (
	client *s3.Client
	config = Config{}
)

func init() {
	cfgsvr.Load(&config)

	client = s3.New(&s3.Config{AccessID: config.AccessID, AccessKey: config.AccessKey, Region: config.Region, Bucket: config.Bucket, Endpoint: config.Endpoint})
}

func TestAll(t *testing.T) {
	fmt.Println("testing S3 with public ACL")
	tests.TestAll(client, t)

	fmt.Println("testing S3 with private ACL")
	privateClient := s3.New(&s3.Config{AccessID: config.AccessID, AccessKey: config.AccessKey, Region: config.Region, Bucket: config.Bucket, ACL: awss3.BucketCannedACLPrivate, Endpoint: config.Endpoint})
	tests.TestAll(privateClient, t)

	fmt.Println("testing S3 with AuthenticatedRead ACL")
	authenticatedReadClient := s3.New(&s3.Config{AccessID: config.AccessID, AccessKey: config.AccessKey, Region: config.Region, Bucket: config.Bucket, ACL: awss3.BucketCannedACLAuthenticatedRead, Endpoint: config.Endpoint})
	tests.TestAll(authenticatedReadClient, t)
}

func TestToRelativePath(t *testing.T) {
	urlMap := map[string]string{
		"https://mybucket.s3.amazonaws.com/myobject.ext": "/myobject.ext",
		"https://example.drive.bhojpur.net/myobject.ext": "/myobject.ext",
		"//mybucket.s3.amazonaws.com/myobject.ext":       "/myobject.ext",
		"http://mybucket.s3.amazonaws.com/myobject.ext":  "/myobject.ext",
		"myobject.ext": "/myobject.ext",
	}

	for url, path := range urlMap {
		if client.ToRelativePath(url) != path {
			t.Errorf("%v's relative path should be %v, but got %v", url, path, client.ToRelativePath(url))
		}
	}
}

func TestToRelativePathWithS3ForcePathStyle(t *testing.T) {
	urlMap := map[string]string{
		"https://s3.amazonaws.com/mybucket/myobject.ext": "/myobject.ext",
		"https://example.drive.bhojpur.net/myobject.ext": "/myobject.ext",
		"//s3.amazonaws.com/mybucket/myobject.ext":       "/myobject.ext",
		"http://s3.amazonaws.com/mybucket/myobject.ext":  "/myobject.ext",
		"/mybucket/myobject.ext":                         "/myobject.ext",
		"myobject.ext":                                   "/myobject.ext",
	}

	client := s3.New(&s3.Config{AccessID: config.AccessID, AccessKey: config.AccessKey, Region: config.Region, Bucket: "mybucket", S3ForcePathStyle: true, Endpoint: config.Endpoint})

	for url, path := range urlMap {
		if client.ToRelativePath(url) != path {
			t.Errorf("%v's relative path should be %v, but got %v", url, path, client.ToRelativePath(url))
		}
	}
}
