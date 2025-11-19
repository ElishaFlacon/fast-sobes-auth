package hello

import (
	"context"

	pb "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/hello/proto"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Implementation struct {
	pb.UnimplementedHelloServer
	usecase usecase.HelloUsecase
}

func NewImplementation(uc usecase.HelloUsecase) *Implementation {
	return &Implementation{
		usecase: uc,
	}
}

func (i *Implementation) RegisterImplementation(grpcServer *grpc.Server) {
	pb.RegisterHelloServer(grpcServer, i)
}

func (i *Implementation) Hello(ctx context.Context, _ *emptypb.Empty) (*pb.HelloResponse, error) {
	response := i.usecase.Hello()
	return &pb.HelloResponse{Message: response}, nil
}
