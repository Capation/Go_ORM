package Go_ORM

import (
	"context"
	"reflect"
	"strings"
)

type Selector[T any] struct {
	table string
}

func (s *Selector[T]) Build() (*Query, error) {
	var sb strings.Builder
	sb.WriteString("SELECT * FROM ")
	// 我怎么把表名拿到 --> 反射
	// 如果用户指定了表名, 我们就用表名
	// 如果用户没有指定表名, 我们就用类型名
	// 决策: 如果用户指定了表名, 就直接使用, 不会使用反引号; 否则使用反引号括起来
	if s.table == "" {
		var t T
		typ := reflect.TypeOf(t)
		sb.WriteByte('`')
		sb.WriteString(typ.Name())
		sb.WriteByte('`')
	} else {
		sb.WriteString(s.table)
	}
	sb.WriteByte(';')
	return &Query{
		SQL: sb.String(),
	}, nil
}

func (s *Selector[T]) From(table string) *Selector[T] {
	s.table = table
	return s
}

func (s *Selector[T]) Get(ctx context.Context) (*T, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Selector[T]) GetMulti(ctx context.Context) ([]*T, error) {
	//TODO implement me
	panic("implement me")
}
