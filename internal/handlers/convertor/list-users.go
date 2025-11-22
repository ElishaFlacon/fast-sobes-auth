package converter

import (
	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func UserListToProto(l *domain.UserList) *proto.ListUsersResponse {
	if l == nil {
		return nil
	}
	users := make([]*proto.User, len(l.Users))
	for i, u := range l.Users {
		users[i] = UserToProto(u)
	}
	return &proto.ListUsersResponse{
		Users:    users,
		Total:    l.Total,
		Page:     l.Page,
		PageSize: l.PageSize,
	}
}
