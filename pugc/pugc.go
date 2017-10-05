package main

import (
	"flag"
	"fmt"
	"os"

	pug "github.com/eknkc/pug"
)

var prettyPrint bool

func init() {
	flag.BoolVar(&prettyPrint, "prettyprint", true, "Use pretty indentation in output html.")
	flag.BoolVar(&prettyPrint, "pp", true, "Use pretty indentation in output html.")

	flag.Parse()
}

func main() {
	input := flag.Arg(0)

	if len(input) == 0 {
		fmt.Fprintln(os.Stderr, "Please provide an input file. (amberc input.amber)")
		os.Exit(1)
	}

	result, err := pug.ParseFile(input, pug.Options{
		PrettyPrint: prettyPrint,
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	_, err = os.Stdout.WriteString(result)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
