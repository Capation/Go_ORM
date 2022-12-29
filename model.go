package Go_ORM

import (
	"Go_ORM/internal/errs"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

const (
	tagTestColumn = "column"
)

type model struct {
	// tableName 结构体对应的表名
	tableName string
	fileMap   map[string]*field
}

// field 字段
type field struct {
	colName string
}

//var models = map[reflect.Type]*model{}

// registry 代表元数据的注册中心
type registry struct {
	// 读写锁
	//lock   sync.RWMutex
	models sync.Map
}

// 全局变量
//var defaultRegistry = &registry{models: map[reflect.Type]*model{}}

func NewRegistry() *registry {
	return &registry{}
}

func (r *registry) Get(val any) (*model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*model), nil
	}
	m, err := r.parseMode(val)
	if err != nil {
		return nil, err
	}
	r.models.Store(typ, m)
	return m.(*model), nil
}

//func (r *registry) Get1(val any) (*model, error) {
//	typ := reflect.TypeOf(val)
//	// 读锁
//	r.lock.RLock()
//	m, ok := r.models[typ]
//	r.lock.RUnlock()
//	if ok {
//		return m, nil
//	}
//
//	// 加写锁，准备解析你的数据
//	r.lock.Lock()
//	defer r.lock.Unlock()
//
//	m, ok = r.models[typ]
//	if ok {
//		return m, nil
//	}
//
//	m, err := r.parseMode(val)
//	if err != nil {
//		return nil, err
//	}
//	r.models[typ] = m
//
//	return m, nil
//}

// 限制只能用一级指针
func (r *registry) parseMode(entity any) (*model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	typ = typ.Elem()
	numField := typ.NumField()
	fieldMap := make(map[string]*field, numField)
	for i := 0; i < numField; i++ {
		fd := typ.Field(i)
		pair, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		colName := pair[tagTestColumn]
		if colName == "" {
			// 用户没有设置，我们就给它转
			colName = underscoreName(fd.Name)
		}
		fieldMap[fd.Name] = &field{
			colName: colName,
		}
	}

	// 接口自定义表名
	var tableName string
	if tbl, ok := entity.(TableName); ok {
		tableName = tbl.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(typ.Name())
	}

	return &model{
		tableName: tableName,
		fileMap:   fieldMap,
	}, nil
}

func (r *registry) parseTag(tag reflect.StructTag) (map[string]string, error) {
	ormTag, ok := tag.Lookup("orm")
	if !ok {
		return map[string]string{}, nil
	}
	pairs := strings.Split(ormTag, ",")
	res := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		segs := strings.Split(pair, "=")
		if len(segs) != 2 {
			return nil, errs.NewErrInvalidTagContent(pair)
		}
		key := segs[0]
		val := segs[1]
		res[key] = val
	}
	return res, nil
}

// underscoreName 驼峰转字符串命名
func underscoreName(tableName string) string {
	var buf []byte
	for i, v := range tableName {
		if unicode.IsUpper(v) {
			if i != 0 {
				buf = append(buf, '_')
			}
			buf = append(buf, byte(unicode.ToLower(v)))
		} else {
			buf = append(buf, byte(v))
		}

	}
	return string(buf)
}
