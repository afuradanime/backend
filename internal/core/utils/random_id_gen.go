package utils

import (
	"crypto/rand"
	"math/big"
	"time"
)

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// generates a random alphanumeric ID, the first 6 chars encode a millisecond timestamp (for rough sortability),
// the remaining chars are cryptographically random
// length: 11 characters
func GenerateRandomID() string {
	const timeChars = 6
	const randChars = 5
	id := make([]byte, timeChars+randChars)

	// Encode timestamp in base-62 (first 6 chars)
	ts := uint64(time.Now().UnixMilli())
	for i := timeChars - 1; i >= 0; i-- {
		id[i] = alphabet[ts%62]
		ts /= 62
	}

	// Fill remaining chars with crypto/rand
	for i := timeChars; i < timeChars+randChars; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		id[i] = alphabet[n.Int64()]
	}

	return string(id)
}
