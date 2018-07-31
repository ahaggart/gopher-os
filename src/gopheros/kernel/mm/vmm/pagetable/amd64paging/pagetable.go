// Package amd64paging implements page table operations for amd64 hardware
package amd64paging

import (
	"gopheros/kernel/mm"
	"gopheros/kernel/mm/vmm/pagetable"
	"unsafe"
)

const (
	// PageSize is the page size in bytes for amd64 processors
	PageSize = 4096

	// EntrySize is the size of a table entry in bytes
	EntrySize = 8
)

// convenience constants for working with page table indices
const (
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

	tableIndexMask        = (1 << tableIndexSize) - 1
	tableEntryAddressMask = TableEntry((1<<AddressSize - 1) << OffsetSize)

	NumEntries       = PageSize / EntrySize
	SelfAddressIndex = NumEntries - 1

	NilMapping = TableEntry(0)

	cr3Mask = ^uint64((1<<48 - 1) << 12)

	pageTableFlags = 0xFFF
)

// PageTable satisfies pagetable.PageTable
type PageTable struct {
	pm *PML4
}

var _ pagetable.PageTable = &PageTable{}

// Walk traverses this page table using recusive mapping. This page table must
//  be the page table in use, or Walk will fail
func (pt *PageTable) Walk(v, flags uint64) (uint64, error) {
	addr := Address(v)
	if pdpt, err := pt.pm.Step(addr, flags); err != nil {
		// step through the pml4 table to get the page directory pointer table
		return 0, err
	} else if pd, err := pdpt.Step(addr, flags); err != nil {
		// step through the pdpt to get the page directory
		return 0, err
	} else if t, err := pd.Step(addr, flags); err != nil {
		// step through the page directory to get the page table
		return 0, err
	} else if te, err := t.Step(addr, flags); err != nil {
		// step through the page table to get a page table entry
		return 0, err
	} else {
		// extract the physical address from the page table entry
		return uint64(te.Address().withOffset(addr.offset())), nil
	}
}

// Use loads this page table into the processor
func (pt *PageTable) Use() {
	Load(pt.pm.Phys())
}

// Load sets a physical address as the current page table to use
func Load(p Address) {
	LoadCR3(p)
}

// LoadCR3 is an assembly implementation for loading a page table into the
//  processor
func LoadCR3(p Address)

// BaseTable is the base page table type
type BaseTable [NumEntries]TableEntry

// Step returns a virtual address the L3 table at the given index
func (bt *BaseTable) Step(idx, flags uint64) (unsafe.Pointer, error) {
	if te := bt[idx]; te == NilMapping {
		// check if the mapping exists
		return nil, pagetable.NotMappedError{}

	} else if !te.HasFlags(flags) {
		// check if the mapping has the appropriate privileges
		return nil, makePermissionsError(te.Flags())

	}
	idx = (idx & tableIndexMask) << L1LSB
	addr := (bt.toAddress() << tableIndexSize) | Address(idx)

	return addr.toPointer(), nil
}

// GetMapping checks for a mapping at the given index, and adds on if it
//  does not exist
func (bt *BaseTable) GetMapping(idx, flags uint64) error {
	if bt[idx] == NilMapping {
		frame, err := mm.AllocFrame()
		if err != nil {
			return err
		}
		bt[idx] = newTableEntry(
			Address(frame),
			flags,
		)
	}
	return nil
}

func (bt *BaseTable) toAddress() Address {
	return Address(uintptr(unsafe.Pointer(bt)))
}

// TableEntry is the base type for amd64 page table entries
type TableEntry uint64

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

// for now just mask off the lower 8 bits
func makePermissionsError(flags uint64) pagetable.PermissionsError {
	return pagetable.PermissionsError(flags & 0xFFF)
}

func newTableEntry(addr Address, flags uint64) (entry TableEntry) {
	te := &entry
	te.SetAddress(addr)
	te.SetFlags(flags)
	return
}
