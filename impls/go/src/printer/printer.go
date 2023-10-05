package printer

import (
	"fmt"
	. "mal/impls/go/src/types"
	"strings"
)

func P_str(astNode MalType) string {
	switch NodeType := astNode.(type) {
	case MalList:
		return p_list(NodeType.Val, "(", ")", " ")
	case MalVector:
		return p_list(NodeType.Val, "[", "]", " ")
	case MalHash:
		return p_list(NodeType.Val, "{", "}", " ")
	case Atoms:
		return P_str(NodeType.Val)
	case Malsymbols:
		return NodeType.Val
	case nil:
		return "nil"
	case string:
		return `"` + NodeType + `"`
	default:
		return fmt.Sprintf("%v", astNode)
	}
}

func p_list(astList []MalType, start string, end string, join string) string {
	list := make([]string, 0, 1)
	for _, node := range astList {
		list = append(list, P_str(node))
	}
	return start + strings.Join(list, join) + end
}
