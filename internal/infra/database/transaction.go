package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

type txkey struct{}

var TXKey = txkey{}

func SetTXContext(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TXKey, tx)
}

func GetTXFromContext(ctx context.Context, defaultConn DBTX) DBTX {
	tx, ok := ctx.Value(TXKey).(pgx.Tx)
	if ok {
		return tx
	}
	return defaultConn
}

type TransactorTX interface {
	WithTransaction(context.Context, func(context.Context) error) error
}

var _ TransactorTX = (*TRXManager)(nil)

type TRXManager struct {
	db  *pgxpool.Pool
	log *zerolog.Logger
}

func NewTRXManager(conn *pgxpool.Pool, log *zerolog.Logger) *TRXManager {
	return &TRXManager{
		db:  conn,
		log: log,
	}
}

func (txdb *TRXManager) WithTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := txdb.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("create transaction")
	}

	ctx = SetTXContext(ctx, tx)
	err = fn(ctx)
	if err != nil {
		if rbErr := tx.Rollback(ctx); rbErr != nil {
			return fmt.Errorf("rollback: %v : original %w", rbErr, err)
		}
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commit:%w", err)
	}
	return nil
}
