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
						name: "If",
					},
					&ruleRefExpr{
						pos:  position{line: 75, col: 3, offset: 1569},
						name: "Each",
					},
					&ruleRefExpr{
						pos:  position{line: 76, col: 3, offset: 1578},
						name: "DocType",
					},
					&ruleRefExpr{
						pos:  position{line: 77, col: 3, offset: 1590},
						name: "Mixin",
					},
					&ruleRefExpr{
						pos:  position{line: 78, col: 3, offset: 1600},
						name: "MixinCall",
					},
					&ruleRefExpr{
						pos:  position{line: 79, col: 3, offset: 1614},
						name: "Assignment",
					},
					&ruleRefExpr{
						pos:  position{line: 80, col: 3, offset: 1629},
						name: "Block",
					},
					&ruleRefExpr{
						pos:  position{line: 81, col: 3, offset: 1639},
						name: "Tag",
					},
					&actionExpr{
						pos: position{line: 82, col: 3, offset: 1647},
						run: (*parser).callonListNode14,
						expr: &seqExpr{
							pos: position{line: 82, col: 4, offset: 1648},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 82, col: 4, offset: 1648},
									name: "_",
								},
								&charClassMatcher{
									pos:        position{line: 82, col: 6, offset: 1650},
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
			pos:         position{line: 85, col: 1, offset: 1685},
			expr: &actionExpr{
				pos: position{line: 85, col: 21, offset: 1705},
				run: (*parser).callonDocType1,
				expr: &seqExpr{
					pos: position{line: 85, col: 21, offset: 1705},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 85, col: 21, offset: 1705},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 85, col: 23, offset: 1707},
							val:        "doctype",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 85, col: 33, offset: 1717},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 85, col: 35, offset: 1719},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 85, col: 39, offset: 1723},
								name: "LineText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 85, col: 48, offset: 1732},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Tag",
			pos:  position{line: 91, col: 1, offset: 1825},
			expr: &actionExpr{
				pos: position{line: 91, col: 8, offset: 1832},
				run: (*parser).callonTag1,
				expr: &seqExpr{
					pos: position{line: 91, col: 8, offset: 1832},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 91, col: 8, offset: 1832},
							label: "tag",
							expr: &ruleRefExpr{
								pos:  position{line: 91, col: 12, offset: 1836},
								name: "TagHeader",
							},
						},
						&labeledExpr{
							pos:   position{line: 91, col: 22, offset: 1846},
							label: "list",
							expr: &zeroOrOneExpr{
								pos: position{line: 91, col: 27, offset: 1851},
								expr: &ruleRefExpr{
									pos:  position{line: 91, col: 27, offset: 1851},
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
			pos:  position{line: 101, col: 1, offset: 1974},
			expr: &choiceExpr{
				pos: position{line: 101, col: 14, offset: 1987},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 101, col: 14, offset: 1987},
						run: (*parser).callonTagHeader2,
						expr: &seqExpr{
							pos: position{line: 101, col: 14, offset: 1987},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 101, col: 14, offset: 1987},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 101, col: 16, offset: 1989},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 101, col: 21, offset: 1994},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 101, col: 29, offset: 2002},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 101, col: 35, offset: 2008},
										expr: &ruleRefExpr{
											pos:  position{line: 101, col: 35, offset: 2008},
											name: "TagAttributes",
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 101, col: 50, offset: 2023},
									label: "tl",
									expr: &zeroOrOneExpr{
										pos: position{line: 101, col: 53, offset: 2026},
										expr: &seqExpr{
											pos: position{line: 101, col: 54, offset: 2027},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 101, col: 54, offset: 2027},
													name: "__",
												},
												&zeroOrOneExpr{
													pos: position{line: 101, col: 57, offset: 2030},
													expr: &ruleRefExpr{
														pos:  position{line: 101, col: 57, offset: 2030},
														name: "TextList",
													},
												},
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 101, col: 69, offset: 2042},
									name: "NL",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 110, col: 5, offset: 2265},
						run: (*parser).callonTagHeader17,
						expr: &seqExpr{
							pos: position{line: 110, col: 5, offset: 2265},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 110, col: 5, offset: 2265},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 110, col: 7, offset: 2267},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 12, offset: 2272},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 110, col: 20, offset: 2280},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 110, col: 26, offset: 2286},
										expr: &ruleRefExpr{
											pos:  position{line: 110, col: 26, offset: 2286},
											name: "TagAttributes",
										},
									},
								},
								&litMatcher{
									pos:        position{line: 110, col: 41, offset: 2301},
									val:        ".",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 110, col: 45, offset: 2305},
									name: "NL",
								},
								&labeledExpr{
									pos:   position{line: 110, col: 48, offset: 2308},
									label: "text",
									expr: &zeroOrOneExpr{
										pos: position{line: 110, col: 53, offset: 2313},
										expr: &ruleRefExpr{
											pos:  position{line: 110, col: 53, offset: 2313},
											name: "IndentedRawText",
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 119, col: 5, offset: 2539},
						run: (*parser).callonTagHeader30,
						expr: &seqExpr{
							pos: position{line: 119, col: 5, offset: 2539},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 119, col: 5, offset: 2539},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 119, col: 7, offset: 2541},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 119, col: 12, offset: 2546},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 119, col: 20, offset: 2554},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 119, col: 26, offset: 2560},
										expr: &ruleRefExpr{
											pos:  position{line: 119, col: 26, offset: 2560},
											name: "TagAttributes",
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 119, col: 41, offset: 2575},
									label: "unescaped",
									expr: &zeroOrOneExpr{
										pos: position{line: 119, col: 51, offset: 2585},
										expr: &litMatcher{
											pos:        position{line: 119, col: 51, offset: 2585},
											val:        "!",
											ignoreCase: false,
										},
									},
								},
								&litMatcher{
									pos:        position{line: 119, col: 56, offset: 2590},
									val:        "=",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 119, col: 60, offset: 2594},
									name: "__",
								},
								&labeledExpr{
									pos:   position{line: 119, col: 63, offset: 2597},
									label: "expr",
									expr: &zeroOrOneExpr{
										pos: position{line: 119, col: 68, offset: 2602},
										expr: &ruleRefExpr{
											pos:  position{line: 119, col: 68, offset: 2602},
											name: "Expression",
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 119, col: 80, offset: 2614},
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
			pos:  position{line: 134, col: 1, offset: 3019},
			expr: &choiceExpr{
				pos: position{line: 134, col: 12, offset: 3030},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 134, col: 12, offset: 3030},
						run: (*parser).callonTagName2,
						expr: &seqExpr{
							pos: position{line: 134, col: 12, offset: 3030},
							exprs: []interface{}{
								&charClassMatcher{
									pos:        position{line: 134, col: 12, offset: 3030},
									val:        "[_a-zA-Z]",
									chars:      []rune{'_'},
									ranges:     []rune{'a', 'z', 'A', 'Z'},
									ignoreCase: false,
									inverted:   false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 134, col: 22, offset: 3040},
									expr: &charClassMatcher{
										pos:        position{line: 134, col: 22, offset: 3040},
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
						pos: position{line: 136, col: 5, offset: 3091},
						expr: &ruleRefExpr{
							pos:  position{line: 136, col: 6, offset: 3092},
							name: "TagAttributeClass",
						},
					},
					&actionExpr{
						pos: position{line: 136, col: 26, offset: 3112},
						run: (*parser).callonTagName9,
						expr: &andExpr{
							pos: position{line: 136, col: 26, offset: 3112},
							expr: &ruleRefExpr{
								pos:  position{line: 136, col: 27, offset: 3113},
								name: "TagAttributeID",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributes",
			pos:  position{line: 140, col: 1, offset: 3153},
			expr: &choiceExpr{
				pos: position{line: 140, col: 18, offset: 3170},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 140, col: 18, offset: 3170},
						run: (*parser).callonTagAttributes2,
						expr: &seqExpr{
							pos: position{line: 140, col: 18, offset: 3170},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 140, col: 18, offset: 3170},
									label: "head",
									expr: &choiceExpr{
										pos: position{line: 140, col: 24, offset: 3176},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 140, col: 24, offset: 3176},
												name: "TagAttributeClass",
											},
											&ruleRefExpr{
												pos:  position{line: 140, col: 44, offset: 3196},
												name: "TagAttributeID",
											},
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 140, col: 60, offset: 3212},
									label: "tail",
									expr: &zeroOrOneExpr{
										pos: position{line: 140, col: 65, offset: 3217},
										expr: &ruleRefExpr{
											pos:  position{line: 140, col: 65, offset: 3217},
											name: "TagAttributes",
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 148, col: 5, offset: 3382},
						run: (*parser).callonTagAttributes11,
						expr: &seqExpr{
							pos: position{line: 148, col: 5, offset: 3382},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 148, col: 5, offset: 3382},
									label: "head",
									expr: &seqExpr{
										pos: position{line: 148, col: 11, offset: 3388},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 148, col: 11, offset: 3388},
												val:        "(",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 148, col: 15, offset: 3392},
												name: "_",
											},
											&seqExpr{
												pos: position{line: 148, col: 18, offset: 3395},
												exprs: []interface{}{
													&ruleRefExpr{
														pos:  position{line: 148, col: 18, offset: 3395},
														name: "TagAttribute",
													},
													&zeroOrMoreExpr{
														pos: position{line: 148, col: 31, offset: 3408},
														expr: &seqExpr{
															pos: position{line: 148, col: 32, offset: 3409},
															exprs: []interface{}{
																&choiceExpr{
																	pos: position{line: 148, col: 33, offset: 3410},
																	alternatives: []interface{}{
																		&ruleRefExpr{
																			pos:  position{line: 148, col: 33, offset: 3410},
																			name: "__",
																		},
																		&seqExpr{
																			pos: position{line: 148, col: 39, offset: 3416},
																			exprs: []interface{}{
																				&ruleRefExpr{
																					pos:  position{line: 148, col: 39, offset: 3416},
																					name: "_",
																				},
																				&litMatcher{
																					pos:        position{line: 148, col: 41, offset: 3418},
																					val:        ",",
																					ignoreCase: false,
																				},
																				&ruleRefExpr{
																					pos:  position{line: 148, col: 45, offset: 3422},
																					name: "_",
																				},
																			},
																		},
																	},
																},
																&ruleRefExpr{
																	pos:  position{line: 148, col: 49, offset: 3426},
																	name: "TagAttribute",
																},
															},
														},
													},
												},
											},
											&ruleRefExpr{
												pos:  position{line: 148, col: 65, offset: 3442},
												name: "_",
											},
											&litMatcher{
												pos:        position{line: 148, col: 67, offset: 3444},
												val:        ")",
												ignoreCase: false,
											},
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 148, col: 72, offset: 3449},
									label: "tail",
									expr: &zeroOrOneExpr{
										pos: position{line: 148, col: 77, offset: 3454},
										expr: &ruleRefExpr{
											pos:  position{line: 148, col: 77, offset: 3454},
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
			pos:  position{line: 172, col: 1, offset: 3900},
			expr: &actionExpr{
				pos: position{line: 172, col: 22, offset: 3921},
				run: (*parser).callonTagAttributeClass1,
				expr: &seqExpr{
					pos: position{line: 172, col: 22, offset: 3921},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 172, col: 22, offset: 3921},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 172, col: 26, offset: 3925},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 172, col: 31, offset: 3930},
								name: "ClassName",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributeID",
			pos:  position{line: 176, col: 1, offset: 4079},
			expr: &actionExpr{
				pos: position{line: 176, col: 19, offset: 4097},
				run: (*parser).callonTagAttributeID1,
				expr: &seqExpr{
					pos: position{line: 176, col: 19, offset: 4097},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 176, col: 19, offset: 4097},
							val:        "#",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 176, col: 23, offset: 4101},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 176, col: 28, offset: 4106},
								name: "IdName",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttribute",
			pos:  position{line: 180, col: 1, offset: 4249},
			expr: &choiceExpr{
				pos: position{line: 180, col: 17, offset: 4265},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 180, col: 17, offset: 4265},
						run: (*parser).callonTagAttribute2,
						expr: &seqExpr{
							pos: position{line: 180, col: 17, offset: 4265},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 180, col: 17, offset: 4265},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 180, col: 22, offset: 4270},
										name: "TagAttributeName",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 180, col: 39, offset: 4287},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 180, col: 41, offset: 4289},
									val:        "=",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 180, col: 45, offset: 4293},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 180, col: 47, offset: 4295},
									label: "value",
									expr: &ruleRefExpr{
										pos:  position{line: 180, col: 53, offset: 4301},
										name: "Expression",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 182, col: 5, offset: 4437},
						run: (*parser).callonTagAttribute11,
						expr: &seqExpr{
							pos: position{line: 182, col: 5, offset: 4437},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 182, col: 5, offset: 4437},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 182, col: 10, offset: 4442},
										name: "TagAttributeName",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 182, col: 27, offset: 4459},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 182, col: 29, offset: 4461},
									val:        "!=",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 182, col: 34, offset: 4466},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 182, col: 36, offset: 4468},
									label: "value",
									expr: &ruleRefExpr{
										pos:  position{line: 182, col: 42, offset: 4474},
										name: "Expression",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 184, col: 5, offset: 4627},
						run: (*parser).callonTagAttribute20,
						expr: &labeledExpr{
							pos:   position{line: 184, col: 5, offset: 4627},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 184, col: 10, offset: 4632},
								name: "TagAttributeName",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributeName",
			pos:  position{line: 188, col: 1, offset: 4746},
			expr: &choiceExpr{
				pos: position{line: 188, col: 21, offset: 4766},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 188, col: 21, offset: 4766},
						run: (*parser).callonTagAttributeName2,
						expr: &seqExpr{
							pos: position{line: 188, col: 21, offset: 4766},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 188, col: 21, offset: 4766},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 188, col: 25, offset: 4770},
									label: "tn",
									expr: &ruleRefExpr{
										pos:  position{line: 188, col: 28, offset: 4773},
										name: "TagAttributeNameLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 188, col: 52, offset: 4797},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 190, col: 5, offset: 4824},
						run: (*parser).callonTagAttributeName8,
						expr: &seqExpr{
							pos: position{line: 190, col: 5, offset: 4824},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 190, col: 5, offset: 4824},
									val:        "[",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 190, col: 9, offset: 4828},
									label: "tn",
									expr: &ruleRefExpr{
										pos:  position{line: 190, col: 12, offset: 4831},
										name: "TagAttributeNameLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 190, col: 36, offset: 4855},
									val:        "]",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 192, col: 5, offset: 4882},
						run: (*parser).callonTagAttributeName14,
						expr: &labeledExpr{
							pos:   position{line: 192, col: 5, offset: 4882},
							label: "tn",
							expr: &ruleRefExpr{
								pos:  position{line: 192, col: 8, offset: 4885},
								name: "TagAttributeNameLiteral",
							},
						},
					},
				},
			},
		},
		{
			name: "ClassName",
			pos:  position{line: 196, col: 1, offset: 4931},
			expr: &actionExpr{
				pos: position{line: 196, col: 14, offset: 4944},
				run: (*parser).callonClassName1,
				expr: &oneOrMoreExpr{
					pos: position{line: 196, col: 14, offset: 4944},
					expr: &charClassMatcher{
						pos:        position{line: 196, col: 14, offset: 4944},
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
			pos:  position{line: 200, col: 1, offset: 4994},
			expr: &actionExpr{
				pos: position{line: 200, col: 11, offset: 5004},
				run: (*parser).callonIdName1,
				expr: &oneOrMoreExpr{
					pos: position{line: 200, col: 11, offset: 5004},
					expr: &charClassMatcher{
						pos:        position{line: 200, col: 11, offset: 5004},
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
			pos:  position{line: 204, col: 1, offset: 5054},
			expr: &actionExpr{
				pos: position{line: 204, col: 28, offset: 5081},
				run: (*parser).callonTagAttributeNameLiteral1,
				expr: &seqExpr{
					pos: position{line: 204, col: 28, offset: 5081},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 204, col: 28, offset: 5081},
							val:        "[@_a-zA-Z]",
							chars:      []rune{'@', '_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 204, col: 39, offset: 5092},
							expr: &charClassMatcher{
								pos:        position{line: 204, col: 39, offset: 5092},
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
			pos:  position{line: 209, col: 1, offset: 5152},
			expr: &actionExpr{
				pos: position{line: 209, col: 7, offset: 5158},
				run: (*parser).callonIf1,
				expr: &seqExpr{
					pos: position{line: 209, col: 7, offset: 5158},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 209, col: 7, offset: 5158},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 209, col: 9, offset: 5160},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 209, col: 14, offset: 5165},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 209, col: 17, offset: 5168},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 209, col: 22, offset: 5173},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 209, col: 33, offset: 5184},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 209, col: 35, offset: 5186},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 209, col: 38, offset: 5189},
							label: "block",
							expr: &ruleRefExpr{
								pos:  position{line: 209, col: 44, offset: 5195},
								name: "IndentedList",
							},
						},
						&labeledExpr{
							pos:   position{line: 209, col: 57, offset: 5208},
							label: "elseNode",
							expr: &zeroOrOneExpr{
								pos: position{line: 209, col: 66, offset: 5217},
								expr: &ruleRefExpr{
									pos:  position{line: 209, col: 66, offset: 5217},
									name: "Else",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Else",
			pos:  position{line: 217, col: 1, offset: 5426},
			expr: &choiceExpr{
				pos: position{line: 217, col: 9, offset: 5434},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 217, col: 9, offset: 5434},
						run: (*parser).callonElse2,
						expr: &seqExpr{
							pos: position{line: 217, col: 9, offset: 5434},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 217, col: 9, offset: 5434},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 217, col: 11, offset: 5436},
									val:        "else",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 217, col: 18, offset: 5443},
									label: "node",
									expr: &ruleRefExpr{
										pos:  position{line: 217, col: 23, offset: 5448},
										name: "If",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 219, col: 5, offset: 5476},
						run: (*parser).callonElse8,
						expr: &seqExpr{
							pos: position{line: 219, col: 5, offset: 5476},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 219, col: 5, offset: 5476},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 219, col: 7, offset: 5478},
									val:        "else",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 219, col: 14, offset: 5485},
									name: "_",
								},
								&ruleRefExpr{
									pos:  position{line: 219, col: 16, offset: 5487},
									name: "NL",
								},
								&labeledExpr{
									pos:   position{line: 219, col: 19, offset: 5490},
									label: "block",
									expr: &ruleRefExpr{
										pos:  position{line: 219, col: 25, offset: 5496},
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
			pos:  position{line: 223, col: 1, offset: 5534},
			expr: &actionExpr{
				pos: position{line: 223, col: 9, offset: 5542},
				run: (*parser).callonEach1,
				expr: &seqExpr{
					pos: position{line: 223, col: 9, offset: 5542},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 223, col: 9, offset: 5542},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 223, col: 11, offset: 5544},
							val:        "each",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 223, col: 18, offset: 5551},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 223, col: 21, offset: 5554},
							label: "v1",
							expr: &ruleRefExpr{
								pos:  position{line: 223, col: 24, offset: 5557},
								name: "Variable",
							},
						},
						&labeledExpr{
							pos:   position{line: 223, col: 33, offset: 5566},
							label: "v2",
							expr: &zeroOrOneExpr{
								pos: position{line: 223, col: 36, offset: 5569},
								expr: &seqExpr{
									pos: position{line: 223, col: 37, offset: 5570},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 223, col: 37, offset: 5570},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 223, col: 39, offset: 5572},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 223, col: 43, offset: 5576},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 223, col: 45, offset: 5578},
											name: "Variable",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 223, col: 56, offset: 5589},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 223, col: 58, offset: 5591},
							val:        "in",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 223, col: 63, offset: 5596},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 223, col: 65, offset: 5598},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 223, col: 70, offset: 5603},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 223, col: 81, offset: 5614},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 223, col: 83, offset: 5616},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 223, col: 86, offset: 5619},
							label: "block",
							expr: &ruleRefExpr{
								pos:  position{line: 223, col: 92, offset: 5625},
								name: "IndentedList",
							},
						},
					},
				},
			},
		},
		{
			name: "Assignment",
			pos:  position{line: 234, col: 1, offset: 5910},
			expr: &actionExpr{
				pos: position{line: 234, col: 15, offset: 5924},
				run: (*parser).callonAssignment1,
				expr: &seqExpr{
					pos: position{line: 234, col: 15, offset: 5924},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 234, col: 15, offset: 5924},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 234, col: 17, offset: 5926},
							val:        "-",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 234, col: 21, offset: 5930},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 234, col: 24, offset: 5933},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 234, col: 24, offset: 5933},
									val:        "var",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 234, col: 32, offset: 5941},
									val:        "let",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 234, col: 40, offset: 5949},
									val:        "const",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 234, col: 49, offset: 5958},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 234, col: 52, offset: 5961},
							label: "vr",
							expr: &ruleRefExpr{
								pos:  position{line: 234, col: 55, offset: 5964},
								name: "Variable",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 234, col: 64, offset: 5973},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 234, col: 66, offset: 5975},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 234, col: 70, offset: 5979},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 234, col: 72, offset: 5981},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 234, col: 77, offset: 5986},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 234, col: 88, offset: 5997},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 234, col: 90, offset: 5999},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Mixin",
			pos:  position{line: 239, col: 1, offset: 6131},
			expr: &actionExpr{
				pos: position{line: 239, col: 10, offset: 6140},
				run: (*parser).callonMixin1,
				expr: &seqExpr{
					pos: position{line: 239, col: 10, offset: 6140},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 239, col: 10, offset: 6140},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 239, col: 12, offset: 6142},
							val:        "mixin",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 239, col: 20, offset: 6150},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 239, col: 23, offset: 6153},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 239, col: 28, offset: 6158},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 239, col: 39, offset: 6169},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 239, col: 41, offset: 6171},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 239, col: 46, offset: 6176},
								expr: &ruleRefExpr{
									pos:  position{line: 239, col: 46, offset: 6176},
									name: "MixinArguments",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 239, col: 62, offset: 6192},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 239, col: 65, offset: 6195},
							label: "list",
							expr: &ruleRefExpr{
								pos:  position{line: 239, col: 70, offset: 6200},
								name: "IndentedList",
							},
						},
					},
				},
			},
		},
		{
			name: "MixinArguments",
			pos:  position{line: 247, col: 1, offset: 6409},
			expr: &choiceExpr{
				pos: position{line: 247, col: 19, offset: 6427},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 247, col: 19, offset: 6427},
						run: (*parser).callonMixinArguments2,
						expr: &seqExpr{
							pos: position{line: 247, col: 19, offset: 6427},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 247, col: 19, offset: 6427},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 247, col: 23, offset: 6431},
									label: "head",
									expr: &ruleRefExpr{
										pos:  position{line: 247, col: 28, offset: 6436},
										name: "MixinArgument",
									},
								},
								&labeledExpr{
									pos:   position{line: 247, col: 42, offset: 6450},
									label: "tail",
									expr: &zeroOrMoreExpr{
										pos: position{line: 247, col: 47, offset: 6455},
										expr: &seqExpr{
											pos: position{line: 247, col: 48, offset: 6456},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 247, col: 48, offset: 6456},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 247, col: 50, offset: 6458},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 247, col: 54, offset: 6462},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 247, col: 56, offset: 6464},
													name: "MixinArgument",
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 247, col: 72, offset: 6480},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 260, col: 5, offset: 6742},
						run: (*parser).callonMixinArguments15,
						expr: &seqExpr{
							pos: position{line: 260, col: 5, offset: 6742},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 260, col: 5, offset: 6742},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 260, col: 9, offset: 6746},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 260, col: 11, offset: 6748},
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
			pos:  position{line: 264, col: 1, offset: 6775},
			expr: &actionExpr{
				pos: position{line: 264, col: 18, offset: 6792},
				run: (*parser).callonMixinArgument1,
				expr: &seqExpr{
					pos: position{line: 264, col: 18, offset: 6792},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 264, col: 18, offset: 6792},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 264, col: 23, offset: 6797},
								name: "Variable",
							},
						},
						&labeledExpr{
							pos:   position{line: 264, col: 32, offset: 6806},
							label: "def",
							expr: &zeroOrOneExpr{
								pos: position{line: 264, col: 36, offset: 6810},
								expr: &seqExpr{
									pos: position{line: 264, col: 37, offset: 6811},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 264, col: 37, offset: 6811},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 264, col: 39, offset: 6813},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 264, col: 43, offset: 6817},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 264, col: 45, offset: 6819},
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
			pos:  position{line: 275, col: 1, offset: 7043},
			expr: &actionExpr{
				pos: position{line: 275, col: 14, offset: 7056},
				run: (*parser).callonMixinCall1,
				expr: &seqExpr{
					pos: position{line: 275, col: 14, offset: 7056},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 275, col: 14, offset: 7056},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 275, col: 16, offset: 7058},
							val:        "+",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 275, col: 20, offset: 7062},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 275, col: 25, offset: 7067},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 275, col: 36, offset: 7078},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 275, col: 38, offset: 7080},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 275, col: 43, offset: 7085},
								expr: &ruleRefExpr{
									pos:  position{line: 275, col: 43, offset: 7085},
									name: "CallArguments",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 275, col: 58, offset: 7100},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "CallArguments",
			pos:  position{line: 283, col: 1, offset: 7271},
			expr: &actionExpr{
				pos: position{line: 283, col: 18, offset: 7288},
				run: (*parser).callonCallArguments1,
				expr: &seqExpr{
					pos: position{line: 283, col: 18, offset: 7288},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 283, col: 18, offset: 7288},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 283, col: 22, offset: 7292},
							label: "head",
							expr: &zeroOrOneExpr{
								pos: position{line: 283, col: 27, offset: 7297},
								expr: &ruleRefExpr{
									pos:  position{line: 283, col: 27, offset: 7297},
									name: "Expression",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 283, col: 39, offset: 7309},
							label: "tail",
							expr: &zeroOrMoreExpr{
								pos: position{line: 283, col: 44, offset: 7314},
								expr: &seqExpr{
									pos: position{line: 283, col: 45, offset: 7315},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 283, col: 45, offset: 7315},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 283, col: 47, offset: 7317},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 283, col: 51, offset: 7321},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 283, col: 53, offset: 7323},
											name: "Expression",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 283, col: 66, offset: 7336},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Import",
			pos:  position{line: 304, col: 1, offset: 7658},
			expr: &actionExpr{
				pos: position{line: 304, col: 11, offset: 7668},
				run: (*parser).callonImport1,
				expr: &seqExpr{
					pos: position{line: 304, col: 11, offset: 7668},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 304, col: 11, offset: 7668},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 304, col: 13, offset: 7670},
							val:        "include",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 304, col: 23, offset: 7680},
							expr: &litMatcher{
								pos:        position{line: 304, col: 23, offset: 7680},
								val:        "s",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 304, col: 28, offset: 7685},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 304, col: 31, offset: 7688},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 304, col: 36, offset: 7693},
								name: "String",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 304, col: 43, offset: 7700},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Extend",
			pos:  position{line: 308, col: 1, offset: 7783},
			expr: &actionExpr{
				pos: position{line: 308, col: 11, offset: 7793},
				run: (*parser).callonExtend1,
				expr: &seqExpr{
					pos: position{line: 308, col: 11, offset: 7793},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 308, col: 11, offset: 7793},
							val:        "extend",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 308, col: 20, offset: 7802},
							expr: &litMatcher{
								pos:        position{line: 308, col: 20, offset: 7802},
								val:        "s",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 308, col: 25, offset: 7807},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 308, col: 28, offset: 7810},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 308, col: 33, offset: 7815},
								name: "String",
							},
						},
					},
				},
			},
		},
		{
			name: "Block",
			pos:  position{line: 312, col: 1, offset: 7902},
			expr: &actionExpr{
				pos: position{line: 312, col: 10, offset: 7911},
				run: (*parser).callonBlock1,
				expr: &seqExpr{
					pos: position{line: 312, col: 10, offset: 7911},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 312, col: 10, offset: 7911},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 312, col: 12, offset: 7913},
							val:        "block",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 312, col: 20, offset: 7921},
							label: "mod",
							expr: &zeroOrOneExpr{
								pos: position{line: 312, col: 24, offset: 7925},
								expr: &seqExpr{
									pos: position{line: 312, col: 25, offset: 7926},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 312, col: 25, offset: 7926},
											name: "__",
										},
										&choiceExpr{
											pos: position{line: 312, col: 29, offset: 7930},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 312, col: 29, offset: 7930},
													val:        "append",
													ignoreCase: false,
												},
												&litMatcher{
													pos:        position{line: 312, col: 40, offset: 7941},
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
							pos:  position{line: 312, col: 53, offset: 7954},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 312, col: 56, offset: 7957},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 312, col: 61, offset: 7962},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 312, col: 72, offset: 7973},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 312, col: 75, offset: 7976},
							label: "list",
							expr: &zeroOrOneExpr{
								pos: position{line: 312, col: 80, offset: 7981},
								expr: &ruleRefExpr{
									pos:  position{line: 312, col: 80, offset: 7981},
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
			pos:  position{line: 333, col: 1, offset: 8346},
			expr: &actionExpr{
				pos: position{line: 333, col: 12, offset: 8357},
				run: (*parser).callonComment1,
				expr: &seqExpr{
					pos: position{line: 333, col: 12, offset: 8357},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 333, col: 12, offset: 8357},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 333, col: 14, offset: 8359},
							val:        "//",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 333, col: 19, offset: 8364},
							label: "silent",
							expr: &zeroOrOneExpr{
								pos: position{line: 333, col: 26, offset: 8371},
								expr: &litMatcher{
									pos:        position{line: 333, col: 26, offset: 8371},
									val:        "-",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 333, col: 31, offset: 8376},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 333, col: 33, offset: 8378},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 333, col: 41, offset: 8386},
								name: "LineText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 333, col: 50, offset: 8395},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "LineText",
			pos:  position{line: 338, col: 1, offset: 8529},
			expr: &actionExpr{
				pos: position{line: 338, col: 13, offset: 8541},
				run: (*parser).callonLineText1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 338, col: 13, offset: 8541},
					expr: &charClassMatcher{
						pos:        position{line: 338, col: 13, offset: 8541},
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
			pos:  position{line: 343, col: 1, offset: 8590},
			expr: &actionExpr{
				pos: position{line: 343, col: 13, offset: 8602},
				run: (*parser).callonPipeText1,
				expr: &seqExpr{
					pos: position{line: 343, col: 13, offset: 8602},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 343, col: 13, offset: 8602},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 343, col: 16, offset: 8605},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 343, col: 16, offset: 8605},
									val:        "|",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 343, col: 22, offset: 8611},
									val:        "<",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 343, col: 27, offset: 8616},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 343, col: 29, offset: 8618},
							label: "tl",
							expr: &ruleRefExpr{
								pos:  position{line: 343, col: 32, offset: 8621},
								name: "TextList",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 343, col: 41, offset: 8630},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "TextList",
			pos:  position{line: 347, col: 1, offset: 8655},
			expr: &choiceExpr{
				pos: position{line: 347, col: 13, offset: 8667},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 347, col: 13, offset: 8667},
						run: (*parser).callonTextList2,
						expr: &seqExpr{
							pos: position{line: 347, col: 13, offset: 8667},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 347, col: 13, offset: 8667},
									label: "intr",
									expr: &ruleRefExpr{
										pos:  position{line: 347, col: 18, offset: 8672},
										name: "Interpolation",
									},
								},
								&labeledExpr{
									pos:   position{line: 347, col: 32, offset: 8686},
									label: "tl",
									expr: &ruleRefExpr{
										pos:  position{line: 347, col: 35, offset: 8689},
										name: "TextList",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 362, col: 5, offset: 9010},
						run: (*parser).callonTextList8,
						expr: &andExpr{
							pos: position{line: 362, col: 5, offset: 9010},
							expr: &ruleRefExpr{
								pos:  position{line: 362, col: 6, offset: 9011},
								name: "NL",
							},
						},
					},
					&actionExpr{
						pos: position{line: 364, col: 5, offset: 9076},
						run: (*parser).callonTextList11,
						expr: &seqExpr{
							pos: position{line: 364, col: 5, offset: 9076},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 364, col: 5, offset: 9076},
									label: "ch",
									expr: &anyMatcher{
										line: 364, col: 8, offset: 9079,
									},
								},
								&labeledExpr{
									pos:   position{line: 364, col: 10, offset: 9081},
									label: "tl",
									expr: &ruleRefExpr{
										pos:  position{line: 364, col: 13, offset: 9084},
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
			pos:  position{line: 381, col: 1, offset: 9480},
			expr: &litMatcher{
				pos:        position{line: 381, col: 11, offset: 9490},
				val:        "\x01",
				ignoreCase: false,
			},
		},
		{
			name: "Outdent",
			pos:  position{line: 382, col: 1, offset: 9499},
			expr: &litMatcher{
				pos:        position{line: 382, col: 12, offset: 9510},
				val:        "\x02",
				ignoreCase: false,
			},
		},
		{
			name: "Interpolation",
			pos:  position{line: 384, col: 1, offset: 9520},
			expr: &actionExpr{
				pos: position{line: 384, col: 18, offset: 9537},
				run: (*parser).callonInterpolation1,
				expr: &seqExpr{
					pos: position{line: 384, col: 18, offset: 9537},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 384, col: 18, offset: 9537},
							label: "mod",
							expr: &choiceExpr{
								pos: position{line: 384, col: 23, offset: 9542},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 384, col: 23, offset: 9542},
										val:        "#",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 384, col: 29, offset: 9548},
										val:        "!",
										ignoreCase: false,
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 384, col: 34, offset: 9553},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 384, col: 38, offset: 9557},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 384, col: 40, offset: 9559},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 384, col: 45, offset: 9564},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 384, col: 56, offset: 9575},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 384, col: 58, offset: 9577},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 392, col: 1, offset: 9761},
			expr: &ruleRefExpr{
				pos:  position{line: 392, col: 15, offset: 9775},
				name: "ExpressionTernery",
			},
		},
		{
			name: "ExpressionTernery",
			pos:  position{line: 394, col: 1, offset: 9794},
			expr: &actionExpr{
				pos: position{line: 394, col: 22, offset: 9815},
				run: (*parser).callonExpressionTernery1,
				expr: &seqExpr{
					pos: position{line: 394, col: 22, offset: 9815},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 394, col: 22, offset: 9815},
							label: "cnd",
							expr: &ruleRefExpr{
								pos:  position{line: 394, col: 26, offset: 9819},
								name: "ExpressionBinOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 394, col: 42, offset: 9835},
							label: "rest",
							expr: &zeroOrOneExpr{
								pos: position{line: 394, col: 47, offset: 9840},
								expr: &seqExpr{
									pos: position{line: 394, col: 48, offset: 9841},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 394, col: 48, offset: 9841},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 394, col: 50, offset: 9843},
											val:        "?",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 394, col: 54, offset: 9847},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 394, col: 56, offset: 9849},
											name: "ExpressionTernery",
										},
										&ruleRefExpr{
											pos:  position{line: 394, col: 74, offset: 9867},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 394, col: 76, offset: 9869},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 394, col: 80, offset: 9873},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 394, col: 82, offset: 9875},
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
			pos:  position{line: 414, col: 1, offset: 10247},
			expr: &actionExpr{
				pos: position{line: 414, col: 20, offset: 10266},
				run: (*parser).callonExpressionBinOp1,
				expr: &seqExpr{
					pos: position{line: 414, col: 20, offset: 10266},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 414, col: 20, offset: 10266},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 414, col: 26, offset: 10272},
								name: "ExpressionCmpOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 414, col: 42, offset: 10288},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 414, col: 47, offset: 10293},
								expr: &seqExpr{
									pos: position{line: 414, col: 49, offset: 10295},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 414, col: 49, offset: 10295},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 414, col: 51, offset: 10297},
											name: "CmpOp",
										},
										&ruleRefExpr{
											pos:  position{line: 414, col: 57, offset: 10303},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 414, col: 59, offset: 10305},
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
			pos:  position{line: 418, col: 1, offset: 10365},
			expr: &actionExpr{
				pos: position{line: 418, col: 20, offset: 10384},
				run: (*parser).callonExpressionCmpOp1,
				expr: &seqExpr{
					pos: position{line: 418, col: 20, offset: 10384},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 418, col: 20, offset: 10384},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 418, col: 26, offset: 10390},
								name: "ExpressionAddOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 418, col: 42, offset: 10406},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 418, col: 47, offset: 10411},
								expr: &seqExpr{
									pos: position{line: 418, col: 49, offset: 10413},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 418, col: 49, offset: 10413},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 418, col: 51, offset: 10415},
											name: "CmpOp",
										},
										&ruleRefExpr{
											pos:  position{line: 418, col: 57, offset: 10421},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 418, col: 59, offset: 10423},
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
			pos:  position{line: 422, col: 1, offset: 10483},
			expr: &actionExpr{
				pos: position{line: 422, col: 20, offset: 10502},
				run: (*parser).callonExpressionAddOp1,
				expr: &seqExpr{
					pos: position{line: 422, col: 20, offset: 10502},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 422, col: 20, offset: 10502},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 422, col: 26, offset: 10508},
								name: "ExpressionMulOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 422, col: 42, offset: 10524},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 422, col: 47, offset: 10529},
								expr: &seqExpr{
									pos: position{line: 422, col: 49, offset: 10531},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 422, col: 49, offset: 10531},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 422, col: 51, offset: 10533},
											name: "AddOp",
										},
										&ruleRefExpr{
											pos:  position{line: 422, col: 57, offset: 10539},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 422, col: 59, offset: 10541},
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
			pos:  position{line: 426, col: 1, offset: 10601},
			expr: &actionExpr{
				pos: position{line: 426, col: 20, offset: 10620},
				run: (*parser).callonExpressionMulOp1,
				expr: &seqExpr{
					pos: position{line: 426, col: 20, offset: 10620},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 426, col: 20, offset: 10620},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 426, col: 26, offset: 10626},
								name: "ExpressionUnaryOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 426, col: 44, offset: 10644},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 426, col: 49, offset: 10649},
								expr: &seqExpr{
									pos: position{line: 426, col: 51, offset: 10651},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 426, col: 51, offset: 10651},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 426, col: 53, offset: 10653},
											name: "MulOp",
										},
										&ruleRefExpr{
											pos:  position{line: 426, col: 59, offset: 10659},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 426, col: 61, offset: 10661},
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
			pos:  position{line: 430, col: 1, offset: 10721},
			expr: &choiceExpr{
				pos: position{line: 430, col: 22, offset: 10742},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 430, col: 22, offset: 10742},
						run: (*parser).callonExpressionUnaryOp2,
						expr: &seqExpr{
							pos: position{line: 430, col: 22, offset: 10742},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 430, col: 22, offset: 10742},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 430, col: 25, offset: 10745},
										name: "UnaryOp",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 430, col: 33, offset: 10753},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 430, col: 35, offset: 10755},
									label: "ex",
									expr: &ruleRefExpr{
										pos:  position{line: 430, col: 38, offset: 10758},
										name: "ExpressionFactor",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 432, col: 5, offset: 10881},
						name: "ExpressionFactor",
					},
				},
			},
		},
		{
			name: "ExpressionFactor",
			pos:  position{line: 434, col: 1, offset: 10899},
			expr: &choiceExpr{
				pos: position{line: 434, col: 21, offset: 10919},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 434, col: 21, offset: 10919},
						run: (*parser).callonExpressionFactor2,
						expr: &seqExpr{
							pos: position{line: 434, col: 21, offset: 10919},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 434, col: 21, offset: 10919},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 434, col: 25, offset: 10923},
									label: "e",
									expr: &ruleRefExpr{
										pos:  position{line: 434, col: 27, offset: 10925},
										name: "Expression",
									},
								},
								&litMatcher{
									pos:        position{line: 434, col: 38, offset: 10936},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 436, col: 5, offset: 10962},
						name: "StringExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 436, col: 24, offset: 10981},
						name: "NumberExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 436, col: 43, offset: 11000},
						name: "BooleanExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 436, col: 63, offset: 11020},
						name: "NilExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 436, col: 79, offset: 11036},
						name: "MemberExpression",
					},
				},
			},
		},
		{
			name: "StringExpression",
			pos:  position{line: 438, col: 1, offset: 11054},
			expr: &actionExpr{
				pos: position{line: 438, col: 21, offset: 11074},
				run: (*parser).callonStringExpression1,
				expr: &labeledExpr{
					pos:   position{line: 438, col: 21, offset: 11074},
					label: "s",
					expr: &ruleRefExpr{
						pos:  position{line: 438, col: 23, offset: 11076},
						name: "String",
					},
				},
			},
		},
		{
			name: "NumberExpression",
			pos:  position{line: 442, col: 1, offset: 11171},
			expr: &actionExpr{
				pos: position{line: 442, col: 21, offset: 11191},
				run: (*parser).callonNumberExpression1,
				expr: &seqExpr{
					pos: position{line: 442, col: 21, offset: 11191},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 442, col: 21, offset: 11191},
							expr: &litMatcher{
								pos:        position{line: 442, col: 21, offset: 11191},
								val:        "-",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 442, col: 26, offset: 11196},
							name: "Integer",
						},
						&labeledExpr{
							pos:   position{line: 442, col: 34, offset: 11204},
							label: "dec",
							expr: &zeroOrOneExpr{
								pos: position{line: 442, col: 38, offset: 11208},
								expr: &seqExpr{
									pos: position{line: 442, col: 40, offset: 11210},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 442, col: 40, offset: 11210},
											val:        ".",
											ignoreCase: false,
										},
										&oneOrMoreExpr{
											pos: position{line: 442, col: 44, offset: 11214},
											expr: &ruleRefExpr{
												pos:  position{line: 442, col: 44, offset: 11214},
												name: "DecimalDigit",
											},
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 442, col: 61, offset: 11231},
							label: "ex",
							expr: &zeroOrOneExpr{
								pos: position{line: 442, col: 64, offset: 11234},
								expr: &ruleRefExpr{
									pos:  position{line: 442, col: 64, offset: 11234},
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
			pos:  position{line: 452, col: 1, offset: 11564},
			expr: &actionExpr{
				pos: position{line: 452, col: 18, offset: 11581},
				run: (*parser).callonNilExpression1,
				expr: &ruleRefExpr{
					pos:  position{line: 452, col: 18, offset: 11581},
					name: "Null",
				},
			},
		},
		{
			name: "BooleanExpression",
			pos:  position{line: 456, col: 1, offset: 11652},
			expr: &actionExpr{
				pos: position{line: 456, col: 22, offset: 11673},
				run: (*parser).callonBooleanExpression1,
				expr: &labeledExpr{
					pos:   position{line: 456, col: 22, offset: 11673},
					label: "b",
					expr: &ruleRefExpr{
						pos:  position{line: 456, col: 24, offset: 11675},
						name: "Bool",
					},
				},
			},
		},
		{
			name: "MemberExpression",
			pos:  position{line: 460, col: 1, offset: 11767},
			expr: &actionExpr{
				pos: position{line: 460, col: 21, offset: 11787},
				run: (*parser).callonMemberExpression1,
				expr: &seqExpr{
					pos: position{line: 460, col: 21, offset: 11787},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 460, col: 21, offset: 11787},
							label: "field",
							expr: &ruleRefExpr{
								pos:  position{line: 460, col: 27, offset: 11793},
								name: "Field",
							},
						},
						&labeledExpr{
							pos:   position{line: 460, col: 33, offset: 11799},
							label: "member",
							expr: &zeroOrMoreExpr{
								pos: position{line: 460, col: 40, offset: 11806},
								expr: &choiceExpr{
									pos: position{line: 460, col: 41, offset: 11807},
									alternatives: []interface{}{
										&seqExpr{
											pos: position{line: 460, col: 42, offset: 11808},
											exprs: []interface{}{
												&litMatcher{
													pos:        position{line: 460, col: 42, offset: 11808},
													val:        ".",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 460, col: 46, offset: 11812},
													name: "Identifier",
												},
											},
										},
										&seqExpr{
											pos: position{line: 460, col: 61, offset: 11827},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 460, col: 61, offset: 11827},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 460, col: 63, offset: 11829},
													name: "CallArguments",
												},
											},
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
			name: "Field",
			pos:  position{line: 477, col: 1, offset: 12287},
			expr: &actionExpr{
				pos: position{line: 477, col: 10, offset: 12296},
				run: (*parser).callonField1,
				expr: &labeledExpr{
					pos:   position{line: 477, col: 10, offset: 12296},
					label: "variable",
					expr: &ruleRefExpr{
						pos:  position{line: 477, col: 19, offset: 12305},
						name: "Variable",
					},
				},
			},
		},
		{
			name: "UnaryOp",
			pos:  position{line: 481, col: 1, offset: 12414},
			expr: &actionExpr{
				pos: position{line: 481, col: 12, offset: 12425},
				run: (*parser).callonUnaryOp1,
				expr: &choiceExpr{
					pos: position{line: 481, col: 14, offset: 12427},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 481, col: 14, offset: 12427},
							val:        "+",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 481, col: 20, offset: 12433},
							val:        "-",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 481, col: 26, offset: 12439},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AddOp",
			pos:  position{line: 485, col: 1, offset: 12479},
			expr: &actionExpr{
				pos: position{line: 485, col: 10, offset: 12488},
				run: (*parser).callonAddOp1,
				expr: &choiceExpr{
					pos: position{line: 485, col: 12, offset: 12490},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 485, col: 12, offset: 12490},
							val:        "+",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 485, col: 18, offset: 12496},
							val:        "-",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "MulOp",
			pos:  position{line: 489, col: 1, offset: 12536},
			expr: &actionExpr{
				pos: position{line: 489, col: 10, offset: 12545},
				run: (*parser).callonMulOp1,
				expr: &choiceExpr{
					pos: position{line: 489, col: 12, offset: 12547},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 489, col: 12, offset: 12547},
							val:        "*",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 489, col: 18, offset: 12553},
							val:        "/",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 489, col: 24, offset: 12559},
							val:        "%",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CmpOp",
			pos:  position{line: 493, col: 1, offset: 12599},
			expr: &actionExpr{
				pos: position{line: 493, col: 10, offset: 12608},
				run: (*parser).callonCmpOp1,
				expr: &choiceExpr{
					pos: position{line: 493, col: 12, offset: 12610},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 493, col: 12, offset: 12610},
							val:        "==",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 493, col: 19, offset: 12617},
							val:        "!=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 493, col: 26, offset: 12624},
							val:        "<",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 493, col: 32, offset: 12630},
							val:        ">",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 493, col: 38, offset: 12636},
							val:        "<=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 493, col: 45, offset: 12643},
							val:        ">=",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "BinOp",
			pos:  position{line: 497, col: 1, offset: 12684},
			expr: &actionExpr{
				pos: position{line: 497, col: 10, offset: 12693},
				run: (*parser).callonBinOp1,
				expr: &choiceExpr{
					pos: position{line: 497, col: 12, offset: 12695},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 497, col: 12, offset: 12695},
							val:        "&&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 497, col: 19, offset: 12702},
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
			pos:         position{line: 501, col: 1, offset: 12743},
			expr: &actionExpr{
				pos: position{line: 501, col: 20, offset: 12762},
				run: (*parser).callonString1,
				expr: &seqExpr{
					pos: position{line: 501, col: 20, offset: 12762},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 501, col: 20, offset: 12762},
							name: "Quote",
						},
						&zeroOrMoreExpr{
							pos: position{line: 501, col: 26, offset: 12768},
							expr: &choiceExpr{
								pos: position{line: 501, col: 28, offset: 12770},
								alternatives: []interface{}{
									&seqExpr{
										pos: position{line: 501, col: 28, offset: 12770},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 501, col: 28, offset: 12770},
												expr: &ruleRefExpr{
													pos:  position{line: 501, col: 29, offset: 12771},
													name: "EscapedChar",
												},
											},
											&anyMatcher{
												line: 501, col: 41, offset: 12783,
											},
										},
									},
									&seqExpr{
										pos: position{line: 501, col: 45, offset: 12787},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 501, col: 45, offset: 12787},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 501, col: 50, offset: 12792},
												name: "EscapeSequence",
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 501, col: 68, offset: 12810},
							name: "Quote",
						},
					},
				},
			},
		},
		{
			name:        "Quote",
			displayName: "\"quote\"",
			pos:         position{line: 505, col: 1, offset: 12862},
			expr: &litMatcher{
				pos:        position{line: 505, col: 18, offset: 12879},
				val:        "\"",
				ignoreCase: false,
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 507, col: 1, offset: 12884},
			expr: &charClassMatcher{
				pos:        position{line: 507, col: 16, offset: 12899},
				val:        "[\\x00-\\x1f\"\\\\]",
				chars:      []rune{'"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 508, col: 1, offset: 12914},
			expr: &choiceExpr{
				pos: position{line: 508, col: 19, offset: 12932},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 508, col: 19, offset: 12932},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 508, col: 38, offset: 12951},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 509, col: 1, offset: 12965},
			expr: &charClassMatcher{
				pos:        position{line: 509, col: 21, offset: 12985},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 510, col: 1, offset: 12997},
			expr: &seqExpr{
				pos: position{line: 510, col: 18, offset: 13014},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 510, col: 18, offset: 13014},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 510, col: 22, offset: 13018},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 510, col: 31, offset: 13027},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 510, col: 40, offset: 13036},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 510, col: 49, offset: 13045},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 512, col: 1, offset: 13055},
			expr: &choiceExpr{
				pos: position{line: 512, col: 12, offset: 13066},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 512, col: 12, offset: 13066},
						val:        "0",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 512, col: 18, offset: 13072},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 512, col: 18, offset: 13072},
								name: "NonZeroDecimalDigit",
							},
							&zeroOrMoreExpr{
								pos: position{line: 512, col: 38, offset: 13092},
								expr: &ruleRefExpr{
									pos:  position{line: 512, col: 38, offset: 13092},
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
			pos:  position{line: 513, col: 1, offset: 13106},
			expr: &seqExpr{
				pos: position{line: 513, col: 13, offset: 13118},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 513, col: 13, offset: 13118},
						val:        "e",
						ignoreCase: true,
					},
					&zeroOrOneExpr{
						pos: position{line: 513, col: 18, offset: 13123},
						expr: &charClassMatcher{
							pos:        position{line: 513, col: 18, offset: 13123},
							val:        "[+-]",
							chars:      []rune{'+', '-'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&oneOrMoreExpr{
						pos: position{line: 513, col: 24, offset: 13129},
						expr: &ruleRefExpr{
							pos:  position{line: 513, col: 24, offset: 13129},
							name: "DecimalDigit",
						},
					},
				},
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 514, col: 1, offset: 13143},
			expr: &charClassMatcher{
				pos:        position{line: 514, col: 17, offset: 13159},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "NonZeroDecimalDigit",
			pos:  position{line: 515, col: 1, offset: 13165},
			expr: &charClassMatcher{
				pos:        position{line: 515, col: 24, offset: 13188},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 516, col: 1, offset: 13194},
			expr: &charClassMatcher{
				pos:        position{line: 516, col: 13, offset: 13206},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "Bool",
			pos:  position{line: 517, col: 1, offset: 13216},
			expr: &choiceExpr{
				pos: position{line: 517, col: 9, offset: 13224},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 517, col: 9, offset: 13224},
						run: (*parser).callonBool2,
						expr: &litMatcher{
							pos:        position{line: 517, col: 9, offset: 13224},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 517, col: 39, offset: 13254},
						run: (*parser).callonBool4,
						expr: &litMatcher{
							pos:        position{line: 517, col: 39, offset: 13254},
							val:        "false",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Null",
			pos:  position{line: 518, col: 1, offset: 13284},
			expr: &actionExpr{
				pos: position{line: 518, col: 9, offset: 13292},
				run: (*parser).callonNull1,
				expr: &choiceExpr{
					pos: position{line: 518, col: 10, offset: 13293},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 518, col: 10, offset: 13293},
							val:        "null",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 518, col: 19, offset: 13302},
							val:        "nil",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Variable",
			pos:  position{line: 520, col: 1, offset: 13330},
			expr: &actionExpr{
				pos: position{line: 520, col: 13, offset: 13342},
				run: (*parser).callonVariable1,
				expr: &labeledExpr{
					pos:   position{line: 520, col: 13, offset: 13342},
					label: "ident",
					expr: &ruleRefExpr{
						pos:  position{line: 520, col: 19, offset: 13348},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 524, col: 1, offset: 13442},
			expr: &actionExpr{
				pos: position{line: 524, col: 15, offset: 13456},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 524, col: 15, offset: 13456},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 524, col: 15, offset: 13456},
							val:        "[a-zA-Z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 524, col: 25, offset: 13466},
							expr: &charClassMatcher{
								pos:        position{line: 524, col: 25, offset: 13466},
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
			pos:  position{line: 528, col: 1, offset: 13514},
			expr: &actionExpr{
				pos: position{line: 528, col: 14, offset: 13527},
				run: (*parser).callonEmptyLine1,
				expr: &seqExpr{
					pos: position{line: 528, col: 14, offset: 13527},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 528, col: 14, offset: 13527},
							name: "_",
						},
						&charClassMatcher{
							pos:        position{line: 528, col: 16, offset: 13529},
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
			pos:         position{line: 532, col: 1, offset: 13557},
			expr: &actionExpr{
				pos: position{line: 532, col: 19, offset: 13575},
				run: (*parser).callon_1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 532, col: 19, offset: 13575},
					expr: &charClassMatcher{
						pos:        position{line: 532, col: 19, offset: 13575},
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
			pos:         position{line: 533, col: 1, offset: 13602},
			expr: &actionExpr{
				pos: position{line: 533, col: 20, offset: 13621},
				run: (*parser).callon__1,
				expr: &charClassMatcher{
					pos:        position{line: 533, col: 20, offset: 13621},
					val:        "[ \\t]",
					chars:      []rune{' ', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "NL",
			pos:  position{line: 534, col: 1, offset: 13648},
			expr: &choiceExpr{
				pos: position{line: 534, col: 7, offset: 13654},
				alternatives: []interface{}{
					&charClassMatcher{
						pos:        position{line: 534, col: 7, offset: 13654},
						val:        "[\\n]",
						chars:      []rune{'\n'},
						ignoreCase: false,
						inverted:   false,
					},
					&andExpr{
						pos: position{line: 534, col: 14, offset: 13661},
						expr: &ruleRefExpr{
							pos:  position{line: 534, col: 15, offset: 13662},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 535, col: 1, offset: 13666},
			expr: &notExpr{
				pos: position{line: 535, col: 8, offset: 13673},
				expr: &anyMatcher{
					line: 535, col: 9, offset: 13674,
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

func (c *current) onListNode14() (interface{}, error) {
	return nil, nil
}

func (p *parser) callonListNode14() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListNode14()
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
		tagElem.Block = list.(*List)
	}

	return tagElem, nil
}

func (p *parser) callonTag1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTag1(stack["tag"], stack["list"])
}

func (c *current) onTagHeader2(name, attrs, tl interface{}) (interface{}, error) {
	tag := &Tag{Name: name.(string), GraphNode: NewNode(pos(c.pos))}
	if attrs != nil {
		tag.Attributes = attrs.([]*Attribute)
	}
	if tl != nil {
		tag.Text = toSlice(tl)[1].(*TextList)
	}
	return tag, nil
}

func (p *parser) callonTagHeader2() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagHeader2(stack["name"], stack["attrs"], stack["tl"])
}

func (c *current) onTagHeader17(name, attrs, text interface{}) (interface{}, error) {
	tag := &Tag{Name: name.(string), GraphNode: NewNode(pos(c.pos))}
	if attrs != nil {
		tag.Attributes = attrs.([]*Attribute)
	}
	if text != nil {
		tag.Block = text.(*Text)
	}
	return tag, nil
}

func (p *parser) callonTagHeader17() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagHeader17(stack["name"], stack["attrs"], stack["text"])
}

func (c *current) onTagHeader30(name, attrs, unescaped, expr interface{}) (interface{}, error) {
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

func (p *parser) callonTagHeader30() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagHeader30(stack["name"], stack["attrs"], stack["unescaped"], stack["expr"])
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
		mSlice := toSlice(m)

		if mSlice[0] != nil && string(mSlice[0].([]byte)) == "." {
			cur = &MemberExpression{X: cur, Name: mSlice[1].(string), GraphNode: NewNode(pos(c.pos))}
		} else {
			cur = &FunctionCallExpression{X: cur, Arguments: mSlice[1].([]Expression), GraphNode: NewNode(pos(c.pos))}
		}
	}

	return cur, nil
}

func (p *parser) callonMemberExpression1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMemberExpression1(stack["field"], stack["member"])
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
