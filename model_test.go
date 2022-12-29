package Go_ORM

import (
	"Go_ORM/internal/errs"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_Register(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantModel *Model
		wantErr   error

		opts []ModelOpt
	}{
		{
			name:    "test Model",
			entity:  TestModel{},
			wantErr: errs.ErrPointerOnly,
		},

		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				tableName: "test_model",
				fileMap: map[string]*Field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
		},

		{
			name:    "map",
			entity:  map[string]string{},
			wantErr: errs.ErrPointerOnly,
		},

		{
			name:    "slice",
			entity:  []int{},
			wantErr: errs.ErrPointerOnly,
		},

		{
			name:    "basic types",
			entity:  0,
			wantErr: errs.ErrPointerOnly,
		},
	}

	r := &registry{}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Register(tc.entity, tc.opts...)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}

func TestRegistry_Get(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantModel *Model
		wantErr   error
	}{
		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &Model{
				tableName: "test_model",
				fileMap: map[string]*Field{
					"Id": {
						colName: "id",
					},
					"FirstName": {
						colName: "first_name",
					},
					"LastName": {
						colName: "last_name",
					},
					"Age": {
						colName: "age",
					},
				},
			},
		},

		{
			name: "tag",
			// 局部匿名方法
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column=first_name_t"`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				tableName: "tag_table",
				fileMap: map[string]*Field{
					"FirstName": {
						colName: "first_name_t",
					},
				},
			},
		},

		{
			name: "empty column",
			// 局部匿名方法
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column="`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				tableName: "tag_table",
				fileMap: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},

		{
			name: "column only",
			// 局部匿名方法
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"column"`
				}
				return &TagTable{}
			}(),
			wantErr: errs.NewErrInvalidTagContent("column"),
		},

		{
			name: "ignore tag",
			// 局部匿名方法
			entity: func() any {
				type TagTable struct {
					FirstName string `orm:"abc=abc"`
				}
				return &TagTable{}
			}(),
			wantModel: &Model{
				tableName: "tag_table",
				fileMap: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},

		{
			name:   "table name",
			entity: &CustomTableName{},
			wantModel: &Model{
				tableName: "custom_table_name_t",
				fileMap: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},

		{
			name:   "table name ptr",
			entity: &CustomTableNamePtr{},
			wantModel: &Model{
				tableName: "custom_table_name_ptr_t",
				fileMap: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},

		{
			name:   "empty table name",
			entity: &EmptyTableName{},
			wantModel: &Model{
				tableName: "empty_table_name",
				fileMap: map[string]*Field{
					"FirstName": {
						colName: "first_name",
					},
				},
			},
		},
	}

	r := NewRegistry()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := r.Get(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)

			typ := reflect.TypeOf(tc.entity)
			cache, ok := r.models.Load(typ)
			assert.True(t, ok)
			assert.Equal(t, tc.wantModel, cache)
		})
	}
}

// Go里面 结构体实现接口 与 结构体指针实现接口 两者是不等价的
type CustomTableName struct {
	FirstName string
}

func (c CustomTableName) TableName() string {
	return "custom_table_name_t"
}

type CustomTableNamePtr struct {
	FirstName string
}

func (c *CustomTableNamePtr) TableName() string {
	return "custom_table_name_ptr_t"
}

type EmptyTableName struct {
	FirstName string
}

func (c *EmptyTableName) TableName() string {
	return ""
}

func TestModelWithTableName(t *testing.T) {
	r := NewRegistry()
	m, err := r.Register(&TestModel{}, ModelWithTableName("test_model_ttt"))
	assert.NoError(t, err)
	assert.Equal(t, "test_model_ttt", m.tableName)
}

func TestModelWithTableName1(t *testing.T) {
	r := NewRegistry()
	m, err := r.Register(&TestModel{}, ModelWithTableName(""))
	assert.NoError(t, err)
	assert.Equal(t, "", m.tableName)
}

func TestModelWithColumnName(t *testing.T) {
	testCases := []struct {
		name    string
		field   string
		colName string

		wantColName string
		wantErr     error
	}{
		{
			name:    "column name",
			field:   "FirstName",
			colName: "first_name_ccc",

			wantColName: "first_name_ccc",
		},

		{
			name:    "invalid column name",
			field:   "xxx",
			colName: "first_name_ccc",

			wantErr: errs.NewErrUnknownField("xxx"),
		},

		{
			name:    "empty column name",
			field:   "FirstName",
			colName: "",

			wantColName: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r := NewRegistry()
			m, err := r.Register(&TestModel{}, ModelWithColumnName(tc.field, tc.colName))
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			fd, ok := m.fileMap[tc.field]
			assert.True(t, ok)
			assert.Equal(t, tc.wantColName, fd.colName)
		})
	}
}
