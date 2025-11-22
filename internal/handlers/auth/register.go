package auth

import (
	"context"

	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	converter "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/convertor"
)

func (i *Implementation) Register(
	ctx context.Context,
	req *proto.RegisterRequest,
) (*proto.RegisterResponse, error) {
	result, err := i.usecase.Register(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}

	return converter.RegisterResultToProto(result), nil
}

func (i *Implementation) VerifyTwoFactorSetup(
	ctx context.Context,
	req *proto.VerifyTwoFactorSetupRequest,
) (*proto.AuthResponse, error) {
	result, err := i.usecase.VerifyTwoFactorSetup(ctx, req.TempToken, req.Code)
	if err != nil {
		return nil, err
	}

	return converter.AuthResultToProto(result), nil
}
