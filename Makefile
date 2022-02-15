add-addi.bin: testdata/add-addi.s
	riscv64-unknown-elf-gcc -march=rv32i -mabi=ilp32 -Wl,-Ttext=0x0 -nostdlib -O0 -o testdata/add-addi testdata/add-addi.s
	riscv64-unknown-elf-objcopy -O binary testdata/add-addi testdata/add-addi.bin
	rm testdata/add-addi

clean:
	rm -f testdata/add-addi
	rm -f testdata/add-addi.bin