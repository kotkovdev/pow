package challenger_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kotkovdev/pow/pkg/challenger"
)

func TestPuzzleSolving(t *testing.T) {
	chal := challenger.NewChallenger(challenger.DefaultSHA256Func)
	msg, err := chal.CreatePuzzle([]byte("some request"), time.Now(), 2)
	assert.NoError(t, err)

	result := chal.SolveRecursive(msg.Source, msg.Target)

	assert.Equal(t, msg.Original, result)
}
