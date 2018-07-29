package amd64paging

const (
	// NumEntriesL4 is the number of entries in the L4 page table
	NumEntriesL4 = PageSize / EntrySize
)

// PML4 is the Go representation of the amd64 L4 page directory table
type PML4 struct {
	entries [NumEntriesL4]PML4Entry
}

// PML4Entry is an entry in the PML4 table, pointing to an L3 table
type PML4Entry uint64

// Bootstrap creates a recursize mapping to allow self-referential page tables
func (pdt *PML4) Bootstrap(p Address) {

}
