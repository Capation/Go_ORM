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

type Registry interface {
	Get(val any) (*Model, error)
	Register(val any, opts ...ModelOpt) (*Model, error)
}

type Model struct {
	// tableName 结构体对应的表名
	tableName string
	fileMap   map[string]*Field
}

type ModelOpt func(m *Model) error

// Field 字段
type Field struct {
	colName string
}

//var models = map[reflect.Type]*Model{}

// registry 代表元数据的注册中心
type registry struct {
	// 读写锁
	//lock   sync.RWMutex
	models sync.Map
}

// 全局变量
//var defaultRegistry = &registry{models: map[reflect.Type]*Model{}}

func NewRegistry() *registry {
	return &registry{}
}

func (r *registry) Get(val any) (*Model, error) {
	typ := reflect.TypeOf(val)
	m, ok := r.models.Load(typ)
	if ok {
		return m.(*Model), nil
	}
	m, err := r.Register(val)
	if err != nil {
		return nil, err
	}

	return m.(*Model), nil
}

//func (r *registry) Get1(val any) (*Model, error) {
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
//	m, err := r.Register(val)
//	if err != nil {
//		return nil, err
//	}
//	r.models[typ] = m
//
//	return m, nil
//}

// Register 限制只能用一级指针
func (r *registry) Register(entity any, opts ...ModelOpt) (*Model, error) {
	typ := reflect.TypeOf(entity)
	if typ.Kind() != reflect.Pointer || typ.Elem().Kind() != reflect.Struct {
		return nil, errs.ErrPointerOnly
	}
	elemTyp := typ.Elem()
	numField := elemTyp.NumField()
	fieldMap := make(map[string]*Field, numField)
	for i := 0; i < numField; i++ {
		fd := elemTyp.Field(i)
		pair, err := r.parseTag(fd.Tag)
		if err != nil {
			return nil, err
		}
		colName := pair[tagTestColumn]
		if colName == "" {
			// 用户没有设置，我们就给它转
			colName = underscoreName(fd.Name)
		}
		fieldMap[fd.Name] = &Field{
			colName: colName,
		}
	}

	// 接口自定义表名
	var tableName string
	if tbl, ok := entity.(TableName); ok {
		tableName = tbl.TableName()
	}
	if tableName == "" {
		tableName = underscoreName(elemTyp.Name())
	}

	res := &Model{
		tableName: tableName,
		fileMap:   fieldMap,
	}

	for _, opt := range opts {
		err := opt(res)
		if err != nil {
			return nil, err
		}
	}

	r.models.Store(typ, res)
	return res, nil
}

func ModelWithTableName(tableName string) ModelOpt {
	return func(m *Model) error {
		m.tableName = tableName
		return nil
	}
}

//field 字段名，，，，，colName 列名
func ModelWithColumnName(field string, colName string) ModelOpt {
	return func(m *Model) error {
		fd, ok := m.fileMap[field]
		if !ok {
			return errs.NewErrUnknownField(field)
		}
		fd.colName = colName
		return nil
	}
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
