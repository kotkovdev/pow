package server_test

import (
	"bufio"
	"errors"
	"net"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"

	"github.com/kotkovdev/pow/internal/server"
	"github.com/kotkovdev/pow/internal/server/mocks"
	"github.com/kotkovdev/pow/internal/util"
	"github.com/kotkovdev/pow/pkg/challenger"
)

type serverSuite struct {
	suite.Suite
	quoteService *mocks.MockQuotesService
	challenger   *mocks.MockChallenger
	client       net.Conn
	server       net.Conn
	srv          server.Server
}

func (s *serverSuite) SetupSuite() {
	ctrl := gomock.NewController(s.T())
	s.quoteService = mocks.NewMockQuotesService(ctrl)
	s.challenger = mocks.NewMockChallenger(ctrl)
	s.srv = server.New(s.quoteService, s.challenger, 1)
}

func (s *serverSuite) SetupTest() {
	s.client, s.server = net.Pipe()
}

func (s *serverSuite) TearDownTest() {
	s.server.Close()
	s.client.Close()
}

func (s *serverSuite) TestServerPipeline() {
	s.challenger.EXPECT().CreatePuzzle(gomock.Any(), gomock.Any(), 1).Return(&challenger.Puzzle{
		Original: []byte{10, 20, 30},
		Source:   []byte{10, 20},
		Target:   []byte{50, 60, 70},
	}, nil)
	s.quoteService.EXPECT().GetRandomQuote().Return("quote", nil)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		s.server.Write([]byte(""))
		wg.Done()
	}()

	var result string
	var err error
	go func() {
		result, err = bufio.NewReader(s.client).ReadString(util.MessageDelimeter)
		wg.Done()
	}()

	s.srv.HandleConnection(s.server)
	wg.Wait()

	s.NoError(err)
	s.Equal("ChQ=|MjxG\n", result)

	wg.Add(1)
	go func() {
		result, err = bufio.NewReader(s.client).ReadString(util.MessageDelimeter)
		wg.Done()
	}()

	s.srv.HandleSolution([]byte("ChQe"), s.server)
	wg.Wait()

	s.NoError(err)
	s.Equal("quote\n", result)
}

func (s *serverSuite) TestServerPipelineCreatePuzzleError() {
	expectedErr := errors.New("generate error")
	s.challenger.EXPECT().CreatePuzzle(gomock.Any(), gomock.Any(), 1).Return(&challenger.Puzzle{}, expectedErr)

	s.srv.HandleConnection(s.server)
}

func (s *serverSuite) TestServerPipelineQuoteError() {
	quoteErr := errors.New("quote error")
	s.challenger.EXPECT().CreatePuzzle(gomock.Any(), gomock.Any(), 1).Return(&challenger.Puzzle{
		Original: []byte{10, 20, 30},
		Source:   []byte{10, 20},
		Target:   []byte{50, 60, 70},
	}, nil)
	s.quoteService.EXPECT().GetRandomQuote().Return("", quoteErr)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		s.server.Write([]byte(""))
		wg.Done()
	}()

	var result string
	var err error
	go func() {
		result, err = bufio.NewReader(s.client).ReadString(util.MessageDelimeter)
		wg.Done()
	}()

	s.srv.HandleConnection(s.server)
	wg.Wait()

	s.NoError(err)
	s.Equal("ChQ=|MjxG\n", result)

	wg.Add(1)
	go func() {
		result, err = bufio.NewReader(s.client).ReadString(util.MessageDelimeter)
		wg.Done()
	}()

	s.srv.HandleSolution([]byte("ChQe"), s.server)
	wg.Wait()

	s.NoError(err)
	s.Empty(result)
}

func TestServer(t *testing.T) {
	suite.Run(t, new(serverSuite))
}
