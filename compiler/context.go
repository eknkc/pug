package compiler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"path/filepath"
	"strings"
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

	define(string, ...func() error) (*Define, error)

	ParseFile(name string) (*Root, error)
	CompileFile(name string) (string, error)

	String() string
	WriteTo(io.Writer) (int64, error)
}

type context struct {
	body         *bytes.Buffer
	dir          Dir
	indentLevel  int
	indentString string
	path         string
	definitions  map[string]*Define
}

func (bw *context) clone() *context {
	return &context{
		body:         &bytes.Buffer{},
		dir:          bw.dir,
		indentString: bw.indentString,
	}
}

func (bw *context) CompileFile(name string) (string, error) {
	root, err := bw.ParseFile(name)

	if err != nil {
		return "", err
	}

	if err = root.Compile(bw, nil); err != nil {
		return "", err
	}

	for _, def := range bw.definitions {
		if err := def.Compile(bw, root); err != nil {
			return "", err
		}
	}

	return bw.String(), nil
}

func (bw *context) ParseFile(name string) (*Root, error) {
	reader, err := bw.dir.Open(name)

	if err != nil {
		return nil, err
	}

	prep, _, err := preprocess(reader)

	if err != nil {
		return nil, err
	}

	ret, err := Parse(name, []byte(prep))

	if err != nil {
		// if errList, ok := err.(errList); ok {
		// 	for _, err := range errList {
		// 		if parseErr, ok := err.(*parserError); ok {
		// 			spew.Dump(parseErr)
		// 		}
		// 	}
		// }

		return nil, err
	}

	root := ret.(*Root)
	root.Filename = name

	if bw.path == "" {
		bw.path = filepath.Clean(name)
	}

	if extend, err := bw.checkExtend(root); err != nil {
		return root, err
	} else if extend != nil {
		extend.Handled = true
		parentRoot, err := bw.ParseFile(filepath.Join(filepath.Dir(root.Filename), extend.File))

		if err != nil {
			return root, err
		}

		parentRoot.List.Nodes = append(parentRoot.List.Nodes, root)
		root.Extends = parentRoot

		root = parentRoot
	}

	return root, nil
}

func (bw *context) checkExtend(root Node) (*Extend, error) {
	rn, ok := root.(*Root)

	if !ok {
		return nil, errors.New("Unexpected root node")
	}

	var ex *Extend

	for _, node := range rn.List.Nodes {
		if extend, ok := node.(*Extend); ok {
			ex = extend
			break
		}
	}

	if ex == nil {
		return nil, nil
	}

	for _, node := range rn.List.Nodes {
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

func (bw *context) define(name string, definer ...func() error) (*Define, error) {
	if len(definer) == 0 {
		return bw.definitions[name], nil
	}

	definerFunc := definer[0]

	body, indentLevel := bw.body, bw.indentLevel

	bw.body = &bytes.Buffer{}
	bw.indentLevel = 1

	defBody := bw.body

	err := definerFunc()

	bw.body, bw.indentLevel = body, indentLevel

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
func (bw *context) beginLine() { bw.write(strings.Repeat(bw.indentString, bw.indentLevel)) }

func (bw *context) endLine() {
	if bw.indentString != "" {
		bw.write("\n")
	}
}

func NewContext(dir Dir, indentString string) Context {
	return &context{
		body:         &bytes.Buffer{},
		dir:          dir,
		indentString: indentString,
	}
}
