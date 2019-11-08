package main

import (
	"connect/handler"
	"connect/model"
	"connect/server"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"

	"crypto/tls"
	"net"

	"github.com/labstack/echo/middleware"
)

const (
	MONGODB_URL = "mongodb://dev:gogox123@gogo-api-dev-shard-00-00-uptzo.mongodb.net:27017,gogo-api-dev-shard-00-01-uptzo.mongodb.net:27017,gogo-api-dev-shard-00-02-uptzo.mongodb.net:27017/test?replicaSet=GOGO-API-DEV-shard-0&authSource=admin"
)

func main() {
	e := echo.New()
	e.Use(middleware.CORS())
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &handler.CustomContext{c}
			return next(cc)
		}
	})
	e.Use(handler.CheckAuth)
	tlsConfig := &tls.Config{}
	tlsConfig.InsecureSkipVerify = true
	dialInfo, err := mgo.ParseURL(MONGODB_URL)
	dialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
		conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
		return conn, err
	}

	db, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		e.Logger.Fatal(err)
	}
	server.Mongodb = db
	model.InitCollections()
	defer db.Close()
	addRoutes()
	h := &handler.Handler{DB: db}

	// User Routes
	e.POST("/user/register", h.Register)
	e.POST("/user/login", h.Login)
	e.GET("/auth/routes", h.GetAllRoutes)
	e.Static("/", "../public")
	e.GET("/ws", h.Connect)
	e.POST("/card/card", h.CreateCard)
	e.GET("/card/card", h.GetCards)
	e.POST("/card/list", h.CreateList)
	e.PUT("/card/list", h.UpdateList)
	e.PUT("/card/sort", h.UpdateCardSort)
	e.Logger.Fatal(e.Start(":3002"))
}
