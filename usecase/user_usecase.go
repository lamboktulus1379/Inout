package usecase

import (
	"context"
	"crypto/md5"
	"fmt"
	"inout/constants"
	"inout/domain/dto"
	"inout/domain/model"
	"inout/domain/repository"
	"log"
	"math"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

type IUserUsecase interface {
	Login(ctx context.Context, req model.ReqLogin) dto.ResLogin
	Logout(ctx context.Context, userID uint64) dto.ResLogout
	Checkin(ctx context.Context, req model.ReqUserActivity) dto.Res
	History(ctx context.Context, reqQueryParama model.ReqQueryParamUserActivity) dto.ResUserActivity
	Create(ctx context.Context, req model.ReqCreateUserActivity) dto.Res
	Update(ctx context.Context, req model.ReqUpdateUserActivity, reqURIParam model.ReqURIParamUserActivity) dto.Res
	Delete(ctx context.Context, reqURIParam model.ReqURIParamUserActivity) dto.Res
}

type UserUsecase struct {
	userRepository         repository.IUser
	userActivityRepository repository.IUserActivity
}

func NewUserUsecase(userRepository repository.IUser, userActivityRepository repository.IUserActivity) IUserUsecase {
	return &UserUsecase{userRepository: userRepository, userActivityRepository: userActivityRepository}
}

func (userUsecase *UserUsecase) Login(ctx context.Context, req model.ReqLogin) dto.ResLogin {
	var res dto.ResLogin

	user, err := userUsecase.userRepository.GetByUserName(ctx, req.UserName)
	log.Printf("Username: %s\n", user.UserName)
	if err != nil {
		log.Printf("User not found. %v\n", err)
		res.ResponseCode = "401"
		res.ResponseMessage = "Unautorized."
		return res
	}
	md5Req := fmt.Sprintf("%x", md5.Sum([]byte(req.Password)))

	if md5Req != user.Password {
		res.ResponseCode = "401"
		res.ResponseMessage = "Unautorized."
		return res
	}

	secretKey := os.Getenv("SECRET_KEY")
	fmt.Println("Secret Key: ", secretKey)
	mySigningKey := []byte(secretKey)

	// Create the Claims
	expiration := time.Now().Add(5 * time.Minute)

	claims := model.UserClaims{
		UserName: req.UserName,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expiration.Unix(),
			Issuer:    fmt.Sprint(user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(mySigningKey)
	fmt.Printf("%v %v", accessToken, err)
	if err != nil {
		res.ResponseCode = "401"
		res.ResponseMessage = "Unautorized"
		return res
	}
	userActivity := model.UserActivity{
		UserID:       user.ID,
		ActivityType: constants.LOGIN,
		CreatedBy:    user.ID,
		UpdatedBy:    user.ID,
	}

	_, err = userUsecase.userActivityRepository.Create(ctx, userActivity)
	if err != nil {
		log.Printf("An error occurred %v", err)
		res.ResponseCode = "500"
		res.ResponseMessage = fmt.Sprintf("An error occurred %v", err)
		return res
	}

	res.ResponseCode = "200"
	res.ResponseMessage = "Success"
	res.Data.AccessToken = accessToken
	res.Data.ExpiresAt = expiration.Unix()

	return res
}

func (userUsecase *UserUsecase) Logout(ctx context.Context, userID uint64) dto.ResLogout {
	var res dto.ResLogout

	userActivity := model.UserActivity{
		UserID:       userID,
		ActivityType: constants.LOGOUT,
		CreatedBy:    userID,
		UpdatedBy:    userID,
	}

	_, err := userUsecase.userActivityRepository.Create(ctx, userActivity)
	if err != nil {
		log.Printf("An error occurred %v", err)
		res.ResponseCode = "500"
		res.ResponseMessage = fmt.Sprintf("An error occurred %v", err)
		return res
	}

	res.ResponseCode = "200"
	res.ResponseMessage = "Success"

	return res
}

func (userUsecase *UserUsecase) Checkin(ctx context.Context, req model.ReqUserActivity) dto.Res {
	var res dto.Res

	userActivity := model.UserActivity{
		UserID:       req.UserID,
		ActivityType: constants.CHECKIN,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		CreatedBy:    req.UserID,
		UpdatedBy:    req.UserID,
	}

	checkout, err := userUsecase.userActivityRepository.Create(ctx, userActivity)
	if err != nil {
		log.Printf("An error occurred %v", err)
		res.ResponseCode = "500"
		res.ResponseMessage = fmt.Sprintf("An error occurred %v", err)
		return res
	}

	log.Printf("%#v", checkout)

	res.ResponseCode = "200"
	res.ResponseMessage = "Success"

	return res
}

func (userUsecase *UserUsecase) History(ctx context.Context, reqQueryParam model.ReqQueryParamUserActivity) dto.ResUserActivity {
	var res dto.ResUserActivity

	result, count, err := userUsecase.userActivityRepository.Get(ctx, reqQueryParam)
	if err != nil {
		log.Printf("An error occurred %v", err)
		res.ResponseCode = "500"
		res.ResponseMessage = "Internal server error"
		return res
	}
	userActivities := make([]dto.UserActivity, 0)
	for _, userActivity := range result {
		ua := dto.UserActivity{
			ID:           userActivity.ID,
			UserID:       userActivity.UserID,
			ActivityType: constants.ActivityType[int(userActivity.ActivityType)],
			Latitude:     userActivity.Latitude,
			Longitude:    userActivity.Longitude,
			CreatedAt:    userActivity.CreatedAt,
			CreatedBy:    userActivity.CreatedBy,
			UpdatedAt:    userActivity.UpdatedAt,
			UpdatedBy:    userActivity.UpdatedBy,
		}
		userActivities = append(userActivities, ua)
	}
	pagination := dto.Pagination{
		PageNumber:  reqQueryParam.PageNumber,
		PerPage:     reqQueryParam.PerPage,
		TotalPage:   int(math.Ceil(float64(count) / float64(reqQueryParam.PerPage))),
		TotalRecord: int(count),
	}
	res.Data = userActivities
	res.Pagination = &pagination
	res.ResponseCode = "200"
	res.ResponseMessage = "Success"

	return res
}

func (userUsecase *UserUsecase) Create(ctx context.Context, req model.ReqCreateUserActivity) dto.Res {
	var res dto.Res

	_, ok := constants.ActivityType[int(req.ActivityType)]
	if !ok {
		res.ResponseCode = "400"
		res.ResponseMessage = "Bad Request"
		return res
	}

	userActivity := model.UserActivity{
		UserID:       req.UserID,
		ActivityType: req.ActivityType,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		CreatedBy:    req.UserID,
		UpdatedBy:    req.UserID,
	}

	create, err := userUsecase.userActivityRepository.Create(ctx, userActivity)
	if err != nil {
		log.Printf("An error occurred %v", err)
		res.ResponseCode = "500"
		res.ResponseMessage = fmt.Sprintf("An error occurred %v", err)
		return res
	}

	log.Printf("%#v", create)

	res.ResponseCode = "200"
	res.ResponseMessage = "Success"

	return res
}

func (userUsecase *UserUsecase) Update(ctx context.Context, req model.ReqUpdateUserActivity, reqURIParam model.ReqURIParamUserActivity) dto.Res {
	var res dto.Res

	_, ok := constants.ActivityType[int(req.ActivityType)]
	if !ok {
		res.ResponseCode = "400"
		res.ResponseMessage = "Bad Request"
		return res
	}

	userActivity := model.UserActivity{
		ActivityType: req.ActivityType,
		Latitude:     req.Latitude,
		Longitude:    req.Longitude,
		CreatedBy:    req.UserID,
		UpdatedBy:    req.UserID,
	}

	userActivity, err := userUsecase.userActivityRepository.GetById(ctx, reqURIParam.ID)
	if err != nil {
		res.ResponseCode = "404"
		res.ResponseMessage = err.Error()
		return res
	}

	if userActivity.UserID != req.UserID {
		res.ResponseCode = "400"
		res.ResponseMessage = "User ID does not match"
		return res
	}

	err = userUsecase.userActivityRepository.Update(ctx, userActivity, reqURIParam.ID)
	if err != nil {
		log.Printf("An error occurred %v", err)
		res.ResponseCode = "500"
		res.ResponseMessage = fmt.Sprintf("An error occurred %v", err)
		return res
	}

	res.ResponseCode = "200"
	res.ResponseMessage = "Success"

	return res
}

func (userUsecase *UserUsecase) Delete(ctx context.Context, reqURIParam model.ReqURIParamUserActivity) dto.Res {
	var res dto.Res

	userActivity, err := userUsecase.userActivityRepository.GetById(ctx, reqURIParam.ID)
	if err != nil {
		res.ResponseCode = "404"
		res.ResponseMessage = err.Error()
		return res
	}

	if userActivity.UserID != reqURIParam.UserID {
		res.ResponseCode = "400"
		res.ResponseMessage = "User ID does not match"
		return res
	}

	err = userUsecase.userActivityRepository.Delete(ctx, reqURIParam.ID)
	if err != nil {
		log.Printf("An error occurred %v", err)
		res.ResponseCode = "500"
		res.ResponseMessage = fmt.Sprintf("An error occurred %v", err)
		return res
	}

	res.ResponseCode = "200"
	res.ResponseMessage = "Success"

	return res
}
