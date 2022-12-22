package Go_ORM

import (
	"Go_ORM/internal/errs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_parseMode(t *testing.T) {
	testCases := []struct {
		name   string
		entity any

		wantModel *model
		wantErr   error
	}{
		{
			name:    "test model",
			entity:  TestModel{},
			wantErr: errs.ErrPointerOnly,
		},

		{
			name:   "pointer",
			entity: &TestModel{},
			wantModel: &model{
				tableName: "test_model",
				fileMap: map[string]*field{
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

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := parseMode(tc.entity)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantModel, m)
		})
	}
}
