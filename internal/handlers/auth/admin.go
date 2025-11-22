package auth

import (
	"context"

	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	converter "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/convertor"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) UpdatePermissions(
	ctx context.Context,
	req *proto.UpdatePermissionsRequest,
) (*proto.User, error) {
	user, err := i.usecase.UpdatePermissions(ctx, req.UserId, req.PermissionLevel)
	if err != nil {
		return nil, err
	}
	return converter.UserToProto(user), nil
}

func (i *Implementation) DisableUser(
	ctx context.Context,
	req *proto.DisableUserRequest,
) (*emptypb.Empty, error) {
	err := i.usecase.DisableUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) EnableUser(
	ctx context.Context,
	req *proto.EnableUserRequest,
) (*emptypb.Empty, error) {
	err := i.usecase.EnableUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (i *Implementation) GetUser(
	ctx context.Context,
	req *proto.GetUserRequest,
) (*proto.User, error) {
	user, err := i.usecase.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return converter.UserToProto(user), nil
}

func (i *Implementation) ListUsers(
	ctx context.Context,
	req *proto.ListUsersRequest,
) (*proto.ListUsersResponse, error) {
	list, err := i.usecase.ListUsers(ctx, req.Page, req.PageSize, req.MinPermissionLevel, *req.IncludeDisabled)
	if err != nil {
		return nil, err
	}
	return converter.UserListToProto(list), nil
}

func (i *Implementation) DeleteUser(
	ctx context.Context,
	req *proto.DeleteUserRequest,
) (*emptypb.Empty, error) {
	err := i.usecase.DeleteUser(ctx, req.UserId)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
