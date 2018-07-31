package amd64paging

import "unsafe"

// PageDirectoryPointerTable is the L3 amd64 page table
type PageDirectoryPointerTable BaseTable

// Map recursively adds a mapping from the given virtual address to the given
//  physical address, with the given access flags
func (pdpt *PageDirectoryPointerTable) Map(v Address, p Address, flags uint64) error {
	idx := v.getL3Index()

	// check that the mapping exists in this table
	if err := (*BaseTable)(pdpt).GetMapping(idx); err != nil {
		return err
	}

	// step into the page table mapped for the given virtual address
	pd, err := pdpt.Step(v, pageTableFlags)
	if err != nil {
		return err
	}

	// add mapping at next paging level
	return pd.Map(v, p, flags)
}

// ptrToPageDirectory is a testing hook to redirect recursive page table traversal
var ptrToPageDirectory = func(ptr unsafe.Pointer) *PageDirectory {
	return (*PageDirectory)(ptr)
}

// Step returns a virtual address the L2 table at the given index
func (pdpt *PageDirectoryPointerTable) Step(v Address, flags uint64) (*PageDirectory, error) {
	// ignore the supplied base and use a fresh recurse base
	pd, err := (*BaseTable)(pdpt).Step(v.getL3Index(), flags)
	return ptrToPageDirectory(pd), err
}

func (pdpt *PageDirectoryPointerTable) toAddress() Address {
	return Address(uintptr(unsafe.Pointer(pdpt)))
}
