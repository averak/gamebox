package transaction

import (
	"context"
	"database/sql"
)

type (
	Transaction = *sql.Tx

	Option  func(*Options)
	Options struct {
		IsolationLevel sql.IsolationLevel
	}
)

// Connection は、データベースとの物理的な接続を表すインターフェースです。
type Connection interface {
	BeginRoTransaction(ctx context.Context, f func(ctx context.Context, tx Transaction) error, opts ...Option) error
	BeginRwTransaction(ctx context.Context, f func(ctx context.Context, tx Transaction) error, opts ...Option) error
	Close() error
}

func DefaultOptions() Options {
	return Options{
		IsolationLevel: sql.LevelReadCommitted,
	}
}

// WithIsolationLevel は、トランザクション分離レベルを設定するオプションです。
func WithIsolationLevel(l sql.IsolationLevel) Option {
	return func(o *Options) {
		o.IsolationLevel = l
	}
}
