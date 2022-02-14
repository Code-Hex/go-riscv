package riscv

// DRAM (Dyanmic random access memory) is our memory that contains
// all the instructions to be executed and the data.
type DRAM struct {
	mem []byte
}

var _ Device = (*DRAM)(nil)

func NewDRAM(size int) *DRAM {
	return &DRAM{
		mem: make([]byte, size),
	}
}

// Read reads any values from dram.
// size specify the bit size. i.e 8, 16, 32, 64 bit...
func (d *DRAM) Read(addr, size uint32) uint32 {
	var result uint32
	for i := uint32(0); i < size; i++ {
		idx := int(addr + i)
		result |= uint32(d.mem[idx] << (8 * i))
	}
	return result
}

// Write writes any values to dram.
// size specify the bit size. i.e 8, 16, 32, 64 bit...
func (d *DRAM) Write(addr, size, value uint32) {
	for i := uint32(0); i < size; i++ {
		idx := int(addr + i)
		d.mem[idx] = byte(value >> ((8 * i) & 0xff)) // 0xff ~ 8 bit masking (byte type == uint8)
	}
}

// StartAddr represents start address for DRAM.
func (d *DRAM) StartAddr() uint32 { return 0x80000000 }

// EndAddr represents end of address for DRAM.
func (d *DRAM) EndAddr() uint32 { return d.StartAddr() + uint32(len(d.mem)) }
