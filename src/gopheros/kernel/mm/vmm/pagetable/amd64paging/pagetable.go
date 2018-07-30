// Package amd64paging implements page table operations for amd64 hardware
package amd64paging

import (
	"gopheros/kernel/mm/vmm/pagetable"
	"unsafe"
)

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

	tableIndexSize = 9

	OffsetMSB  = 11
	OffsetSize = OffsetMSB + 1

	AddressSize = 48

	tableIndexMask = (1 << tableIndexSize) - 1

	NumEntries       = PageSize / EntrySize
	SelfAddressIndex = NumEntries - 1

	NilMapping = TableEntry(0)
)

// NotMappedError is an error returned when requested page is not mapped
type NotMappedError struct{}

func (nme NotMappedError) Error() string {
	return "Error: Page not mapped"
}

// PermissionsError is an error returned when a page is requested with the wrong
//  permission flags
type PermissionsError struct {
	entry TableEntry
	flags uint64
}

func (pe *PermissionsError) Error() string {
	return "Error: Improper page access permissions"
}

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
	Load(pt.pm.Phys())
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

func (addr Address) toPointer() unsafe.Pointer {
	return unsafe.Pointer(uintptr(addr))
}

// BaseTable is the base page table type
type BaseTable [NumEntries]TableEntry

// Step returns a virtual address the L3 table at the given index
func (bt *BaseTable) Step(idx, flags uint64) (unsafe.Pointer, error) {
	if te := bt[idx]; te == NilMapping {
		// check if the mapping exists
		return nil, NotMappedError{}

	} else if !te.HasFlags(flags) {
		// check if the mapping has the appropriate privileges
		return nil, &PermissionsError{te, flags}

	}
	idx = (idx & tableIndexMask) << L1LSB
	addr := (bt.toAddress() << tableIndexSize) | Address(idx)

	return addr.toPointer(), nil
}

func (bt *BaseTable) toAddress() Address {
	return Address(uintptr(unsafe.Pointer(bt)))
}

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
func (te *TableEntry) SetFlags(flags uint64) {
	*te &= ^TableEntry(0xFFF)
	*te |= TableEntry(flags)
}

// HasFlags checks if the entry has all the flags supplied
func (te *TableEntry) HasFlags(flags uint64) bool {
	flags &= 0xFFF
	return flags^te.Flags() == 0
}

// Flags returns the flags for the table entry
func (te *TableEntry) Flags() uint64 {
	return uint64(*te & 0xFFF)
}
