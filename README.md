# RISC-V emulator


## memo

- [Writing a RISC-V Emulator in Rust](https://book.rvemu.app/)
- [HackerNews - Writing a RISC-V Emulator from Scratch](https://news.ycombinator.com/item?id=23033517)
- [The Adventures of OS: Making a RISC-V Operating System using Rust](http://osblog.stephenmarz.com/)
- [RISC-V ELF psABI](https://github.com/riscv-non-isa/riscv-elf-psabi-doc)
- [riscv-emu (Go)](https://github.com/LMMilewski/riscv-emu)
- [RISC-V Instruction Formats](https://inst.eecs.berkeley.edu/~cs61c/resources/su18_lec/Lecture7.pdf)
- [github.com/johnwinans/rvalp - RISC-V Assemly Language Programming](https://github.com/johnwinans/rvalp)
- [Writing a simple RISC-V emulator in plain C](https://fmash16.github.io/content/posts/riscv-emulator-in-c.html)
- Japanese
  - [RISC-Vを実装してみる](https://kamiyaowl.github.io/presentation/pdf/lets-impl-rv32i.pdf)

## Install tool-chain

- macOS should install using homebrew.
  - see: https://github.com/riscv-software-src/homebrew-riscv

### Instructions

RISC-V has been defined 6-type instructions

- R-format
  - instructions using 3 register inputs
  - add, xor, mul - arithmetic/logical ops
- I-format
  - instructions with immediates, loads
  - addi, lw, jalr, slli
- S-format
  - store instructions
  - sw, sb
- U-format
  - instructions with upper immediates
  - lui, auipc - upper immediate is 20-bits
- SB-format
  - branch instructions
  - beq, bge
- UJ-format
  - jump instructions
  - jal
