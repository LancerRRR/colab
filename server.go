package main

import (
	"connect/handler"
	"connect/server"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"

	"crypto/tls"
	"net"

	"github.com/labstack/echo/middleware"
)

const (
	MONGODB_URL = "mongodb://admin:Wk950523@cluster0-shard-00-00-b0tff.mongodb.net:27017,cluster0-shard-00-01-b0tff.mongodb.net:27017,cluster0-shard-00-02-b0tff.mongodb.net:27017/test?replicaSet=Cluster0-shard-0&authSource=admin"
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
	defer db.Close()
	addRoutes()
	h := &handler.Handler{DB: db}
	// User Routes
	e.POST("/user/register", h.Register)
	e.GET("/auth/routes", h.GetAllRoutes)
	e.Logger.Fatal(e.Start(":8000"))
}
