package sqlex

import (
	"database/sql"
	"errors"
	"fmt"
)

// WithTransaction runs a function within a transaction.
// If the process returns an error, it will be rolled back.
func WithTransaction(db *sql.DB, block func(tx *sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		rbErr := tx.Rollback()
		if rbErr == nil || errors.Is(rbErr, sql.ErrTxDone) {
			return
		}
		if err == nil {
			err = rbErr
			return
		}
		err = fmt.Errorf("%v: %w", rbErr, err)
	}()
	err = block(tx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
