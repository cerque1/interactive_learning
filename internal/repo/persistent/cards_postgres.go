package persistent

import (
	"database/sql"
	"errors"
	"interactive_learning/internal/entity"
)

type CardsRepo struct {
	db *sql.DB
}

func NewCardsRepo(db *sql.DB) *CardsRepo {
	return &CardsRepo{db: db}
}

func (cr *CardsRepo) GetCardsByModule(moduleId int) ([]entity.Card, error) {
	rows, err := cr.db.Query("SELECT * FROM cards WHERE module_id = $1", moduleId)
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

func (cr *CardsRepo) GetCardById(cardId int) (entity.Card, error) {
	row := cr.db.QueryRow("SELECT * FROM cards WHERE id = $1", cardId)
	c := entity.Card{}
	err := row.Scan(&c.Id,
		&c.ParentModule,
		&c.Term.Lang,
		&c.Term.Text,
		&c.Definition.Lang,
		&c.Definition.Text)
	if err != nil {
		return entity.Card{}, err
	}
	return c, nil
}

func (cr *CardsRepo) GetLastInsertedCardId() (int, error) {
	row := cr.db.QueryRow("SELECT MAX(id) FROM cards")
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (cr *CardsRepo) GetParentModuleId(cardId int) (int, error) {
	row := cr.db.QueryRow("SELECT module_id FROM cards WHERE id = $1", cardId)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (cr *CardsRepo) InsertCard(card entity.Card) error {
	result, err := cr.db.Exec("INSERT INTO cards(module_id, term_lang, term_text, def_lang, def_text) "+
		"VALUES($1, $2, $3, $4, $5)", card.ParentModule, card.Term.Lang, card.Term.Text, card.Definition.Lang, card.Definition.Text)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("insert card error")
	}
	return nil
}

func (cr *CardsRepo) UpdateCard(card entity.Card) error {
	result, err := cr.db.Exec("UPDATE cards "+
		"SET term_lang = $1, term_text = $2, def_lang = $3, def_text = $4 "+
		"WHERE id = $5", card.Term.Lang, card.Term.Text, card.Definition.Lang, card.Definition.Text, card.Id)
	if err != nil {
		return err
	}
	if count, _ := result.RowsAffected(); count == 0 {
		return errors.New("update card error")
	}
	return nil
}

func (cr *CardsRepo) DeleteCard(cardId int) error {
	_, err := cr.db.Exec("DELETE FROM cards WHERE id = $1", cardId)
	if err != nil {
		return err
	}
	return nil
}

func (cr *CardsRepo) DeleteCardsToParentModule(moduleId int) error {
	_, err := cr.db.Exec("DELETE FROM cards WHERE module_id = $1", moduleId)
	if err != nil {
		return err
	}
	return nil
}
