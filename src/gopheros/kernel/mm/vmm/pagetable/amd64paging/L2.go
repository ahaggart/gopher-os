package amd64paging

import "unsafe"

// PageDirectory is the amd64 L2 page table
type PageDirectory BaseTable

// Step returns a virtual address the L1 table at the given index
func (pd *PageDirectory) Step(base Address, idx, flags uint64) (*Table, error) {
	// ignore the supplied base and use a fresh recurse base
	t, err := (*BaseTable)(pd).Step(base, idx, flags)
	return (*Table)(t), err
}

func (pd *PageDirectory) toAddress() Address {
	return Address(uintptr(unsafe.Pointer(pd)))
}
