package repository

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pepshot/SoftPlace/shared/logger"
)

type TxManager struct {
	db     *pgxpool.Pool
	logger *logger.Logger
}

func NewTxManager(db *pgxpool.Pool, logger *logger.Logger) *TxManager {
	return &TxManager{
		db:     db,
		logger: logger,
	}
}

func (m *TxManager) RunInTx(ctx context.Context, fn func(ctx context.Context, tx DBExecutor) error) error {
	tx, err := m.db.Begin(ctx)
	if err != nil {
		m.logger.Error("failed to begin transaction", "error", err)
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		}
	}()

	if err := fn(ctx, tx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			m.logger.Error("failed to rollback transaction", "error", rollbackErr)
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		m.logger.Error("failed to commit transaction", "error", err)
		return err
	}

	return nil
}
