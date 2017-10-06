

# pug
`import "github.com/eknkc/pug"`

* [Overview](#pkg-overview)
* [Index](#pkg-index)
* [Subdirectories](#pkg-subdirectories)

## <a name="pkg-overview">Overview</a>
Package pug.go is an elegant templating engine for Go Programming Language.
It is a port of Pug template engine, previously known as Jade.

### Tags
A tag is simply a word:


	html

is converted to


	<html></html>

It is possible to add ID and CLASS attributes to tags:


	div#main
	span.time

are converted to


	<div id="main"></div>
	<span class="time"></span>

Any arbitrary attribute name / value pair can be added this way:


	a(href="<a href="http://www.google.com">http://www.google.com</a>")

You can mix multiple attributes together


	a#someid(href="/" title="Main Page").main.link Click Link

gets converted to


	<a id="someid" class="main link" href="/" title="Main Page">Click Link</a>

### Doctypes
To add a doctype, use `doctype` keyword:


	doctype transitional
	// <!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Transitional//EN" "<a href="http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd">http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd</a>">
	
	doctype html
	// <!DOCTYPE html>

Available options: `html`, `xml`, `transitional`, `strict`, `frameset`, `1.1`, `basic`, `mobile`, `plist`

### Tag Content
For single line tag text, you can just append the text after tag name:


	p Testing!

would yield


	<p>Testing!</p>

For multi line tag text, or nested tags, use indentation:


	html
	    head
	        title Page Title
	    body
	        div#content
	            p
	                | This is a long page content
	                | These lines are all part of the parent p
	
	                a(href="/") Go To Main Page

### Data
Input template data can be reached by key names directly. For example, assuming the template has been
executed with following JSON data:


	{
	    "Name": "Ekin",
	    "LastName": "Koc",
	    "Repositories": [
	        "pug",
	        "dateformat"
	    ],
	    "Avatar": "/images/ekin.jpg",
	    "Friends": 17
	}

It is possible to interpolate fields using `#{}`


	p Welcome #{Name}!

would print


	<p>Welcome Ekin!</p>

Attributes can have field names as well


	a(title=Name href="/ekin.koc")

would print


	<a title="Ekin" href="/ekin.koc"></a>

### Expressions
Pug can expand basic expressions. For example, it is possible to concatenate strings with + operator:


	p Welcome #{Name + " " + LastName}

Arithmetic expressions are also supported:


	p You need #{50 - Friends} more friends to reach 50!

Expressions can be used within attributes


	img(alt=Name + " " + LastName src=Avatar)

### Variables
It is possible to define dynamic variables within templates


	div
	    fullname = Name + " " + LastName
	    p Welcome #{fullname}

### Conditions
For conditional blocks, it is possible to use `if <expression>`


	div
	    if Friends > 10
	        p You have more than 10 friends
	    else if Friends > 5
	        p You have more than 5 friends
	    else
	        p You need more friends

Again, it is possible to use arithmetic and boolean operators


	div
	    if Name == "Ekin" && LastName == "Koc"
	        p Hey! I know you..

### Iterations
It is possible to iterate over arrays and maps using `each`:


	each repo in Repositories
	    p #{repo}

would print


	p pug
	p dateformat

It is also possible to iterate over values and indexes at the same time


	each i, repo in Repositories
	    p(class=i % 2 == 0 ? "even" : "odd") #{repo}

### Includes
A template can include other templates using `include`:


	a.pug
	    p this is template a
	
	b.pug
	    p this is template b
	
	c.pug
	    div
	        include a
	        include b

gets compiled to


	div
	    p this is template a
	    p this is template b

### Inheritance
A template can inherit other templates. In order to inherit another template, an `extends` keyword should be used.
Parent template can define several named blocks and child template can modify the blocks.


	master.pug
	    doctype html
	    html
	        head
	            block meta
	                meta(name="description" content="This is a great website")
	
	            title
	                block title
	                    | Default title
	        body
	            block content
	
	subpage.pug
	    extends master
	
	    block title
	        | Some sub page!
	
	    block append meta
	        // This will be added after the description meta tag. It is also possible
	        // to prepend something to an existing block
	        meta(name="keywords" content="foo bar")
	
	    block content
	        div#main
	            p Some content here

License
(The MIT License)

Copyright (c) 2017 Ekin Koc <ekin@eknkc.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.




## <a name="pkg-index">Index</a>
* [func CompileFile(filename string, options ...Options) (*template.Template, error)](#CompileFile)
* [func CompileString(input string, options ...Options) (*template.Template, error)](#CompileString)
* [func ParseFile(filename string, options ...Options) (string, error)](#ParseFile)
* [func ParseString(input string, options ...Options) (string, error)](#ParseString)
* [type Options](#Options)


#### <a name="pkg-files">Package files</a>
[doc.go](/src/github.com/eknkc/pug/doc.go) [pug.go](/src/github.com/eknkc/pug/pug.go) 





## <a name="CompileFile">func</a> [CompileFile](/src/target/pug.go?s=900:981#L46)
``` go
func CompileFile(filename string, options ...Options) (*template.Template, error)
```


## <a name="CompileString">func</a> [CompileString](/src/target/pug.go?s=1178:1258#L55)
``` go
func CompileString(input string, options ...Options) (*template.Template, error)
```


## <a name="ParseFile">func</a> [ParseFile](/src/target/pug.go?s=1454:1521#L64)
``` go
func ParseFile(filename string, options ...Options) (string, error)
```


## <a name="ParseString">func</a> [ParseString](/src/target/pug.go?s=1601:1667#L68)
``` go
func ParseString(input string, options ...Options) (string, error)
```



## <a name="Options">type</a> [Options](/src/target/pug.go?s=108:427#L10)
``` go
type Options struct {
    // Setting if pretty printing is enabled.
    // Pretty printing ensures that the output html is properly indented and in human readable form.
    // If disabled, produced HTML is compact. This might be more suitable in production environments.
    // Default: false
    PrettyPrint bool

    Dir compiler.Dir
}
```













- - -
Generated by [godoc2md](http://godoc.org/github.com/davecheney/godoc2md)
