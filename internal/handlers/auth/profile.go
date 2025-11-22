package auth

import (
	"context"

	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) ChangePassword(
	ctx context.Context,
	req *proto.ChangePasswordRequest,
) (*emptypb.Empty, error) {
	err := i.usecase.ChangePassword(ctx, req.UserId, req.OldPassword, req.NewPassword)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) ChangeEmail(
	ctx context.Context,
	req *proto.ChangeEmailRequest,
) (*emptypb.Empty, error) {
	err := i.usecase.ChangeEmail(ctx, req.UserId, req.NewEmail, req.Password)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) VerifyEmailChange(
	ctx context.Context,
	req *proto.VerifyEmailChangeRequest,
) (*emptypb.Empty, error) {
	err := i.usecase.VerifyEmailChange(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
