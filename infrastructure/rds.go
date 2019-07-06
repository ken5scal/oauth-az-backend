package infrastructure

import (
	"database/sql"
	"github.com/go-redis/redis"
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

type AuthorzInfoRepositoryImpl struct {
	db *sql.DB
	r  *redis.Conn // TODO 適当
	tx *sql.Tx
}

//NewAuthzInfoRepository
func NewAuthzInfoRepositoryImpl(db *sql.DB) *AuthorzInfoRepositoryImpl {
	return &AuthorzInfoRepositoryImpl{
		db: db,
	}
}

// Transact is a wrapper to handle transaction properly
// ref: https://stackoverflow.com/questions/16184238/database-sql-tx-detecting-commit-or-rollback
func Transact(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		// catch panic to make sure db rolls back
		if p := recover(); p != nil {
			// If we did not handle panics the transaction would be rolled back eventually.
			// A non-commited transaction gets rolled back by the database
			// when the client disconnects or when the transaction gets garbage collected.
			// However, waiting for the transaction to resolve on its own could cause other (undefined) issues.
			// So it's better to resolve it as quickly as possible.
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// We don;t return err because it will override existing err
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}
