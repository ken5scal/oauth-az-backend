package infrastructure

import (
	"database/sql"
	"github.com/ken5scal/oauth-az/domain"
)

type tokenRepositoryImpl struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) domain.TokenRepository {
	return &tokenRepositoryImpl{
		db: db,
	}
}

func (t *tokenRepositoryImpl) GetByID(tokenID string) (*domain.Token, error) {
	return nil, nil
}

func (t *tokenRepositoryImpl) Insert(token *domain.Token) error {
	// Use ORM to insert
	return nil
}

func (t *tokenRepositoryImpl) Update(token *domain.Token) error {
	return nil
}

func (t *tokenRepositoryImpl) Delete(token *domain.Token) error {
	return nil
}
