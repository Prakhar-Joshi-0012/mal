package main

import (
	"fmt"
	"io"
	"os"

	"github.com/essentialkaos/go-linenoise/v3"
)

func READ(s string) string {
	return s
}

func EVAL(s string) string {
	return s
}

func PRINT(s string) string {
	return s
}

const PROMPT = "user> "

func REP(s string) string {
	return PRINT(EVAL(READ(s)))
}

func main() {
	out := os.Stdout
	for {
		input, err := linenoise.Line(PROMPT)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return
		}
		repo := REP(input)
		io.WriteString(out, repo)
		io.WriteString(out, "\n")
		linenoise.AddHistory(repo)
	}
}
