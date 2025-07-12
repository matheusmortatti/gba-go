package vram

import (
	"runtime/volatile"
	"unsafe"

	"github.com/matheusmortatti/gba-go/lib/memory"
)

// VRAMDebugInfo provides debugging information about VRAM usage
type VRAMDebugInfo struct {
	TotalSize      int
	UsedSize       int
	FreeSize       int
	Mode           int
	FrameSize      int
	SupportsDouble bool
	CurrentPage    int
}

// GetVRAMDebugInfo returns debugging information about VRAM
func GetVRAMDebugInfo(vm *VRAMManager) *VRAMDebugInfo {
	frameSize := vm.GetFrameSize()
	var usedSize int
	
	if vm.SupportsDoubleBuffering() {
		usedSize = frameSize * 2 // Two frame buffers
	} else {
		usedSize = frameSize // Single frame buffer
	}
	
	return &VRAMDebugInfo{
		TotalSize:      VRAM_SIZE,
		UsedSize:       usedSize,
		FreeSize:       VRAM_SIZE - usedSize,
		Mode:           vm.GetMode(),
		FrameSize:      frameSize,
		SupportsDouble: vm.SupportsDoubleBuffering(),
		CurrentPage:    vm.GetCurrentPage(),
	}
}

// VRAMPattern represents a test pattern for VRAM testing
type VRAMPattern int

const (
	PATTERN_SOLID VRAMPattern = iota
	PATTERN_CHECKERBOARD
	PATTERN_GRADIENT_H
	PATTERN_GRADIENT_V
	PATTERN_STRIPES_H
	PATTERN_STRIPES_V
	PATTERN_RAINBOW
	PATTERN_NOISE
)

// FillPattern fills a buffer with a test pattern
func FillPattern(buffer *BitmapBuffer, pattern VRAMPattern, color1, color2 uint16) {
	width := buffer.GetWidth()
	height := buffer.GetHeight()
	
	switch pattern {
	case PATTERN_SOLID:
		buffer.Clear(color1)
		
	case PATTERN_CHECKERBOARD:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if (x/8+y/8)%2 == 0 {
					buffer.PlotPixelFast(x, y, color1)
				} else {
					buffer.PlotPixelFast(x, y, color2)
				}
			}
		}
		
	case PATTERN_GRADIENT_H:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				intensity := (x * 31) / width
				if buffer.GetBPP() == 16 {
					color := uint16(intensity | (intensity << 5) | (intensity << 10))
					buffer.PlotPixelFast(x, y, color)
				} else {
					color := uint16((intensity * 255) / 31)
					buffer.PlotPixelFast(x, y, color)
				}
			}
		}
		
	case PATTERN_GRADIENT_V:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				intensity := (y * 31) / height
				if buffer.GetBPP() == 16 {
					color := uint16(intensity | (intensity << 5) | (intensity << 10))
					buffer.PlotPixelFast(x, y, color)
				} else {
					color := uint16((intensity * 255) / 31)
					buffer.PlotPixelFast(x, y, color)
				}
			}
		}
		
	case PATTERN_STRIPES_H:
		for y := 0; y < height; y++ {
			color := color1
			if (y/4)%2 == 1 {
				color = color2
			}
			for x := 0; x < width; x++ {
				buffer.PlotPixelFast(x, y, color)
			}
		}
		
	case PATTERN_STRIPES_V:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				color := color1
				if (x/4)%2 == 1 {
					color = color2
				}
				buffer.PlotPixelFast(x, y, color)
			}
		}
		
	case PATTERN_RAINBOW:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if buffer.GetBPP() == 16 {
					r := (x * 31) / width
					g := (y * 31) / height
					b := ((x + y) * 31) / (width + height)
					color := uint16(r | (g << 5) | (b << 10))
					buffer.PlotPixelFast(x, y, color)
				} else {
					color := uint16(((x + y) * 255) / (width + height))
					buffer.PlotPixelFast(x, y, color)
				}
			}
		}
		
	case PATTERN_NOISE:
		// Simple pseudo-random pattern
		seed := uint32(0x12345678)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				seed = seed*1103515245 + 12345
				if buffer.GetBPP() == 16 {
					color := uint16(seed & 0x7FFF)
					buffer.PlotPixelFast(x, y, color)
				} else {
					color := uint16(seed & 0xFF)
					buffer.PlotPixelFast(x, y, color)
				}
			}
		}
	}
}

// MemoryUsageInfo provides information about memory usage
type MemoryUsageInfo struct {
	VRAMUsed      int
	VRAMFree      int
	CharDataUsed  int
	ScreenDataUsed int
	BitmapDataUsed int
}

// GetMemoryUsage calculates memory usage for different VRAM regions
func GetMemoryUsage(mode int) *MemoryUsageInfo {
	info := &MemoryUsageInfo{}
	
	switch mode {
	case memory.MODE_0, memory.MODE_1, memory.MODE_2:
		// Tile modes - character and screen data
		info.CharDataUsed = MAX_CHAR_BLOCKS * CHAR_BASE_SIZE
		info.ScreenDataUsed = MAX_SCREEN_BLOCKS * SCREEN_BASE_SIZE
		info.VRAMUsed = info.CharDataUsed + info.ScreenDataUsed
		
	case memory.MODE_3:
		info.BitmapDataUsed = MODE3_FRAME_SIZE
		info.VRAMUsed = info.BitmapDataUsed
		
	case memory.MODE_4:
		info.BitmapDataUsed = MODE4_FRAME_SIZE * 2 // Two frames
		info.VRAMUsed = info.BitmapDataUsed
		
	case memory.MODE_5:
		info.BitmapDataUsed = MODE5_FRAME_SIZE * 2 // Two frames
		info.VRAMUsed = info.BitmapDataUsed
	}
	
	info.VRAMFree = VRAM_SIZE - info.VRAMUsed
	return info
}

// ValidateVRAMAccess checks if a memory access is valid
func ValidateVRAMAccess(addr uintptr, size int) bool {
	return addr >= VRAM_BASE && addr+uintptr(size) <= VRAM_BASE+VRAM_SIZE
}

// DumpVRAMRegion dumps a region of VRAM to a byte slice for debugging
func DumpVRAMRegion(offset, size uintptr) []uint8 {
	if !ValidateVRAMAccess(VRAM_BASE+offset, int(size)) {
		return nil
	}
	
	data := make([]uint8, size)
	for i := uintptr(0); i < size; i += 2 {
		addr := VRAM_BASE + offset + i
		value := (*volatile.Register16)(unsafe.Pointer(addr)).Get()
		data[i] = uint8(value & 0xFF)
		if i+1 < size {
			data[i+1] = uint8(value >> 8)
		}
	}
	
	return data
}

// LoadVRAMRegion loads data into a VRAM region
func LoadVRAMRegion(offset uintptr, data []uint8) bool {
	if !ValidateVRAMAccess(VRAM_BASE+offset, len(data)) {
		return false
	}
	
	for i := 0; i < len(data); i += 2 {
		addr := VRAM_BASE + offset + uintptr(i)
		var value uint16
		if i+1 < len(data) {
			value = uint16(data[i]) | (uint16(data[i+1]) << 8)
		} else {
			value = uint16(data[i])
		}
		(*volatile.Register16)(unsafe.Pointer(addr)).Set(value)
	}
	
	return true
}

// CalculateCharBlock calculates which character block an address belongs to
func CalculateCharBlock(addr uintptr) int {
	if addr < VRAM_BASE || addr >= VRAM_BASE+VRAM_SIZE {
		return -1
	}
	
	offset := addr - VRAM_BASE
	return int(offset / CHAR_BASE_SIZE)
}

// CalculateScreenBlock calculates which screen block an address belongs to
func CalculateScreenBlock(addr uintptr) int {
	if addr < VRAM_BASE || addr >= VRAM_BASE+VRAM_SIZE {
		return -1
	}
	
	offset := addr - VRAM_BASE
	
	// Screen blocks start after character data in some layouts
	// This is a simplified calculation
	if offset >= 0x8000 { // Screen data typically starts at 0x06008000
		screenOffset := offset - 0x8000
		return int(screenOffset / SCREEN_BASE_SIZE)
	}
	
	return -1
}

// Performance measurement utilities
type PerformanceCounter struct {
	operations int
	startTime  int // This would need a proper timer implementation
}

// NewPerformanceCounter creates a new performance counter
func NewPerformanceCounter() *PerformanceCounter {
	return &PerformanceCounter{
		operations: 0,
		startTime:  0, // Would need actual timer
	}
}

// Start begins performance measurement
func (pc *PerformanceCounter) Start() {
	pc.operations = 0
	// pc.startTime = getCurrentTime() // Would need timer implementation
}

// AddOperation increments the operation counter
func (pc *PerformanceCounter) AddOperation() {
	pc.operations++
}

// GetOperations returns the number of operations performed
func (pc *PerformanceCounter) GetOperations() int {
	return pc.operations
}

// VRAMAddressInfo provides information about a VRAM address
type VRAMAddressInfo struct {
	Valid       bool
	Offset      uintptr
	CharBlock   int
	ScreenBlock int
	TileIndex   int
	PixelX      int
	PixelY      int
}

// AnalyzeVRAMAddress analyzes a VRAM address and provides information about it
func AnalyzeVRAMAddress(addr uintptr, mode int) *VRAMAddressInfo {
	info := &VRAMAddressInfo{
		Valid:       ValidateVRAMAccess(addr, 1),
		Offset:      addr - VRAM_BASE,
		CharBlock:   -1,
		ScreenBlock: -1,
		TileIndex:   -1,
		PixelX:      -1,
		PixelY:      -1,
	}
	
	if !info.Valid {
		return info
	}
	
	switch mode {
	case memory.MODE_3:
		// Calculate pixel coordinates for Mode 3
		pixelOffset := info.Offset / 2 // 16-bit pixels
		info.PixelY = int(pixelOffset) / memory.SCREEN_WIDTH
		info.PixelX = int(pixelOffset) % memory.SCREEN_WIDTH
		
	case memory.MODE_4:
		// Calculate pixel coordinates for Mode 4
		if info.Offset < MODE4_FRAME_SIZE {
			// Frame 0
			info.PixelY = int(info.Offset) / memory.SCREEN_WIDTH
			info.PixelX = int(info.Offset) % memory.SCREEN_WIDTH
		} else if info.Offset < MODE4_FRAME_SIZE*2 {
			// Frame 1
			frameOffset := info.Offset - MODE4_FRAME_SIZE
			info.PixelY = int(frameOffset) / memory.SCREEN_WIDTH
			info.PixelX = int(frameOffset) % memory.SCREEN_WIDTH
		}
		
	case memory.MODE_5:
		// Calculate pixel coordinates for Mode 5
		pixelOffset := info.Offset / 2 // 16-bit pixels
		if pixelOffset < 160*128 {
			// Frame 0
			info.PixelY = int(pixelOffset) / 160
			info.PixelX = int(pixelOffset) % 160
		} else if pixelOffset < 160*128*2 {
			// Frame 1
			frameOffset := pixelOffset - 160*128
			info.PixelY = int(frameOffset) / 160
			info.PixelX = int(frameOffset) % 160
		}
		
	default:
		// Tile modes
		info.CharBlock = CalculateCharBlock(addr)
		info.ScreenBlock = CalculateScreenBlock(addr)
		
		if info.CharBlock >= 0 {
			charOffset := info.Offset - uintptr(info.CharBlock*CHAR_BASE_SIZE)
			info.TileIndex = int(charOffset / TILE_4BPP_SIZE)
		}
	}
	
	return info
}