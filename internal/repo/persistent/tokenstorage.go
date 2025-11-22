package persistent

import (
	"errors"
	"interactive_learning/internal/utils/pair"
	"interactive_learning/internal/utils/tokengenerator"
	"sync"
	"time"
)

type TokenStorage struct {
	user_id_to_token map[int]pair.Pair[tokengenerator.Token, time.Time]
	token_to_user_id map[tokengenerator.Token]int
	m                sync.Mutex
}

func NewTokenStorage() *TokenStorage {
	return &TokenStorage{
		make(map[int]pair.Pair[tokengenerator.Token, time.Time]),
		make(map[tokengenerator.Token]int),
		sync.Mutex{},
	}
}

func (t *TokenStorage) AddTokenToUser(id int) tokengenerator.Token {
	t.m.Lock()
	defer t.m.Unlock()

	token := tokengenerator.GenerateToken()
	_, ok := t.token_to_user_id[token]
	for ok {
		token = tokengenerator.GenerateToken()
		_, ok = t.token_to_user_id[token]
	}

	t.user_id_to_token[id] =
		pair.Pair[tokengenerator.Token, time.Time]{
			First:  token,
			Second: time.Now(),
		}
	t.token_to_user_id[token] = id

	return token
}

func (t *TokenStorage) DeleteTokenToUser(id int) error {
	t.m.Lock()
	defer t.m.Unlock()

	token, ok := t.user_id_to_token[id]

	if !ok {
		return errors.New("no such user to delete token")
	}

	delete(t.user_id_to_token, id)
	delete(t.token_to_user_id, token.First)

	return nil
}

func (t *TokenStorage) IsValidToken(token tokengenerator.Token) (int, error) {
	t.m.Lock()
	defer t.m.Unlock()

	user_id, ok := t.token_to_user_id[token]
	if !ok {
		return -1, errors.New("invalid token")
	} else if time.Until(t.user_id_to_token[user_id].Second).Hours() < -1 {
		return -1, errors.New("token is expired")
	}

	return user_id, nil
}
