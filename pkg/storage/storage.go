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
	"github.com/bhojpur/drive/pkg/model"
)

func GetStorageProvider(providerType string, clientId string, clientSecret string, region string, bucket string, endpoint string) model.StorageInterface {
	switch providerType {
	case "Local File System":
		return NewLocalFileSystemStorageProvider(clientId, clientSecret, region, bucket, endpoint)
	case "AWS S3":
		return NewAwsS3StorageProvider(clientId, clientSecret, region, bucket, endpoint)
	case "Aliyun OSS":
		return NewAliyunOssStorageProvider(clientId, clientSecret, region, bucket, endpoint)
	case "Tencent Cloud COS":
		return NewTencentCloudCosStorageProvider(clientId, clientSecret, region, bucket, endpoint)
	}

	return nil
}
