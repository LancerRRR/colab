package main

import (
	"connect/http/request"
	"connect/model"
	"connect/server"
	"reflect"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var routes []model.Route

func init() {
	// ----------   User   ----------
	route := model.Route{
		Method:   "POST",
		Path:     "/user/register",
		Service:  "user",
		Request:  InterfaceToType(request.UserRegister{}),
		Required: []string{"email", "password"},
	}
	routes = append(routes, route)

	route = model.Route{
		Method:   "POST",
		Path:     "/user/login",
		Service:  "user",
		Request:  InterfaceToType(request.UserRegister{}),
		Required: []string{"username", "password"},
	}
	routes = append(routes, route)

	// -----------   Auth   --------
	route = model.Route{
		Method:   "POST",
		Path:     "/auth/routes",
		Service:  "auth",
		IsQuery:  true,
		Response: InterfaceToType(model.Route{}),
	}
	routes = append(routes, route)
}

func addRoutes() {
	for _, route := range routes {
		r := model.Route{}
		err := server.Mongodb.DB("Auth").C("routes").Find(bson.M{"service": route.Service, "path": route.Path}).One(&r)
		if err == mgo.ErrNotFound {
			route.ID = bson.NewObjectId()
			server.Mongodb.DB("Auth").C("routes").Insert(&route)
		}
	}
}

func InterfaceToType(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	switch v.(type) {
	case time.Time:
		return "string"
	}
	switch reflect.TypeOf(v).Kind() {
	case reflect.Struct:
		val := reflect.ValueOf(v)
		typeOfTstObj := val.Type()
		out := make(map[string]interface{}, 0)
		for i := 0; i < val.NumField(); i++ {
			fieldType := val.Field(i)
			value := InterfaceToType(fieldType.Interface())
			out[typeOfTstObj.Field(i).Tag.Get("json")] = value
		}
		return out
	case reflect.Slice:
		val := reflect.ValueOf(v)
		if val.Kind() != reflect.Struct || val.Kind() != reflect.Slice {
			return reflect.TypeOf(v).String()
		}
		return []interface{}{InterfaceToType(val.Index(0).Interface())}
	default:
		return reflect.TypeOf(v).String()
	}
}
