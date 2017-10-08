package runtime

import (
	"errors"
	"fmt"
	"html/template"
)

var FuncMap template.FuncMap = template.FuncMap{
	"__pug_binaryop":     Binary,
	"__pug_unaryop":      Unary,
	"__pug_nil":          Nil,
	"__pug_slice":        Slice,
	"__pug_unescape":     Unescape,
	"__pug_unescapeattr": UnescapeAttr,
}

func Binary(op string, x, y interface{}) (interface{}, error) {
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
					return nil, errors.New("Unexpected operator")
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
					return nil, fmt.Errorf("unsupported operator %s between two integers", op)
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
					return nil, fmt.Errorf("unsupported operator %s between two floats", op)
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
					return nil, fmt.Errorf("unsupported operator %s between two strings", op)
				}
			}
		}
	}

	return res, nil
}

func Unary(op string, x interface{}) interface{} {
	if op == "-" {
		if ix, ok := makeInt(x); ok {
			x = -1 * ix
		} else if fx, ok := makeFloat(x); ok {
			x = -1 * fx
		}
	}

	return x
}

func Slice(variables ...interface{}) interface{} {
	return variables
}

func Unescape(val string) template.HTML {
	return template.HTML(val)
}

func UnescapeAttr(name, val string) template.HTMLAttr {
	return template.HTMLAttr(fmt.Sprintf(` %s="%s"`, name, val))
}

func Nil() interface{} {
	return nil
}
