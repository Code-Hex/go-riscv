package riscv

import (
	"fmt"
	"math"
	"os"
	"testing"
)

func TestCPU(t *testing.T) {
	code, err := os.ReadFile("testdata/add-addi.bin")
	if err != nil {
		t.Fatal(err)
	}
	cpu := NewCPU(code)
	for cpu.pc < uint64(len(cpu.dram)) {
		// 1. Fetch.
		inst := cpu.Fetch()
		// 2. Decode.
		decoded := cpu.Decode(inst)
		// 3. Execute.
		cpu.Execute(decoded)
	}
	cpu.DumpRegisters()
}

func TestSignExtend(t *testing.T) {
	type args struct {
		a       uint32
		bitSize uint32
	}
	cases := []struct {
		args args
		want int32
	}{
		{
			args: args{
				a:       0b11110110, // 8bit, -10
				bitSize: 8,
			},
			want: -10,
		},
		{
			args: args{
				a:       0b00001010, // 8bit, 10
				bitSize: 8,
			},
			want: 10,
		},
		{
			args: args{
				a:       0b1110_00111000, // 12bit, -456
				bitSize: 12,
			},
			want: -456, // 0b0001_11000111 + 1 == 0b0001_11001000
		},
		{
			args: args{
				a:       math.MaxInt32, // 32bit, max int32
				bitSize: 32,
			},
			want: math.MaxInt32,
		},
		{
			args: args{
				a:       0b10000000_00000000_00000000_00000000, // 32bit, min int32
				bitSize: 32,
			},
			want: math.MinInt32,
		},
		{
			args: args{
				a:       0b00000000_00000000_00000000_10000000,
				bitSize: 8,
			},
			want: math.MinInt8,
		},
		{
			args: args{
				a:       0,
				bitSize: 8,
			},
			want: 0,
		},
	}
	for _, tc := range cases {
		tc := tc
		name := fmt.Sprintf("(%#v)", tc.args)
		t.Run(name, func(t *testing.T) {
			got := SignExtend(tc.args.a, tc.args.bitSize)
			if tc.want != got {
				t.Errorf("want %d but got %d", tc.want, got)
			}
		})
	}

}
