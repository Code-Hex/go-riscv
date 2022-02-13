package alu

import "fmt"

const (
	ADD  = "add"
	SUB  = "sub"
	OR   = "or"
	XOR  = "xor"
	AND  = "and"
	SLL  = "sll"
	SLT  = "slt"
	SLTU = "sltu"
	SRL  = "srl"
	SRA  = "sra"
)

// https://sites.pitt.edu/~kmram/CoE0147/lectures/datapath3.pdf
// https://msyksphinz-self.github.io/riscv-isadoc/html/rvi.html#slti
func Compute(op string, rs1, rs2 uint32) uint32 {
	switch op {
	case ADD:
		return rs1 + rs2
	case SUB:
		return rs1 - rs2
	case OR:
		return rs1 | rs2
	case XOR:
		return rs1 ^ rs2
	case AND:
		return rs1 & rs2
	case SLL:
		return rs1 << rs2
	case SLT:
		if int32(rs1) < int32(rs2) {
			return 1
		}
		return 0
	case SLTU:
		if rs1 < rs2 {
			return 1
		}
		return 0
	case SRL:
		return rs1 >> rs2
	case SRA:
		return rs1 >> int32(rs2)
	}
	panic(fmt.Errorf("invalid ALU operation: %q", op))
}
