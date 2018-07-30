package amd64paging

import "unsafe"

// PageDirectoryPointerTable is the L3 amd64 page table
type PageDirectoryPointerTable BaseTable

// ptrToPageDirectory is a testing hook to redirect recursive page table traversal
var ptrToPageDirectory = func(ptr unsafe.Pointer) *PageDirectory {
	return (*PageDirectory)(ptr)
}

// Step returns a virtual address the L2 table at the given index
func (pdpt *PageDirectoryPointerTable) Step(idx, flags uint64) (*PageDirectory, error) {
	// ignore the supplied base and use a fresh recurse base
	pd, err := (*BaseTable)(pdpt).Step(idx, flags)
	return ptrToPageDirectory(pd), err
}

func (pdpt *PageDirectoryPointerTable) toAddress() Address {
	return Address(uintptr(unsafe.Pointer(pdpt)))
}
