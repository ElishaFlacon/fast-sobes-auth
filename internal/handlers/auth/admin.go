package auth

import "context"

func (i *Implementation) Admin(ctx context.Context) error {
    return i.usecase.Admin()
}
