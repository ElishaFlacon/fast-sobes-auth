package hello

import (
	"buf.build/gen/go/fast-sobes/proto/grpc/go/test/testgrpc"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"

	"google.golang.org/grpc"
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
