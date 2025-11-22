package converter

import (
	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func TokenInfoToProto(t *domain.TokenInfo) *proto.VerifyTokenResponse {
	if t == nil {
		return nil
	}
	return &proto.VerifyTokenResponse{
		Valid:           t.Valid,
		UserId:          t.UserID,
		PermissionLevel: t.PermissionLevel,
		Email:           t.Email,
	}
}

func ProtoToTokenInfo(p *proto.VerifyTokenResponse) *domain.TokenInfo {
	if p == nil {
		return nil
	}
	return &domain.TokenInfo{
		Valid:           p.Valid,
		UserID:          p.UserId,
		PermissionLevel: p.PermissionLevel,
		Email:           p.Email,
	}
}
