package converter

import (
	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func RegisterResultToProto(r *domain.RegisterResult) *proto.RegisterResponse {
	if r == nil {
		return nil
	}
	return &proto.RegisterResponse{
		User:              UserToProto(r.User),
		TwoFactorSecret:   r.TwoFactorSecret,
		QrCodeUrl:         r.QRCodeURL,
		RequiresTwoFactor: r.RequiresTwoFactor,
		TempToken:         r.TempToken,
	}
}

func AuthResultToProto(a *domain.AuthResult) *proto.AuthResponse {
	if a == nil {
		return nil
	}
	return &proto.AuthResponse{
		AccessToken:  a.AccessToken,
		RefreshToken: a.RefreshToken,
		ExpiresIn:    a.ExpiresIn,
		User:         UserToProto(a.User),
	}
}

func LoginResultToProto(l *domain.LoginResult) *proto.LoginResponse {
	if l == nil {
		return nil
	}
	var auth *proto.AuthResponse
	if l.Auth != nil {
		auth = AuthResultToProto(l.Auth)
	}
	return &proto.LoginResponse{
		RequiresTwoFactor: l.RequiresTwoFactor,
		TempToken:         l.TempToken,
		Auth:              auth,
	}
}
