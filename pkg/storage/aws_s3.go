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
	awss3 "github.com/aws/aws-sdk-go/service/s3"
	"github.com/bhojpur/drive/pkg/model"
	"github.com/bhojpur/drive/pkg/provider/s3"
)

func NewAwsS3StorageProvider(clientId string, clientSecret string, region string, bucket string, endpoint string) model.StorageInterface {
	sp := s3.New(&s3.Config{
		AccessID:  clientId,
		AccessKey: clientSecret,
		Region:    region,
		Bucket:    bucket,
		Endpoint:  endpoint,
		ACL:       awss3.BucketCannedACLPublicRead,
	})

	return sp
}
