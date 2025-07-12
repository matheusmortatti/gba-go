package vram

import (
	"testing"
	
	"github.com/matheusmortatti/gba-go/lib/memory"
)

func TestVRAMManager(t *testing.T) {
	// Test Mode 3
	vm := NewVRAMManager(memory.MODE_3)
	if vm.GetMode() != memory.MODE_3 {
		t.Errorf("Expected mode %d, got %d", memory.MODE_3, vm.GetMode())
	}
	
	width, height := vm.GetDimensions()
	if width != memory.SCREEN_WIDTH || height != memory.SCREEN_HEIGHT {
		t.Errorf("Expected dimensions %dx%d, got %dx%d", 
			memory.SCREEN_WIDTH, memory.SCREEN_HEIGHT, width, height)
	}
	
	if vm.GetBPP() != 16 {
		t.Errorf("Expected 16 BPP for Mode 3, got %d", vm.GetBPP())
	}
	
	if vm.SupportsDoubleBuffering() {
		t.Error("Mode 3 should not support double buffering")
	}
	
	// Test Mode 4
	vm.SetMode(memory.MODE_4)
	if vm.GetMode() != memory.MODE_4 {
		t.Errorf("Expected mode %d, got %d", memory.MODE_4, vm.GetMode())
	}
	
	if vm.GetBPP() != 8 {
		t.Errorf("Expected 8 BPP for Mode 4, got %d", vm.GetBPP())
	}
	
	if !vm.SupportsDoubleBuffering() {
		t.Error("Mode 4 should support double buffering")
	}
	
	// Test Mode 5
	vm.SetMode(memory.MODE_5)
	width, height = vm.GetDimensions()
	if width != 160 || height != 128 {
		t.Errorf("Expected dimensions 160x128 for Mode 5, got %dx%d", width, height)
	}
}

func TestBitmapBuffer(t *testing.T) {
	vm := NewVRAMManager(memory.MODE_3)
	buffer := vm.GetCurrentBuffer()
	
	if buffer == nil {
		t.Fatal("Buffer should not be nil")
	}
	
	// Test bounds checking
	if !buffer.InBounds(0, 0) {
		t.Error("Origin should be in bounds")
	}
	
	if !buffer.InBounds(memory.SCREEN_WIDTH-1, memory.SCREEN_HEIGHT-1) {
		t.Error("Max coordinates should be in bounds")
	}
	
	if buffer.InBounds(-1, 0) {
		t.Error("Negative coordinates should be out of bounds")
	}
	
	if buffer.InBounds(memory.SCREEN_WIDTH, memory.SCREEN_HEIGHT) {
		t.Error("Max+1 coordinates should be out of bounds")
	}
	
	// Test pixel operations
	color := uint16(0x7FFF) // White
	err := buffer.PlotPixel(100, 80, color)
	if err != nil {
		t.Errorf("PlotPixel failed: %v", err)
	}
	
	// Note: We can't easily test GetPixel in this environment without actual hardware
	// but we can test that it doesn't crash
	_, err = buffer.GetPixel(100, 80)
	if err != nil {
		t.Errorf("GetPixel failed: %v", err)
	}
	
	// Test out of bounds error
	err = buffer.PlotPixel(300, 200, color)
	if err == nil {
		t.Error("Expected error for out of bounds pixel")
	}
}

func TestDoubleBuffering(t *testing.T) {
	vm := NewVRAMManager(memory.MODE_4)
	
	front := vm.GetCurrentBuffer()
	back := vm.GetBackBuffer()
	
	if front == nil || back == nil {
		t.Fatal("Buffers should not be nil")
	}
	
	// Buffers should be different
	if front.GetBase() == back.GetBase() {
		t.Error("Front and back buffers should have different base addresses")
	}
	
	// Test buffer swapping
	originalFrontBase := front.GetBase()
	originalBackBase := back.GetBase()
	
	vm.SwapBuffers()
	
	newFront := vm.GetCurrentBuffer()
	newBack := vm.GetBackBuffer()
	
	if newFront.GetBase() != originalBackBase {
		t.Error("After swap, front buffer should be original back buffer")
	}
	
	if newBack.GetBase() != originalFrontBase {
		t.Error("After swap, back buffer should be original front buffer")
	}
}

func TestTileData(t *testing.T) {
	td := GetCharacterData(0, 4) // Character block 0, 4 BPP
	
	if td == nil {
		t.Fatal("TileData should not be nil")
	}
	
	if td.GetCharBlock() != 0 {
		t.Errorf("Expected character block 0, got %d", td.GetCharBlock())
	}
	
	if td.GetBPP() != 4 {
		t.Errorf("Expected 4 BPP, got %d", td.GetBPP())
	}
	
	expectedMaxTiles := CHAR_BASE_SIZE / TILE_4BPP_SIZE
	if td.GetMaxTiles() != expectedMaxTiles {
		t.Errorf("Expected %d max tiles, got %d", expectedMaxTiles, td.GetMaxTiles())
	}
	
	// Test invalid character block
	invalidTd := GetCharacterData(-1, 4)
	if invalidTd != nil {
		t.Error("Invalid character block should return nil")
	}
	
	invalidTd = GetCharacterData(MAX_CHAR_BLOCKS, 4)
	if invalidTd != nil {
		t.Error("Out of range character block should return nil")
	}
}

func TestScreenData(t *testing.T) {
	sd := GetScreenData(0, 32, 32) // Screen block 0, 32x32 tiles
	
	if sd == nil {
		t.Fatal("ScreenData should not be nil")
	}
	
	if sd.GetScreenBlock() != 0 {
		t.Errorf("Expected screen block 0, got %d", sd.GetScreenBlock())
	}
	
	width, height := sd.GetDimensions()
	if width != 32 || height != 32 {
		t.Errorf("Expected dimensions 32x32, got %dx%d", width, height)
	}
	
	// Test bounds checking
	if !sd.InBounds(0, 0) {
		t.Error("Origin should be in bounds")
	}
	
	if !sd.InBounds(31, 31) {
		t.Error("Max coordinates should be in bounds")
	}
	
	if sd.InBounds(-1, 0) {
		t.Error("Negative coordinates should be out of bounds")
	}
	
	if sd.InBounds(32, 32) {
		t.Error("Max+1 coordinates should be out of bounds")
	}
	
	// Test invalid screen block
	invalidSd := GetScreenData(-1, 32, 32)
	if invalidSd != nil {
		t.Error("Invalid screen block should return nil")
	}
}

func TestVRAMUtilities(t *testing.T) {
	// Test address validation
	if !InVRAMBounds(VRAM_BASE) {
		t.Error("VRAM base should be in bounds")
	}
	
	if !InVRAMBounds(VRAM_BASE + VRAM_SIZE - 1) {
		t.Error("VRAM end should be in bounds")
	}
	
	if InVRAMBounds(VRAM_BASE - 1) {
		t.Error("Before VRAM should be out of bounds")
	}
	
	if InVRAMBounds(VRAM_BASE + VRAM_SIZE) {
		t.Error("After VRAM should be out of bounds")
	}
	
	// Test address alignment
	aligned := AlignAddress(0x06000001, 4)
	if aligned != 0x06000004 {
		t.Errorf("Expected aligned address 0x06000004, got 0x%08X", aligned)
	}
	
	aligned = AlignAddress(0x06000004, 4)
	if aligned != 0x06000004 {
		t.Errorf("Already aligned address should remain the same, got 0x%08X", aligned)
	}
}

func TestMemoryUsage(t *testing.T) {
	// Test Mode 3 memory usage
	usage := GetMemoryUsage(memory.MODE_3)
	if usage.BitmapDataUsed != MODE3_FRAME_SIZE {
		t.Errorf("Expected bitmap data used %d, got %d", 
			MODE3_FRAME_SIZE, usage.BitmapDataUsed)
	}
	
	// Test Mode 4 memory usage
	usage = GetMemoryUsage(memory.MODE_4)
	expected := MODE4_FRAME_SIZE * 2 // Two frames
	if usage.BitmapDataUsed != expected {
		t.Errorf("Expected bitmap data used %d, got %d", 
			expected, usage.BitmapDataUsed)
	}
	
	// Test that free size calculation is correct
	if usage.VRAMFree != VRAM_SIZE-usage.VRAMUsed {
		t.Error("Free size calculation is incorrect")
	}
}

func TestVRAMAddressAnalysis(t *testing.T) {
	// Test Mode 3 address analysis
	addr := uintptr(VRAM_BASE + 100*2) // Pixel at offset 100
	info := AnalyzeVRAMAddress(addr, memory.MODE_3)
	
	if !info.Valid {
		t.Error("Valid VRAM address should be marked as valid")
	}
	
	expectedY := 100 / memory.SCREEN_WIDTH
	expectedX := 100 % memory.SCREEN_WIDTH
	
	if info.PixelX != expectedX || info.PixelY != expectedY {
		t.Errorf("Expected pixel coordinates (%d, %d), got (%d, %d)",
			expectedX, expectedY, info.PixelX, info.PixelY)
	}
	
	// Test invalid address
	invalidAddr := uintptr(0x05000000) // Palette RAM, not VRAM
	info = AnalyzeVRAMAddress(invalidAddr, memory.MODE_3)
	
	if info.Valid {
		t.Error("Invalid address should be marked as invalid")
	}
}

func TestTileAttributes(t *testing.T) {
	// Test palette attribute
	palAttr := SetTilePalette(5)
	expected := uint16(5 << 12)
	if palAttr != expected {
		t.Errorf("Expected palette attribute 0x%04X, got 0x%04X", expected, palAttr)
	}
	
	// Test flip attributes
	flipAttr := SetTileFlip(true, false)
	if flipAttr != TILE_HFLIP {
		t.Errorf("Expected horizontal flip attribute 0x%04X, got 0x%04X", 
			TILE_HFLIP, flipAttr)
	}
	
	flipAttr = SetTileFlip(false, true)
	if flipAttr != TILE_VFLIP {
		t.Errorf("Expected vertical flip attribute 0x%04X, got 0x%04X", 
			TILE_VFLIP, flipAttr)
	}
	
	flipAttr = SetTileFlip(true, true)
	if flipAttr != (TILE_HFLIP|TILE_VFLIP) {
		t.Errorf("Expected both flip attributes 0x%04X, got 0x%04X", 
			TILE_HFLIP|TILE_VFLIP, flipAttr)
	}
	
	// Test combining attributes
	combined := CombineTileAttributes(SetTilePalette(3), SetTileFlip(true, false))
	expected = (3 << 12) | TILE_HFLIP
	if combined != expected {
		t.Errorf("Expected combined attributes 0x%04X, got 0x%04X", expected, combined)
	}
}

// Benchmark tests
func BenchmarkPixelPlot(b *testing.B) {
	vm := NewVRAMManager(memory.MODE_3)
	buffer := vm.GetCurrentBuffer()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := i % memory.SCREEN_WIDTH
		y := (i / memory.SCREEN_WIDTH) % memory.SCREEN_HEIGHT
		buffer.PlotPixelFast(x, y, uint16(i))
	}
}

func BenchmarkFastClear(b *testing.B) {
	vm := NewVRAMManager(memory.MODE_3)
	buffer := vm.GetCurrentBuffer()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.FastClear(0x0000)
	}
}

func BenchmarkRegularClear(b *testing.B) {
	vm := NewVRAMManager(memory.MODE_3)
	buffer := vm.GetCurrentBuffer()
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buffer.Clear(0x0000)
	}
}