package handler

import (
	"connect/http/request"
	"connect/model"
	"connect/service"
	"time"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (h *Handler) Register(c echo.Context) error {
	cc := c.(*CustomContext)
	userReq := model.User{}
	err := c.Bind(&userReq)
	if err != nil {
		return cc.RespondError(500, err)
	}
	if len(userReq.Email) == 0 || len(userReq.Password) < 6 {
		return cc.RespondError(400, "invalid format!")
	}
	_, err = model.GetUserByQuery(bson.M{"email": userReq.Email}, nil)
	if err == nil {
		return cc.RespondError(400, "the username is already registered!")
	}
	if err == mgo.ErrNotFound {
		userReq.ID = bson.NewObjectId()
		hashedPassword, _ := service.HashPassword(userReq.Password)
		userReq.Password = hashedPassword
		userReq.CreatedAt = time.Now()
		err := model.CreateUser(userReq)
		if err != nil {
			return cc.RespondError(500, err)
		}
		return cc.Respond(200, "created")
	}
	return cc.RespondError(500, err)
}

func (h *Handler) Login(c echo.Context) error {
	cc := c.(*CustomContext)
	userLogin := request.UserLogin{}
	err := c.Bind(&userLogin)
	if err != nil {
		return cc.RespondError(500, err)
	}
	resp, err := service.Login(userLogin.UserName, userLogin.Password)
	if err == mgo.ErrNotFound {
		return cc.RespondError(400, "Invalid username")
	}
	if err != nil {
		return cc.RespondError(400, err)
	}
	return cc.Respond(200, resp)
}
