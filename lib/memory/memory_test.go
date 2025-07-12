package memory

import (
	"testing"
)

func TestMemoryConstants(t *testing.T) {
	// Verify memory layout constants
	if VRAM_BASE != 0x06000000 {
		t.Errorf("Expected VRAM_BASE to be 0x06000000, got 0x%08X", VRAM_BASE)
	}
	
	if VRAM_SIZE != 0x18000 {
		t.Errorf("Expected VRAM_SIZE to be 0x18000, got 0x%08X", VRAM_SIZE)
	}
	
	if SCREEN_WIDTH != 240 {
		t.Errorf("Expected SCREEN_WIDTH to be 240, got %d", SCREEN_WIDTH)
	}
	
	if SCREEN_HEIGHT != 160 {
		t.Errorf("Expected SCREEN_HEIGHT to be 160, got %d", SCREEN_HEIGHT)
	}
	
	if SCREEN_PIXELS != 38400 {
		t.Errorf("Expected SCREEN_PIXELS to be 38400, got %d", SCREEN_PIXELS)
	}
}

func TestColorOperations(t *testing.T) {
	// Test RGB15 color construction and extraction
	red := RGB15(31, 0, 0)
	if red.R() != 31 {
		t.Errorf("Expected red component to be 31, got %d", red.R())
	}
	if red.G() != 0 {
		t.Errorf("Expected green component to be 0, got %d", red.G())
	}
	if red.B() != 0 {
		t.Errorf("Expected blue component to be 0, got %d", red.B())
	}
	
	// Test RGB color construction (8-bit to 5-bit conversion)
	green := RGB(0, 255, 0)
	if green.R() != 0 {
		t.Errorf("Expected red component to be 0, got %d", green.R())
	}
	if green.G() != 31 {
		t.Errorf("Expected green component to be 31, got %d", green.G())
	}
	if green.B() != 0 {
		t.Errorf("Expected blue component to be 0, got %d", green.B())
	}
	
	// Test color constants
	if BLACK != RGB15(0, 0, 0) {
		t.Errorf("BLACK constant incorrect")
	}
	if WHITE != RGB15(31, 31, 31) {
		t.Errorf("WHITE constant incorrect")
	}
	if RED != RGB15(31, 0, 0) {
		t.Errorf("RED constant incorrect")
	}
}

func TestColorConversion(t *testing.T) {
	// Test ToRGB conversion
	color := RGB15(15, 16, 8)
	r, g, b := color.ToRGB()
	
	// 5-bit to 8-bit conversion: (value << 3) | (value >> 2)
	expectedR := uint8((15 << 3) | (15 >> 2)) // 123
	expectedG := uint8((16 << 3) | (16 >> 2)) // 132
	expectedB := uint8((8 << 3) | (8 >> 2))   // 66
	
	if r != expectedR {
		t.Errorf("Expected red to be %d, got %d", expectedR, r)
	}
	if g != expectedG {
		t.Errorf("Expected green to be %d, got %d", expectedG, g)
	}
	if b != expectedB {
		t.Errorf("Expected blue to be %d, got %d", expectedB, b)
	}
}

func TestPalette(t *testing.T) {
	var palette Palette
	
	// Test setting and getting colors
	palette.SetColor(0, RED)
	palette.SetColor(1, GREEN)
	palette.SetColor(255, BLUE)
	
	if palette.GetColor(0) != RED {
		t.Errorf("Expected color at index 0 to be RED")
	}
	if palette.GetColor(1) != GREEN {
		t.Errorf("Expected color at index 1 to be GREEN")
	}
	if palette.GetColor(255) != BLUE {
		t.Errorf("Expected color at index 255 to be BLUE")
	}
	
	// Test bounds checking - should return BLACK for out of bounds
	if palette.GetColor(100) != BLACK {
		t.Errorf("Expected unset color to be BLACK, got %v", palette.GetColor(100))
	}
}

func TestMemoryRegion(t *testing.T) {
	// Create a test memory region
	region := NewMemoryRegion(0x06000000, 1024)
	
	// Test bounds checking
	if !region.InBounds(0) {
		t.Errorf("Expected offset 0 to be in bounds")
	}
	if !region.InBounds(1023) {
		t.Errorf("Expected offset 1023 to be in bounds")
	}
	if region.InBounds(1024) {
		t.Errorf("Expected offset 1024 to be out of bounds")
	}
	
	// Test size and base
	if region.Size() != 1024 {
		t.Errorf("Expected size to be 1024, got %d", region.Size())
	}
	if region.Base() != 0x06000000 {
		t.Errorf("Expected base to be 0x06000000, got 0x%08X", region.Base())
	}
}

func TestGlobalMemoryRegions(t *testing.T) {
	// Test global VRAM instance
	vram := GetVRAM()
	if vram.Base() != VRAM_BASE {
		t.Errorf("Expected VRAM base to be 0x%08X, got 0x%08X", VRAM_BASE, vram.Base())
	}
	if vram.Size() != VRAM_SIZE {
		t.Errorf("Expected VRAM size to be 0x%08X, got 0x%08X", VRAM_SIZE, vram.Size())
	}
	
	// Test global OAM instance
	oam := GetOAM()
	if oam.Base() != OAM_BASE {
		t.Errorf("Expected OAM base to be 0x%08X, got 0x%08X", OAM_BASE, oam.Base())
	}
	if oam.Size() != OAM_SIZE {
		t.Errorf("Expected OAM size to be 0x%08X, got 0x%08X", OAM_SIZE, oam.Size())
	}
	
	// Test global Palette RAM instance
	palRAM := GetPaletteRAM()
	if palRAM.Base() != PALETTE_BASE {
		t.Errorf("Expected Palette RAM base to be 0x%08X, got 0x%08X", PALETTE_BASE, palRAM.Base())
	}
	if palRAM.Size() != PALETTE_SIZE {
		t.Errorf("Expected Palette RAM size to be 0x%08X, got 0x%08X", PALETTE_SIZE, palRAM.Size())
	}
}

func TestVRAMLayout(t *testing.T) {
	// Test Mode 3 region size (240x160x2 bytes = 76,800 bytes)
	expectedMode3Size := uintptr(SCREEN_PIXELS * 2)
	if VRAMRegions.Mode3.Size() != expectedMode3Size {
		t.Errorf("Expected Mode3 size to be %d, got %d", expectedMode3Size, VRAMRegions.Mode3.Size())
	}
	
	// Test Mode 4 region size (240x160x1 byte = 38,400 bytes)
	expectedMode4Size := uintptr(SCREEN_PIXELS)
	if VRAMRegions.Mode4.Size() != expectedMode4Size {
		t.Errorf("Expected Mode4 size to be %d, got %d", expectedMode4Size, VRAMRegions.Mode4.Size())
	}
	
	// Test Mode 5 region size (160x128x2 bytes = 40,960 bytes)
	expectedMode5Size := uintptr(160 * 128 * 2)
	if VRAMRegions.Mode5.Size() != expectedMode5Size {
		t.Errorf("Expected Mode5 size to be %d, got %d", expectedMode5Size, VRAMRegions.Mode5.Size())
	}
}