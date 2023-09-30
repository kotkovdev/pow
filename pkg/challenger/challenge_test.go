package challenger_test

import (
	"encoding/hex"
	"fmt"
	"testing"
	"time"

	"github.com/kotkovdev/pow/pkg/challenger"
)

func TestHash(t *testing.T) {
	chal := challenger.NewChallenger(challenger.DefaultSHA256Func)
	msg := chal.CreatePuzzle([]byte("some request"), time.Now())

	fmt.Printf("source hash: %s\n", hex.EncodeToString(msg.Source))
	fmt.Printf("target hash: %s\n", hex.EncodeToString(msg.Target))
	fmt.Printf("original hash: %s\n", hex.EncodeToString(msg.Original))
	fmt.Println("===========================")
	result := chal.SolveRecursive(msg.Source, msg.Target)
	fmt.Printf("found result: %s\n", hex.EncodeToString(result))
}
