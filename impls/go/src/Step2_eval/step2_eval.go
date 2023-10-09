package main

import (
	"errors"
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

func EVAL(ast MalType, env map[string]MalType) (MalType, error) {
	switch ast.(type) {
	case MalList:
	default:
		return eval_ast(ast, env)
	}
	if len(ast.(MalList).Val) == 0 {
		return ast, nil
	}
	lst, e := eval_ast(ast, env)
	if e != nil {
		return nil, e
	}
	f, ok := lst.(MalList).Val[0].(func([]MalType) (MalType, error))
	if !ok {
		return nil, errors.New("Non function symbol")
	}
	return f(lst.(MalList).Val[1:])
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
	res, e := EVAL(exp, repl_env)
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
		linenoise.LoadHistory("history.txt")
		io.WriteString(out, fmt.Sprintf("%v\n", repo))
		linenoise.AddHistory(fmt.Sprintf("%v", input))
		linenoise.SaveHistory("history.txt")
	}
}

var repl_env = map[string]MalType{
	"+": func(A []MalType) (MalType, error) {
		if len(A) != 2 {
			return nil, errors.New("wrong # operands given to opertor +\n")
		}
		arg1 := A[0].(Atoms).Val
		arg2 := A[1].(Atoms).Val
		return Atoms{Val: arg1.(int) + arg2.(int)}, nil
	},
	"-": func(A []MalType) (MalType, error) {
		if len(A) != 2 {
			return nil, errors.New("wrong # operands given to opertor -\n")
		}
		arg1 := A[0].(Atoms).Val
		arg2 := A[1].(Atoms).Val
		return Atoms{Val: arg1.(int) - arg2.(int)}, nil

	},
	"*": func(A []MalType) (MalType, error) {
		if len(A) != 2 {
			return nil, errors.New("wrong # operands given to opertor *\n")
		}
		arg1 := A[0].(Atoms).Val
		arg2 := A[1].(Atoms).Val
		return Atoms{Val: arg1.(int) * arg2.(int)}, nil
	},
	"/": func(A []MalType) (MalType, error) {
		if len(A) != 2 {
			return nil, errors.New("wrong # operands given to operator /\n")
		}
		arg1 := A[0].(Atoms).Val
		arg2 := A[1].(Atoms).Val
		return Atoms{Val: arg1.(int) / arg2.(int)}, nil
	},
}

func eval_ast(ast MalType, repl_env map[string]MalType) (MalType, error) {
	switch ast.(type) {
	case MalList:
		lst := []MalType{}
		for _, a := range ast.(MalList).Val {
			exp, e := EVAL(a, repl_env)
			if e != nil {
				return nil, e
			}
			lst = append(lst, exp)
		}
		return MalList{Val: lst, Meta: nil}, nil
	case Malsymbols:
		k := ast.(Malsymbols).Val
		exp, e := repl_env[k]
		if e == false {
			return nil, errors.New("Wrong Symbol given '" + k + "'\n")
		}
		return exp, nil
	case MalVector:
		lst := []MalType{}
		for _, a := range ast.(MalVector).Val {
			exp, e := EVAL(a, repl_env)
			if e != nil {
				return nil, e
			}
			lst = append(lst, exp)
		}
		return MalVector{Val: lst, Meta: nil}, nil
	case MalHash:
		hm := ast.(MalHash)
		new_hm := MalHash{Val: map[string]MalType{}, Meta: nil}
		for k, v := range hm.Val {
			new_val, e := EVAL(v, repl_env)
			if e != nil {
				return nil, e
			}
			new_hm.Val[k] = new_val
		}
		return new_hm, nil
	default:
		return ast, nil
	}
}
