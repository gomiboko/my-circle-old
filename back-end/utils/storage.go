package utils

import (
	"crypto/sha256"
	"fmt"
)

func CreateHashedStorageKey(dir string, prefix string, id uint) string {
	key := fmt.Sprintf("%s%d", prefix, id)
	hashedKey := sha256.Sum256([]byte(key))

	return fmt.Sprintf("%s/%x", dir, hashedKey)
}
