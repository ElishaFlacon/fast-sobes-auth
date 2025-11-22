package auth

import (
	"buf.build/gen/go/fast-sobes/proto/grpc/go/auth/authgrpc"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"

	"google.golang.org/grpc"
)

type Implementation struct {
	authgrpc.UnimplementedAuthServiceServer
	usecase usecase.AuthUsecase
}

func NewImplementation(uc usecase.AuthUsecase) *Implementation {
	return &Implementation{
		usecase: uc,
	}
}

func (i *Implementation) RegisterImplementation(grpcServer *grpc.Server) {
	authgrpc.RegisterAuthServiceServer(grpcServer, i)
}
