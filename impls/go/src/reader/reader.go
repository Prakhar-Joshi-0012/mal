package reader

import (
	"errors"
	"fmt"
	. "mal/impls/go/src/types"
	"regexp"
	"strconv"
	"strings"
)

type Reader interface {
	next() *string
	peek() *string
}

type TokenReader struct {
	tokens   []string
	position int
}

func (trdr *TokenReader) next() *string {
	if trdr.position >= len(trdr.tokens) {
		return nil
	}
	token := trdr.tokens[trdr.position]
	trdr.position++
	return &token
}

func (trdr *TokenReader) peek() *string {
	if trdr.position >= len(trdr.tokens) {
		return nil
	}
	return &trdr.tokens[trdr.position]
}

func Read_str(s string) (MalType, error) {
	tokens := tokenize(s)
	if len(tokens) == 0 {
		return nil, errors.New("empty input string")
	}
	return read_form(&TokenReader{tokens: tokens, position: 0})
}

func tokenize(s string) []string {
	result := make([]string, 0, 1)
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" +
		`,;)]*)`)
	for _, group := range re.FindAllStringSubmatch(s, -1) {
		if group[1] == "" || group[1][0] == ';' {
			continue
		}
		result = append(result, group[1])
	}
	return result
}

func read_form(rdr Reader) (MalType, error) {
	token := rdr.peek()
	if token == nil {
		return nil, errors.New("empty Token read")
	}
	switch *token {
	case `'`:
		rdr.next()
		form, e := read_form(rdr)
		if e != nil {
			return nil, e
		}
		return MalList{Val: []MalType{Malsymbols{Val: "quote"}, form}, Meta: nil}, nil
	case "`":
		rdr.next()
		form, e := read_form(rdr)
		if e != nil {
			return nil, e
		}
		return MalList{Val: []MalType{Malsymbols{Val: "quasiquote"}, form}, Meta: nil}, nil
	case `~@`:
		rdr.next()
		form, e := read_form(rdr)
		if e != nil {
			return nil, e
		}
		return MalList{Val: []MalType{Malsymbols{Val: "splice-unquote"}, form}, Meta: nil}, nil
	case `^`:
		rdr.next()
		meta, e := read_form(rdr)
		if e != nil {
			return nil, e
		}
		form, e := read_form(rdr)
		if e != nil {
			return nil, e
		}
		return MalList{Val: []MalType{Malsymbols{Val: "with-meta"}, form, meta}}, nil
	case `@`:
		rdr.next()
		form, e := read_form(rdr)
		if e != nil {
			return nil, e
		}
		return MalList{Val: []MalType{Malsymbols{Val: "deref"}, form}, Meta: nil}, nil
	case `~`:
		rdr.next()
		form, e := read_form(rdr)
		if e != nil {
			return nil, e
		}
		return MalList{Val: []MalType{Malsymbols{Val: "unquote"}, form}, Meta: nil}, nil
	//readlist
	case "(":
		return read_list(rdr, "(", ")")
	case ")":
		return nil, errors.New("unexpected `)` ")

	// read vector
	case "[":
		return read_vector(rdr)
	case "]":
		return nil, errors.New("unexpected ']")

	// read hash
	case "{":
		return read_hash(rdr)
	case "}":
		return nil, errors.New("unexpected '}'")

	default:
		return read_atom(rdr) // read object
	}
}

func read_list(rdr Reader, start, end string) (MalType, error) {
	ast := []MalType{}
	token := rdr.next()
	if token == nil {
		return nil, errors.New(fmt.Sprintf("empty tokens list read"))
	}
	if *token != start {
		return nil, errors.New(fmt.Sprintf("TokenLiteral doesn't match with %v, got=%v", start, *token))
	}
	token = rdr.peek()
	for ; true; token = rdr.peek() {
		if token == nil {
			return nil, errors.New(fmt.Sprintf("expected %v. got%v", ")", "EOF"))
		}
		if *token == end {
			break
		}
		t, e := read_form(rdr)
		if e != nil {
			return nil, e
		}
		ast = append(ast, t)
	}
	rdr.next()
	return MalList{Val: ast}, nil
}

func read_atom(rdr Reader) (MalType, error) {
	token := rdr.next()
	if token == nil {
		return nil, errors.New(fmt.Sprintf("empty tokens atom read"))
	}
	if match, _ := regexp.MatchString(`^-?[0-9]+$`, *token); match {
		var i int
		var e error
		if i, e = strconv.Atoi(*token); e != nil {
			return nil, errors.New("integer parse error")
		}
		return Atoms{Val: i}, nil
	} else if match, _ := regexp.MatchString(`^"(?:\\.|[^\\"])*"$`, *token); match {
		str := (*token)[1 : len(*token)-1]
		str = strings.Replace(str, `\\`, "\u029e", -1)
		str = strings.Replace(str, `\"`, `\"`, -1)
		str = strings.Replace(str, `\n`, `\n`, -1)
		str = strings.Replace(str, "\u029e", `\\`, -1)
		return str, nil

	} else if (*token)[0] == '"' {
		return nil, errors.New("expected '\"', got EOF")
	} else if *token == "nil" {
		return nil, nil
	} else if *token == "true" {
		return true, nil
	} else if *token == "false" {
		return false, nil
	} else {
		return Malsymbols{Val: *token}, nil
	}
}

func read_vector(rdr Reader) (MalType, error) {
	lst, e := read_list(rdr, "[", "]")
	if e != nil {
		return nil, e
	}
	return MalVector{Val: lst.(MalList).Val, Meta: nil}, nil
}

func read_hash(rdr Reader) (MalType, error) {
	lst, e := read_list(rdr, "{", "}")
	if e != nil {
		return nil, e
	}
	return MalHash{Val: lst.(MalList).Val, Meta: nil}, nil
}
