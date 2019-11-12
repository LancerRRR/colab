package model

import (
	"connect/http/response"
	"connect/server"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

type Route struct {
	ID           bson.ObjectId `json:"id" bson:"_id"`
	Description  string        `json:"description" bson:"description"`
	Path         string        `json:"path" bson:"path"`
	Method       string        `json:"method" bson:"method"`
	IsQuery      bool          `json:"isQuery" bson:"isQuery"`
	Deprecated   bool          `json:"deprecated" bson:"deprecated"`
	Service      string        `json:"service" bson:"service"`
	Auth         []string      `json:"auth" bson:"auth"`
	Request      interface{}   `json:"request" bson:"request"`
	Response     interface{}   `json:"response" bson:"response"`
	ResponseJSON interface{}   `json:"responseJSON" bson:"responseJSON"`
}

type Request struct {
	Type        string `json:"type" bson:"type"`
	Description string `json:"description" bson:"description"`
	IsRequired  bool   `json:"isRequired" bson:"isRequired"`
}

type Response struct {
	Type        string `json:"type" bson:"type"`
	Description string `json:"description" bson:"description"`
}

type RequestNested struct {
	Type        string      `json:"type" bson:"type"`
	Description string      `json:"description" bson:"description"`
	IsRequired  bool        `json:"isRequired" bson:"isRequired"`
	Nested      interface{} `json:"nested" bson:"nested"`
}

type ResponseNested struct {
	Type        string      `json:"type" bson:"type"`
	Description string      `json:"description" bson:"description"`
	Nested      interface{} `json:"nested" bson:"nested"`
}

func GetRouteAuth(path, method string) ([]string, error) {
	route := Route{}
	err := server.Mongodb.DB("Auth_colab").C("routes").Find(bson.M{"path": path, "method": method}).One(&route)
	return route.Auth, err
}

type Claim struct {
	UserID bson.ObjectId `json:"id"`
	Exp    int64         `json:"exp"`
	jwt.StandardClaims
}

func GetAllRoutes(param bson.M) ([]response.Route, error) {
	routes := []response.Route{}
	err := server.Mongodb.DB("Auth_colab").C("routes").Find(nil).Select(param).All(&routes)
	return routes, err
}
