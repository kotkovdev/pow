package challenger_test

import (
	"testing"
	"time"

	"github.com/kotkovdev/pow/pkg/challenger"
	"github.com/stretchr/testify/assert"
)

func TestHash(t *testing.T) {
	chal := challenger.NewChallenger(challenger.DefaultSHA256Func)
	msg := chal.CreatePuzzle([]byte("some request"), time.Now())

	result := chal.SolveRecursive(msg.Source, msg.Target)

	assert.Equal(t, msg.Original, result)
}
