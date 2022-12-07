package Go_ORM

import (
	"context"
	"database/sql"
)

// Querier 用于 SELECT 语句
type Querier[T any] interface {
	Get(ctx context.Context) (*T, error)
	// GetMulti 批量查询
	GetMulti(ctx context.Context) ([]*T, error)
}

// Executor 用于 INSERT, DELETE 和 UPDATE
type Executor interface {
	Exec(ctx context.Context) (sql.Result, error)
}

// QueryBuilder 代表 SQL构造 的过程
type QueryBuilder interface {
	Build() (*Query, error)
}

type Query struct {
	SQL  string
	Args []any
}
