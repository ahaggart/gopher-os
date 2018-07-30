package amd64paging

import "unsafe"

// PageDirectoryPointerTable is the L3 amd64 page table
type PageDirectoryPointerTable BaseTable

// Step returns a virtual address the L2 table at the given index
func (pdpt *PageDirectoryPointerTable) Step(base Address, idx, flags uint64) (*PageDirectory, error) {
	// ignore the supplied base and use a fresh recurse base
	pd, err := (*BaseTable)(pdpt).Step(base, idx, flags)
	return (*PageDirectory)(pd), err
}

func (pdpt *PageDirectoryPointerTable) toAddress() Address {
	return Address(uintptr(unsafe.Pointer(pdpt)))
}
