// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//	http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package cache

import (
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/stretchr/testify/assert"
)

const (
	testRegion        = "test-region"
	testCacheFilename = "cache.json"
	testAccessKey     = "accessKey"
	testSecretKey     = "secretKey"
	testToken         = "token"
	// base64 MD5 sum of Credentials struct
	testCredentialHash = "YWNjZXNzS2V51B2M2Y8AsgTpgAmY7PhCfg=="
)

func TestFactoryBuildFileCache(t *testing.T) {
	config := aws.Config{
		Region:      testRegion,
		Credentials: credentials.NewStaticCredentialsProvider(testAccessKey, testSecretKey, testToken),
	}

	cache := BuildCredentialsCache(config, "")
	assert.NotNil(t, cache)

	fileCache, ok := cache.(*fileCredentialCache)

	assert.True(t, ok, "built cache is not a fileCredentialsCache")
	assert.Equal(t, fileCache.cachePrefixKey, fmt.Sprintf("%s-%s-", testRegion, testCredentialHash))
	assert.Equal(t, fileCache.filename, testCacheFilename)
}

func TestFactoryBuildNullCacheWithoutCredentials(t *testing.T) {
	config := aws.Config{
		Region:      testRegion,
		Credentials: aws.AnonymousCredentials{},
	}

	cache := BuildCredentialsCache(config, "")
	assert.NotNil(t, cache)

	_, ok := cache.(*nullCredentialsCache)
	assert.True(t, ok, "built cache is a nullCredentialsCache")
}

func TestFactoryBuildNullCache(t *testing.T) {
	os.Setenv("AWS_ECR_DISABLE_CACHE", "1")
	defer os.Setenv("AWS_ECR_DISABLE_CACHE", "1")

	config := aws.Config{Region: testRegion}

	cache := BuildCredentialsCache(config, "")
	assert.NotNil(t, cache)
	_, ok := cache.(*nullCredentialsCache)
	assert.True(t, ok, "built cache is a nullCredentialsCache")
}
