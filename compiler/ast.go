package compiler

import (
	"fmt"
	"strconv"
)

type Node interface {
	Compile(w Context, parent Node) error

	variable(v *Variable, lookup ...bool) *Variable
	graphNode() *GraphNode
	root() *Root
}

type Expression interface {
	Node
	RawValue(w Context, parent Node) *string
}

type Position struct {
	Line   int
	Col    int
	Offset int
}

type GraphNode struct {
	Position
	Parent Node
	Scope  *Scope
}

func (node *GraphNode) graphNode() *GraphNode {
	return node
}

func (node *GraphNode) root() *Root {
	if root, ok := node.Parent.(*Root); ok {
		return root
	}

	return node.Parent.root()
}

func (node *GraphNode) variable(v *Variable, lookup ...bool) *Variable {
	n := node

	for n != nil {
		if n.Scope != nil {
			if variable, ok := n.Scope.Variables[v.Name]; ok {
				return variable
			}
		}

		if n.Parent != nil {
			n = n.Parent.graphNode()
		} else {
			n = nil
		}
	}

	if len(lookup) > 0 && lookup[0] {
		return nil
	}

	if node.Scope == nil {
		node.Scope = &Scope{
			Variables: make(map[string]*Variable),
		}
	}

	node.Scope.Variables[v.Name] = v

	return v
}

func NewNode(pos Position) *GraphNode {
	return &GraphNode{
		Position: pos,
	}
}

type Scope struct {
	Owner     *List
	Variables map[string]*Variable
	Node      Node
}

type Root struct {
	*GraphNode
	Extends  *Root
	List     *List
	Filename string
}

type List struct {
	*GraphNode
	Nodes   []Node
	Append  *List
	Prepend *List
}

type TextList struct {
	*GraphNode
	Nodes []Node
}

type Text struct {
	*GraphNode
	Value string
}

type Interpolation struct {
	*GraphNode
	Expr      Expression
	Unescaped bool
}

type StringExpression struct {
	*GraphNode
	Value string
}

func (s StringExpression) RawValue(w Context, parent Node) *string {
	s.GraphNode.Compile(w, parent)

	quoted := strconv.Quote(s.Value)
	return &quoted
}

type BooleanExpression struct {
	*GraphNode
	Value bool
}

func (s BooleanExpression) RawValue(w Context, parent Node) *string {
	s.GraphNode.Compile(w, parent)

	val := "false"

	if s.Value {
		val = "true"
	}

	return &val
}

type NilExpression struct {
	*GraphNode
}

func (s NilExpression) RawValue(w Context, parent Node) *string {
	s.GraphNode.Compile(w, parent)

	return nil
}

type FloatExpression struct {
	*GraphNode
	Value float64
}

func (s FloatExpression) RawValue(w Context, parent Node) *string {
	s.GraphNode.Compile(w, parent)

	val := fmt.Sprintf("%f", s.Value)
	return &val
}

type IntegerExpression struct {
	*GraphNode
	Value int64
}

func (s IntegerExpression) RawValue(w Context, parent Node) *string {
	s.GraphNode.Compile(w, parent)

	val := fmt.Sprintf("%d", s.Value)
	return &val
}

type ArrayExpression struct {
	*GraphNode
	Expressions []Expression
}

func (s ArrayExpression) RawValue(w Context, parent Node) *string {
	return nil
}

type BinaryExpression struct {
	*GraphNode
	Op string
	X  Expression
	Y  Expression
}

func (s BinaryExpression) RawValue(w Context, parent Node) *string {
	s.GraphNode.Compile(w, parent)
	return nil
}

type UnaryExpression struct {
	*GraphNode
	Op string
	X  Expression
}

func (s UnaryExpression) RawValue(w Context, parent Node) *string {
	s.GraphNode.Compile(w, parent)
	return nil
}

type MemberExpression struct {
	*GraphNode
	X    Expression
	Name string
}

func (s MemberExpression) RawValue(w Context, parent Node) *string {
	s.GraphNode.Compile(w, parent)
	return nil
}

type IndexExpression struct {
	*GraphNode
	X     Expression
	Index Expression
}

func (s IndexExpression) RawValue(w Context, parent Node) *string {
	s.GraphNode.Compile(w, parent)
	return nil
}

type FieldExpression struct {
	*GraphNode
	Variable *Variable
}

func (s FieldExpression) RawValue(w Context, parent Node) *string {
	s.GraphNode.Compile(w, parent)
	return nil
}

type FunctionCallExpression struct {
	*GraphNode
	X         Expression
	Arguments []Expression
}

func (s FunctionCallExpression) RawValue(w Context, parent Node) *string {
	s.GraphNode.Compile(w, parent)
	return nil
}

type Comment struct {
	*GraphNode
	Value  string
	Silent bool
}

type Import struct {
	*GraphNode
	File string
}

type Extend struct {
	*GraphNode
	File    string
	Handled bool
}

type Block struct {
	*GraphNode
	ParentBlock *Block
	Name        string
	Modifier    string
	GlobalName  string
	Block       Node
}

type Attribute struct {
	*GraphNode
	Name      string
	Value     Expression
	Unescaped bool
}

type Tag struct {
	*GraphNode
	Block      Node
	Text       *TextList
	Name       string
	Raw        bool
	SelfClose  bool
	Attributes []*Attribute
}

type MixinArgument struct {
	*GraphNode
	Name    *Variable
	Default Expression
}

type Mixin struct {
	*GraphNode
	Name      string
	Arguments []MixinArgument
	Block     Node
}

type MixinCall struct {
	*GraphNode
	Name      string
	Arguments []Expression
}

type If struct {
	*GraphNode
	PositiveBlock Node
	NegativeBlock Node
	Condition     Expression
}

type Each struct {
	*GraphNode
	IndexVariable   *Variable
	ElementVariable *Variable
	Container       Expression
	Block           Node
}

type Assignment struct {
	*GraphNode
	Variable   *Variable
	Expression Expression
}

type Variable struct {
	*GraphNode
	Name string
}

type DocType struct {
	*GraphNode
	Value string
}

type Define struct {
	*GraphNode
	Name   string
	Tpl    string
	Hidden bool
}
