package card

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/usecase"
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
		"new_id": id,
	})
}
