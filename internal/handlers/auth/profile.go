package auth

import "context"

func (i *Implementation) Profile(ctx context.Context) error {
    return i.usecase.Profile()
}
