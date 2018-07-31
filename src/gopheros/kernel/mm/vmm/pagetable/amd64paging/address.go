package amd64paging

import "unsafe"

//Address wraps a 64bit address
type Address uint64

func (addr Address) withBits(insertion, mask Address) Address {
	insertion &= mask
	addr &= ^mask
	return addr | insertion
}

// offset returns the offset bits of the address
func (addr Address) offset() Address {
	return addr & 0xFFF
}

func (addr Address) withOffset(offset Address) Address {
	return addr.withBits(offset, 0xFFF)
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
