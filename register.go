package riscv

// https://riscv.org/wp-content/uploads/2015/01/riscv-calling.pdf
var xregsABINames = []string{
	"x0 (zero)",
	"x1 (ra)", // Return address
	"x2 (sp)", // Stack pointer
	"x3 (gp)",
	"x4 (tp)",
	"x5 (t0)",
	"x6 (t1)",
	"x7 (t2)",
	"x8 (s0/fp)",
	"x9 (s1)",
	"x10 (a0)",
	"x11 (a1)",
	"x12 (a2)",
	"x13 (a3)",
	"x14 (a4)",
	"x15 (a5)",
	"x16 (a6)",
	"x17 (a7)",
	"x18 (s2)",
	"x19 (s3)",
	"x20 (s4)",
	"x21 (s5)",
	"x22 (s6)",
	"x23 (s7)",
	"x24 (s8)",
	"x25 (s9)",
	"x26 (s10)",
	"x27 (s11)",
	"x28 (t3)",
	"x29 (t4)",
	"x30 (t5)",
	"x31 (t6)",
}