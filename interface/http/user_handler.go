package http

import (
	"fmt"
	"inout/constants"
	"inout/domain/dto"
	"inout/domain/model"
	"inout/usecase"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type IUserHandler interface {
	Login(c *gin.Context)

	// TODO
	Logout(c *gin.Context)
	Checkin(c *gin.Context)
	Checkout(c *gin.Context)
	History(c *gin.Context)
	Create(c *gin.Context)
	Update(c *gin.Context)
	Delete(c *gin.Context)
}

type UserHandler struct {
	userUsecase usecase.IUserUsecase
}

func NewUserHandler(userUsecase usecase.IUserUsecase) IUserHandler {
	return &UserHandler{userUsecase: userUsecase}
}

func (userHandler *UserHandler) Login(c *gin.Context) {
	var req model.ReqLogin

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("An error occurred: %v", err)
		c.JSON(http.StatusBadRequest, fmt.Sprintf("An error occurred: %v", err.Error()))
	}

	res := userHandler.userUsecase.Login(c.Request.Context(), req)

	c.JSON(constants.ResponseCode[res.ResponseCode], res)
}

func (userHandler *UserHandler) Logout(c *gin.Context) {
	userIDStr := c.MustGet("user_id").(string)
	intVar, _ := strconv.ParseInt(userIDStr, 0, 32)
	userID := uint64(intVar)
	res := userHandler.userUsecase.Logout(c.Request.Context(), userID)

	c.JSON(constants.ResponseCode[res.ResponseCode], res)
}

func (userHandler *UserHandler) Checkin(c *gin.Context) {
	var req model.ReqUserActivity
	var res dto.Res

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("An error occurred: %v", err)
		res.ResponseCode = "400"
		res.ResponseMessage = fmt.Sprintf("An error occurred: %v", err.Error())
		c.JSON(http.StatusBadRequest, res)
	}

	userIDStr := c.MustGet("user_id").(string)
	intVar, _ := strconv.ParseInt(userIDStr, 0, 32)
	req.UserID = uint64(intVar)
	res = userHandler.userUsecase.Checkin(c.Request.Context(), req)

	c.JSON(constants.ResponseCode[res.ResponseCode], res)
}

func (userHandler *UserHandler) Checkout(c *gin.Context) {
	var req model.ReqUserActivity
	var res dto.Res

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("An error occurred: %v", err)
		res.ResponseCode = "400"
		res.ResponseMessage = fmt.Sprintf("An error occurred: %v", err.Error())
		c.JSON(http.StatusBadRequest, res)
	}

	userIDStr := c.MustGet("user_id").(string)
	intVar, _ := strconv.ParseInt(userIDStr, 0, 32)
	req.UserID = uint64(intVar)
	res = userHandler.userUsecase.Checkin(c.Request.Context(), req)

	c.JSON(constants.ResponseCode[res.ResponseCode], res)
}

func (userHandler *UserHandler) History(c *gin.Context) {
	var res dto.ResUserActivity
	var reqQueryParam model.ReqQueryParamUserActivity

	fmt.Println("User Login: ", c.MustGet("user_id"))
	if err := c.ShouldBindQuery(&reqQueryParam); err != nil {
		log.Printf("An error occurred: %v", err)
		res.ResponseCode = "400"
		res.ResponseMessage = fmt.Sprintf("An error occurred: %v", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	log.Printf("Request: %#v\n", reqQueryParam)
	if reqQueryParam.PageNumber == 0 {
		reqQueryParam.PageNumber = 1
	}
	if reqQueryParam.PerPage == 0 {
		reqQueryParam.PerPage = 10
	}
	if reqQueryParam.ActivityType == "" {
		reqQueryParam.ActivityType = "3,4"
	}
	layout := "2006-01-02"
	layoutDB := "2006-01-02 15:04:05"
	a, err := time.Parse(layout, reqQueryParam.StartDate)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		res.ResponseCode = "400"
		res.ResponseMessage = fmt.Sprintf("An error occurred: %v", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	reqQueryParam.StartDate = a.Format(layoutDB)

	b, err := time.Parse(layout, reqQueryParam.EndDate)
	if err != nil {
		log.Printf("An error occurred: %v", err)
		res.ResponseCode = "400"
		res.ResponseMessage = fmt.Sprintf("An error occurred: %v", err.Error())
		c.JSON(http.StatusBadRequest, res)
		return
	}

	b = b.Add(time.Hour*23 + time.Minute*59 + time.Second*59 + time.Millisecond*999)
	reqQueryParam.EndDate = b.Format(layoutDB)

	userIDStr := c.MustGet("user_id").(string)
	intVar, _ := strconv.ParseInt(userIDStr, 0, 64)
	reqQueryParam.UserID = uint64(intVar)
	res = userHandler.userUsecase.History(c.Request.Context(), reqQueryParam)

	c.JSON(constants.ResponseCode[res.ResponseCode], res)
}

func (userHandler *UserHandler) Create(c *gin.Context) {
	var req model.ReqCreateUserActivity
	var res dto.Res

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("An error occurred: %v", err)
		res.ResponseCode = "400"
		res.ResponseMessage = fmt.Sprintf("An error occurred: %v", err.Error())
		c.JSON(http.StatusBadRequest, res)
	}
	userIDStr := c.MustGet("user_id").(string)
	intVar, _ := strconv.ParseInt(userIDStr, 0, 64)
	req.UserID = uint64(intVar)
	res = userHandler.userUsecase.Create(c.Request.Context(), req)

	c.JSON(constants.ResponseCode[res.ResponseCode], res)
}

func (userHandler *UserHandler) Update(c *gin.Context) {
	var req model.ReqUpdateUserActivity
	var reqURIParam model.ReqURIParamUserActivity
	var res dto.Res

	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("An error occurred: %v", err)
		res.ResponseCode = "400"
		res.ResponseMessage = fmt.Sprintf("An error occurred: %v", err.Error())
		c.JSON(http.StatusBadRequest, res)
	}
	if err := c.ShouldBindUri(&reqURIParam); err != nil {
		log.Printf("An error occurred: %v", err)
		res.ResponseCode = "400"
		res.ResponseMessage = fmt.Sprintf("An error occurred: %v", err.Error())
		c.JSON(http.StatusBadRequest, res)
	}
	userIDStr := c.MustGet("user_id").(string)
	intVar, _ := strconv.ParseInt(userIDStr, 0, 64)
	req.UserID = uint64(intVar)

	res = userHandler.userUsecase.Update(c.Request.Context(), req, reqURIParam)

	c.JSON(constants.ResponseCode[res.ResponseCode], res)
}

func (userHandler *UserHandler) Delete(c *gin.Context) {
	var reqURIParam model.ReqURIParamUserActivity
	var res dto.Res

	if err := c.ShouldBindUri(&reqURIParam); err != nil {
		log.Printf("An error occurred: %v", err)
		res.ResponseCode = "400"
		res.ResponseMessage = fmt.Sprintf("An error occurred: %v", err.Error())
		c.JSON(http.StatusBadRequest, res)
	}
	userIDStr := c.MustGet("user_id").(string)
	intVar, _ := strconv.ParseInt(userIDStr, 0, 64)
	reqURIParam.UserID = uint64(intVar)

	res = userHandler.userUsecase.Delete(c.Request.Context(), reqURIParam)

	c.JSON(constants.ResponseCode[res.ResponseCode], res)
}
