package compiler

import (
	"errors"
	"fmt"
	"html"
	"path/filepath"
	"strconv"
)

func (n *Root) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	if err := n.List.Compile(w, n); err != nil {
		return err
	}

	return nil
}

func (n *Define) Compile(w Context, parent Node) (err error) {
	if n.Hidden {
		return nil
	}

	w.writeLinef("{{ define %s }}", strconv.Quote(n.Name))
	w.write(n.Tpl)
	w.writeLine("{{ end }}")

	return
}

func (n *List) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	for _, node := range n.Nodes {
		if err := node.Compile(w, n); err != nil {
			return err
		}
	}

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
		w.write(`}}"`)
	}

	return
}

func (n Interpolation) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	rawValue := n.Expr.RawValue(w, n)

	if rawValue != nil && !n.Unescaped {
		w.write(html.EscapeString(*rawValue))
	} else {
		w.write("{{ ")

		if n.Unescaped {
			w.write("__pug_unescape ")
		}

		if err := n.Expr.Compile(w, n); err != nil {
			return err
		}
		w.write(" }}")
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
	w.write(" }}")
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
		n.IndexVariable = n.Parent.variable(n.IndexVariable)
		w.writef("$%s, ", n.IndexVariable.Name)
	}

	n.ElementVariable = n.Parent.variable(n.ElementVariable)
	w.writef("$%s := ", n.ElementVariable.Name)

	if err := n.Container.Compile(w, n); err != nil {
		return err
	}

	w.write(" }}")
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
		for i, arg := range n.Arguments {
			arg.Name = n.Parent.variable(arg.Name)
			w.writeLinef("{{ $%s := index . %d }}", arg.Name.Name, i+1)
		}

		w.writeLine(`{{ with index . 0 }}`)
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
	w.writef("{{ template %s (__pug_slice .", strconv.Quote(mixin.Name))

	for _, arg := range n.Arguments {
		w.write(" ")
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

	if n.variable(n.Variable, true) != nil {
		w.write("$")
	} else {
		w.write(".")
	}

	w.write(n.Variable.Name)

	return
}

func (n *FunctionCallExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.write("(")

	if field, ok := n.X.(*FieldExpression); ok && n.variable(field.Variable, true) == nil {
		w.write(field.Variable.Name)
	} else {
		w.write("call (")

		if err := n.X.Compile(w, n); err != nil {
			return err
		}

		w.write(")")
	}

	for _, a := range n.Arguments {
		w.write(" ")

		if err := a.Compile(w, n); err != nil {
			return err
		}
	}

	w.writef(")")

	return
}

func (n *MemberExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	if err := n.X.Compile(w, n); err != nil {
		return err
	}

	w.write(".")
	w.write(n.Name)

	return
}

func (n *StringExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.writef("%s", strconv.Quote(n.Value))
	return
}

func (n *FloatExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.writef("%f", n.Value)
	return
}

func (n *IntegerExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.writef("%d", n.Value)
	return
}

func (n *BooleanExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	if n.Value {
		w.write("true")
	} else {
		w.write("false")
	}

	return
}

func (n *NilExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.write("(__pug_nil)")
	return
}

func (n *UnaryExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.writef(`(__pug_unaryop %s `, strconv.Quote(n.Op))

	if err := n.X.Compile(w, parent); err != nil {
		return err
	}

	w.write(`)`)

	return
}

func (n *BinaryExpression) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	w.writef(`(__pug_binaryop %s `, strconv.Quote(n.Op))

	if err := n.X.Compile(w, n); err != nil {
		return err
	}

	w.write(" ")

	if err := n.Y.Compile(w, n); err != nil {
		return err
	}

	w.write(")")

	return
}

func (n *Assignment) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	n.Variable = n.Parent.variable(n.Variable)

	if rawExpr := n.Expression.RawValue(w, n); rawExpr != nil {
		w.writeLinef("{{ $%s := %s }}", n.Variable.Name, *rawExpr)
		return
	}

	w.beginLine()
	w.writef("{{ $%s := ", n.Variable.Name)
	if err := n.Expression.Compile(w, n); err != nil {
		return err
	}
	w.write(" }}")
	w.endLine()

	return
}

func (n *Import) Compile(w Context, parent Node) (err error) {
	if err := n.GraphNode.Compile(w, parent); err != nil {
		return err
	}

	root := n.root()
	file := filepath.Join(filepath.Dir(root.Filename), n.File)

	if def, err := w.define(file); err != nil {
		return err
	} else if def != nil {
		w.writeLinef("{{ template %s . }}", strconv.Quote(def.Name))
		return nil
	}

	if ast, err := w.ParseFile(file); err != nil {
		return err
	} else {
		_, err = w.define(file, func() error {
			return ast.Compile(w, root)
		})

		if err != nil {
			return err
		}

		w.writeLinef("{{ template %s . }}", strconv.Quote(file))
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

	root := n.root()

	renderBlock := func() error {
		if n.Block != nil {
			if err := n.Block.Compile(w, n); err != nil {
				return err
			}
		}

		return nil
	}

	if root.Extends == nil {
		_, err := w.define(fmt.Sprintf("%s", n.Name), func() error {
			return renderBlock()
		})

		if err != nil {
			return err
		}

		w.writeLinef("{{ template %s . }}", strconv.Quote(n.Name))
	} else {
		if def, err := w.define(fmt.Sprintf("%s", n.Name)); err != nil {
			return err
		} else if def != nil {
			_, err := w.define(fmt.Sprintf("%s", n.Name), func() error {
				if n.Modifier == "prepend" {
					if err := renderBlock(); err != nil {
						return err
					}

					w.write(def.Tpl)
					return nil
				}

				if n.Modifier == "append" {
					w.write(def.Tpl)

					if err := renderBlock(); err != nil {
						return err
					}

					return nil
				}

				return renderBlock()
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (n *GraphNode) Compile(w Context, parent Node) (err error) {
	if n == nil {
		return
	}

	n.Parent = parent
	return
}
