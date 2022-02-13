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
	xregs [32]uint32
	// Program counter to hold the dram address of the
	// next instruction that would be executed
	pc     uint64
	nextpc uint64
	// Computer dram to store executable instructions.
	dram []byte
}

const dramSize = 1024 * 1024 * 128 // (128MiB).

func NewCPU(code []byte) *CPU {
	regs := [32]uint32{
		2: dramSize,
	}
	return &CPU{
		xregs: regs,
		pc:    0,
		dram:  code,
	}
}

func (c *CPU) Next() bool {
	c.pc = c.nextpc
	c.nextpc += 4
	return c.pc < uint64(len(c.dram))
}

func (c *CPU) Run() {
	for c.Next() {
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
	return binary.LittleEndian.Uint32(c.dram[c.pc:c.nextpc])
}

func (c *CPU) Decode(rawInst uint32) *Instruction {
	// 2.2 Base Instruction Formats
	//
	// The RISC-V ISA keeps the source (rs1 and rs2) and destination (rd) registers
	// at the same position
	opcode := rawInst & 0b1111111
	funct3 := (rawInst >> 12) & 0b111
	fotmat := detectInstructionFormat(opcode, funct3)
	return &Instruction{
		raw:    rawInst,
		opcode: opcode,                      // bits 0 to 6
		rd:     (rawInst >> 7) & 0b11111,    // bits 7 to 11
		funct3: funct3,                      // bits 12 to 14
		rs1:    (rawInst >> 15) & 0b11111,   // bits 15 to 19
		rs2:    (rawInst >> 20) & 0b11111,   // bits 20 to 24
		funct7: (rawInst >> 25) & 0b1111111, // bits 25 to 31
		format: fotmat,
		imm:    decodeImmediate(rawInst, fotmat),
	}
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
		c.xregs[rd] = c.xregs[rs1] | inst.imm
	case 0x33: // add
		c.xregs[rd] = c.xregs[rs1] + c.xregs[rs2]
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
			strconv.FormatUint(uint64(xreg), 10),
			fmt.Sprintf("0x%08x", xreg),
			fmt.Sprintf("0b%032b", xreg),
		})
	}
	table.Render()
	fmt.Println(buf.String())
}

// SignedExtend extends value to be bits length as signed to 32 bit
func SignedExtend(value, bitSize uint32) uint32 {
	tmp := 32 - bitSize
	return uint32(int32(value) << tmp >> tmp)
}

// UnSignedExtend extends value to be bits length as unsigned to 32 bit
func UnSignedExtend(value, bitSize uint32) uint32 {
	tmp := 32 - bitSize
	return value << tmp >> tmp
}
