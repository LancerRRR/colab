package request

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type UserRegister struct {
	Phone       string    `json:"phone" bson:"phone"`
	Email       string    `json:"email" bson:"email"`
	NickName    string    `json:"nickName" bson:"nickName"`
	Avatar      string    `json:"avatar" bson:"avatar"`
	Password    string    `json:"password" bson:"password"`
	IsVerified  bool      `json:"isVerified" bson:"isVerified"`
	Institution string    `json:"institution" bson:"institution"`
	Skills      []string  `json:"skills" bson:"skills"`
	Role        string    `json:"role" bson:"role"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
}

type UserLogin struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type Parameters struct {
	Param bson.M `json:"param"`
}
