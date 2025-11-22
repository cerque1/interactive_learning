package interactivelearning

import (
	"interactive_learning/internal/entity"
	"interactive_learning/internal/repo"
	"interactive_learning/internal/repo/persistent"
	"interactive_learning/internal/utils/tokengenerator"
)

type UseCase struct {
	userRepo     repo.UsersRepo
	tokenStorage repo.TokenStorage
	cardRepo     repo.CardRepo
	moduleRepo   repo.ModuleRepo
}

func New(userRepo repo.UsersRepo) *UseCase {
	return &UseCase{userRepo: userRepo, tokenStorage: persistent.NewTokenStorage()}
}

func (u *UseCase) AddTokenToUser(id int) tokengenerator.Token {
	return u.tokenStorage.AddTokenToUser(id)
}

func (u *UseCase) DeleteTokenToUser(id int) error {
	return u.tokenStorage.DeleteTokenToUser(id)
}

func (u *UseCase) IsValidToken(token tokengenerator.Token) (int, error) {
	return u.tokenStorage.IsValidToken(token)
}

func (u *UseCase) GetUserByLogin(login string) (entity.User, error) {
	return u.userRepo.GetUserByLogin(login)
}

func (u *UseCase) InsertUser(user entity.User) (int, error) {
	err := u.userRepo.InsertUser(user)
	if err != nil {
		return -1, err
	}
	new_user, err := u.userRepo.GetUserByLogin(user.Login)
	if err != nil {
		return -1, err
	}
	return new_user.Id, nil
}

func (u *UseCase) GetCardsByModule(module_id int) ([]entity.Card, error) {
	return u.cardRepo.GetCardsByModule(module_id)
}

func (u *UseCase) InsertCard(card entity.Card) (int, error) {
	err := u.cardRepo.InsertCard(card)
	if err != nil {
		return -1, err
	}
	return u.cardRepo.GetLastInsertedCardId()
}

func (u *UseCase) InsertCards(cards []entity.Card) ([]int, error) {
	var ids []int
	var err error
	var cur_id int

	for _, card := range cards {
		err = u.cardRepo.InsertCard(card)
		if err != nil {
			return []int{}, err
		}
		cur_id, err = u.cardRepo.GetLastInsertedCardId()
		if err != nil {
			return []int{}, err
		}
		ids = append(ids, cur_id)
	}
	return ids, nil
}

func (u *UseCase) DeleteCard(card_id int) error {
	return u.cardRepo.DeleteCard(card_id)
}

func (u *UseCase) GetModulesByUser(user_id int) ([]entity.Module, error) {
	return u.moduleRepo.GetModulesByUser(user_id)
}

func (u *UseCase) GetModulesWithCardsByUser(user_id int) ([]entity.Module, error) {
	modules, err := u.moduleRepo.GetModulesByUser(user_id)
	if err != nil {
		return []entity.Module{}, err
	}

	for _, module := range modules {
		cards, err := u.cardRepo.GetCardsByModule(user_id)
		if err != nil {
			return []entity.Module{}, err
		}
		module.Cards = cards
	}

	return modules, nil
}

func (u *UseCase) InsertModule(module entity.Module) (int, []int, error) {
	err := u.moduleRepo.InsertModule(module)
	if err != nil {
		return -1, []int{}, err
	}
	insert_ids, err := u.InsertCards(module.Cards)
	if err != nil {
		return -1, []int{}, err
	}
	id, err := u.moduleRepo.GetLastInsertedModuleId()
	if err != nil {
		return -1, []int{}, err
	}
	return id, insert_ids, nil
}

func (u *UseCase) DeleteModule(module_id int) error {
	return u.moduleRepo.DeleteModule(module_id)
}
