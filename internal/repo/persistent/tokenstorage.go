package persistent

import (
	"interactive_learning/internal/repo"
	"interactive_learning/internal/utils/pair"
	"interactive_learning/internal/utils/tokengenerator"
	"sync"
	"time"
)

type TokenStorage struct {
	userIdToToken map[int]pair.Pair[tokengenerator.Token, time.Time]
	tokenToUserId map[tokengenerator.Token]int
	m             sync.Mutex
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
	_, ok := t.tokenToUserId[token]
	for ok {
		token = tokengenerator.GenerateToken()
		_, ok = t.tokenToUserId[token]
	}

	t.userIdToToken[id] =
		pair.Pair[tokengenerator.Token, time.Time]{
			First:  token,
			Second: time.Now(),
		}
	t.tokenToUserId[token] = id

	return token
}

func (t *TokenStorage) DeleteTokenToUser(id int) error {
	t.m.Lock()
	defer t.m.Unlock()

	token, ok := t.userIdToToken[id]

	if !ok {
		return repo.NoSuchRecordToDelete
	}

	delete(t.userIdToToken, id)
	delete(t.tokenToUserId, token.First)

	return nil
}

func (t *TokenStorage) IsValidToken(token tokengenerator.Token) (int, error) {
	t.m.Lock()
	defer t.m.Unlock()

	userId, ok := t.tokenToUserId[token]
	if !ok {
		return -1, repo.InvalidToken
	} else if time.Until(t.userIdToToken[userId].Second).Hours() < -1 {
		return -1, repo.ExpiredToken
	}

	return userId, nil
}
