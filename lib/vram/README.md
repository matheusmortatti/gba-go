# VRAM Management System

A comprehensive Video RAM (VRAM) management system for Game Boy Advance development, providing efficient and safe access to graphics memory across all GBA video modes.

## Overview

This package provides high-level abstractions for GBA VRAM operations, including bitmap graphics, tile rendering, and double buffering support. It's designed for optimal performance while maintaining memory safety.

## Features

- **Multi-mode Support**: All GBA video modes (0-5) with mode-specific VRAM layouts
- **Double Buffering**: Smooth animation support for bitmap modes 4 & 5
- **DMA Acceleration**: Hardware-accelerated bulk operations for maximum performance
- **Memory Safety**: Comprehensive bounds checking and validation
- **Tile Management**: Complete tile and screen data handling for tile-based modes
- **Debug Utilities**: Memory analysis and performance profiling tools

## File Structure

```
lib/vram/
├── README.md           # This documentation
├── manager.go          # Core VRAM manager and mode handling
├── bitmap.go           # Bitmap buffer operations and pixel manipulation
├── fast.go             # DMA-optimized bulk operations
├── tiles.go            # Tile data and screen map management
├── buffering.go        # Double and triple buffering support
├── utils.go            # Utility functions and debugging tools
└── vram_test.go        # Comprehensive unit tests
```

## Core Components

### VRAMManager (`manager.go`)

The central manager for VRAM operations across different video modes.

```go
// Create a VRAM manager for Mode 4 (8-bit bitmap with double buffering)
vm := vram.NewVRAMManager(memory.MODE_4)

// Get current framebuffer
buffer := vm.GetCurrentBuffer()

// Check capabilities
if vm.SupportsDoubleBuffering() {
    backBuffer := vm.GetBackBuffer()
    // Draw to back buffer...
    vm.SwapBuffers()
}
```

**Key Features:**
- Mode-specific configuration (0-5)
- Buffer allocation and tracking
- Double buffering support detection
- Memory layout management

### BitmapBuffer (`bitmap.go`)

High-level bitmap operations with bounds checking and safety.

```go
// Basic pixel operations
buffer.PlotPixel(x, y, color)          // Safe with bounds checking
buffer.PlotPixelFast(x, y, color)      // Unsafe but faster

// Drawing primitives
buffer.Clear(color)
buffer.FillRect(x, y, width, height, color)
buffer.DrawLine(x1, y1, x2, y2, color)
buffer.DrawCircle(centerX, centerY, radius, color)
buffer.FillCircle(centerX, centerY, radius, color)

// Buffer operations
buffer.CopyFrom(srcBuffer, srcX, srcY, dstX, dstY, width, height)
```

**Key Features:**
- 8-bit and 16-bit pixel support
- Bounds checking on all operations
- Optimized drawing primitives
- Buffer-to-buffer copying

### Fast Operations (`fast.go`)

DMA-accelerated operations for maximum performance.

```go
// Fast bulk operations using DMA
buffer.FastClear(color)
buffer.FastFill(x, y, width, height, color)
buffer.FastCopy(srcBuffer)

// DMA status checking
if vram.GetDMAStatus() {
    vram.WaitForDMA()
}
```

**Key Features:**
- DMA3 channel utilization
- Automatic fallback to CPU for small operations
- 32-bit aligned memory access
- Hardware acceleration

### Tile Management (`tiles.go`)

Complete tile data and screen map management for tile-based modes.

```go
// Character data management
tileData := vram.GetCharacterData(charBlock, bpp)
tileData.LoadTile(tileIndex, data)
tileData.ClearTile(tileIndex)

// Screen map management
screenData := vram.GetScreenData(screenBlock, width, height)
screenData.SetTile(x, y, tileIndex, attributes)
screenData.FillScreen(tileIndex, attributes)

// Tile attributes
attrs := vram.CombineTileAttributes(
    vram.SetTilePalette(3),
    vram.SetTileFlip(true, false),
)
```

**Key Features:**
- Character block management (0-3)
- Screen block management (0-31)
- 4bpp and 8bpp tile support
- Tile attribute handling

### Double Buffering (`buffering.go`)

Advanced buffering systems for smooth animation.

```go
// Double buffering
db := vram.NewDoubleBuffer(memory.MODE_4)
if db != nil {
    backBuffer := db.GetBackBuffer()
    // Draw to back buffer...
    db.Present() // Swap buffers with VSync
}

// Triple buffering (advanced)
tb := vram.NewTripleBuffer(memory.MODE_4)
if tb != nil {
    tb.RequestSwap()
    tb.Update() // Handle pending swaps
}
```

**Key Features:**
- Automatic page flipping
- VSync synchronization
- Hardware register updates
- Triple buffering support

### Utilities (`utils.go`)

Debugging and analysis tools.

```go
// Memory analysis
info := vram.GetVRAMDebugInfo(vm)
usage := vram.GetMemoryUsage(memory.MODE_4)

// Address analysis
addrInfo := vram.AnalyzeVRAMAddress(addr, mode)

// Test patterns
vram.FillPattern(buffer, vram.PATTERN_CHECKERBOARD, color1, color2)
vram.FillPattern(buffer, vram.PATTERN_RAINBOW, 0, 0)

// Memory utilities
data := vram.DumpVRAMRegion(offset, size)
vram.LoadVRAMRegion(offset, data)
```

**Key Features:**
- Memory usage analysis
- Test pattern generation
- Address validation
- Performance profiling

## Usage Examples

### Basic Bitmap Graphics (Mode 3)

```go
func main() {
    // Initialize for Mode 3 (16-bit bitmap)
    vm := vram.NewVRAMManager(memory.MODE_3)
    registers.Lcd.DISPCNT.Set(memory.MODE_3 | (1 << 10))
    
    buffer := vm.GetCurrentBuffer()
    
    // Clear screen to black
    buffer.Clear(0x0000)
    
    // Draw some graphics
    buffer.FillRect(50, 50, 100, 60, 0x7FFF)        // White rectangle
    buffer.DrawCircle(120, 80, 30, 0x001F)          // Red circle
    buffer.DrawLine(0, 0, 239, 159, 0x03E0)         // Green diagonal line
    
    for {
        video.VSync()
        // Static display - no animation needed
    }
}
```

### Double Buffered Animation (Mode 4)

```go
func main() {
    // Initialize double buffering
    db := vram.NewDoubleBuffer(memory.MODE_4)
    registers.Lcd.DISPCNT.Set(memory.MODE_4 | (1 << 10))
    
    // Set up palette
    paletteManager := palette.NewPaletteManager()
    // ... palette setup
    
    frame := 0
    for {
        // Get back buffer for drawing
        backBuffer := db.GetBackBuffer()
        
        // Clear and draw animated content
        backBuffer.FastClear(0)
        drawAnimatedScene(backBuffer, frame)
        
        // Present the frame
        db.Present()
        frame++
    }
}
```

### Tile-Based Graphics (Mode 0)

```go
func main() {
    // Set up tile mode
    registers.Lcd.DISPCNT.Set(memory.MODE_0 | (1 << 8)) // Mode 0 + BG0
    
    // Load character data
    tileData := vram.GetCharacterData(0, 4) // Block 0, 4bpp
    for i, tile := range tileGraphics {
        tileData.LoadTile(i, tile)
    }
    
    // Set up screen map
    screenData := vram.GetScreenData(8, 32, 32) // Block 8, 32x32 tiles
    
    // Draw tilemap
    for y := 0; y < 20; y++ {
        for x := 0; x < 30; x++ {
            tileIndex := levelData[y][x]
            screenData.SetTile(x, y, tileIndex, 0)
        }
    }
    
    for {
        video.VSync()
        input.Poll()
        // Handle input and update tilemap as needed
    }
}
```

## Performance Considerations

### DMA Usage
- DMA operations are automatically used for large fills/copies (>64 pixels)
- Small operations use CPU for better efficiency
- Always check `GetDMAStatus()` before starting new DMA operations

### Memory Alignment
- All memory operations are optimized for 16-bit and 32-bit alignment
- Use `AlignAddress()` utility for custom memory operations

### Bounds Checking
- Use `PlotPixelFast()` for performance-critical inner loops
- Regular `PlotPixel()` includes bounds checking for safety
- `InBounds()` method available for manual checking

## Testing

The package includes comprehensive unit tests covering:

```bash
# Run all VRAM tests
tinygo test -target gameboy-advance github.com/matheusmortatti/gba-go/lib/vram

# Run benchmarks
tinygo test -target gameboy-advance -bench=. github.com/matheusmortatti/gba-go/lib/vram
```

Test coverage includes:
- All video modes
- Pixel operations and bounds checking
- Double buffering functionality
- Tile data management
- Memory utilities
- Performance benchmarks

## Constants Reference

### Video Modes
- `memory.MODE_0` - Tile mode, 4 backgrounds
- `memory.MODE_1` - Mixed tile/affine mode
- `memory.MODE_2` - Affine tile mode
- `memory.MODE_3` - 16-bit bitmap (240x160)
- `memory.MODE_4` - 8-bit bitmap with palette (240x160)
- `memory.MODE_5` - 16-bit small bitmap (160x128)

### VRAM Layout
- `VRAM_BASE` - 0x06000000
- `VRAM_SIZE` - 96KB
- `FRAME0_ADDR` - Frame 0 address
- `FRAME1_ADDR` - Frame 1 address (modes 4&5)
- `CHAR_BASE_SIZE` - 16KB per character block
- `SCREEN_BASE_SIZE` - 2KB per screen block

### Pattern Types
- `PATTERN_SOLID` - Solid color fill
- `PATTERN_CHECKERBOARD` - Checkerboard pattern
- `PATTERN_GRADIENT_H/V` - Horizontal/vertical gradients
- `PATTERN_STRIPES_H/V` - Horizontal/vertical stripes
- `PATTERN_RAINBOW` - Rainbow pattern
- `PATTERN_NOISE` - Pseudo-random noise

## Integration

This VRAM system integrates with other GBA library components:

- **Memory System**: Uses `lib/memory` constants and utilities
- **Palette System**: Works with `lib/palette` for color management
- **Registers**: Uses `lib/registers` for hardware control
- **Video Sync**: Integrates with `lib/video` for VSync
- **Input**: Compatible with `lib/input` for interactive applications

## Error Handling

All functions return appropriate error values:
- Bounds checking errors for out-of-range operations
- Validation errors for invalid parameters
- Memory access errors for invalid addresses

Always check error returns in production code:

```go
if err := buffer.PlotPixel(x, y, color); err != nil {
    // Handle error appropriately
}
```

## Best Practices

1. **Always initialize display properly**:
   ```go
   registers.Lcd.DISPCNT.Set(mode | (1 << 10))
   ```

2. **Use appropriate buffer for your needs**:
   - Mode 3: Single buffer, immediate display
   - Mode 4/5: Double buffering for smooth animation

3. **Optimize performance-critical sections**:
   - Use `PlotPixelFast()` in inner loops
   - Use DMA operations for large fills
   - Check bounds manually when needed

4. **Test with real hardware**:
   ```bash
   mgba your-program.gba
   ```

For more examples, see `/examples/vram-demo.go` which demonstrates all major features of the VRAM system.