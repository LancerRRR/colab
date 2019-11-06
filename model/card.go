package model

import (
	"connect/server"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	CollectionList *mgo.Collection
	CollectionCard *mgo.Collection
)

func InitCollections() {
	CollectionList = server.Mongodb.DB("Card").C("List")
	CollectionCard = server.Mongodb.DB("Card").C("Card")
}

type List struct {
	ID    bson.ObjectId   `json:"id" bson:"_id"`
	Title string          `json:"title" bson:"title"`
	Sort  []bson.ObjectId `json:"sort" bson:"sort"`
}

type Card struct {
	ID          bson.ObjectId  `json:"id" bson:"_id"`
	Title       string         `json:"title" bson:"title"`
	Description string         `json:"description" bson:"description"`
	Mindmap     string         `json:"mindmap" bson:"mindmap"`
	Functions   []FunctionInfo `json:"functions" bson:"functions"`
	List        string         `json:"list" bson:"list"`
	ListID      bson.ObjectId  `json:"listID" bson:"listID"`
	Members     []Member       `json:"members" bson:"members"`
	CreatedAt   time.Time      `json:"createdAt" bson:"createdAt"`
}

type FunctionInfo struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name" bson:"name"`
	Designed  bool          `json:"designed" bson:"designed"`
	Developed bool          `json:"developed" bson:"developed"`
}

type Member struct {
	ID     bson.ObjectId `json:"id" bson:"_id"`
	Name   string        `json:"name" bson:"name"`
	Avatar string        `json:"avatar" bson:"avatar"`
	Role   string        `json:"role" bson:"role"`
}

func CreateCard(card Card) error {
	return CollectionCard.Insert(card)
}

func GetAllCards() ([]Card, error) {
	cards := []Card{}
	err := CollectionCard.Find(nil).All(&cards)
	return cards, err
}

func CreateList(list List) error {
	return CollectionList.Insert(list)
}

func GetSingleList(id bson.ObjectId) (List, error) {
	list := List{}
	err := CollectionList.FindId(id).One(&list)
	return list, err
}

func GetAllList() ([]List, error) {
	lists := []List{}
	err := CollectionList.Find(nil).All(&lists)
	return lists, err
}

func UpdateListByID(id bson.ObjectId, update interface{}) error {
	return CollectionList.UpdateId(id, update)
}
