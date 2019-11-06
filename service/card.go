package service

import (
	"connect/model"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func CreateCard(card model.Card) (model.Card, error) {
	card.ID = bson.NewObjectId()
	card.CreatedAt = time.Now()
	err := model.CreateCard(card)
	if err != nil {
		return card, err
	}
	list, err := model.GetSingleList(card.ListID)
	if err != nil {
		return card, err
	}
	newSort := append(list.Sort, card.ID)
	err = model.UpdateListByID(list.ID, bson.M{"$set": bson.M{"sort": newSort}})
	if err != nil {
		return card, err
	}
	return card, nil
}

func GetAllCards() ([]map[string]interface{}, error) {
	lists, err := model.GetAllList()
	if err != nil {
		return nil, err
	}
	cards, err := model.GetAllCards()
	if err != nil {
		return nil, err
	}
	cardMap := make(map[bson.ObjectId]model.Card, 0)
	for _, card := range cards {
		cardMap[card.ID] = card
	}
	out := make([]map[string]interface{}, 0)
	for _, list := range lists {
		listInfo := make(map[string]interface{}, 0)
		listInfo["id"] = list.ID
		listInfo["title"] = list.Title
		cardsInList := []model.Card{}
		for _, cardID := range list.Sort {
			cardsInList = append(cardsInList, cardMap[cardID])
		}
		listInfo["cards"] = cardsInList
		out = append(out, listInfo)
	}
	return out, nil
}
