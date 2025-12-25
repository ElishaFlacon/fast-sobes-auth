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
		UserId:          t.UserId,
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
		UserId:          p.UserId,
		PermissionLevel: p.PermissionLevel,
		Email:           p.Email,
	}
}
