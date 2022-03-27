package riscv

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Int32ToUint32(v int32) uint32 {
	return uint32(v)
}

func TestCPU(t *testing.T) {
	t.Run("binary data", func(t *testing.T) {
		cases := []struct {
			name      string
			code      []byte
			wantXregs [32]uint32
		}{
			{
				name: "lb rd offset rs1",
				code: []byte{
					0x13, 0x08, 0x50, 0x00, // addi x16, x0, 5
					0x93, 0x08, 0x30, 0x00, // addi x17, x0, 3
					0x03, 0x09, 0x40, 0x00, // lb x18, 4(x0)
				},
				wantXregs: [32]uint32{
					2:  dramStartAddress + dramSize,
					15: 0,
					16: 5,
					17: 3,
					18: Int32ToUint32(-109),
				},
			},
		}
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				cpu := NewCPU(tc.code)
				cpu.debug = true
				if err := cpu.Run(); err != nil {
					t.Fatal(err)
				}

				// if tc.name == "want" {
				// 	cpu.DumpRegisters()
				// }

				if diff := cmp.Diff(tc.wantXregs, cpu.xregs); diff != "" {
					t.Fatalf("(-want, +got)\n%s", diff)
				}
			})
		}
	})

	t.Run("read bin files", func(t *testing.T) {
		cases := []struct {
			name      string
			wantXregs [32]uint32
		}{
			{
				name: "add-addi",
				wantXregs: [32]uint32{
					2:  dramStartAddress + dramSize,
					29: 5,
					30: 37,
					31: 42,
				},
			},
		}
		for _, tc := range cases {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				filename := filepath.Join("testdata", tc.name, tc.name+".bin")
				code, err := os.ReadFile(filename)
				if err != nil {
					t.Fatal(err)
				}
				cpu := NewCPU(code)
				if err := cpu.Run(); err != nil {
					t.Fatal(err)
				}

				// if tc.name == "want" {
				// 	cpu.DumpRegisters()
				// }

				if diff := cmp.Diff(tc.wantXregs, cpu.xregs); diff != "" {
					t.Fatalf("(-want, +got)\n%s", diff)
				}
			})
		}
	})
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
