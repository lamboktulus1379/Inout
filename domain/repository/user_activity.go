package repository

import (
	"context"

	"inout/domain/model"
)

//go:generate mockgen -destination ../../mocks/repositories/mock_iuser_activity_repository.go -package=mocks inout/domain/repository IUserActivity
type IUserActivity interface {
	GetById(ctx context.Context, id uint64) (model.UserActivity, error)
	GetByUser(ctx context.Context, userID int) (model.UserActivity, error)
	Create(ctx context.Context, userActivity model.UserActivity) (model.UserActivity, error)
	Get(ctx context.Context, reqQueryParam model.ReqQueryParamUserActivity) ([]model.UserActivity, int64, error)
	Update(ctx context.Context, userActivity model.UserActivity, ID uint64) error
	Delete(ctx context.Context, ID uint64) error
	GetLastLogin(ctx context.Context, UserID uint64) (model.UserActivity, error)
	GetLastLogout(ctx context.Context, UserID uint64) (model.UserActivity, error)
}
