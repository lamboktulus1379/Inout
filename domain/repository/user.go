package repository

import (
	"context"

	"inout/domain/model"
)

//go:generate mockgen -destination ../../mocks/repositories/mock_iuser_repository.go -package=mocks inout/domain/repository IUser
type IUser interface {
	GetById(ctx context.Context, id int) (user model.User, err error)
	GetByUserName(ctx context.Context, userName string) (model.User, error)
}
