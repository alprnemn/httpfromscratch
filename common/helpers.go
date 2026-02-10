package common

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
)

// IsAllUpperBytes reports whether all characters in the given byte slice
// are uppercase ASCII letters (A–Z).
func IsAllUpperBytes(b []byte) bool {
	for _, c := range b {
		if c < 'A' || c > 'Z' {
			return false
		}
	}
	return true
}

// GenerateRandomBytes returns a slice of cryptographically secure random bytes
// of the specified size. It panics if the system random number generator fails.
func GenerateRandomBytes(size int) []byte {
	b := make([]byte, size)

	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Sprintf("Failed to generate random bytes: %v", err))
	}

	return b
}

// FailOnError logs the given message and error using log.Fatalf if err is not nil.
// Intended for use in fatal application initialization paths.
func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// ConvertDecToHex converts a non-negative decimal integer into its hexadecimal
// string representation. It returns an error for negative inputs.
func ConvertDecToHex(numDec int) (string, error) {

	if numDec < 0 {
		return "", fmt.Errorf("negative numbers not supported")
	}

	if numDec == 0 {
		return "0", nil
	}

	quotient := numDec
	remainder := 0
	res := ""

	for quotient != 0 {
		remainder = quotient % 16
		quotient /= 16
		if remainder < 10 {
			res += string(rune('0' + remainder))
		} else {
			res += string(rune('A' + (remainder - 10)))
		}
	}

	return reverseString(res)
}

// reverseString reverses the given string and returns the result.
// It returns an error if the input string is empty.
func reverseString(val string) (string, error) {
	if len(val) == 0 {
		return "", errors.New("reverse string: empty input")
	}

	b := make([]byte, len(val))
	for i := 0; i < len(val); i++ {
		b[i] = val[len(val)-1-i]
	}
	return string(b), nil
}
