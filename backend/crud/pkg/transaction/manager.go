package transaction

import (
	"errors"

	"gorm.io/gorm"
)

func Begin(db *gorm.DB) (*Transaction, error) {
	if db == nil {
		return nil, errors.New("db is nul")
	}

	tx := db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &Transaction{tx: tx}, nil

}

func (t *Transaction) Rollback() error {
	if t.tx == nil {
		return errors.New("transaction is null")
	}

	return t.tx.Rollback().Error
}

func (t *Transaction) Commit() error {
	if t.tx == nil {
		return errors.New("transaction is null")
	}
	return t.tx.Commit().Error
}

func (t *Transaction) DB() *gorm.DB {
	return t.tx
}
