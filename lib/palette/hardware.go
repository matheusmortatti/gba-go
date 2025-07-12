package palette

import (
	"runtime/volatile"
	"unsafe"
)

const (
	// Palette RAM layout
	BG_PALETTE_BASE  = 0x05000000
	OBJ_PALETTE_BASE = 0x05000200
	PALETTE_SIZE     = 0x400 // 1KB total
	BG_PALETTE_SIZE  = 0x200 // 512 bytes
	OBJ_PALETTE_SIZE = 0x200 // 512 bytes
)

// PaletteBank represents the hardware palette memory
type PaletteBank struct {
	bgBank  *[BG_PALETTE_SIZE / 2]volatile.Register16
	objBank *[OBJ_PALETTE_SIZE / 2]volatile.Register16
}

var paletteBank *PaletteBank

// GetPaletteBank returns the global palette bank instance
func GetPaletteBank() *PaletteBank {
	if paletteBank == nil {
		paletteBank = &PaletteBank{
			bgBank:  (*[BG_PALETTE_SIZE / 2]volatile.Register16)(unsafe.Pointer(uintptr(BG_PALETTE_BASE))),
			objBank: (*[OBJ_PALETTE_SIZE / 2]volatile.Register16)(unsafe.Pointer(uintptr(OBJ_PALETTE_BASE))),
		}
	}
	return paletteBank
}

// SetBGColor sets a background palette color at the specified palette and color index
func (pb *PaletteBank) SetBGColor(paletteIndex, colorIndex int, color Color) {
	if paletteIndex < 0 || paletteIndex >= MAX_BG_PALETTES ||
		colorIndex < 0 || colorIndex >= COLORS_PER_PALETTE_16 {
		return // Silently ignore invalid indices
	}

	offset := paletteIndex*COLORS_PER_PALETTE_16 + colorIndex
	pb.bgBank[offset].Set(uint16(color))
}

// GetBGColor gets a background palette color at the specified palette and color index
func (pb *PaletteBank) GetBGColor(paletteIndex, colorIndex int) Color {
	if paletteIndex < 0 || paletteIndex >= MAX_BG_PALETTES ||
		colorIndex < 0 || colorIndex >= COLORS_PER_PALETTE_16 {
		return BLACK // Return black for invalid indices
	}

	offset := paletteIndex*COLORS_PER_PALETTE_16 + colorIndex
	return Color(pb.bgBank[offset].Get())
}

// SetOBJColor sets a sprite palette color at the specified palette and color index
func (pb *PaletteBank) SetOBJColor(paletteIndex, colorIndex int, color Color) {
	if paletteIndex < 0 || paletteIndex >= MAX_OBJ_PALETTES ||
		colorIndex < 0 || colorIndex >= COLORS_PER_PALETTE_16 {
		return // Silently ignore invalid indices
	}

	offset := paletteIndex*COLORS_PER_PALETTE_16 + colorIndex
	pb.objBank[offset].Set(uint16(color))
}

// GetOBJColor gets a sprite palette color at the specified palette and color index
func (pb *PaletteBank) GetOBJColor(paletteIndex, colorIndex int) Color {
	if paletteIndex < 0 || paletteIndex >= MAX_OBJ_PALETTES ||
		colorIndex < 0 || colorIndex >= COLORS_PER_PALETTE_16 {
		return BLACK // Return black for invalid indices
	}

	offset := paletteIndex*COLORS_PER_PALETTE_16 + colorIndex
	return Color(pb.objBank[offset].Get())
}

// SetBG256Color sets a color in 256-color background palette mode
func (pb *PaletteBank) SetBG256Color(colorIndex int, color Color) {
	if colorIndex < 0 || colorIndex >= COLORS_PER_PALETTE_256 {
		return // Silently ignore invalid indices
	}

	pb.bgBank[colorIndex].Set(uint16(color))
}

// GetBG256Color gets a color from 256-color background palette mode
func (pb *PaletteBank) GetBG256Color(colorIndex int) Color {
	if colorIndex < 0 || colorIndex >= COLORS_PER_PALETTE_256 {
		return BLACK // Return black for invalid indices
	}

	return Color(pb.bgBank[colorIndex].Get())
}

// SetOBJ256Color sets a color in 256-color sprite palette mode
func (pb *PaletteBank) SetOBJ256Color(colorIndex int, color Color) {
	if colorIndex < 0 || colorIndex >= COLORS_PER_PALETTE_256 {
		return // Silently ignore invalid indices
	}

	pb.objBank[colorIndex].Set(uint16(color))
}

// GetOBJ256Color gets a color from 256-color sprite palette mode
func (pb *PaletteBank) GetOBJ256Color(colorIndex int) Color {
	if colorIndex < 0 || colorIndex >= COLORS_PER_PALETTE_256 {
		return BLACK // Return black for invalid indices
	}

	return Color(pb.objBank[colorIndex].Get())
}

// LoadBGPalette16 loads a 16-color palette to the specified background palette slot
func (pb *PaletteBank) LoadBGPalette16(paletteIndex int, palette *Palette16) {
	if paletteIndex < 0 || paletteIndex >= MAX_BG_PALETTES {
		return
	}

	for i := 0; i < COLORS_PER_PALETTE_16; i++ {
		pb.SetBGColor(paletteIndex, i, palette[i])
	}
}

// LoadBGPalette256 loads a 256-color palette to background palette memory
func (pb *PaletteBank) LoadBGPalette256(palette *Palette256) {
	for i := 0; i < COLORS_PER_PALETTE_256; i++ {
		pb.SetBG256Color(i, palette[i])
	}
}

// LoadOBJPalette16 loads a 16-color palette to the specified sprite palette slot
func (pb *PaletteBank) LoadOBJPalette16(paletteIndex int, palette *Palette16) {
	if paletteIndex < 0 || paletteIndex >= MAX_OBJ_PALETTES {
		return
	}

	for i := 0; i < COLORS_PER_PALETTE_16; i++ {
		pb.SetOBJColor(paletteIndex, i, palette[i])
	}
}

// LoadOBJPalette256 loads a 256-color palette to sprite palette memory
func (pb *PaletteBank) LoadOBJPalette256(palette *Palette256) {
	for i := 0; i < COLORS_PER_PALETTE_256; i++ {
		pb.SetOBJ256Color(i, palette[i])
	}
}
