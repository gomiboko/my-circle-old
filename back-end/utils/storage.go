package utils

import (
	"crypto/sha256"
	"fmt"
	"os"

	"github.com/gomiboko/my-circle/consts"
)

func CreateHashedStorageKey(dir string, prefix string, id uint) string {
	key := fmt.Sprintf("%s%d", prefix, id)
	hashedKey := sha256.Sum256([]byte(key))

	return fmt.Sprintf("%s/%x", dir, hashedKey)
}

func CreateStorageUrl(key string) string {
	return getStorageEndpoint() + "/" + consts.AwsS3BucketName + "/" + key
}

func getStorageEndpoint() string {
	localStackEndpoint := os.Getenv(consts.EnvLocalStackEndpointForFront)

	if localStackEndpoint != "" {
		return localStackEndpoint
	} else {
		return "https://s3." + os.Getenv(consts.EnvAWSRegion) + ".amazonaws.com"
	}
}
