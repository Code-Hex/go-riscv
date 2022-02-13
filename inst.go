package riscv

import "fmt"

// Green Card
// https://www.cl.cam.ac.uk/teaching/1617/ECAD+Arch/files/docs/RISCVGreenCardv8-20151013.pdf

// Instruction represents decoded instruction which is fetched
type Instruction struct {
	raw    uint32
	opcode uint32
	rd     uint32
	funct3 uint32
	rs1    uint32
	rs2    uint32
	funct7 uint32
	format InstFormat
	imm    uint32
}

// Instformat is a format of instruction.
type InstFormat string

const (
	// RType represent R instruction format.
	RType InstFormat = "R"
	// IType represent I instruction format.
	IType InstFormat = "I"
	// SType represent S instruction format.
	SType InstFormat = "S"
	// BType represent B instruction format.
	BType InstFormat = "B"
	// UType represent U instruction format.
	UType InstFormat = "U"
	// JType represent J instruction format.
	JType InstFormat = "J"
)

// Chapter 19. RV32/64G Instruction Set Listings
// And I referenced https://guillaume-savaton-eseo.github.io/emulsiV/doc/
//
// +---------------+--------------------+-------------+------------+------------+-------------+--------+
// | Format / Bits | 31:25              | 24:20       | 19:15      | 14:12      | 11:7        | 6:0    |
// +---------------+--------------------+-------------+------------+------------+-------------+--------+
// | R             | funct7             | rs2         | rs1        | funct3     | rd          | opcode |
// +---------------+--------------------+-------------+------------+------------+-------------+--------+
// | I             | imm[11:5]`/`funct7 | imm[4:0]    | rs1        | funct3     | rd          | opcode |
// +---------------+--------------------+-------------+------------+------------+-------------+--------+
// | S             | imm[11:5]          | rs2         | rs1        | funct3     | imm[4:0]    | opcode |
// +---------------+--------------------+-------------+------------+------------+-------------+--------+
// | B             | imm[12,10:5]       | rs2         | rs1        | funct3     | imm[4:0,11] | opcode |
// +---------------+--------------------+-------------+------------+------------+-------------+--------+
// | U             | imm[31:25]         | imm[24:20]  | imm[19:15] | imm[14:12] | rd          | opcode |
// +---------------+--------------------+-------------+------------+------------+-------------+--------+
// | J             | imm[20,10:5]       | imm[4:1,11] | imm[19:15] | imm[14:12] | rd          | opcode |
// +---------------+--------------------+-------------+------------+------------+-------------+--------+
func detectInstructionFormat(opcode, funct3 uint32) InstFormat {
	// RV32I Base Instruction Set
	switch opcode {
	case 0b0110011: // in RV32I Base Instruction Set
		return RType
	case 0b1100111, 0b0000011, 0b0010011, 0b0001111:
		return IType
	case 0b0100011:
		return SType
	case 0b1100011:
		return BType
	case 0b0110111, 0b0010111: // LUI, AUIPC
		return UType
	case 0b1101111: // JAL
		return JType
	}

	panic(fmt.Errorf("unexpected opcode: 0b%07b", opcode))
}

// I have referenced https://guillaume-savaton-eseo.github.io/emulsiV/doc/
//
// +─────────+─────────────+─────────────+──────────+─────────────+─────────────+──────────+──────────────+──────────────+───────────+
// | Format  | imm[31:25]  | imm[24:21]  | imm[20]  | imm[19:15]  | imm[14:12]  | imm[11]  | imm[10:5]    | imm[4:1]     | imm[0]    |
// +─────────+─────────────+─────────────+──────────+─────────────+─────────────+──────────+──────────────+──────────────+───────────+
// | I       | inst[31]    | inst[31]    | inst[31] | inst[31]    | inst[31]    | inst[31] | inst[30:25]  | inst[24:21]  | inst[20]  |
// | S       | inst[31]    | inst[31]    | inst[31] | inst[31]    | inst[31]    | inst[31] | inst[30:25]  | inst[11:8]   | inst[7]   |
// | B       | inst[31]    | inst[31]    | inst[31] | inst[31]    | inst[31]    | inst[7]  | inst[30:25]  | inst[11:8]   | 0         |
// | U       | inst[31:25] | inst[24:21] | inst[20] | inst[19:15] | inst[14:12] | 0        | 0            | 0            | 0         |
// | J       | inst[31]    | inst[31]    | inst[31] | inst[19:15] | inst[14:12] | inst[20] | inst[30:25]  | inst[24:21]  | 0         |
// +─────────+─────────────+─────────────+──────────+─────────────+─────────────+──────────+──────────────+──────────────+───────────+
func decodeImmediate(rawInst uint32, fType InstFormat) uint32 {
	switch fType {
	case IType:
		return SignedExtend(rawInst>>20, 12) // 12 bit range (imm[11] is the MSB)
	case SType:
		imm0 := UnSignedExtend(rawInst>>25, 12) << 5 // start inst from 25 bit, start imm from 5 bit
		imm1 := UnSignedExtend(rawInst>>7, 5)        // 5 bit range
		return SignedExtend(imm0|imm1, 12)           // 12 bit range (imm[11] is the MSB)
	case BType:
		imm0 := UnSignedExtend(rawInst>>31, 1) << 12 // start inst from 31 bit, start imm from 12 bit
		imm1 := UnSignedExtend(rawInst>>7, 1) << 11  // start inst from 7 bit, start imm from 11 bit, 1 bit range
		imm2 := UnSignedExtend(rawInst>>25, 6) << 5  // start inst from 25 bit, start imm from 5 bit, 6 bit range
		imm3 := UnSignedExtend(rawInst>>8, 4) << 1   // start inst from 8 bit, start imm from 1 bit, 4 bit range
		return SignedExtend(imm0|imm1|imm2|imm3, 13)
	case UType:
		return SignedExtend(rawInst>>12, 20) << 12 // start inst from 12 bit, start imm from 12 bit
	case JType:
		imm0 := UnSignedExtend(rawInst>>31, 1) << 20 // start inst from 31 bit, start imm from 20 bit
		imm1 := UnSignedExtend(rawInst>>12, 8) << 12 // start inst from 12 bit, start imm from 12 bit, 7 bit range
		imm2 := UnSignedExtend(rawInst>>20, 1) << 11 // start inst from 20 bit, start imm from 11 bit, 1 bit range
		imm3 := UnSignedExtend(rawInst>>21, 10) << 1 // start inst from 21 bit, start imm from 1 bit, 10 bit range
		return SignedExtend(imm0|imm1|imm2|imm3, 21) // 21 bit range (imm[20] is the MSB)
	}
	return 0
}
