package common

import (
	"crypto/rand"
	"fmt"
	"log"
)

// IsAllUpper function checks string to all chars are upper
func IsAllUpperBytes(b []byte) bool {
	for _, c := range b {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

// GenerateRandomBytes func generates bytes with given size
func GenerateRandomBytes(size int) []byte {
	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate random bytes: %v", err))
	}

	return b
}

// FailOnError fatal error check
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
