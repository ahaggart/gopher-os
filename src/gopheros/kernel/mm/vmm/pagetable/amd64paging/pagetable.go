// Package amd64paging implements page table operations for amd64 hardware
package amd64paging

const (
	// PageSize is the page size for amd64 processors
	PageSize = 4096

	// L4 table index bits and size
	L4MSB  = 47
	L4LSB  = 39
	L4Size = L4MSB - L4LSB

	L3MSB  = 38
	L3LSB  = 30
	L3Size = L3MSB - L3LSB

	L2MSB  = 29
	L2LSB  = 21
	L2Size = L2MSB - L2LSB

	L1MSB  = 20
	L1LSB  = 12
	L1Size = L1MSB - L1LSB

	OffsetMSB  = 11
	OffsetSize = OffsetMSB
)

// PML4 is the Go representation of the amd64 L4 page directory table
type PML4 struct {
}

// PML4Entry is an entry in the PML4 table, pointing to an L3 table
type PML4Entry uint64

// PageDirectoryPointerTable is the L3 amd64 page table
type PageDirectoryPointerTable struct {
}

// PDPTEntry is an entry in the PDPT table, pointing to an L2 table
type PDPTEntry uint64

// PageDirectory is the amd64 L2 page table
type PageDirectory struct {
}

// PDEntry is an entry in the Page-Directory table, pointing to an L1 table
type PDEntry uint64

//Address wraps a 64bit address
type Address uint64

func (addr Address) getOffset() Address {
	return addr & 0xFF
}
