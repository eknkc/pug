package pug

import (
	"bytes"
	"strings"
	"testing"
)

func Test_Doctype(t *testing.T) {
	res, err := run(`doctype html`, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<!DOCTYPE html>`, t)
	}
}

func Test_Nesting(t *testing.T) {
	res, err := run(`html
						head
							title
						body`, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<html><head><title></title></head><body></body></html>`, t)
	}
}

func Test_Id(t *testing.T) {
	res, err := run(`div#test`, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<div id="test"></div>`, t)
	}
}

func Test_Class(t *testing.T) {
	res, err := run(`div.test`, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<div class="test"></div>`, t)
	}
}

func Test_MultiClass(t *testing.T) {
	res, err := run(`div.test.foo.bar(class="baz")`, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<div class="test foo bar baz"></div>`, t)
	}
}

func Test_Attribute(t *testing.T) {
	res, err := run(`
div(name="Test" @foo.bar="baz").testclass
	p(style="text-align: center; color: maroon")`, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<div name="Test" @foo.bar="baz" class="testclass"><p style="text-align: center; color: maroon"></p></div>`, t)
	}
}

func Test_EmptyAttribute(t *testing.T) {
	res, err := run(`div(name)`, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<div name></div>`, t)
	}
}

func Test_Empty(t *testing.T) {
	res, err := run(``, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, ``, t)
	}
}

func Test_ArithmeticExpression(t *testing.T) {
	res, err := run(`| #{A + B * C}`, map[string]int{"A": 2, "B": 3, "C": 4})

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `14`, t)
	}
}

func Test_BooleanExpression(t *testing.T) {
	res, err := run(`| #{C - A < B}`, map[string]int{"A": 2, "B": 3, "C": 4})

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `true`, t)
	}
}

func Test_TerneryExpression(t *testing.T) {
	res, err := run(`| #{ B > A ? A > B ? "x" : "y" : "z" }`, map[string]int{"A": 2, "B": 3})

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `y`, t)
	}
}

func Test_TerneryClass(t *testing.T) {
	res, err := run(`
each item, i in Items
	p(class=i % 2 == 0 ? "even" : "odd") #{item}`, testStruct{Items: []string{"test1", "test2"}})

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<p class="even">test1</p><p class="odd">test2</p>`, t)
	}
}

func Test_NilClass(t *testing.T) {
	res, err := run(`p(class=nil)`, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<p class=""></p>`, t)
	}
}

func Test_MapAccess(t *testing.T) {
	res, err := run(`p #{a.b().c}`, map[string]interface{}{
		"a": map[string]interface{}{
			"b": func() interface{} {
				return map[string]interface{}{
					"c": "d",
				}
			},
		},
	})

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<p>d</p>`, t)
	}
}

func Test_Dollar_In_TagAttributes(t *testing.T) {
	res, err := run(`input(placeholder="$ per "+kwh)`, map[string]interface{}{
		"kwh": "kWh",
	})

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<input placeholder="$ per kWh" />`, t)
	}
}

func Test_MixinBasic(t *testing.T) {
	res, err := run(`
mixin test()
	p #{Key}

+test()
`, testStruct{Key: "value"})

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<p>value</p>`, t)
	}
}

func Test_MixinWithArgs(t *testing.T) {
	res, err := run(`
mixin test(arg, arg2)
	p #{Key} #{arg} #{arg2}

+test(15, 1+1)
`, testStruct{Key: "value"})

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<p>value 15 2</p>`, t)
	}
}

func Test_Each(t *testing.T) {
	res, err := run(`
each v in Items
		p #{v}
		`, testStruct{
		Items: []string{"t1", "t2"},
	})

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<p>t1</p><p>t2</p>`, t)
	}
}

func Test_Assignment(t *testing.T) {
	res, err := run(`
vrb = "test"
vrb = "test2"
p #{vrb}
`, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<p>test2</p>`, t)
	}
}

func Test_Block(t *testing.T) {
	res, err := run(`
block deneme
		p Test
		`, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, `<p>Test</p>`, t)
	}
}

func Test_RawText(t *testing.T) {
	res, err := run(`
style.
  body{ color: red }
p a
`, nil)

	if err != nil {
		t.Fatal(err.Error())
	} else {
		expect(res, "<style>  body{ color: red }\n</style><p>a</p>", t)
	}
}

func Test_Import(t *testing.T) {
	tpl, err := CompileFile("examples/import/import.pug")

	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}

	if err := tpl.Execute(buf, nil); err != nil {
		t.Fatal(err)
	}

	expect(string(buf.Bytes()), "<p>Main<p>import1</p><p>import2</p></p><p>import2</p>", t)
}

func Test_Extend(t *testing.T) {
	tpl, err := CompileFile("examples/extend/extend.pug")

	if err != nil {
		t.Fatal(err)
	}

	buf := &bytes.Buffer{}

	if err := tpl.Execute(buf, nil); err != nil {
		t.Fatal(err)
	}

	expect(string(buf.Bytes()), "<body><p>extend-test1</p><p>mid-test2</p><p>base-test3</p><p>extend-test3-append</p></body>", t)
}

func Benchmark_Parse(b *testing.B) {
	code := `
	!!! 5
	html
		head
			title Test Title
		body
			nav#mainNav[data-foo="bar"]
			div#content
				div.left
				div.center
					block center
						p Main Content
							.long ? somevar && someothervar
				div.right`

	for i := 0; i < b.N; i++ {
		CompileString(code)
	}
}

func expect(cur, expected string, t *testing.T) {
	if cur != expected {
		t.Fatalf("Expected {%s} got {%s}.", expected, cur)
	}
}

func run(tpl string, data interface{}) (string, error) {
	t, err := CompileString(tpl)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err = t.Execute(&buf, data); err != nil {
		return "", err
	}
	return strings.TrimSpace(buf.String()), nil
}

type testStruct struct {
	Key   string
	Items []string
}
