package persistence

import (
	"context"
	"log"

	"inout/domain/model"
	"inout/domain/repository"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(DB *gorm.DB) repository.IUser {
	return &UserRepository{DB}
}

func (repo *UserRepository) GetById(ctx context.Context, id int) (user model.User, err error) {
	var result model.User
	if err = repo.DB.Debug().WithContext(ctx).First(&user, id).Error; err != nil {
		log.Printf("Failed Get User With Error : %v", err)
		return result, err
	}

	return result, nil
}

func (repo *UserRepository) GetByUserName(ctx context.Context, userName string) (model.User, error) {
	var user model.User
	if err := repo.DB.Debug().WithContext(ctx).Where("user_name = ?", userName).First(&user).Error; err != nil {
		log.Printf("Failed Get User With Error : %v", err)
		return user, err
	}

	return user, nil
}
