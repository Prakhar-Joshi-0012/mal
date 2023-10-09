package main

import (
	"errors"
	"fmt"
	"io"
	. "mal/impls/go/src/env"
	"mal/impls/go/src/printer"
	"mal/impls/go/src/reader"
	. "mal/impls/go/src/types"
	"os"

	"github.com/essentialkaos/go-linenoise/v3"
)

func READ(s string) (MalType, error) {
	return reader.Read_str(s)
}

func EVAL(ast MalType, env MalEnv) (MalType, error) {
	switch ast.(type) {
	case MalList:
	default:
		return eval_ast(ast, env)
	}
	if len(ast.(MalList).Val) == 0 {
		return ast, nil
	}
	// apply list
	arg0 := ast.(MalList).Val[0]
	var arg1 MalType = nil
	var arg2 MalType = nil
	switch len(ast.(MalList).Val) {
	case 1:
	case 2:
		arg1 = ast.(MalList).Val[1]
	default:
		arg1 = ast.(MalList).Val[1]
		arg2 = ast.(MalList).Val[2]
	}
	arg0Sym := "non_function"
	if _, ok := arg0.(Malsymbols); ok {
		arg0Sym = arg0.(Malsymbols).Val
	}
	switch arg0Sym {
	case "def!":
		val, e := EVAL(arg2, env)
		if e != nil {
			return nil, e
		}
		return env.Set(arg1.(Malsymbols), val), nil
	case "let*":
		env_new := NewEnv(env)
		lst, e := GetSlice(arg1)
		if e != nil {
			return nil, e
		}
		for i := 0; i < len(lst); i += 2 {
			if _, ok := lst[i].(Malsymbols); !ok {
				return nil, errors.New("non-symbol bind value")
			}
			exp, e := EVAL(lst[i+1], env_new)
			if e != nil {
				return nil, e
			}
			env_new.Set(lst[i].(Malsymbols), exp)
		}
		return EVAL(arg2, env_new)
	default:
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

var repl_env = NewEnv(nil)

func main() {
	// Env Setup
	repl_env.Set(Malsymbols{Val: "+"}, func(A []MalType) (MalType, error) {
		if len(A) != 2 {
			return nil, errors.New("wrong # operands given to opertor +\n")
		}
		arg1 := A[0].(Atoms).Val
		arg2 := A[1].(Atoms).Val
		return Atoms{Val: arg1.(int) + arg2.(int)}, nil
	})
	repl_env.Set(Malsymbols{Val: "-"}, func(A []MalType) (MalType, error) {
		if len(A) != 2 {
			return nil, errors.New("wrong # operands given to operator -\n")
		}
		arg1 := A[0].(Atoms).Val
		arg2 := A[1].(Atoms).Val
		return Atoms{Val: arg1.(int) - arg2.(int)}, nil
	})
	repl_env.Set(Malsymbols{Val: "*"}, func(A []MalType) (MalType, error) {
		if len(A) != 2 {
			return nil, errors.New("wrong # operands given to operator *\n")
		}
		arg1 := A[0].(Atoms).Val
		arg2 := A[1].(Atoms).Val
		return Atoms{Val: arg1.(int) * arg2.(int)}, nil
	})
	repl_env.Set(Malsymbols{Val: "/"}, func(A []MalType) (MalType, error) {
		if len(A) != 2 {
			return nil, errors.New("wrong # operands given to operator\n")
		}
		arg1 := A[0].(Atoms).Val
		arg2 := A[1].(Atoms).Val
		return Atoms{Val: arg1.(int) / arg2.(int)}, nil
	})

	// Repl Setup
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

func eval_ast(ast MalType, repl_env MalEnv) (MalType, error) {
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
		exp, e := repl_env.Get(ast.(Malsymbols))
		if e != nil {
			return nil, errors.New(fmt.Sprintf("%v not found.", ast.(Malsymbols).Val))
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
