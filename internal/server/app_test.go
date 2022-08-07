package server

import (
	"analytic-service/config"
	"analytic-service/internal/logger"
	"context"
	"log"
	"net"
	"os"
	"testing"

	pb "github.com/castaneai/grpc-testing-with-bufconn"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

type server struct{}

func (*server) SayHello(context.Context, *pb.HelloRequest) (*pb.HelloReply, error) {
	panic("unimplemented")
}

var lis *bufconn.Listener

const bufSize = 1024 * 1024

func bufDialer(ctx context.Context, address string) (net.Conn, error) {
	return lis.Dial()
}

func TestApplication_Mock_New(t *testing.T) {
	var a *Application

	t.Run("App_NEW_with_Mocks", func(t *testing.T) {
		os.Setenv("DB_CONNECTION_STRING", "test:test")
		os.Setenv("DB", "MEM")
		os.Setenv("AUTHSERVICE_ADDRESS", "bufnet")

		lis = bufconn.Listen(bufSize)
		s := grpc.NewServer()
		pb.RegisterGreeterServer(s, &server{})
		go func() {
			if err := s.Serve(lis); err != nil {
				log.Fatal(err)
			}
		}()

		cfgMock, _ := config.Init()
		logger, _ := logger.New(cfgMock)
		logger.Info("LOGGER GET SUCCESS")
		a = New(cfgMock, logger, grpc.WithContextDialer(bufDialer))
		if a == nil {
			t.Error("App not created")
		}
	})

}

func TestApplication_New(t *testing.T) {
	var a *Application

	t.Run("App_NEW_with_MEM", func(t *testing.T) {
		os.Setenv("DB_CONNECTION_STRING", "test:test")
		os.Setenv("DB", "MEM")

		cfg, _ := config.Init()
		logger, _ := logger.New(cfg)
		logger.Info("LOGGER GET SUCCESS")
		a = New(cfg, logger, grpc.WithBlock())
		if a == nil {
			t.Error("App not created")
		}
	})

	t.Run("App_RUN_with_MEM", func(t *testing.T) {
		a.Run()
		a.Dispose()
	})

}
