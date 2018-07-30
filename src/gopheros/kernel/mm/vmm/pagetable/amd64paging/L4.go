package amd64paging

const (
	// SelfAddressBase is the base for building recursive self-referential
	//  pointers
	SelfAddressBase = ^Address(1<<48 - 1)
)

// PML4 is the Go representation of the amd64 L4 page directory table
type PML4 BaseTable

// Bootstrap creates a recursive mapping to allow self-referential page tables
func (pm *PML4) Bootstrap(p Address) {
	entry := &pm[SelfAddressIndex]
	entry.SetAddress(p)
}

// Self returns the physical address of this page table
func (pm *PML4) Self() (p Address) {
	return pm[SelfAddressIndex].Address()
}

func (pm *PML4) getL3() *PageDirectoryPointerTable {

}

func buildSelfAddress(level int) Address {
	return SelfAddressBase.withL4Index(
		0,
	).withL3Index(
		0,
	).withL2Index(
		0,
	).withL1Index(
		0,
	)
}
