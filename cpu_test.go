package riscv

import (
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

		// 2. Add 4 to the program counter.
		cpu.pc += 4

		// 3. Decode.
		// 4. Execute.
		cpu.Execute(inst)
	}
	cpu.DumpRegisters()
}
