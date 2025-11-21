package auth

import "context"

func (i *Implementation) Register(ctx context.Context) error {
    return i.usecase.Register()
}
