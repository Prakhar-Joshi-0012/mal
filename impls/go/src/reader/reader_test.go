package reader

import (
	"testing"
)

func TestRead_str(t *testing.T) {
	input := `{  :a  {:b   {  :cde     3   }  }}`
	tokens := tokenize(input)
	rdr := &TokenReader{tokens: tokens, position: 0}
	t.Errorf("%v size of tokens =%v\n", tokens, len(tokens))
	MalList, err := read_form(rdr)
	if err != nil {
		t.Errorf("Error: ReadForm Error %v \n", err)
	}
	t.Errorf("%v", MalList)
}
