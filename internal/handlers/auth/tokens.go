package auth

import (
	"context"

	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	converter "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/convertor"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) VerifyToken(
	ctx context.Context,
	req *proto.VerifyTokenRequest,
) (*proto.VerifyTokenResponse, error) {
	result, err := i.usecase.VerifyToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return converter.TokenInfoToProto(result), nil
}

func (i *Implementation) RefreshToken(
	ctx context.Context,
	req *proto.RefreshTokenRequest,
) (*proto.AuthResponse, error) {
	result, err := i.usecase.RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return converter.AuthResultToProto(result), nil
}

func (i *Implementation) Logout(
	ctx context.Context,
	req *proto.LogoutRequest,
) (*emptypb.Empty, error) {
	err := i.usecase.Logout(ctx, req.Token, req.RefreshToken)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
