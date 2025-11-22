package converter

import (
	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func UserToProto(u *domain.User) *proto.User {
	if u == nil {
		return nil
	}
	return &proto.User{
		Id:               u.ID,
		Email:            u.Email,
		PermissionLevel:  u.PermissionLevel,
		Disabled:         u.Disabled,
		TwoFactorEnabled: u.TwoFactorEnabled,
		CreatedAt:        timestamppb.New(u.CreatedAt),
		UpdatedAt:        timestamppb.New(u.UpdatedAt),
	}
}

func ProtoToUser(u *proto.User) *domain.User {
	if u == nil {
		return nil
	}
	return &domain.User{
		ID:               u.Id,
		Email:            u.Email,
		PermissionLevel:  u.PermissionLevel,
		Disabled:         u.Disabled,
		TwoFactorEnabled: u.TwoFactorEnabled,
		CreatedAt:        u.CreatedAt.AsTime(),
		UpdatedAt:        u.UpdatedAt.AsTime(),
	}
}
