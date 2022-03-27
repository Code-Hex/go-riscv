add-addi.bin: testdata/add-addi/add-addi.s
	riscv64-unknown-elf-gcc -march=rv32i -mabi=ilp32 -Wl,-Ttext=0x0 -nostdlib -O0 -o testdata/add-addi/add-addi testdata/add-addi/add-addi.s
	riscv64-unknown-elf-objcopy -O binary testdata/add-addi/add-addi testdata/add-addi/add-addi.bin
	rm testdata/add-addi/add-addi

lb.bin: testdata/lb/lb.s
	riscv64-unknown-elf-gcc -march=rv32i -mabi=ilp32 -Wl,-Ttext=0x0 -nostdlib -O0 -o testdata/lb/lb testdata/lb/lb.s
	riscv64-unknown-elf-objcopy -O binary testdata/lb/lb testdata/lb/lb.bin
	rm testdata/lb/lb

clean:
	rm -f testdata/add-addi/add-addi
	rm -f testdata/add-addi/add-addi.bin