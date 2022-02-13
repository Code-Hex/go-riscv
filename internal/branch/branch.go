package branch

import "fmt"

const (
	EQ  = "eq"
	NE  = "ne"
	LT  = "lt"
	GE  = "ge"
	LTU = "ltu"
	GEU = "geu"
)

func Comparator(cmp string, rs1, rs2 uint32) bool {
	switch cmp {
	case EQ:
		return rs1 == rs2
	case NE:
		return rs1 != rs2
	case LT:
		return int32(rs1) < int32(rs2)
	case GE:
		return int32(rs1) >= int32(rs2)
	case LTU:
		return rs1 < rs2
	case GEU:
		return rs1 >= rs2
	}
	panic(fmt.Errorf("invalid branch comparator: %q", cmp))
}
