package memory

// VRAMLayout provides access to VRAM regions for different video modes
type VRAMLayout struct {
	Mode0 *MemoryRegion // Character data
	Mode3 *MemoryRegion // 16-bit bitmap
	Mode4 *MemoryRegion // 8-bit bitmap  
	Mode5 *MemoryRegion // 16-bit small bitmap
}

// Global VRAM layout instance
var VRAMRegions = &VRAMLayout{
	Mode0: NewMemoryRegion(VRAM_BASE, VRAM_SIZE),           // Full VRAM for character data
	Mode3: NewMemoryRegion(VRAM_BASE, SCREEN_PIXELS*2),     // 240x160x16bit = 76,800 bytes
	Mode4: NewMemoryRegion(VRAM_BASE, SCREEN_PIXELS),       // 240x160x8bit = 38,400 bytes
	Mode5: NewMemoryRegion(VRAM_BASE, 160*128*2),           // 160x128x16bit = 40,960 bytes
}

// Mode3 utilities - 16-bit bitmap mode (240x160)
func SetPixelMode3(x, y int, color Color) {
	if x >= 0 && x < SCREEN_WIDTH && y >= 0 && y < SCREEN_HEIGHT {
		offset := uintptr((y*SCREEN_WIDTH + x) * 2)
		VRAMRegions.Mode3.WriteColor(offset, color)
	}
}

func GetPixelMode3(x, y int) Color {
	if x >= 0 && x < SCREEN_WIDTH && y >= 0 && y < SCREEN_HEIGHT {
		offset := uintptr((y*SCREEN_WIDTH + x) * 2)
		return VRAMRegions.Mode3.ReadColor(offset)
	}
	return BLACK
}

// Mode4 utilities - 8-bit bitmap mode (240x160)  
func SetPixelMode4(x, y int, paletteIndex uint8) {
	if x >= 0 && x < SCREEN_WIDTH && y >= 0 && y < SCREEN_HEIGHT {
		offset := uintptr(y*SCREEN_WIDTH + x)
		// Mode 4 uses 8-bit palette indices, but we need to handle alignment
		if offset%2 == 0 {
			// Even offset - lower byte
			current := VRAMRegions.Mode4.Read16(offset & ^uintptr(1))
			VRAMRegions.Mode4.Write16(offset & ^uintptr(1), (current&0xFF00)|uint16(paletteIndex))
		} else {
			// Odd offset - upper byte  
			current := VRAMRegions.Mode4.Read16(offset & ^uintptr(1))
			VRAMRegions.Mode4.Write16(offset & ^uintptr(1), (current&0x00FF)|(uint16(paletteIndex)<<8))
		}
	}
}

func GetPixelMode4(x, y int) uint8 {
	if x >= 0 && x < SCREEN_WIDTH && y >= 0 && y < SCREEN_HEIGHT {
		offset := uintptr(y*SCREEN_WIDTH + x)
		value := VRAMRegions.Mode4.Read16(offset & ^uintptr(1))
		if offset%2 == 0 {
			return uint8(value & 0xFF)
		} else {
			return uint8(value >> 8)
		}
	}
	return 0
}

// Mode5 utilities - 16-bit small bitmap mode (160x128)
func SetPixelMode5(x, y int, color Color) {
	if x >= 0 && x < 160 && y >= 0 && y < 128 {
		offset := uintptr((y*160 + x) * 2)
		VRAMRegions.Mode5.WriteColor(offset, color)
	}
}

func GetPixelMode5(x, y int) Color {
	if x >= 0 && x < 160 && y >= 0 && y < 128 {
		offset := uintptr((y*160 + x) * 2)
		return VRAMRegions.Mode5.ReadColor(offset)
	}
	return BLACK
}

// Drawing utilities
func DrawRectMode3(x, y, width, height int, color Color) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			SetPixelMode3(x+dx, y+dy, color)
		}
	}
}

func FillScreenMode3(color Color) {
	VRAMRegions.Mode3.FillColor(color)
}

func DrawRectMode5(x, y, width, height int, color Color) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			SetPixelMode5(x+dx, y+dy, color)
		}
	}
}

func FillScreenMode5(color Color) {
	VRAMRegions.Mode5.FillColor(color)
}

// Palette management for Mode 4
func SetBackgroundPalette(index uint8, color Color) {
	offset := uintptr(index) * 2
	PaletteRAM.WriteColor(offset, color)
}

func GetBackgroundPalette(index uint8) Color {
	offset := uintptr(index) * 2
	return PaletteRAM.ReadColor(offset)
}

func SetSpritePalette(index uint8, color Color) {
	// Sprite palettes start at offset 512 (256 colors * 2 bytes)
	offset := uintptr(512) + uintptr(index)*2
	PaletteRAM.WriteColor(offset, color)
}

func GetSpritePalette(index uint8) Color {
	offset := uintptr(512) + uintptr(index)*2
	return PaletteRAM.ReadColor(offset)
}

// Clear all VRAM
func ClearVRAM() {
	VRAM.Clear()
}

// Clear all palettes
func ClearPalettes() {
	PaletteRAM.Clear()
}