// package quotes provides random phrases.
package quotes

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"os"
)

type Service struct {
	quotes []string
}

// New returns new service.
func New(path string) (*Service, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	rows := bytes.Split(content, []byte("\n"))
	quotes := make([]string, len(rows))
	for idx, row := range rows {
		quotes[idx] = string(row)
	}

	return &Service{
		quotes: quotes,
	}, nil
}

// GetRandomQuote returns random quote.
func (s *Service) GetRandomQuote() (string, error) {
	num, err := rand.Int(rand.Reader, big.NewInt(int64(len(s.quotes)-1)))
	if err != nil {
		return "", err
	}
	return s.quotes[num.Int64()], nil
}
