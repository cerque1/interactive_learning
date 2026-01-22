package interactivelearning

import (
	"errors"
	"interactive_learning/internal/entity"
	"interactive_learning/internal/uow"
	"interactive_learning/internal/usecase"
)

func (u *UseCase) getCardOwnerId(cardId int, uow uow.UnitOfWork) (int, error) {
	cardsRepoRead := u.cardsRepoRead
	moduleRepoRead := u.moduleRepoRead
	if uow != nil {
		cardsRepoRead = uow.GetCardRepoReader()
		moduleRepoRead = uow.GetModuleRepoReader()
	}

	parentModule, err := cardsRepoRead.GetParentModuleId(cardId)
	if err != nil {
		return -1, u.errorsMapper.DBErrorToApp(err)
	}
	moduleOwner, err := moduleRepoRead.GetModuleOwnerId(parentModule)
	if err != nil {
		return -1, u.errorsMapper.DBErrorToApp(err)
	}
	return moduleOwner, nil
}

func (u *UseCase) GetCardById(cardId int) (entity.Card, error) {
	card, err := u.cardsRepoRead.GetCardById(cardId)
	if err != nil {
		return entity.Card{}, u.errorsMapper.DBErrorToApp(err)
	}
	return card, nil
}

func (u *UseCase) GetCardsByModule(moduleId int, userId int) ([]entity.Card, error) {
	_, err := u.GetModuleById(moduleId, userId)
	if err != nil { // нужно для проверки доступен ли модуль
		return []entity.Card{}, err
	}
	cards, err := u.cardsRepoRead.GetCardsByModule(moduleId)
	if err != nil {
		return []entity.Card{}, u.errorsMapper.DBErrorToApp(err)
	}
	return cards, nil
}

func (u *UseCase) InsertCard(card entity.Card) (int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return -1, usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	err := uow.GetCardRepoWriter().InsertCard(card)
	if err != nil {
		return -1, u.errorsMapper.DBErrorToApp(err)
	}

	insertedId, err := uow.GetCardRepoReader().GetLastInsertedCardId()
	if err != nil {
		return -1, u.errorsMapper.DBErrorToApp(err)
	}

	if err = uow.Commit(); err != nil {
		return -1, usecase.NewInternalError(err)
	}

	return insertedId, nil
}

func (u *UseCase) InsertCards(cards entity.CardsToAdd) ([]int, error) {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return []int{}, usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	var ids []int
	var err error
	var curId int

	for _, card := range cards.Cards {
		err = uow.GetCardRepoWriter().InsertCard(entity.Card{ParentModule: cards.ParentModule, Term: card.Term, Definition: card.Definition})
		if err != nil {
			return []int{}, u.errorsMapper.DBErrorToApp(err)
		}
		curId, err = uow.GetCardRepoReader().GetLastInsertedCardId()
		if err != nil {
			return []int{}, u.errorsMapper.DBErrorToApp(err)
		}
		ids = append(ids, curId)
	}

	if err = uow.Commit(); err != nil {
		return []int{}, usecase.NewInternalError(err)
	}

	return ids, nil
}

func (u *UseCase) UpdateCard(userId int, card entity.Card) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	ownerId, err := u.getCardOwnerId(card.Id, uow)
	if err != nil {
		return err
	} else if ownerId != userId {
		return usecase.NewNotAvailableError("card", card.Id)
	}

	err = uow.GetCardRepoWriter().UpdateCard(card)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	if err = uow.Commit(); err != nil {
		return usecase.NewInternalError(err)
	}
	return nil
}

func (u *UseCase) DeleteCard(userId int, cardId int) error {
	uow := u.unitOfWorkFactory()
	if err := uow.Begin(); err != nil {
		return usecase.NewInternalError(err)
	}
	defer uow.Rollback()

	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	ownerId, err := u.getCardOwnerId(cardId, uow)
	if err != nil {
		return err
	} else if ownerId != userId {
		return usecase.NewNotAvailableError("card", cardId)
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	err = uow.GetCardsResultsRepoWriter().DeleteResultsToCard(cardId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	err = uow.GetCardRepoWriter().DeleteCard(cardId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	if err = uow.Commit(); err != nil {
		return usecase.NewInternalError(err)
	}
	return nil
}

func (u *UseCase) deleteCardsToParentModule(moduleId int, uow uow.UnitOfWork) error {
	if uow == nil {
		return usecase.NewInternalError(errors.New("uow is null"))
	}

	module, err := uow.GetModuleRepoReader().GetModuleById(moduleId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	u.cardsResultsMutex.Lock()
	defer u.cardsResultsMutex.Unlock()

	for _, card := range module.Cards {
		err = uow.GetCardsResultsRepoWriter().DeleteResultsToCard(card.Id)
		if err != nil {
			return u.errorsMapper.DBErrorToApp(err)
		}
	}

	u.cardMutex.Lock()
	defer u.cardMutex.Unlock()

	err = uow.GetCardRepoWriter().DeleteCardsToParentModule(moduleId)
	if err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}

	return nil
}
