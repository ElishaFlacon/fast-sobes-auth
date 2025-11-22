package hello

import (
	"context"

	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/test"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Hello(ctx context.Context, _ *emptypb.Empty) (*proto.HelloResponse, error) {
	response := i.usecase.Hello()
	return &proto.HelloResponse{Message: response}, nil
}
