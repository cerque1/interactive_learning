package card

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/usecase"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type CardRoutes struct {
	CardUC usecase.Cards
}

func NewCardRoutes(cardUc usecase.Cards) *CardRoutes {
	return &CardRoutes{CardUC: cardUc}
}

func (cr *CardRoutes) GetCardsByModule(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad id",
		})
	}

	cards, err := cr.CardUC.GetCardsByModule(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"cards": cards,
	})
}

func (cr *CardRoutes) GetCardById(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad id",
		})
	}

	card, err := cr.CardUC.GetCardById(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"card": card,
	})
}

func (cr *CardRoutes) InsertCards(c echo.Context) error {
	cards := entity.CardsToAdd{}
	if err := c.Bind(&cards); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": err.Error(),
		})
	}

	id, err := cr.CardUC.InsertCards(cards)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"new_ids": id,
	})
}

func (cr *CardRoutes) UpdateCard(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	idStr := c.Param("id")
	cardId, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad id",
		})
	}

	card := entity.Card{}
	if err = c.Bind(&card); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "error parse data: " + err.Error(),
		})
	}
	card.Id = cardId

	err = cr.CardUC.UpdateCard(userId, card)
	if err != nil {
		log.Println("update card error", err.Error())
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "update card error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}

func (cr *CardRoutes) DeleteCard(c echo.Context) error {
	userId, err := strconv.Atoi(c.QueryParam("user_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad user id",
		})
	}

	idStr := c.Param("id")
	cardId, err := strconv.Atoi(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad id",
		})
	}

	err = cr.CardUC.DeleteCard(userId, cardId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": "delete card error " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{})
}
