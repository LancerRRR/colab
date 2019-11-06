package model

import (
	"connect/server"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var dbUser = server.Mongodb.DB("user")

type User struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Phone      string        `json:"phone" bson:"phone"`
	Email      string        `json:"email" bson:"email"`
	NickName   string        `json:"nickName" bson:"nickName"`
	Avatar     string        `json:"avatar" bson:"avatar"`
	Password   string        `json:"password" bson:"password"`
	IsVerified bool          `json:"isVerified" bson:"isVerified"`
	Auth       []string      `json:"auth" bson:"auth"`
	CreatedAt  time.Time     `json:"createdAt" bson:"createdAt"`
}

func CreateUser(user User) error {
	err := server.Mongodb.DB("user").C("user").Insert(user)
	return err
}

func GetUserByQuery(query bson.M, param bson.M) (User, error) {
	user := User{}
	err := server.Mongodb.DB("user").C("user").Find(query).Select(param).One(&user)
	return user, err
}
