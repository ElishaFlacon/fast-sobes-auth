package converter

import (
	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SettingsToProto(s *domain.Settings) *proto.Settings {
	if s == nil {
		return nil
	}
	return &proto.Settings{
		Id:                        s.Id,
		RequireTwoFactor:          s.RequireTwoFactor,
		TokenTtlMinutes:           s.TokenTTLMinutes,
		RefreshTokenTtlDays:       s.RefreshTokenTTLDays,
		MinPasswordLength:         s.MinPasswordLength,
		RequirePasswordComplexity: s.RequirePasswordComplexity,
		UpdatedAt:                 timestamppb.New(s.UpdatedAt),
	}
}

func ProtoToSettings(s *proto.Settings) *domain.Settings {
	if s == nil {
		return nil
	}
	return &domain.Settings{
		Id:                        s.Id,
		RequireTwoFactor:          s.RequireTwoFactor,
		TokenTTLMinutes:           s.TokenTtlMinutes,
		RefreshTokenTTLDays:       s.RefreshTokenTtlDays,
		MinPasswordLength:         s.MinPasswordLength,
		RequirePasswordComplexity: s.RequirePasswordComplexity,
		UpdatedAt:                 s.UpdatedAt.AsTime(),
	}
}
