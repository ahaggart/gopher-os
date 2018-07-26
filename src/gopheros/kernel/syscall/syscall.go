package syscall

import (
	"gopheros/kernel/gate"
	"gopheros/kernel/kfmt"
)

//InitSyscalls installs a syscall interrupt handler on the IDT
func InitSyscalls() {
	gate.HandleInterrupt(
		gate.Syscall,
		0,
		handleSyscall,
	)
}

//Interrupt executes a syscall instruction
func Interrupt()

func handleSyscall(regs *gate.Registers) {
	kfmt.Printf("handling syscall\n")
}
