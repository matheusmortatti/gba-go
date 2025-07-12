package vram

import (
	"errors"
	"runtime/volatile"
	"unsafe"
)

// TileData represents character data in VRAM
type TileData struct {
	base      uintptr
	charBlock int
	bpp       int
	maxTiles  int
}

// ScreenData represents tile map data
type ScreenData struct {
	base        uintptr
	screenBlock int
	width       int // in tiles
	height      int // in tiles
	maxEntries  int
}

// GetCharacterData returns a TileData structure for the specified character block
func GetCharacterData(charBlock int, bpp int) *TileData {
	if charBlock < 0 || charBlock >= MAX_CHAR_BLOCKS {
		return nil
	}

	base := VRAM_BASE + uintptr(charBlock*CHAR_BASE_SIZE)
	var maxTiles int
	
	if bpp == 4 {
		maxTiles = CHAR_BASE_SIZE / TILE_4BPP_SIZE
	} else {
		maxTiles = CHAR_BASE_SIZE / TILE_8BPP_SIZE
	}

	return &TileData{
		base:      base,
		charBlock: charBlock,
		bpp:       bpp,
		maxTiles:  maxTiles,
	}
}

// GetScreenData returns a ScreenData structure for the specified screen block
func GetScreenData(screenBlock int, width, height int) *ScreenData {
	if screenBlock < 0 || screenBlock >= MAX_SCREEN_BLOCKS {
		return nil
	}

	base := VRAM_BASE + uintptr(screenBlock*SCREEN_BASE_SIZE)
	maxEntries := SCREEN_BASE_SIZE / 2 // 2 bytes per screen entry

	return &ScreenData{
		base:        base,
		screenBlock: screenBlock,
		width:       width,
		height:      height,
		maxEntries:  maxEntries,
	}
}

// GetBase returns the base address of the tile data
func (td *TileData) GetBase() uintptr {
	return td.base
}

// GetCharBlock returns the character block index
func (td *TileData) GetCharBlock() int {
	return td.charBlock
}

// GetBPP returns the bits per pixel
func (td *TileData) GetBPP() int {
	return td.bpp
}

// GetMaxTiles returns the maximum number of tiles that can be stored
func (td *TileData) GetMaxTiles() int {
	return td.maxTiles
}

// LoadTile loads tile data into the specified tile index
func (td *TileData) LoadTile(tileIndex int, data []uint8) error {
	if tileIndex < 0 || tileIndex >= td.maxTiles {
		return errors.New("tile index out of range")
	}

	var tileSize int
	if td.bpp == 4 {
		tileSize = TILE_4BPP_SIZE
	} else {
		tileSize = TILE_8BPP_SIZE
	}

	if len(data) != tileSize {
		return errors.New("invalid tile data size")
	}

	tileAddr := td.base + uintptr(tileIndex*tileSize)
	
	// Copy data 16 bits at a time for efficiency
	for i := 0; i < tileSize; i += 2 {
		var value uint16
		if i+1 < len(data) {
			value = uint16(data[i]) | (uint16(data[i+1]) << 8)
		} else {
			value = uint16(data[i])
		}
		
		(*volatile.Register16)(unsafe.Pointer(tileAddr + uintptr(i))).Set(value)
	}

	return nil
}

// GetTile retrieves tile data from the specified tile index
func (td *TileData) GetTile(tileIndex int) ([]uint8, error) {
	if tileIndex < 0 || tileIndex >= td.maxTiles {
		return nil, errors.New("tile index out of range")
	}

	var tileSize int
	if td.bpp == 4 {
		tileSize = TILE_4BPP_SIZE
	} else {
		tileSize = TILE_8BPP_SIZE
	}

	data := make([]uint8, tileSize)
	tileAddr := td.base + uintptr(tileIndex*tileSize)
	
	// Read data 16 bits at a time
	for i := 0; i < tileSize; i += 2 {
		value := (*volatile.Register16)(unsafe.Pointer(tileAddr + uintptr(i))).Get()
		data[i] = uint8(value & 0xFF)
		if i+1 < tileSize {
			data[i+1] = uint8(value >> 8)
		}
	}

	return data, nil
}

// ClearTile clears the specified tile (fills with zeros)
func (td *TileData) ClearTile(tileIndex int) error {
	if tileIndex < 0 || tileIndex >= td.maxTiles {
		return errors.New("tile index out of range")
	}

	var tileSize int
	if td.bpp == 4 {
		tileSize = TILE_4BPP_SIZE
	} else {
		tileSize = TILE_8BPP_SIZE
	}

	tileAddr := td.base + uintptr(tileIndex*tileSize)
	
	// Clear 16 bits at a time
	for i := 0; i < tileSize; i += 2 {
		(*volatile.Register16)(unsafe.Pointer(tileAddr + uintptr(i))).Set(0)
	}

	return nil
}

// GetBase returns the base address of the screen data
func (sd *ScreenData) GetBase() uintptr {
	return sd.base
}

// GetScreenBlock returns the screen block index
func (sd *ScreenData) GetScreenBlock() int {
	return sd.screenBlock
}

// GetDimensions returns the width and height in tiles
func (sd *ScreenData) GetDimensions() (int, int) {
	return sd.width, sd.height
}

// GetMaxEntries returns the maximum number of screen entries
func (sd *ScreenData) GetMaxEntries() int {
	return sd.maxEntries
}

// InBounds checks if tile coordinates are within bounds
func (sd *ScreenData) InBounds(x, y int) bool {
	return x >= 0 && x < sd.width && y >= 0 && y < sd.height
}

// SetTile sets a tile at the specified coordinates with attributes
func (sd *ScreenData) SetTile(x, y, tileIndex int, attributes uint16) error {
	if !sd.InBounds(x, y) {
		return errors.New("tile coordinates out of bounds")
	}

	entryIndex := y*sd.width + x
	if entryIndex >= sd.maxEntries {
		return errors.New("screen entry index out of range")
	}

	// Screen entry format: 
	// Bits 0-9: Tile index (0-1023)
	// Bits 10-11: Horizontal/Vertical flip
	// Bits 12-15: Palette selection (4bpp mode)
	screenEntry := uint16(tileIndex&0x3FF) | (attributes & 0xFC00)
	
	entryAddr := sd.base + uintptr(entryIndex*2)
	(*volatile.Register16)(unsafe.Pointer(entryAddr)).Set(screenEntry)

	return nil
}

// GetTile gets the tile index and attributes at the specified coordinates
func (sd *ScreenData) GetTile(x, y int) (int, uint16, error) {
	if !sd.InBounds(x, y) {
		return 0, 0, errors.New("tile coordinates out of bounds")
	}

	entryIndex := y*sd.width + x
	if entryIndex >= sd.maxEntries {
		return 0, 0, errors.New("screen entry index out of range")
	}

	entryAddr := sd.base + uintptr(entryIndex*2)
	screenEntry := (*volatile.Register16)(unsafe.Pointer(entryAddr)).Get()
	
	tileIndex := int(screenEntry & 0x3FF)
	attributes := screenEntry & 0xFC00

	return tileIndex, attributes, nil
}

// ClearScreen clears all screen entries (sets all tiles to 0)
func (sd *ScreenData) ClearScreen() {
	for i := 0; i < sd.maxEntries; i++ {
		entryAddr := sd.base + uintptr(i*2)
		(*volatile.Register16)(unsafe.Pointer(entryAddr)).Set(0)
	}
}

// FillScreen fills the entire screen with the specified tile and attributes
func (sd *ScreenData) FillScreen(tileIndex int, attributes uint16) {
	screenEntry := uint16(tileIndex&0x3FF) | (attributes & 0xFC00)
	
	totalEntries := sd.width * sd.height
	if totalEntries > sd.maxEntries {
		totalEntries = sd.maxEntries
	}
	
	for i := 0; i < totalEntries; i++ {
		entryAddr := sd.base + uintptr(i*2)
		(*volatile.Register16)(unsafe.Pointer(entryAddr)).Set(screenEntry)
	}
}

// SetTileRect sets a rectangular area of tiles
func (sd *ScreenData) SetTileRect(x, y, width, height, tileIndex int, attributes uint16) {
	// Clamp to bounds
	if x < 0 {
		width += x
		x = 0
	}
	if y < 0 {
		height += y
		y = 0
	}
	if x+width > sd.width {
		width = sd.width - x
	}
	if y+height > sd.height {
		height = sd.height - y
	}

	if width <= 0 || height <= 0 {
		return
	}

	screenEntry := uint16(tileIndex&0x3FF) | (attributes & 0xFC00)
	
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			entryIndex := (y+dy)*sd.width + (x+dx)
			if entryIndex < sd.maxEntries {
				entryAddr := sd.base + uintptr(entryIndex*2)
				(*volatile.Register16)(unsafe.Pointer(entryAddr)).Set(screenEntry)
			}
		}
	}
}

// Tile attribute constants for convenience
const (
	TILE_HFLIP    = 1 << 10 // Horizontal flip
	TILE_VFLIP    = 1 << 11 // Vertical flip
	TILE_PAL_MASK = 0xF000  // Palette mask (bits 12-15)
)

// SetTilePalette creates tile attributes with the specified palette
func SetTilePalette(palette int) uint16 {
	return uint16((palette & 0xF) << 12)
}

// SetTileFlip creates tile attributes with flip flags
func SetTileFlip(hflip, vflip bool) uint16 {
	var attrs uint16
	if hflip {
		attrs |= TILE_HFLIP
	}
	if vflip {
		attrs |= TILE_VFLIP
	}
	return attrs
}

// CombineTileAttributes combines multiple tile attributes
func CombineTileAttributes(attrs ...uint16) uint16 {
	var result uint16
	for _, attr := range attrs {
		result |= attr
	}
	return result
}