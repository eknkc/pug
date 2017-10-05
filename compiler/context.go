package compiler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

type Context interface {
	write(s string)
	writeLine(s string)
	writef(pattern string, data ...interface{})
	writeLinef(pattern string, data ...interface{})

	indent()
	outdent()
	beginLine()
	endLine()

	pushScope()
	popScope()

	variable(string) *Variable
	setVariable(*Variable) *Variable

	include(name string) (string, error)
	define(string, ...func() error) (*Define, error)
	block(string) *Define

	ParseFile(name string) (Node, error)
	CompileFile(name string) (string, error)

	String() string
	WriteTo(io.Writer) (int64, error)
}

type scope struct {
	Variables map[string]*Variable
	Parent    *scope
}

func newScope() *scope {
	return &scope{
		Variables: make(map[string]*Variable),
	}
}

type context struct {
	body         *bytes.Buffer
	Dir          Dir
	Scope        *scope
	indentLevel  int
	IndentString string
	path         string
	parent       *context
	blocks       map[string][]*Block
	definitions  map[string]*Define
	extend       *Extend
}

func (bw *context) clone() *context {
	return &context{
		body:         &bytes.Buffer{},
		Dir:          bw.Dir,
		IndentString: bw.IndentString,
		parent:       bw,
		Scope:        newScope(),
	}
}

func (bw *context) CompileFile(name string) (string, error) {
	root, err := bw.ParseFile(name)

	if err != nil {
		return "", err
	}

	tplstring := ""

	if extend, err := bw.checkExtend(root); err != nil {
		return "", err
	} else if extend != nil {
		extend.Handled = true
		extend.File = filepath.Join(filepath.Dir(bw.path), extend.File)
		bw.extend = extend

		clone := bw.clone()
		if tpl, err := clone.CompileFile(extend.File); err != nil {
			return "", err
		} else {
			tplstring += tpl
			bw.definitions = clone.definitions
		}
	}

	if err = root.Compile(bw, nil); err != nil {
		return "", err
	}

	if bw.parent == nil {
		for _, def := range bw.definitions {
			if err := def.Compile(bw, root); err != nil {
				return "", err
			}
		}
	}

	tplstring += bw.String()

	return tplstring, nil
}

func (bw *context) ParseFile(name string) (Node, error) {
	reader, err := bw.Dir.Open(name)

	if err != nil {
		return nil, err
	}

	prep, _, err := preprocess(reader)

	if err != nil {
		return nil, err
	}

	ret, err := Parse(name, []byte(prep))

	if err != nil {
		if errList, ok := err.(errList); ok {
			for _, err := range errList {
				if parseErr, ok := err.(*parserError); ok {
					spew.Dump(parseErr)
				}
			}
		}

		return nil, err
	}

	bw.path = name

	return ret.(Node), nil
}

func (bw *context) rootContext() *context {
	for bw.parent != nil {
		bw = bw.parent
	}

	return bw
}

func (bw *context) checkExtend(root Node) (*Extend, error) {
	rn, ok := root.(*Root)

	if !ok {
		return nil, errors.New("Unexpected root node")
	}

	var ex *Extend

	for _, node := range rn.Nodes {
		if extend, ok := node.(*Extend); ok {
			ex = extend
			break
		}
	}

	if ex == nil {
		return nil, nil
	}

	for _, node := range rn.Nodes {
		switch node.(type) {
		case *Extend, *Mixin, *Block, *Comment:
			continue
		default:
			return nil, errors.New("extending templates can only contain mixin definitions and blocks on root level")
		}
	}

	return ex, nil
}

func (bw *context) write(s string) {
	bw.body.WriteString(s)
}

func (bw *context) writeLine(s string) {
	bw.beginLine()
	bw.write(s)
	bw.endLine()
}

func (bw *context) writef(pattern string, data ...interface{}) {
	bw.body.WriteString(fmt.Sprintf(pattern, data...))
}

func (bw *context) writeLinef(pattern string, data ...interface{}) {
	bw.beginLine()
	bw.writef(pattern, data...)
	bw.endLine()
}

func (bw *context) String() string {
	return string(bw.body.Bytes())
}

func (bw *context) WriteTo(w io.Writer) (int64, error) {
	return bw.body.WriteTo(w)
}

func (bw *context) pushScope() {
	ns := newScope()
	ns.Parent = bw.Scope
	bw.Scope = ns
}

func (bw *context) popScope() {
	if bw.Scope.Parent != nil {
		bw.Scope = bw.Scope.Parent
	}
}

func (bw *context) variable(name string) *Variable {
	s := bw.Scope

	for s != nil {
		if v, ok := s.Variables[name]; ok {
			return v
		}

		s = s.Parent
	}

	return nil
}

func (bw *context) block(name string) *Define {
	return nil
}

func (bw *context) setVariable(v *Variable) *Variable {
	if bw.variable(v.Name) == nil {
		bw.Scope.Variables[v.Name] = v
	}

	return v
}

func (bw *context) include(name string) (string, error) {
	root := bw.rootContext()
	name = filepath.Join(filepath.Dir(bw.path), name)

	if _, ok := root.definitions[name]; ok {
		return name, nil
	}

	clone := bw.clone()
	clone.indent()

	if tpl, err := clone.CompileFile(name); err != nil {
		return "", err
	} else {
		if root.definitions == nil {
			root.definitions = make(map[string]*Define)
		}

		root.definitions[name] = &Define{Name: name, Tpl: tpl}

		return name, nil
	}
}

func (bw *context) define(name string, definer ...func() error) (*Define, error) {
	if len(definer) == 0 {
		return bw.definitions[name], nil
	}

	definerFunc := definer[0]

	body, indentLevel, scope := bw.body, bw.indentLevel, bw.Scope

	bw.body = &bytes.Buffer{}
	bw.Scope = newScope()
	bw.indentLevel = 1

	defBody := bw.body

	err := definerFunc()

	bw.body, bw.indentLevel, bw.Scope = body, indentLevel, scope

	if err != nil {
		return nil, err
	}

	def := &Define{Name: name, Tpl: string(defBody.Bytes())}

	if bw.definitions == nil {
		bw.definitions = make(map[string]*Define)
	}

	bw.definitions[name] = def

	return def, nil
}

func (bw *context) indent()    { bw.indentLevel++ }
func (bw *context) outdent()   { bw.indentLevel-- }
func (bw *context) beginLine() { bw.write(strings.Repeat(bw.IndentString, bw.indentLevel)) }

func (bw *context) endLine() {
	if bw.IndentString != "" {
		bw.write("\n")
	}
}

func NewContext(dir Dir, indentString string) Context {
	return &context{
		body:         &bytes.Buffer{},
		Dir:          dir,
		Scope:        newScope(),
		IndentString: indentString,
	}
}
