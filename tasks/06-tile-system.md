# Task 06: Tile System

## Objective
Implement a comprehensive tile system for the GBA that manages character data (tiles), tile maps, and provides utilities for loading, manipulating, and rendering tile-based graphics for both backgrounds and sprites.

## Background
The GBA uses a tile-based graphics system for efficient rendering. Tiles are 8x8 pixel graphics stored in character blocks, while tile maps define which tiles appear where on screen. This system is fundamental for creating backgrounds and optimized sprite graphics.

### GBA Tile System Overview
- **Tiles**: 8x8 pixel graphics in 4-bit (16 colors) or 8-bit (256 colors) format
- **Character Blocks**: 16KB regions in VRAM storing tile data (4 blocks for BG, 2 for sprites)
- **Screen Blocks**: 2KB regions storing tile map data (32 blocks available)
- **Tile Maps**: Arrays of 16-bit values specifying tile index and attributes
- **Memory Layout**: Character data at 0x06000000, screen data overlaps at 0x06008000

### Tile Formats
- **4-bit tiles**: 32 bytes per tile (4 bits per pixel, 16 colors, uses palette)
- **8-bit tiles**: 64 bytes per tile (8 bits per pixel, 256 colors, uses palette)
- **Tile attributes**: Palette selection, flipping, priority

## Requirements

### Functional Requirements
1. Load and manage tile graphics data in VRAM character blocks
2. Create and manipulate tile maps for backgrounds
3. Support both 4-bit and 8-bit tile formats
4. Provide tile map editing and rendering utilities
5. Efficient tile data compression and storage
6. Integration with palette and VRAM systems
7. Sprite tile management for animated characters
8. Collision detection based on tile properties

### Technical Requirements
- Direct VRAM character block access
- Optimized tile loading using DMA
- Memory-efficient tile storage
- Support for all 4 character blocks (BG) and 2 sprite blocks
- Tile map scrolling and wrapping
- Integration with existing memory management

## API Design

### Constants
```go
package tiles

import "github.com/matheusmortatti/gba-go/lib/memory"

const (
    // Tile dimensions
    TILE_SIZE         = 8    // 8x8 pixels
    TILE_PIXELS       = 64   // 8 * 8
    
    // Tile data sizes
    TILE_4BPP_SIZE    = 32   // 4 bits per pixel * 64 pixels / 8 bits per byte
    TILE_8BPP_SIZE    = 64   // 8 bits per pixel * 64 pixels / 8 bits per byte
    
    // Character blocks
    CHAR_BLOCK_SIZE   = 0x4000  // 16KB per block
    CHAR_BLOCKS_BG    = 4       // Blocks 0-3 for backgrounds
    CHAR_BLOCKS_OBJ   = 2       // Blocks 4-5 for sprites
    TILES_PER_BLOCK_4BPP = CHAR_BLOCK_SIZE / TILE_4BPP_SIZE  // 512 tiles
    TILES_PER_BLOCK_8BPP = CHAR_BLOCK_SIZE / TILE_8BPP_SIZE  // 256 tiles
    
    // Screen blocks
    SCREEN_BLOCK_SIZE = 0x800   // 2KB per block
    SCREEN_BLOCKS     = 32      // Total screen blocks
    SCREEN_ENTRIES    = SCREEN_BLOCK_SIZE / 2  // 1024 16-bit entries
    
    // Screen dimensions in tiles
    SCREEN_TILES_X    = 32
    SCREEN_TILES_Y    = 32
    
    // Tile map entry flags
    TILE_HFLIP        = 0x0400  // Horizontal flip
    TILE_VFLIP        = 0x0800  // Vertical flip
    TILE_PALETTE_MASK = 0xF000  // Palette bits (4-bit mode)
    TILE_INDEX_MASK   = 0x03FF  // Tile index bits
    
    // Character block base addresses
    CHAR_BLOCK_0      = memory.VRAM_BASE + 0x0000
    CHAR_BLOCK_1      = memory.VRAM_BASE + 0x4000
    CHAR_BLOCK_2      = memory.VRAM_BASE + 0x8000
    CHAR_BLOCK_3      = memory.VRAM_BASE + 0xC000
    CHAR_BLOCK_4      = memory.VRAM_BASE + 0x10000  // Sprite tiles
    CHAR_BLOCK_5      = memory.VRAM_BASE + 0x14000  // Sprite tiles
    
    // Screen block base
    SCREEN_BLOCK_BASE = memory.VRAM_BASE + 0x8000
)
```

### Data Structures
```go
// Tile represents 8x8 pixel tile data
type Tile struct {
    data     []uint8  // Raw pixel data
    bpp      int      // Bits per pixel (4 or 8)
    palette  int      // Palette index (4-bit mode only)
    metadata map[string]interface{} // Custom properties
}

// TileMap represents a grid of tile references
type TileMap struct {
    width     int
    height    int
    entries   []uint16  // Tile map entries
    scrollX   int
    scrollY   int
    wrapX     bool
    wrapY     bool
}

// CharacterBlock manages tiles in a VRAM character block
type CharacterBlock struct {
    blockId     int
    baseAddr    uintptr
    bpp         int
    tiles       []*Tile
    nextFreeId  int
    isSprite    bool  // true for sprite blocks (4-5)
}

// ScreenBlock manages a tile map in VRAM
type ScreenBlock struct {
    blockId   int
    baseAddr  uintptr
    tileMap   *TileMap
    charBlock int  // Associated character block
}

// TileManager coordinates all tile operations
type TileManager struct {
    charBlocks   [6]*CharacterBlock  // 4 BG + 2 sprite blocks
    screenBlocks [32]*ScreenBlock
    tileCache    map[string]*Tile    // Named tile cache
}

// TileMapEntry represents a single tile map entry
type TileMapEntry struct {
    tileIndex uint16
    palette   uint8
    hFlip     bool
    vFlip     bool
    priority  bool
}
```

### Core Functions
```go
// Tile manager
func NewTileManager() *TileManager
func (tm *TileManager) GetCharacterBlock(blockId int) *CharacterBlock
func (tm *TileManager) GetScreenBlock(blockId int) *ScreenBlock
func (tm *TileManager) LoadTileSet(name string, data []byte, bpp int) error

// Character block operations
func (cb *CharacterBlock) LoadTile(tileId int, tile *Tile) error
func (cb *CharacterBlock) GetTile(tileId int) *Tile
func (cb *CharacterBlock) AllocateTile() int
func (cb *CharacterBlock) FreeTile(tileId int)
func (cb *CharacterBlock) Clear()
func (cb *CharacterBlock) GetFreeSpace() int

// Tile operations
func NewTile(data []uint8, bpp int) *Tile
func (t *Tile) GetPixel(x, y int) uint8
func (t *Tile) SetPixel(x, y int, colorIndex uint8)
func (t *Tile) FlipHorizontal() *Tile
func (t *Tile) FlipVertical() *Tile
func (t *Tile) Rotate90() *Tile
func (t *Tile) ToBytes() []uint8

// Tile map operations
func NewTileMap(width, height int) *TileMap
func (tm *TileMap) SetTile(x, y int, entry TileMapEntry)
func (tm *TileMap) GetTile(x, y int) TileMapEntry
func (tm *TileMap) Fill(entry TileMapEntry)
func (tm *TileMap) FillRegion(x, y, width, height int, entry TileMapEntry)
func (tm *TileMap) Copy(src *TileMap, srcX, srcY, dstX, dstY, width, height int)
func (tm *TileMap) Scroll(dx, dy int)
func (tm *TileMap) SetScrolling(x, y int)

// Screen block operations
func (sb *ScreenBlock) LoadTileMap(tileMap *TileMap)
func (sb *ScreenBlock) UpdateRegion(x, y, width, height int)
func (sb *ScreenBlock) SetCharacterBlock(blockId int)
func (sb *ScreenBlock) Clear()

// Utility functions
func LoadTilesFromImage(imageData []byte) ([]*Tile, error)
func ConvertTile4To8(tile4bpp *Tile) *Tile
func ConvertTile8To4(tile8bpp *Tile, palette int) *Tile
func CreateSolidColorTile(colorIndex uint8, bpp int) *Tile
func CreateCheckerboardTile(color1, color2 uint8, bpp int) *Tile

// File I/O
func LoadTileMapFromFile(filename string) (*TileMap, error)
func SaveTileMapToFile(tileMap *TileMap, filename string) error
func ExportTileSetToImage(charBlock *CharacterBlock, filename string) error
```

## Implementation Details

### Step 1: Tile Data Management
Create `lib/tiles/tile.go`:
- Tile struct with pixel manipulation
- Support for 4-bit and 8-bit formats
- Tile transformation utilities (flip, rotate)

### Step 2: Character Block Management
Create `lib/tiles/charblock.go`:
- Direct VRAM character block access
- Tile allocation and deallocation
- Efficient tile loading using DMA

### Step 3: Tile Map System
Create `lib/tiles/tilemap.go`:
- TileMap data structure and operations
- Scrolling and wrapping support
- Region-based updates for performance

### Step 4: Screen Block Management
Create `lib/tiles/screenblock.go`:
- Screen block VRAM interface
- Tile map to hardware conversion
- Efficient batch updates

### Step 5: Tile Manager Coordination
Create `lib/tiles/manager.go`:
- Central tile system coordinator
- Resource management and caching
- Integration with other systems

### Step 6: Utilities and File I/O
Create `lib/tiles/utils.go`:
- Tile creation helpers
- Image conversion utilities
- Debug and visualization tools

## Testing Strategy

### Unit Tests
```go
func TestTileCreation(t *testing.T) {
    // Create 4-bit tile
    data := make([]uint8, 32)
    for i := range data {
        data[i] = uint8(i % 16) // Fill with color indices 0-15
    }
    
    tile := NewTile(data, 4)
    assert.Equal(t, 4, tile.bpp)
    assert.Equal(t, 32, len(tile.data))
    
    // Test pixel access
    pixel := tile.GetPixel(0, 0)
    assert.Equal(t, uint8(0), pixel)
    
    pixel = tile.GetPixel(1, 0)
    assert.Equal(t, uint8(1), pixel)
}

func TestCharacterBlock(t *testing.T) {
    manager := NewTileManager()
    charBlock := manager.GetCharacterBlock(0)
    
    // Create test tile
    data := make([]uint8, 32)
    tile := NewTile(data, 4)
    
    // Load tile
    err := charBlock.LoadTile(0, tile)
    assert.NoError(t, err)
    
    // Retrieve tile
    retrievedTile := charBlock.GetTile(0)
    assert.NotNil(t, retrievedTile)
    assert.Equal(t, 4, retrievedTile.bpp)
}

func TestTileMap(t *testing.T) {
    tileMap := NewTileMap(32, 32)
    
    // Set tile entry
    entry := TileMapEntry{
        tileIndex: 5,
        palette:   2,
        hFlip:     true,
        vFlip:     false,
    }
    
    tileMap.SetTile(10, 15, entry)
    retrieved := tileMap.GetTile(10, 15)
    
    assert.Equal(t, uint16(5), retrieved.tileIndex)
    assert.Equal(t, uint8(2), retrieved.palette)
    assert.True(t, retrieved.hFlip)
    assert.False(t, retrieved.vFlip)
}

func TestTileMapScrolling(t *testing.T) {
    tileMap := NewTileMap(32, 32)
    
    // Fill with test pattern
    for y := 0; y < 32; y++ {
        for x := 0; x < 32; x++ {
            entry := TileMapEntry{tileIndex: uint16(x + y*32)}
            tileMap.SetTile(x, y, entry)
        }
    }
    
    // Test scrolling
    originalTile := tileMap.GetTile(5, 5)
    tileMap.Scroll(1, 1)
    scrolledTile := tileMap.GetTile(4, 4)
    
    assert.Equal(t, originalTile.tileIndex, scrolledTile.tileIndex)
}
```

### Integration Tests
```go
func TestTileSystemIntegration(t *testing.T) {
    manager := NewTileManager()
    
    // Load tiles into character block
    charBlock := manager.GetCharacterBlock(0)
    tiles := createTestTiles()
    for i, tile := range tiles {
        charBlock.LoadTile(i, tile)
    }
    
    // Create tile map
    tileMap := NewTileMap(32, 32)
    for y := 0; y < 32; y++ {
        for x := 0; x < 32; x++ {
            entry := TileMapEntry{tileIndex: uint16((x + y) % len(tiles))}
            tileMap.SetTile(x, y, entry)
        }
    }
    
    // Load into screen block
    screenBlock := manager.GetScreenBlock(0)
    screenBlock.SetCharacterBlock(0)
    screenBlock.LoadTileMap(tileMap)
    
    // Verify data was written to VRAM
    // (This would require actual hardware/emulator testing)
}
```

## Example Program
```go
package main

import (
    "github.com/matheusmortatti/gba-go/lib/tiles"
    "github.com/matheusmortatti/gba-go/lib/memory"
    "github.com/matheusmortatti/gba-go/lib/palette"
    "github.com/matheusmortatti/gba-go/lib/registers"
    "github.com/matheusmortatti/gba-go/lib/video"
    "github.com/matheusmortatti/gba-go/lib/input"
)

func main() {
    // Initialize tile system
    tileManager := tiles.NewTileManager()
    paletteManager := palette.NewPaletteManager()
    
    // Set video mode 0 (tile mode)
    registers.Lcd.DISPCNT.SetBits(memory.MODE_0)
    registers.Lcd.DISPCNT.SetBits(1 << 8)  // Enable BG0
    
    // Configure background 0
    registers.Lcd.BG0CNT.SetBits(0 << 2)   // Character base block 0
    registers.Lcd.BG0CNT.SetBits(8 << 8)   // Screen base block 8
    registers.Lcd.BG0CNT.SetBits(0 << 14)  // Size 256x256
    
    // Set up palette
    setupPalette(paletteManager)
    
    // Create and load tiles
    createTileSet(tileManager)
    
    // Create scrolling background
    createScrollingBackground(tileManager)
    
    scrollX, scrollY := 0, 0
    
    for {
        video.VSync()
        input.Poll()
        
        // Handle scrolling input
        if input.BtnDown(input.KeyLeft) {
            scrollX -= 2
        }
        if input.BtnDown(input.KeyRight) {
            scrollX += 2
        }
        if input.BtnDown(input.KeyUp) {
            scrollY -= 2
        }
        if input.BtnDown(input.KeyDown) {
            scrollY += 2
        }
        
        // Update background scrolling
        registers.Lcd.BG0HOFS.Set(uint16(scrollX))
        registers.Lcd.BG0VOFS.Set(uint16(scrollY))
        
        // Demonstrate dynamic tile updates
        if input.BtnClicked(input.KeyA) {
            updateTileAnimation(tileManager)
        }
    }
}

func setupPalette(manager *palette.PaletteManager) {
    // Create a simple palette for tiles
    bgPalette := &palette.Palette16{}
    bgPalette.SetColor(0, palette.BLACK)      // Background
    bgPalette.SetColor(1, palette.WHITE)      // White
    bgPalette.SetColor(2, palette.RED)        // Red
    bgPalette.SetColor(3, palette.GREEN)      // Green
    bgPalette.SetColor(4, palette.BLUE)       // Blue
    bgPalette.SetColor(5, palette.YELLOW)     // Yellow
    bgPalette.SetColor(6, palette.CYAN)       // Cyan
    bgPalette.SetColor(7, palette.MAGENTA)    // Magenta
    bgPalette.SetColor(8, palette.GRAY)       // Gray
    bgPalette.SetColor(9, palette.LIGHT_GRAY) // Light gray
    bgPalette.SetColor(10, palette.DARK_GRAY) // Dark gray
    
    manager.LoadBGPalette16(0, bgPalette)
}

func createTileSet(manager *tiles.TileManager) {
    charBlock := manager.GetCharacterBlock(0)
    
    // Create solid color tiles
    colors := []uint8{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
    for i, color := range colors {
        tile := tiles.CreateSolidColorTile(color, 4)
        charBlock.LoadTile(i, tile)
    }
    
    // Create patterned tiles
    checkerTile := tiles.CreateCheckerboardTile(1, 2, 4)
    charBlock.LoadTile(11, checkerTile)
    
    // Create gradient tile
    gradientTile := createGradientTile()
    charBlock.LoadTile(12, gradientTile)
    
    // Create border tiles
    borderTiles := createBorderTiles()
    for i, tile := range borderTiles {
        charBlock.LoadTile(13+i, tile)
    }
}

func createScrollingBackground(manager *tiles.TileManager) {
    screenBlock := manager.GetScreenBlock(8)
    screenBlock.SetCharacterBlock(0)
    
    // Create a tile map with pattern
    tileMap := tiles.NewTileMap(32, 32)
    
    // Create a landscape pattern
    for y := 0; y < 32; y++ {
        for x := 0; x < 32; x++ {
            var tileIndex uint16
            
            if y < 5 {
                // Sky
                tileIndex = 4 // Blue
            } else if y < 8 {
                // Clouds (checkerboard pattern)
                if (x+y)%3 == 0 {
                    tileIndex = 11 // Checkerboard
                } else {
                    tileIndex = 4 // Blue
                }
            } else if y < 20 {
                // Ground
                if (x+y)%4 == 0 {
                    tileIndex = 3 // Green
                } else if (x+y)%7 == 0 {
                    tileIndex = 7 // Magenta (flowers)
                } else {
                    tileIndex = 2 // Red (dirt)
                }
            } else {
                // Underground
                tileIndex = 10 // Dark gray
            }
            
            entry := tiles.TileMapEntry{
                tileIndex: tileIndex,
                palette:   0,
                hFlip:     false,
                vFlip:     false,
            }
            
            tileMap.SetTile(x, y, entry)
        }
    }
    
    // Add some random details
    addRandomDetails(tileMap)
    
    screenBlock.LoadTileMap(tileMap)
}

func createGradientTile() *tiles.Tile {
    data := make([]uint8, 32) // 4-bit tile
    
    for y := 0; y < 8; y++ {
        for x := 0; x < 8; x++ {
            // Create gradient from corner
            distance := int(math.Sqrt(float64(x*x + y*y)))
            color := uint8(distance % 8 + 1) // Colors 1-8
            
            byteIndex := (y*8 + x) / 2
            if x%2 == 0 {
                data[byteIndex] = (data[byteIndex] & 0xF0) | color
            } else {
                data[byteIndex] = (data[byteIndex] & 0x0F) | (color << 4)
            }
        }
    }
    
    return tiles.NewTile(data, 4)
}

func createBorderTiles() []*tiles.Tile {
    // Create tiles for borders and edges
    borderTiles := make([]*tiles.Tile, 8)
    
    // Top-left corner
    borderTiles[0] = createBorderTile([][]uint8{
        {1,1,1,1,1,1,1,1},
        {1,0,0,0,0,0,0,0},
        {1,0,0,0,0,0,0,0},
        {1,0,0,0,0,0,0,0},
        {1,0,0,0,0,0,0,0},
        {1,0,0,0,0,0,0,0},
        {1,0,0,0,0,0,0,0},
        {1,0,0,0,0,0,0,0},
    })
    
    // Top edge
    borderTiles[1] = createBorderTile([][]uint8{
        {1,1,1,1,1,1,1,1},
        {0,0,0,0,0,0,0,0},
        {0,0,0,0,0,0,0,0},
        {0,0,0,0,0,0,0,0},
        {0,0,0,0,0,0,0,0},
        {0,0,0,0,0,0,0,0},
        {0,0,0,0,0,0,0,0},
        {0,0,0,0,0,0,0,0},
    })
    
    // Add more border variations...
    
    return borderTiles
}

func createBorderTile(pattern [][]uint8) *tiles.Tile {
    data := make([]uint8, 32) // 4-bit tile
    
    for y := 0; y < 8; y++ {
        for x := 0; x < 8; x++ {
            color := pattern[y][x]
            byteIndex := (y*8 + x) / 2
            if x%2 == 0 {
                data[byteIndex] = (data[byteIndex] & 0xF0) | color
            } else {
                data[byteIndex] = (data[byteIndex] & 0x0F) | (color << 4)
            }
        }
    }
    
    return tiles.NewTile(data, 4)
}

func addRandomDetails(tileMap *tiles.TileMap) {
    // Add some random decorative elements
    for i := 0; i < 20; i++ {
        x := rand.Intn(32)
        y := 8 + rand.Intn(12) // Only in ground area
        
        entry := tiles.TileMapEntry{
            tileIndex: uint16(12 + rand.Intn(4)), // Random decorative tiles
            palette:   0,
            hFlip:     rand.Intn(2) == 1,
            vFlip:     rand.Intn(2) == 1,
        }
        
        tileMap.SetTile(x, y, entry)
    }
}

func updateTileAnimation(manager *tiles.TileManager) {
    // Demonstrate dynamic tile updates
    charBlock := manager.GetCharacterBlock(0)
    
    // Create animated water tile
    animatedTile := createAnimatedWaterTile()
    charBlock.LoadTile(5, animatedTile) // Update tile 5
}

func createAnimatedWaterTile() *tiles.Tile {
    data := make([]uint8, 32)
    
    // Create wavy pattern based on current time
    frame := int(time.Now().UnixMilli() / 100) % 8
    
    for y := 0; y < 8; y++ {
        for x := 0; x < 8; x++ {
            // Create wave effect
            wave := int(math.Sin(float64(x+frame)*0.5) * 2)
            if y >= 4+wave && y <= 6+wave {
                color := uint8(6) // Cyan
            } else {
                color := uint8(4) // Blue
            }
            
            byteIndex := (y*8 + x) / 2
            if x%2 == 0 {
                data[byteIndex] = (data[byteIndex] & 0xF0) | color
            } else {
                data[byteIndex] = (data[byteIndex] & 0x0F) | (color << 4)
            }
        }
    }
    
    return tiles.NewTile(data, 4)
}
```

## File Structure
```
lib/tiles/
├── tile.go         // Individual tile data and operations
├── charblock.go    // Character block management
├── tilemap.go      // Tile map data structure
├── screenblock.go  // Screen block VRAM interface
├── manager.go      // Central tile system coordinator
└── utils.go        // Utilities and file I/O
```

## Integration Points
- Use memory constants from Task 01
- Integrate with palette system from Task 02
- Coordinate with VRAM management from Task 03
- Support sprite system from Task 04
- Enable background system from Task 07

## Resources
- [GBATEK Tile Data](https://problemkaputt.de/gbatek.htm#gbavideobgscreendisplay)
- [Tonc Tile System](https://www.coranac.com/tonc/text/regbg.htm)
- [GBA Tile Programming Guide](https://www.cs.rit.edu/~tjh8300/CowBite/CowBiteSpec.htm#Tiles)

## Success Criteria
- All character blocks can be managed independently
- Tile maps render correctly with scrolling
- Both 4-bit and 8-bit tiles work properly
- Dynamic tile updates happen without visual artifacts
- Memory usage is optimized for tile storage
- Example program demonstrates all tile features
- Performance suitable for real-time tile updates
- Comprehensive test coverage (>90%)