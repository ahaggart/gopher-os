package amd64paging

import (
	"gopheros/kernel"
	"gopheros/kernel/mm"
	"testing"
)

func TestBaseTableMap(t *testing.T) {
	addr := mm.Frame(uintptr(0x0000123456789000))
	flags := uint64(0xFFF)

	defer mm.SetFrameAllocator(nil)

	mm.SetFrameAllocator(func() (mm.Frame, *kernel.Error) {
		return addr, nil
	})

	bt := &BaseTable{}

	if err := bt.GetMapping(0, flags); err != nil {
		t.Fatalf("BaseTable: Mapping for index 0 failed: %v", err)
	}

	te := bt[0]

	if test, real := te.Address(), Address(addr); test != real {
		t.Fatalf(
			"BaseTable: Mapped address does not match: %x != %x",
			test, real,
		)
	}

	if test := te.Flags(); test != flags {
		t.Fatalf(
			"BaseTable: Mapped flags do not match: %x != %x",
			test, flags,
		)
	}

	mm.SetFrameAllocator(func() (mm.Frame, *kernel.Error) {
		return mm.Frame(uintptr(0x0000987654321000)), nil
	})

	if err := bt.GetMapping(0, flags); err != nil {
		t.Fatalf("BaseTable: Mapping for index 0 failed: %v", err)
	}

	te = bt[0]

	if test, real := te.Address(), Address(addr); test != real {
		t.Fatalf(
			"BaseTable: Mapped address does not match: %x != %x",
			test, real,
		)
	}

	if test := te.Flags(); test != flags {
		t.Fatalf(
			"BaseTable: Mapped flags do not match: %x != %x",
			test, flags,
		)
	}

}
