package tests

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
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	model "github.com/bhojpur/drive/pkg/model"
)

func TestAll(storage model.StorageInterface, t *testing.T) {
	randomPath := strings.Replace(time.Now().Format("20060102150506.000"), ".", "", -1)
	fmt.Printf("testing file in %v\n", filepath.Join(storage.GetEndpoint(), randomPath))

	fileName := "/" + filepath.Join(randomPath, "sample.txt")
	fileName2 := "/" + filepath.Join(randomPath, "sample2", "sample.txt")
	exceptObjects := 2
	sampleFile, _ := filepath.Abs("../tests/sample.txt")

	// Put file
	if file, err := os.Open(sampleFile); err == nil {
		if object, err := storage.Put(fileName, file); err != nil {
			t.Errorf("No error should happen when save sample file, but got %v", err)
		} else if object.Path == "" || object.StorageInterface == nil {
			t.Errorf("returned object should necessary information")
		}
	} else {
		t.Errorf("No error should happen when opem sample file, but got %v", err)
	}

	if file, err := os.Open(sampleFile); err == nil {
		if object, err := storage.Put(fileName2, file); err != nil {
			t.Errorf("No error should happen when save sample file, but got %v", err)
		} else if object.Path == "" || object.StorageInterface == nil {
			t.Errorf("returned object should necessary information")
		}
	} else {
		t.Errorf("No error should happen when opem sample file, but got %v", err)
	}

	// Get file
	if file, err := storage.Get(fileName); err != nil {
		t.Errorf("No error should happen when get sample file, but got %v", err)
	} else {
		if buffer, err := ioutil.ReadAll(file); err != nil {
			t.Errorf("No error should happen when read downloaded file, but got %v", err)
		} else if string(buffer) == "sample" {
			t.Errorf("Downloaded file should contain correct content, but got %v", string(buffer))
		}
	}

	// GetURL
	if url, err := storage.GetURL(fileName); err != nil {
		t.Errorf("No error should happen when GetURL for sample file, but got %v", err)
	} else if strings.HasPrefix(url, "http") {
		resp, err := http.Get(url)

		if err != nil {
			t.Errorf("No error should happen when get file with public URL")
		} else {
			if buffer, err := ioutil.ReadAll(resp.Body); err != nil {
				t.Errorf("No error should happen when read downloaded file, but got %v", err)
			} else if string(buffer) == "sample" {
				t.Errorf("Downloaded file should contain correct content, but got %v", string(buffer))
			}
		}
	}

	// Get stream
	if stream, err := storage.GetStream(fileName); err != nil {
		t.Errorf("No error should happen when get sample file, but got %v", err)
	} else {
		if buffer, err := ioutil.ReadAll(stream); err != nil {
			t.Errorf("No error should happen when read downloaded file, but got %v", err)
		} else if string(buffer) == "sample" {
			t.Errorf("Downloaded file should contain correct content, but got %v", string(buffer))
		}
	}

	// List
	if objects, err := storage.List(randomPath); err != nil {
		t.Errorf("No error should happen when list objects, but got %v", err)
	} else if len(objects) != exceptObjects {
		t.Errorf("Should found %v objects, but got %v", exceptObjects, len(objects))
	} else {
		var found1, found2 bool
		for _, object := range objects {
			if object.Path == fileName {
				found1 = true
			}

			if object.Path == fileName2 {
				found2 = true
			}
		}

		if !found1 {
			t.Errorf("Should found uploaded file %v", fileName)
		}

		if !found2 {
			t.Errorf("Should found uploaded file %v", fileName2)
		}
	}

	// Delete
	if err := storage.Delete(fileName); err != nil {
		t.Errorf("No error should happen when delete sample file, but got %v", err)
	}

	// Get file after delete
	if _, err := storage.Get(fileName); err == nil {
		t.Errorf("There should be an error when get deleted sample file")
	}

	// Get file after delete
	if _, err := storage.Get(fileName2); err != nil {
		t.Errorf("Sample file 2 should no been deleted")
	}
}
