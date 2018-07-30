package amd64paging

// L1Flags is the Go representation of the permissions flags for an L1 PTE
type L1Flags uint8

// Table is the amd64 L1 page table
type Table BaseTable

// Map points a virtual address to a physical address
func (t *Table) Map(v Address, p Address, flags L1Flags) {
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
