package persistent

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
)

type CardsRepo struct {
	psql repo.PSQL
}

func NewCardsRepo(psql repo.PSQL) *CardsRepo {
	return &CardsRepo{psql: psql}
}

func (cr *CardsRepo) GetCardsByModule(moduleId int) ([]entity.Card, error) {
	rows, err := cr.psql.Query("SELECT * FROM cards WHERE module_id = $1", moduleId)
	if err != nil {
		return []entity.Card{}, repo.NewDBError("cards", "select", err)
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
			return []entity.Card{}, repo.NewDBError("cards", "select", err)
		}
		cards = append(cards, c)
	}

	return cards, nil
}

func (cr *CardsRepo) GetCardById(cardId int) (entity.Card, error) {
	row := cr.psql.QueryRow("SELECT * FROM cards WHERE id = $1", cardId)
	c := entity.Card{}
	err := row.Scan(&c.Id,
		&c.ParentModule,
		&c.Term.Lang,
		&c.Term.Text,
		&c.Definition.Lang,
		&c.Definition.Text)
	if err != nil {
		return entity.Card{}, repo.NewDBError("cards", "select", err)
	}
	return c, nil
}

func (cr *CardsRepo) GetLastInsertedCardId() (int, error) {
	row := cr.psql.QueryRow("SELECT MAX(id) FROM cards")
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, repo.NewDBError("cards", "select", err)
	}
	return id, nil
}

func (cr *CardsRepo) GetParentModuleId(cardId int) (int, error) {
	row := cr.psql.QueryRow("SELECT module_id FROM cards WHERE id = $1", cardId)
	var id int
	err := row.Scan(&id)
	if err != nil {
		return -1, repo.NewDBError("cards", "select", err)
	}
	return id, nil
}

func (cr *CardsRepo) InsertCard(card entity.Card) error {
	result, err := cr.psql.Exec("INSERT INTO cards(module_id, term_lang, term_text, def_lang, def_text) "+
		"VALUES($1, $2, $3, $4, $5)", card.ParentModule, card.Term.Lang, card.Term.Text, card.Definition.Lang, card.Definition.Text)
	if err != nil {
		return repo.NewDBError("cards", "insert", err)
	} else if count, _ := result.RowsAffected(); count == 0 {
		return repo.InsertRecordError
	}
	return nil
}

func (cr *CardsRepo) UpdateCard(card entity.Card) error {
	result, err := cr.psql.Exec("UPDATE cards "+
		"SET term_lang = $1, term_text = $2, def_lang = $3, def_text = $4 "+
		"WHERE id = $5", card.Term.Lang, card.Term.Text, card.Definition.Lang, card.Definition.Text, card.Id)
	if err != nil {
		return repo.NewDBError("cards", "update", err)
	} else if count, _ := result.RowsAffected(); count == 0 {
		return repo.NoSuchRecordToUpdate
	}
	return nil
}

func (cr *CardsRepo) DeleteCard(cardId int) error {
	result, err := cr.psql.Exec("DELETE FROM cards WHERE id = $1", cardId)
	if err != nil {
		return repo.NewDBError("cards", "delete", err)
	} else if count, _ := result.RowsAffected(); count < 1 {
		return repo.NoSuchRecordToDelete
	}
	return nil
}

func (cr *CardsRepo) DeleteCardsToParentModule(moduleId int) error {
	_, err := cr.psql.Exec("DELETE FROM cards WHERE module_id = $1", moduleId)
	if err != nil {
		return repo.NewDBError("cards", "delete", err)
	}
	return nil
}
