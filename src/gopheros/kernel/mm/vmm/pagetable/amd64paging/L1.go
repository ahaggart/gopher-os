package amd64paging

import "unsafe"

// Table is the amd64 L1 page table
type Table BaseTable

// Map points a virtual address to a physical address
func (t *Table) Map(v Address, p Address, flags uint64) {
	entry := t.getTableEntry(v)
	entry.SetAddress(p)
	entry.SetFlags(flags)
}

// walk returns the physical address for a given virtual address, mirroring
//  the behavior of the CPU page walk for an amd64 processor
func (t *Table) walk(v Address) (p Address) {
	offset := v.getOffset()
	te := t.getTableEntry(v)
	return createAddress(*te, offset)
}

// getTableEntry returns the TableEntry given a virtual address
//  by extracting the L1 index bits from the virtual address
func (t *Table) getTableEntry(v Address) *TableEntry {
	return &t[((v >> L1LSB) & (1<<L1Size - 1))]
}

func createAddress(te TableEntry, offset Address) Address {
	frame := Address(te & ^TableEntry(0xFFF))
	return frame | (offset & 0xFFF)
}

// Step returns a virtual address the L3 table at the given index
func (t *Table) Step(base Address, idx, flags uint64) (*TableEntry, error) {
	// ignore the supplied base and use a fresh recurse base
	te, err := (*BaseTable)(t).Step(base, idx, flags)
	return (*TableEntry)(te), err
}

func (t *Table) toAddress() Address {
	return Address(uintptr(unsafe.Pointer(t)))
}
