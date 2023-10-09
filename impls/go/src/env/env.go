package env

import (
	"errors"
	"fmt"
	. "mal/impls/go/src/types"
)

type Env struct {
	data     map[string]MalType
	outerEnv MalEnv
}

func NewEnv(outer MalEnv) MalEnv {
	return Env{data: map[string]MalType{}, outerEnv: outer}
}

func (e Env) Set(key Malsymbols, val MalType) MalType {
	e.data[key.Val] = val
	return val
}

func (e Env) Find(key Malsymbols) MalEnv {
	if _, ok := e.data[key.Val]; ok {
		return e
	} else if e.outerEnv != nil {
		return e.outerEnv.Find(key)
	} else {
		return nil
	}
}

func (e Env) Get(key Malsymbols) (MalType, error) {
	env := e.Find(key)
	if env == nil {
		return nil, errors.New(fmt.Sprintf("Symbol %v:  not found", key.Val))
	}
	return env.(Env).data[key.Val], nil
}
