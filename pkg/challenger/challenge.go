// package challenger implements challenge response.
package challenger

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"time"
)

const (
	// maxComplexity is a max depth for solving puzzle recursion.
	maxComplexity = 3

	// alphabet is a alphabet for generating hash salt.
	alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
)

// Puzzle is a generated puzzle.
type Puzzle struct {
	Original []byte
	Target   []byte
	Source   []byte
}

// HashFunc is a function for hashing provided bytes.
type HashFunc func(body []byte) []byte

// DefaultSHA256Func returns sha256 hashed bytes.
func DefaultSHA256Func(body []byte) []byte {
	hasher := sha256.New()
	hasher.Write(body)
	return hasher.Sum(nil)
}

// Challenger is a challenge instance.
type Challenger struct {
	hashFn HashFunc
}

// NewChallenger returns challenger instance.
func NewChallenger(hashFunc HashFunc) Challenger {
	return Challenger{
		hashFn: hashFunc,
	}
}

// CreatePuzzle creates puzzle challenge for client.
func (c Challenger) CreatePuzzle(req []byte, timestamp time.Time, size int) (*Puzzle, error) {
	payload := append(req, []byte(timestamp.String())...)
	salt := make([]byte, 10)
	for idx := range salt {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet)-1)))
		if err != nil {
			return nil, err
		}
		salt[idx] = byte(alphabet[num.Int64()])
	}

	original := c.hashFn(append(payload, salt...))
	originalCopy := make([]byte, len(original))
	copy(originalCopy, original)
	source := originalCopy[:len(originalCopy)-size]
	target := c.hashFn(original)
	msg := &Puzzle{
		Source:   source,
		Target:   target,
		Original: original,
	}
	return msg, nil
}

// SolveRecursive calculates source hash during it not equal target and complexity less than max coplexity.
func (c Challenger) SolveRecursive(source, target []byte) (result []byte) {
	var check func(source []byte, current, depth int)
	check = func(source []byte, current, depth int) {
	loop:
		for i := 0; i <= 255; i++ {
			generatedHash := append(source, byte(i))
			calculatedHash := c.hashFn(generatedHash)
			if bytes.Equal(calculatedHash, target) {
				generatedHashCopy := make([]byte, len(generatedHash))
				copy(generatedHashCopy, generatedHash)
				result = generatedHashCopy
				break loop
			}
			if current < depth {
				check(append(source, byte(i)), current+1, depth)
			}
		}
	}

	check(source, 1, maxComplexity)
	return
}
