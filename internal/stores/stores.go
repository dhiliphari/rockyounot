package stores

import "context"

type PasswordStore interface {
	Add(ctx context.Context, passwords []any) error
	Exists(ctx context.Context, password []byte) (bool, error)
}
