package riscv

import (
	"testing"
)

func TestDecodeImmediate(t *testing.T) {
	type args struct {
		rawInst uint32
		fType   InstFormat
	}
	// generated by https://guillaume-savaton-eseo.github.io/emulsiV/
	cases := []struct {
		name string
		args args
		want uint32
	}{
		{
			name: "addi x1, x0, 32",
			args: args{
				rawInst: 0x02000093,
				fType:   IType,
			},
			want: 32,
		},
		{
			name: "addi x1, x0, -40",
			args: args{
				rawInst: 0xfd800093,
				fType:   IType,
			},
			want: 0xffffffd8,
		},
		{
			name: "lui x2, 0xc0000000",
			args: args{
				rawInst: 0xc0000137,
				fType:   UType,
			},
			want: 0xc0000000,
		},
		{
			name: "lui x6, -4096",
			args: args{
				rawInst: 0xfffff337,
				fType:   UType,
			},
			want: 0xfffff000,
		},
		{
			name: "beq x3, x0, +16",
			args: args{
				rawInst: 0x00018863,
				fType:   BType,
			},
			want: 0x00000010,
		},
		{
			name: "bge x3, x2, -16",
			args: args{
				rawInst: 0xfe21d8e3,
				fType:   BType,
			},
			want: 0xfffffff0,
		},
		{
			name: "jal x0, -16",
			args: args{
				rawInst: 0xff1ff06f,
				fType:   JType,
			},
			want: 0xfffffff0,
		},
		{
			name: "jal x0, 0",
			args: args{
				rawInst: 0x0000006f,
				fType:   JType,
			},
			want: 0,
		},
		{
			name: "jal x0, +28",
			args: args{
				rawInst: 0x01c0006f,
				fType:   JType,
			},
			want: 0x0000001c,
		},
		{
			name: "sb x2, 0(x1)",
			args: args{
				rawInst: 0x00208023,
				fType:   SType,
			},
			want: 0,
		},
		{
			name: "sb x2, -30(x1)",
			args: args{
				rawInst: 0xfe208123,
				fType:   SType,
			},
			want: 0xffffffe2,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := decodeImmediate(tc.args.rawInst, tc.args.fType)
			if tc.want != got {
				t.Errorf("want %d but got %d", tc.want, got)
			}
		})
	}
}

func Test_detectInstructionFormat(t *testing.T) {
	type args struct {
		opcode uint32
		funct3 uint32
	}
	tests := []struct {
		name string
		args args
		want InstFormat
	}{
		{
			name: "lui",
			args: args{
				opcode: 0b0110111,
			},
			want: UType,
		},
		{
			name: "auipc",
			args: args{
				opcode: 0b0010111,
			},
			want: UType,
		},
		{
			name: "jal",
			args: args{
				opcode: 0b1101111,
			},
			want: JType,
		},
		{
			name: "jalr",
			args: args{
				opcode: 0b1100111,
			},
			want: IType,
		},
		{
			name: "beq",
			args: args{
				opcode: 0b1100011,
				funct3: 0b000,
			},
			want: BType,
		},
		{
			name: "bne",
			args: args{
				opcode: 0b1100011,
				funct3: 0b001,
			},
			want: BType,
		},
		{
			name: "blt",
			args: args{
				opcode: 0b1100011,
				funct3: 0b100,
			},
			want: BType,
		},
		{
			name: "bge",
			args: args{
				opcode: 0b1100011,
				funct3: 0b101,
			},
			want: BType,
		},
		{
			name: "bltu",
			args: args{
				opcode: 0b1100011,
				funct3: 0b110,
			},
			want: BType,
		},
		{
			name: "bgeu",
			args: args{
				opcode: 0b1100011,
				funct3: 0b111,
			},
			want: BType,
		},
		{
			name: "lb",
			args: args{
				opcode: 0b0000011,
				funct3: 0b000,
			},
			want: IType,
		},
		{
			name: "addi",
			args: args{
				opcode: 0b0010011,
				funct3: 0b000,
			},
			want: IType,
		},
		{
			name: "sll",
			args: args{
				opcode: 0b0110011,
				funct3: 0b001,
			},
			want: RType,
		},
		{
			name: "slli",
			args: args{
				opcode: 0b0010011,
				funct3: 0b001,
			},
			want: IType,
		},
		{
			name: "add",
			args: args{
				opcode: 0b0110011,
				funct3: 0b000,
			},
			want: RType,
		},
		{
			name: "ecall",
			args: args{
				opcode: 0b1110011,
				funct3: 0b000,
			},
			want: IType,
		},
		{
			name: "sb",
			args: args{
				opcode: 0b0100011,
				funct3: 0b000,
			},
			want: SType,
		},
		{
			name: "sh",
			args: args{
				opcode: 0b0100011,
				funct3: 0b001,
			},
			want: SType,
		},
		{
			name: "sw",
			args: args{
				opcode: 0b0100011,
				funct3: 0b010,
			},
			want: SType,
		},
		{
			name: "mul",
			args: args{
				opcode: 0b0110011,
				funct3: 0b000,
			},
			want: RType,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := detectInstructionFormat(tt.args.opcode, tt.args.funct3); got != tt.want {
				t.Errorf("detectInstructionFormat() = %v, want %v", got, tt.want)
			}
		})
	}
}
