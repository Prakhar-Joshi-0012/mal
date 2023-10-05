package main

import (
	"fmt"
	"io"
	"mal/impls/go/src/printer"
	"mal/impls/go/src/reader"
	. "mal/impls/go/src/types"
	"os"

	"github.com/essentialkaos/go-linenoise/v3"
)

func READ(s string) (MalType, error) {
	return reader.Read_str(s)
}

func EVAL(ast MalType) (MalType, error) {
	return ast, nil
}

func PRINT(ast MalType) (MalType, error) {
	return printer.P_str(ast), nil
}

const PROMPT = "user> "

func REP(s string) (MalType, error) {
	exp, e := READ(s)
	if e != nil {
		return nil, e
	}
	res, e := EVAL(exp)
	if e != nil {
		return nil, e
	}
	out, e := PRINT(res)
	if e != nil {
		return nil, e
	}
	return out, nil
}

func main() {
	out := os.Stdout
	for {
		input, err := linenoise.Line(PROMPT)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		repo, e := REP(input)
		if e != nil {
			fmt.Printf("Error: %v\n", e)
			continue
		}
		io.WriteString(out, fmt.Sprintf("%v\n", repo))
		linenoise.AddHistory(fmt.Sprintf("%v", repo))
	}
}
