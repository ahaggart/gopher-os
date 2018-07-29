package pagetable

// PageTable is the interface to a hardware-specific page table implementation,
// 	and handles setup and maintenance of page tables in memory
type PageTable interface {
	// Walk traverses the page table until a translation is found or an error
	//  occurs
	Walk()

	// Use causes this PageTable to be used for virtual address translations
	Use()
}
