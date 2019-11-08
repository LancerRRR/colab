package handler

import (
	"connect/model"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	Channels   = make(map[string]map[string]Channel, 0)
	UserRoutes = make(map[string][]string, 0)
	upgrader   = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
)

type Channel struct {
	UserID     bson.ObjectId
	Connection *websocket.Conn
}

type SuccessMsg struct {
	Code   int64                  `json:"code"`
	Route  string                 `json:"route"`
	Method string                 `json:"method"`
	Data   map[string]interface{} `json:data"`
}

type ErrorMsg struct {
	Code   int64  `json:"code"`
	Route  string `json:"route"`
	Method string `json:"method"`
	Error  string `json:"error"`
}

type (
	Handler struct {
		DB *mgo.Session
	}
)

type CustomContext struct {
	echo.Context
}

const (
	// Key (Should come from somewhere else).
	Key = "secret"
)

func (c *CustomContext) Respond(code int, response interface{}) error {
	resp := make(map[string]interface{}, 0)
	resp["data"] = response
	return c.JSON(code, resp)
}

func (c *CustomContext) RespondError(code int, response interface{}) error {
	resp := make(map[string]interface{}, 0)
	resp["error"] = response
	return c.JSON(code, resp)
}

func (h *Handler) Connect(c echo.Context) error {
	tokenStr := c.QueryParam("token")
	cc := c.(*CustomContext)
	claim := model.Claim{}
	token, err := jwt.ParseWithClaims(tokenStr, &claim, func(token *jwt.Token) (interface{}, error) {
		return []byte("shawn"), nil
	})
	if err != nil {
		return echo.ErrUnauthorized
	}
	if !token.Valid {
		return echo.ErrUnauthorized
	}
	if time.Now().Unix() > claim.Exp {
		return cc.RespondError(401, "expired")
	}
	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println(err)
		cc.RespondError(500, err)
	}
	channel := Channel{
		UserID:     claim.UserID,
		Connection: conn,
	}
	for {
		t, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			currentRoutes := UserRoutes[tokenStr]
			for _, route := range currentRoutes {
				delete(Channels[route], tokenStr)
			}
			delete(UserRoutes, tokenStr)
			conn.Close()
			return cc.RespondError(500, err)
		}
		if string(msg) == "ping" {
			conn.WriteMessage(t, []byte("pong"))
		} else {
			routes := []string{}
			json.Unmarshal(msg, &routes)
			log.Println(routes)
			currentRoutes, ok := UserRoutes[tokenStr]
			if !ok {
				UserRoutes[tokenStr] = routes
				for _, route := range routes {
					if _, ok := Channels[route]; !ok {
						Channels[route] = make(map[string]Channel, 0)
					}
					Channels[route][tokenStr] = channel
				}
			} else {
				for _, route := range currentRoutes {
					delete(Channels[route], tokenStr)
				}
				UserRoutes[tokenStr] = routes
				for _, route := range routes {
					if _, ok := Channels[route]; !ok {
						Channels[route] = make(map[string]Channel, 0)
					}
					Channels[route][tokenStr] = channel
				}
			}
			//conn.WriteMessage(t, []byte(claim.UserID.Hex()+": "+string(msg)))
		}
		continue
	}
}

func PushToUsers(data interface{}, route string) {
	for _, channel := range Channels[route] {
		dataMap := make(map[string]interface{}, 0)
		dataMap["handler"] = route
		dataMap["data"] = data
		channel.Connection.WriteJSON(dataMap)
	}
}

func (h *Handler) Test(c echo.Context) error {
	cc := c.(*CustomContext)
	test := make(map[string]interface{}, 0)
	test["data"] = "1111111"
	PushToUsers(test, "a")
	cc.Respond(200, "success")
	return nil
}
