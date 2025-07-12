package vram

import (
	"github.com/matheusmortatti/gba-go/lib/memory"
)

const (
	// VRAM layout constants
	VRAM_BASE = memory.VRAM_BASE
	VRAM_SIZE = memory.VRAM_SIZE

	// Video mode specific sizes
	MODE3_FRAME_SIZE = memory.SCREEN_WIDTH * memory.SCREEN_HEIGHT * 2 // 16-bit
	MODE4_FRAME_SIZE = memory.SCREEN_WIDTH * memory.SCREEN_HEIGHT     // 8-bit
	MODE5_FRAME_SIZE = 160 * 128 * 2                                  // 16-bit small

	// Buffer addresses for bitmap modes
	FRAME0_ADDR = VRAM_BASE
	FRAME1_ADDR = VRAM_BASE + 0xA000 // Frame 1 for modes 4&5

	// Character and screen base addresses (tile modes)
	CHAR_BASE_SIZE      = 0x4000 // 16KB per character base block
	SCREEN_BASE_SIZE    = 0x800  // 2KB per screen base block
	MAX_CHAR_BLOCKS     = 4      // 0-3
	MAX_SCREEN_BLOCKS   = 32     // 0-31

	// Tile constants
	TILE_4BPP_SIZE       = 32                                 // 4 bits per pixel
	TILE_8BPP_SIZE       = 64                                 // 8 bits per pixel
	TILES_PER_CHAR_BLOCK = CHAR_BASE_SIZE / TILE_4BPP_SIZE
)

// VRAMManager provides mode-specific VRAM access
type VRAMManager struct {
	mode        int
	currentPage int // For double buffering
	frameSize   int
	width       int
	height      int
	bpp         int // bits per pixel
}

// NewVRAMManager creates a new VRAM manager for the specified video mode
func NewVRAMManager(mode int) *VRAMManager {
	vm := &VRAMManager{
		mode:        mode,
		currentPage: 0,
	}
	vm.SetMode(mode)
	return vm
}

// SetMode configures the manager for a specific video mode
func (vm *VRAMManager) SetMode(mode int) {
	vm.mode = mode
	
	switch mode {
	case memory.MODE_3:
		vm.frameSize = MODE3_FRAME_SIZE
		vm.width = memory.SCREEN_WIDTH
		vm.height = memory.SCREEN_HEIGHT
		vm.bpp = 16
	case memory.MODE_4:
		vm.frameSize = MODE4_FRAME_SIZE
		vm.width = memory.SCREEN_WIDTH
		vm.height = memory.SCREEN_HEIGHT
		vm.bpp = 8
	case memory.MODE_5:
		vm.frameSize = MODE5_FRAME_SIZE
		vm.width = 160
		vm.height = 128
		vm.bpp = 16
	default:
		// Default to mode 3
		vm.frameSize = MODE3_FRAME_SIZE
		vm.width = memory.SCREEN_WIDTH
		vm.height = memory.SCREEN_HEIGHT
		vm.bpp = 16
	}
}

// GetMode returns the current video mode
func (vm *VRAMManager) GetMode() int {
	return vm.mode
}

// GetCurrentPage returns the current buffer page (0 or 1)
func (vm *VRAMManager) GetCurrentPage() int {
	return vm.currentPage
}

// GetFrameSize returns the size in bytes of one frame buffer
func (vm *VRAMManager) GetFrameSize() int {
	return vm.frameSize
}

// GetDimensions returns the width and height of the current mode
func (vm *VRAMManager) GetDimensions() (int, int) {
	return vm.width, vm.height
}

// GetBPP returns the bits per pixel for the current mode
func (vm *VRAMManager) GetBPP() int {
	return vm.bpp
}

// SupportsDoubleBuffering returns true if the current mode supports double buffering
func (vm *VRAMManager) SupportsDoubleBuffering() bool {
	return vm.mode == memory.MODE_4 || vm.mode == memory.MODE_5
}

// GetCurrentBuffer returns the current active framebuffer
func (vm *VRAMManager) GetCurrentBuffer() *BitmapBuffer {
	if !vm.SupportsDoubleBuffering() {
		return &BitmapBuffer{
			base:   FRAME0_ADDR,
			width:  vm.width,
			height: vm.height,
			bpp:    vm.bpp,
		}
	}
	
	var base uintptr
	if vm.currentPage == 0 {
		base = FRAME0_ADDR
	} else {
		base = FRAME1_ADDR
	}
	
	return &BitmapBuffer{
		base:   base,
		width:  vm.width,
		height: vm.height,
		bpp:    vm.bpp,
	}
}

// GetBackBuffer returns the back buffer for double buffering
func (vm *VRAMManager) GetBackBuffer() *BitmapBuffer {
	if !vm.SupportsDoubleBuffering() {
		// No double buffering, return same buffer
		return vm.GetCurrentBuffer()
	}
	
	var base uintptr
	if vm.currentPage == 0 {
		base = FRAME1_ADDR
	} else {
		base = FRAME0_ADDR
	}
	
	return &BitmapBuffer{
		base:   base,
		width:  vm.width,
		height: vm.height,
		bpp:    vm.bpp,
	}
}

// SwapBuffers swaps the front and back buffers for double buffering
func (vm *VRAMManager) SwapBuffers() {
	if vm.SupportsDoubleBuffering() {
		vm.currentPage = 1 - vm.currentPage
	}
}

// VRAMAddr returns a VRAM address with offset validation
func VRAMAddr(offset uintptr) uintptr {
	addr := VRAM_BASE + offset
	if InVRAMBounds(addr) {
		return addr
	}
	return VRAM_BASE // Return base if out of bounds
}

// InVRAMBounds checks if an address is within VRAM bounds
func InVRAMBounds(addr uintptr) bool {
	return addr >= VRAM_BASE && addr < VRAM_BASE+VRAM_SIZE
}

// AlignAddress aligns an address to the specified boundary
func AlignAddress(addr uintptr, alignment int) uintptr {
	mask := uintptr(alignment - 1)
	return (addr + mask) &^ mask
}