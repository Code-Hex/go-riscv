package riscv

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// Fetch-decode-execute Cycle
//
// Running image
// ```
// while cpu.pc < cpu.dram.len() as u64 {
// 	// 1. Fetch.
// 	let inst = cpu.fetch();
//
// 	// 2. Add 4 to the program counter.
// 	cpu.pc = cpu.pc + 4;
//
// 	// 3. Decode.
// 	// 4. Execute.
// 	cpu.execute(inst);
// }
// ```
type CPU struct {
	// 32 64-bit integer registers.
	// RV32I, RV64I
	xregs [32]uint64
	// Program counter to hold the dram address of the
	// next instruction that would be executed
	pc uint64
	// Computer dram to store executable instructions.
	dram []byte
}

const dramSize = 1024 * 1024 * 128 // (128MiB).

func NewCPU(code []byte) *CPU {
	regs := [32]uint64{}
	regs[2] = dramSize
	return &CPU{
		xregs: regs,
		pc:    0,
		dram:  code,
	}
}

// Fetch reads the next instruction to be executed from the memory where the program is stored.
//
// see: https://book.rvemu.app/hardware-components/01-cpu.html#fetch-stage
func (c *CPU) Fetch() uint32 {
	index := c.pc
	return binary.LittleEndian.Uint32(c.dram[index : index+4])
}

// Execute performs the action required by the instruction.
func (c *CPU) Execute(inst uint32) {
	opcode := inst & 0x7f

	// 2.2 Base Instruction Formats
	//
	// The RISC-V ISA keeps the source (rs1 and rs2) and destination (rd) registers
	// at the same position
	rd := (inst >> 7) & 0x1f
	rs1 := (inst >> 15) & 0x1f
	rs2 := (inst >> 20) & 0x1f

	c.xregs[0] = 0

	switch opcode {
	case 0x13: // addi
		imm := (inst & 0xfff00000) >> 20
		c.xregs[rd] = c.xregs[rs1] | uint64(imm)
	case 0x33: // add
		c.xregs[rd] = c.xregs[rs1] + uint64(c.xregs[rs2])
	default:
		panic(fmt.Sprintf("unimplemented opcode: %d", opcode))

	}
}

func (c *CPU) DumpRegisters() {
	abi := []string{
		"zero", " ra ", " sp ", " gp ", " tp ", " t0 ", " t1 ", " t2 ", " s0 ", " s1 ", " a0 ",
		" a1 ", " a2 ", " a3 ", " a4 ", " a5 ", " a6 ", " a7 ", " s2 ", " s3 ", " s4 ", " s5 ",
		" s6 ", " s7 ", " s8 ", " s9 ", " s10", " s11", " t3 ", " t4 ", " t5 ", " t6 ",
	}

	var buf strings.Builder
	// 32 for 32-bit. x0-x31
	for i := 0; i < 32; i += 4 {
		fmt.Fprintf(&buf,
			"x%02d(%s) = 0x%08x  x%02d(%s) = 0x%08x  x%02d(%s) = 0x%08x  x%02d(%s) = 0x%08x \n",
			i, abi[i], c.xregs[i],
			i+1, abi[i+1], c.xregs[i+1],
			i+2, abi[i+2], c.xregs[i+2],
			i+3, abi[i+3], c.xregs[i+3],
		)
	}
	fmt.Println(buf.String())
}
