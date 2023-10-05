package reader

import (
	"testing"
)

func TestRead_str(t *testing.T) {
	input := `"abc\"def"`
	tokens := tokenize(input)
	rdr := &TokenReader{tokens: tokens, position: 0}
	t.Errorf("%v size of tokens =%v\n", tokens, len(tokens))
	MalList, err := read_form(rdr)
	if err != nil {
		t.Errorf("Error: ReadForm Error\n")
	}
	t.Errorf("%v", MalList)
}
