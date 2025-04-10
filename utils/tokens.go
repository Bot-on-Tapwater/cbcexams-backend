package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// GenerateRandomToken generates a random token of the specified length.
// It creates a byte slice of the given length, fills it with random data
// using the crypto/rand package, and then encodes it as a hexadecimal string.
// If an error occurs during random data generation, the function will panic.
//
// Parameters:
//   - length: The desired length of the random token in bytes.
//
// Returns:
//   A hexadecimal string representation of the random token.
func GenerateRandomToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}