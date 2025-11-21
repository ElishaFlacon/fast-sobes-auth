package settings

import "context"

func (i *Implementation) SettingsHandler(ctx context.Context) error {
	return i.usecase.SettingsUseCase()
}
