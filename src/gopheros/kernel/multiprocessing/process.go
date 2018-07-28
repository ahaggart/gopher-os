package multiprocessing

import (
	"gopheros/kernel/gate"
)

const (
	// ProcessRunning is the status code for a running process
	ProcessRunning = iota

	// ProcessSuspended is the status code for a suspended process
	ProcessSuspended
)

// ProcessControlBlock stores the complete state of a process
type ProcessControlBlock struct {
	// Process status info
	pid    uint64
	status uint8
	regs   gate.Registers

	// Linked Lists for process queueing
	nextProcess *ProcessControlBlock
	prevProcess *ProcessControlBlock

	// Process ancestry info
	parent   *ProcessControlBlock
	children []*ProcessControlBlock
}
