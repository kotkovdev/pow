package challenger

import (
	"bytes"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"time"
)

const (
	maxComplexity = 3
	alphabet      = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-"
)

type Puzzle struct {
	Original []byte
	Target   []byte
	Source   []byte
}

type HashFunc func(body []byte) []byte

func DefaultSHA256Func(body []byte) []byte {
	hasher := sha256.New()
	hasher.Write(body)
	return hasher.Sum(nil)
}

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
func (c *Challenger) CreatePuzzle(req []byte, timestamp time.Time, size int) (*Puzzle, error) {
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
	source := original[:len(original)-size]
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
	var check func(source []byte, current, deep int)
	check = func(source []byte, current, deep int) {
		for i := 0; i <= 255; i++ {
			generatedHash := append(source, byte(i))
			calculatedHash := c.hashFn(generatedHash)
			if bytes.Equal(calculatedHash, target) {
				result = generatedHash
			}
		}
		for i := 0; i <= 255; i++ {
			if current < deep {
				generatedHash := append(source, byte(i))
				check(generatedHash, current+1, deep)
			}
		}
	}

	check(source, 1, maxComplexity)
	return
}
