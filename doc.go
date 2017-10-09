/*
Package pug.go is an elegant templating engine for Go Programming Language.
It is a port of Pug template engine, previously known as Jade.

Pug.go compiles .pug templates to standard go templates (https://golang.org/pkg/html/template/) and returns a `*template.Template` instance.

While there is no JavaScript environment present, Pug.go provides basic expression support over go template syntax. Such as `a(href="/user/" + UserId)` would concatenate two strings. You can use arithmetic, logical and comparison operators as well as ternery if operator.

Please check *Pug Language Reference* for details: https://pugjs.org/api/getting-started.html.

Differences between Pug and Pug.go (items with checkboxes are planned, just not present yet)

- [ ] Multiline attributes are not supported
- [ ] `&attributes` syntax is not supported
- [ ] `case` statement is not supported
- [ ] Filters are not supported
- [ ] Mixin rest arguments are not supported.
- `while` loops are not supported as Go templates do not provide it. We could use recursive templates or channel range loops etc but that would be unnecessary complexity.
- Unbuffered code blocks are not possible as we don't have a JS environment. However it is possible to define variables using `- var x = "foo"` syntax as an exception.

Apart from these missing features, everything in the language reference should be supported.

License

(The MIT License)

Copyright (c) 2017 Ekin Koc <ekin@eknkc.com>

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the 'Software'), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/
package pug
