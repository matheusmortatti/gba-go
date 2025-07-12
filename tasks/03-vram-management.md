# Task 03: VRAM Management System

## Objective
Create a comprehensive Video RAM (VRAM) management system that provides efficient, safe access to graphics memory across all GBA video modes, with optimized operations for bitmap graphics, tile rendering, and double buffering.

## Background
The GBA's 96KB VRAM is organized differently depending on the video mode. Understanding and properly managing VRAM layout is crucial for graphics performance and avoiding visual artifacts.

### VRAM Layout by Video Mode
- **Mode 0-2 (Tile modes)**: Character data, tile maps, sprite tiles
- **Mode 3**: Single 240x160 16-bit bitmap (77KB used)
- **Mode 4**: Two 240x160 8-bit bitmaps with palette indices (38KB each)
- **Mode 5**: Two 160x128 16-bit bitmaps (40KB each)

### VRAM Memory Map
- **Total VRAM**: 0x06000000 - 0x06017FFF (96KB)
- **Character Data**: 0x06000000 - 0x0600FFFF (64KB for tiles)
- **Screen Data**: 0x06008000 - 0x0600FFFF (32KB for tile maps)
- **Bitmap Data**: 0x06000000 onwards (size depends on mode)

## Requirements

### Functional Requirements
1. Support all video modes with appropriate VRAM layouts
2. Provide fast pixel plotting and bitmap operations
3. Implement double buffering for smooth animation
4. Safe bounds checking and memory access
5. Efficient bulk operations (clear, copy, fill)
6. Integration with palette and tile systems

### Technical Requirements
- Mode-specific VRAM addressing
- 16-bit and 8-bit pixel access
- DMA-compatible memory operations
- Memory-aligned access for performance
- Minimal overhead for critical operations

## API Design

### Constants
```go
package vram

import "github.com/matheusmortatti/gba-go/lib/memory"

const (
    // VRAM layout constants
    VRAM_BASE           = memory.VRAM_BASE
    VRAM_SIZE           = memory.VRAM_SIZE
    
    // Video mode specific sizes
    MODE3_FRAME_SIZE    = memory.SCREEN_WIDTH * memory.SCREEN_HEIGHT * 2  // 16-bit
    MODE4_FRAME_SIZE    = memory.SCREEN_WIDTH * memory.SCREEN_HEIGHT      // 8-bit
    MODE5_FRAME_SIZE    = 160 * 128 * 2                                   // 16-bit small
    
    // Buffer addresses for bitmap modes
    FRAME0_ADDR         = VRAM_BASE
    FRAME1_ADDR         = VRAM_BASE + 0xA000  // Frame 1 for modes 4&5
    
    // Character and screen base addresses (tile modes)
    CHAR_BASE_SIZE      = 0x4000  // 16KB per character base block
    SCREEN_BASE_SIZE    = 0x800   // 2KB per screen base block
    MAX_CHAR_BLOCKS     = 4       // 0-3
    MAX_SCREEN_BLOCKS   = 32      // 0-31
    
    // Tile constants
    TILE_4BPP_SIZE      = 32      // 4 bits per pixel
    TILE_8BPP_SIZE      = 64      // 8 bits per pixel
    TILES_PER_CHAR_BLOCK = CHAR_BASE_SIZE / TILE_4BPP_SIZE
)
```

### Data Structures
```go
// VRAMManager provides mode-specific VRAM access
type VRAMManager struct {
    mode        int
    currentPage int  // For double buffering
    frameSize   int
}

// BitmapBuffer represents a framebuffer for bitmap modes
type BitmapBuffer struct {
    base   uintptr
    width  int
    height int
    bpp    int  // bits per pixel (8 or 16)
}

// TileData represents character data in VRAM
type TileData struct {
    base      uintptr
    charBlock int
    bpp       int
}

// ScreenData represents tile map data
type ScreenData struct {
    base        uintptr
    screenBlock int
    width       int  // in tiles
    height      int  // in tiles
}
```

### Core Functions
```go
// VRAM Manager
func NewVRAMManager(mode int) *VRAMManager
func (vm *VRAMManager) SetMode(mode int)
func (vm *VRAMManager) GetCurrentBuffer() *BitmapBuffer
func (vm *VRAMManager) GetBackBuffer() *BitmapBuffer
func (vm *VRAMManager) SwapBuffers()

// Bitmap operations
func (bb *BitmapBuffer) PlotPixel(x, y int, color uint16) error
func (bb *BitmapBuffer) GetPixel(x, y int) (uint16, error)
func (bb *BitmapBuffer) Clear(color uint16)
func (bb *BitmapBuffer) FillRect(x, y, width, height int, color uint16)
func (bb *BitmapBuffer) DrawLine(x1, y1, x2, y2 int, color uint16)
func (bb *BitmapBuffer) CopyFrom(src *BitmapBuffer, srcX, srcY, dstX, dstY, width, height int)

// Fast bulk operations
func (bb *BitmapBuffer) FastClear(color uint16)
func (bb *BitmapBuffer) FastFill(x, y, width, height int, color uint16)
func (bb *BitmapBuffer) FastCopy(src *BitmapBuffer)

// Tile data access
func GetCharacterData(charBlock int, bpp int) *TileData
func GetScreenData(screenBlock int, width, height int) *ScreenData
func (td *TileData) LoadTile(tileIndex int, data []uint8) error
func (sd *ScreenData) SetTile(x, y, tileIndex int, attributes uint16) error

// Memory utilities
func VRAMAddr(offset uintptr) uintptr
func InVRAMBounds(addr uintptr) bool
func AlignAddress(addr uintptr, alignment int) uintptr
```

## Implementation Details

### Step 1: VRAM Manager Foundation
Create `lib/vram/manager.go`:
- Mode-specific VRAM layout management
- Buffer allocation and tracking
- Safe memory access validation

### Step 2: Bitmap Buffer Implementation
Create `lib/vram/bitmap.go`:
- BitmapBuffer struct with pixel operations
- Bounds checking for all operations
- Support for both 8-bit and 16-bit pixels

### Step 3: Fast Operations
Create `lib/vram/fast.go`:
- DMA-optimized bulk operations
- Memory-aligned operations for performance
- Assembly optimizations where needed

### Step 4: Tile Data Management
Create `lib/vram/tiles.go`:
- Character block management
- Screen block management
- Tile loading and manipulation utilities

### Step 5: Double Buffering
Create `lib/vram/buffering.go`:
- Page flipping for smooth animation
- Buffer synchronization
- VSYNC integration

### Step 6: Utility Functions
Create `lib/vram/utils.go`:
- Memory debugging utilities
- VRAM visualization helpers
- Performance profiling tools

## Testing Strategy

### Unit Tests
```go
func TestVRAMManager(t *testing.T) {
    // Test mode switching
    vm := NewVRAMManager(memory.MODE_3)
    assert.Equal(t, memory.MODE_3, vm.mode)
    
    buffer := vm.GetCurrentBuffer()
    assert.NotNil(t, buffer)
    assert.Equal(t, memory.SCREEN_WIDTH, buffer.width)
    assert.Equal(t, memory.SCREEN_HEIGHT, buffer.height)
    assert.Equal(t, 16, buffer.bpp)
}

func TestBitmapOperations(t *testing.T) {
    vm := NewVRAMManager(memory.MODE_3)
    buffer := vm.GetCurrentBuffer()
    
    // Test pixel plotting
    err := buffer.PlotPixel(100, 80, 0x7FFF) // White
    assert.NoError(t, err)
    
    pixel, err := buffer.GetPixel(100, 80)
    assert.NoError(t, err)
    assert.Equal(t, uint16(0x7FFF), pixel)
    
    // Test bounds checking
    err = buffer.PlotPixel(300, 200, 0x7FFF) // Out of bounds
    assert.Error(t, err)
}

func TestDoubleBuffering(t *testing.T) {
    vm := NewVRAMManager(memory.MODE_4)
    
    front := vm.GetCurrentBuffer()
    back := vm.GetBackBuffer()
    
    // Buffers should be different
    assert.NotEqual(t, front.base, back.base)
    
    // Test buffer swapping
    vm.SwapBuffers()
    newFront := vm.GetCurrentBuffer()
    assert.Equal(t, back.base, newFront.base)
}
```

### Performance Tests
```go
func BenchmarkPixelPlot(b *testing.B) {
    vm := NewVRAMManager(memory.MODE_3)
    buffer := vm.GetCurrentBuffer()
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        x := i % memory.SCREEN_WIDTH
        y := (i / memory.SCREEN_WIDTH) % memory.SCREEN_HEIGHT
        buffer.PlotPixel(x, y, uint16(i))
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
```

## Example Program
```go
package main

import (
    "github.com/matheusmortatti/gba-go/lib/vram"
    "github.com/matheusmortatti/gba-go/lib/memory"
    "github.com/matheusmortatti/gba-go/lib/palette"
    "github.com/matheusmortatti/gba-go/lib/registers"
    "github.com/matheusmortatti/gba-go/lib/video"
    "github.com/matheusmortatti/gba-go/lib/input"
)

func main() {
    // Initialize VRAM manager for Mode 4 (8-bit bitmap with double buffering)
    vramManager := vram.NewVRAMManager(memory.MODE_4)
    
    // Set up video mode
    registers.Lcd.DISPCNT.SetBits(memory.MODE_4)
    registers.Lcd.DISPCNT.SetBits(1 << 10) // Enable BG2
    
    // Set up palette
    paletteManager := palette.NewPaletteManager()
    rainbowPalette := createSimplePalette()
    paletteManager.LoadBGPalette256(rainbowPalette)
    
    frame := 0
    
    for {
        video.VSync()
        input.Poll()
        
        // Get back buffer for drawing
        backBuffer := vramManager.GetBackBuffer()
        
        // Clear back buffer
        backBuffer.Clear(0) // Black background
        
        // Draw animated content
        drawAnimatedScene(backBuffer, frame)
        
        // Handle input
        if input.BtnDown(input.KeyA) {
            drawPlayerSprite(backBuffer, frame)
        }
        
        // Swap buffers for smooth animation
        vramManager.SwapBuffers()
        
        // Update display control to show new buffer
        if vramManager.currentPage == 1 {
            registers.Lcd.DISPCNT.SetBits(1 << 4) // Show page 1
        } else {
            registers.Lcd.DISPCNT.ClearBits(1 << 4) // Show page 0
        }
        
        frame++
    }
}

func drawAnimatedScene(buffer *vram.BitmapBuffer, frame int) {
    // Draw moving circles
    for i := 0; i < 5; i++ {
        x := (frame*2 + i*40) % memory.SCREEN_WIDTH
        y := 80 + int(30*math.Sin(float64(frame+i*20)*0.1))
        color := uint8(i + 1) // Palette color index
        
        drawCircle(buffer, x, y, 15, color)
    }
    
    // Draw scrolling background pattern
    drawPattern(buffer, frame)
}

func drawCircle(buffer *vram.BitmapBuffer, centerX, centerY, radius int, color uint8) {
    for y := centerY - radius; y <= centerY + radius; y++ {
        for x := centerX - radius; x <= centerX + radius; x++ {
            dx := x - centerX
            dy := y - centerY
            if dx*dx + dy*dy <= radius*radius {
                buffer.PlotPixel(x, y, uint16(color))
            }
        }
    }
}

func drawPattern(buffer *vram.BitmapBuffer, frame int) {
    for y := 0; y < memory.SCREEN_HEIGHT; y += 8 {
        for x := 0; x < memory.SCREEN_WIDTH; x += 8 {
            color := uint8(((x + y + frame) / 16) % 16)
            buffer.FillRect(x, y, 8, 8, uint16(color))
        }
    }
}

func drawPlayerSprite(buffer *vram.BitmapBuffer, frame int) {
    // Simple animated player sprite
    spriteX := memory.SCREEN_WIDTH / 2
    spriteY := memory.SCREEN_HEIGHT / 2
    spriteColor := uint8(15) // Bright color
    
    // Draw a simple character (8x8 sprite)
    spriteData := [][]uint8{
        {0,0,1,1,1,1,0,0},
        {0,1,1,1,1,1,1,0},
        {1,1,2,1,1,2,1,1},
        {1,1,1,1,1,1,1,1},
        {1,2,1,1,1,1,2,1},
        {1,1,2,2,2,2,1,1},
        {0,1,1,1,1,1,1,0},
        {0,0,1,1,1,1,0,0},
    }
    
    for y := 0; y < 8; y++ {
        for x := 0; x < 8; x++ {
            if spriteData[y][x] > 0 {
                color := spriteColor + spriteData[y][x] - 1
                buffer.PlotPixel(spriteX+x-4, spriteY+y-4, uint16(color))
            }
        }
    }
}

func createSimplePalette() *palette.Palette256 {
    pal := &palette.Palette256{}
    
    // Create a simple color palette
    for i := 0; i < 256; i++ {
        r := uint8((i * 31) / 255)
        g := uint8(((i * 7) % 32) * 31 / 31)
        b := uint8(((i * 3) % 32) * 31 / 31)
        color := palette.RGB15(r, g, b)
        pal.SetColor(i, color)
    }
    
    return pal
}
```

## Advanced Features

### DMA-Accelerated Operations
```go
// Fast VRAM operations using DMA
func (bb *BitmapBuffer) DMAFill(x, y, width, height int, color uint16) {
    // Use DMA3 for fast memory fills
    if width*height < 32 {
        // Use CPU for small operations
        bb.FillRect(x, y, width, height, color)
        return
    }
    
    // Set up DMA for bulk fill operation
    startAddr := bb.getPixelAddr(x, y)
    wordCount := (width * height * bb.bpp) / 16 // 16-bit words
    
    dma.Fill32(startAddr, uint32(color|(color<<16)), wordCount)
}
```

### Memory-Mapped Access Optimization
```go
// Direct memory access for critical performance
func (bb *BitmapBuffer) getPixelAddr(x, y int) uintptr {
    if bb.bpp == 16 {
        return bb.base + uintptr((y*bb.width+x)*2)
    } else {
        return bb.base + uintptr(y*bb.width+x)
    }
}

func (bb *BitmapBuffer) PlotPixelFast(x, y int, color uint16) {
    addr := bb.getPixelAddr(x, y)
    if bb.bpp == 16 {
        *(*volatile.Register16)(unsafe.Pointer(addr)) = volatile.Register16(color)
    } else {
        *(*volatile.Register8)(unsafe.Pointer(addr)) = volatile.Register8(color)
    }
}
```

## File Structure
```
lib/vram/
├── manager.go       // VRAM manager and mode handling
├── bitmap.go        // Bitmap buffer operations
├── fast.go          // DMA-optimized operations
├── tiles.go         // Tile data management
├── buffering.go     // Double buffering support
└── utils.go         // Utility functions and debugging
```

## Integration Points
- Use memory constants from Task 01
- Integrate with palette system from Task 02
- Prepare for sprite system (Task 04)
- Support tile system requirements (Task 06)

## Resources
- [GBATEK VRAM Memory](https://problemkaputt.de/gbatek.htm#gbavideomemory)
- [Tonc Video Modes](https://www.coranac.com/tonc/text/video.htm)
- [GBA Video Programming Guide](https://www.cs.rit.edu/~tjh8300/CowBite/CowBiteSpec.htm#Video%20Memory)

## Success Criteria
- All video modes properly supported
- Pixel operations work correctly in all modes
- Double buffering provides smooth animation
- Performance suitable for real-time graphics
- Memory safety maintained throughout
- Example program demonstrates all features
- Comprehensive test coverage (>90%)
- DMA acceleration provides measurable performance improvement