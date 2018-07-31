package amd64paging

import (
	"testing"
	"unsafe"
)

func TestL4Phys(t *testing.T) {
	pm := &PML4{}

	testAddress := Address(0x0000123456789876000)

	pm.Bootstrap(testAddress)

	physAddress := pm.Phys()

	if testAddress != physAddress {
		t.Fatalf(
			"PML4 does not return its physical address pointer: %x != %x",
			testAddress,
			physAddress,
		)
	}

}

func TestL4Map(t *testing.T) {

}

func TestL4Step(t *testing.T) {
	// hook the casting functions to prevent a segfault
	var (
		origPtrToPDPT          = ptrToPDPT
		origPtrToPageDirectory = ptrToPageDirectory
		origPtrToTable         = ptrToTable
		origPtrToTableEntry    = ptrToTableEntry
	)
	defer func() {
		ptrToPDPT = origPtrToPDPT
		ptrToPageDirectory = origPtrToPageDirectory
		ptrToTable = origPtrToTable
		ptrToTableEntry = origPtrToTableEntry
	}()

	pm := &PML4{}

	testAddress := Address(0x0000123456789876000)

	pm.Bootstrap(testAddress)

	pmRecursive, err := pm.Step(SelfAddressIndex<<L4LSB, 0)

	if err != nil {
		t.Fatalf("PML4 Step returned error: %v", err)
	}

	recursivePtr := Address(uintptr(unsafe.Pointer(pmRecursive)))

	if recursivePtr != testAddress {
		t.Fatalf(
			"PML4 Recursive address is incorrect: %x != %x",
			testAddress,
			recursivePtr,
		)
	}

}
