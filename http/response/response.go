package response

import "gopkg.in/mgo.v2/bson"

type UserLoginResp struct {
	UserID bson.ObjectId `json:"userID" bson:"userID"`
	Token  string        `json:"token" bson:"token"`
}

type Route struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Path         string        `json:"path,omitempty" bson:"path"`
	Description  string        `json:"description,omitempty" bson:"description"`
	Method       string        `json:"method,omitempty" bson:"method"`
	IsQuery      bool          `json:"isQuery,omitempty" bson:"isQuery"`
	Deprecated   bool          `json:"deprecated,omitempty" bson:"deprecated"`
	Service      string        `json:"service,omitempty" bson:"service"`
	Auth         []string      `json:"auth,omitempty" bson:"auth"`
	Request      interface{}   `json:"request,omitempty" bson:"request"`
	Response     interface{}   `json:"response,omitempty" bson:"response"`
	ResponseJSON interface{}   `json:"responseJSON,omitempty" bson:"responseJSON"`
}
