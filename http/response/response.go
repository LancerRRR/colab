package response

import "gopkg.in/mgo.v2/bson"

type UserLoginResp struct {
	UserID bson.ObjectId `json:"userID"`
	Token  string        `json:"token"`
}
