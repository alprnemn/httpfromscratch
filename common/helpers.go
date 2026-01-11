package common

import (
	"crypto/rand"
	"fmt"
	"log"
	"unicode"
)

// IsAllUpper function checks string to all chars are upper
func IsAllUpper(s string) bool {
	for _, r := range s {
		if unicode.IsLetter(r) && !unicode.IsUpper(r) {
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
