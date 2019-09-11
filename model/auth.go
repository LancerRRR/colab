package model

import (
	"connect/server"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/mgo.v2/bson"
)

type Route struct {
	ID         bson.ObjectId `json:"id" bson:"_id"`
	Path       string        `json:"path" bson:"path"`
	Method     string        `json:"method" bson:"Method"`
	IsQuery    bool          `json:"isQuery" bson:"isQuery"`
	Deprecated bool          `json:"deprecated" bson:"deprecated"`
	Service    string        `json:"service" bson:"service"`
	Auth       []string      `json:"auth" bson:"auth"`
	Required   []string      `json:"required" bson:"required"`
	Request    interface{}   `json:"request" bson:"request"`
	Response   interface{}   `json:"response" bson:"response"`
}

func GetRouteAuth(path string) ([]string, error) {
	route := Route{}
	err := server.Mongodb.DB("Auth").C("routes").Find(bson.M{"path": path}).One(&route)
	return route.Auth, err
}

type Claim struct {
	UserID bson.ObjectId `json:"id"`
	Exp    int64         `json:"exp"`
	jwt.StandardClaims
}

func GetAllRoutes(param bson.M) ([]map[string]interface{}, error) {
	routes := make([]map[string]interface{}, 0)
	err := server.Mongodb.DB("Auth").C("routes").Find(nil).Select(param).All(&routes)
	return routes, err
}
