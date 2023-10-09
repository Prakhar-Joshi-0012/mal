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
		lst := make([]string, 0, len(NodeType.Val))
		for k, v := range NodeType.Val {
			lst = append(lst, P_str(k))
			lst = append(lst, P_str(v))
		}
		return "{" + strings.Join(lst, " ") + "}"
	case Atoms:
		return P_str(NodeType.Val)
	case Malsymbols:
		return NodeType.Val
	case nil:
		return "nil"
	case string:
		if strings.HasPrefix(NodeType, "\u029e") {
			return ":" + NodeType[2:]
		}
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
