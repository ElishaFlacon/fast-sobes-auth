package auth

import "context"

func (i *Implementation) Twofactor(ctx context.Context) error {
    return i.usecase.Twofactor()
}
