package reflect

import (
	"errors"
	"reflect"
)

// IterateFields 遍历字段
func IterateFields(entity any) (map[string]any, error) {
	if entity == nil {
		return nil, errors.New("不支持 nil")
	}

	typ := reflect.TypeOf(entity)
	val := reflect.ValueOf(entity)

	if val.IsZero() {
		return nil, errors.New("不支持零值")
	}

	for typ.Kind() == reflect.Pointer {
		// 拿到指针指向的对象
		typ = typ.Elem()
		val = val.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, errors.New("不支持的数据类型")
	}

	// 看它有多少个字段
	numFields := typ.NumField()
	res := make(map[string]any, numFields)
	for i := 0; i < numFields; i++ {
		// 字段的类型
		fieldType := typ.Field(i)
		// 字段的值
		fieldValue := val.Field(i)

		if fieldType.IsExported() {
			res[fieldType.Name] = fieldValue.Interface()
		} else {
			res[fieldType.Name] = reflect.Zero(fieldType.Type).Interface()
		}
	}

	return res, nil
}
