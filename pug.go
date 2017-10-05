package pug

import (
	"html/template"

	"github.com/eknkc/pug/compiler"
	"github.com/eknkc/pug/runtime"
)

type Options struct {
	// Setting if pretty printing is enabled.
	// Pretty printing ensures that the output html is properly indented and in human readable form.
	// If disabled, produced HTML is compact. This might be more suitable in production environments.
	// Default: false
	PrettyPrint bool

	Dir compiler.Dir
}

func newContext(dir compiler.Dir, options ...Options) compiler.Context {
	opt := Options{}

	if len(options) > 0 {
		opt = options[0]
	}

	indentString := ""
	if opt.PrettyPrint {
		indentString = "  "
	}

	if opt.Dir != nil {
		dir = opt.Dir
	}

	context := compiler.NewContext(dir, indentString)

	return context
}

func compileTemplate(name string, tplstring string) (*template.Template, error) {
	return template.New(name).Funcs(runtime.FuncMap).Parse(tplstring)
}

func CompileFile(filename string, options ...Options) (*template.Template, error) {
	ctx := newContext(compiler.FsDir("."), options...)
	if tplstring, err := ctx.CompileFile(filename); err != nil {
		return nil, err
	} else {
		return compileTemplate(filename, tplstring)
	}
}

func CompileString(input string, options ...Options) (*template.Template, error) {
	ctx := newContext(compiler.StringInputDir(input), options...)
	if tplstring, err := ctx.CompileFile(""); err != nil {
		return nil, err
	} else {
		return compileTemplate("", tplstring)
	}
}

func ParseFile(filename string, options ...Options) (string, error) {
	return newContext(compiler.FsDir("."), options...).CompileFile(filename)
}

func ParseString(input string, options ...Options) (string, error) {
	return newContext(compiler.StringInputDir(input), options...).CompileFile("")
}
