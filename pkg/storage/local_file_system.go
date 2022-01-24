package storage

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
	"os"
	"path/filepath"
	"strings"

	"github.com/bhojpur/drive/pkg/model"
)

var baseFolder = "files"

// FileSystem file system storage
type FileSystem struct {
	Base string
}

// NewFileSystem initialize the local file system storage
func NewFileSystem(base string) *FileSystem {
	absBase, err := filepath.Abs(base)
	if err != nil {
		panic("local file system storage's base folder is not initialized")
	}

	return &FileSystem{Base: absBase}
}

// GetFullPath get full path from absolute/relative path
func (fileSystem FileSystem) GetFullPath(path string) string {
	fullPath := path
	if !strings.HasPrefix(path, fileSystem.Base) {
		fullPath, _ = filepath.Abs(filepath.Join(fileSystem.Base, path))
	}
	return fullPath
}

// Get receive file with given path
func (fileSystem FileSystem) Get(path string) (*os.File, error) {
	return os.Open(fileSystem.GetFullPath(path))
}

// GetStream get file as stream
func (fileSystem FileSystem) GetStream(path string) (io.ReadCloser, error) {
	return os.Open(fileSystem.GetFullPath(path))
}

// Put store a reader into given path
func (fileSystem FileSystem) Put(path string, reader io.Reader) (*model.Object, error) {
	var (
		fullPath = fileSystem.GetFullPath(path)
		err      = os.MkdirAll(filepath.Dir(fullPath), os.ModePerm)
	)

	if err != nil {
		return nil, err
	}

	dst, err := os.Create(fullPath)

	if err == nil {
		if seeker, ok := reader.(io.ReadSeeker); ok {
			seeker.Seek(0, 0)
		}
		_, err = io.Copy(dst, reader)
	}

	return &model.Object{Path: path, Name: filepath.Base(path), StorageInterface: fileSystem}, err
}

// Delete delete file
func (fileSystem FileSystem) Delete(path string) error {
	return os.Remove(fileSystem.GetFullPath(path))
}

// List list all objects under current path
func (fileSystem FileSystem) List(path string) ([]*model.Object, error) {
	var (
		objects  []*model.Object
		fullPath = fileSystem.GetFullPath(path)
	)

	filepath.Walk(fullPath, func(path string, info os.FileInfo, err error) error {
		if path == fullPath {
			return nil
		}

		if err == nil && !info.IsDir() {
			modTime := info.ModTime()
			objects = append(objects, &model.Object{
				Path:             strings.TrimPrefix(path, fileSystem.Base),
				Name:             info.Name(),
				LastModified:     &modTime,
				StorageInterface: fileSystem,
			})
		}
		return nil
	})

	return objects, nil
}

// GetEndpoint get endpoint, FileSystem's endpoint is /
func (fileSystem FileSystem) GetEndpoint() string {
	return "/"
}

// GetURL get public accessible URL
func (fileSystem FileSystem) GetURL(path string) (url string, err error) {
	return path, nil
}

func NewLocalFileSystemStorageProvider(clientId string, clientSecret string, region string, bucket string, endpoint string) model.StorageInterface {
	return NewFileSystem(baseFolder)
}
