package compiler

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"regexp"
	"strings"
)

var rgxIndent = regexp.MustCompile(`^([ \t]+)`)
var indentToken = "\u0001"
var dedentToken = "\u0002"

type lineOffsets map[int]int

func preprocess(in io.Reader) (string, lineOffsets, error) {
	bufin := bufio.NewReader(in)
	indentStack := list.New()
	indentStack.PushBack("")

	lines := 1

	out := ""
	offsets := lineOffsets(make(map[int]int))

	for {
		line, err := bufin.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				if line == "" {
					break
				} else {
					line = line + "\n"
				}
			} else {
				return "", offsets, err
			}
		}

		if strings.TrimSpace(line) != "" {
			curindent := indentStack.Back().Value.(string)
			indent := rgxIndent.FindString(line)

			offsets[lines] = 0

			if indent != curindent && strings.HasPrefix(indent, curindent) {
				out += indentToken
				offsets[lines] += len([]byte(indentToken))

				indentStack.PushBack(indent)
				curindent = indent
			} else {
				for indent != curindent {
					out += dedentToken
					offsets[lines] += len(dedentToken)

					indentStack.Remove(indentStack.Back())

					if indentStack.Back() == nil {
						return "", offsets, fmt.Errorf("Inconsistent indentation found at line %d", lines)
					}

					curindent = indentStack.Back().Value.(string)
				}
			}
		}

		lines++
		out += line
	}

	out += "\n"

	for indentStack.Back() != nil {
		if indentStack.Back().Value.(string) != "" {
			out += dedentToken
		}

		indentStack.Remove(indentStack.Back())
	}

	return out, offsets, nil
}
