package amd64paging

import (
	"unsafe"
)

// Convenience constants for building recursive addresses
const (
	recurseBase = ^Address((1<<OffsetSize - 1))
)

// ptrToPDPT casts an unsafe pointer to an L3 page table pointer
//  we split the code up this way to allow tests to redirect a pointer
//  and prevent segfaults
var ptrToPDPT = func(ptr unsafe.Pointer) *PageDirectoryPointerTable {
	return (*PageDirectoryPointerTable)(ptr)
}

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

// Map recursively adds a mapping from the given virtual address to the given
//  physical address, with the given access flags
func (pm *PML4) Map(v Address, p Address, flags uint64) error {
	idx := v.getL4Index()

	// check that the mapping exists in this table
	if err := (*BaseTable)(pm).GetMapping(idx, flags); err != nil {
		return err
	}

	// step into the page table mapped for the given virtual address
	pdpt, err := pm.Step(v, pageTableFlags)
	if err != nil {
		return err
	}

	// add mapping at next paging level
	return pdpt.Map(v, p, flags)
}

// Step returns a virtual address the L3 table at the given index
func (pm *PML4) Step(v Address, flags uint64) (*PageDirectoryPointerTable, error) {
	// ignore the supplied base and use a fresh recurse base
	pdpt, err := (*BaseTable)(pm).Step(v.getL4Index(), flags)
	return ptrToPDPT(pdpt), err
}
