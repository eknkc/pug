package compiler

import (
	"errors"
	"fmt"
	"html"
	"strconv"
)

func (n *Root) Compile(w Context, parent Node) (err error) {
	return n.List.Compile(w, parent)
}

func (n *Define) Compile(w Context, parent Node) (err error) {
	if n.Written {
		return nil
	}

	w.writeLinef("{{ define %s }}", strconv.Quote(n.Name))
	w.write(n.Tpl)
	w.writeLine("{{ end }}")

	n.Written = true

	return
}

func (n *List) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.pushScope()

	for _, node := range n.Nodes {
		if err := node.Compile(w, n); err != nil {
			return err
		}
	}

	w.popScope()

	return nil
}

func (n *TextList) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.beginLine()
	for _, node := range n.Nodes {
		if err := node.Compile(w, n); err != nil {
			return err
		}
	}
	w.endLine()

	return nil
}

func (n *Text) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.write(n.Value)
	return
}

func (n *DocType) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	if defined, ok := doctypes[n.Value]; ok {
		w.writeLine(defined)
	} else {
		w.writeLinef("<!DOCTYPE %s >", n.Value)
	}

	return
}

func (n *Comment) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	if !n.Silent {
		w.writeLinef("<!-- %s -->", n.Value)
	}

	return
}

func (n *Tag) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	selfClosing := n.Text == nil && n.Block == nil && selfClosingTags[n.Name]

	w.beginLine()
	w.writef("<%s", n.Name)

	n.fixAttributes()

	for _, attr := range n.Attributes {
		w.write(" ")

		if err := attr.Compile(w, parent); err != nil {
			return err
		}
	}

	if !selfClosing {
		w.write(">")
		w.endLine()

		w.indent()
		if n.Text != nil {
			if err := n.Text.Compile(w, n); err != nil {
				return err
			}
		}

		if n.Block != nil {
			if err := n.Block.Compile(w, n); err != nil {
				return err
			}
		}
		w.outdent()

		w.writeLinef("</%s>", n.Name)
	} else {
		w.write(" />")
		w.endLine()
	}

	return
}

func (n *Tag) fixAttributes() {
	var class *Attribute
	var attrs []*Attribute

	for _, attr := range n.Attributes {
		if attr.Name == "class" {
			if class == nil {
				class = attr
				attrs = append(attrs, attr)
			} else {
				strclass, sok := class.Value.(*StringExpression)
				strattr, aok := attr.Value.(*StringExpression)

				if sok && aok {
					strclass.Value = strclass.Value + " " + strattr.Value
				} else {
					class.Value = &BinaryExpression{
						GraphNode: NewNode(class.Position),
						Op:        "|",
						X:         class.Value,
						Y:         attr.Value,
					}
				}
			}
		} else {
			attrs = append(attrs, attr)
		}
	}

	n.Attributes = attrs
}

func (n *Attribute) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	if n.Value == nil {
		w.write(n.Name)
		return
	}

	strAttribute, ok := n.Value.(*StringExpression)

	if ok {
		w.writef(`%s="%s"`, n.Name, html.EscapeString(strAttribute.Value))
	} else {
		w.writef(`%s="{{`, n.Name)
		if err := n.Value.Compile(w, n); err != nil {
			return err
		}
		w.write(` | __pug_pop}}"`)
	}

	return
}

func (n Interpolation) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	rawValue := n.Expr.RawValue(w, n)

	if rawValue != nil {
		w.write(html.EscapeString(*rawValue))
	} else {
		w.write("{{ ")
		if err := n.Expr.Compile(w, n); err != nil {
			return err
		}
		w.write(" | __pug_pop }}")
	}

	return
}

func (n *If) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.beginLine()
	w.write("{{ if ")
	if err := n.Condition.Compile(w, n); err != nil {
		return err
	}
	w.write(" | __pug_pop }}")
	w.endLine()

	w.indent()
	if err := n.PositiveBlock.Compile(w, n); err != nil {
		return err
	}
	w.outdent()

	if n.NegativeBlock != nil {
		w.writeLine("{{ else }}")

		w.indent()
		if err := n.NegativeBlock.Compile(w, n); err != nil {
			return err
		}
		w.outdent()
	}

	w.writeLine("{{ end }}")

	return
}

func (n *Each) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.beginLine()
	w.write("{{ range ")

	if n.IndexVariable != nil {
		n.IndexVariable = w.setVariable(n.IndexVariable)
		w.writef("$%s, ", n.IndexVariable.Name)
	}

	n.ElementVariable = w.setVariable(n.ElementVariable)
	w.writef("$%s := ", n.ElementVariable.Name)

	if err := n.Container.Compile(w, n); err != nil {
		return err
	}

	w.write(" | __pug_pop }}")
	w.endLine()

	w.indent()
	if err := n.Block.Compile(w, n); err != nil {
		return err
	}
	w.outdent()

	w.writeLine("{{ end }}")

	return
}

func (n *Mixin) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	_, err = w.define(fmt.Sprintf("mixin-%s", n.Name), func() error {
		for _, arg := range n.Arguments {
			w.beginLine()

			arg.Name = w.setVariable(arg.Name)
			w.writef("{{ $%s := __pug_pop . ", arg.Name.Name)

			if arg.Default != nil {
				raw := arg.Default.RawValue(w, n)

				if raw != nil {
					w.write(*raw)
					w.write(" ")
				} else {
					w.write("(")
					if err := arg.Default.Compile(w, n); err != nil {
						return err
					}
					w.write(" | __pug_pop) ")
				}
			}

			w.write("}}")
			w.endLine()
		}

		w.writeLine(`{{ with (__pug_binaryop "||" true .) | __pug_pop }}`)
		if n.Block != nil {
			w.indent()
			if err := n.Block.Compile(w, n); err != nil {
				return err
			}
			w.outdent()
		}
		w.writeLine("{{ end }}")

		return nil
	})

	return err
}

func (n *MixinCall) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	name := fmt.Sprintf("mixin-%s", n.Name)

	mixin, err := w.define(name)

	if err != nil {
		return err
	}

	if mixin == nil {
		return fmt.Errorf("mixin not defined: %s", n.Name)
	}

	w.beginLine()
	w.writef("{{ template %s (__pug_push . ", strconv.Quote(mixin.Name))

	for i := len(n.Arguments) - 1; i >= 0; i-- {
		arg := n.Arguments[i]

		w.write(" | ")
		if err := arg.Compile(w, n); err != nil {
			return err
		}
	}

	w.write(") }}")
	w.endLine()

	return nil
}

func (n *FieldExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.write("__pug_push ")

	if w.variable(n.Name) != nil {
		w.write("$")
	} else {
		w.write(".")
	}

	w.write(n.Name)

	return
}

func (n *FunctionCallExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	if err := n.X.Compile(w, n); err != nil {
		return err
	}

	for _, a := range n.Arguments {
		w.write(" | ")

		if err := a.Compile(w, n); err != nil {
			return err
		}
	}

	w.writef(" | __pug_call %d", len(n.Arguments))

	return
}

func (n *MemberExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	var cur Expression = n

	for {
		if member, ok := cur.(*MemberExpression); ok {
			cur = member.X
		} else {
			if err := n.X.Compile(w, n); err != nil {
				return err
			}

			if _, ok := cur.(*FieldExpression); ok {
				w.write(".")
				w.write(n.Name)
				break
			} else {
				w.writef(" | __pug_field %s", strconv.Quote(n.Name))
				break
			}
		}
	}

	return
}

func (n *StringExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.writef("__pug_push %s", strconv.Quote(n.Value))
	return
}

func (n *FloatExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.writef("__pug_push %f", n.Value)
	return
}

func (n *IntegerExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.writef("__pug_push %d", n.Value)
	return
}

func (n *BooleanExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	if n.Value {
		w.write("__pug_push true")
	} else {
		w.write("__pug_push false")
	}

	return
}

func (n *NilExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.write("__pug_push nil")
	return
}

func (n *UnaryExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	if err := n.X.Compile(w, parent); err != nil {
		return err
	}

	w.writef(` | __pug_unaryop %s`, strconv.Quote(n.Op))
	return
}

func (n *BinaryExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	rawX := n.X.RawValue(w, n)
	rawY := n.Y.RawValue(w, n)

	if rawX != nil && rawY != nil {
		w.writef(`__pug_binaryop %s %s %s`, strconv.Quote(n.Op), *rawX, *rawY)
	} else {
		if err := n.X.Compile(w, n); err != nil {
			return err
		}

		if rawY != nil {
			w.writef(` | __pug_binaryop %s %s`, strconv.Quote(n.Op), *rawY)
		} else {
			w.write(" | ")

			if err := n.Y.Compile(w, n); err != nil {
				return err
			}

			w.writef(` | __pug_binaryop %s`, strconv.Quote(n.Op))
		}
	}

	return
}

func (n *Assignment) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	n.Variable = w.setVariable(n.Variable)

	if rawExpr := n.Expression.RawValue(w, n); rawExpr != nil {
		w.writeLinef("{{ $%s := %s }}", n.Variable.Name, *rawExpr)
		return
	}

	w.beginLine()
	w.writef("{{ $%s := ", n.Variable.Name)
	if err := n.Expression.Compile(w, n); err != nil {
		return err
	}
	w.write(" | __pug_pop }}")
	w.endLine()

	return
}

func (n *ExpressionList) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	for i, ex := range n.Expressions {
		if i > 0 {
			w.write(n.Seperator)
		}

		if err := ex.Compile(w, n); err != nil {
			return err
		}
	}

	return
}

func (n *Import) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	if name, err := w.include(n.File); err != nil {
		return err
	} else {
		w.writeLinef("{{ template %s . }}", strconv.Quote(name))
	}

	return
}

func (n *Extend) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	if !n.Handled {
		return errors.New("extend statements can only be present on root node")
	}

	return nil
}

func (n *Block) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	// parentBlock, _ := w.define(n.Name)

	// if n.Modifier != "" && n.Block != nil && parentBlock != nil {
	// 	_, err = w.define(fmt.Sprintf("%s-%s", parentBlock.Name, n.Modifier), func() error {
	// 		return n.Block.Compile(w, n)
	// 	})

	// 	return err
	// }

	if _, err := w.define(fmt.Sprintf("%s-prepend", n.Name), func() error { return nil }); err != nil {
		return err
	}

	if _, err := w.define(fmt.Sprintf("%s-append", n.Name), func() error { return nil }); err != nil {
		return err
	}

	if _, err := w.define(n.Name, func() error { return nil }); err != nil {
		if n.Block != nil {
			w.indent()
			if err := n.Block.Compile(w, n); err != nil {
				return err
			}
			w.outdent()
		}
	}

	w.writeLinef(`{{ template %s . }}`, strconv.Quote(fmt.Sprintf("%s-prepend", n.Name)))
	w.writeLinef(`{{ template %s . }}`, strconv.Quote(n.Name))
	w.writeLinef(`{{ template %s . }}`, strconv.Quote(fmt.Sprintf("%s-append", n.Name)))

	// w.writeLinef(`{{ define %s }}{{ end }}`, strconv.Quote(fmt.Sprintf("%s_prepend", n.GlobalName)))
	// w.writeLinef(`{{ define %s }}{{ end }}`, strconv.Quote(fmt.Sprintf("%s_append", n.GlobalName)))
	// w.writeLinef(`{{ define %s }}`, strconv.Quote(n.GlobalName))
	// if n.Modifier == "" && n.Block != nil {
	// 	w.indent()
	// 	if err := n.Block.Compile(w, n); err != nil {
	// 		return err
	// 	}
	// 	w.outdent()
	// }
	// w.writeLine(`{{ end }}`)

	// if pb != nil {
	// 	w.writeLinef(`{{ define %s }}`, strconv.Quote(fmt.Sprintf("%s", pb.GlobalName)))
	// 	w.indent()
	// }

	// w.writeLinef(`{{ template %s . }}`, strconv.Quote(fmt.Sprintf("%s_prepend", n.GlobalName)))
	// w.writeLinef(`{{ template %s . }}`, strconv.Quote(n.GlobalName))
	// w.writeLinef(`{{ template %s . }}`, strconv.Quote(fmt.Sprintf("%s_append", n.GlobalName)))

	// if pb != nil {
	// 	w.outdent()
	// 	w.writeLine(`{{ end }}`)
	// }

	return nil
}

func (n *GraphNode) Compile(w Context, parent Node) (err error) {
	if n == nil {
		return
	}

	n.Parent = parent
	return
}
