package types

import (
	"errors"
	"fmt"
)

type MalType interface {
}

// MalList
type MalList struct {
	Val  []MalType
	Meta MalType
}

// MalSymbols
type Malsymbols struct {
	Val string
}

// Atoms
type Atoms struct {
	Val  MalType
	Meta MalType
}

// Vectors
type MalVector struct {
	Val  []MalType
	Meta MalType
}
type MalHash struct {
	Val  map[string]MalType
	Meta MalType
}

func NewKeyword(s string) (MalType, error) {
	return "\u029e" + s, nil
}

func NewHashMap(seq MalType) (MalType, error) {
	var lst []MalType
	switch seqType := seq.(type) {
	case MalList:
		lst = seqType.Val
	case MalVector:
		lst = seqType.Val
	default:
		return nil, errors.New("Not a sequence")
	}
	if len(lst)%2 == 1 {
		return nil, errors.New("wrong number of arguments given to hash")
	}
	m := map[string]MalType{}
	for i := 0; i < len(lst); i += 2 {
		str, ok := lst[i].(string)
		if !ok {
			return nil, errors.New(fmt.Sprintf("is unexpected string %v given and %v is list", lst[i], lst))
		}
		m[str] = lst[i+1]
	}
	return MalHash{Val: m, Meta: nil}, nil

}
