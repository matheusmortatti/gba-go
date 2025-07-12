# Task 02: Palette Management System

## Objective
Implement a comprehensive palette management system for the GBA, supporting both background and sprite palettes, with utilities for loading, manipulating, and managing color palettes efficiently.

## Background
The GBA uses a palette-based color system with 15-bit colors (32768 possible colors). The system has separate palettes for backgrounds and sprites (objects), with support for both 16-color and 256-color modes.

### GBA Palette System
- **Palette RAM**: 1KB total (0x05000000 - 0x050003FF)
- **Background Palettes**: 0x05000000 - 0x050001FF (256 colors)
- **Sprite Palettes**: 0x05000200 - 0x050003FF (256 colors)
- **15-bit Color Format**: 0bGGGGGRRRRRBBBBB (5 bits each for RGB)
- **Palette Modes**: 16-color (16 palettes) or 256-color (1 palette)

### Color Layout
- Each color is 2 bytes (16-bit), but only 15 bits are used
- Background palettes: 16 sub-palettes of 16 colors each OR 1 palette of 256 colors
- Sprite palettes: 16 sub-palettes of 16 colors each OR 1 palette of 256 colors
- Color 0 in each 16-color palette is transparent

## Requirements

### Functional Requirements
1. Load and manage background and sprite palettes
2. Support both 16-color and 256-color palette modes
3. Provide color format conversion utilities (RGB24 to RGB15)
4. Handle transparency (color 0 in 16-color palettes)
5. Palette animation and cycling support
6. Memory-efficient palette operations

### Technical Requirements
- Direct hardware palette RAM access
- Type-safe palette operations
- Support for common image formats (later integration)
- Efficient bulk palette operations
- Bounds checking for palette indices

## API Design

### Constants
```go
package palette

const (
    // Palette RAM layout
    BG_PALETTE_BASE    = 0x05000000
    OBJ_PALETTE_BASE   = 0x05000200
    PALETTE_SIZE       = 0x400       // 1KB total
    BG_PALETTE_SIZE    = 0x200       // 512 bytes
    OBJ_PALETTE_SIZE   = 0x200       // 512 bytes

    // Palette configurations
    COLORS_PER_PALETTE_16  = 16
    COLORS_PER_PALETTE_256 = 256
    SUB_PALETTES_16        = 16
    MAX_BG_PALETTES        = 16
    MAX_OBJ_PALETTES       = 16

    // Special indices
    TRANSPARENT_COLOR_INDEX = 0

    // Common colors (RGB15 format)
    BLACK       = 0x0000
    WHITE       = 0x7FFF
    RED         = 0x001F
    GREEN       = 0x03E0
    BLUE        = 0x7C00
    YELLOW      = 0x03FF
    CYAN        = 0x7FE0
    MAGENTA     = 0x7C1F
    GRAY        = 0x39CE
    DARK_GRAY   = 0x1CE7
    LIGHT_GRAY  = 0x5AD6
)
```

### Data Structures
```go
// Color represents a 15-bit GBA color
type Color uint16

// Palette16 represents a 16-color palette
type Palette16 [COLORS_PER_PALETTE_16]Color

// Palette256 represents a 256-color palette
type Palette256 [COLORS_PER_PALETTE_256]Color

// PaletteManager handles background and sprite palettes
type PaletteManager struct {
    bgPalettes  [MAX_BG_PALETTES]*Palette16
    objPalettes [MAX_OBJ_PALETTES]*Palette16
    bg256       *Palette256
    obj256      *Palette256
}

// PaletteBank represents the hardware palette memory
type PaletteBank struct {
    bgBank  *[BG_PALETTE_SIZE/2]volatile.Register16
    objBank *[OBJ_PALETTE_SIZE/2]volatile.Register16
}
```

### Core Functions
```go
// Color operations
func RGB15(r, g, b uint8) Color
func RGB24ToRGB15(r, g, b uint8) Color
func (c Color) ToRGB24() (r, g, b uint8)
func (c Color) R() uint8
func (c Color) G() uint8
func (c Color) B() uint8

// Palette creation and management
func NewPaletteManager() *PaletteManager
func (pm *PaletteManager) LoadBGPalette16(index int, palette *Palette16) error
func (pm *PaletteManager) LoadBGPalette256(palette *Palette256) error
func (pm *PaletteManager) LoadOBJPalette16(index int, palette *Palette16) error
func (pm *PaletteManager) LoadOBJPalette256(palette *Palette256) error

// Hardware palette access
func GetPaletteBank() *PaletteBank
func (pb *PaletteBank) SetBGColor(paletteIndex, colorIndex int, color Color)
func (pb *PaletteBank) GetBGColor(paletteIndex, colorIndex int) Color
func (pb *PaletteBank) SetOBJColor(paletteIndex, colorIndex int, color Color)
func (pb *PaletteBank) GetOBJColor(paletteIndex, colorIndex int) Color

// Utility functions
func (p *Palette16) SetColor(index int, color Color) error
func (p *Palette16) GetColor(index int) Color
func (p *Palette256) SetColor(index int, color Color) error
func (p *Palette256) GetColor(index int) Color
func CreateGradient(start, end Color, steps int) []Color
func CreateGrayscalePalette() *Palette16

// Palette effects
func (pm *PaletteManager) FadeTo(targetPalette *Palette16, steps int, paletteIndex int)
func (pm *PaletteManager) RotatePalette(paletteIndex int, startColor, endColor int)
func BlendColors(color1, color2 Color, ratio float32) Color
```

## Implementation Details

### Step 1: Color System Foundation
Create `lib/palette/color.go` with:
- 15-bit color type and operations
- RGB conversion functions (24-bit to 15-bit)
- Color arithmetic and blending operations
- Common color constants

### Step 2: Palette Data Structures
Create `lib/palette/palette.go` with:
- Palette16 and Palette256 types
- Safe color access with bounds checking
- Palette creation and initialization utilities

### Step 3: Hardware Interface
Create `lib/palette/hardware.go` with:
- Direct palette RAM access using volatile registers
- PaletteBank structure for hardware interaction
- Efficient bulk palette loading operations

### Step 4: Palette Manager
Create `lib/palette/manager.go` with:
- High-level palette management
- Mode switching between 16-color and 256-color
- Palette slot allocation and tracking

### Step 5: Palette Effects
Create `lib/palette/effects.go` with:
- Fade effects and transitions
- Palette rotation and cycling
- Color interpolation and gradients

### Step 6: Utility Functions
Create `lib/palette/utils.go` with:
- Common palette generators (grayscale, rainbow, etc.)
- Palette validation and debugging utilities
- Export/import functions for palette data

## Testing Strategy

### Unit Tests
```go
func TestColorConversion(t *testing.T) {
    // Test RGB24 to RGB15 conversion
    color := RGB24ToRGB15(255, 128, 64)
    r, g, b := color.ToRGB24()
    
    // Allow for slight precision loss in conversion
    assert.InDelta(t, 255, r, 8) // 5-bit precision = ~8 unit tolerance
    assert.InDelta(t, 128, g, 8)
    assert.InDelta(t, 64, b, 8)
}

func TestPaletteOperations(t *testing.T) {
    palette := &Palette16{}
    
    // Test color setting and getting
    red := RGB15(31, 0, 0)
    err := palette.SetColor(1, red)
    assert.NoError(t, err)
    assert.Equal(t, red, palette.GetColor(1))
    
    // Test bounds checking
    err = palette.SetColor(16, red) // Out of bounds
    assert.Error(t, err)
}

func TestPaletteManager(t *testing.T) {
    manager := NewPaletteManager()
    palette := &Palette16{}
    palette.SetColor(1, RGB15(31, 0, 0))
    
    err := manager.LoadBGPalette16(0, palette)
    assert.NoError(t, err)
    
    // Verify hardware was updated
    bank := GetPaletteBank()
    color := bank.GetBGColor(0, 1)
    assert.Equal(t, RGB15(31, 0, 0), color)
}
```

### Integration Tests
- Test palette loading with actual sprite/background rendering
- Verify color accuracy on hardware/emulator
- Test palette effects during gameplay

## Example Program
```go
package main

import (
    "github.com/matheusmortatti/gba-go/lib/palette"
    "github.com/matheusmortatti/gba-go/lib/memory"
    "github.com/matheusmortatti/gba-go/lib/registers"
    "github.com/matheusmortatti/gba-go/lib/video"
)

func main() {
    // Initialize palette manager
    paletteManager := palette.NewPaletteManager()
    
    // Create a rainbow palette for backgrounds
    rainbowPalette := createRainbowPalette()
    paletteManager.LoadBGPalette16(0, rainbowPalette)
    
    // Create a grayscale palette for sprites
    grayPalette := palette.CreateGrayscalePalette()
    paletteManager.LoadOBJPalette16(0, grayPalette)
    
    // Set video mode 3 for demonstration
    registers.Lcd.DISPCNT.SetBits(memory.MODE_3)
    registers.Lcd.DISPCNT.SetBits(1 << 10) // Enable BG2
    
    vram := memory.GetVRAM()
    colorIndex := 0
    
    for {
        video.VSync()
        
        // Cycle through palette colors
        color := rainbowPalette.GetColor(colorIndex % 16)
        
        // Fill screen with current palette color
        for y := 0; y < memory.SCREEN_HEIGHT; y++ {
            for x := 0; x < memory.SCREEN_WIDTH; x++ {
                offset := uintptr((y*memory.SCREEN_WIDTH + x) * 2)
                vram.Write16(offset, uint16(color))
            }
        }
        
        colorIndex++
        
        // Wait a bit between color changes
        for i := 0; i < 30; i++ {
            video.VSync()
        }
    }
}

func createRainbowPalette() *palette.Palette16 {
    pal := &palette.Palette16{}
    
    // Create rainbow colors
    colors := []palette.Color{
        palette.BLACK,      // 0 - transparent
        palette.RED,        // 1
        palette.RGB15(31, 15, 0),  // Orange
        palette.YELLOW,     // 3
        palette.GREEN,      // 4
        palette.RGB15(0, 31, 15),  // Cyan-green
        palette.CYAN,       // 6
        palette.BLUE,       // 7
        palette.RGB15(15, 0, 31),  // Purple
        palette.MAGENTA,    // 9
        palette.WHITE,      // 10
        palette.LIGHT_GRAY, // 11
        palette.GRAY,       // 12
        palette.DARK_GRAY,  // 13
        palette.RGB15(15, 7, 3),   // Brown
        palette.RGB15(31, 20, 31), // Pink
    }
    
    for i, color := range colors {
        pal.SetColor(i, color)
    }
    
    return pal
}
```

## Advanced Features

### Palette Animation Example
```go
func animatePalette(manager *palette.PaletteManager) {
    // Create base palette
    basePalette := createRainbowPalette()
    
    // Animate palette by rotating colors
    for frame := 0; frame < 60; frame++ {
        animatedPalette := &palette.Palette16{}
        
        // Rotate colors based on frame
        for i := 1; i < 16; i++ {
            sourceIndex := ((i + frame) % 15) + 1 // Skip color 0 (transparent)
            color := basePalette.GetColor(sourceIndex)
            animatedPalette.SetColor(i, color)
        }
        
        manager.LoadBGPalette16(0, animatedPalette)
        video.VSync()
    }
}
```

## File Structure
```
lib/palette/
├── color.go         // Color type and conversion utilities
├── palette.go       // Palette data structures
├── hardware.go      // Hardware palette RAM interface
├── manager.go       // High-level palette management
├── effects.go       // Palette animation and effects
└── utils.go         // Utility functions and generators
```

## Integration with Existing Code
- Import memory constants from Task 01
- Update existing drawing code to use palette colors
- Ensure compatibility with register definitions
- Provide migration helpers for direct color usage

## Resources
- [GBATEK Palette RAM](https://problemkaputt.de/gbatek.htm#gbalcdcolorpalettes)
- [Tonc Palettes and Colors](https://www.coranac.com/tonc/text/bitmaps.htm)
- [GBA Color Format Specification](https://mgba.io/gbatek/#lcdcolorpalettes)

## Success Criteria
- All palette operations function correctly
- Color conversion maintains visual accuracy
- Hardware palette loading works without issues
- Palette effects render smoothly
- Example program demonstrates full functionality
- Comprehensive test coverage (>90%)
- Performance suitable for real-time palette animation