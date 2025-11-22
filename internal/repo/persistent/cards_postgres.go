package persistent

import (
	"database/sql"
	"errors"
	"interactive_learning/internal/entity"
)

type CardsRepoImpl struct {
	db *sql.DB
}

func (cr *CardsRepoImpl) GetCardsByModule(module_id int) ([]entity.Card, error) {
	rows, err := cr.db.Query("SELECT * FROM cards WHERE module_id = $1", module_id)
	if err != nil {
		return []entity.Card{}, err
	}
	defer rows.Close()

	cards := []entity.Card{}
	for rows.Next() {
		c := entity.Card{}
		err = rows.Scan(&c.Id,
			&c.ParentModule,
			&c.Term.Lang,
			&c.Term.Text,
			&c.Definition.Lang,
			&c.Definition.Text)
		if err != nil {
			return []entity.Card{}, err
		}
		cards = append(cards, c)
	}

	return cards, nil
}

func (cr *CardsRepoImpl) GetLastInsertedCardId() (int, error) {
	row := cr.db.QueryRow("SELECT MAX(id) FROM cards")
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, errors.ErrUnsupported
	}
	return id, nil
}

func (cr *CardsRepoImpl) InsertCard(card entity.Card) error {
	result, err := cr.db.Exec("INSERT INTO cards(module_id, term_lang, term_text, def_lang, def_text) " +
		"VALUES($1, $2, $3, $4, $5)")
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("insert card error")
	}
	return nil
}

func (cr *CardsRepoImpl) DeleteCard(card_id int) error {
	result, err := cr.db.Exec("DELETE FROM cards WHERE id = $1", card_id)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("delete card error")
	}
	return nil
}
