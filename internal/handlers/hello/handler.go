package hello

import (
	"context"

	"buf.build/gen/go/fast-sobes/proto/grpc/go/test/testgrpc"
	"buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/test"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Implementation struct {
	testgrpc.UnimplementedHelloServer
	usecase usecase.HelloUsecase
}

func NewImplementation(uc usecase.HelloUsecase) *Implementation {
	return &Implementation{
		usecase: uc,
	}
}

func (i *Implementation) RegisterImplementation(grpcServer *grpc.Server) {
	testgrpc.RegisterHelloServer(grpcServer, i)
}

func (i *Implementation) Hello(ctx context.Context, _ *emptypb.Empty) (*test.HelloResponse, error) {
	response := i.usecase.Hello()
	return &test.HelloResponse{Message: response}, nil
}
