package handler

import (
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
)

type (
	Handler struct {
		DB *mgo.Session
	}
)

type CustomContext struct {
	echo.Context
}

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

func (c *CustomContext) Respond(code int, response interface{}) error {
	resp := make(map[string]interface{}, 0)
	resp["data"] = response
	return c.JSON(code, resp)
}

func (c *CustomContext) RespondError(code int, response interface{}) error {
	resp := make(map[string]interface{}, 0)
	resp["error"] = response
	return c.JSON(code, resp)
}
