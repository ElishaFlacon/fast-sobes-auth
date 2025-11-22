package converter

import (
	proto "buf.build/gen/go/fast-sobes/proto/protocolbuffers/go/auth"
	"github.com/ElishaFlacon/fast-sobes-auth/internal/domain"
)

func TwoFactorSetupToProto(tf *domain.TwoFactorSetup) *proto.EnableTwoFactorResponse {
	if tf == nil {
		return nil
	}
	return &proto.EnableTwoFactorResponse{
		Secret:      tf.Secret,
		QrCodeUrl:   tf.QRCodeURL,
		BackupCodes: tf.BackupCodes,
	}
}
