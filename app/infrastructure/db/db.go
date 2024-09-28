package db

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/XSAM/otelsql"
	"github.com/averak/gamebox/app/core/config"
	"github.com/averak/gamebox/app/core/logger"
	"github.com/averak/gamebox/app/domain/repository/transaction"
	"github.com/averak/gamebox/app/infrastructure/trace"
	_ "github.com/lib/pq"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

var (
	conn *Connection
	once sync.Once
)

func NewConnection() (transaction.Connection, error) {
	var initErr error
	once.Do(func() {
		c, err := open()
		if err != nil {
			initErr = err
			return
		}
		conn = &c
	})
	if initErr != nil {
		return nil, initErr
	}
	return conn, nil
}

// open は DB に接続し、コネクションを返す。
// Connection は規定量のコネクションプールを保有するため、複数作成されないよう singleton を返す NewConnection を公開している。
func open() (Connection, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=UTC",
		config.Get().GetPostgres().GetHost(),
		config.Get().GetPostgres().GetUser(),
		config.Get().GetPostgres().GetPassword(),
		config.Get().GetPostgres().GetDatabase(),
		config.Get().GetPostgres().GetPort(),
	)
	db, err := otelsql.Open("postgres", dsn, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
	if err != nil {
		return Connection{}, err
	}
	err = db.Ping()
	if err != nil {
		return Connection{}, err
	}
	err = otelsql.RegisterDBStatsMetrics(db, otelsql.WithAttributes(semconv.DBSystemPostgreSQL))
	if err != nil {
		return Connection{}, err
	}
	return Connection{db}, nil
}

type Connection struct {
	db *sql.DB
}

func (t Connection) BeginRoTransaction(ctx context.Context, f func(ctx context.Context, tx transaction.Transaction) error, opts ...transaction.Option) error {
	ctx, span := trace.StartSpan(ctx, "db.BeginRoTransaction")
	defer span.End()

	options := transaction.DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	tx, err := t.db.BeginTx(
		ctx,
		&sql.TxOptions{
			Isolation: options.IsolationLevel,
			ReadOnly:  true,
		},
	)
	if err != nil {
		return fmt.Errorf("db.BeginTx: %w", err)
	}

	err = f(ctx, tx)
	if err != nil {
		logger.Debug(ctx, map[string]interface{}{
			"message": "Rollback transaction.",
			"error":   err.Error(),
		})
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %w", err, rbErr)
		}
		return fmt.Errorf("transaction error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}

func (t Connection) BeginRwTransaction(ctx context.Context, f func(ctx context.Context, tx transaction.Transaction) error, opts ...transaction.Option) error {
	ctx, span := trace.StartSpan(ctx, "db.BeginRwTransaction")
	defer span.End()

	options := transaction.DefaultOptions()
	for _, opt := range opts {
		opt(&options)
	}
	tx, err := t.db.BeginTx(
		ctx,
		&sql.TxOptions{
			Isolation: options.IsolationLevel,
			ReadOnly:  false,
		},
	)
	if err != nil {
		return fmt.Errorf("db.BeginTx: %w", err)
	}

	err = f(ctx, tx)
	if err != nil {
		logger.Debug(ctx, map[string]interface{}{
			"message": "Rollback transaction.",
			"error":   err.Error(),
		})
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("transaction error: %w, rollback error: %w", err, rbErr)
		}
		return fmt.Errorf("transaction error: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("tx.Commit: %w", err)
	}
	return nil
}

func (t Connection) Close() error {
	return t.db.Close()
}
