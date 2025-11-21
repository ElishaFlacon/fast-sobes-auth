package auth

import "context"

func (i *Implementation) Tokens(ctx context.Context) error {
    return i.usecase.Tokens()
}
