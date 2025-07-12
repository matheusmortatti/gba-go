# Task 04: OAM Sprite System

## Objective
Implement a complete Object Attribute Memory (OAM) sprite system for the GBA, providing high-level APIs for sprite creation, management, positioning, animation, and hardware acceleration with full support for all sprite features.

## Background
The GBA sprite system (also called "objects") uses Object Attribute Memory (OAM) to store sprite attributes. The hardware can display up to 128 sprites simultaneously with various sizes, transformations, and effects.

### OAM Structure
- **OAM Location**: 0x07000000 - 0x070003FF (1KB)
- **128 Objects**: Each object uses 8 bytes (6 bytes attributes + 2 bytes padding)
- **OAM Size**: 128 objects × 8 bytes = 1024 bytes
- **Sprite Tiles**: Stored in VRAM character blocks 4-5

### Sprite Attributes (per object)
**Attribute 0 (Y position, object mode, effects)**
- Bits 0-7: Y coordinate
- Bits 8-9: Object mode (Normal, Semi-transparent, Object window, Prohibited)
- Bits 10-11: GFX mode (Normal, Affine, Hide, Double-size affine)
- Bits 12: Mosaic enable
- Bit 13: Color mode (16-color/256-color)
- Bits 14-15: Object shape (Square, Horizontal, Vertical)

**Attribute 1 (X position, size, flip)**
- Bits 0-8: X coordinate
- Bits 9-13: Affine parameter selection (for affine sprites)
- Bit 12: Horizontal flip (non-affine)
- Bit 13: Vertical flip (non-affine)
- Bits 14-15: Object size

**Attribute 2 (Tile index, priority, palette)**
- Bits 0-9: Character name (tile index)
- Bits 10-11: Priority (0=highest, 3=lowest)
- Bits 12-15: Palette number (16-color mode only)

## Requirements

### Functional Requirements
1. Complete OAM management for all 128 hardware sprites
2. Support all sprite sizes (8x8 to 64x64)
3. Sprite positioning, scaling, rotation, and flipping
4. Priority management and depth sorting
5. Palette assignment for 16-color and 256-color modes
6. Sprite visibility and culling
7. Hardware affine transformation support
8. Collision detection utilities

### Technical Requirements
- Direct OAM hardware access using volatile registers
- Efficient sprite updating and batch operations
- Memory management for sprite tile data
- Integration with VRAM and palette systems
- Performance optimization for 60fps operation

## API Design

### Constants
```go
package sprites

const (
    // OAM constants
    OAM_BASE            = 0x07000000
    OAM_SIZE            = 0x400      // 1KB
    MAX_SPRITES         = 128
    OAM_ENTRY_SIZE      = 8          // 6 bytes attributes + 2 padding
    
    // Sprite tile data in VRAM
    SPRITE_TILES_BASE   = 0x06010000 // Character blocks 4-5
    SPRITE_TILES_SIZE   = 0x8000     // 32KB
    
    // Object modes
    OBJ_MODE_NORMAL     = 0
    OBJ_MODE_BLEND      = 1
    OBJ_MODE_WINDOW     = 2
    OBJ_MODE_PROHIBITED = 3
    
    // GFX modes
    GFX_MODE_NORMAL     = 0
    GFX_MODE_AFFINE     = 1
    GFX_MODE_HIDE       = 2
    GFX_MODE_DOUBLE     = 3
    
    // Color modes
    COLOR_16            = 0  // 16 colors, 16 palettes
    COLOR_256           = 1  // 256 colors, 1 palette
    
    // Object shapes
    SHAPE_SQUARE        = 0
    SHAPE_HORIZONTAL    = 1
    SHAPE_VERTICAL      = 2
    
    // Object sizes (combined with shape)
    SIZE_8x8    = 0  // Square
    SIZE_16x16  = 1
    SIZE_32x32  = 2
    SIZE_64x64  = 3
    SIZE_16x8   = 0  // Horizontal
    SIZE_32x8   = 1
    SIZE_32x16  = 2
    SIZE_64x32  = 3
    SIZE_8x16   = 0  // Vertical
    SIZE_8x32   = 1
    SIZE_16x32  = 2
    SIZE_32x64  = 3
)
```

### Data Structures
```go
// SpriteAttributes represents the three attribute words for a sprite
type SpriteAttributes struct {
    attr0 volatile.Register16
    attr1 volatile.Register16
    attr2 volatile.Register16
    dummy volatile.Register16  // Padding
}

// Sprite represents a high-level sprite object
type Sprite struct {
    id          int
    x, y        int
    tileIndex   int
    paletteId   int
    priority    int
    width       int
    height      int
    visible     bool
    hFlip       bool
    vFlip       bool
    colorMode   int
    objMode     int
    gfxMode     int
    
    // Affine transformation (if enabled)
    affineId    int
    scaleX      float32
    scaleY      float32
    rotation    float32
}

// AffineMatrix represents transformation parameters for affine sprites
type AffineMatrix struct {
    pa int16  // dx
    pb int16  // dmx  
    pc int16  // dy
    pd int16  // dmy
}

// SpriteManager manages all sprites and OAM operations
type SpriteManager struct {
    oam          *[MAX_SPRITES]SpriteAttributes
    sprites      [MAX_SPRITES]*Sprite
    affineParams [32]AffineMatrix  // 32 affine parameter sets
    nextFreeId   int
    spriteCount  int
}
```

### Core Functions
```go
// Sprite manager
func NewSpriteManager() *SpriteManager
func (sm *SpriteManager) Update() // Apply all sprite changes to OAM
func (sm *SpriteManager) Clear()  // Hide all sprites

// Sprite creation and management
func (sm *SpriteManager) CreateSprite(x, y, tileIndex int, size SpriteSize) *Sprite
func (sm *SpriteManager) DestroySprite(sprite *Sprite)
func (sm *SpriteManager) GetSprite(id int) *Sprite
func (sm *SpriteManager) GetFreeSprite() *Sprite

// Sprite positioning and visibility
func (s *Sprite) SetPosition(x, y int)
func (s *Sprite) GetPosition() (int, int)
func (s *Sprite) SetVisible(visible bool)
func (s *Sprite) IsVisible() bool
func (s *Sprite) Move(dx, dy int)

// Sprite appearance
func (s *Sprite) SetTile(tileIndex int)
func (s *Sprite) SetPalette(paletteId int)
func (s *Sprite) SetPriority(priority int)
func (s *Sprite) SetFlip(hFlip, vFlip bool)
func (s *Sprite) SetColorMode(mode int)

// Affine transformations
func (s *Sprite) SetAffine(enable bool) error
func (s *Sprite) SetScale(scaleX, scaleY float32) error
func (s *Sprite) SetRotation(angle float32) error
func (s *Sprite) SetTransform(scaleX, scaleY, rotation float32) error

// Collision detection
func (s *Sprite) GetBounds() (x, y, width, height int)
func (s *Sprite) CollidesWith(other *Sprite) bool
func (s *Sprite) CollidesWithPoint(x, y int) bool
func (s *Sprite) CollidesWithRect(x, y, width, height int) bool

// Utility functions
func (sm *SpriteManager) SortByPriority()
func (sm *SpriteManager) CountVisibleSprites() int
func (sm *SpriteManager) FindSpritesAt(x, y int) []*Sprite
func GetSpriteSize(shape, size int) (width, height int)
```

## Implementation Details

### Step 1: OAM Hardware Interface
Create `lib/sprites/oam.go`:
- Direct OAM memory access using volatile registers
- SpriteAttributes struct mapping to hardware layout
- Low-level attribute manipulation functions

### Step 2: Sprite Management
Create `lib/sprites/sprite.go`:
- High-level Sprite struct with game-friendly API
- Sprite lifecycle management (create/destroy)
- Position, visibility, and appearance control

### Step 3: Sprite Manager
Create `lib/sprites/manager.go`:
- SpriteManager for coordinating all sprites
- Batch OAM updates for performance
- Sprite ID allocation and tracking

### Step 4: Affine Transformations
Create `lib/sprites/affine.go`:
- Affine parameter calculation and management
- Scale, rotation, and transformation utilities
- Hardware affine parameter allocation

### Step 5: Collision Detection
Create `lib/sprites/collision.go`:
- Bounding box collision detection
- Point-in-sprite testing
- Spatial partitioning for performance

### Step 6: Utilities and Helpers
Create `lib/sprites/utils.go`:
- Sprite size calculation utilities
- Priority sorting algorithms
- Debug and visualization helpers

## Testing Strategy

### Unit Tests
```go
func TestSpriteCreation(t *testing.T) {
    manager := NewSpriteManager()
    
    sprite := manager.CreateSprite(100, 50, 0, SIZE_16x16)
    assert.NotNil(t, sprite)
    assert.Equal(t, 100, sprite.x)
    assert.Equal(t, 50, sprite.y)
    assert.Equal(t, 0, sprite.tileIndex)
    assert.True(t, sprite.visible)
}

func TestSpritePositioning(t *testing.T) {
    manager := NewSpriteManager()
    sprite := manager.CreateSprite(0, 0, 0, SIZE_16x16)
    
    sprite.SetPosition(50, 75)
    x, y := sprite.GetPosition()
    assert.Equal(t, 50, x)
    assert.Equal(t, 75, y)
    
    sprite.Move(10, -5)
    x, y = sprite.GetPosition()
    assert.Equal(t, 60, x)
    assert.Equal(t, 70, y)
}

func TestAffineTransformation(t *testing.T) {
    manager := NewSpriteManager()
    sprite := manager.CreateSprite(100, 100, 0, SIZE_32x32)
    
    err := sprite.SetAffine(true)
    assert.NoError(t, err)
    
    err = sprite.SetScale(1.5, 0.8)
    assert.NoError(t, err)
    
    err = sprite.SetRotation(45.0)
    assert.NoError(t, err)
}

func TestCollisionDetection(t *testing.T) {
    manager := NewSpriteManager()
    sprite1 := manager.CreateSprite(50, 50, 0, SIZE_16x16)
    sprite2 := manager.CreateSprite(60, 60, 1, SIZE_16x16)
    sprite3 := manager.CreateSprite(100, 100, 2, SIZE_16x16)
    
    // Overlapping sprites should collide
    assert.True(t, sprite1.CollidesWith(sprite2))
    assert.True(t, sprite2.CollidesWith(sprite1))
    
    // Non-overlapping sprites should not collide
    assert.False(t, sprite1.CollidesWith(sprite3))
    assert.False(t, sprite3.CollidesWith(sprite1))
}
```

### Integration Tests
- Test sprite rendering with actual graphics
- Verify OAM updates are applied correctly
- Test performance with maximum sprite count

## Example Program
```go
package main

import (
    "github.com/matheusmortatti/gba-go/lib/sprites"
    "github.com/matheusmortatti/gba-go/lib/memory"
    "github.com/matheusmortatti/gba-go/lib/palette"
    "github.com/matheusmortatti/gba-go/lib/registers"
    "github.com/matheusmortatti/gba-go/lib/video"
    "github.com/matheusmortatti/gba-go/lib/input"
    "math"
)

func main() {
    // Initialize sprite system
    spriteManager := sprites.NewSpriteManager()
    
    // Set video mode 0 (tile mode with sprites)
    registers.Lcd.DISPCNT.SetBits(memory.MODE_0)
    registers.Lcd.DISPCNT.SetBits(1 << 12) // Enable sprites
    
    // Set up sprite palette
    setupSpritePalette()
    
    // Load sprite tiles into VRAM
    loadSpriteTiles()
    
    // Create player sprite
    player := spriteManager.CreateSprite(120, 80, 0, sprites.SIZE_16x16)
    player.SetPalette(0)
    player.SetPriority(0) // Highest priority
    
    // Create enemy sprites
    enemies := make([]*sprites.Sprite, 5)
    for i := 0; i < 5; i++ {
        x := 20 + i*40
        y := 20
        enemies[i] = spriteManager.CreateSprite(x, y, 1, sprites.SIZE_16x16)
        enemies[i].SetPalette(1)
        enemies[i].SetPriority(1)
    }
    
    // Create rotating sprite with affine transformation
    rotatingSprite := spriteManager.CreateSprite(200, 120, 2, sprites.SIZE_32x32)
    rotatingSprite.SetAffine(true)
    rotatingSprite.SetPalette(2)
    
    frame := 0
    
    for {
        video.VSync()
        input.Poll()
        
        // Handle player input
        handlePlayerInput(player)
        
        // Update enemy movement
        updateEnemies(enemies, frame)
        
        // Update rotating sprite
        angle := float32(frame) * 2.0
        rotatingSprite.SetRotation(angle)
        scale := 1.0 + 0.3*float32(math.Sin(float64(frame)*0.1))
        rotatingSprite.SetScale(scale, scale)
        
        // Check collisions
        checkCollisions(player, enemies)
        
        // Apply all sprite updates to hardware
        spriteManager.Update()
        
        frame++
    }
}

func handlePlayerInput(player *sprites.Sprite) {
    x, y := player.GetPosition()
    speed := 2
    
    if input.BtnDown(input.KeyLeft) && x > 0 {
        player.SetPosition(x-speed, y)
    }
    if input.BtnDown(input.KeyRight) && x < memory.SCREEN_WIDTH-16 {
        player.SetPosition(x+speed, y)
    }
    if input.BtnDown(input.KeyUp) && y > 0 {
        player.SetPosition(x, y-speed)
    }
    if input.BtnDown(input.KeyDown) && y < memory.SCREEN_HEIGHT-16 {
        player.SetPosition(x, y+speed)
    }
    
    // Flip sprite based on movement direction
    if input.BtnDown(input.KeyLeft) {
        player.SetFlip(true, false)
    } else if input.BtnDown(input.KeyRight) {
        player.SetFlip(false, false)
    }
}

func updateEnemies(enemies []*sprites.Sprite, frame int) {
    for i, enemy := range enemies {
        if enemy == nil {
            continue
        }
        
        // Simple sine wave movement
        x := 20 + i*40 + int(20*math.Sin(float64(frame+i*20)*0.1))
        y := 20 + int(10*math.Cos(float64(frame+i*30)*0.15))
        
        enemy.SetPosition(x, y)
        
        // Animate sprite by changing tiles
        if frame%30 == 0 {
            tileIndex := 1 + (frame/30+i)%4 // Cycle through 4 animation frames
            enemy.SetTile(tileIndex)
        }
    }
}

func checkCollisions(player *sprites.Sprite, enemies []*sprites.Sprite) {
    for i, enemy := range enemies {
        if enemy == nil || !enemy.IsVisible() {
            continue
        }
        
        if player.CollidesWith(enemy) {
            // Hide enemy on collision
            enemy.SetVisible(false)
            enemies[i] = nil
            
            // Flash player sprite
            player.SetPalette((player.paletteId + 1) % 4)
        }
    }
}

func setupSpritePalette() {
    paletteManager := palette.NewPaletteManager()
    
    // Player palette (palette 0)
    playerPalette := &palette.Palette16{}
    playerPalette.SetColor(0, palette.TRANSPARENT) // Transparent
    playerPalette.SetColor(1, palette.BLUE)        // Main color
    playerPalette.SetColor(2, palette.WHITE)       // Highlight
    playerPalette.SetColor(3, palette.BLACK)       // Outline
    paletteManager.LoadOBJPalette16(0, playerPalette)
    
    // Enemy palette (palette 1)
    enemyPalette := &palette.Palette16{}
    enemyPalette.SetColor(0, palette.TRANSPARENT)
    enemyPalette.SetColor(1, palette.RED)
    enemyPalette.SetColor(2, palette.YELLOW)
    enemyPalette.SetColor(3, palette.BLACK)
    paletteManager.LoadOBJPalette16(1, enemyPalette)
    
    // Rotating sprite palette (palette 2)
    rotatingPalette := &palette.Palette16{}
    rotatingPalette.SetColor(0, palette.TRANSPARENT)
    rotatingPalette.SetColor(1, palette.GREEN)
    rotatingPalette.SetColor(2, palette.CYAN)
    rotatingPalette.SetColor(3, palette.MAGENTA)
    paletteManager.LoadOBJPalette16(2, rotatingPalette)
}

func loadSpriteTiles() {
    // In a real implementation, this would load actual sprite graphics
    // For now, we'll assume tile data is already in VRAM
    // This would typically use the tile system from Task 06
}
```

## Advanced Features

### Sprite Animation System
```go
type SpriteAnimation struct {
    frames    []int  // Tile indices
    delays    []int  // Frame delays
    loop      bool
    current   int
    timer     int
}

func (s *Sprite) PlayAnimation(anim *SpriteAnimation) {
    // Animation playback system
}
```

### Hardware Optimization
```go
func (sm *SpriteManager) BatchUpdate() {
    // Batch all OAM updates for efficiency
    // Sort sprites by memory access pattern
    // Use DMA for bulk OAM updates
}
```

## File Structure
```
lib/sprites/
├── oam.go           // Low-level OAM hardware interface
├── sprite.go        // High-level sprite object
├── manager.go       // Sprite manager and coordination
├── affine.go        // Affine transformation system
├── collision.go     // Collision detection utilities
└── utils.go         // Helper functions and debugging
```

## Integration Points
- Use memory constants from Task 01
- Integrate with palette system from Task 02
- Coordinate with VRAM management from Task 03
- Prepare for tile system integration (Task 06)

## Resources
- [GBATEK OAM and Sprites](https://problemkaputt.de/gbatek.htm#lcdobjoverview)
- [Tonc Sprites and Backgrounds](https://www.coranac.com/tonc/text/objbg.htm)
- [GBA Sprite Programming Guide](https://www.cs.rit.edu/~tjh8300/CowBite/CowBiteSpec.htm#OAM%20(Sprite%20Attribute%20Table))

## Success Criteria
- All 128 hardware sprites can be controlled independently
- Affine transformations work correctly (scale, rotate)
- Collision detection is accurate and performant
- Sprite rendering displays correctly in all modes
- Performance maintains 60fps with full sprite usage
- Example program demonstrates all sprite features
- Memory management prevents leaks or corruption
- Comprehensive test coverage (>90%)