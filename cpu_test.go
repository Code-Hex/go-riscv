package riscv

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCPU(t *testing.T) {
	cases := []struct {
		name      string
		wantXregs [32]uint32
	}{
		{
			name: "add-addi",
			wantXregs: [32]uint32{
				2:  dramSize,
				29: 5,
				30: 37,
				31: 42,
			},
		},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			filename := filepath.Join("testdata", tc.name+".bin")
			code, err := os.ReadFile(filename)
			if err != nil {
				t.Fatal(err)
			}
			cpu := NewCPU(code)
			cpu.Run()

			// if tc.name == "want" {
			// 	cpu.DumpRegisters()
			// }

			if diff := cmp.Diff(tc.wantXregs, cpu.xregs); diff != "" {
				t.Fatalf("(-want, +got)\n%s", diff)
			}
		})
	}
}

func TestSignExtend(t *testing.T) {
	type args struct {
		a       uint32
		bitSize uint32
	}
	cases := []struct {
		args args
		want uint32
	}{
		{
			args: args{
				a:       0b11110110, // 8bit, -10
				bitSize: 8,
			},
			want: math.MaxUint32 - 10 + 1,
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
			want: math.MaxUint32 - 456 + 1, // 0b0001_11000111 + 1 == 0b0001_11001000
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
			want: 0b10000000_00000000_00000000_00000000,
		},
		{
			args: args{
				a:       0b00000000_00000000_00000000_10000000,
				bitSize: 8,
			},
			want: math.MaxUint32 - math.MaxInt8,
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
			got := SignedExtend(tc.args.a, tc.args.bitSize)
			if tc.want != got {
				t.Errorf("want %d but got %d", tc.want, got)
			}
		})
	}
}

func TestUnSignedExtend(t *testing.T) {
	type args struct {
		a       uint32
		bitSize uint32
	}
	cases := []struct {
		args args
		want uint32
	}{
		{
			args: args{
				a:       0b11110110,
				bitSize: 8,
			},
			want: 0b11110110,
		},
		{
			args: args{
				a:       0b00001010,
				bitSize: 8,
			},
			want: 0b00001010,
		},
		{
			args: args{
				a:       0b1110_00111000, // 12bit
				bitSize: 8,
			},
			want: 0b0000_00111000,
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
			want: 0b10000000_00000000_00000000_00000000,
		},
		{
			args: args{
				a:       0b11111111_11111111_11111111_10000000,
				bitSize: 8,
			},
			want: 0b00000000_00000000_00000000_10000000,
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
			got := UnSignedExtend(tc.args.a, tc.args.bitSize)
			if tc.want != got {
				t.Errorf("want %d but got %d", tc.want, got)
			}
		})
	}
}
