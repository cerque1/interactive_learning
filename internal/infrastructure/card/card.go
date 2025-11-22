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
	id_str := c.Param("id")
	id, err := strconv.Atoi(id_str)
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

func (cr *CardRoutes) InsertCard(c echo.Context) error {
	card := entity.Card{}
	if err := c.Bind(card); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"message": "bad card info",
		})
	}

	id, err := cr.CardUC.InsertCard(card)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"new_id": id,
	})
}
