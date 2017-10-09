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
*/
package pug
