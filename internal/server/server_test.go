package server_test

import (
	"bufio"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/kotkovdev/pow/internal/netutil"
	"github.com/kotkovdev/pow/internal/server"
	"github.com/kotkovdev/pow/internal/server/mocks"
	"github.com/kotkovdev/pow/pkg/challenger"
)

type serverSuite struct {
	suite.Suite
	quoteService *mocks.MockQuotesService
	challenger   *mocks.MockChallenger
	client       net.Conn
	srv          server.Server
}

func (s *serverSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.quoteService = mocks.NewMockQuotesService(ctrl)
	s.challenger = mocks.NewMockChallenger(ctrl)
	s.srv = server.New(s.quoteService, s.challenger, 1)
}

func (s *serverSuite) SetupTest() {
	go func() {
		err := s.srv.Serve(":8081")
		s.Require().NoError(err)
	}()

	// await to start server.
	time.Sleep(time.Second)
	var err error
	s.client, err = net.Dial("tcp", ":8081")
	s.NoError(err)
}

func (s *serverSuite) TearDownTest() {
	s.client.Close()
}

func (s *serverSuite) TestServerSuccess() {
	s.challenger.EXPECT().CreatePuzzle(gomock.Any(), gomock.Any(), 1).Return(&challenger.Puzzle{
		Original: []byte{10, 20, 30},
		Source:   []byte{10, 20},
		Target:   []byte{50, 60, 70},
	}, nil)
	s.quoteService.EXPECT().GetRandomQuote().Return("quote", nil)

	_, err := s.client.Write([]byte("\n"))
	s.NoError(err)

	result, err := bufio.NewReader(s.client).ReadString(netutil.MessageDelimeter)

	s.NoError(err)
	s.Equal("ChQ=|MjxG\n", result)

	_, err = s.client.Write([]byte("ChQe\n"))
	s.NoError(err)

	result, err = bufio.NewReader(s.client).ReadString(netutil.MessageDelimeter)

	s.NoError(err)
	s.Equal("quote\n", result)
}

func TestServer(t *testing.T) {
	suite.Run(t, new(serverSuite))
}
