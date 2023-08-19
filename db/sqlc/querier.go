// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package db

import (
	"context"
)

type Querier interface {
	CreateImage(ctx context.Context, arg CreateImageParams) (Image, error)
	CreateUser(ctx context.Context, arg CreateUserParams) (User, error)
	GetImages(ctx context.Context, userID int64) ([]Image, error)
	GetUser(ctx context.Context, username string) (User, error)
}

var _ Querier = (*Queries)(nil)
