package aliyun

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
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	aliyun "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/bhojpur/drive/pkg/model"
)

// Client Aliyun storage
type Client struct {
	*aliyun.Bucket
	Config *Config
}

// Config Aliyun client config
type Config struct {
	AccessID      string
	AccessKey     string
	Region        string
	Bucket        string
	Endpoint      string
	ACL           aliyun.ACLType
	ClientOptions []aliyun.ClientOption
	UseCname      bool
}

// New initialize Aliyun storage
func New(config *Config) *Client {
	var (
		err    error
		client = &Client{Config: config}
	)

	if config.Endpoint == "" {
		config.Endpoint = "oss-cn-hangzhou.aliyuncs.com"
	}

	if config.ACL == "" {
		config.ACL = aliyun.ACLPublicRead
	}

	if config.UseCname {
		config.ClientOptions = append(config.ClientOptions, aliyun.UseCname(config.UseCname))
	}

	Aliyun, err := aliyun.New(config.Endpoint, config.AccessID, config.AccessKey, config.ClientOptions...)

	if err == nil {
		client.Bucket, err = Aliyun.Bucket(config.Bucket)
	}

	if err != nil {
		panic(err)
	}

	return client
}

// Get receive file with given path
func (client Client) Get(path string) (file *os.File, err error) {
	readCloser, err := client.GetStream(path)

	if err == nil {
		if file, err = ioutil.TempFile("/tmp", "ali"); err == nil {
			defer readCloser.Close()
			_, err = io.Copy(file, readCloser)
			file.Seek(0, 0)
		}
	}

	return file, err
}

// GetStream get file as stream
func (client Client) GetStream(path string) (io.ReadCloser, error) {
	return client.Bucket.GetObject(client.ToRelativePath(path))
}

// Put store a reader into given path
func (client Client) Put(urlPath string, reader io.Reader) (*model.Object, error) {
	if seeker, ok := reader.(io.ReadSeeker); ok {
		seeker.Seek(0, 0)
	}

	err := client.Bucket.PutObject(client.ToRelativePath(urlPath), reader, aliyun.ACL(client.Config.ACL))
	now := time.Now()

	return &model.Object{
		Path:             urlPath,
		Name:             filepath.Base(urlPath),
		LastModified:     &now,
		StorageInterface: client,
	}, err
}

// Delete delete file
func (client Client) Delete(path string) error {
	return client.Bucket.DeleteObject(client.ToRelativePath(path))
}

// List list all objects under current path
func (client Client) List(path string) ([]*model.Object, error) {
	var objects []*model.Object

	results, err := client.Bucket.ListObjects(aliyun.Prefix(path))

	if err == nil {
		for _, obj := range results.Objects {
			objects = append(objects, &model.Object{
				Path:             "/" + client.ToRelativePath(obj.Key),
				Name:             filepath.Base(obj.Key),
				LastModified:     &obj.LastModified,
				StorageInterface: client,
			})
		}
	}

	return objects, err
}

// GetEndpoint get endpoint, FileSystem's endpoint is /
func (client Client) GetEndpoint() string {
	if client.Config.Endpoint != "" {
		if strings.HasSuffix(client.Config.Endpoint, "aliyuncs.com") {
			return client.Config.Bucket + "." + client.Config.Endpoint
		}
		return client.Config.Endpoint
	}

	endpoint := client.Bucket.Client.Config.Endpoint
	for _, prefix := range []string{"https://", "http://"} {
		endpoint = strings.TrimPrefix(endpoint, prefix)
	}

	return client.Config.Bucket + "." + endpoint
}

var urlRegexp = regexp.MustCompile(`(https?:)?//((\w+).)+(\w+)/`)

// ToRelativePath process path to relative path
func (client Client) ToRelativePath(urlPath string) string {
	if urlRegexp.MatchString(urlPath) {
		if u, err := url.Parse(urlPath); err == nil {
			return strings.TrimPrefix(u.Path, "/")
		}
	}

	return strings.TrimPrefix(urlPath, "/")
}

// GetURL get public accessible URL
func (client Client) GetURL(path string) (url string, err error) {
	if client.Config.ACL == aliyun.ACLPrivate {
		return client.Bucket.SignURL(client.ToRelativePath(path), aliyun.HTTPGet, 60*60) // 1 hour
	}
	return path, nil
}
