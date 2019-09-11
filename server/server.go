package server

import (
	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
)

var Mongodb *mgo.Session
var Cl = redis.NewClient(&redis.Options{
	Addr:     "127.0.0.1:6379",
	Password: "",
	DB:       0,
})
