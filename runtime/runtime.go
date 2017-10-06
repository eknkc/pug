package runtime

import (
	"errors"
	"fmt"
	"html/template"
	"reflect"
)

var FuncMap template.FuncMap = template.FuncMap{
	"__pug_push":     Push,
	"__pug_pop":      Pop,
	"__pug_binaryop": Binary,
	"__pug_unaryop":  Unary,
	"__pug_field":    Field,
	"__pug_call":     Call,
}

func toSlice(v interface{}) *stack {
	if v == nil {
		return &stack{}
	}

	if vslice, ok := v.(*stack); ok {
		return vslice
	}

	return &stack{values: []interface{}{v}}
}

func Push(values ...interface{}) *stack {
	s := &stack{}

	if len(values) == 0 {
		return s
	}

	if oldstack, ok := values[len(values)-1].(*stack); ok {
		s = oldstack
		values = values[:len(values)-1]
	}

	for _, val := range values {
		s.push(val)
	}

	return s
}

func Pop(s *stack, def ...interface{}) (interface{}, error) {
	if s.len() == 0 {
		if len(def) > 0 {
			return def[0], nil
		}

		return nil, errors.New("Stack empty")
	}

	return s.pop(), nil
}

func Binary(op string, values ...interface{}) (*stack, error) {
	stack := Push(values...)

	y := stack.pop()
	x := stack.pop()

	var res interface{}

	if op == "||" || op == "&&" {
		bx := makeBool(x)
		by := makeBool(y)

		switch op {
		case "||":
			if bx {
				res = x
			} else if by {
				res = y
			} else {
				res = false
			}
		case "&&":
			if bx && by {
				res = y
			} else {
				res = false
			}
		}
	} else {
		if ix, ok := makeInt(x); ok {
			if iy, ok := makeInt(y); ok {
				switch op {
				case "+":
					res = ix + iy
				case "-":
					res = ix - iy
				case "*":
					res = ix * iy
				case "/":
					res = ix / iy
				case "%":
					res = ix % iy
				case "<":
					res = ix < iy
				case ">":
					res = ix > iy
				case ">=":
					res = ix >= iy
				case "<=":
					res = ix <= iy
				case "==":
					res = ix == iy
				case "!=":
					res = ix != iy
				default:
					return stack, errors.New("Unexpected operator")
				}
			} else if fy, ok := makeFloat(y); ok {
				fx := float64(ix)

				switch op {
				case "+":
					res = fx + fy
				case "-":
					res = fx - fy
				case "*":
					res = fx * fy
				case "/":
					res = fx / fy
				case "<":
					res = fx < fy
				case ">":
					res = fx > fy
				case ">=":
					res = fx >= fy
				case "<=":
					res = fx <= fy
				case "==":
					res = fx == fy
				case "!=":
					res = fx != fy
				default:
					return stack, fmt.Errorf("unsupported operator %s between two integers", op)
				}
			}
		} else if fx, ok := makeFloat(x); ok {
			if fy, ok := makeFloat(y); ok {
				switch op {
				case "+":
					res = fx + fy
				case "-":
					res = fx - fy
				case "*":
					res = fx * fy
				case "/":
					res = fx / fy
				case "<":
					res = fx < fy
				case ">":
					res = fx > fy
				case ">=":
					res = fx >= fy
				case "<=":
					res = fx <= fy
				case "==":
					res = fx == fy
				case "!=":
					res = fx != fy
				default:
					return stack, fmt.Errorf("unsupported operator %s between two floats", op)
				}
			}
		}
	}

	if res == nil {
		if sx, ok := makeString(x); ok {
			if sy, ok := makeString(y); ok {
				switch op {
				case "+":
					res = sx + sy
				case "|":
					res = sx + " " + sy
				case "==":
					res = sx == sy
				case "!=":
					res = sx != sy
				default:
					return stack, fmt.Errorf("unsupported operator %s between two strings", op)
				}
			}
		}
	}

	if res == nil && op == "+" {

	}

	stack.push(res)
	return stack, nil
}

func Unary(op string, values ...interface{}) *stack {
	stack := Push(values...)
	x := stack.pop()

	if op == "-" {
		if ix, ok := makeInt(x); ok {
			x = -1 * ix
		} else if fx, ok := makeFloat(x); ok {
			x = -1 * fx
		}
	}

	stack.push(x)

	return stack
}

func Field(name string, values ...interface{}) (*stack, error) {
	stack := Push(values...)
	x := stack.pop()

	xVal := reflect.ValueOf(x)

	if !xVal.IsValid() {
		return stack, fmt.Errorf("nil pointer evaluating %s.%s", xVal.Type(), name)
	}

	xVal, isNil := indirect(xVal)

	switch xVal.Kind() {
	case reflect.Struct:
		if tField, ok := xVal.Type().FieldByName(name); ok {
			if tField.PkgPath != "" {
				return stack, fmt.Errorf("%s is an unexported field of struct type %s", name, xVal.Type())
			}

			if isNil {
				return stack, fmt.Errorf("nil pointer evaluating %s.%s", xVal.Type(), name)
			}

			vField := xVal.FieldByIndex(tField.Index)

			if vField.IsValid() {
				stack.push(vField.Interface())
			} else {
				return stack, fmt.Errorf("invalid field access: %s", name)
			}
		}
	case reflect.Map:
		if isNil {
			return stack, fmt.Errorf("nil pointer evaluating %s.%s", xVal.Type(), name)
		}

		nameVal := reflect.ValueOf(name)
		if nameVal.Type().AssignableTo(xVal.Type().Key()) {
			result := xVal.MapIndex(nameVal)
			if !result.IsValid() {
				result = reflect.Zero(xVal.Type().Elem())
			}
			stack.push(result.Interface())
		}
	default:
		return stack, fmt.Errorf("can't evaluate field %s in type %s", name, xVal.Type())
	}

	return stack, nil
}

func Call(argCount int, values ...interface{}) (*stack, error) {
	stack := Push(values...)
	args := make([]interface{}, argCount)

	for i := 0; i < argCount; i++ {
		args[i] = stack.pop()
	}

	x := stack.pop()
	v := indirectInterface(reflect.ValueOf(x))

	if !v.IsValid() {
		return stack, fmt.Errorf("call of nil")
	}

	typ := v.Type()
	if typ.Kind() != reflect.Func {
		return stack, fmt.Errorf("non-function of type %s", typ)
	}
	if !goodFunc(typ) {
		return stack, fmt.Errorf("function called with %d args; should be 1 or 2", typ.NumOut())
	}

	numIn := typ.NumIn()
	var dddType reflect.Type

	if typ.IsVariadic() {
		if argCount < numIn-1 {
			return stack, fmt.Errorf("wrong number of args: got %d want at least %d", argCount, numIn-1)
		}
		dddType = typ.In(numIn - 1).Elem()
	} else {
		if argCount != numIn {
			return stack, fmt.Errorf("wrong number of args: got %d want %d", argCount, numIn)
		}
	}

	argv := make([]reflect.Value, argCount)

	for i := 0; i < argCount; i++ {
		arg := args[len(args)-1-i]
		value := indirectInterface(reflect.ValueOf(arg))

		var argType reflect.Type
		if !typ.IsVariadic() || i < numIn-1 {
			argType = typ.In(i)
		} else {
			argType = dddType
		}

		var err error
		if argv[i], err = prepareArg(value, argType); err != nil {
			return stack, fmt.Errorf("arg %d: %s", i, err)
		}
	}

	result := v.Call(argv)
	stack.push(result[0].Interface())

	if len(result) == 2 && !result[1].IsNil() {
		return stack, result[1].Interface().(error)
	}

	return stack, nil
}
