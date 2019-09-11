package handler

import (
	"connect/http/request"
	"connect/model"
	"connect/server"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func CheckAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(*CustomContext)
		routeAuth, err := model.GetRouteAuth(c.Path())
		if err != nil {
			return echo.ErrNotFound
		}
		for _, auth := range routeAuth {
			if auth == "guest" {
				return next(c)
			}
		}
		claim := model.Claim{}
		tokenStr := c.Request().Header.Get("Authorization")
		fmt.Println(tokenStr)
		token, err := jwt.ParseWithClaims(tokenStr, &claim, func(token *jwt.Token) (interface{}, error) {
			return []byte("shawn"), nil
		})
		if err != nil {
			return echo.ErrUnauthorized
		}
		if !token.Valid {
			return echo.ErrUnauthorized
		}
		if time.Now().Unix() > claim.Exp {
			return cc.RespondError(401, "expired")
		}
		code, err := server.Cl.Get(claim.UserID.Hex()).Result()
		if err != nil {
			return cc.RespondError(500, err)
		}
		if code != tokenStr {
			return echo.ErrUnauthorized
		}
		user, err := model.GetUserByQuery(bson.M{"_id": claim.UserID}, nil)
		if err != nil {
			return cc.RespondError(400, "user not found")
		}
		for _, auth := range user.Auth {
			if auth == "developer" {
				return next(c)
			}
		}
		for _, auth1 := range routeAuth {
			if auth1 == "user" {
				return next(c)
			}
			for _, auth := range user.Auth {
				if auth1 == auth {
					return next(c)
				}
			}
		}
		return echo.ErrUnauthorized
	}
}

func (h *Handler) GetAllRoutes(c echo.Context) error {
	cc := c.(*CustomContext)
	params := request.Parameters{}
	err := c.Bind(&params)
	if err != nil {
		return cc.RespondError(500, err)
	}
	routes, err := model.GetAllRoutes(params.Param)
	if err != nil {
		return cc.RespondError(500, err)
	}
	return cc.Respond(200, routes)
}
