package challenger

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"
)

const (
	size          = 3
	maxComplexity = 3
)

type Request struct{}

type Message struct {
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

func NewChallenger(hashFunc HashFunc) Challenger {
	return Challenger{
		hashFn: hashFunc,
	}
}

func (c *Challenger) CreatePuzzle(req []byte, timestamp time.Time) Message {
	original := c.hashFn(append(req, []byte(timestamp.String())...))
	source := original[:len(original)-size]
	target := c.hashFn(original)
	msg := Message{
		Source:   source,
		Target:   target,
		Original: original,
	}
	return msg
}

// SolveRecursive calculates source hash during it not equal target and complexity less than max coplexity.
func (c Challenger) SolveRecursive(source, target []byte) []byte {
	var check func(source []byte, current, deep int) []byte
	check = func(source []byte, current, deep int) []byte {
		for i := 0; i <= 255; i++ {
			generatedHash := append(source, byte(i))
			calculatedHash := c.hashFn(generatedHash)
			if bytes.Equal(calculatedHash, target) {
				fmt.Printf("found hash: %s", hex.EncodeToString(generatedHash))
				return generatedHash
			}
		}

		for i := 0; i <= 255; i++ {
			if current < deep {
				generatedHash := append(source, byte(i))
				check(generatedHash, current+1, deep)
			}
		}

		return []byte{}
	}

	return check(source, 1, maxComplexity)
}
