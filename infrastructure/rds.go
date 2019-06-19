package infrastructure

import (
	"database/sql"
	"github.com/ken5scal/oauth-az/domain"
)

type TokenRepositoryImpl struct {
	db *sql.DB
	tx *sql.Tx
}

func NewTokenRepository(db *sql.DB) *TokenRepositoryImpl {
	return &TokenRepositoryImpl{
		db: db,
	}
}

func (t *TokenRepositoryImpl) GetAccessTokenByID(tokenID string) (*domain.Token, error) {
	return nil, nil
}

func (t *TokenRepositoryImpl) Insert(token *domain.Token) error {
	// Use ORM to insert
	return nil
}

func (t *TokenRepositoryImpl) Update(token *domain.Token) error {
	return nil
}

func (t *TokenRepositoryImpl) Delete(token *domain.Token) error {
	return nil
}

func (t *TokenRepositoryImpl) BeginTransaction() (*sql.Tx, error) {
	if tx, err := t.db.Begin(); err != nil {
		return nil, err
	} else {
		t.tx = tx
		return tx, nil
	}
}

func (t *TokenRepositoryImpl) Rollback() error {
	return t.tx.Rollback()
}

func (t *TokenRepositoryImpl) Commit() error {
	return t.tx.Commit()
}

type RDSAuthorzInfoRepositoryImpl struct {
	db *sql.DB
	tx *sql.Tx
}

//NewAuthzInfoRepository
func NewRDSAuthzInfoRepositoryImpl(db *sql.DB) *RDSAuthorzInfoRepositoryImpl {
	return &RDSAuthorzInfoRepositoryImpl{
		db: db,
	}
}
