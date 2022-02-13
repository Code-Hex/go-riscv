package riscv

import (
	"encoding/binary"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Code-Hex/go-riscv/internal/alu"
	"github.com/Code-Hex/go-riscv/internal/branch"
	"github.com/olekukonko/tablewriter"
)

type CPU struct {
	// 32 64-bit integer registers.
	// RV32I, RV64I
	xregs [32]uint32
	// Program counter to hold the dram address of the
	// next instruction that would be executed
	pc     uint32
	nextpc uint32
	// Computer dram to store executable instructions.
	dram []byte

	debug bool
}

func (c *CPU) debugf(format string, v ...interface{}) {
	if c.debug {
		log.Printf(format, v...)
	}
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
	return c.pc < uint32(len(c.dram))
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

	rd := inst.rd
	rs1 := inst.rs1
	rs2 := inst.rs2

	// Chapter 19 RV32/64G Instruction Set Listings
	switch inst.opcode {
	case OPIMM:
		switch inst.funct3 {
		case 0b000:
			c.debugf("addi rd, rs1=%d, imm=%d", c.xregs[rs1], inst.imm)
			c.xregs[rd] = alu.Compute(alu.ADD, c.xregs[rs1], inst.imm)
			return
		case 0b001:
			c.debugf("slli rd, rs1=%d, shamt=%d", c.xregs[rs1], inst.imm)
			c.xregs[rd] = alu.Compute(alu.SLL, c.xregs[rs1], inst.imm)
			return
		case 0b101:
			switch inst.funct7 {
			case 0b0000000:
				c.debugf("srli rd, rs1=%d, shamt=%d", c.xregs[rs1], inst.imm)
				c.xregs[rd] = alu.Compute(alu.SRL, c.xregs[rs1], inst.imm)
				return
			case 0b0100000:
				shamt := SignedExtend(inst.imm, 4) // shamt ~ 4 bit range
				c.debugf("srai rd, rs1=%d, shamt=%d", c.xregs[rs1], shamt)
				c.xregs[rd] = alu.Compute(alu.SRA, c.xregs[rs1], shamt)
				return
			}
		case 0b010:
			c.debugf("slti rd, rs1=%d, imm=%d", c.xregs[rs1], inst.imm)
			c.xregs[rd] = alu.Compute(alu.SLT, c.xregs[rs1], inst.imm)
			return
		case 0b011:
			c.debugf("sltiu rd, rs1=%d, imm=%d", c.xregs[rs1], inst.imm)
			c.xregs[rd] = alu.Compute(alu.SLTU, c.xregs[rs1], inst.imm)
			return
		case 0b100:
			c.debugf("xori rd, rs1=%d, imm=%d", c.xregs[rs1], inst.imm)
			c.xregs[rd] = alu.Compute(alu.XOR, c.xregs[rs1], inst.imm)
			return
		case 0b110:
			c.debugf("ori rd, rs1=%d, imm=%d", c.xregs[rs1], inst.imm)
			c.xregs[rd] = alu.Compute(alu.OR, c.xregs[rs1], inst.imm)
			return
		case 0b111:
			c.debugf("andi rd, rs1=%d, imm=%d", c.xregs[rs1], inst.imm)
			c.xregs[rd] = alu.Compute(alu.AND, c.xregs[rs1], inst.imm)
			return
		}
	case OPREG:
		switch inst.funct3 {
		case 0b000:
			switch inst.funct7 {
			case 0b0000000:
				c.debugf("add rd, rs1=%d, rs2=%d", c.xregs[rs1], c.xregs[rs2])
				c.xregs[rd] = alu.Compute(alu.ADD, c.xregs[rs1], c.xregs[rs2])
				return
			case 0b0100000:
				c.debugf("sub rd, rs1=%d, rs2=%d", c.xregs[rs1], c.xregs[rs2])
				c.xregs[rd] = alu.Compute(alu.SUB, c.xregs[rs1], c.xregs[rs2])
				return
			}
		case 0b001:
			c.debugf("sll rd, rs1=%d, rs2=%d", c.xregs[rs1], c.xregs[rs2])
			c.xregs[rd] = alu.Compute(alu.SLL, c.xregs[rs1], c.xregs[rs2])
			return
		case 0b010:
			c.debugf("slt rd, rs1=%d, rs2=%d", c.xregs[rs1], c.xregs[rs2])
			c.xregs[rd] = alu.Compute(alu.SLT, c.xregs[rs1], c.xregs[rs2])
			return
		case 0b011:
			c.debugf("sltu rd, rs1=%d, rs2=%d", c.xregs[rs1], c.xregs[rs2])
			c.xregs[rd] = alu.Compute(alu.SLTU, c.xregs[rs1], c.xregs[rs2])
			return
		case 0b100:
			c.debugf("xor rd, rs1=%d, rs2=%d", c.xregs[rs1], c.xregs[rs2])
			c.xregs[rd] = alu.Compute(alu.XOR, c.xregs[rs1], c.xregs[rs2])
			return
		case 0b101:
			switch inst.funct7 {
			case 0b0000000:
				c.debugf("srl rd, rs1=%d, rs2=%d", c.xregs[rs1], c.xregs[rs2])
				c.xregs[rd] = alu.Compute(alu.SRL, c.xregs[rs1], c.xregs[rs2])
				return
			case 0b0100000:
				c.debugf("sra rd, rs1=%d, rs2=%d", c.xregs[rs1], c.xregs[rs2])
				c.xregs[rd] = alu.Compute(alu.SRA, c.xregs[rs1], c.xregs[rs2])
				return
			}
		case 0b110:
			c.debugf("or rd, rs1=%d, rs2=%d", c.xregs[rs1], c.xregs[rs2])
			c.xregs[rd] = alu.Compute(alu.OR, c.xregs[rs1], c.xregs[rs2])
			return
		case 0b111:
			c.debugf("and rd, rs1=%d, rs2=%d", c.xregs[rs1], c.xregs[rs2])
			c.xregs[rd] = alu.Compute(alu.AND, c.xregs[rs1], c.xregs[rs2])
			return
		}
	case OPAUIPC:
		c.debugf("auipc rd, imm=%d", inst.imm)
		c.xregs[rd] = alu.Compute(alu.ADD, c.pc, inst.imm)
		return
	case OPLUI:
		c.debugf("lui rd, imm=%d", inst.imm)
		c.xregs[rd] = inst.imm
		return
	case OPJAL:
		c.debugf("jal rd, offset=%d", inst.imm)
		c.xregs[rd] = c.pc + 4
		c.pc += inst.imm
		return
	case OPJALR:
		c.debugf("jalr rd, rs1=%d, offset=%d", c.xregs[rs1], inst.imm)
		t := c.pc + 4
		c.xregs[rd] = t
		c.pc = (c.xregs[rs1] + inst.imm) &^ 1
		return
	case OPBRANCH:
		switch inst.funct3 {
		case 0b000:
			c.debugf("beq rs1=%d, rs2=%d, offset=%d", c.xregs[rs1], c.xregs[rs2], inst.imm)
			if branch.Comparator(branch.EQ, c.xregs[rs1], c.xregs[rs2]) {
				c.pc += inst.imm
			}
			return
		case 0b001:
			c.debugf("bne rs1=%d, rs2=%d, offset=%d", c.xregs[rs1], c.xregs[rs2], inst.imm)
			if branch.Comparator(branch.NE, c.xregs[rs1], c.xregs[rs2]) {
				c.pc += inst.imm
			}
			return
		case 0b100:
			c.debugf("blt rs1=%d, rs2=%d, offset=%d", c.xregs[rs1], c.xregs[rs2], inst.imm)
			if branch.Comparator(branch.LT, c.xregs[rs1], c.xregs[rs2]) {
				c.pc += inst.imm
			}
			return
		case 0b101:
			c.debugf("bge rs1=%d, rs2=%d, offset=%d", c.xregs[rs1], c.xregs[rs2], inst.imm)
			if branch.Comparator(branch.GE, c.xregs[rs1], c.xregs[rs2]) {
				c.pc += inst.imm
			}
			return
		case 0b110:
			c.debugf("bltu rs1=%d, rs2=%d, offset=%d", c.xregs[rs1], c.xregs[rs2], inst.imm)
			if branch.Comparator(branch.LTU, c.xregs[rs1], c.xregs[rs2]) {
				c.pc += inst.imm
			}
			return
		case 0b111:
			c.debugf("bgeu rs1=%d, rs2=%d, offset=%d", c.xregs[rs1], c.xregs[rs2], inst.imm)
			if branch.Comparator(branch.GEU, c.xregs[rs1], c.xregs[rs2]) {
				c.pc += inst.imm
			}
			return
		}
	case OPLOAD:
	case OPSTORE:
	case OPSYSTEM:
	}
	panic(fmt.Sprintf("unimplemented opcode: %d", inst.opcode))
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
