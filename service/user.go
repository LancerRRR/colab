package service

import (
	"connect/http/response"
	"connect/model"
	"connect/server"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Login(userName, password string) (*response.UserLoginResp, error) {
	user, err := model.GetUserByQuery(bson.M{"email": userName}, nil)
	if err != nil {
		return nil, err
	}
	if !CheckPasswordHash(password, user.Password) {
		return nil, errors.New("Invalid password")
	}
	out := &response.UserLoginResp{}
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = user.ID
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	t, err := token.SignedString([]byte("shawn"))
	if err != nil {
		return nil, err
	}
	_, err = server.Cl.Set(user.ID.Hex(), t, time.Hour*72).Result()
	if err != nil {
		return nil, err
	}
	out.Token = t
	out.UserID = user.ID
	return out, nil
}
