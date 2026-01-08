package persistent

import (
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
)

type CardsResultsRepo struct {
	psql repo.PSQL
}

func NewCardsResultsRepo(psql repo.PSQL) *CardsResultsRepo {
	return &CardsResultsRepo{psql: psql}
}

func (crr *CardsResultsRepo) GetCardsResultById(resultId int) ([]entity.CardsResult, error) {
	rows, err := crr.psql.Query("SELECT card_id, result FROM cards_results WHERE result_id = $1", resultId)
	if err != nil {
		return []entity.CardsResult{}, err
	}

	cards_results := []entity.CardsResult{}
	for rows.Next() {
		card_result := entity.CardsResult{}
		err := rows.Scan(&card_result.CardId,
			&card_result.Result)
		if err != nil {
			return []entity.CardsResult{}, err
		}

		cards_results = append(cards_results, card_result)
	}

	return cards_results, nil
}

func (crr *CardsResultsRepo) InsertCardResult(resultId, cardId int, resultStr string) error {
	result, err := crr.psql.Exec("INSERT INTO cards_results(result_id, card_id, result) "+
		"VALUES($1, $2, $3)", resultId, cardId, resultStr)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("insert card result error")
	}
	return nil
}

func (crr *CardsResultsRepo) DeleteCardResult(resultId, cardId int) error {
	_, err := crr.psql.Exec("DELETE FROM cards_results WHERE result_id = $1 AND card_id = $2", resultId, cardId)
	if err != nil {
		return err
	}
	return nil
}

func (crr *CardsResultsRepo) DeleteCardsToResult(resultId int) error {
	_, err := crr.psql.Exec("DELETE FROM cards_results WHERE result_id = $1", resultId)
	if err != nil {
		return err
	}
	return nil
}

func (crr *CardsResultsRepo) DeleteResultsToCard(cardId int) error {
	_, err := crr.psql.Exec("DELETE FROM cards_results WHERE card_id = $1", cardId)
	if err != nil {
		return err
	}
	return nil
}
