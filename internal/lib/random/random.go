package random

import (
	"math/rand"
	"time"
)

const ALPHABET = "ABCDEFGHIJKMNLOPQRSTUVWXYZabcdefghijkmnlopqrstuvwxyz0123456789"

// actually math/rand is not safety because of it generates fake-random numbers
// it means our alias can be hacked, but it is just link shortener, what can be wrong?
func NewRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune(ALPHABET)

	b := make([]rune, size)
	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))]
	}
	return string(b)
}
