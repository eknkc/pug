package compiler

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func toSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func pos(p position) Position {
	return Position{
		Line:   p.line,
		Col:    p.col,
		Offset: p.offset,
	}
}

func binary(first, rest interface{}, curpos position) (Expression, error) {
	restElem := toSlice(rest)

	if len(restElem) == 0 {
		return first.(Expression), nil
	}

	var cur Expression = first.(Expression)

	for _, x := range restElem {
		elem := toSlice(x)
		cur = &BinaryExpression{X: cur, Y: elem[3].(Expression), Op: elem[1].(string), GraphNode: NewNode(pos(curpos))}
	}

	return cur, nil
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Input",
			pos:  position{line: 37, col: 1, offset: 708},
			expr: &actionExpr{
				pos: position{line: 37, col: 10, offset: 717},
				run: (*parser).callonInput1,
				expr: &seqExpr{
					pos: position{line: 37, col: 10, offset: 717},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 37, col: 10, offset: 717},
							label: "l",
							expr: &ruleRefExpr{
								pos:  position{line: 37, col: 12, offset: 719},
								name: "List",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 37, col: 17, offset: 724},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "List",
			pos:  position{line: 41, col: 1, offset: 802},
			expr: &choiceExpr{
				pos: position{line: 41, col: 9, offset: 810},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 41, col: 9, offset: 810},
						run: (*parser).callonList2,
						expr: &seqExpr{
							pos: position{line: 41, col: 9, offset: 810},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 41, col: 9, offset: 810},
									label: "node",
									expr: &ruleRefExpr{
										pos:  position{line: 41, col: 14, offset: 815},
										name: "ListNode",
									},
								},
								&labeledExpr{
									pos:   position{line: 41, col: 23, offset: 824},
									label: "list",
									expr: &ruleRefExpr{
										pos:  position{line: 41, col: 28, offset: 829},
										name: "List",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 49, col: 5, offset: 1015},
						run: (*parser).callonList8,
						expr: &andExpr{
							pos: position{line: 49, col: 5, offset: 1015},
							expr: &choiceExpr{
								pos: position{line: 49, col: 7, offset: 1017},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 49, col: 7, offset: 1017},
										name: "Outdent",
									},
									&ruleRefExpr{
										pos:  position{line: 49, col: 17, offset: 1027},
										name: "EOF",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "IndentedList",
			pos:  position{line: 53, col: 1, offset: 1089},
			expr: &actionExpr{
				pos: position{line: 53, col: 17, offset: 1105},
				run: (*parser).callonIndentedList1,
				expr: &seqExpr{
					pos: position{line: 53, col: 17, offset: 1105},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 53, col: 17, offset: 1105},
							name: "Indent",
						},
						&labeledExpr{
							pos:   position{line: 53, col: 24, offset: 1112},
							label: "list",
							expr: &ruleRefExpr{
								pos:  position{line: 53, col: 29, offset: 1117},
								name: "List",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 53, col: 34, offset: 1122},
							name: "Outdent",
						},
					},
				},
			},
		},
		{
			name: "IndentedRawText",
			pos:  position{line: 57, col: 1, offset: 1154},
			expr: &actionExpr{
				pos: position{line: 57, col: 20, offset: 1173},
				run: (*parser).callonIndentedRawText1,
				expr: &seqExpr{
					pos: position{line: 57, col: 20, offset: 1173},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 57, col: 20, offset: 1173},
							name: "Indent",
						},
						&labeledExpr{
							pos:   position{line: 57, col: 27, offset: 1180},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 57, col: 29, offset: 1182},
								name: "RawText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 57, col: 37, offset: 1190},
							name: "Outdent",
						},
					},
				},
			},
		},
		{
			name: "RawText",
			pos:  position{line: 61, col: 1, offset: 1274},
			expr: &choiceExpr{
				pos: position{line: 61, col: 12, offset: 1285},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 61, col: 12, offset: 1285},
						run: (*parser).callonRawText2,
						expr: &seqExpr{
							pos: position{line: 61, col: 12, offset: 1285},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 61, col: 12, offset: 1285},
									name: "Indent",
								},
								&labeledExpr{
									pos:   position{line: 61, col: 19, offset: 1292},
									label: "rt",
									expr: &ruleRefExpr{
										pos:  position{line: 61, col: 22, offset: 1295},
										name: "RawText",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 61, col: 30, offset: 1303},
									name: "Outdent",
								},
								&labeledExpr{
									pos:   position{line: 61, col: 38, offset: 1311},
									label: "tail",
									expr: &ruleRefExpr{
										pos:  position{line: 61, col: 43, offset: 1316},
										name: "RawText",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 63, col: 5, offset: 1372},
						run: (*parser).callonRawText10,
						expr: &choiceExpr{
							pos: position{line: 63, col: 6, offset: 1373},
							alternatives: []interface{}{
								&andExpr{
									pos: position{line: 63, col: 6, offset: 1373},
									expr: &ruleRefExpr{
										pos:  position{line: 63, col: 7, offset: 1374},
										name: "Outdent",
									},
								},
								&andExpr{
									pos: position{line: 63, col: 17, offset: 1384},
									expr: &ruleRefExpr{
										pos:  position{line: 63, col: 18, offset: 1385},
										name: "EOF",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 65, col: 5, offset: 1413},
						run: (*parser).callonRawText16,
						expr: &seqExpr{
							pos: position{line: 65, col: 5, offset: 1413},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 65, col: 5, offset: 1413},
									label: "head",
									expr: &anyMatcher{
										line: 65, col: 10, offset: 1418,
									},
								},
								&labeledExpr{
									pos:   position{line: 65, col: 12, offset: 1420},
									label: "tail",
									expr: &ruleRefExpr{
										pos:  position{line: 65, col: 17, offset: 1425},
										name: "RawText",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name:        "ListNode",
			displayName: "\"listnode\"",
			pos:         position{line: 69, col: 1, offset: 1490},
			expr: &choiceExpr{
				pos: position{line: 70, col: 3, offset: 1515},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 70, col: 3, offset: 1515},
						name: "Comment",
					},
					&ruleRefExpr{
						pos:  position{line: 71, col: 3, offset: 1527},
						name: "Import",
					},
					&ruleRefExpr{
						pos:  position{line: 72, col: 3, offset: 1538},
						name: "Extend",
					},
					&ruleRefExpr{
						pos:  position{line: 73, col: 3, offset: 1549},
						name: "PipeText",
					},
					&ruleRefExpr{
						pos:  position{line: 74, col: 3, offset: 1562},
						name: "PipeExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 75, col: 3, offset: 1581},
						name: "If",
					},
					&ruleRefExpr{
						pos:  position{line: 76, col: 3, offset: 1588},
						name: "Unless",
					},
					&ruleRefExpr{
						pos:  position{line: 77, col: 3, offset: 1599},
						name: "Each",
					},
					&ruleRefExpr{
						pos:  position{line: 78, col: 3, offset: 1608},
						name: "DocType",
					},
					&ruleRefExpr{
						pos:  position{line: 79, col: 3, offset: 1620},
						name: "Mixin",
					},
					&ruleRefExpr{
						pos:  position{line: 80, col: 3, offset: 1630},
						name: "MixinCall",
					},
					&ruleRefExpr{
						pos:  position{line: 81, col: 3, offset: 1644},
						name: "Assignment",
					},
					&ruleRefExpr{
						pos:  position{line: 82, col: 3, offset: 1659},
						name: "Block",
					},
					&ruleRefExpr{
						pos:  position{line: 83, col: 3, offset: 1669},
						name: "Tag",
					},
					&actionExpr{
						pos: position{line: 84, col: 3, offset: 1677},
						run: (*parser).callonListNode16,
						expr: &seqExpr{
							pos: position{line: 84, col: 4, offset: 1678},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 84, col: 4, offset: 1678},
									name: "_",
								},
								&charClassMatcher{
									pos:        position{line: 84, col: 6, offset: 1680},
									val:        "[\\n]",
									chars:      []rune{'\n'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
					},
				},
			},
		},
		{
			name:        "DocType",
			displayName: "\"doctype\"",
			pos:         position{line: 87, col: 1, offset: 1715},
			expr: &actionExpr{
				pos: position{line: 87, col: 21, offset: 1735},
				run: (*parser).callonDocType1,
				expr: &seqExpr{
					pos: position{line: 87, col: 21, offset: 1735},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 87, col: 21, offset: 1735},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 87, col: 23, offset: 1737},
							val:        "doctype",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 87, col: 33, offset: 1747},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 87, col: 35, offset: 1749},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 87, col: 39, offset: 1753},
								name: "LineText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 87, col: 48, offset: 1762},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Tag",
			pos:  position{line: 93, col: 1, offset: 1855},
			expr: &actionExpr{
				pos: position{line: 93, col: 8, offset: 1862},
				run: (*parser).callonTag1,
				expr: &seqExpr{
					pos: position{line: 93, col: 8, offset: 1862},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 93, col: 8, offset: 1862},
							label: "tag",
							expr: &ruleRefExpr{
								pos:  position{line: 93, col: 12, offset: 1866},
								name: "TagHeader",
							},
						},
						&labeledExpr{
							pos:   position{line: 93, col: 22, offset: 1876},
							label: "list",
							expr: &zeroOrOneExpr{
								pos: position{line: 93, col: 27, offset: 1881},
								expr: &ruleRefExpr{
									pos:  position{line: 93, col: 27, offset: 1881},
									name: "IndentedList",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TagHeader",
			pos:  position{line: 107, col: 1, offset: 2158},
			expr: &choiceExpr{
				pos: position{line: 107, col: 14, offset: 2171},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 107, col: 14, offset: 2171},
						run: (*parser).callonTagHeader2,
						expr: &seqExpr{
							pos: position{line: 107, col: 14, offset: 2171},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 107, col: 14, offset: 2171},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 107, col: 16, offset: 2173},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 107, col: 21, offset: 2178},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 107, col: 29, offset: 2186},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 107, col: 35, offset: 2192},
										expr: &ruleRefExpr{
											pos:  position{line: 107, col: 35, offset: 2192},
											name: "TagAttributes",
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 107, col: 50, offset: 2207},
									label: "selfClose",
									expr: &zeroOrOneExpr{
										pos: position{line: 107, col: 60, offset: 2217},
										expr: &litMatcher{
											pos:        position{line: 107, col: 60, offset: 2217},
											val:        "/",
											ignoreCase: false,
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 107, col: 65, offset: 2222},
									label: "tl",
									expr: &zeroOrOneExpr{
										pos: position{line: 107, col: 68, offset: 2225},
										expr: &seqExpr{
											pos: position{line: 107, col: 69, offset: 2226},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 107, col: 69, offset: 2226},
													name: "__",
												},
												&zeroOrOneExpr{
													pos: position{line: 107, col: 72, offset: 2229},
													expr: &ruleRefExpr{
														pos:  position{line: 107, col: 72, offset: 2229},
														name: "TextList",
													},
												},
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 107, col: 84, offset: 2241},
									name: "NL",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 119, col: 5, offset: 2517},
						run: (*parser).callonTagHeader20,
						expr: &seqExpr{
							pos: position{line: 119, col: 5, offset: 2517},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 119, col: 5, offset: 2517},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 119, col: 7, offset: 2519},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 119, col: 12, offset: 2524},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 119, col: 20, offset: 2532},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 119, col: 26, offset: 2538},
										expr: &ruleRefExpr{
											pos:  position{line: 119, col: 26, offset: 2538},
											name: "TagAttributes",
										},
									},
								},
								&litMatcher{
									pos:        position{line: 119, col: 41, offset: 2553},
									val:        ".",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 119, col: 45, offset: 2557},
									name: "NL",
								},
								&labeledExpr{
									pos:   position{line: 119, col: 48, offset: 2560},
									label: "text",
									expr: &zeroOrOneExpr{
										pos: position{line: 119, col: 53, offset: 2565},
										expr: &ruleRefExpr{
											pos:  position{line: 119, col: 53, offset: 2565},
											name: "IndentedRawText",
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 128, col: 5, offset: 2791},
						run: (*parser).callonTagHeader33,
						expr: &seqExpr{
							pos: position{line: 128, col: 5, offset: 2791},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 128, col: 5, offset: 2791},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 128, col: 7, offset: 2793},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 128, col: 12, offset: 2798},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 128, col: 20, offset: 2806},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 128, col: 26, offset: 2812},
										expr: &ruleRefExpr{
											pos:  position{line: 128, col: 26, offset: 2812},
											name: "TagAttributes",
										},
									},
								},
								&litMatcher{
									pos:        position{line: 128, col: 41, offset: 2827},
									val:        ":",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 128, col: 45, offset: 2831},
									label: "block",
									expr: &ruleRefExpr{
										pos:  position{line: 128, col: 51, offset: 2837},
										name: "ListNode",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 137, col: 5, offset: 3056},
						run: (*parser).callonTagHeader44,
						expr: &seqExpr{
							pos: position{line: 137, col: 5, offset: 3056},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 137, col: 5, offset: 3056},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 137, col: 7, offset: 3058},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 137, col: 12, offset: 3063},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 137, col: 20, offset: 3071},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 137, col: 26, offset: 3077},
										expr: &ruleRefExpr{
											pos:  position{line: 137, col: 26, offset: 3077},
											name: "TagAttributes",
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 137, col: 41, offset: 3092},
									label: "unescaped",
									expr: &zeroOrOneExpr{
										pos: position{line: 137, col: 51, offset: 3102},
										expr: &litMatcher{
											pos:        position{line: 137, col: 51, offset: 3102},
											val:        "!",
											ignoreCase: false,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 137, col: 56, offset: 3107},
									val:        "=",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 137, col: 60, offset: 3111},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 137, col: 63, offset: 3114},
									label: "expr",
									expr: &zeroOrOneExpr{
										pos: position{line: 137, col: 68, offset: 3119},
										expr: &ruleRefExpr{
											pos:  position{line: 137, col: 68, offset: 3119},
											name: "Expression",
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 137, col: 80, offset: 3131},
									name: "NL",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TagName",
			pos:  position{line: 152, col: 1, offset: 3536},
			expr: &choiceExpr{
				pos: position{line: 152, col: 12, offset: 3547},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 152, col: 12, offset: 3547},
						run: (*parser).callonTagName2,
						expr: &seqExpr{
							pos: position{line: 152, col: 12, offset: 3547},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 152, col: 12, offset: 3547},
									val:        "[_a-zA-Z]",
									chars:      []rune{'_'},
									ranges:     []rune{'a', 'z', 'A', 'Z'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 152, col: 22, offset: 3557},
									expr: &charClassMatcher{
										pos:        position{line: 152, col: 22, offset: 3557},
										val:        "[_-:a-zA-Z0-9]",
										ranges:     []rune{'_', ':', 'a', 'z', 'A', 'Z', '0', '9'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
					&andExpr{
						pos: position{line: 154, col: 5, offset: 3608},
						expr: &ruleRefExpr{
							pos:  position{line: 154, col: 6, offset: 3609},
							name: "TagAttributeClass",
						},
					},
					&actionExpr{
						pos: position{line: 154, col: 26, offset: 3629},
						run: (*parser).callonTagName9,
						expr: &andExpr{
							pos: position{line: 154, col: 26, offset: 3629},
							expr: &ruleRefExpr{
								pos:  position{line: 154, col: 27, offset: 3630},
								name: "TagAttributeID",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributes",
			pos:  position{line: 158, col: 1, offset: 3670},
			expr: &choiceExpr{
				pos: position{line: 158, col: 18, offset: 3687},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 158, col: 18, offset: 3687},
						run: (*parser).callonTagAttributes2,
						expr: &seqExpr{
							pos: position{line: 158, col: 18, offset: 3687},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 158, col: 18, offset: 3687},
									label: "head",
									expr: &choiceExpr{
										pos: position{line: 158, col: 24, offset: 3693},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 158, col: 24, offset: 3693},
												name: "TagAttributeClass",
											},
											&ruleRefExpr{
												pos:  position{line: 158, col: 44, offset: 3713},
												name: "TagAttributeID",
											},
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 158, col: 60, offset: 3729},
									label: "tail",
									expr: &zeroOrOneExpr{
										pos: position{line: 158, col: 65, offset: 3734},
										expr: &ruleRefExpr{
											pos:  position{line: 158, col: 65, offset: 3734},
											name: "TagAttributes",
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 166, col: 5, offset: 3899},
						run: (*parser).callonTagAttributes11,
						expr: &seqExpr{
							pos: position{line: 166, col: 5, offset: 3899},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 166, col: 5, offset: 3899},
									label: "head",
									expr: &seqExpr{
										pos: position{line: 166, col: 11, offset: 3905},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 166, col: 11, offset: 3905},
												val:        "(",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 166, col: 15, offset: 3909},
												name: "_",
											},
											&seqExpr{
												pos: position{line: 166, col: 18, offset: 3912},
												exprs: []interface{}{
													&ruleRefExpr{
														pos:  position{line: 166, col: 18, offset: 3912},
														name: "TagAttribute",
													},
													&zeroOrMoreExpr{
														pos: position{line: 166, col: 31, offset: 3925},
														expr: &seqExpr{
															pos: position{line: 166, col: 32, offset: 3926},
															exprs: []interface{}{
																&choiceExpr{
																	pos: position{line: 166, col: 33, offset: 3927},
																	alternatives: []interface{}{
																		&ruleRefExpr{
																			pos:  position{line: 166, col: 33, offset: 3927},
																			name: "__",
																		},
																		&seqExpr{
																			pos: position{line: 166, col: 39, offset: 3933},
																			exprs: []interface{}{
																				&ruleRefExpr{
																					pos:  position{line: 166, col: 39, offset: 3933},
																					name: "_",
																				},
																				&litMatcher{
																					pos:        position{line: 166, col: 41, offset: 3935},
																					val:        ",",
																					ignoreCase: false,
																				},
																				&ruleRefExpr{
																					pos:  position{line: 166, col: 45, offset: 3939},
																					name: "_",
																				},
																			},
																		},
																	},
																},
																&ruleRefExpr{
																	pos:  position{line: 166, col: 49, offset: 3943},
																	name: "TagAttribute",
																},
															},
														},
													},
												},
											},
											&ruleRefExpr{
												pos:  position{line: 166, col: 65, offset: 3959},
												name: "_",
											},
											&litMatcher{
												pos:        position{line: 166, col: 67, offset: 3961},
												val:        ")",
												ignoreCase: false,
											},
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 166, col: 72, offset: 3966},
									label: "tail",
									expr: &zeroOrOneExpr{
										pos: position{line: 166, col: 77, offset: 3971},
										expr: &ruleRefExpr{
											pos:  position{line: 166, col: 77, offset: 3971},
											name: "TagAttributes",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributeClass",
			pos:  position{line: 190, col: 1, offset: 4417},
			expr: &actionExpr{
				pos: position{line: 190, col: 22, offset: 4438},
				run: (*parser).callonTagAttributeClass1,
				expr: &seqExpr{
					pos: position{line: 190, col: 22, offset: 4438},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 190, col: 22, offset: 4438},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 190, col: 26, offset: 4442},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 190, col: 31, offset: 4447},
								name: "ClassName",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributeID",
			pos:  position{line: 194, col: 1, offset: 4596},
			expr: &actionExpr{
				pos: position{line: 194, col: 19, offset: 4614},
				run: (*parser).callonTagAttributeID1,
				expr: &seqExpr{
					pos: position{line: 194, col: 19, offset: 4614},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 194, col: 19, offset: 4614},
							val:        "#",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 194, col: 23, offset: 4618},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 194, col: 28, offset: 4623},
								name: "IdName",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttribute",
			pos:  position{line: 198, col: 1, offset: 4766},
			expr: &choiceExpr{
				pos: position{line: 198, col: 17, offset: 4782},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 198, col: 17, offset: 4782},
						run: (*parser).callonTagAttribute2,
						expr: &seqExpr{
							pos: position{line: 198, col: 17, offset: 4782},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 198, col: 17, offset: 4782},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 198, col: 22, offset: 4787},
										name: "TagAttributeName",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 198, col: 39, offset: 4804},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 198, col: 41, offset: 4806},
									val:        "=",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 198, col: 45, offset: 4810},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 198, col: 47, offset: 4812},
									label: "value",
									expr: &ruleRefExpr{
										pos:  position{line: 198, col: 53, offset: 4818},
										name: "Expression",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 200, col: 5, offset: 4954},
						run: (*parser).callonTagAttribute11,
						expr: &seqExpr{
							pos: position{line: 200, col: 5, offset: 4954},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 200, col: 5, offset: 4954},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 200, col: 10, offset: 4959},
										name: "TagAttributeName",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 200, col: 27, offset: 4976},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 200, col: 29, offset: 4978},
									val:        "!=",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 200, col: 34, offset: 4983},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 200, col: 36, offset: 4985},
									label: "value",
									expr: &ruleRefExpr{
										pos:  position{line: 200, col: 42, offset: 4991},
										name: "Expression",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 202, col: 5, offset: 5144},
						run: (*parser).callonTagAttribute20,
						expr: &labeledExpr{
							pos:   position{line: 202, col: 5, offset: 5144},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 202, col: 10, offset: 5149},
								name: "TagAttributeName",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributeName",
			pos:  position{line: 206, col: 1, offset: 5263},
			expr: &choiceExpr{
				pos: position{line: 206, col: 21, offset: 5283},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 206, col: 21, offset: 5283},
						run: (*parser).callonTagAttributeName2,
						expr: &seqExpr{
							pos: position{line: 206, col: 21, offset: 5283},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 206, col: 21, offset: 5283},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 206, col: 25, offset: 5287},
									label: "tn",
									expr: &ruleRefExpr{
										pos:  position{line: 206, col: 28, offset: 5290},
										name: "TagAttributeNameLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 206, col: 52, offset: 5314},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 208, col: 5, offset: 5341},
						run: (*parser).callonTagAttributeName8,
						expr: &seqExpr{
							pos: position{line: 208, col: 5, offset: 5341},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 208, col: 5, offset: 5341},
									val:        "[",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 208, col: 9, offset: 5345},
									label: "tn",
									expr: &ruleRefExpr{
										pos:  position{line: 208, col: 12, offset: 5348},
										name: "TagAttributeNameLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 208, col: 36, offset: 5372},
									val:        "]",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 210, col: 5, offset: 5399},
						run: (*parser).callonTagAttributeName14,
						expr: &labeledExpr{
							pos:   position{line: 210, col: 5, offset: 5399},
							label: "tn",
							expr: &ruleRefExpr{
								pos:  position{line: 210, col: 8, offset: 5402},
								name: "TagAttributeNameLiteral",
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 212, col: 5, offset: 5449},
						name: "String",
					},
				},
			},
		},
		{
			name: "ClassName",
			pos:  position{line: 214, col: 1, offset: 5457},
			expr: &actionExpr{
				pos: position{line: 214, col: 14, offset: 5470},
				run: (*parser).callonClassName1,
				expr: &oneOrMoreExpr{
					pos: position{line: 214, col: 14, offset: 5470},
					expr: &charClassMatcher{
						pos:        position{line: 214, col: 14, offset: 5470},
						val:        "[_-:a-zA-Z0-9]",
						ranges:     []rune{'_', ':', 'a', 'z', 'A', 'Z', '0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "IdName",
			pos:  position{line: 218, col: 1, offset: 5520},
			expr: &actionExpr{
				pos: position{line: 218, col: 11, offset: 5530},
				run: (*parser).callonIdName1,
				expr: &oneOrMoreExpr{
					pos: position{line: 218, col: 11, offset: 5530},
					expr: &charClassMatcher{
						pos:        position{line: 218, col: 11, offset: 5530},
						val:        "[_-:a-zA-Z0-9]",
						ranges:     []rune{'_', ':', 'a', 'z', 'A', 'Z', '0', '9'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name: "TagAttributeNameLiteral",
			pos:  position{line: 222, col: 1, offset: 5580},
			expr: &actionExpr{
				pos: position{line: 222, col: 28, offset: 5607},
				run: (*parser).callonTagAttributeNameLiteral1,
				expr: &seqExpr{
					pos: position{line: 222, col: 28, offset: 5607},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 222, col: 28, offset: 5607},
							val:        "[@_a-zA-Z]",
							chars:      []rune{'@', '_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 222, col: 39, offset: 5618},
							expr: &charClassMatcher{
								pos:        position{line: 222, col: 39, offset: 5618},
								val:        "[._-:a-zA-Z0-9]",
								chars:      []rune{'.'},
								ranges:     []rune{'_', ':', 'a', 'z', 'A', 'Z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "If",
			pos:  position{line: 227, col: 1, offset: 5678},
			expr: &actionExpr{
				pos: position{line: 227, col: 7, offset: 5684},
				run: (*parser).callonIf1,
				expr: &seqExpr{
					pos: position{line: 227, col: 7, offset: 5684},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 227, col: 7, offset: 5684},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 227, col: 9, offset: 5686},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 227, col: 14, offset: 5691},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 227, col: 17, offset: 5694},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 227, col: 22, offset: 5699},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 227, col: 33, offset: 5710},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 227, col: 35, offset: 5712},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 227, col: 38, offset: 5715},
							label: "block",
							expr: &ruleRefExpr{
								pos:  position{line: 227, col: 44, offset: 5721},
								name: "IndentedList",
							},
						},
						&labeledExpr{
							pos:   position{line: 227, col: 57, offset: 5734},
							label: "elseNode",
							expr: &zeroOrOneExpr{
								pos: position{line: 227, col: 66, offset: 5743},
								expr: &ruleRefExpr{
									pos:  position{line: 227, col: 66, offset: 5743},
									name: "Else",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Unless",
			pos:  position{line: 235, col: 1, offset: 5952},
			expr: &actionExpr{
				pos: position{line: 235, col: 11, offset: 5962},
				run: (*parser).callonUnless1,
				expr: &seqExpr{
					pos: position{line: 235, col: 11, offset: 5962},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 235, col: 11, offset: 5962},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 235, col: 13, offset: 5964},
							val:        "unless",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 235, col: 22, offset: 5973},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 235, col: 25, offset: 5976},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 235, col: 30, offset: 5981},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 235, col: 41, offset: 5992},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 235, col: 43, offset: 5994},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 235, col: 46, offset: 5997},
							label: "block",
							expr: &ruleRefExpr{
								pos:  position{line: 235, col: 52, offset: 6003},
								name: "IndentedList",
							},
						},
					},
				},
			},
		},
		{
			name: "Else",
			pos:  position{line: 244, col: 1, offset: 6199},
			expr: &choiceExpr{
				pos: position{line: 244, col: 9, offset: 6207},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 244, col: 9, offset: 6207},
						run: (*parser).callonElse2,
						expr: &seqExpr{
							pos: position{line: 244, col: 9, offset: 6207},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 244, col: 9, offset: 6207},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 244, col: 11, offset: 6209},
									val:        "else",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 244, col: 18, offset: 6216},
									label: "node",
									expr: &ruleRefExpr{
										pos:  position{line: 244, col: 23, offset: 6221},
										name: "If",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 246, col: 5, offset: 6249},
						run: (*parser).callonElse8,
						expr: &seqExpr{
							pos: position{line: 246, col: 5, offset: 6249},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 246, col: 5, offset: 6249},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 246, col: 7, offset: 6251},
									val:        "else",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 246, col: 14, offset: 6258},
									name: "_",
								},
								&ruleRefExpr{
									pos:  position{line: 246, col: 16, offset: 6260},
									name: "NL",
								},
								&labeledExpr{
									pos:   position{line: 246, col: 19, offset: 6263},
									label: "block",
									expr: &ruleRefExpr{
										pos:  position{line: 246, col: 25, offset: 6269},
										name: "IndentedList",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Each",
			pos:  position{line: 250, col: 1, offset: 6307},
			expr: &actionExpr{
				pos: position{line: 250, col: 9, offset: 6315},
				run: (*parser).callonEach1,
				expr: &seqExpr{
					pos: position{line: 250, col: 9, offset: 6315},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 250, col: 9, offset: 6315},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 250, col: 11, offset: 6317},
							val:        "each",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 250, col: 18, offset: 6324},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 250, col: 21, offset: 6327},
							label: "v1",
							expr: &ruleRefExpr{
								pos:  position{line: 250, col: 24, offset: 6330},
								name: "Variable",
							},
						},
						&labeledExpr{
							pos:   position{line: 250, col: 33, offset: 6339},
							label: "v2",
							expr: &zeroOrOneExpr{
								pos: position{line: 250, col: 36, offset: 6342},
								expr: &seqExpr{
									pos: position{line: 250, col: 37, offset: 6343},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 250, col: 37, offset: 6343},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 250, col: 39, offset: 6345},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 250, col: 43, offset: 6349},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 250, col: 45, offset: 6351},
											name: "Variable",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 250, col: 56, offset: 6362},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 250, col: 58, offset: 6364},
							val:        "in",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 250, col: 63, offset: 6369},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 250, col: 65, offset: 6371},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 250, col: 70, offset: 6376},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 250, col: 81, offset: 6387},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 250, col: 83, offset: 6389},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 250, col: 86, offset: 6392},
							label: "block",
							expr: &ruleRefExpr{
								pos:  position{line: 250, col: 92, offset: 6398},
								name: "IndentedList",
							},
						},
					},
				},
			},
		},
		{
			name: "Assignment",
			pos:  position{line: 261, col: 1, offset: 6683},
			expr: &actionExpr{
				pos: position{line: 261, col: 15, offset: 6697},
				run: (*parser).callonAssignment1,
				expr: &seqExpr{
					pos: position{line: 261, col: 15, offset: 6697},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 261, col: 15, offset: 6697},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 261, col: 17, offset: 6699},
							val:        "-",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 261, col: 21, offset: 6703},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 261, col: 24, offset: 6706},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 261, col: 24, offset: 6706},
									val:        "var",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 261, col: 32, offset: 6714},
									val:        "let",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 261, col: 40, offset: 6722},
									val:        "const",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 261, col: 49, offset: 6731},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 261, col: 52, offset: 6734},
							label: "vr",
							expr: &ruleRefExpr{
								pos:  position{line: 261, col: 55, offset: 6737},
								name: "Variable",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 261, col: 64, offset: 6746},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 261, col: 66, offset: 6748},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 261, col: 70, offset: 6752},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 261, col: 72, offset: 6754},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 261, col: 77, offset: 6759},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 261, col: 88, offset: 6770},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 261, col: 90, offset: 6772},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Mixin",
			pos:  position{line: 266, col: 1, offset: 6904},
			expr: &actionExpr{
				pos: position{line: 266, col: 10, offset: 6913},
				run: (*parser).callonMixin1,
				expr: &seqExpr{
					pos: position{line: 266, col: 10, offset: 6913},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 266, col: 10, offset: 6913},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 266, col: 12, offset: 6915},
							val:        "mixin",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 266, col: 20, offset: 6923},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 266, col: 23, offset: 6926},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 266, col: 28, offset: 6931},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 266, col: 39, offset: 6942},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 266, col: 41, offset: 6944},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 266, col: 46, offset: 6949},
								expr: &ruleRefExpr{
									pos:  position{line: 266, col: 46, offset: 6949},
									name: "MixinArguments",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 266, col: 62, offset: 6965},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 266, col: 65, offset: 6968},
							label: "list",
							expr: &ruleRefExpr{
								pos:  position{line: 266, col: 70, offset: 6973},
								name: "IndentedList",
							},
						},
					},
				},
			},
		},
		{
			name: "MixinArguments",
			pos:  position{line: 274, col: 1, offset: 7182},
			expr: &choiceExpr{
				pos: position{line: 274, col: 19, offset: 7200},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 274, col: 19, offset: 7200},
						run: (*parser).callonMixinArguments2,
						expr: &seqExpr{
							pos: position{line: 274, col: 19, offset: 7200},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 274, col: 19, offset: 7200},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 274, col: 23, offset: 7204},
									label: "head",
									expr: &ruleRefExpr{
										pos:  position{line: 274, col: 28, offset: 7209},
										name: "MixinArgument",
									},
								},
								&labeledExpr{
									pos:   position{line: 274, col: 42, offset: 7223},
									label: "tail",
									expr: &zeroOrMoreExpr{
										pos: position{line: 274, col: 47, offset: 7228},
										expr: &seqExpr{
											pos: position{line: 274, col: 48, offset: 7229},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 274, col: 48, offset: 7229},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 274, col: 50, offset: 7231},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 274, col: 54, offset: 7235},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 274, col: 56, offset: 7237},
													name: "MixinArgument",
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 274, col: 72, offset: 7253},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 287, col: 5, offset: 7515},
						run: (*parser).callonMixinArguments15,
						expr: &seqExpr{
							pos: position{line: 287, col: 5, offset: 7515},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 287, col: 5, offset: 7515},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 287, col: 9, offset: 7519},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 287, col: 11, offset: 7521},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "MixinArgument",
			pos:  position{line: 291, col: 1, offset: 7548},
			expr: &actionExpr{
				pos: position{line: 291, col: 18, offset: 7565},
				run: (*parser).callonMixinArgument1,
				expr: &seqExpr{
					pos: position{line: 291, col: 18, offset: 7565},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 291, col: 18, offset: 7565},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 291, col: 23, offset: 7570},
								name: "Variable",
							},
						},
						&labeledExpr{
							pos:   position{line: 291, col: 32, offset: 7579},
							label: "def",
							expr: &zeroOrOneExpr{
								pos: position{line: 291, col: 36, offset: 7583},
								expr: &seqExpr{
									pos: position{line: 291, col: 37, offset: 7584},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 291, col: 37, offset: 7584},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 291, col: 39, offset: 7586},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 291, col: 43, offset: 7590},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 291, col: 45, offset: 7592},
											name: "Expression",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "MixinCall",
			pos:  position{line: 302, col: 1, offset: 7816},
			expr: &actionExpr{
				pos: position{line: 302, col: 14, offset: 7829},
				run: (*parser).callonMixinCall1,
				expr: &seqExpr{
					pos: position{line: 302, col: 14, offset: 7829},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 302, col: 14, offset: 7829},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 302, col: 16, offset: 7831},
							val:        "+",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 302, col: 20, offset: 7835},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 302, col: 25, offset: 7840},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 302, col: 36, offset: 7851},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 302, col: 38, offset: 7853},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 302, col: 43, offset: 7858},
								expr: &ruleRefExpr{
									pos:  position{line: 302, col: 43, offset: 7858},
									name: "CallArguments",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 302, col: 58, offset: 7873},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "CallArguments",
			pos:  position{line: 310, col: 1, offset: 8044},
			expr: &actionExpr{
				pos: position{line: 310, col: 18, offset: 8061},
				run: (*parser).callonCallArguments1,
				expr: &seqExpr{
					pos: position{line: 310, col: 18, offset: 8061},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 310, col: 18, offset: 8061},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 310, col: 22, offset: 8065},
							label: "head",
							expr: &zeroOrOneExpr{
								pos: position{line: 310, col: 27, offset: 8070},
								expr: &ruleRefExpr{
									pos:  position{line: 310, col: 27, offset: 8070},
									name: "Expression",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 310, col: 39, offset: 8082},
							label: "tail",
							expr: &zeroOrMoreExpr{
								pos: position{line: 310, col: 44, offset: 8087},
								expr: &seqExpr{
									pos: position{line: 310, col: 45, offset: 8088},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 310, col: 45, offset: 8088},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 310, col: 47, offset: 8090},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 310, col: 51, offset: 8094},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 310, col: 53, offset: 8096},
											name: "Expression",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 310, col: 66, offset: 8109},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Import",
			pos:  position{line: 331, col: 1, offset: 8431},
			expr: &actionExpr{
				pos: position{line: 331, col: 11, offset: 8441},
				run: (*parser).callonImport1,
				expr: &seqExpr{
					pos: position{line: 331, col: 11, offset: 8441},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 331, col: 11, offset: 8441},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 331, col: 13, offset: 8443},
							val:        "include",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 331, col: 23, offset: 8453},
							expr: &litMatcher{
								pos:        position{line: 331, col: 23, offset: 8453},
								val:        "s",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 28, offset: 8458},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 331, col: 31, offset: 8461},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 331, col: 36, offset: 8466},
								name: "LineText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 331, col: 45, offset: 8475},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Extend",
			pos:  position{line: 335, col: 1, offset: 8558},
			expr: &actionExpr{
				pos: position{line: 335, col: 11, offset: 8568},
				run: (*parser).callonExtend1,
				expr: &seqExpr{
					pos: position{line: 335, col: 11, offset: 8568},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 335, col: 11, offset: 8568},
							val:        "extend",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 335, col: 20, offset: 8577},
							expr: &litMatcher{
								pos:        position{line: 335, col: 20, offset: 8577},
								val:        "s",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 335, col: 25, offset: 8582},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 335, col: 28, offset: 8585},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 335, col: 33, offset: 8590},
								name: "LineText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 335, col: 42, offset: 8599},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Block",
			pos:  position{line: 339, col: 1, offset: 8682},
			expr: &actionExpr{
				pos: position{line: 339, col: 10, offset: 8691},
				run: (*parser).callonBlock1,
				expr: &seqExpr{
					pos: position{line: 339, col: 10, offset: 8691},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 339, col: 10, offset: 8691},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 339, col: 12, offset: 8693},
							val:        "block",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 339, col: 20, offset: 8701},
							label: "mod",
							expr: &zeroOrOneExpr{
								pos: position{line: 339, col: 24, offset: 8705},
								expr: &seqExpr{
									pos: position{line: 339, col: 25, offset: 8706},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 339, col: 25, offset: 8706},
											name: "__",
										},
										&choiceExpr{
											pos: position{line: 339, col: 29, offset: 8710},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 339, col: 29, offset: 8710},
													val:        "append",
													ignoreCase: false,
												},
												&litMatcher{
													pos:        position{line: 339, col: 40, offset: 8721},
													val:        "prepend",
													ignoreCase: false,
												},
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 339, col: 53, offset: 8734},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 339, col: 56, offset: 8737},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 339, col: 61, offset: 8742},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 339, col: 72, offset: 8753},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 339, col: 75, offset: 8756},
							label: "list",
							expr: &zeroOrOneExpr{
								pos: position{line: 339, col: 80, offset: 8761},
								expr: &ruleRefExpr{
									pos:  position{line: 339, col: 80, offset: 8761},
									name: "IndentedList",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Comment",
			pos:  position{line: 360, col: 1, offset: 9126},
			expr: &actionExpr{
				pos: position{line: 360, col: 12, offset: 9137},
				run: (*parser).callonComment1,
				expr: &seqExpr{
					pos: position{line: 360, col: 12, offset: 9137},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 360, col: 12, offset: 9137},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 360, col: 14, offset: 9139},
							val:        "//",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 360, col: 19, offset: 9144},
							label: "silent",
							expr: &zeroOrOneExpr{
								pos: position{line: 360, col: 26, offset: 9151},
								expr: &litMatcher{
									pos:        position{line: 360, col: 26, offset: 9151},
									val:        "-",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 360, col: 31, offset: 9156},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 360, col: 33, offset: 9158},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 360, col: 41, offset: 9166},
								name: "LineText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 360, col: 50, offset: 9175},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "LineText",
			pos:  position{line: 365, col: 1, offset: 9309},
			expr: &actionExpr{
				pos: position{line: 365, col: 13, offset: 9321},
				run: (*parser).callonLineText1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 365, col: 13, offset: 9321},
					expr: &charClassMatcher{
						pos:        position{line: 365, col: 13, offset: 9321},
						val:        "[^\\n]",
						chars:      []rune{'\n'},
						ignoreCase: false,
						inverted:   true,
					},
				},
			},
		},
		{
			name: "PipeText",
			pos:  position{line: 370, col: 1, offset: 9370},
			expr: &actionExpr{
				pos: position{line: 370, col: 13, offset: 9382},
				run: (*parser).callonPipeText1,
				expr: &seqExpr{
					pos: position{line: 370, col: 13, offset: 9382},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 370, col: 13, offset: 9382},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 370, col: 16, offset: 9385},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 370, col: 16, offset: 9385},
									val:        "|",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 370, col: 22, offset: 9391},
									val:        "<",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 370, col: 27, offset: 9396},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 370, col: 29, offset: 9398},
							label: "tl",
							expr: &ruleRefExpr{
								pos:  position{line: 370, col: 32, offset: 9401},
								name: "TextList",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 370, col: 41, offset: 9410},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "PipeExpression",
			pos:  position{line: 374, col: 1, offset: 9435},
			expr: &actionExpr{
				pos: position{line: 374, col: 19, offset: 9453},
				run: (*parser).callonPipeExpression1,
				expr: &seqExpr{
					pos: position{line: 374, col: 19, offset: 9453},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 374, col: 19, offset: 9453},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 374, col: 21, offset: 9455},
							label: "mod",
							expr: &choiceExpr{
								pos: position{line: 374, col: 26, offset: 9460},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 374, col: 26, offset: 9460},
										val:        "=",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 374, col: 32, offset: 9466},
										val:        "!=",
										ignoreCase: false,
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 374, col: 38, offset: 9472},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 374, col: 40, offset: 9474},
							label: "ex",
							expr: &ruleRefExpr{
								pos:  position{line: 374, col: 43, offset: 9477},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 374, col: 54, offset: 9488},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "TextList",
			pos:  position{line: 384, col: 1, offset: 9663},
			expr: &choiceExpr{
				pos: position{line: 384, col: 13, offset: 9675},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 384, col: 13, offset: 9675},
						run: (*parser).callonTextList2,
						expr: &seqExpr{
							pos: position{line: 384, col: 13, offset: 9675},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 384, col: 13, offset: 9675},
									label: "intr",
									expr: &ruleRefExpr{
										pos:  position{line: 384, col: 18, offset: 9680},
										name: "Interpolation",
									},
								},
								&labeledExpr{
									pos:   position{line: 384, col: 32, offset: 9694},
									label: "tl",
									expr: &ruleRefExpr{
										pos:  position{line: 384, col: 35, offset: 9697},
										name: "TextList",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 399, col: 5, offset: 10018},
						run: (*parser).callonTextList8,
						expr: &andExpr{
							pos: position{line: 399, col: 5, offset: 10018},
							expr: &ruleRefExpr{
								pos:  position{line: 399, col: 6, offset: 10019},
								name: "NL",
							},
						},
					},
					&actionExpr{
						pos: position{line: 401, col: 5, offset: 10084},
						run: (*parser).callonTextList11,
						expr: &seqExpr{
							pos: position{line: 401, col: 5, offset: 10084},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 401, col: 5, offset: 10084},
									label: "ch",
									expr: &anyMatcher{
										line: 401, col: 8, offset: 10087,
									},
								},
								&labeledExpr{
									pos:   position{line: 401, col: 10, offset: 10089},
									label: "tl",
									expr: &ruleRefExpr{
										pos:  position{line: 401, col: 13, offset: 10092},
										name: "TextList",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Indent",
			pos:  position{line: 418, col: 1, offset: 10488},
			expr: &litMatcher{
				pos:        position{line: 418, col: 11, offset: 10498},
				val:        "\x01",
				ignoreCase: false,
			},
		},
		{
			name: "Outdent",
			pos:  position{line: 419, col: 1, offset: 10507},
			expr: &litMatcher{
				pos:        position{line: 419, col: 12, offset: 10518},
				val:        "\x02",
				ignoreCase: false,
			},
		},
		{
			name: "Interpolation",
			pos:  position{line: 421, col: 1, offset: 10528},
			expr: &actionExpr{
				pos: position{line: 421, col: 18, offset: 10545},
				run: (*parser).callonInterpolation1,
				expr: &seqExpr{
					pos: position{line: 421, col: 18, offset: 10545},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 421, col: 18, offset: 10545},
							label: "mod",
							expr: &choiceExpr{
								pos: position{line: 421, col: 23, offset: 10550},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 421, col: 23, offset: 10550},
										val:        "#",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 421, col: 29, offset: 10556},
										val:        "!",
										ignoreCase: false,
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 421, col: 34, offset: 10561},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 421, col: 38, offset: 10565},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 421, col: 40, offset: 10567},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 421, col: 45, offset: 10572},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 421, col: 56, offset: 10583},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 421, col: 58, offset: 10585},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 429, col: 1, offset: 10769},
			expr: &ruleRefExpr{
				pos:  position{line: 429, col: 15, offset: 10783},
				name: "ExpressionTernery",
			},
		},
		{
			name: "ExpressionTernery",
			pos:  position{line: 431, col: 1, offset: 10802},
			expr: &actionExpr{
				pos: position{line: 431, col: 22, offset: 10823},
				run: (*parser).callonExpressionTernery1,
				expr: &seqExpr{
					pos: position{line: 431, col: 22, offset: 10823},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 431, col: 22, offset: 10823},
							label: "cnd",
							expr: &ruleRefExpr{
								pos:  position{line: 431, col: 26, offset: 10827},
								name: "ExpressionBinOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 431, col: 42, offset: 10843},
							label: "rest",
							expr: &zeroOrOneExpr{
								pos: position{line: 431, col: 47, offset: 10848},
								expr: &seqExpr{
									pos: position{line: 431, col: 48, offset: 10849},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 431, col: 48, offset: 10849},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 431, col: 50, offset: 10851},
											val:        "?",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 431, col: 54, offset: 10855},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 431, col: 56, offset: 10857},
											name: "ExpressionTernery",
										},
										&ruleRefExpr{
											pos:  position{line: 431, col: 74, offset: 10875},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 431, col: 76, offset: 10877},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 431, col: 80, offset: 10881},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 431, col: 82, offset: 10883},
											name: "ExpressionTernery",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ExpressionBinOp",
			pos:  position{line: 451, col: 1, offset: 11255},
			expr: &actionExpr{
				pos: position{line: 451, col: 20, offset: 11274},
				run: (*parser).callonExpressionBinOp1,
				expr: &seqExpr{
					pos: position{line: 451, col: 20, offset: 11274},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 451, col: 20, offset: 11274},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 451, col: 26, offset: 11280},
								name: "ExpressionCmpOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 451, col: 42, offset: 11296},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 451, col: 47, offset: 11301},
								expr: &seqExpr{
									pos: position{line: 451, col: 49, offset: 11303},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 451, col: 49, offset: 11303},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 451, col: 51, offset: 11305},
											name: "CmpOp",
										},
										&ruleRefExpr{
											pos:  position{line: 451, col: 57, offset: 11311},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 451, col: 59, offset: 11313},
											name: "ExpressionBinOp",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ExpressionCmpOp",
			pos:  position{line: 455, col: 1, offset: 11373},
			expr: &actionExpr{
				pos: position{line: 455, col: 20, offset: 11392},
				run: (*parser).callonExpressionCmpOp1,
				expr: &seqExpr{
					pos: position{line: 455, col: 20, offset: 11392},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 455, col: 20, offset: 11392},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 455, col: 26, offset: 11398},
								name: "ExpressionAddOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 455, col: 42, offset: 11414},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 455, col: 47, offset: 11419},
								expr: &seqExpr{
									pos: position{line: 455, col: 49, offset: 11421},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 455, col: 49, offset: 11421},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 455, col: 51, offset: 11423},
											name: "CmpOp",
										},
										&ruleRefExpr{
											pos:  position{line: 455, col: 57, offset: 11429},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 455, col: 59, offset: 11431},
											name: "ExpressionCmpOp",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ExpressionAddOp",
			pos:  position{line: 459, col: 1, offset: 11491},
			expr: &actionExpr{
				pos: position{line: 459, col: 20, offset: 11510},
				run: (*parser).callonExpressionAddOp1,
				expr: &seqExpr{
					pos: position{line: 459, col: 20, offset: 11510},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 459, col: 20, offset: 11510},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 459, col: 26, offset: 11516},
								name: "ExpressionMulOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 459, col: 42, offset: 11532},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 459, col: 47, offset: 11537},
								expr: &seqExpr{
									pos: position{line: 459, col: 49, offset: 11539},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 459, col: 49, offset: 11539},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 459, col: 51, offset: 11541},
											name: "AddOp",
										},
										&ruleRefExpr{
											pos:  position{line: 459, col: 57, offset: 11547},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 459, col: 59, offset: 11549},
											name: "ExpressionAddOp",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ExpressionMulOp",
			pos:  position{line: 463, col: 1, offset: 11609},
			expr: &actionExpr{
				pos: position{line: 463, col: 20, offset: 11628},
				run: (*parser).callonExpressionMulOp1,
				expr: &seqExpr{
					pos: position{line: 463, col: 20, offset: 11628},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 463, col: 20, offset: 11628},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 463, col: 26, offset: 11634},
								name: "ExpressionUnaryOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 463, col: 44, offset: 11652},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 463, col: 49, offset: 11657},
								expr: &seqExpr{
									pos: position{line: 463, col: 51, offset: 11659},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 463, col: 51, offset: 11659},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 463, col: 53, offset: 11661},
											name: "MulOp",
										},
										&ruleRefExpr{
											pos:  position{line: 463, col: 59, offset: 11667},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 463, col: 61, offset: 11669},
											name: "ExpressionMulOp",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ExpressionUnaryOp",
			pos:  position{line: 467, col: 1, offset: 11729},
			expr: &choiceExpr{
				pos: position{line: 467, col: 22, offset: 11750},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 467, col: 22, offset: 11750},
						run: (*parser).callonExpressionUnaryOp2,
						expr: &seqExpr{
							pos: position{line: 467, col: 22, offset: 11750},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 467, col: 22, offset: 11750},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 467, col: 25, offset: 11753},
										name: "UnaryOp",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 467, col: 33, offset: 11761},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 467, col: 35, offset: 11763},
									label: "ex",
									expr: &ruleRefExpr{
										pos:  position{line: 467, col: 38, offset: 11766},
										name: "ExpressionFactor",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 469, col: 5, offset: 11889},
						name: "ExpressionFactor",
					},
				},
			},
		},
		{
			name: "ExpressionFactor",
			pos:  position{line: 471, col: 1, offset: 11907},
			expr: &choiceExpr{
				pos: position{line: 471, col: 21, offset: 11927},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 471, col: 21, offset: 11927},
						run: (*parser).callonExpressionFactor2,
						expr: &seqExpr{
							pos: position{line: 471, col: 21, offset: 11927},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 471, col: 21, offset: 11927},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 471, col: 25, offset: 11931},
									label: "e",
									expr: &ruleRefExpr{
										pos:  position{line: 471, col: 27, offset: 11933},
										name: "Expression",
									},
								},
								&litMatcher{
									pos:        position{line: 471, col: 38, offset: 11944},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 473, col: 5, offset: 11970},
						name: "StringExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 473, col: 24, offset: 11989},
						name: "NumberExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 473, col: 43, offset: 12008},
						name: "BooleanExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 473, col: 63, offset: 12028},
						name: "NilExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 473, col: 79, offset: 12044},
						name: "MemberExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 473, col: 98, offset: 12063},
						name: "ArrayExpression",
					},
				},
			},
		},
		{
			name: "StringExpression",
			pos:  position{line: 475, col: 1, offset: 12080},
			expr: &actionExpr{
				pos: position{line: 475, col: 21, offset: 12100},
				run: (*parser).callonStringExpression1,
				expr: &labeledExpr{
					pos:   position{line: 475, col: 21, offset: 12100},
					label: "s",
					expr: &ruleRefExpr{
						pos:  position{line: 475, col: 23, offset: 12102},
						name: "String",
					},
				},
			},
		},
		{
			name: "NumberExpression",
			pos:  position{line: 479, col: 1, offset: 12197},
			expr: &actionExpr{
				pos: position{line: 479, col: 21, offset: 12217},
				run: (*parser).callonNumberExpression1,
				expr: &seqExpr{
					pos: position{line: 479, col: 21, offset: 12217},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 479, col: 21, offset: 12217},
							expr: &litMatcher{
								pos:        position{line: 479, col: 21, offset: 12217},
								val:        "-",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 479, col: 26, offset: 12222},
							name: "Integer",
						},
						&labeledExpr{
							pos:   position{line: 479, col: 34, offset: 12230},
							label: "dec",
							expr: &zeroOrOneExpr{
								pos: position{line: 479, col: 38, offset: 12234},
								expr: &seqExpr{
									pos: position{line: 479, col: 40, offset: 12236},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 479, col: 40, offset: 12236},
											val:        ".",
											ignoreCase: false,
										},
										&oneOrMoreExpr{
											pos: position{line: 479, col: 44, offset: 12240},
											expr: &ruleRefExpr{
												pos:  position{line: 479, col: 44, offset: 12240},
												name: "DecimalDigit",
											},
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 479, col: 61, offset: 12257},
							label: "ex",
							expr: &zeroOrOneExpr{
								pos: position{line: 479, col: 64, offset: 12260},
								expr: &ruleRefExpr{
									pos:  position{line: 479, col: 64, offset: 12260},
									name: "Exponent",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "NilExpression",
			pos:  position{line: 489, col: 1, offset: 12590},
			expr: &actionExpr{
				pos: position{line: 489, col: 18, offset: 12607},
				run: (*parser).callonNilExpression1,
				expr: &ruleRefExpr{
					pos:  position{line: 489, col: 18, offset: 12607},
					name: "Null",
				},
			},
		},
		{
			name: "BooleanExpression",
			pos:  position{line: 493, col: 1, offset: 12678},
			expr: &actionExpr{
				pos: position{line: 493, col: 22, offset: 12699},
				run: (*parser).callonBooleanExpression1,
				expr: &labeledExpr{
					pos:   position{line: 493, col: 22, offset: 12699},
					label: "b",
					expr: &ruleRefExpr{
						pos:  position{line: 493, col: 24, offset: 12701},
						name: "Bool",
					},
				},
			},
		},
		{
			name: "MemberExpression",
			pos:  position{line: 497, col: 1, offset: 12793},
			expr: &actionExpr{
				pos: position{line: 497, col: 21, offset: 12813},
				run: (*parser).callonMemberExpression1,
				expr: &seqExpr{
					pos: position{line: 497, col: 21, offset: 12813},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 497, col: 21, offset: 12813},
							label: "field",
							expr: &choiceExpr{
								pos: position{line: 497, col: 28, offset: 12820},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 497, col: 28, offset: 12820},
										name: "Field",
									},
									&ruleRefExpr{
										pos:  position{line: 497, col: 36, offset: 12828},
										name: "ArrayExpression",
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 497, col: 53, offset: 12845},
							label: "member",
							expr: &zeroOrMoreExpr{
								pos: position{line: 497, col: 60, offset: 12852},
								expr: &choiceExpr{
									pos: position{line: 497, col: 61, offset: 12853},
									alternatives: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 497, col: 61, offset: 12853},
											name: "MemberField",
										},
										&ruleRefExpr{
											pos:  position{line: 497, col: 75, offset: 12867},
											name: "MemberIndex",
										},
										&ruleRefExpr{
											pos:  position{line: 497, col: 89, offset: 12881},
											name: "MemberCall",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "MemberField",
			pos:  position{line: 515, col: 1, offset: 13380},
			expr: &actionExpr{
				pos: position{line: 515, col: 16, offset: 13395},
				run: (*parser).callonMemberField1,
				expr: &seqExpr{
					pos: position{line: 515, col: 16, offset: 13395},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 515, col: 16, offset: 13395},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 515, col: 20, offset: 13399},
							label: "ident",
							expr: &ruleRefExpr{
								pos:  position{line: 515, col: 26, offset: 13405},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "MemberIndex",
			pos:  position{line: 519, col: 1, offset: 13441},
			expr: &actionExpr{
				pos: position{line: 519, col: 16, offset: 13456},
				run: (*parser).callonMemberIndex1,
				expr: &seqExpr{
					pos: position{line: 519, col: 16, offset: 13456},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 519, col: 16, offset: 13456},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 519, col: 18, offset: 13458},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 519, col: 22, offset: 13462},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 519, col: 24, offset: 13464},
							label: "i",
							expr: &ruleRefExpr{
								pos:  position{line: 519, col: 26, offset: 13466},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 519, col: 37, offset: 13477},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 519, col: 39, offset: 13479},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "MemberCall",
			pos:  position{line: 523, col: 1, offset: 13504},
			expr: &actionExpr{
				pos: position{line: 523, col: 15, offset: 13518},
				run: (*parser).callonMemberCall1,
				expr: &seqExpr{
					pos: position{line: 523, col: 15, offset: 13518},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 523, col: 15, offset: 13518},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 523, col: 17, offset: 13520},
							label: "arg",
							expr: &ruleRefExpr{
								pos:  position{line: 523, col: 21, offset: 13524},
								name: "CallArguments",
							},
						},
					},
				},
			},
		},
		{
			name: "ArrayExpression",
			pos:  position{line: 527, col: 1, offset: 13561},
			expr: &actionExpr{
				pos: position{line: 527, col: 20, offset: 13580},
				run: (*parser).callonArrayExpression1,
				expr: &seqExpr{
					pos: position{line: 527, col: 20, offset: 13580},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 527, col: 20, offset: 13580},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 527, col: 24, offset: 13584},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 527, col: 26, offset: 13586},
							label: "head",
							expr: &zeroOrOneExpr{
								pos: position{line: 527, col: 31, offset: 13591},
								expr: &ruleRefExpr{
									pos:  position{line: 527, col: 31, offset: 13591},
									name: "Expression",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 527, col: 43, offset: 13603},
							label: "tail",
							expr: &zeroOrMoreExpr{
								pos: position{line: 527, col: 48, offset: 13608},
								expr: &seqExpr{
									pos: position{line: 527, col: 49, offset: 13609},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 527, col: 49, offset: 13609},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 527, col: 51, offset: 13611},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 527, col: 55, offset: 13615},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 527, col: 57, offset: 13617},
											name: "Expression",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 527, col: 70, offset: 13630},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 527, col: 72, offset: 13632},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 540, col: 1, offset: 13998},
			expr: &actionExpr{
				pos: position{line: 540, col: 10, offset: 14007},
				run: (*parser).callonField1,
				expr: &labeledExpr{
					pos:   position{line: 540, col: 10, offset: 14007},
					label: "variable",
					expr: &ruleRefExpr{
						pos:  position{line: 540, col: 19, offset: 14016},
						name: "Variable",
					},
				},
			},
		},
		{
			name: "UnaryOp",
			pos:  position{line: 544, col: 1, offset: 14125},
			expr: &actionExpr{
				pos: position{line: 544, col: 12, offset: 14136},
				run: (*parser).callonUnaryOp1,
				expr: &choiceExpr{
					pos: position{line: 544, col: 14, offset: 14138},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 544, col: 14, offset: 14138},
							val:        "+",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 544, col: 20, offset: 14144},
							val:        "-",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 544, col: 26, offset: 14150},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AddOp",
			pos:  position{line: 548, col: 1, offset: 14190},
			expr: &actionExpr{
				pos: position{line: 548, col: 10, offset: 14199},
				run: (*parser).callonAddOp1,
				expr: &choiceExpr{
					pos: position{line: 548, col: 12, offset: 14201},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 548, col: 12, offset: 14201},
							val:        "+",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 548, col: 18, offset: 14207},
							val:        "-",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "MulOp",
			pos:  position{line: 552, col: 1, offset: 14247},
			expr: &actionExpr{
				pos: position{line: 552, col: 10, offset: 14256},
				run: (*parser).callonMulOp1,
				expr: &choiceExpr{
					pos: position{line: 552, col: 12, offset: 14258},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 552, col: 12, offset: 14258},
							val:        "*",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 552, col: 18, offset: 14264},
							val:        "/",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 552, col: 24, offset: 14270},
							val:        "%",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CmpOp",
			pos:  position{line: 556, col: 1, offset: 14310},
			expr: &actionExpr{
				pos: position{line: 556, col: 10, offset: 14319},
				run: (*parser).callonCmpOp1,
				expr: &choiceExpr{
					pos: position{line: 556, col: 12, offset: 14321},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 556, col: 12, offset: 14321},
							val:        "==",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 556, col: 19, offset: 14328},
							val:        "!=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 556, col: 26, offset: 14335},
							val:        "<",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 556, col: 32, offset: 14341},
							val:        ">",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 556, col: 38, offset: 14347},
							val:        "<=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 556, col: 45, offset: 14354},
							val:        ">=",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "BinOp",
			pos:  position{line: 560, col: 1, offset: 14395},
			expr: &actionExpr{
				pos: position{line: 560, col: 10, offset: 14404},
				run: (*parser).callonBinOp1,
				expr: &choiceExpr{
					pos: position{line: 560, col: 12, offset: 14406},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 560, col: 12, offset: 14406},
							val:        "&&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 560, col: 19, offset: 14413},
							val:        "||",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name:        "String",
			displayName: "\"string\"",
			pos:         position{line: 564, col: 1, offset: 14454},
			expr: &actionExpr{
				pos: position{line: 564, col: 20, offset: 14473},
				run: (*parser).callonString1,
				expr: &seqExpr{
					pos: position{line: 564, col: 20, offset: 14473},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 564, col: 20, offset: 14473},
							name: "Quote",
						},
						&zeroOrMoreExpr{
							pos: position{line: 564, col: 26, offset: 14479},
							expr: &choiceExpr{
								pos: position{line: 564, col: 28, offset: 14481},
								alternatives: []interface{}{
									&seqExpr{
										pos: position{line: 564, col: 28, offset: 14481},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 564, col: 28, offset: 14481},
												expr: &ruleRefExpr{
													pos:  position{line: 564, col: 29, offset: 14482},
													name: "EscapedChar",
												},
											},
											&anyMatcher{
												line: 564, col: 41, offset: 14494,
											},
										},
									},
									&seqExpr{
										pos: position{line: 564, col: 45, offset: 14498},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 564, col: 45, offset: 14498},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 564, col: 50, offset: 14503},
												name: "EscapeSequence",
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 564, col: 68, offset: 14521},
							name: "Quote",
						},
					},
				},
			},
		},
		{
			name: "Index",
			pos:  position{line: 568, col: 1, offset: 14573},
			expr: &actionExpr{
				pos: position{line: 568, col: 10, offset: 14582},
				run: (*parser).callonIndex1,
				expr: &ruleRefExpr{
					pos:  position{line: 568, col: 10, offset: 14582},
					name: "Integer",
				},
			},
		},
		{
			name:        "Quote",
			displayName: "\"quote\"",
			pos:         position{line: 572, col: 1, offset: 14645},
			expr: &litMatcher{
				pos:        position{line: 572, col: 18, offset: 14662},
				val:        "\"",
				ignoreCase: false,
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 574, col: 1, offset: 14667},
			expr: &charClassMatcher{
				pos:        position{line: 574, col: 16, offset: 14682},
				val:        "[\\x00-\\x1f\"\\\\]",
				chars:      []rune{'"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 575, col: 1, offset: 14697},
			expr: &choiceExpr{
				pos: position{line: 575, col: 19, offset: 14715},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 575, col: 19, offset: 14715},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 575, col: 38, offset: 14734},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 576, col: 1, offset: 14748},
			expr: &charClassMatcher{
				pos:        position{line: 576, col: 21, offset: 14768},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 577, col: 1, offset: 14780},
			expr: &seqExpr{
				pos: position{line: 577, col: 18, offset: 14797},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 577, col: 18, offset: 14797},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 577, col: 22, offset: 14801},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 577, col: 31, offset: 14810},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 577, col: 40, offset: 14819},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 577, col: 49, offset: 14828},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 579, col: 1, offset: 14838},
			expr: &choiceExpr{
				pos: position{line: 579, col: 12, offset: 14849},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 579, col: 12, offset: 14849},
						val:        "0",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 579, col: 18, offset: 14855},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 579, col: 18, offset: 14855},
								name: "NonZeroDecimalDigit",
							},
							&zeroOrMoreExpr{
								pos: position{line: 579, col: 38, offset: 14875},
								expr: &ruleRefExpr{
									pos:  position{line: 579, col: 38, offset: 14875},
									name: "DecimalDigit",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Exponent",
			pos:  position{line: 580, col: 1, offset: 14889},
			expr: &seqExpr{
				pos: position{line: 580, col: 13, offset: 14901},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 580, col: 13, offset: 14901},
						val:        "e",
						ignoreCase: true,
					},
					&zeroOrOneExpr{
						pos: position{line: 580, col: 18, offset: 14906},
						expr: &charClassMatcher{
							pos:        position{line: 580, col: 18, offset: 14906},
							val:        "[+-]",
							chars:      []rune{'+', '-'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&oneOrMoreExpr{
						pos: position{line: 580, col: 24, offset: 14912},
						expr: &ruleRefExpr{
							pos:  position{line: 580, col: 24, offset: 14912},
							name: "DecimalDigit",
						},
					},
				},
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 581, col: 1, offset: 14926},
			expr: &charClassMatcher{
				pos:        position{line: 581, col: 17, offset: 14942},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "NonZeroDecimalDigit",
			pos:  position{line: 582, col: 1, offset: 14948},
			expr: &charClassMatcher{
				pos:        position{line: 582, col: 24, offset: 14971},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 583, col: 1, offset: 14977},
			expr: &charClassMatcher{
				pos:        position{line: 583, col: 13, offset: 14989},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "Bool",
			pos:  position{line: 584, col: 1, offset: 14999},
			expr: &choiceExpr{
				pos: position{line: 584, col: 9, offset: 15007},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 584, col: 9, offset: 15007},
						run: (*parser).callonBool2,
						expr: &litMatcher{
							pos:        position{line: 584, col: 9, offset: 15007},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 584, col: 39, offset: 15037},
						run: (*parser).callonBool4,
						expr: &litMatcher{
							pos:        position{line: 584, col: 39, offset: 15037},
							val:        "false",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Null",
			pos:  position{line: 585, col: 1, offset: 15067},
			expr: &actionExpr{
				pos: position{line: 585, col: 9, offset: 15075},
				run: (*parser).callonNull1,
				expr: &choiceExpr{
					pos: position{line: 585, col: 10, offset: 15076},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 585, col: 10, offset: 15076},
							val:        "null",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 585, col: 19, offset: 15085},
							val:        "nil",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Variable",
			pos:  position{line: 587, col: 1, offset: 15113},
			expr: &actionExpr{
				pos: position{line: 587, col: 13, offset: 15125},
				run: (*parser).callonVariable1,
				expr: &labeledExpr{
					pos:   position{line: 587, col: 13, offset: 15125},
					label: "ident",
					expr: &ruleRefExpr{
						pos:  position{line: 587, col: 19, offset: 15131},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 591, col: 1, offset: 15225},
			expr: &actionExpr{
				pos: position{line: 591, col: 15, offset: 15239},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 591, col: 15, offset: 15239},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 591, col: 15, offset: 15239},
							val:        "[a-zA-Z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 591, col: 25, offset: 15249},
							expr: &charClassMatcher{
								pos:        position{line: 591, col: 25, offset: 15249},
								val:        "[a-zA-Z0-9_]",
								chars:      []rune{'_'},
								ranges:     []rune{'a', 'z', 'A', 'Z', '0', '9'},
								ignoreCase: false,
								inverted:   false,
							},
						},
					},
				},
			},
		},
		{
			name: "EmptyLine",
			pos:  position{line: 595, col: 1, offset: 15297},
			expr: &actionExpr{
				pos: position{line: 595, col: 14, offset: 15310},
				run: (*parser).callonEmptyLine1,
				expr: &seqExpr{
					pos: position{line: 595, col: 14, offset: 15310},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 595, col: 14, offset: 15310},
							name: "_",
						},
						&charClassMatcher{
							pos:        position{line: 595, col: 16, offset: 15312},
							val:        "[\\n]",
							chars:      []rune{'\n'},
							ignoreCase: false,
							inverted:   false,
						},
					},
				},
			},
		},
		{
			name:        "_",
			displayName: "\"whitespace\"",
			pos:         position{line: 599, col: 1, offset: 15340},
			expr: &actionExpr{
				pos: position{line: 599, col: 19, offset: 15358},
				run: (*parser).callon_1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 599, col: 19, offset: 15358},
					expr: &charClassMatcher{
						pos:        position{line: 599, col: 19, offset: 15358},
						val:        "[ \\t]",
						chars:      []rune{' ', '\t'},
						ignoreCase: false,
						inverted:   false,
					},
				},
			},
		},
		{
			name:        "__",
			displayName: "\"whitespace\"",
			pos:         position{line: 600, col: 1, offset: 15385},
			expr: &actionExpr{
				pos: position{line: 600, col: 20, offset: 15404},
				run: (*parser).callon__1,
				expr: &charClassMatcher{
					pos:        position{line: 600, col: 20, offset: 15404},
					val:        "[ \\t]",
					chars:      []rune{' ', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "NL",
			pos:  position{line: 601, col: 1, offset: 15431},
			expr: &choiceExpr{
				pos: position{line: 601, col: 7, offset: 15437},
				alternatives: []interface{}{
					&charClassMatcher{
						pos:        position{line: 601, col: 7, offset: 15437},
						val:        "[\\n]",
						chars:      []rune{'\n'},
						ignoreCase: false,
						inverted:   false,
					},
					&andExpr{
						pos: position{line: 601, col: 14, offset: 15444},
						expr: &ruleRefExpr{
							pos:  position{line: 601, col: 15, offset: 15445},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 602, col: 1, offset: 15449},
			expr: &notExpr{
				pos: position{line: 602, col: 8, offset: 15456},
				expr: &anyMatcher{
					line: 602, col: 9, offset: 15457,
				},
			},
		},
	},
}

func (c *current) onInput1(l interface{}) (interface{}, error) {
	return &Root{List: l.(*List), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonInput1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInput1(stack["l"])
}

func (c *current) onList2(node, list interface{}) (interface{}, error) {
	listItem := list.(*List)
	if node != nil {
		listItem.Nodes = append([]Node{node.(Node)}, listItem.Nodes...)
	}
	listItem.Position = pos(c.pos)

	return listItem, nil
}

func (p *parser) callonList2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onList2(stack["node"], stack["list"])
}

func (c *current) onList8() (interface{}, error) {
	return &List{GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonList8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onList8()
}

func (c *current) onIndentedList1(list interface{}) (interface{}, error) {
	return list, nil
}

func (p *parser) callonIndentedList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIndentedList1(stack["list"])
}

func (c *current) onIndentedRawText1(t interface{}) (interface{}, error) {
	return &Text{Value: t.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonIndentedRawText1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIndentedRawText1(stack["t"])
}

func (c *current) onRawText2(rt, tail interface{}) (interface{}, error) {
	return rt.(string) + tail.(string), nil
}

func (p *parser) callonRawText2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRawText2(stack["rt"], stack["tail"])
}

func (c *current) onRawText10() (interface{}, error) {
	return "", nil
}

func (p *parser) callonRawText10() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRawText10()
}

func (c *current) onRawText16(head, tail interface{}) (interface{}, error) {
	return string(head.([]byte)) + tail.(string), nil
}

func (p *parser) callonRawText16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onRawText16(stack["head"], stack["tail"])
}

func (c *current) onListNode16() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonListNode16() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListNode16()
}

func (c *current) onDocType1(val interface{}) (interface{}, error) {
	return &DocType{Value: val.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonDocType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDocType1(stack["val"])
}

func (c *current) onTag1(tag, list interface{}) (interface{}, error) {
	tagElem := tag.(*Tag)

	if list != nil {
		if tagElem.Block != nil {
			tagElem.Block = &List{GraphNode: NewNode(pos(c.pos)), Nodes: []Node{tagElem.Block, list.(Node)}}
		} else {
			tagElem.Block = list.(*List)
		}
	}

	return tagElem, nil
}

func (p *parser) callonTag1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTag1(stack["tag"], stack["list"])
}

func (c *current) onTagHeader2(name, attrs, selfClose, tl interface{}) (interface{}, error) {
	tag := &Tag{Name: name.(string), GraphNode: NewNode(pos(c.pos))}
	if attrs != nil {
		tag.Attributes = attrs.([]*Attribute)
	}
	if tl != nil {
		tag.Text = toSlice(tl)[1].(*TextList)
	}
	if selfClose != nil {
		tag.SelfClose = true
	}
	return tag, nil
}

func (p *parser) callonTagHeader2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagHeader2(stack["name"], stack["attrs"], stack["selfClose"], stack["tl"])
}

func (c *current) onTagHeader20(name, attrs, text interface{}) (interface{}, error) {
	tag := &Tag{Name: name.(string), GraphNode: NewNode(pos(c.pos))}
	if attrs != nil {
		tag.Attributes = attrs.([]*Attribute)
	}
	if text != nil {
		tag.Block = text.(*Text)
	}
	return tag, nil
}

func (p *parser) callonTagHeader20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagHeader20(stack["name"], stack["attrs"], stack["text"])
}

func (c *current) onTagHeader33(name, attrs, block interface{}) (interface{}, error) {
	tag := &Tag{Name: name.(string), GraphNode: NewNode(pos(c.pos))}
	if attrs != nil {
		tag.Attributes = attrs.([]*Attribute)
	}
	if block != nil {
		tag.Block = block.(Node)
	}
	return tag, nil
}

func (p *parser) callonTagHeader33() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagHeader33(stack["name"], stack["attrs"], stack["block"])
}

func (c *current) onTagHeader44(name, attrs, unescaped, expr interface{}) (interface{}, error) {
	tag := &Tag{Name: name.(string), GraphNode: NewNode(pos(c.pos))}
	if attrs != nil {
		tag.Attributes = attrs.([]*Attribute)
	}
	if expr != nil {
		intr := &Interpolation{Expr: expr.(Expression), GraphNode: NewNode(pos(c.pos))}
		if unescaped != nil {
			intr.Unescaped = true
		}
		tag.Block = &TextList{Nodes: []Node{intr}, GraphNode: NewNode(pos(c.pos))}
	}
	return tag, nil
}

func (p *parser) callonTagHeader44() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagHeader44(stack["name"], stack["attrs"], stack["unescaped"], stack["expr"])
}

func (c *current) onTagName2() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonTagName2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagName2()
}

func (c *current) onTagName9() (interface{}, error) {
	return "div", nil
}

func (p *parser) callonTagName9() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagName9()
}

func (c *current) onTagAttributes2(head, tail interface{}) (interface{}, error) {
	tailElem := []*Attribute{}

	if tail != nil {
		tailElem = tail.([]*Attribute)
	}

	return append(head.([]*Attribute), tailElem...), nil
}

func (p *parser) callonTagAttributes2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributes2(stack["head"], stack["tail"])
}

func (c *current) onTagAttributes11(head, tail interface{}) (interface{}, error) {
	tailElem := []*Attribute{}

	if tail != nil {
		tailElem = tail.([]*Attribute)
	}

	vals := toSlice(toSlice(head)[2])

	if len(vals) == 0 {
		return tailElem, nil
	}

	headAttrs := vals[0].([]*Attribute)
	restAttrs := toSlice(vals[1])

	for _, a := range restAttrs {
		restAttr := toSlice(a)
		headAttrs = append(headAttrs, restAttr[1].([]*Attribute)...)
	}

	return append(headAttrs, tailElem...), nil
}

func (p *parser) callonTagAttributes11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributes11(stack["head"], stack["tail"])
}

func (c *current) onTagAttributeClass1(name interface{}) (interface{}, error) {
	return []*Attribute{&Attribute{Name: "class", Value: &StringExpression{Value: name.(string)}, GraphNode: NewNode(pos(c.pos))}}, nil
}

func (p *parser) callonTagAttributeClass1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeClass1(stack["name"])
}

func (c *current) onTagAttributeID1(name interface{}) (interface{}, error) {
	return []*Attribute{&Attribute{Name: "id", Value: &StringExpression{Value: name.(string)}, GraphNode: NewNode(pos(c.pos))}}, nil
}

func (p *parser) callonTagAttributeID1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeID1(stack["name"])
}

func (c *current) onTagAttribute2(name, value interface{}) (interface{}, error) {
	return []*Attribute{&Attribute{Name: name.(string), Value: value.(Expression), GraphNode: NewNode(pos(c.pos))}}, nil
}

func (p *parser) callonTagAttribute2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttribute2(stack["name"], stack["value"])
}

func (c *current) onTagAttribute11(name, value interface{}) (interface{}, error) {
	return []*Attribute{&Attribute{Name: name.(string), Value: value.(Expression), Unescaped: true, GraphNode: NewNode(pos(c.pos))}}, nil
}

func (p *parser) callonTagAttribute11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttribute11(stack["name"], stack["value"])
}

func (c *current) onTagAttribute20(name interface{}) (interface{}, error) {
	return []*Attribute{&Attribute{Name: name.(string), GraphNode: NewNode(pos(c.pos))}}, nil
}

func (p *parser) callonTagAttribute20() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttribute20(stack["name"])
}

func (c *current) onTagAttributeName2(tn interface{}) (interface{}, error) {
	return tn, nil
}

func (p *parser) callonTagAttributeName2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeName2(stack["tn"])
}

func (c *current) onTagAttributeName8(tn interface{}) (interface{}, error) {
	return tn, nil
}

func (p *parser) callonTagAttributeName8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeName8(stack["tn"])
}

func (c *current) onTagAttributeName14(tn interface{}) (interface{}, error) {
	return tn, nil
}

func (p *parser) callonTagAttributeName14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeName14(stack["tn"])
}

func (c *current) onClassName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonClassName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onClassName1()
}

func (c *current) onIdName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIdName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdName1()
}

func (c *current) onTagAttributeNameLiteral1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonTagAttributeNameLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttributeNameLiteral1()
}

func (c *current) onIf1(expr, block, elseNode interface{}) (interface{}, error) {
	ifElem := &If{Condition: expr.(Expression), PositiveBlock: block.(Node), GraphNode: NewNode(pos(c.pos))}
	if elseNode != nil {
		ifElem.NegativeBlock = elseNode.(Node)
	}
	return ifElem, nil
}

func (p *parser) callonIf1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIf1(stack["expr"], stack["block"], stack["elseNode"])
}

func (c *current) onUnless1(expr, block interface{}) (interface{}, error) {
	condition := &UnaryExpression{
		Op: "!",
		X:  expr.(Expression),
	}

	return &If{Condition: condition, PositiveBlock: block.(Node), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonUnless1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnless1(stack["expr"], stack["block"])
}

func (c *current) onElse2(node interface{}) (interface{}, error) {
	return node, nil
}

func (p *parser) callonElse2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onElse2(stack["node"])
}

func (c *current) onElse8(block interface{}) (interface{}, error) {
	return block, nil
}

func (p *parser) callonElse8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onElse8(stack["block"])
}

func (c *current) onEach1(v1, v2, expr, block interface{}) (interface{}, error) {
	eachElem := &Each{GraphNode: NewNode(pos(c.pos)), Block: block.(Node), ElementVariable: v1.(*Variable), Container: expr.(Expression)}
	v2Slice := toSlice(v2)

	if len(v2Slice) != 0 {
		eachElem.IndexVariable = v2Slice[3].(*Variable)
	}

	return eachElem, nil
}

func (p *parser) callonEach1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEach1(stack["v1"], stack["v2"], stack["expr"], stack["block"])
}

func (c *current) onAssignment1(vr, expr interface{}) (interface{}, error) {
	return &Assignment{Variable: vr.(*Variable), Expression: expr.(Expression), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonAssignment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAssignment1(stack["vr"], stack["expr"])
}

func (c *current) onMixin1(name, args, list interface{}) (interface{}, error) {
	mixinElem := &Mixin{Name: name.(string), Block: list.(Node), GraphNode: NewNode(pos(c.pos))}
	if args != nil {
		mixinElem.Arguments = args.([]MixinArgument)
	}
	return mixinElem, nil
}

func (p *parser) callonMixin1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMixin1(stack["name"], stack["args"], stack["list"])
}

func (c *current) onMixinArguments2(head, tail interface{}) (interface{}, error) {
	args := []MixinArgument{head.(MixinArgument)}

	if tail != nil {
		tailSlice := toSlice(tail)

		for _, arg := range tailSlice {
			argSlice := toSlice(arg)
			args = append(args, argSlice[3].(MixinArgument))
		}
	}

	return args, nil
}

func (p *parser) callonMixinArguments2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMixinArguments2(stack["head"], stack["tail"])
}

func (c *current) onMixinArguments15() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonMixinArguments15() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMixinArguments15()
}

func (c *current) onMixinArgument1(name, def interface{}) (interface{}, error) {
	argElem := MixinArgument{Name: name.(*Variable), GraphNode: NewNode(pos(c.pos))}

	if def != nil {
		defSlice := toSlice(def)
		argElem.Default = defSlice[3].(Expression)
	}

	return argElem, nil
}

func (p *parser) callonMixinArgument1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMixinArgument1(stack["name"], stack["def"])
}

func (c *current) onMixinCall1(name, args interface{}) (interface{}, error) {
	mcElem := &MixinCall{Name: name.(string), GraphNode: NewNode(pos(c.pos))}
	if args != nil {
		mcElem.Arguments = args.([]Expression)
	}
	return mcElem, nil
}

func (p *parser) callonMixinCall1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMixinCall1(stack["name"], stack["args"])
}

func (c *current) onCallArguments1(head, tail interface{}) (interface{}, error) {
	args := []Expression{}

	if head != nil {
		args = append(args, head.(Expression))
	}

	if tail != nil {
		tailSlice := toSlice(tail)

		for _, arg := range tailSlice {
			argSlice := toSlice(arg)
			args = append(args, argSlice[3].(Expression))
		}
	}

	return args, nil
}

func (p *parser) callonCallArguments1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCallArguments1(stack["head"], stack["tail"])
}

func (c *current) onImport1(file interface{}) (interface{}, error) {
	return &Import{File: file.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonImport1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onImport1(stack["file"])
}

func (c *current) onExtend1(file interface{}) (interface{}, error) {
	return &Extend{File: file.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonExtend1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExtend1(stack["file"])
}

func (c *current) onBlock1(mod, name, list interface{}) (interface{}, error) {
	block := &Block{Name: name.(string), GraphNode: NewNode(pos(c.pos))}

	if mod != nil {
		modSlice := toSlice(mod)

		if string(modSlice[1].([]byte)) == "append" {
			block.Modifier = "append"
		} else {
			block.Modifier = "prepend"
		}
	}

	if list != nil {
		block.Block = list.(*List)
	}

	return block, nil
}

func (p *parser) callonBlock1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBlock1(stack["mod"], stack["name"], stack["list"])
}

func (c *current) onComment1(silent, comment interface{}) (interface{}, error) {
	isSilent := silent != nil
	return &Comment{Value: comment.(string), Silent: isSilent, GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonComment1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onComment1(stack["silent"], stack["comment"])
}

func (c *current) onLineText1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonLineText1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLineText1()
}

func (c *current) onPipeText1(tl interface{}) (interface{}, error) {
	return tl, nil
}

func (p *parser) callonPipeText1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPipeText1(stack["tl"])
}

func (c *current) onPipeExpression1(mod, ex interface{}) (interface{}, error) {
	intr := &Interpolation{Expr: ex.(Expression), GraphNode: NewNode(pos(c.pos))}

	if string(mod.([]byte)) == "!=" {
		intr.Unescaped = true
	}

	return intr, nil
}

func (p *parser) callonPipeExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onPipeExpression1(stack["mod"], stack["ex"])
}

func (c *current) onTextList2(intr, tl interface{}) (interface{}, error) {
	intNode := intr.(*Interpolation)

	if tl != nil {
		tlnode := tl.(*TextList)
		return &TextList{
			Nodes:     append([]Node{intNode}, tlnode.Nodes...),
			GraphNode: NewNode(pos(c.pos)),
		}, nil
	}

	return TextList{
		Nodes:     []Node{intNode},
		GraphNode: NewNode(pos(c.pos)),
	}, nil
}

func (p *parser) callonTextList2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTextList2(stack["intr"], stack["tl"])
}

func (c *current) onTextList8() (interface{}, error) {
	return &TextList{GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonTextList8() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTextList8()
}

func (c *current) onTextList11(ch, tl interface{}) (interface{}, error) {
	tlnode := tl.(*TextList)
	chstr := string(ch.([]byte))

	if len(tlnode.Nodes) > 0 {
		if tn, ok := tlnode.Nodes[0].(*Text); ok {
			tlnode.Nodes[0] = &Text{Value: chstr + tn.Value}
			return tlnode, nil
		}
	}

	tlnode.Nodes = append([]Node{&Text{Value: chstr, GraphNode: NewNode(pos(c.pos))}}, tlnode.Nodes...)
	tlnode.Position = pos(c.pos)

	return tlnode, nil
}

func (p *parser) callonTextList11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTextList11(stack["ch"], stack["tl"])
}

func (c *current) onInterpolation1(mod, expr interface{}) (interface{}, error) {
	intElem := &Interpolation{Expr: expr.(Expression), GraphNode: NewNode(pos(c.pos))}
	if string(mod.([]byte)) == "!" {
		intElem.Unescaped = true
	}
	return intElem, nil
}

func (p *parser) callonInterpolation1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInterpolation1(stack["mod"], stack["expr"])
}

func (c *current) onExpressionTernery1(cnd, rest interface{}) (interface{}, error) {
	if rest == nil {
		return cnd, nil
	}

	restSlice := toSlice(rest)

	return &BinaryExpression{
		X: &BinaryExpression{
			X:         cnd.(Expression),
			Y:         restSlice[3].(Expression),
			Op:        "&&",
			GraphNode: NewNode(pos(c.pos)),
		},
		Y:         restSlice[7].(Expression),
		Op:        "||",
		GraphNode: NewNode(pos(c.pos)),
	}, nil
}

func (p *parser) callonExpressionTernery1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionTernery1(stack["cnd"], stack["rest"])
}

func (c *current) onExpressionBinOp1(first, rest interface{}) (interface{}, error) {
	return binary(first, rest, c.pos)
}

func (p *parser) callonExpressionBinOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionBinOp1(stack["first"], stack["rest"])
}

func (c *current) onExpressionCmpOp1(first, rest interface{}) (interface{}, error) {
	return binary(first, rest, c.pos)
}

func (p *parser) callonExpressionCmpOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionCmpOp1(stack["first"], stack["rest"])
}

func (c *current) onExpressionAddOp1(first, rest interface{}) (interface{}, error) {
	return binary(first, rest, c.pos)
}

func (p *parser) callonExpressionAddOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionAddOp1(stack["first"], stack["rest"])
}

func (c *current) onExpressionMulOp1(first, rest interface{}) (interface{}, error) {
	return binary(first, rest, c.pos)
}

func (p *parser) callonExpressionMulOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionMulOp1(stack["first"], stack["rest"])
}

func (c *current) onExpressionUnaryOp2(op, ex interface{}) (interface{}, error) {
	return &UnaryExpression{X: ex.(Expression), Op: op.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonExpressionUnaryOp2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionUnaryOp2(stack["op"], stack["ex"])
}

func (c *current) onExpressionFactor2(e interface{}) (interface{}, error) {
	return e, nil
}

func (p *parser) callonExpressionFactor2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onExpressionFactor2(stack["e"])
}

func (c *current) onStringExpression1(s interface{}) (interface{}, error) {
	return &StringExpression{Value: s.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonStringExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStringExpression1(stack["s"])
}

func (c *current) onNumberExpression1(dec, ex interface{}) (interface{}, error) {
	if dec != nil || ex != nil {
		val, err := strconv.ParseFloat(string(c.text), 64)
		return &FloatExpression{Value: val, GraphNode: NewNode(pos(c.pos))}, err
	} else {
		val, err := strconv.ParseInt(string(c.text), 10, 64)
		return &IntegerExpression{Value: val, GraphNode: NewNode(pos(c.pos))}, err
	}
}

func (p *parser) callonNumberExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNumberExpression1(stack["dec"], stack["ex"])
}

func (c *current) onNilExpression1() (interface{}, error) {
	return &NilExpression{GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonNilExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNilExpression1()
}

func (c *current) onBooleanExpression1(b interface{}) (interface{}, error) {
	return &BooleanExpression{Value: b.(bool), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonBooleanExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBooleanExpression1(stack["b"])
}

func (c *current) onMemberExpression1(field, member interface{}) (interface{}, error) {
	memberSlice := toSlice(member)
	cur := field.(Expression)

	for _, m := range memberSlice {
		switch sub := m.(type) {
		case string:
			cur = &MemberExpression{X: cur, Name: sub, GraphNode: NewNode(pos(c.pos))}
		case Expression:
			cur = &IndexExpression{X: cur, Index: sub, GraphNode: NewNode(pos(c.pos))}
		case []Expression:
			cur = &FunctionCallExpression{X: cur, Arguments: sub, GraphNode: NewNode(pos(c.pos))}
		}
	}

	return cur, nil
}

func (p *parser) callonMemberExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMemberExpression1(stack["field"], stack["member"])
}

func (c *current) onMemberField1(ident interface{}) (interface{}, error) {
	return ident, nil
}

func (p *parser) callonMemberField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMemberField1(stack["ident"])
}

func (c *current) onMemberIndex1(i interface{}) (interface{}, error) {
	return i, nil
}

func (p *parser) callonMemberIndex1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMemberIndex1(stack["i"])
}

func (c *current) onMemberCall1(arg interface{}) (interface{}, error) {
	return arg, nil
}

func (p *parser) callonMemberCall1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMemberCall1(stack["arg"])
}

func (c *current) onArrayExpression1(head, tail interface{}) (interface{}, error) {
	expressions := []Expression{}
	if head != nil {
		expressions = append(expressions, head.(Expression))
	}
	tailSlice := toSlice(tail)
	for _, ex := range tailSlice {
		exSlice := toSlice(ex)
		expressions = append(expressions, exSlice[3].(Expression))
	}
	return &ArrayExpression{Expressions: expressions, GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonArrayExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onArrayExpression1(stack["head"], stack["tail"])
}

func (c *current) onField1(variable interface{}) (interface{}, error) {
	return &FieldExpression{Variable: variable.(*Variable), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["variable"])
}

func (c *current) onUnaryOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonUnaryOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onUnaryOp1()
}

func (c *current) onAddOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonAddOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onAddOp1()
}

func (c *current) onMulOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonMulOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMulOp1()
}

func (c *current) onCmpOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonCmpOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCmpOp1()
}

func (c *current) onBinOp1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonBinOp1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBinOp1()
}

func (c *current) onString1() (interface{}, error) {
	return strconv.Unquote(string(c.text))
}

func (p *parser) callonString1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onString1()
}

func (c *current) onIndex1() (interface{}, error) {
	return strconv.ParseInt(string(c.text), 10, 64)
}

func (p *parser) callonIndex1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIndex1()
}

func (c *current) onBool2() (interface{}, error) {
	return true, nil
}

func (p *parser) callonBool2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBool2()
}

func (c *current) onBool4() (interface{}, error) {
	return false, nil
}

func (p *parser) callonBool4() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBool4()
}

func (c *current) onNull1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonNull1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNull1()
}

func (c *current) onVariable1(ident interface{}) (interface{}, error) {
	return &Variable{Name: ident.(string), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonVariable1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onVariable1(stack["ident"])
}

func (c *current) onIdentifier1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1()
}

func (c *current) onEmptyLine1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonEmptyLine1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEmptyLine1()
}

func (c *current) on_1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callon_1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.on_1()
}

func (c *current) on__1() (interface{}, error) {
	return nil, nil
}

func (p *parser) callon__1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.on__1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errMaxExprCnt is used to signal that the maximum number of
	// expressions have been parsed.
	errMaxExprCnt = errors.New("max number of expresssions parsed")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// MaxExpressions creates an Option to stop parsing after the provided
// number of expressions have been parsed, if the value is 0 then the parser will
// parse for as many steps as needed (possibly an infinite number).
//
// The default for maxExprCnt is 0.
func MaxExpressions(maxExprCnt uint64) Option {
	return func(p *parser) Option {
		oldMaxExprCnt := p.maxExprCnt
		p.maxExprCnt = maxExprCnt
		return MaxExpressions(oldMaxExprCnt)
	}
}

// Statistics adds a user provided Stats struct to the parser to allow
// the user to process the results after the parsing has finished.
// Also the key for the "no match" counter is set.
//
// Example usage:
//
//     input := "input"
//     stats := Stats{}
//     _, err := Parse("input-file", []byte(input), Statistics(&stats, "no match"))
//     if err != nil {
//         log.Panicln(err)
//     }
//     b, err := json.MarshalIndent(stats.ChoiceAltCnt, "", "  ")
//     if err != nil {
//         log.Panicln(err)
//     }
//     fmt.Println(string(b))
//
func Statistics(stats *Stats, choiceNoMatch string) Option {
	return func(p *parser) Option {
		oldStats := p.Stats
		p.Stats = stats
		oldChoiceNoMatch := p.choiceNoMatch
		p.choiceNoMatch = choiceNoMatch
		if p.Stats.ChoiceAltCnt == nil {
			p.Stats.ChoiceAltCnt = make(map[string]map[string]int)
		}
		return Statistics(oldStats, oldChoiceNoMatch)
	}
}

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// GlobalStore creates an Option to set a key to a certain value in
// the globalStore.
func GlobalStore(key string, value interface{}) Option {
	return func(p *parser) Option {
		old := p.cur.globalStore[key]
		p.cur.globalStore[key] = value
		return GlobalStore(key, old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (i interface{}, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := f.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match

	// the globalStore allows the parser to store arbitrary values
	globalStore map[string]interface{}
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type recoveryExpr struct {
	pos          position
	expr         interface{}
	recoverExpr  interface{}
	failureLabel []string
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type throwExpr struct {
	pos   position
	label string
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos             position
	val             string
	basicLatinChars [128]bool
	chars           []rune
	ranges          []rune
	classes         []*unicode.RangeTable
	ignoreCase      bool
	inverted        bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner    error
	pos      position
	prefix   string
	expected []string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	stats := Stats{
		ChoiceAltCnt: make(map[string]map[string]int),
	}

	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
		cur: current{
			globalStore: make(map[string]interface{}),
		},
		maxFailPos:      position{col: 1, line: 1},
		maxFailExpected: make([]string, 0, 20),
		Stats:           &stats,
	}
	p.setOptions(opts)

	if p.maxExprCnt == 0 {
		p.maxExprCnt = math.MaxUint64
	}

	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

const choiceNoMatch = -1

// Stats stores some statistics, gathered during parsing
type Stats struct {
	// ExprCnt counts the number of expressions processed during parsing
	// This value is compared to the maximum number of expressions allowed
	// (set by the MaxExpressions option).
	ExprCnt uint64

	// ChoiceAltCnt is used to count for each ordered choice expression,
	// which alternative is used how may times.
	// These numbers allow to optimize the order of the ordered choice expression
	// to increase the performance of the parser
	//
	// The outer key of ChoiceAltCnt is composed of the name of the rule as well
	// as the line and the column of the ordered choice.
	// The inner key of ChoiceAltCnt is the number (one-based) of the matching alternative.
	// For each alternative the number of matches are counted. If an ordered choice does not
	// match, a special counter is incremented. The name of this counter is set with
	// the parser option Statistics.
	// For an alternative to be included in ChoiceAltCnt, it has to match at least once.
	ChoiceAltCnt map[string]map[string]int
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	depth   int
	recover bool
	debug   bool

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// parse fail
	maxFailPos            position
	maxFailExpected       []string
	maxFailInvertExpected bool

	// max number of expressions to be parsed
	maxExprCnt uint64

	*Stats

	choiceNoMatch string
	// recovery expression stack, keeps track of the currently available recovery expression, these are traversed in reverse
	recoveryStack []map[string]interface{}
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

// push a recovery expression with its labels to the recoveryStack
func (p *parser) pushRecovery(labels []string, expr interface{}) {
	if cap(p.recoveryStack) == len(p.recoveryStack) {
		// create new empty slot in the stack
		p.recoveryStack = append(p.recoveryStack, nil)
	} else {
		// slice to 1 more
		p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)+1]
	}

	m := make(map[string]interface{}, len(labels))
	for _, fl := range labels {
		m[fl] = expr
	}
	p.recoveryStack[len(p.recoveryStack)-1] = m
}

// pop a recovery expression from the recoveryStack
func (p *parser) popRecovery() {
	// GC that map
	p.recoveryStack[len(p.recoveryStack)-1] = nil

	p.recoveryStack = p.recoveryStack[:len(p.recoveryStack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position, []string{})
}

func (p *parser) addErrAt(err error, pos position, expected []string) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, pos: pos, prefix: buf.String(), expected: expected}
	p.errs.add(pe)
}

func (p *parser) failAt(fail bool, pos position, want string) {
	// process fail if parsing fails and not inverted or parsing succeeds and invert is set
	if fail == p.maxFailInvertExpected {
		if pos.offset < p.maxFailPos.offset {
			return
		}

		if pos.offset > p.maxFailPos.offset {
			p.maxFailPos = pos
			p.maxFailExpected = p.maxFailExpected[:0]
		}

		if p.maxFailInvertExpected {
			want = "!" + want
		}
		p.maxFailExpected = append(p.maxFailExpected, want)
	}
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n == 1 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// If parsing fails, but no errors have been recorded, the expected values
			// for the farthest parser position are returned as error.
			maxFailExpectedMap := make(map[string]struct{}, len(p.maxFailExpected))
			for _, v := range p.maxFailExpected {
				maxFailExpectedMap[v] = struct{}{}
			}
			expected := make([]string, 0, len(maxFailExpectedMap))
			eof := false
			if _, ok := maxFailExpectedMap["!."]; ok {
				delete(maxFailExpectedMap, "!.")
				eof = true
			}
			for k := range maxFailExpectedMap {
				expected = append(expected, k)
			}
			sort.Strings(expected)
			if eof {
				expected = append(expected, "EOF")
			}
			p.addErrAt(errors.New("no match found, expected: "+listJoin(expected, ", ", "or")), p.maxFailPos, expected)
		}

		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func listJoin(list []string, sep string, lastSep string) string {
	switch len(list) {
	case 0:
		return ""
	case 1:
		return list[0]
	default:
		return fmt.Sprintf("%s %s %s", strings.Join(list[:len(list)-1], sep), lastSep, list[len(list)-1])
	}
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.ExprCnt++
	if p.ExprCnt > p.maxExprCnt {
		panic(errMaxExprCnt)
	}

	var val interface{}
	var ok bool
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *recoveryExpr:
		val, ok = p.parseRecoveryExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *throwExpr:
		val, ok = p.parseThrowExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position, []string{})
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		p.failAt(true, start.position, ".")
		return p.sliceFrom(start), true
	}
	p.failAt(false, p.pt.position, ".")
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt

	// can't match EOF
	if cur == utf8.RuneError {
		p.failAt(false, start.position, chr.val)
		return nil, false
	}

	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				p.failAt(false, start.position, chr.val)
				return nil, false
			}
			p.read()
			p.failAt(true, start.position, chr.val)
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		p.failAt(true, start.position, chr.val)
		return p.sliceFrom(start), true
	}
	p.failAt(false, start.position, chr.val)
	return nil, false
}

func (p *parser) incChoiceAltCnt(ch *choiceExpr, altI int) {
	choiceIdent := fmt.Sprintf("%s %d:%d", p.rstack[len(p.rstack)-1].name, ch.pos.line, ch.pos.col)
	m := p.ChoiceAltCnt[choiceIdent]
	if m == nil {
		m = make(map[string]int)
		p.ChoiceAltCnt[choiceIdent] = m
	}
	// We increment altI by 1, so the keys do not start at 0
	alt := strconv.Itoa(altI + 1)
	if altI == choiceNoMatch {
		alt = p.choiceNoMatch
	}
	m[alt]++
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for altI, alt := range ch.alternatives {
		// dummy assignment to prevent compile error if optimized
		_ = altI

		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			p.incChoiceAltCnt(ch, altI)
			return val, ok
		}
	}
	p.incChoiceAltCnt(ch, choiceNoMatch)
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	ignoreCase := ""
	if lit.ignoreCase {
		ignoreCase = "i"
	}
	val := fmt.Sprintf("%q%s", lit.val, ignoreCase)
	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.failAt(false, start.position, val)
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	p.failAt(true, start.position, val)
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	_, ok := p.parseExpr(not.expr)
	p.maxFailInvertExpected = !p.maxFailInvertExpected
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRecoveryExpr(recover *recoveryExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRecoveryExpr (" + strings.Join(recover.failureLabel, ",") + ")"))
	}

	p.pushRecovery(recover.failureLabel, recover.recoverExpr)
	val, ok := p.parseExpr(recover.expr)
	p.popRecovery()

	return val, ok
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	vals := make([]interface{}, 0, len(seq.exprs))

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseThrowExpr(expr *throwExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseThrowExpr"))
	}

	for i := len(p.recoveryStack) - 1; i >= 0; i-- {
		if recoverExpr, ok := p.recoveryStack[i][expr.label]; ok {
			if val, ok := p.parseExpr(recoverExpr); ok {
				return val, ok
			}
		}
	}

	return nil, false
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}
