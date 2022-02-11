package riscv

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"

	"github.com/olekukonko/tablewriter"
)

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

func (c *CPU) Run() {
	for c.pc < uint64(len(c.dram)) {
		// 1. Fetch.
		inst := c.Fetch()
		// 2. Decode.
		decoded := c.Decode(inst)
		// 3. Execute.
		c.Execute(decoded)
	}
}

// Fetch reads the next instruction to be executed from the memory where the program is stored.
//
// see: https://book.rvemu.app/hardware-components/01-cpu.html#fetch-stage
func (c *CPU) Fetch() uint32 {
	current := c.pc
	c.pc += 4
	return binary.LittleEndian.Uint32(c.dram[current:c.pc])
}

func (c *CPU) Decode(rawInst uint32) *Instruction {
	// 2.2 Base Instruction Formats
	//
	// The RISC-V ISA keeps the source (rs1 and rs2) and destination (rd) registers
	// at the same position
	return &Instruction{
		raw:    rawInst,
		opcode: rawInst & 0x7f,         // bits 0 to 6
		rd:     (rawInst >> 7) & 0x1f,  // bits 7 to 11
		funct3: (rawInst >> 12) & 0x1f, // bits 12 to 14
		rs1:    (rawInst >> 15) & 0x1f, // bits 15 to 19
		rs2:    (rawInst >> 20) & 0x1f, // bits 20 to 24
		funct7: (rawInst >> 20) & 0x1f, // bits 25 to 31
	}
}

type Instruction struct {
	raw    uint32
	opcode uint32
	rd     uint32
	funct3 uint32
	rs1    uint32
	rs2    uint32
	funct7 uint32
}

// Execute performs the action required by the instruction.
func (c *CPU) Execute(inst *Instruction) {
	c.xregs[0] = 0

	opcode := inst.opcode
	rd := inst.rd
	rs1 := inst.rs1
	rs2 := inst.rs2

	// Chapter 19 RV32/64G Instruction Set Listings
	switch opcode {
	case 0b0010011: // addi
		imm := (inst.raw & 0xfff00000) >> 20
		c.xregs[rd] = c.xregs[rs1] | uint64(imm)
	case 0x33: // add
		c.xregs[rd] = c.xregs[rs1] + uint64(c.xregs[rs2])
	default:
		panic(fmt.Sprintf("unimplemented opcode: %d", opcode))

	}
}

func (c *CPU) DumpRegisters() {
	var buf strings.Builder
	table := tablewriter.NewWriter(&buf)
	table.SetHeader([]string{
		"Register",
		"Decimal",
		"Hex",
		"Binary",
	})

	for i, xreg := range c.xregs {
		table.Append([]string{
			xregsABINames[i],
			strconv.FormatUint(xreg, 10),
			fmt.Sprintf("0x%08x", xreg),
			fmt.Sprintf("0b%032b", xreg),
		})
	}
	table.Render()
	fmt.Println(buf.String())
}

// SignedExtend assumes value to be bits length and sign extends to 32 bit
func SignExtend(value, bitSize uint32) int32 {
	tmp := 32 - bitSize
	return int32(value) << tmp >> tmp
}
