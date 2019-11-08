package handler

import (
	"connect/model"
	"connect/service"

	"github.com/labstack/echo"
	"gopkg.in/mgo.v2/bson"
)

func (h *Handler) CreateCard(c echo.Context) error {
	cc := c.(*CustomContext)
	card := model.Card{}
	err := cc.Bind(&card)
	if err != nil {
		return cc.RespondError(400, err.Error())
	}
	cardCreated, err := service.CreateCard(card)
	if err != nil {
		return cc.RespondError(500, err.Error())
	}
	PushToUsers(cardCreated, "CardNew")
	return cc.Respond(200, cardCreated)
}

func (h *Handler) CreateList(c echo.Context) error {
	cc := c.(*CustomContext)
	list := model.List{}
	err := cc.Bind(&list)
	if err != nil {
		return cc.RespondError(400, err.Error())
	}
	list.ID = bson.NewObjectId()
	err = model.CreateList(list)
	if err != nil {
		return cc.RespondError(500, err.Error())
	}
	PushToUsers(list, "ListNew")
	return cc.Respond(200, list)
}

func (h *Handler) GetCards(c echo.Context) error {
	cc := c.(*CustomContext)
	cards, err := service.GetAllCards()
	if err != nil {
		return cc.RespondError(500, err.Error())
	}
	return cc.Respond(200, cards)
}

func (h *Handler) UpdateList(c echo.Context) error {
	cc := c.(*CustomContext)
	req := model.List{}
	err := cc.Bind(&req)
	if err != nil {
		return cc.RespondError(400, err.Error())
	}
	err = model.UpdateListByID(req.ID, req)
	if err != nil {
		return cc.RespondError(500, err.Error())
	}
	PushToUsers(req, "ListUpdate")
	return cc.Respond(200, req)
}

func (h *Handler) UpdateCardSort(c echo.Context) error {
	cc := c.(*CustomContext)
	req := service.UpdateCardSortReq{}
	err := cc.Bind(&req)
	if err != nil {
		return cc.RespondError(400, err.Error())
	}
	resp, err := service.UpdateCardSort(req)
	if err != nil {
		return cc.RespondError(500, err.Error())
	}
	PushToUsers(resp, "CardSortUpdate")
	return cc.Respond(200, "ok")
}
