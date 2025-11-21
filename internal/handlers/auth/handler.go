package auth

import (
	"github.com/ElishaFlacon/fast-sobes-auth/internal/usecase"

	"google.golang.org/grpc"
)

type Implementation struct {
	usecase usecase.AuthUsecase
}

func NewImplementation(uc usecase.AuthUsecase) *Implementation {
	return &Implementation{
		usecase: uc,
	}
}

func (i *Implementation) RegisterImplementation(grpcServer *grpc.Server) {
	// Register your gRPC service here
}
