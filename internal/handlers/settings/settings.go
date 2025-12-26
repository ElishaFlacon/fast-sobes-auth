package settings

import (
	"context"

	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	converter "github.com/ElishaFlacon/fast-sobes-auth/internal/handlers/convertor"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) GetSettings(ctx context.Context, _ *emptypb.Empty) (*proto.Settings, error) {
	settings, err := i.usecase.GetSettings(ctx)
	if err != nil {
		return nil, err
	}
	return converter.SettingsToProto(settings), nil
}

func (i *Implementation) UpdateSettings(ctx context.Context, req *proto.UpdateSettingsRequest) (*proto.Settings, error) {
	updated, err := i.usecase.UpdateSettings(
		ctx,
		req.RequireTwoFactor,
		req.TokenTtlMinutes,
		req.MinPasswordLength,
		req.RequirePasswordComplexity,
	)
	if err != nil {
		return nil, err
	}
	return converter.SettingsToProto(updated), nil
}

func (i *Implementation) ResetSettings(ctx context.Context, _ *emptypb.Empty) (*proto.Settings, error) {
	res, err := i.usecase.ResetSettings(ctx)
	if err != nil {
		return nil, err
	}
	return converter.SettingsToProto(res), nil
}
