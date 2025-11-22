package hello

import (
	"context"

	"buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/test"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Hello(ctx context.Context, _ *emptypb.Empty) (*test.HelloResponse, error) {
	response := i.usecase.Hello()
	return &test.HelloResponse{Message: response}, nil
}
