package service

import (
	"connect/model"
	"errors"
	"time"

	"gopkg.in/mgo.v2/bson"
)

type UpdateCardSortReq struct {
	ID        bson.ObjectId `json:"id"`
	NewListID bson.ObjectId `json:"newListID"`
	Sort      int           `json:"sort"`
}

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
	err = model.UpdateListByID(card.ListID, bson.M{"$set": bson.M{"sort": newSort}})
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
		for _, id := range list.Sort {
			cardsInList = append(cardsInList, cardMap[id])
		}
		listInfo["cards"] = cardsInList
		out = append(out, listInfo)
	}
	return out, nil
}

func UpdateCardSort(data UpdateCardSortReq) (map[string]interface{}, error) {
	card, err := model.GetCardByID(data.ID)
	if err != nil {
		return nil, err
	}
	oldList, err := model.GetSingleList(card.ListID)
	if err != nil {
		return nil, err
	}
	out := make(map[string]interface{}, 0)
	out["id"] = data.ID
	out["oldListID"] = oldList.ID
	out["newListID"] = data.NewListID
	out["newSort"] = data.Sort
	if data.NewListID == card.ListID {
		oldSort := int(0)
		for i, id := range oldList.Sort {
			if id == card.ID {
				oldSort = i
				break
			}
		}
		out["oldSort"] = oldSort
		newSort := []bson.ObjectId{}
		if oldSort < data.Sort {
			for i, id := range oldList.Sort {
				if i == oldSort {
					continue
				}
				if i == data.Sort {
					newSort = append(newSort, id)
					newSort = append(newSort, data.ID)
					continue
				}
				newSort = append(newSort, id)
			}
		} else if oldSort == data.Sort {
			return nil, errors.New("no change")
		} else {
			for i, id := range oldList.Sort {
				if i == oldSort {
					continue
				}
				if i == data.Sort {
					newSort = append(newSort, data.ID)
					newSort = append(newSort, id)
					continue
				}
				newSort = append(newSort, id)
			}
		}
		err = model.UpdateListByID(card.ListID, bson.M{"$set": bson.M{"sort": newSort}})
		if err != nil {
			return nil, err
		}
	} else {
		newSort1 := []bson.ObjectId{}
		for i, id := range oldList.Sort {
			if id != data.ID {
				newSort1 = append(newSort1, id)
			} else {
				out["oldSort"] = i
			}
		}
		err = model.UpdateListByID(card.ListID, bson.M{"$set": bson.M{"sort": newSort1}})
		if err != nil {
			return nil, err
		}
		newList, err := model.GetSingleList(data.NewListID)
		if err != nil {
			return nil, err
		}
		newSort2 := []bson.ObjectId{}
		for i, id := range newList.Sort {
			if i == data.Sort {
				newSort2 = append(newSort2, card.ID)
			}
			newSort2 = append(newSort2, id)
		}
		if data.Sort >= len(newList.Sort) {
			newSort2 = append(newSort2, card.ID)
		}
		err = model.UpdateListByID(newList.ID, bson.M{"$set": bson.M{"sort": newSort2}})
		if err != nil {
			return nil, err
		}
		err = model.UpdateCardByID(data.ID, bson.M{"$set": bson.M{"listID": newList.ID, "list": newList.Title}})
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}
