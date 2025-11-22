package auth

import (
	"context"

	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	converter "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/convertor"
)

func (i *Implementation) Login(
	ctx context.Context,
	req *proto.LoginRequest,
) (*proto.LoginResponse, error) {
	result, err := i.usecase.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, err
	}
	return converter.LoginResultToProto(result), nil
}

func (i *Implementation) VerifyTwoFactor(
	ctx context.Context,
	req *proto.VerifyTwoFactorRequest,
) (*proto.AuthResponse, error) {
	result, err := i.usecase.VerifyTwoFactor(ctx, req.TempToken, req.Code)
	if err != nil {
		return nil, err
	}
	return converter.AuthResultToProto(result), nil
}
