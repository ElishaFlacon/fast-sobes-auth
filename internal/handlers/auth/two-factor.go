package auth

import (
	"context"

	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	converter "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/convertor"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) EnableTwoFactor(
	ctx context.Context,
	req *proto.EnableTwoFactorRequest,
) (*proto.EnableTwoFactorResponse, error) {
	result, err := i.usecase.EnableTwoFactor(ctx, req.UserId, req.Password)
	if err != nil {
		return nil, err
	}
	return converter.TwoFactorSetupToProto(result), nil
}

func (i *Implementation) DisableTwoFactor(
	ctx context.Context,
	req *proto.DisableTwoFactorRequest,
) (*emptypb.Empty, error) {
	err := i.usecase.DisableTwoFactor(ctx, req.UserId, req.Password, req.Code)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
