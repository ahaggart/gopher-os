package pagetable

// PageTable is the interface to a hardware-specific page table implementation,
// 	and handles setup and maintenance of page tables in memory
type PageTable interface {
	// Walk traverses the page table until a translation is found or an error
	//  occurs
	Walk(v, flags uint64) (uint64, error)

	// Use causes this PageTable to be used for virtual address translations
	Use()
}

// NotMappedError is an error returned when requested page is not mapped
type NotMappedError struct{}

func (nme NotMappedError) Error() string {
	return "Error: Page not mapped"
}

// PermissionsError is an error returned when a page is requested with the wrong
//  permission flags
type PermissionsError uint8

func (pe PermissionsError) Error() string {
	return "Error: Improper page access permissions"
}
