package riscv

import (
	"fmt"
)

type Device interface {
	StartAddr() uint32
	EndAddr() uint32
	Read(addr, size uint32) uint32
	Write(addr, size, value uint32)
}

// Bus represents a system bus which is a single computer bus that connects the major components
// of a computer system, combining the functions of a data bus to carry information, an address bus
// to determine where it should be sent or read from, and a control bus to determine its operation.
//
// This Bus struct contains memory-mapped I/O as into memory bus.
//
// ref: https://fmash16.github.io/content/posts/riscv-emulator-in-c.html#:~:text=Writing%20a%20DRAM%20struct
//
//                 +---------------+
//                 | Address space |
//                 |   +-------+   |
//                 |   |  ROM  |   |
//                 |   +-------+   |
// +-------+address|   |       |   |
// |       |------>|   |  RAM  |   |
// |  CPU  |  bus  |   |       |   |
// |       |<----->|   +-------+   |
// +-------+ data  |   |       |   |
//                 |   |  I/O  |   |
//                 |   +-------+   |
//                 +---------------+
//
// Memory-mapped I/O (MMIO) is two complementary methods of performing input/output (I/O) between
// the central processing unit (CPU) and peripheral devices in a computer.
// Memory-mapped I/O uses the same address space to address both memory and I/O devices.
//
// The dram region address is started from 0x80000000.
// qemu: https://github.com/qemu/qemu/blob/5e9d14f2bea6df89c0675df953f9c839560d2266/hw/riscv/virt.c#L61
//
// static const MemMapEntry virt_memmap[] = {
//     [VIRT_DEBUG] =       {        0x0,         0x100 },
//     [VIRT_MROM] =        {     0x1000,        0xf000 },
//     [VIRT_TEST] =        {   0x100000,        0x1000 },
//     [VIRT_RTC] =         {   0x101000,        0x1000 },
//     [VIRT_CLINT] =       {  0x2000000,       0x10000 },
//     [VIRT_ACLINT_SSWI] = {  0x2F00000,        0x4000 },
//     [VIRT_PCIE_PIO] =    {  0x3000000,       0x10000 },
//     [VIRT_PLIC] =        {  0xc000000, VIRT_PLIC_SIZE(VIRT_CPUS_MAX * 2) },
//     [VIRT_UART0] =       { 0x10000000,         0x100 },
//     [VIRT_VIRTIO] =      { 0x10001000,        0x1000 },
//     [VIRT_FW_CFG] =      { 0x10100000,          0x18 },
//     [VIRT_FLASH] =       { 0x20000000,     0x4000000 },
//     [VIRT_PCIE_ECAM] =   { 0x30000000,    0x10000000 },
//     [VIRT_PCIE_MMIO] =   { 0x40000000,    0x40000000 },
//     [VIRT_DRAM] =        { 0x80000000,           0x0 },
// };
type Bus struct {
	devices []Device
}

func NewBus(devices ...Device) *Bus {
	return &Bus{devices: devices}
}

func (b *Bus) findDevice(addr, size uint32) (Device, error) {
	for _, dev := range b.devices {
		useAddrLen := addr + size - 1
		if dev.StartAddr() <= addr && useAddrLen <= dev.EndAddr() {
			return dev, nil
		}
	}
	return nil, fmt.Errorf("device is not found (addr: %08x, size: %d)", addr, size)
}

func (b *Bus) Read(addr, size uint32) (uint32, error) {
	device, err := b.findDevice(addr, size)
	if err != nil {
		return 0, err
	}
	return device.Read(addr, size), nil
}

func (b *Bus) Write(addr, size, value uint32) error {
	device, err := b.findDevice(addr, size)
	if err != nil {
		return err
	}
	device.Write(addr, size, value)
	return nil
}
