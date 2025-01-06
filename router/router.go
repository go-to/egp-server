package router

import (
	"context"
	"fmt"
	"github.com/go-to/egp-protobuf/pb"
	"github.com/go-to/egp_backend/usecase"
	"github.com/go-to/egp_backend/usecase/input"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
)

type Server struct {
	pb.UnimplementedEgpServiceServer
	Usecase Usecase
}

type Usecase struct {
	Shop usecase.ShopUsecase
}

func New(port int, u Usecase) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()

	pb.RegisterEgpServiceServer(s, NewServer(u))

	reflection.Register(s)

	go func() {
		log.Printf("grpc server listening on port %d", port)
		err := s.Serve(listener)
		if err != nil {
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutting down server...")
	s.GracefulStop()
}

func NewServer(u Usecase) *Server {
	return &Server{Usecase: u}
}

func (s *Server) GetShops(ctx context.Context, req *pb.ShopsRequest) (*pb.ShopsResponse, error) {
	in := input.ShopsInput{
		ShopsRequest: req,
	}

	out, err := s.Usecase.Shop.GetShops(&in)
	if err != nil {
		return nil, err
	}

	return &out.ShopsResponse, nil
}
