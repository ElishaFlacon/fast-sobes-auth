package auth

import "context"

func (i *Implementation) Login(ctx context.Context) error {
    return i.usecase.Login()
}
