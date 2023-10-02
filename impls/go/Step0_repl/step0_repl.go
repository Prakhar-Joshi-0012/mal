package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
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
	in := os.Stdin
	out := os.Stdout

	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		repo := REP(scanner.Text())
		io.WriteString(out, repo)
		io.WriteString(out, "\n")
	}
}
