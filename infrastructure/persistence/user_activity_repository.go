package persistence

import (
	"context"
	"fmt"
	"log"

	"inout/constants"
	"inout/domain/model"
	"inout/domain/repository"

	"gorm.io/gorm"
)

type UserActivityRepository struct {
	DB *gorm.DB
}

func NewUserActivityRepository(DB *gorm.DB) repository.IUserActivity {
	return &UserActivityRepository{DB}
}

func (repo *UserActivityRepository) GetById(ctx context.Context, id uint64) (model.UserActivity, error) {
	var userActivity model.UserActivity
	if err := repo.DB.Debug().WithContext(ctx).First(&userActivity, id).Error; err != nil {
		log.Printf("Failed Get User Activity With Error : %v", err)
		return userActivity, err
	}

	return userActivity, nil
}

func (repo *UserActivityRepository) GetByUser(ctx context.Context, userID int) (model.UserActivity, error) {
	var userActivity model.UserActivity
	if err := repo.DB.Debug().WithContext(ctx).Where("user_id = ?", userID).First(&userActivity).Error; err != nil {
		log.Printf("Failed Get User With Error : %v", err)
		return userActivity, err
	}

	return userActivity, nil
}

func (repo *UserActivityRepository) Create(ctx context.Context, userActivity model.UserActivity) (model.UserActivity, error) {
	if err := repo.DB.Debug().WithContext(ctx).Create(&userActivity).Error; err != nil {
		log.Printf("Failed Get User With Error : %v", err)
		return userActivity, err
	}

	return userActivity, nil
}

func (repo *UserActivityRepository) Get(ctx context.Context, reqQueryParam model.ReqQueryParamUserActivity) ([]model.UserActivity, int64, error) {
	var count int64
	var userActivities []model.UserActivity

	queryCount := fmt.Sprintf("SELECT COUNT(*) FROM useractivities INNER JOIN users ON useractivities.user_id=users.id WHERE useractivities.user_id=%d AND useractivities.activity_type IN(%v) AND useractivities.created_at BETWEEN '%s' AND '%s'", reqQueryParam.UserID, reqQueryParam.ActivityType, reqQueryParam.StartDate, reqQueryParam.EndDate)

	query := fmt.Sprintf("SELECT useractivities.* FROM useractivities INNER JOIN users ON useractivities.user_id=users.id WHERE useractivities.user_id=%d AND useractivities.activity_type IN(%v) AND useractivities.created_at BETWEEN '%s' AND '%s' ORDER BY useractivities.created_at DESC OFFSET %d ROWS FETCH NEXT %d ROWS ONLY", reqQueryParam.UserID, reqQueryParam.ActivityType, reqQueryParam.StartDate, reqQueryParam.EndDate, (reqQueryParam.PageNumber-1)*reqQueryParam.PerPage, reqQueryParam.PerPage)

	if err := repo.DB.Debug().WithContext(ctx).Raw(queryCount).Scan(&count).Error; err != nil {
		log.Printf("Failed Get User With Error : %v", err)
		return userActivities, 0, err
	}
	if err := repo.DB.Debug().WithContext(ctx).Raw(query).Find(&userActivities).Error; err != nil {
		log.Printf("Failed Get User With Error : %v", err)
		return userActivities, 0, err
	}

	return userActivities, count, nil
}

func (repo *UserActivityRepository) Update(ctx context.Context, userActivity model.UserActivity, ID uint64) error {
	if err := repo.DB.Debug().WithContext(ctx).Model(model.UserActivity{}).Where("id = ?", ID).Updates(&userActivity).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserActivityRepository) Delete(ctx context.Context, ID uint64) error {
	if err := repo.DB.Debug().WithContext(ctx).Delete(&model.UserActivity{}, ID).Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserActivityRepository) GetLastLogin(ctx context.Context, UserID uint64) (model.UserActivity, error) {
	var userActivity model.UserActivity

	query := fmt.Sprintf("SELECT TOP 1 useractivities.* FROM useractivities INNER JOIN users ON useractivities.user_id=users.id WHERE useractivities.user_id=%d AND useractivities.activity_type=%d ORDER BY useractivities.created_at DESC", UserID, constants.LOGIN)
	if err := repo.DB.Debug().WithContext(ctx).Raw(query).Find(&userActivity).Error; err != nil {
		return userActivity, err
	}

	return userActivity, nil
}

func (repo *UserActivityRepository) GetLastLogout(ctx context.Context, UserID uint64) (model.UserActivity, error) {
	var userActivity model.UserActivity

	query := fmt.Sprintf("SELECT TOP 1 useractivities.* FROM useractivities INNER JOIN users ON useractivities.user_id=users.id WHERE useractivities.user_id=%d AND useractivities.activity_type=%d ORDER BY useractivities.created_at DESC", UserID, constants.LOGOUT)
	if err := repo.DB.Debug().WithContext(ctx).Raw(query).Find(&userActivity).Error; err != nil {
		return userActivity, err
	}

	return userActivity, nil
}
