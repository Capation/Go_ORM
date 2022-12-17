package reflect

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIterateFields(t *testing.T) {

	type User struct {
		Name string
		age  int
	}

	testCases := []struct {
		name    string
		entity  any
		wantErr error
		wantRes map[string]any
	}{
		{
			name: "struct",
			entity: User{
				Name: "Tom",
				age:  18,
			},
			wantRes: map[string]any{
				"Name": "Tom",
				// age 是私有的，拿不到，最终我们创建了 0 值来填充
				"age": 0,
			},
		},

		{
			name: "pointer",
			entity: &User{
				Name: "Tom",
				age:  18,
			},
			wantRes: map[string]any{
				"Name": "Tom",
				// age 是私有的，拿不到，最终我们创建了 0 值来填充
				"age": 0,
			},
		},

		{
			name:    "basic type",
			entity:  18,
			wantErr: errors.New("不支持的数据类型"),
		},

		{
			// 多级指针
			name: "multiple pointer",
			// 局部匿名函数
			entity: func() **User {
				// res 是一级指针
				res := &User{
					Name: "Tom",
					age:  18,
				}
				return &res
			}(),
			wantRes: map[string]any{
				"Name": "Tom",
				// age 是私有的，拿不到，最终我们创建了 0 值来填充
				"age": 0,
			},
		},

		{
			name:    "nil entity",
			entity:  nil,
			wantErr: errors.New("不支持 nil"),
		},

		{
			name:    "user nil",
			entity:  (*User)(nil),
			wantErr: errors.New("不支持零值"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			res, err := IterateFields(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantRes, res)
		})
	}
}
