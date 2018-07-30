package amd64paging

// Convenience constants for building recursive addresses
const (
	recurseBase = ^Address((1<<OffsetSize - 1))
)

// PML4 is the Go representation of the amd64 L4 page directory table
type PML4 BaseTable

// Bootstrap creates a recursive mapping to allow self-referential page tables
func (pm *PML4) Bootstrap(p Address) {
	entry := &pm[SelfAddressIndex]
	entry.SetAddress(p)
}

// Phys returns the physical address of this page table
func (pm *PML4) Phys() (p Address) {
	return pm[SelfAddressIndex].Address()
}

// Step returns a virtual address the L3 table at the given index
func (pm *PML4) Step(base Address, idx, flags uint64) (*PageDirectoryPointerTable, error) {
	// ignore the supplied base and use a fresh recurse base
	pdpt, err := (*BaseTable)(pm).Step(recurseBase, idx, flags)
	return (*PageDirectoryPointerTable)(pdpt), err
}
