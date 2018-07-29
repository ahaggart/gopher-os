package amd64paging

// L1Flags is the Go representation of the permissions flags for an L1 PTE
type L1Flags uint8

// Table is the amd64 L1 page table
type Table struct {
	entries [PageSize / 8]TableEntry
}

func BootstrapL1(t *Table, p Address) {

}

// Map points a virtual address to a physical address
func (t *Table) Map(v Address, p Address, flags L1Flags) {
	entry := t.getTableEntry(v)
	entry.setAddress(p)
	entry.setFlags(flags)
}

// walk returns the physical address for a given virtual address, mirroring
//  the behavior of the CPU page walk for an amd64 processor
func (t *Table) walk(v Address) (p Address) {
	offset := v.getOffset()
	te := t.getTableEntry(v)
	return te.createAddress(offset)
}

// getTableEntry returns the TableEntry given a virtual address
//  by extracting the L1 index bits from the virtual address
func (t *Table) getTableEntry(v Address) *TableEntry {
	return &t.entries[((v >> L1LSB) & (1<<L1Size - 1))]
}

//TableEntry is a page table entry, pointing to the physical frame for a
//  virtual address
type TableEntry uint64

func (te *TableEntry) setAddress(p Address) {
	// clear off all address bits
	*te &= 0xFFF

	//clear offset bits
	p &= 0xFFF

	*te |= TableEntry(p)
}

func (te *TableEntry) setFlags(flags L1Flags) {
	*te &= ^TableEntry(0xFF)
	*te |= TableEntry(flags)
}

func (te *TableEntry) createAddress(offset Address) Address {
	frame := Address(*te & ^TableEntry(0xFFF))
	return frame | (offset & 0xFFF)
}
