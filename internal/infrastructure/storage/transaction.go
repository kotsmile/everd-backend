package storage

import (
	"context"
	"database/sql"
	"errors"

	"github.com/kotsmile/everd-backend/internal/util"
)

var ErrInvalidTransaction = errors.New("invalid transaction type")

type SQLTransaction struct {
	tx *sql.Tx
}

type SQLTransactionFactory struct {
	db *sql.DB
}

func (f *SQLTransactionFactory) WithTransaction(ctx context.Context, fn func(util.Transaction) error) error {
	tx, err := f.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(&SQLTransaction{tx: tx}); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}

		return err
	}

	return tx.Commit()
}

func GetDBExecutor(tx util.Transaction, db *sql.DB) (Executor, error) {
	if tx == nil {
		return db, nil
	}

	sqlTx, ok := tx.(*SQLTransaction)
	if !ok {
		return nil, ErrInvalidTransaction
	}

	return sqlTx.tx, nil
}

func GetTxOrCreateTx(
	ctx context.Context,
	tx util.Transaction,
	db *sql.DB,
) (Executor, func() error, func() error, error) {
	nilFn := func() error { return nil }

	if tx != nil {
		sqlTx, ok := tx.(*SQLTransaction)
		if !ok {
			return nil, nilFn, nilFn, ErrInvalidTransaction
		}

		return sqlTx.tx, nilFn, nilFn, nil
	}

	sqlTx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nilFn, nilFn, err
	}

	return sqlTx, sqlTx.Commit, sqlTx.Rollback, nil
}

type Executor interface {
	QueryRow(query string, args ...any) *sql.Row
	Query(query string, args ...any) (*sql.Rows, error)
	Exec(query string, args ...any) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
}
