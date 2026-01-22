package interactivelearning

import "interactive_learning/internal/utils/tokengenerator"

func (u *UseCase) AddTokenToUser(id int) tokengenerator.Token {
	return u.tokenStorage.AddTokenToUser(id)
}

func (u *UseCase) DeleteTokenToUser(id int) error {
	if err := u.tokenStorage.DeleteTokenToUser(id); err != nil {
		return u.errorsMapper.DBErrorToApp(err)
	}
	return nil
}

func (u *UseCase) IsValidToken(token tokengenerator.Token) (int, error) {
	userId, err := u.tokenStorage.IsValidToken(token)
	if err != nil {
		return -1, u.errorsMapper.DBErrorToApp(err)
	}
	return userId, nil
}
