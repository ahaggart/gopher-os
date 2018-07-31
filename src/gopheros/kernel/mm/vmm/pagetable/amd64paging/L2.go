package amd64paging

import "unsafe"

// PageDirectory is the amd64 L2 page table
type PageDirectory BaseTable

// ptrToTable is a testing hook for redirecting recursive page table traversal
var ptrToTable = func(ptr unsafe.Pointer) *Table {
	return (*Table)(ptr)
}

// Map recursively adds a mapping from the given virtual address to the given
//  physical address, with the given access flags
func (pd *PageDirectory) Map(v Address, p Address, flags uint64) error {
	idx := v.getL2Index()

	// check that the mapping exists in this table
	if err := (*BaseTable)(pd).GetMapping(idx, flags); err != nil {
		return err
	}

	// step into the page table mapped for the given virtual address
	t, err := pd.Step(v, pageTableFlags)
	if err != nil {
		return err
	}

	// add mapping at next paging level
	return t.Map(v, p, flags)
}

// Step returns a virtual address the L1 table at the given index
func (pd *PageDirectory) Step(v Address, flags uint64) (*Table, error) {
	// ignore the supplied base and use a fresh recurse base
	t, err := (*BaseTable)(pd).Step(v.getL2Index(), flags)
	return ptrToTable(t), err
}

func (pd *PageDirectory) toAddress() Address {
	return Address(uintptr(unsafe.Pointer(pd)))
}
