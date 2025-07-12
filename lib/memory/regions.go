package memory

import (
	"runtime/volatile"
	"unsafe"
)

// MemoryRegion provides safe access to a memory region
type MemoryRegion struct {
	base uintptr
	size uintptr
}

// NewMemoryRegion creates a new memory region with bounds checking
func NewMemoryRegion(base, size uintptr) *MemoryRegion {
	return &MemoryRegion{
		base: base,
		size: size,
	}
}

// InBounds checks if an offset is within the memory region bounds
func (r *MemoryRegion) InBounds(offset uintptr) bool {
	return offset < r.size
}

// Size returns the size of the memory region
func (r *MemoryRegion) Size() uintptr {
	return r.size
}

// Base returns the base address of the memory region
func (r *MemoryRegion) Base() uintptr {
	return r.base
}

// Read16 reads a 16-bit value from the memory region with bounds checking
func (r *MemoryRegion) Read16(offset uintptr) uint16 {
	if !r.InBounds(offset + 1) {
		return 0
	}
	ptr := (*volatile.Register16)(unsafe.Pointer(r.base + offset))
	return ptr.Get()
}

// Write16 writes a 16-bit value to the memory region with bounds checking
func (r *MemoryRegion) Write16(offset uintptr, value uint16) {
	if !r.InBounds(offset + 1) {
		return
	}
	ptr := (*volatile.Register16)(unsafe.Pointer(r.base + offset))
	ptr.Set(value)
}

// Read32 reads a 32-bit value from the memory region with bounds checking
func (r *MemoryRegion) Read32(offset uintptr) uint32 {
	if !r.InBounds(offset + 3) {
		return 0
	}
	ptr := (*volatile.Register32)(unsafe.Pointer(r.base + offset))
	return ptr.Get()
}

// Write32 writes a 32-bit value to the memory region with bounds checking
func (r *MemoryRegion) Write32(offset uintptr, value uint32) {
	if !r.InBounds(offset + 3) {
		return
	}
	ptr := (*volatile.Register32)(unsafe.Pointer(r.base + offset))
	ptr.Set(value)
}

// ReadColor reads a Color (16-bit) from the memory region
func (r *MemoryRegion) ReadColor(offset uintptr) Color {
	return Color(r.Read16(offset))
}

// WriteColor writes a Color (16-bit) to the memory region
func (r *MemoryRegion) WriteColor(offset uintptr, color Color) {
	r.Write16(offset, uint16(color))
}

// Clear fills the memory region with zeros
func (r *MemoryRegion) Clear() {
	for offset := uintptr(0); offset < r.size; offset += 4 {
		if offset+3 < r.size {
			r.Write32(offset, 0)
		} else {
			r.Write16(offset, 0)
		}
	}
}

// Fill16 fills the memory region with a 16-bit value
func (r *MemoryRegion) Fill16(value uint16) {
	for offset := uintptr(0); offset < r.size; offset += 2 {
		r.Write16(offset, value)
	}
}

// FillColor fills the memory region with a color
func (r *MemoryRegion) FillColor(color Color) {
	r.Fill16(uint16(color))
}

// Global memory region instances
var (
	VRAM       = NewMemoryRegion(VRAM_BASE, VRAM_SIZE)
	OAM        = NewMemoryRegion(OAM_BASE, OAM_SIZE)
	PaletteRAM = NewMemoryRegion(PALETTE_BASE, PALETTE_SIZE)
)

// GetVRAM returns the VRAM memory region
func GetVRAM() *MemoryRegion {
	return VRAM
}

// GetOAM returns the OAM memory region
func GetOAM() *MemoryRegion {
	return OAM
}

// GetPaletteRAM returns the Palette RAM memory region
func GetPaletteRAM() *MemoryRegion {
	return PaletteRAM
}