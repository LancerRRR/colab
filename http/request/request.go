package request

import (
	"gopkg.in/mgo.v2/bson"
)

type UserRegister struct {
	Phone       string  `json:"phone" bson:"phone" required:"1" description:"电话号码"`
	Email       string  `json:"email" bson:"email" description:"邮箱，用于登录"`
	NickName    string  `json:"nickName" bson:"nickName" description:"用户昵称"`
	Avatar      string  `json:"avatar" bson:"avatar" description:"头像"`
	Password    string  `json:"password" bson:"password" description:"密码"`
	IsVerified  bool    `json:"isVerified" bson:"isVerified" description:"是否认证成功"`
	Institution string  `json:"institution" bson:"institution" description:"学校信息"`
	Skills      []Skill `json:"skills" bson:"skills" description:"掌握的技能"`
	Role        string  `json:"role" bson:"role" description:"用户角色, 项目发布者or参与者"`
}

type Skill struct {
	Name string `json:"name"`
}

type UserLogin struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type Parameters struct {
	Param bson.M `json:"param"`
}
