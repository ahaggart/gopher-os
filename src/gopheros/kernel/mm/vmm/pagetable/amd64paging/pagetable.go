// Package amd64paging implements page table operations for amd64 hardware
package amd64paging

import "gopheros/kernel/mm/vmm/pagetable"

const (
	// PageSize is the page size in bytes for amd64 processors
	PageSize = 4096

	// EntrySize is the size of a table entry in bytes
	EntrySize = 8

	// L4 table index bits and size
	L4MSB  = 47
	L4LSB  = 39
	L4Size = L4MSB - L4LSB + 1
	L4Mask = uint64(tableIndexMask << L4LSB)

	L3MSB  = 38
	L3LSB  = 30
	L3Size = L3MSB - L3LSB + 1
	L3Mask = uint64(tableIndexMask << L3LSB)

	L2MSB  = 29
	L2LSB  = 21
	L2Size = L2MSB - L2LSB + 1
	L2Mask = uint64(tableIndexMask << L2LSB)

	L1MSB  = 20
	L1LSB  = 12
	L1Size = L1MSB - L1LSB + 1
	L1Mask = uint64(tableIndexMask << L1LSB)

	OffsetMSB  = 11
	OffsetSize = OffsetMSB + 1

	AddressSize = 48

	tableIndexMask = 0x1FF //9 bits

	NumEntries       = PageSize / EntrySize
	SelfAddressIndex = NumEntries - 1
)

// LoadCR3 is an assembly implementation for loading a page table into the
//  processor
func LoadCR3(p Address)

// Load sets a physical address as the current page table to use
func Load(p Address) {
	LoadCR3(p)
}

// PageTable satisfies pagetable.PageTable
type PageTable struct {
	pm *PML4
}

var _ pagetable.PageTable = &PageTable{}

// Walk traverses this page table using recusive mapping. This page table must
//  be the page table in use, or Walk will fail
func (pt *PageTable) Walk() {

}

// Use loads this page table into the processor
func (pt *PageTable) Use() {
	Load(pt.pm.Self())
}

//Address wraps a 64bit address
type Address uint64

func (addr Address) withBits(insertion, mask Address) Address {
	insertion &= mask
	addr &= ^mask
	return addr | insertion
}

func (addr Address) getOffset() Address {
	return addr & 0xFFF
}

func (addr Address) getL4Index() uint64 {
	return uint64((addr >> L4LSB) & tableIndexMask)
}

func (addr Address) withL4Index(idx uint64) Address {
	return addr.withBits(Address(idx), Address(L4Mask))
}

func (addr Address) getL3Index() uint64 {
	return uint64((addr >> L3LSB) & tableIndexMask)
}

func (addr Address) withL3Index(idx uint64) Address {
	return addr.withBits(Address(idx), Address(L3Mask))
}

func (addr Address) getL2Index() uint64 {
	return uint64((addr >> L2LSB) & tableIndexMask)
}

func (addr Address) withL2Index(idx uint64) Address {
	return addr.withBits(Address(idx), Address(L2Mask))
}

func (addr Address) getL1Index() uint64 {
	return uint64((addr >> L1LSB) & tableIndexMask)
}

func (addr Address) withL1Index(idx uint64) Address {
	return addr.withBits(Address(idx), Address(L1Mask))
}

// BaseTable is the base page table type
type BaseTable [NumEntries]TableEntry

// TableEntry is the base type for amd64 page table entries
type TableEntry uint64

const (
	tableEntryAddressMask = TableEntry((1<<AddressSize - 1) << OffsetSize)
)

// Address returns the embedded 48-bit address in a table entry
func (te *TableEntry) Address() Address {
	// mask off the offset / flag bits and the high bits
	return Address(*te & tableEntryAddressMask)
}

// SetAddress returns the embedded 48-bit address in a table entry
func (te *TableEntry) SetAddress(addr Address) {
	*te &= tableEntryAddressMask
	*te |= TableEntry(addr) & tableEntryAddressMask
}

// SetFlags sets the flags for a table entry
func (te *TableEntry) SetFlags(flags L1Flags) {
	*te &= ^TableEntry(0xFFF)
	*te |= TableEntry(flags)
}
