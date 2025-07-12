# Task 01: Memory Constants and Layout

## Objective
Create a comprehensive memory management module that defines all GBA memory regions, constants, and utility functions for safe memory access. This forms the foundation for all other graphics and hardware operations.

## Background
The Game Boy Advance has a specific memory layout with different regions for VRAM, OAM, palettes, and other hardware components. Understanding and properly defining these memory regions is crucial for all GBA development.

### GBA Memory Map
- **BIOS ROM**: 0x00000000 - 0x00003FFF (16KB)
- **Work RAM**: 0x02000000 - 0x0203FFFF (256KB, fast)
- **Work RAM**: 0x03000000 - 0x03007FFF (32KB, faster)
- **I/O Registers**: 0x04000000 - 0x040003FE
- **Palette RAM**: 0x05000000 - 0x050003FF (1KB)
- **VRAM**: 0x06000000 - 0x06017FFF (96KB)
- **OAM**: 0x07000000 - 0x070003FF (1KB)
- **Game Pak ROM**: 0x08000000 - 0x09FFFFFF (32MB)

## Requirements

### Functional Requirements
1. Define all memory base addresses and sizes as constants
2. Provide type-safe memory region structs
3. Create utility functions for memory bounds checking
4. Define video mode constants and screen dimensions
5. Provide color format conversion utilities

### Technical Requirements
- Use `unsafe.Pointer` for hardware memory access
- Leverage `runtime/volatile` for hardware registers
- Ensure memory-aligned access where required
- Support both 16-bit and 32-bit memory operations

## API Design

### Constants
```go
package memory

const (
    // Memory Base Addresses
    BIOS_BASE       = 0x00000000
    EWRAM_BASE      = 0x02000000
    IWRAM_BASE      = 0x03000000
    IO_BASE         = 0x04000000
    PALETTE_BASE    = 0x05000000
    VRAM_BASE       = 0x06000000
    OAM_BASE        = 0x07000000
    ROM_BASE        = 0x08000000

    // Memory Sizes
    BIOS_SIZE       = 0x4000     // 16KB
    EWRAM_SIZE      = 0x40000    // 256KB
    IWRAM_SIZE      = 0x8000     // 32KB
    PALETTE_SIZE    = 0x400      // 1KB
    VRAM_SIZE       = 0x18000    // 96KB
    OAM_SIZE        = 0x400      // 1KB

    // Screen Constants
    SCREEN_WIDTH    = 240
    SCREEN_HEIGHT   = 160
    SCREEN_PIXELS   = SCREEN_WIDTH * SCREEN_HEIGHT

    // Video Modes
    MODE_0          = 0  // Text mode, 4 backgrounds
    MODE_1          = 1  // Text + Affine, 3 backgrounds
    MODE_2          = 2  // Affine mode, 2 backgrounds
    MODE_3          = 3  // Bitmap 16-bit, 1 background
    MODE_4          = 4  // Bitmap 8-bit, 1 background
    MODE_5          = 5  // Bitmap 16-bit small, 1 background

    // Tile Constants
    TILE_SIZE       = 8   // 8x8 pixels
    TILE_SIZE_BYTES = 32  // 4 bits per pixel * 64 pixels / 2
    TILES_PER_ROW   = 32
    TILES_PER_COL   = 32

    // Color Constants
    COLOR_DEPTH     = 15  // 15-bit color (32768 colors)
    COLORS_PER_PALETTE = 256
    BG_PALETTES     = 16
    OBJ_PALETTES    = 16
)
```

### Data Structures
```go
// Color represents a 15-bit GBA color
type Color uint16

// RGB constructs a Color from 5-bit RGB components
func RGB(r, g, b uint8) Color

// Palette represents a color palette
type Palette [COLORS_PER_PALETTE]Color

// MemoryRegion provides safe access to a memory region
type MemoryRegion struct {
    base uintptr
    size uintptr
}

// VRAM regions for different video modes
type VRAMLayout struct {
    Mode0 *MemoryRegion // Character data
    Mode3 *MemoryRegion // 16-bit bitmap
    Mode4 *MemoryRegion // 8-bit bitmap
    Mode5 *MemoryRegion // 16-bit small bitmap
}
```

### Core Functions
```go
// Memory region access
func GetVRAM() *MemoryRegion
func GetOAM() *MemoryRegion
func GetPaletteRAM() *MemoryRegion

// Bounds checking
func (r *MemoryRegion) InBounds(offset uintptr) bool
func (r *MemoryRegion) Size() uintptr

// Memory operations
func (r *MemoryRegion) Read16(offset uintptr) uint16
func (r *MemoryRegion) Write16(offset uintptr, value uint16)
func (r *MemoryRegion) Read32(offset uintptr) uint32
func (r *MemoryRegion) Write32(offset uintptr, value uint32)

// Color utilities
func (c Color) R() uint8 // Extract red component (0-31)
func (c Color) G() uint8 // Extract green component (0-31)
func (c Color) B() uint8 // Extract blue component (0-31)
func RGB15(r, g, b uint8) Color // Create color from 5-bit components
```

## Implementation Details

### Step 1: Create Memory Constants
Create `lib/memory/constants.go` with all memory layout definitions, ensuring proper alignment and size calculations.

### Step 2: Implement Color System
- Define 15-bit color type with RGB extraction/construction
- Implement color conversion utilities
- Create common color constants (BLACK, WHITE, RED, etc.)

### Step 3: Memory Region Abstraction
- Create MemoryRegion struct with bounds checking
- Implement safe read/write operations
- Add debug assertions for development builds

### Step 4: VRAM Layout Management
- Define VRAM regions for different video modes
- Create mode-specific memory layout helpers
- Implement double-buffering support for bitmap modes

### Step 5: Integration Points
- Ensure compatibility with existing register definitions
- Provide migration path from direct pointer access
- Create performance-optimized paths for critical operations

## Testing Strategy

### Unit Tests
```go
func TestMemoryConstants(t *testing.T) {
    // Verify memory layout constants
    assert.Equal(t, uintptr(0x06000000), VRAM_BASE)
    assert.Equal(t, uintptr(0x18000), VRAM_SIZE)
}

func TestColorOperations(t *testing.T) {
    // Test RGB color construction and extraction
    red := RGB15(31, 0, 0)
    assert.Equal(t, uint8(31), red.R())
    assert.Equal(t, uint8(0), red.G())
    assert.Equal(t, uint8(0), red.B())
}

func TestMemoryBounds(t *testing.T) {
    vram := GetVRAM()
    assert.True(t, vram.InBounds(0))
    assert.True(t, vram.InBounds(VRAM_SIZE-1))
    assert.False(t, vram.InBounds(VRAM_SIZE))
}
```

### Integration Tests
- Test memory region access with actual hardware registers
- Verify color display on screen
- Test bounds checking prevents crashes

## Example Program
```go
package main

import (
    "github.com/matheusmortatti/gba-go/lib/memory"
    "github.com/matheusmortatti/gba-go/lib/registers"
)

func main() {
    // Set video mode 3 (16-bit bitmap)
    registers.Lcd.DISPCNT.SetBits(memory.MODE_3)
    registers.Lcd.DISPCNT.SetBits(1 << 10) // Enable BG2

    // Get VRAM access
    vram := memory.GetVRAM()

    // Draw colored pixels using safe memory access
    colors := []memory.Color{
        memory.RGB15(31, 0, 0),  // Red
        memory.RGB15(0, 31, 0),  // Green
        memory.RGB15(0, 0, 31),  // Blue
        memory.RGB15(31, 31, 0), // Yellow
    }

    // Draw 4 colored squares
    for i, color := range colors {
        x := (i % 2) * 120
        y := (i / 2) * 80
        
        for dy := 0; dy < 80; dy++ {
            for dx := 0; dx < 120; dx++ {
                offset := uintptr(((y+dy)*memory.SCREEN_WIDTH + (x+dx)) * 2)
                if vram.InBounds(offset) {
                    vram.Write16(offset, uint16(color))
                }
            }
        }
    }

    // Main loop
    for {
        // Wait for next frame
    }
}
```

## File Structure
```
lib/memory/
├── constants.go     // Memory layout and video mode constants
├── color.go         // Color type and utilities
├── regions.go       // Memory region abstraction
└── vram.go         // VRAM-specific utilities
```

## Integration with Existing Code
- Update `examples/main.go` to use new memory constants
- Modify existing VRAM access to use safe memory operations
- Ensure backwards compatibility during transition

## Resources
- [GBATEK Memory Map](https://problemkaputt.de/gbatek.htm#gbamemorymap)
- [Tonc GBA Programming Guide](https://www.coranac.com/tonc/text/hardware.htm)
- [GBA Hardware Manual](https://www.akkit.org/info/gbatek.htm)

## Success Criteria
- All memory constants correctly defined and tested
- Safe memory access functions working without crashes
- Color system producing correct visual output
- Integration with existing codebase complete
- Example program demonstrates functionality
- Comprehensive test coverage (>90%)