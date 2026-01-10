package common

import (
	"crypto/rand"
	"fmt"
	"log"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func GenerateRandomBytes(size int) []byte {
	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate random bytes: %v", err))
	}

	return b
}
