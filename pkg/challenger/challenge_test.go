package challenger_test

import (
	"crypto/sha512"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/kotkovdev/pow/pkg/challenger"
)

func TestPuzzleSolving(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		t.Parallel()

		t.Run("invalid hash provided", func(t *testing.T) {
			t.Parallel()

			chal := challenger.NewChallenger(challenger.DefaultSHA256Func, challenger.DefaultSaltGenerateFunc)
			msg, err := chal.CreatePuzzle([]byte("some request"), time.Now(), 2)
			assert.NoError(t, err)

			chal2 := challenger.NewChallenger(func(body []byte) []byte {
				hasher := sha512.New()
				hasher.Write(body)
				return hasher.Sum(nil)
			}, challenger.DefaultSaltGenerateFunc)

			result := chal2.SolveRecursive(msg.Source, msg.Target)

			assert.Nil(t, result)
		})

		t.Run("multiple cases", func(t *testing.T) {
			t.Parallel()

			for i := 0; i < 5; i++ {
				t.Run(fmt.Sprintf("test_case_%d", i), func(t *testing.T) {
					t.Parallel()

					chal := challenger.NewChallenger(challenger.DefaultSHA256Func, challenger.DefaultSaltGenerateFunc)
					msg, err := chal.CreatePuzzle([]byte("some request"), time.Now(), 2)
					assert.NoError(t, err)

					result := chal.SolveRecursive(msg.Source, msg.Target)

					assert.Equal(t, msg.Original, result)
				})
			}
		})
	})

	t.Run("fail", func(t *testing.T) {
		t.Parallel()
		t.Run("generate puzzle returns error", func(t *testing.T) {
			t.Parallel()

			expectedErr := errors.New("error generate puzzle")

			chal := challenger.NewChallenger(challenger.DefaultSHA256Func, func() ([]byte, error) {
				return nil, expectedErr
			})

			msg, err := chal.CreatePuzzle([]byte("some request"), time.Now(), 2)
			assert.Nil(t, msg)
			assert.ErrorIs(t, err, expectedErr)
		})
	})
}
