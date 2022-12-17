package reflect

import "reflect"

func IterateFunc(entity any) (map[string]FuncInfo, error) {
	// 先拿到它的 类型信息
	typ := reflect.TypeOf(entity)
	numMethod := typ.NumMethod()
	res := make(map[string]FuncInfo, numMethod)
	// 开始遍历它的方法
	for i := 0; i < numMethod; i++ {
		// 拿到第n个方法
		method := typ.Method(i)
		fn := method.Func

		// 访问它有多少个参数
		numIn := fn.Type().NumIn()
		input := make([]reflect.Type, 0, numIn)

		// 构建准备调用发起的参数
		inputValues := make([]reflect.Value, 0, numIn)

		// 第0个的时候是拿它 entity 本身附加上去
		input = append(input, reflect.TypeOf(entity))
		inputValues = append(inputValues, reflect.ValueOf(entity))

		// 输入
		for j := 1; j < numIn; j++ {
			fnInType := fn.Type().In(j)
			input = append(input, fnInType)
			inputValues = append(inputValues, reflect.Zero(fnInType)) // 都用 零值
		}

		numOut := fn.Type().NumOut()
		output := make([]reflect.Type, 0, numOut)
		for j := 0; j < numOut; j++ {
			// 把输出的类型记一下
			output = append(output, fn.Type().Out(j))
		}

		// resValues 它是一个反射的values
		resValues := fn.Call(inputValues)
		result := make([]any, 0, len(resValues))
		for _, v := range resValues {
			result = append(result, v.Interface())
		}
		res[method.Name] = FuncInfo{
			Name:        method.Name,
			InputTypes:  input,
			OutputTypes: output,
			Result:      result,
		}
	}

	return res, nil
}

type FuncInfo struct {
	Name        string
	InputTypes  []reflect.Type
	OutputTypes []reflect.Type
	Result      []any
}
