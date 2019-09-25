package main

import (
	"connect/http/request"
	"connect/http/response"
	"connect/model"
	"connect/server"
	"fmt"
	"reflect"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var routes []model.Route
var isRequired bool
var description string
var isRequest bool
var accessed int
var kind string

func init() {
	// ----------   User   ----------
	route := model.Route{
		Method:      "POST",
		Description: "用户注册",
		Path:        "/user/register",
		Service:     "user",
		Request: request.UserRegister{
			Skills: []request.Skill{request.Skill{}},
		},
	}
	routes = append(routes, route)

	route = model.Route{
		Method:      "POST",
		Description: "用户登录",
		Path:        "/user/login",
		Service:     "user",
		Request:     request.UserLogin{},
		Response:    response.UserLoginResp{},
	}
	routes = append(routes, route)

	// -----------   Auth   --------
	route = model.Route{
		Method:      "POST",
		Description: "获取所有Routes",
		Path:        "/auth/routes",
		Auth:        []string{"guest"},
		Service:     "auth",
		IsQuery:     true,
		Response:    []response.Route{response.Route{ID: bson.NewObjectId()}},
	}
	routes = append(routes, route)
}

func addRoutes() {
	currentRoutes, err := model.GetAllRoutes(nil)
	if err != nil {
		panic(err)
	}
	routesMap := make(map[string]response.Route, 0)
	for _, currentRoute := range currentRoutes {
		routesMap[currentRoute.Path] = currentRoute
	}
	for _, route := range routes {
		if route.Request != nil {
			route.Request = InterfaceToType(route.Request)
		}
		if route.Response != nil {
			route.ResponseJSON = InterfaceToJSON(route.Response)
			route.Response = InterfaceToType(route.Response)
		}
		r, ok := routesMap[route.Path]
		if !ok {
			route.ID = bson.NewObjectId()
			err := server.Mongodb.DB("Auth").C("routes").Insert(&route)
			if err != nil {
				panic(err)
			}
		} else {
			route.ID = r.ID
			err := server.Mongodb.DB("Auth").C("routes").UpdateId(route.ID, route)
			if err != nil {
				panic(err)
			}
			delete(routesMap, route.Path)
		}
	}

	for _, v := range routesMap {
		v.Deprecated = true
		err := server.Mongodb.DB("Auth").C("routes").UpdateId(v.ID, v)
		if err != nil {
			panic(err)
		}
	}
}

func InterfaceToJSON(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	switch v.(type) {
	case time.Time:
		return "string"
	}
	fmt.Println(reflect.TypeOf(v).Kind() == reflect.Slice)
	switch reflect.TypeOf(v).Kind() {
	case reflect.Struct:
		val := reflect.ValueOf(v)
		typeOfTstObj := val.Type()
		out := make(map[string]interface{}, 0)
		for i := 0; i < val.NumField(); i++ {
			fieldType := val.Field(i)
			jsonTag := typeOfTstObj.Field(i).Tag.Get("json")
			value := InterfaceToJSON(fieldType.Interface())
			key := ""
			for i := 0; i < len(jsonTag); i++ {
				if string(jsonTag[i]) == "," {
					break
				}
				key += string(jsonTag[i])
			}
			out[key] = value
		}
		return out
	case reflect.Slice:
		val := reflect.ValueOf(v)
		if val.Len() == 0 {
			return []interface{}{reflect.TypeOf(v).String()[2:]}
		}
		return []interface{}{InterfaceToJSON(val.Index(0).Interface())}
	default:
		return reflect.TypeOf(v).String()
	}
}

func InterfaceToType(v interface{}) interface{} {
	if v == nil {
		if isRequest {
			field := model.Request{}
			field.Type = "interface"
			return field
		} else {
			field := model.Response{}
			field.Type = "interface"
			return field
		}
	}
	// switch v.(type) {
	// case time.Time:
	// 	return "string"
	// }
	switch reflect.TypeOf(v).Kind() {
	case reflect.Struct:
		if (len(reflect.TypeOf(v).String()) >= 7 && reflect.TypeOf(v).String()[:7] == "request") || (len(reflect.TypeOf(v).String()) >= 8 && reflect.TypeOf(v).String()[:8] == "response") {
			accessed1 := accessed
			kind1 := kind
			if reflect.TypeOf(v).String()[:7] == "request" {
				isRequest = true
			} else {
				isRequest = false
			}
			val := reflect.ValueOf(v)
			typeOfTstObj := val.Type()
			out := make(map[string]interface{}, 0)
			for i := 0; i < val.NumField(); i++ {
				fieldType := val.Field(i)
				isRequired = false
				description = ""
				if typeOfTstObj.Field(i).Tag.Get("required") == "1" {
					isRequired = true
				}
				description = typeOfTstObj.Field(i).Tag.Get("description")
				kind = "object"
				accessed = 1
				value := InterfaceToType(fieldType.Interface())
				jsonTag := typeOfTstObj.Field(i).Tag.Get("json")
				key := ""
				for i := 0; i < len(jsonTag); i++ {
					if string(jsonTag[i]) == "," {
						break
					}
					key += string(jsonTag[i])
				}
				out[key] = value
			}
			if !isRequest {
				output := model.RequestNested{}
				output.Nested = out
				output.Type = kind1
				output.Description = description
				output.IsRequired = isRequired
				return output
			}
			if accessed1 == 1 {
				if isRequest {
					output := model.RequestNested{}
					output.Nested = out
					output.Type = kind1
					output.Description = description
					output.IsRequired = isRequired
					return output
				} else {
					output := model.ResponseNested{}
					output.Nested = out
					output.Type = kind1
					output.Description = description
					return output
				}
			}
			accessed = 0
			return out
		} else {
			if isRequest {
				field := model.Request{}
				field.Type = reflect.TypeOf(v).String()
				field.Description = description
				field.IsRequired = isRequired
				return field
			} else {
				field := model.Response{}
				field.Type = reflect.TypeOf(v).String()
				field.Description = description
				return field
			}
		}
	case reflect.Slice:
		val := reflect.ValueOf(v)
		if val.Len() == 0 {
			if isRequest {
				field := model.Request{}
				field.Type = reflect.TypeOf(v).String()
				field.Description = description
				field.IsRequired = isRequired
				return field
			} else {
				field := model.Response{}
				field.Type = reflect.TypeOf(v).String()
				field.Description = description
				return field
			}
		}
		kind = "array"
		accessed = 1
		return InterfaceToType(val.Index(0).Interface())
	default:
		if isRequest {
			field := model.Request{}
			field.Type = reflect.TypeOf(v).String()
			field.Description = description
			field.IsRequired = isRequired
			return field
		} else {
			field := model.Response{}
			field.Type = reflect.TypeOf(v).String()
			field.Description = description
			return field
		}
	}
}
