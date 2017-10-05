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
			pos:  position{line: 41, col: 1, offset: 770},
			expr: &choiceExpr{
				pos: position{line: 41, col: 9, offset: 778},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 41, col: 9, offset: 778},
						run: (*parser).callonList2,
						expr: &seqExpr{
							pos: position{line: 41, col: 9, offset: 778},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 41, col: 9, offset: 778},
									label: "node",
									expr: &ruleRefExpr{
										pos:  position{line: 41, col: 14, offset: 783},
										name: "ListNode",
									},
								},
								&labeledExpr{
									pos:   position{line: 41, col: 23, offset: 792},
									label: "list",
									expr: &ruleRefExpr{
										pos:  position{line: 41, col: 28, offset: 797},
										name: "List",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 49, col: 5, offset: 983},
						run: (*parser).callonList8,
						expr: &andExpr{
							pos: position{line: 49, col: 5, offset: 983},
							expr: &choiceExpr{
								pos: position{line: 49, col: 7, offset: 985},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 49, col: 7, offset: 985},
										name: "Outdent",
									},
									&ruleRefExpr{
										pos:  position{line: 49, col: 17, offset: 995},
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
			pos:  position{line: 53, col: 1, offset: 1057},
			expr: &actionExpr{
				pos: position{line: 53, col: 17, offset: 1073},
				run: (*parser).callonIndentedList1,
				expr: &seqExpr{
					pos: position{line: 53, col: 17, offset: 1073},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 53, col: 17, offset: 1073},
							name: "Indent",
						},
						&labeledExpr{
							pos:   position{line: 53, col: 24, offset: 1080},
							label: "list",
							expr: &ruleRefExpr{
								pos:  position{line: 53, col: 29, offset: 1085},
								name: "List",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 53, col: 34, offset: 1090},
							name: "Outdent",
						},
					},
				},
			},
		},
		{
			name: "IndentedRawText",
			pos:  position{line: 57, col: 1, offset: 1122},
			expr: &actionExpr{
				pos: position{line: 57, col: 20, offset: 1141},
				run: (*parser).callonIndentedRawText1,
				expr: &seqExpr{
					pos: position{line: 57, col: 20, offset: 1141},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 57, col: 20, offset: 1141},
							name: "Indent",
						},
						&labeledExpr{
							pos:   position{line: 57, col: 27, offset: 1148},
							label: "t",
							expr: &ruleRefExpr{
								pos:  position{line: 57, col: 29, offset: 1150},
								name: "RawText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 57, col: 37, offset: 1158},
							name: "Outdent",
						},
					},
				},
			},
		},
		{
			name: "RawText",
			pos:  position{line: 61, col: 1, offset: 1242},
			expr: &choiceExpr{
				pos: position{line: 61, col: 12, offset: 1253},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 61, col: 12, offset: 1253},
						run: (*parser).callonRawText2,
						expr: &seqExpr{
							pos: position{line: 61, col: 12, offset: 1253},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 61, col: 12, offset: 1253},
									name: "Indent",
								},
								&labeledExpr{
									pos:   position{line: 61, col: 19, offset: 1260},
									label: "rt",
									expr: &ruleRefExpr{
										pos:  position{line: 61, col: 22, offset: 1263},
										name: "RawText",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 61, col: 30, offset: 1271},
									name: "Outdent",
								},
								&labeledExpr{
									pos:   position{line: 61, col: 38, offset: 1279},
									label: "tail",
									expr: &ruleRefExpr{
										pos:  position{line: 61, col: 43, offset: 1284},
										name: "RawText",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 63, col: 5, offset: 1340},
						run: (*parser).callonRawText10,
						expr: &choiceExpr{
							pos: position{line: 63, col: 6, offset: 1341},
							alternatives: []interface{}{
								&andExpr{
									pos: position{line: 63, col: 6, offset: 1341},
									expr: &ruleRefExpr{
										pos:  position{line: 63, col: 7, offset: 1342},
										name: "Outdent",
									},
								},
								&andExpr{
									pos: position{line: 63, col: 17, offset: 1352},
									expr: &ruleRefExpr{
										pos:  position{line: 63, col: 18, offset: 1353},
										name: "EOF",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 65, col: 5, offset: 1381},
						run: (*parser).callonRawText16,
						expr: &seqExpr{
							pos: position{line: 65, col: 5, offset: 1381},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 65, col: 5, offset: 1381},
									label: "head",
									expr: &anyMatcher{
										line: 65, col: 10, offset: 1386,
									},
								},
								&labeledExpr{
									pos:   position{line: 65, col: 12, offset: 1388},
									label: "tail",
									expr: &ruleRefExpr{
										pos:  position{line: 65, col: 17, offset: 1393},
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
			pos:         position{line: 69, col: 1, offset: 1458},
			expr: &choiceExpr{
				pos: position{line: 70, col: 3, offset: 1483},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 70, col: 3, offset: 1483},
						name: "Comment",
					},
					&ruleRefExpr{
						pos:  position{line: 71, col: 3, offset: 1495},
						name: "Import",
					},
					&ruleRefExpr{
						pos:  position{line: 72, col: 3, offset: 1506},
						name: "Extend",
					},
					&ruleRefExpr{
						pos:  position{line: 73, col: 3, offset: 1517},
						name: "PipeText",
					},
					&ruleRefExpr{
						pos:  position{line: 74, col: 3, offset: 1530},
						name: "If",
					},
					&ruleRefExpr{
						pos:  position{line: 75, col: 3, offset: 1537},
						name: "Each",
					},
					&ruleRefExpr{
						pos:  position{line: 76, col: 3, offset: 1546},
						name: "DocType",
					},
					&ruleRefExpr{
						pos:  position{line: 77, col: 3, offset: 1558},
						name: "Mixin",
					},
					&ruleRefExpr{
						pos:  position{line: 78, col: 3, offset: 1568},
						name: "MixinCall",
					},
					&ruleRefExpr{
						pos:  position{line: 79, col: 3, offset: 1582},
						name: "Assignment",
					},
					&ruleRefExpr{
						pos:  position{line: 80, col: 3, offset: 1597},
						name: "Block",
					},
					&ruleRefExpr{
						pos:  position{line: 81, col: 3, offset: 1607},
						name: "Tag",
					},
					&actionExpr{
						pos: position{line: 82, col: 3, offset: 1615},
						run: (*parser).callonListNode14,
						expr: &seqExpr{
							pos: position{line: 82, col: 4, offset: 1616},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 82, col: 4, offset: 1616},
									name: "_",
								},
								&charClassMatcher{
									pos:        position{line: 82, col: 6, offset: 1618},
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
			pos:         position{line: 85, col: 1, offset: 1653},
			expr: &actionExpr{
				pos: position{line: 85, col: 21, offset: 1673},
				run: (*parser).callonDocType1,
				expr: &seqExpr{
					pos: position{line: 85, col: 21, offset: 1673},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 85, col: 21, offset: 1673},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 85, col: 23, offset: 1675},
							val:        "doctype",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 85, col: 33, offset: 1685},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 85, col: 35, offset: 1687},
							label: "val",
							expr: &ruleRefExpr{
								pos:  position{line: 85, col: 39, offset: 1691},
								name: "LineText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 85, col: 48, offset: 1700},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Tag",
			pos:  position{line: 91, col: 1, offset: 1793},
			expr: &actionExpr{
				pos: position{line: 91, col: 8, offset: 1800},
				run: (*parser).callonTag1,
				expr: &seqExpr{
					pos: position{line: 91, col: 8, offset: 1800},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 91, col: 8, offset: 1800},
							label: "tag",
							expr: &ruleRefExpr{
								pos:  position{line: 91, col: 12, offset: 1804},
								name: "TagHeader",
							},
						},
						&labeledExpr{
							pos:   position{line: 91, col: 22, offset: 1814},
							label: "list",
							expr: &zeroOrOneExpr{
								pos: position{line: 91, col: 27, offset: 1819},
								expr: &ruleRefExpr{
									pos:  position{line: 91, col: 27, offset: 1819},
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
			pos:  position{line: 101, col: 1, offset: 1942},
			expr: &choiceExpr{
				pos: position{line: 101, col: 14, offset: 1955},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 101, col: 14, offset: 1955},
						run: (*parser).callonTagHeader2,
						expr: &seqExpr{
							pos: position{line: 101, col: 14, offset: 1955},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 101, col: 14, offset: 1955},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 101, col: 16, offset: 1957},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 101, col: 21, offset: 1962},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 101, col: 29, offset: 1970},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 101, col: 35, offset: 1976},
										expr: &ruleRefExpr{
											pos:  position{line: 101, col: 35, offset: 1976},
											name: "TagAttributes",
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 101, col: 50, offset: 1991},
									label: "tl",
									expr: &zeroOrOneExpr{
										pos: position{line: 101, col: 53, offset: 1994},
										expr: &seqExpr{
											pos: position{line: 101, col: 54, offset: 1995},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 101, col: 54, offset: 1995},
													name: "__",
												},
												&zeroOrOneExpr{
													pos: position{line: 101, col: 57, offset: 1998},
													expr: &ruleRefExpr{
														pos:  position{line: 101, col: 57, offset: 1998},
														name: "TextList",
													},
												},
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 101, col: 69, offset: 2010},
									name: "NL",
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 110, col: 5, offset: 2233},
						run: (*parser).callonTagHeader17,
						expr: &seqExpr{
							pos: position{line: 110, col: 5, offset: 2233},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 110, col: 5, offset: 2233},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 110, col: 7, offset: 2235},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 110, col: 12, offset: 2240},
										name: "TagName",
									},
								},
								&labeledExpr{
									pos:   position{line: 110, col: 20, offset: 2248},
									label: "attrs",
									expr: &zeroOrOneExpr{
										pos: position{line: 110, col: 26, offset: 2254},
										expr: &ruleRefExpr{
											pos:  position{line: 110, col: 26, offset: 2254},
											name: "TagAttributes",
										},
									},
								},
								&litMatcher{
									pos:        position{line: 110, col: 41, offset: 2269},
									val:        ".",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 110, col: 45, offset: 2273},
									name: "NL",
								},
								&labeledExpr{
									pos:   position{line: 110, col: 48, offset: 2276},
									label: "text",
									expr: &zeroOrOneExpr{
										pos: position{line: 110, col: 53, offset: 2281},
										expr: &ruleRefExpr{
											pos:  position{line: 110, col: 53, offset: 2281},
											name: "IndentedRawText",
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
			name: "TagName",
			pos:  position{line: 121, col: 1, offset: 2506},
			expr: &actionExpr{
				pos: position{line: 121, col: 12, offset: 2517},
				run: (*parser).callonTagName1,
				expr: &seqExpr{
					pos: position{line: 121, col: 12, offset: 2517},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 121, col: 12, offset: 2517},
							val:        "[_a-zA-Z]",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 121, col: 22, offset: 2527},
							expr: &charClassMatcher{
								pos:        position{line: 121, col: 22, offset: 2527},
								val:        "[_-:a-zA-Z0-9]",
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
			name: "TagAttributes",
			pos:  position{line: 125, col: 1, offset: 2577},
			expr: &choiceExpr{
				pos: position{line: 125, col: 18, offset: 2594},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 125, col: 18, offset: 2594},
						run: (*parser).callonTagAttributes2,
						expr: &seqExpr{
							pos: position{line: 125, col: 18, offset: 2594},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 125, col: 18, offset: 2594},
									label: "head",
									expr: &choiceExpr{
										pos: position{line: 125, col: 24, offset: 2600},
										alternatives: []interface{}{
											&ruleRefExpr{
												pos:  position{line: 125, col: 24, offset: 2600},
												name: "TagAttributeClass",
											},
											&ruleRefExpr{
												pos:  position{line: 125, col: 44, offset: 2620},
												name: "TagAttributeID",
											},
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 125, col: 60, offset: 2636},
									label: "tail",
									expr: &zeroOrOneExpr{
										pos: position{line: 125, col: 65, offset: 2641},
										expr: &ruleRefExpr{
											pos:  position{line: 125, col: 65, offset: 2641},
											name: "TagAttributes",
										},
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 133, col: 5, offset: 2806},
						run: (*parser).callonTagAttributes11,
						expr: &seqExpr{
							pos: position{line: 133, col: 5, offset: 2806},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 133, col: 5, offset: 2806},
									label: "head",
									expr: &seqExpr{
										pos: position{line: 133, col: 11, offset: 2812},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 133, col: 11, offset: 2812},
												val:        "(",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 133, col: 15, offset: 2816},
												name: "_",
											},
											&seqExpr{
												pos: position{line: 133, col: 18, offset: 2819},
												exprs: []interface{}{
													&ruleRefExpr{
														pos:  position{line: 133, col: 18, offset: 2819},
														name: "TagAttribute",
													},
													&zeroOrMoreExpr{
														pos: position{line: 133, col: 31, offset: 2832},
														expr: &seqExpr{
															pos: position{line: 133, col: 32, offset: 2833},
															exprs: []interface{}{
																&ruleRefExpr{
																	pos:  position{line: 133, col: 32, offset: 2833},
																	name: "__",
																},
																&ruleRefExpr{
																	pos:  position{line: 133, col: 35, offset: 2836},
																	name: "TagAttribute",
																},
															},
														},
													},
												},
											},
											&ruleRefExpr{
												pos:  position{line: 133, col: 51, offset: 2852},
												name: "_",
											},
											&litMatcher{
												pos:        position{line: 133, col: 53, offset: 2854},
												val:        ")",
												ignoreCase: false,
											},
										},
									},
								},
								&labeledExpr{
									pos:   position{line: 133, col: 58, offset: 2859},
									label: "tail",
									expr: &zeroOrOneExpr{
										pos: position{line: 133, col: 63, offset: 2864},
										expr: &ruleRefExpr{
											pos:  position{line: 133, col: 63, offset: 2864},
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
			pos:  position{line: 157, col: 1, offset: 3310},
			expr: &actionExpr{
				pos: position{line: 157, col: 22, offset: 3331},
				run: (*parser).callonTagAttributeClass1,
				expr: &seqExpr{
					pos: position{line: 157, col: 22, offset: 3331},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 157, col: 22, offset: 3331},
							val:        ".",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 157, col: 26, offset: 3335},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 157, col: 31, offset: 3340},
								name: "ClassName",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributeID",
			pos:  position{line: 161, col: 1, offset: 3489},
			expr: &actionExpr{
				pos: position{line: 161, col: 19, offset: 3507},
				run: (*parser).callonTagAttributeID1,
				expr: &seqExpr{
					pos: position{line: 161, col: 19, offset: 3507},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 161, col: 19, offset: 3507},
							val:        "#",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 161, col: 23, offset: 3511},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 161, col: 28, offset: 3516},
								name: "TagAttributeNameLiteral",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttribute",
			pos:  position{line: 165, col: 1, offset: 3676},
			expr: &choiceExpr{
				pos: position{line: 165, col: 17, offset: 3692},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 165, col: 17, offset: 3692},
						run: (*parser).callonTagAttribute2,
						expr: &seqExpr{
							pos: position{line: 165, col: 17, offset: 3692},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 165, col: 17, offset: 3692},
									label: "name",
									expr: &ruleRefExpr{
										pos:  position{line: 165, col: 22, offset: 3697},
										name: "TagAttributeName",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 165, col: 39, offset: 3714},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 165, col: 41, offset: 3716},
									val:        "=",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 165, col: 45, offset: 3720},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 165, col: 47, offset: 3722},
									label: "value",
									expr: &ruleRefExpr{
										pos:  position{line: 165, col: 53, offset: 3728},
										name: "Expression",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 167, col: 5, offset: 3864},
						run: (*parser).callonTagAttribute11,
						expr: &labeledExpr{
							pos:   position{line: 167, col: 5, offset: 3864},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 167, col: 10, offset: 3869},
								name: "TagAttributeName",
							},
						},
					},
				},
			},
		},
		{
			name: "TagAttributeName",
			pos:  position{line: 171, col: 1, offset: 3983},
			expr: &choiceExpr{
				pos: position{line: 171, col: 21, offset: 4003},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 171, col: 21, offset: 4003},
						run: (*parser).callonTagAttributeName2,
						expr: &seqExpr{
							pos: position{line: 171, col: 21, offset: 4003},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 171, col: 21, offset: 4003},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 171, col: 25, offset: 4007},
									label: "tn",
									expr: &ruleRefExpr{
										pos:  position{line: 171, col: 28, offset: 4010},
										name: "TagAttributeNameLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 171, col: 52, offset: 4034},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 173, col: 5, offset: 4061},
						run: (*parser).callonTagAttributeName8,
						expr: &seqExpr{
							pos: position{line: 173, col: 5, offset: 4061},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 173, col: 5, offset: 4061},
									val:        "[",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 173, col: 9, offset: 4065},
									label: "tn",
									expr: &ruleRefExpr{
										pos:  position{line: 173, col: 12, offset: 4068},
										name: "TagAttributeNameLiteral",
									},
								},
								&litMatcher{
									pos:        position{line: 173, col: 36, offset: 4092},
									val:        "]",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 175, col: 5, offset: 4119},
						run: (*parser).callonTagAttributeName14,
						expr: &labeledExpr{
							pos:   position{line: 175, col: 5, offset: 4119},
							label: "tn",
							expr: &ruleRefExpr{
								pos:  position{line: 175, col: 8, offset: 4122},
								name: "TagAttributeNameLiteral",
							},
						},
					},
				},
			},
		},
		{
			name: "ClassName",
			pos:  position{line: 179, col: 1, offset: 4168},
			expr: &actionExpr{
				pos: position{line: 179, col: 14, offset: 4181},
				run: (*parser).callonClassName1,
				expr: &oneOrMoreExpr{
					pos: position{line: 179, col: 14, offset: 4181},
					expr: &charClassMatcher{
						pos:        position{line: 179, col: 14, offset: 4181},
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
			pos:  position{line: 183, col: 1, offset: 4231},
			expr: &actionExpr{
				pos: position{line: 183, col: 28, offset: 4258},
				run: (*parser).callonTagAttributeNameLiteral1,
				expr: &seqExpr{
					pos: position{line: 183, col: 28, offset: 4258},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 183, col: 28, offset: 4258},
							val:        "[@_a-zA-Z]",
							chars:      []rune{'@', '_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 183, col: 39, offset: 4269},
							expr: &charClassMatcher{
								pos:        position{line: 183, col: 39, offset: 4269},
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
			pos:  position{line: 188, col: 1, offset: 4329},
			expr: &actionExpr{
				pos: position{line: 188, col: 7, offset: 4335},
				run: (*parser).callonIf1,
				expr: &seqExpr{
					pos: position{line: 188, col: 7, offset: 4335},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 188, col: 7, offset: 4335},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 188, col: 9, offset: 4337},
							val:        "if",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 14, offset: 4342},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 188, col: 17, offset: 4345},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 188, col: 22, offset: 4350},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 33, offset: 4361},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 35, offset: 4363},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 188, col: 38, offset: 4366},
							label: "block",
							expr: &ruleRefExpr{
								pos:  position{line: 188, col: 44, offset: 4372},
								name: "IndentedList",
							},
						},
						&labeledExpr{
							pos:   position{line: 188, col: 57, offset: 4385},
							label: "elseNode",
							expr: &zeroOrOneExpr{
								pos: position{line: 188, col: 66, offset: 4394},
								expr: &ruleRefExpr{
									pos:  position{line: 188, col: 66, offset: 4394},
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
			pos:  position{line: 196, col: 1, offset: 4603},
			expr: &choiceExpr{
				pos: position{line: 196, col: 9, offset: 4611},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 196, col: 9, offset: 4611},
						run: (*parser).callonElse2,
						expr: &seqExpr{
							pos: position{line: 196, col: 9, offset: 4611},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 196, col: 9, offset: 4611},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 196, col: 11, offset: 4613},
									val:        "else",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 196, col: 18, offset: 4620},
									label: "node",
									expr: &ruleRefExpr{
										pos:  position{line: 196, col: 23, offset: 4625},
										name: "If",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 198, col: 5, offset: 4653},
						run: (*parser).callonElse8,
						expr: &seqExpr{
							pos: position{line: 198, col: 5, offset: 4653},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 198, col: 5, offset: 4653},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 198, col: 7, offset: 4655},
									val:        "else",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 198, col: 14, offset: 4662},
									name: "_",
								},
								&ruleRefExpr{
									pos:  position{line: 198, col: 16, offset: 4664},
									name: "NL",
								},
								&labeledExpr{
									pos:   position{line: 198, col: 19, offset: 4667},
									label: "block",
									expr: &ruleRefExpr{
										pos:  position{line: 198, col: 25, offset: 4673},
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
			pos:  position{line: 202, col: 1, offset: 4711},
			expr: &actionExpr{
				pos: position{line: 202, col: 9, offset: 4719},
				run: (*parser).callonEach1,
				expr: &seqExpr{
					pos: position{line: 202, col: 9, offset: 4719},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 202, col: 9, offset: 4719},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 202, col: 11, offset: 4721},
							val:        "each",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 18, offset: 4728},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 202, col: 21, offset: 4731},
							label: "v1",
							expr: &ruleRefExpr{
								pos:  position{line: 202, col: 24, offset: 4734},
								name: "Variable",
							},
						},
						&labeledExpr{
							pos:   position{line: 202, col: 33, offset: 4743},
							label: "v2",
							expr: &zeroOrOneExpr{
								pos: position{line: 202, col: 36, offset: 4746},
								expr: &seqExpr{
									pos: position{line: 202, col: 37, offset: 4747},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 202, col: 37, offset: 4747},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 202, col: 39, offset: 4749},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 202, col: 43, offset: 4753},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 202, col: 45, offset: 4755},
											name: "Variable",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 56, offset: 4766},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 202, col: 58, offset: 4768},
							val:        "in",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 63, offset: 4773},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 202, col: 65, offset: 4775},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 202, col: 70, offset: 4780},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 81, offset: 4791},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 202, col: 83, offset: 4793},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 202, col: 86, offset: 4796},
							label: "block",
							expr: &ruleRefExpr{
								pos:  position{line: 202, col: 92, offset: 4802},
								name: "IndentedList",
							},
						},
					},
				},
			},
		},
		{
			name: "Assignment",
			pos:  position{line: 213, col: 1, offset: 5087},
			expr: &actionExpr{
				pos: position{line: 213, col: 15, offset: 5101},
				run: (*parser).callonAssignment1,
				expr: &seqExpr{
					pos: position{line: 213, col: 15, offset: 5101},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 213, col: 15, offset: 5101},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 213, col: 17, offset: 5103},
							label: "vr",
							expr: &ruleRefExpr{
								pos:  position{line: 213, col: 20, offset: 5106},
								name: "Variable",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 213, col: 29, offset: 5115},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 213, col: 31, offset: 5117},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 213, col: 35, offset: 5121},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 213, col: 37, offset: 5123},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 213, col: 42, offset: 5128},
								name: "Expression",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 213, col: 53, offset: 5139},
							name: "_",
						},
						&ruleRefExpr{
							pos:  position{line: 213, col: 55, offset: 5141},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Mixin",
			pos:  position{line: 218, col: 1, offset: 5273},
			expr: &actionExpr{
				pos: position{line: 218, col: 10, offset: 5282},
				run: (*parser).callonMixin1,
				expr: &seqExpr{
					pos: position{line: 218, col: 10, offset: 5282},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 218, col: 10, offset: 5282},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 218, col: 12, offset: 5284},
							val:        "mixin",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 218, col: 20, offset: 5292},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 218, col: 23, offset: 5295},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 218, col: 28, offset: 5300},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 218, col: 39, offset: 5311},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 218, col: 41, offset: 5313},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 218, col: 46, offset: 5318},
								expr: &ruleRefExpr{
									pos:  position{line: 218, col: 46, offset: 5318},
									name: "MixinArguments",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 218, col: 62, offset: 5334},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 218, col: 65, offset: 5337},
							label: "list",
							expr: &ruleRefExpr{
								pos:  position{line: 218, col: 70, offset: 5342},
								name: "IndentedList",
							},
						},
					},
				},
			},
		},
		{
			name: "MixinArguments",
			pos:  position{line: 226, col: 1, offset: 5551},
			expr: &choiceExpr{
				pos: position{line: 226, col: 19, offset: 5569},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 226, col: 19, offset: 5569},
						run: (*parser).callonMixinArguments2,
						expr: &seqExpr{
							pos: position{line: 226, col: 19, offset: 5569},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 226, col: 19, offset: 5569},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 226, col: 23, offset: 5573},
									label: "head",
									expr: &ruleRefExpr{
										pos:  position{line: 226, col: 28, offset: 5578},
										name: "MixinArgument",
									},
								},
								&labeledExpr{
									pos:   position{line: 226, col: 42, offset: 5592},
									label: "tail",
									expr: &zeroOrMoreExpr{
										pos: position{line: 226, col: 47, offset: 5597},
										expr: &seqExpr{
											pos: position{line: 226, col: 48, offset: 5598},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 226, col: 48, offset: 5598},
													name: "_",
												},
												&litMatcher{
													pos:        position{line: 226, col: 50, offset: 5600},
													val:        ",",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 226, col: 54, offset: 5604},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 226, col: 56, offset: 5606},
													name: "MixinArgument",
												},
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 226, col: 72, offset: 5622},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 239, col: 5, offset: 5884},
						run: (*parser).callonMixinArguments15,
						expr: &seqExpr{
							pos: position{line: 239, col: 5, offset: 5884},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 239, col: 5, offset: 5884},
									val:        "(",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 239, col: 9, offset: 5888},
									name: "_",
								},
								&litMatcher{
									pos:        position{line: 239, col: 11, offset: 5890},
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
			pos:  position{line: 243, col: 1, offset: 5917},
			expr: &actionExpr{
				pos: position{line: 243, col: 18, offset: 5934},
				run: (*parser).callonMixinArgument1,
				expr: &seqExpr{
					pos: position{line: 243, col: 18, offset: 5934},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 243, col: 18, offset: 5934},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 243, col: 23, offset: 5939},
								name: "Variable",
							},
						},
						&labeledExpr{
							pos:   position{line: 243, col: 32, offset: 5948},
							label: "def",
							expr: &zeroOrOneExpr{
								pos: position{line: 243, col: 36, offset: 5952},
								expr: &seqExpr{
									pos: position{line: 243, col: 37, offset: 5953},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 243, col: 37, offset: 5953},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 243, col: 39, offset: 5955},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 243, col: 43, offset: 5959},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 243, col: 45, offset: 5961},
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
			pos:  position{line: 254, col: 1, offset: 6185},
			expr: &actionExpr{
				pos: position{line: 254, col: 14, offset: 6198},
				run: (*parser).callonMixinCall1,
				expr: &seqExpr{
					pos: position{line: 254, col: 14, offset: 6198},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 254, col: 14, offset: 6198},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 254, col: 16, offset: 6200},
							val:        "+",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 254, col: 20, offset: 6204},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 254, col: 25, offset: 6209},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 254, col: 36, offset: 6220},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 254, col: 38, offset: 6222},
							label: "args",
							expr: &zeroOrOneExpr{
								pos: position{line: 254, col: 43, offset: 6227},
								expr: &ruleRefExpr{
									pos:  position{line: 254, col: 43, offset: 6227},
									name: "CallArguments",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 254, col: 58, offset: 6242},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "CallArguments",
			pos:  position{line: 262, col: 1, offset: 6413},
			expr: &actionExpr{
				pos: position{line: 262, col: 18, offset: 6430},
				run: (*parser).callonCallArguments1,
				expr: &seqExpr{
					pos: position{line: 262, col: 18, offset: 6430},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 262, col: 18, offset: 6430},
							val:        "(",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 262, col: 22, offset: 6434},
							label: "head",
							expr: &zeroOrOneExpr{
								pos: position{line: 262, col: 27, offset: 6439},
								expr: &ruleRefExpr{
									pos:  position{line: 262, col: 27, offset: 6439},
									name: "Expression",
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 262, col: 39, offset: 6451},
							label: "tail",
							expr: &zeroOrMoreExpr{
								pos: position{line: 262, col: 44, offset: 6456},
								expr: &seqExpr{
									pos: position{line: 262, col: 45, offset: 6457},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 262, col: 45, offset: 6457},
											name: "_",
										},
										&litMatcher{
											pos:        position{line: 262, col: 47, offset: 6459},
											val:        ",",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 262, col: 51, offset: 6463},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 262, col: 53, offset: 6465},
											name: "Expression",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 262, col: 66, offset: 6478},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Import",
			pos:  position{line: 283, col: 1, offset: 6800},
			expr: &actionExpr{
				pos: position{line: 283, col: 11, offset: 6810},
				run: (*parser).callonImport1,
				expr: &seqExpr{
					pos: position{line: 283, col: 11, offset: 6810},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 283, col: 11, offset: 6810},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 283, col: 13, offset: 6812},
							val:        "include",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 283, col: 23, offset: 6822},
							expr: &litMatcher{
								pos:        position{line: 283, col: 23, offset: 6822},
								val:        "s",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 283, col: 28, offset: 6827},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 283, col: 31, offset: 6830},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 283, col: 36, offset: 6835},
								name: "String",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 283, col: 43, offset: 6842},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "Extend",
			pos:  position{line: 287, col: 1, offset: 6925},
			expr: &actionExpr{
				pos: position{line: 287, col: 11, offset: 6935},
				run: (*parser).callonExtend1,
				expr: &seqExpr{
					pos: position{line: 287, col: 11, offset: 6935},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 287, col: 11, offset: 6935},
							val:        "extend",
							ignoreCase: false,
						},
						&zeroOrOneExpr{
							pos: position{line: 287, col: 20, offset: 6944},
							expr: &litMatcher{
								pos:        position{line: 287, col: 20, offset: 6944},
								val:        "s",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 287, col: 25, offset: 6949},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 287, col: 28, offset: 6952},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 287, col: 33, offset: 6957},
								name: "String",
							},
						},
					},
				},
			},
		},
		{
			name: "Block",
			pos:  position{line: 291, col: 1, offset: 7044},
			expr: &actionExpr{
				pos: position{line: 291, col: 10, offset: 7053},
				run: (*parser).callonBlock1,
				expr: &seqExpr{
					pos: position{line: 291, col: 10, offset: 7053},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 291, col: 10, offset: 7053},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 291, col: 12, offset: 7055},
							val:        "block",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 291, col: 20, offset: 7063},
							label: "mod",
							expr: &zeroOrOneExpr{
								pos: position{line: 291, col: 24, offset: 7067},
								expr: &seqExpr{
									pos: position{line: 291, col: 25, offset: 7068},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 291, col: 25, offset: 7068},
											name: "__",
										},
										&choiceExpr{
											pos: position{line: 291, col: 29, offset: 7072},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 291, col: 29, offset: 7072},
													val:        "append",
													ignoreCase: false,
												},
												&litMatcher{
													pos:        position{line: 291, col: 40, offset: 7083},
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
							pos:  position{line: 291, col: 53, offset: 7096},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 291, col: 56, offset: 7099},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 291, col: 61, offset: 7104},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 291, col: 72, offset: 7115},
							name: "NL",
						},
						&labeledExpr{
							pos:   position{line: 291, col: 75, offset: 7118},
							label: "list",
							expr: &zeroOrOneExpr{
								pos: position{line: 291, col: 80, offset: 7123},
								expr: &ruleRefExpr{
									pos:  position{line: 291, col: 80, offset: 7123},
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
			pos:  position{line: 312, col: 1, offset: 7488},
			expr: &actionExpr{
				pos: position{line: 312, col: 12, offset: 7499},
				run: (*parser).callonComment1,
				expr: &seqExpr{
					pos: position{line: 312, col: 12, offset: 7499},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 312, col: 12, offset: 7499},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 312, col: 14, offset: 7501},
							val:        "//",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 312, col: 19, offset: 7506},
							label: "silent",
							expr: &zeroOrOneExpr{
								pos: position{line: 312, col: 26, offset: 7513},
								expr: &litMatcher{
									pos:        position{line: 312, col: 26, offset: 7513},
									val:        "-",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 312, col: 31, offset: 7518},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 312, col: 33, offset: 7520},
							label: "comment",
							expr: &ruleRefExpr{
								pos:  position{line: 312, col: 41, offset: 7528},
								name: "LineText",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 312, col: 50, offset: 7537},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "LineText",
			pos:  position{line: 317, col: 1, offset: 7671},
			expr: &actionExpr{
				pos: position{line: 317, col: 13, offset: 7683},
				run: (*parser).callonLineText1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 317, col: 13, offset: 7683},
					expr: &charClassMatcher{
						pos:        position{line: 317, col: 13, offset: 7683},
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
			pos:  position{line: 322, col: 1, offset: 7732},
			expr: &actionExpr{
				pos: position{line: 322, col: 13, offset: 7744},
				run: (*parser).callonPipeText1,
				expr: &seqExpr{
					pos: position{line: 322, col: 13, offset: 7744},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 322, col: 13, offset: 7744},
							name: "_",
						},
						&choiceExpr{
							pos: position{line: 322, col: 16, offset: 7747},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 322, col: 16, offset: 7747},
									val:        "|",
									ignoreCase: false,
								},
								&litMatcher{
									pos:        position{line: 322, col: 22, offset: 7753},
									val:        "<",
									ignoreCase: false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 322, col: 27, offset: 7758},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 322, col: 29, offset: 7760},
							label: "tl",
							expr: &ruleRefExpr{
								pos:  position{line: 322, col: 32, offset: 7763},
								name: "TextList",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 322, col: 41, offset: 7772},
							name: "NL",
						},
					},
				},
			},
		},
		{
			name: "TextList",
			pos:  position{line: 326, col: 1, offset: 7797},
			expr: &choiceExpr{
				pos: position{line: 326, col: 13, offset: 7809},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 326, col: 13, offset: 7809},
						run: (*parser).callonTextList2,
						expr: &seqExpr{
							pos: position{line: 326, col: 13, offset: 7809},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 326, col: 13, offset: 7809},
									label: "intr",
									expr: &ruleRefExpr{
										pos:  position{line: 326, col: 18, offset: 7814},
										name: "Interpolation",
									},
								},
								&labeledExpr{
									pos:   position{line: 326, col: 32, offset: 7828},
									label: "tl",
									expr: &ruleRefExpr{
										pos:  position{line: 326, col: 35, offset: 7831},
										name: "TextList",
									},
								},
							},
						},
					},
					&actionExpr{
						pos: position{line: 341, col: 5, offset: 8152},
						run: (*parser).callonTextList8,
						expr: &andExpr{
							pos: position{line: 341, col: 5, offset: 8152},
							expr: &ruleRefExpr{
								pos:  position{line: 341, col: 6, offset: 8153},
								name: "NL",
							},
						},
					},
					&actionExpr{
						pos: position{line: 343, col: 5, offset: 8218},
						run: (*parser).callonTextList11,
						expr: &seqExpr{
							pos: position{line: 343, col: 5, offset: 8218},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 343, col: 5, offset: 8218},
									label: "ch",
									expr: &anyMatcher{
										line: 343, col: 8, offset: 8221,
									},
								},
								&labeledExpr{
									pos:   position{line: 343, col: 10, offset: 8223},
									label: "tl",
									expr: &ruleRefExpr{
										pos:  position{line: 343, col: 13, offset: 8226},
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
			pos:  position{line: 360, col: 1, offset: 8622},
			expr: &litMatcher{
				pos:        position{line: 360, col: 11, offset: 8632},
				val:        "\x01",
				ignoreCase: false,
			},
		},
		{
			name: "Outdent",
			pos:  position{line: 361, col: 1, offset: 8641},
			expr: &litMatcher{
				pos:        position{line: 361, col: 12, offset: 8652},
				val:        "\x02",
				ignoreCase: false,
			},
		},
		{
			name: "Interpolation",
			pos:  position{line: 363, col: 1, offset: 8662},
			expr: &actionExpr{
				pos: position{line: 363, col: 18, offset: 8679},
				run: (*parser).callonInterpolation1,
				expr: &seqExpr{
					pos: position{line: 363, col: 18, offset: 8679},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 363, col: 18, offset: 8679},
							label: "mod",
							expr: &choiceExpr{
								pos: position{line: 363, col: 23, offset: 8684},
								alternatives: []interface{}{
									&litMatcher{
										pos:        position{line: 363, col: 23, offset: 8684},
										val:        "#",
										ignoreCase: false,
									},
									&litMatcher{
										pos:        position{line: 363, col: 29, offset: 8690},
										val:        "!",
										ignoreCase: false,
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 363, col: 34, offset: 8695},
							val:        "{",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 363, col: 38, offset: 8699},
							label: "expr",
							expr: &ruleRefExpr{
								pos:  position{line: 363, col: 43, offset: 8704},
								name: "Expression",
							},
						},
						&litMatcher{
							pos:        position{line: 363, col: 54, offset: 8715},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Expression",
			pos:  position{line: 371, col: 1, offset: 8899},
			expr: &ruleRefExpr{
				pos:  position{line: 371, col: 15, offset: 8913},
				name: "ExpressionBinOp",
			},
		},
		{
			name: "ExpressionBinOp",
			pos:  position{line: 373, col: 1, offset: 8930},
			expr: &actionExpr{
				pos: position{line: 373, col: 20, offset: 8949},
				run: (*parser).callonExpressionBinOp1,
				expr: &seqExpr{
					pos: position{line: 373, col: 20, offset: 8949},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 373, col: 20, offset: 8949},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 373, col: 26, offset: 8955},
								name: "ExpressionCmpOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 373, col: 42, offset: 8971},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 373, col: 47, offset: 8976},
								expr: &seqExpr{
									pos: position{line: 373, col: 49, offset: 8978},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 373, col: 49, offset: 8978},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 373, col: 51, offset: 8980},
											name: "CmpOp",
										},
										&ruleRefExpr{
											pos:  position{line: 373, col: 57, offset: 8986},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 373, col: 59, offset: 8988},
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
			pos:  position{line: 377, col: 1, offset: 9048},
			expr: &actionExpr{
				pos: position{line: 377, col: 20, offset: 9067},
				run: (*parser).callonExpressionCmpOp1,
				expr: &seqExpr{
					pos: position{line: 377, col: 20, offset: 9067},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 377, col: 20, offset: 9067},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 377, col: 26, offset: 9073},
								name: "ExpressionAddOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 377, col: 42, offset: 9089},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 377, col: 47, offset: 9094},
								expr: &seqExpr{
									pos: position{line: 377, col: 49, offset: 9096},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 377, col: 49, offset: 9096},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 377, col: 51, offset: 9098},
											name: "CmpOp",
										},
										&ruleRefExpr{
											pos:  position{line: 377, col: 57, offset: 9104},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 377, col: 59, offset: 9106},
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
			pos:  position{line: 381, col: 1, offset: 9166},
			expr: &actionExpr{
				pos: position{line: 381, col: 20, offset: 9185},
				run: (*parser).callonExpressionAddOp1,
				expr: &seqExpr{
					pos: position{line: 381, col: 20, offset: 9185},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 381, col: 20, offset: 9185},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 381, col: 26, offset: 9191},
								name: "ExpressionMulOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 381, col: 42, offset: 9207},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 381, col: 47, offset: 9212},
								expr: &seqExpr{
									pos: position{line: 381, col: 49, offset: 9214},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 381, col: 49, offset: 9214},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 381, col: 51, offset: 9216},
											name: "AddOp",
										},
										&ruleRefExpr{
											pos:  position{line: 381, col: 57, offset: 9222},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 381, col: 59, offset: 9224},
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
			pos:  position{line: 385, col: 1, offset: 9284},
			expr: &actionExpr{
				pos: position{line: 385, col: 20, offset: 9303},
				run: (*parser).callonExpressionMulOp1,
				expr: &seqExpr{
					pos: position{line: 385, col: 20, offset: 9303},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 385, col: 20, offset: 9303},
							label: "first",
							expr: &ruleRefExpr{
								pos:  position{line: 385, col: 26, offset: 9309},
								name: "ExpressionUnaryOp",
							},
						},
						&labeledExpr{
							pos:   position{line: 385, col: 44, offset: 9327},
							label: "rest",
							expr: &zeroOrMoreExpr{
								pos: position{line: 385, col: 49, offset: 9332},
								expr: &seqExpr{
									pos: position{line: 385, col: 51, offset: 9334},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 385, col: 51, offset: 9334},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 385, col: 53, offset: 9336},
											name: "MulOp",
										},
										&ruleRefExpr{
											pos:  position{line: 385, col: 59, offset: 9342},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 385, col: 61, offset: 9344},
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
			pos:  position{line: 389, col: 1, offset: 9404},
			expr: &choiceExpr{
				pos: position{line: 389, col: 22, offset: 9425},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 389, col: 22, offset: 9425},
						run: (*parser).callonExpressionUnaryOp2,
						expr: &seqExpr{
							pos: position{line: 389, col: 22, offset: 9425},
							exprs: []interface{}{
								&labeledExpr{
									pos:   position{line: 389, col: 22, offset: 9425},
									label: "op",
									expr: &ruleRefExpr{
										pos:  position{line: 389, col: 25, offset: 9428},
										name: "UnaryOp",
									},
								},
								&ruleRefExpr{
									pos:  position{line: 389, col: 33, offset: 9436},
									name: "_",
								},
								&labeledExpr{
									pos:   position{line: 389, col: 35, offset: 9438},
									label: "ex",
									expr: &ruleRefExpr{
										pos:  position{line: 389, col: 38, offset: 9441},
										name: "ExpressionFactor",
									},
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 391, col: 5, offset: 9564},
						name: "ExpressionFactor",
					},
				},
			},
		},
		{
			name: "ExpressionFactor",
			pos:  position{line: 393, col: 1, offset: 9582},
			expr: &choiceExpr{
				pos: position{line: 393, col: 21, offset: 9602},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 393, col: 21, offset: 9602},
						run: (*parser).callonExpressionFactor2,
						expr: &seqExpr{
							pos: position{line: 393, col: 21, offset: 9602},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 393, col: 21, offset: 9602},
									val:        "(",
									ignoreCase: false,
								},
								&labeledExpr{
									pos:   position{line: 393, col: 25, offset: 9606},
									label: "e",
									expr: &ruleRefExpr{
										pos:  position{line: 393, col: 27, offset: 9608},
										name: "Expression",
									},
								},
								&litMatcher{
									pos:        position{line: 393, col: 38, offset: 9619},
									val:        ")",
									ignoreCase: false,
								},
							},
						},
					},
					&ruleRefExpr{
						pos:  position{line: 395, col: 5, offset: 9645},
						name: "StringExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 395, col: 24, offset: 9664},
						name: "NumberExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 395, col: 43, offset: 9683},
						name: "BooleanExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 395, col: 63, offset: 9703},
						name: "NilExpression",
					},
					&ruleRefExpr{
						pos:  position{line: 395, col: 79, offset: 9719},
						name: "MemberExpression",
					},
				},
			},
		},
		{
			name: "StringExpression",
			pos:  position{line: 397, col: 1, offset: 9737},
			expr: &actionExpr{
				pos: position{line: 397, col: 21, offset: 9757},
				run: (*parser).callonStringExpression1,
				expr: &labeledExpr{
					pos:   position{line: 397, col: 21, offset: 9757},
					label: "s",
					expr: &ruleRefExpr{
						pos:  position{line: 397, col: 23, offset: 9759},
						name: "String",
					},
				},
			},
		},
		{
			name: "NumberExpression",
			pos:  position{line: 401, col: 1, offset: 9854},
			expr: &actionExpr{
				pos: position{line: 401, col: 21, offset: 9874},
				run: (*parser).callonNumberExpression1,
				expr: &seqExpr{
					pos: position{line: 401, col: 21, offset: 9874},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 401, col: 21, offset: 9874},
							expr: &litMatcher{
								pos:        position{line: 401, col: 21, offset: 9874},
								val:        "-",
								ignoreCase: false,
							},
						},
						&ruleRefExpr{
							pos:  position{line: 401, col: 26, offset: 9879},
							name: "Integer",
						},
						&labeledExpr{
							pos:   position{line: 401, col: 34, offset: 9887},
							label: "dec",
							expr: &zeroOrOneExpr{
								pos: position{line: 401, col: 38, offset: 9891},
								expr: &seqExpr{
									pos: position{line: 401, col: 40, offset: 9893},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 401, col: 40, offset: 9893},
											val:        ".",
											ignoreCase: false,
										},
										&oneOrMoreExpr{
											pos: position{line: 401, col: 44, offset: 9897},
											expr: &ruleRefExpr{
												pos:  position{line: 401, col: 44, offset: 9897},
												name: "DecimalDigit",
											},
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 401, col: 61, offset: 9914},
							label: "ex",
							expr: &zeroOrOneExpr{
								pos: position{line: 401, col: 64, offset: 9917},
								expr: &ruleRefExpr{
									pos:  position{line: 401, col: 64, offset: 9917},
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
			pos:  position{line: 411, col: 1, offset: 10247},
			expr: &actionExpr{
				pos: position{line: 411, col: 18, offset: 10264},
				run: (*parser).callonNilExpression1,
				expr: &ruleRefExpr{
					pos:  position{line: 411, col: 18, offset: 10264},
					name: "Null",
				},
			},
		},
		{
			name: "BooleanExpression",
			pos:  position{line: 415, col: 1, offset: 10335},
			expr: &actionExpr{
				pos: position{line: 415, col: 22, offset: 10356},
				run: (*parser).callonBooleanExpression1,
				expr: &labeledExpr{
					pos:   position{line: 415, col: 22, offset: 10356},
					label: "b",
					expr: &ruleRefExpr{
						pos:  position{line: 415, col: 24, offset: 10358},
						name: "Bool",
					},
				},
			},
		},
		{
			name: "MemberExpression",
			pos:  position{line: 419, col: 1, offset: 10450},
			expr: &actionExpr{
				pos: position{line: 419, col: 21, offset: 10470},
				run: (*parser).callonMemberExpression1,
				expr: &seqExpr{
					pos: position{line: 419, col: 21, offset: 10470},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 419, col: 21, offset: 10470},
							label: "field",
							expr: &ruleRefExpr{
								pos:  position{line: 419, col: 27, offset: 10476},
								name: "Field",
							},
						},
						&labeledExpr{
							pos:   position{line: 419, col: 33, offset: 10482},
							label: "member",
							expr: &zeroOrMoreExpr{
								pos: position{line: 419, col: 40, offset: 10489},
								expr: &choiceExpr{
									pos: position{line: 419, col: 41, offset: 10490},
									alternatives: []interface{}{
										&seqExpr{
											pos: position{line: 419, col: 42, offset: 10491},
											exprs: []interface{}{
												&litMatcher{
													pos:        position{line: 419, col: 42, offset: 10491},
													val:        ".",
													ignoreCase: false,
												},
												&ruleRefExpr{
													pos:  position{line: 419, col: 46, offset: 10495},
													name: "Identifier",
												},
											},
										},
										&seqExpr{
											pos: position{line: 419, col: 61, offset: 10510},
											exprs: []interface{}{
												&ruleRefExpr{
													pos:  position{line: 419, col: 61, offset: 10510},
													name: "_",
												},
												&ruleRefExpr{
													pos:  position{line: 419, col: 63, offset: 10512},
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
			pos:  position{line: 436, col: 1, offset: 10970},
			expr: &actionExpr{
				pos: position{line: 436, col: 10, offset: 10979},
				run: (*parser).callonField1,
				expr: &labeledExpr{
					pos:   position{line: 436, col: 10, offset: 10979},
					label: "id",
					expr: &ruleRefExpr{
						pos:  position{line: 436, col: 13, offset: 10982},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "UnaryOp",
			pos:  position{line: 440, col: 1, offset: 11083},
			expr: &actionExpr{
				pos: position{line: 440, col: 12, offset: 11094},
				run: (*parser).callonUnaryOp1,
				expr: &choiceExpr{
					pos: position{line: 440, col: 14, offset: 11096},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 440, col: 14, offset: 11096},
							val:        "+",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 440, col: 20, offset: 11102},
							val:        "-",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 440, col: 26, offset: 11108},
							val:        "!",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "AddOp",
			pos:  position{line: 444, col: 1, offset: 11148},
			expr: &actionExpr{
				pos: position{line: 444, col: 10, offset: 11157},
				run: (*parser).callonAddOp1,
				expr: &choiceExpr{
					pos: position{line: 444, col: 12, offset: 11159},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 444, col: 12, offset: 11159},
							val:        "+",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 444, col: 18, offset: 11165},
							val:        "-",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "MulOp",
			pos:  position{line: 448, col: 1, offset: 11205},
			expr: &actionExpr{
				pos: position{line: 448, col: 10, offset: 11214},
				run: (*parser).callonMulOp1,
				expr: &choiceExpr{
					pos: position{line: 448, col: 12, offset: 11216},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 448, col: 12, offset: 11216},
							val:        "*",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 448, col: 18, offset: 11222},
							val:        "/",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 448, col: 24, offset: 11228},
							val:        "%",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CmpOp",
			pos:  position{line: 452, col: 1, offset: 11268},
			expr: &actionExpr{
				pos: position{line: 452, col: 10, offset: 11277},
				run: (*parser).callonCmpOp1,
				expr: &choiceExpr{
					pos: position{line: 452, col: 12, offset: 11279},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 452, col: 12, offset: 11279},
							val:        "==",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 452, col: 19, offset: 11286},
							val:        "!=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 452, col: 26, offset: 11293},
							val:        "<",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 452, col: 32, offset: 11299},
							val:        ">",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 452, col: 38, offset: 11305},
							val:        "<=",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 452, col: 45, offset: 11312},
							val:        ">=",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "BinOp",
			pos:  position{line: 456, col: 1, offset: 11353},
			expr: &actionExpr{
				pos: position{line: 456, col: 10, offset: 11362},
				run: (*parser).callonBinOp1,
				expr: &choiceExpr{
					pos: position{line: 456, col: 12, offset: 11364},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 456, col: 12, offset: 11364},
							val:        "&&",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 456, col: 19, offset: 11371},
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
			pos:         position{line: 460, col: 1, offset: 11412},
			expr: &actionExpr{
				pos: position{line: 460, col: 20, offset: 11431},
				run: (*parser).callonString1,
				expr: &seqExpr{
					pos: position{line: 460, col: 20, offset: 11431},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 460, col: 20, offset: 11431},
							name: "Quote",
						},
						&zeroOrMoreExpr{
							pos: position{line: 460, col: 26, offset: 11437},
							expr: &choiceExpr{
								pos: position{line: 460, col: 28, offset: 11439},
								alternatives: []interface{}{
									&seqExpr{
										pos: position{line: 460, col: 28, offset: 11439},
										exprs: []interface{}{
											&notExpr{
												pos: position{line: 460, col: 28, offset: 11439},
												expr: &ruleRefExpr{
													pos:  position{line: 460, col: 29, offset: 11440},
													name: "EscapedChar",
												},
											},
											&anyMatcher{
												line: 460, col: 41, offset: 11452,
											},
										},
									},
									&seqExpr{
										pos: position{line: 460, col: 45, offset: 11456},
										exprs: []interface{}{
											&litMatcher{
												pos:        position{line: 460, col: 45, offset: 11456},
												val:        "\\",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 460, col: 50, offset: 11461},
												name: "EscapeSequence",
											},
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 460, col: 68, offset: 11479},
							name: "Quote",
						},
					},
				},
			},
		},
		{
			name:        "Quote",
			displayName: "\"quote\"",
			pos:         position{line: 464, col: 1, offset: 11531},
			expr: &litMatcher{
				pos:        position{line: 464, col: 18, offset: 11548},
				val:        "\"",
				ignoreCase: false,
			},
		},
		{
			name: "EscapedChar",
			pos:  position{line: 466, col: 1, offset: 11553},
			expr: &charClassMatcher{
				pos:        position{line: 466, col: 16, offset: 11568},
				val:        "[\\x00-\\x1f\"\\\\]",
				chars:      []rune{'"', '\\'},
				ranges:     []rune{'\x00', '\x1f'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EscapeSequence",
			pos:  position{line: 467, col: 1, offset: 11583},
			expr: &choiceExpr{
				pos: position{line: 467, col: 19, offset: 11601},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 467, col: 19, offset: 11601},
						name: "SingleCharEscape",
					},
					&ruleRefExpr{
						pos:  position{line: 467, col: 38, offset: 11620},
						name: "UnicodeEscape",
					},
				},
			},
		},
		{
			name: "SingleCharEscape",
			pos:  position{line: 468, col: 1, offset: 11634},
			expr: &charClassMatcher{
				pos:        position{line: 468, col: 21, offset: 11654},
				val:        "[\"\\\\/bfnrt]",
				chars:      []rune{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "UnicodeEscape",
			pos:  position{line: 469, col: 1, offset: 11666},
			expr: &seqExpr{
				pos: position{line: 469, col: 18, offset: 11683},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 469, col: 18, offset: 11683},
						val:        "u",
						ignoreCase: false,
					},
					&ruleRefExpr{
						pos:  position{line: 469, col: 22, offset: 11687},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 469, col: 31, offset: 11696},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 469, col: 40, offset: 11705},
						name: "HexDigit",
					},
					&ruleRefExpr{
						pos:  position{line: 469, col: 49, offset: 11714},
						name: "HexDigit",
					},
				},
			},
		},
		{
			name: "Integer",
			pos:  position{line: 471, col: 1, offset: 11724},
			expr: &choiceExpr{
				pos: position{line: 471, col: 12, offset: 11735},
				alternatives: []interface{}{
					&litMatcher{
						pos:        position{line: 471, col: 12, offset: 11735},
						val:        "0",
						ignoreCase: false,
					},
					&seqExpr{
						pos: position{line: 471, col: 18, offset: 11741},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 471, col: 18, offset: 11741},
								name: "NonZeroDecimalDigit",
							},
							&zeroOrMoreExpr{
								pos: position{line: 471, col: 38, offset: 11761},
								expr: &ruleRefExpr{
									pos:  position{line: 471, col: 38, offset: 11761},
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
			pos:  position{line: 472, col: 1, offset: 11775},
			expr: &seqExpr{
				pos: position{line: 472, col: 13, offset: 11787},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 472, col: 13, offset: 11787},
						val:        "e",
						ignoreCase: true,
					},
					&zeroOrOneExpr{
						pos: position{line: 472, col: 18, offset: 11792},
						expr: &charClassMatcher{
							pos:        position{line: 472, col: 18, offset: 11792},
							val:        "[+-]",
							chars:      []rune{'+', '-'},
							ignoreCase: false,
							inverted:   false,
						},
					},
					&oneOrMoreExpr{
						pos: position{line: 472, col: 24, offset: 11798},
						expr: &ruleRefExpr{
							pos:  position{line: 472, col: 24, offset: 11798},
							name: "DecimalDigit",
						},
					},
				},
			},
		},
		{
			name: "DecimalDigit",
			pos:  position{line: 473, col: 1, offset: 11812},
			expr: &charClassMatcher{
				pos:        position{line: 473, col: 17, offset: 11828},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "NonZeroDecimalDigit",
			pos:  position{line: 474, col: 1, offset: 11834},
			expr: &charClassMatcher{
				pos:        position{line: 474, col: 24, offset: 11857},
				val:        "[1-9]",
				ranges:     []rune{'1', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "HexDigit",
			pos:  position{line: 475, col: 1, offset: 11863},
			expr: &charClassMatcher{
				pos:        position{line: 475, col: 13, offset: 11875},
				val:        "[0-9a-f]i",
				ranges:     []rune{'0', '9', 'a', 'f'},
				ignoreCase: true,
				inverted:   false,
			},
		},
		{
			name: "Bool",
			pos:  position{line: 476, col: 1, offset: 11885},
			expr: &choiceExpr{
				pos: position{line: 476, col: 9, offset: 11893},
				alternatives: []interface{}{
					&actionExpr{
						pos: position{line: 476, col: 9, offset: 11893},
						run: (*parser).callonBool2,
						expr: &litMatcher{
							pos:        position{line: 476, col: 9, offset: 11893},
							val:        "true",
							ignoreCase: false,
						},
					},
					&actionExpr{
						pos: position{line: 476, col: 39, offset: 11923},
						run: (*parser).callonBool4,
						expr: &litMatcher{
							pos:        position{line: 476, col: 39, offset: 11923},
							val:        "false",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Null",
			pos:  position{line: 477, col: 1, offset: 11953},
			expr: &actionExpr{
				pos: position{line: 477, col: 9, offset: 11961},
				run: (*parser).callonNull1,
				expr: &choiceExpr{
					pos: position{line: 477, col: 10, offset: 11962},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 477, col: 10, offset: 11962},
							val:        "null",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 477, col: 19, offset: 11971},
							val:        "nil",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Variable",
			pos:  position{line: 479, col: 1, offset: 11999},
			expr: &actionExpr{
				pos: position{line: 479, col: 13, offset: 12011},
				run: (*parser).callonVariable1,
				expr: &labeledExpr{
					pos:   position{line: 479, col: 13, offset: 12011},
					label: "ident",
					expr: &ruleRefExpr{
						pos:  position{line: 479, col: 19, offset: 12017},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 483, col: 1, offset: 12111},
			expr: &actionExpr{
				pos: position{line: 483, col: 15, offset: 12125},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 483, col: 15, offset: 12125},
					exprs: []interface{}{
						&charClassMatcher{
							pos:        position{line: 483, col: 15, offset: 12125},
							val:        "[a-zA-Z_]",
							chars:      []rune{'_'},
							ranges:     []rune{'a', 'z', 'A', 'Z'},
							ignoreCase: false,
							inverted:   false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 483, col: 25, offset: 12135},
							expr: &charClassMatcher{
								pos:        position{line: 483, col: 25, offset: 12135},
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
			pos:  position{line: 487, col: 1, offset: 12183},
			expr: &actionExpr{
				pos: position{line: 487, col: 14, offset: 12196},
				run: (*parser).callonEmptyLine1,
				expr: &seqExpr{
					pos: position{line: 487, col: 14, offset: 12196},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 487, col: 14, offset: 12196},
							name: "_",
						},
						&charClassMatcher{
							pos:        position{line: 487, col: 16, offset: 12198},
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
			pos:         position{line: 491, col: 1, offset: 12226},
			expr: &actionExpr{
				pos: position{line: 491, col: 19, offset: 12244},
				run: (*parser).callon_1,
				expr: &zeroOrMoreExpr{
					pos: position{line: 491, col: 19, offset: 12244},
					expr: &charClassMatcher{
						pos:        position{line: 491, col: 19, offset: 12244},
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
			pos:         position{line: 492, col: 1, offset: 12271},
			expr: &actionExpr{
				pos: position{line: 492, col: 20, offset: 12290},
				run: (*parser).callon__1,
				expr: &charClassMatcher{
					pos:        position{line: 492, col: 20, offset: 12290},
					val:        "[ \\t]",
					chars:      []rune{' ', '\t'},
					ignoreCase: false,
					inverted:   false,
				},
			},
		},
		{
			name: "NL",
			pos:  position{line: 493, col: 1, offset: 12317},
			expr: &choiceExpr{
				pos: position{line: 493, col: 7, offset: 12323},
				alternatives: []interface{}{
					&charClassMatcher{
						pos:        position{line: 493, col: 7, offset: 12323},
						val:        "[\\n]",
						chars:      []rune{'\n'},
						ignoreCase: false,
						inverted:   false,
					},
					&andExpr{
						pos: position{line: 493, col: 14, offset: 12330},
						expr: &ruleRefExpr{
							pos:  position{line: 493, col: 15, offset: 12331},
							name: "EOF",
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 494, col: 1, offset: 12335},
			expr: &notExpr{
				pos: position{line: 494, col: 8, offset: 12342},
				expr: &anyMatcher{
					line: 494, col: 9, offset: 12343,
				},
			},
		},
	},
}

func (c *current) onInput1(l interface{}) (interface{}, error) {
	return &Root{List: l.(*List)}, nil
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

func (c *current) onTagName1() (interface{}, error) {
	return string(c.text), nil
}

func (p *parser) callonTagName1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagName1()
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

func (c *current) onTagAttribute11(name interface{}) (interface{}, error) {
	return []*Attribute{&Attribute{Name: name.(string), GraphNode: NewNode(pos(c.pos))}}, nil
}

func (p *parser) callonTagAttribute11() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTagAttribute11(stack["name"])
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

func (c *current) onField1(id interface{}) (interface{}, error) {
	return &FieldExpression{Name: string(c.text), GraphNode: NewNode(pos(c.pos))}, nil
}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["id"])
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
