# Task 11: Special Effects System

## Objective
Implement a comprehensive special effects system for the GBA that provides visual effects like blending, mosaic, windowing, and other hardware-accelerated graphics effects to enhance the visual appeal of games.

## Background
The GBA hardware includes several special effect capabilities that can dramatically enhance the visual presentation of games. These effects include alpha blending, mosaic effects, windowing, and various blending modes that operate at the hardware level for optimal performance.

### GBA Special Effects Hardware
- **Alpha Blending**: Transparency effects between layers
- **Brightness Control**: Fade to black/white effects  
- **Mosaic Effects**: Pixelated appearance for backgrounds and sprites
- **Window Effects**: Selective rendering areas
- **Color Special Effects**: Various blending modes

## Requirements

### Functional Requirements
1. Alpha blending between backgrounds and sprites
2. Brightness fade effects (fade in/out)
3. Mosaic effects for retro aesthetics
4. Window masking and effects
5. Color palette effects and transitions
6. Screen transitions and wipe effects
7. Particle effect support
8. Visual effect animation and sequencing

### Technical Requirements
- Hardware register manipulation for effects
- Integration with background and sprite systems
- Performance optimization for real-time effects
- Memory-efficient effect management

## API Design

### Constants
```go
package effects

const (
    // Blend modes
    BLEND_ALPHA     = 1
    BLEND_BRIGHTEN  = 2
    BLEND_DARKEN    = 3
    
    // Blend targets
    TARGET_BG0      = 0x01
    TARGET_BG1      = 0x02
    TARGET_BG2      = 0x04
    TARGET_BG3      = 0x08
    TARGET_OBJ      = 0x10
    TARGET_BACKDROP = 0x20
    
    // Window types
    WINDOW_0        = 0
    WINDOW_1        = 1
    WINDOW_OBJ      = 2
    WINDOW_OUT      = 3
    
    // Mosaic sizes
    MOSAIC_1x1      = 0
    MOSAIC_2x2      = 1
    MOSAIC_4x4      = 3
    MOSAIC_8x8      = 7
    MOSAIC_16x16    = 15
    
    // Effect timing
    EFFECT_INSTANT  = 0
    EFFECT_LINEAR   = 1
    EFFECT_EASE_IN  = 2
    EFFECT_EASE_OUT = 3
)
```

### Data Structures
```go
// BlendEffect represents alpha blending configuration
type BlendEffect struct {
    mode        int
    firstTarget uint8
    secondTarget uint8
    alphaA      uint8  // First target coefficient
    alphaB      uint8  // Second target coefficient
    brightness  uint8  // Brightness level
    active      bool
}

// MosaicEffect represents mosaic configuration
type MosaicEffect struct {
    bgHSize     uint8  // Background horizontal size
    bgVSize     uint8  // Background vertical size
    objHSize    uint8  // Object horizontal size
    objVSize    uint8  // Object vertical size
    active      bool
}

// WindowEffect represents window masking
type WindowEffect struct {
    id          int
    left        uint8
    right       uint8
    top         uint8
    bottom      uint8
    bgEnable    uint8  // Which backgrounds to show
    objEnable   bool   // Show objects
    effectEnable bool  // Enable special effects
    active      bool
}

// EffectManager coordinates all special effects
type EffectManager struct {
    blendEffect   BlendEffect
    mosaicEffect  MosaicEffect
    windows       [4]WindowEffect
    animations    []*EffectAnimation
    enabled       bool
}

// EffectAnimation represents a time-based effect change
type EffectAnimation struct {
    target      interface{}
    property    string
    startValue  float32
    endValue    float32
    duration    int
    elapsed     int
    timing      int
    onComplete  func()
    active      bool
}
```

### Core Functions
```go
// Effect manager
func NewEffectManager() *EffectManager
func (em *EffectManager) Update()
func (em *EffectManager) Enable(enable bool)
func (em *EffectManager) Reset()

// Blend effects
func (em *EffectManager) SetBlendMode(mode int)
func (em *EffectManager) SetBlendTargets(first, second uint8)
func (em *EffectManager) SetAlphaBlend(alphaA, alphaB uint8)
func (em *EffectManager) SetBrightness(level uint8)
func (em *EffectManager) FadeToBlack(duration int, callback func())
func (em *EffectManager) FadeToWhite(duration int, callback func())
func (em *EffectManager) FadeIn(duration int, callback func())

// Mosaic effects
func (em *EffectManager) SetMosaic(bgH, bgV, objH, objV uint8)
func (em *EffectManager) EnableMosaic(enable bool)
func (em *EffectManager) AnimateMosaic(targetH, targetV uint8, duration int)

// Window effects
func (em *EffectManager) SetWindow(id int, left, right, top, bottom uint8)
func (em *EffectManager) SetWindowContent(id int, bgMask uint8, showObj, showEffects bool)
func (em *EffectManager) EnableWindow(id int, enable bool)
func (em *EffectManager) AnimateWindow(id int, targetLeft, targetRight, targetTop, targetBottom uint8, duration int)

// Effect animations
func (em *EffectManager) AnimateProperty(target interface{}, property string, endValue float32, duration, timing int) *EffectAnimation
func (em *EffectManager) StopAnimation(anim *EffectAnimation)
func (em *EffectManager) StopAllAnimations()

// Screen transitions
func (em *EffectManager) WipeLeft(duration int, callback func())
func (em *EffectManager) WipeRight(duration int, callback func())
func (em *EffectManager) WipeUp(duration int, callback func())
func (em *EffectManager) WipeDown(duration int, callback func())
func (em *EffectManager) CircleWipe(centerX, centerY, duration int, callback func())

// Utility effects
func (em *EffectManager) Flash(color palette.Color, intensity uint8, duration int)
func (em *EffectManager) Shake(intensity, duration int)
func (em *EffectManager) Ripple(centerX, centerY, intensity, duration int)
```

## Implementation Details

### Step 1: Hardware Interface
Create `lib/effects/hardware.go`:
- Direct special effects register access
- Hardware capability detection
- Register state management

### Step 2: Blend Effects
Create `lib/effects/blending.go`:
- Alpha blending configuration
- Brightness control
- Fade effect implementations

### Step 3: Mosaic Effects
Create `lib/effects/mosaic.go`:
- Mosaic size control
- Animation support
- Background/sprite targeting

### Step 4: Window Effects  
Create `lib/effects/windows.go`:
- Window boundary setup
- Content masking control
- Window animation support

### Step 5: Effect Manager
Create `lib/effects/manager.go`:
- Central effect coordination
- Effect animation system
- State management

### Step 6: Screen Transitions
Create `lib/effects/transitions.go`:
- Various screen wipe effects
- Circle and pattern transitions
- Transition sequencing

## Testing Strategy

### Unit Tests
```go
func TestEffectManager(t *testing.T) {
    manager := NewEffectManager()
    
    // Test blend effect
    manager.SetBlendMode(BLEND_ALPHA)
    manager.SetBlendTargets(TARGET_BG0, TARGET_BG1)
    manager.SetAlphaBlend(8, 8)
    
    assert.Equal(t, BLEND_ALPHA, manager.blendEffect.mode)
    assert.Equal(t, uint8(TARGET_BG0), manager.blendEffect.firstTarget)
}

func TestMosaicEffect(t *testing.T) {
    manager := NewEffectManager()
    
    manager.SetMosaic(4, 4, 2, 2)
    manager.EnableMosaic(true)
    
    assert.Equal(t, uint8(4), manager.mosaicEffect.bgHSize)
    assert.True(t, manager.mosaicEffect.active)
}

func TestWindowEffect(t *testing.T) {
    manager := NewEffectManager()
    
    manager.SetWindow(WINDOW_0, 50, 200, 30, 130)
    manager.SetWindowContent(WINDOW_0, TARGET_BG0|TARGET_BG1, true, true)
    manager.EnableWindow(WINDOW_0, true)
    
    window := manager.windows[WINDOW_0]
    assert.Equal(t, uint8(50), window.left)
    assert.Equal(t, uint8(200), window.right)
    assert.True(t, window.active)
}
```

## Example Program
```go
package main

import (
    \"github.com/matheusmortatti/gba-go/lib/effects\"
    \"github.com/matheusmortatti/gba-go/lib/backgrounds\"
    \"github.com/matheusmortatti/gba-go/lib/sprites\"
    \"github.com/matheusmortatti/gba-go/lib/memory\"
    \"github.com/matheusmortatti/gba-go/lib/video\"
    \"github.com/matheusmortatti/gba-go/lib/input\"
)

func main() {
    // Initialize systems
    effectManager := effects.NewEffectManager()
    bgManager := backgrounds.NewBackgroundManager(memory.MODE_0)
    spriteManager := sprites.NewSpriteManager()
    
    // Set up backgrounds and sprites
    setupScene(bgManager, spriteManager)\n    \n    currentEffect := 0\n    effectNames := []string{\n        \"Alpha Blend\",\n        \"Fade Effects\", \n        \"Mosaic\",\n        \"Window\",\n        \"Transitions\",\n    }\n    \n    for {\n        video.VSync()\n        input.Poll()\n        \n        // Cycle through effects\n        if input.BtnClicked(input.KeyA) {\n            currentEffect = (currentEffect + 1) % len(effectNames)\n            switchEffect(effectManager, currentEffect)\n        }\n        \n        // Trigger effect actions\n        if input.BtnClicked(input.KeyB) {\n            triggerEffectAction(effectManager, currentEffect)\n        }\n        \n        // Update effect system\n        effectManager.Update()\n        bgManager.Update()\n        spriteManager.Update()\n    }\n}\n\nfunc setupScene(bgManager *backgrounds.BackgroundManager, spriteManager *sprites.SpriteManager) {\n    // Create backgrounds for demonstration\n    config1 := backgrounds.CreateTextConfig(0, 8, backgrounds.BG_SIZE_256x256, backgrounds.PRIORITY_1)\n    config2 := backgrounds.CreateTextConfig(1, 9, backgrounds.BG_SIZE_256x256, backgrounds.PRIORITY_0)\n    \n    bg1 := bgManager.CreateTextBackground(backgrounds.BG0, config1)\n    bg2 := bgManager.CreateTextBackground(backgrounds.BG1, config2)\n    \n    // Create some sprites\n    for i := 0; i < 10; i++ {\n        x := 50 + i*20\n        y := 80\n        sprite := spriteManager.CreateSprite(x, y, i%4, sprites.SIZE_16x16)\n        sprite.SetPriority(1)\n    }\n}\n\nfunc switchEffect(manager *effects.EffectManager, effectType int) {\n    // Reset all effects first\n    manager.Reset()\n    \n    switch effectType {\n    case 0: // Alpha Blend\n        setupAlphaBlend(manager)\n    case 1: // Fade Effects\n        setupFadeEffects(manager)\n    case 2: // Mosaic\n        setupMosaicEffect(manager)\n    case 3: // Window\n        setupWindowEffect(manager)\n    case 4: // Transitions\n        setupTransitionEffect(manager)\n    }\n}\n\nfunc setupAlphaBlend(manager *effects.EffectManager) {\n    manager.SetBlendMode(effects.BLEND_ALPHA)\n    manager.SetBlendTargets(effects.TARGET_BG0, effects.TARGET_BG1)\n    manager.SetAlphaBlend(8, 8) // 50% blend\n}\n\nfunc setupFadeEffects(manager *effects.EffectManager) {\n    manager.SetBlendMode(effects.BLEND_DARKEN)\n    manager.SetBlendTargets(effects.TARGET_BG0|effects.TARGET_BG1|effects.TARGET_OBJ, 0)\n    manager.SetBrightness(0) // Normal brightness\n}\n\nfunc setupMosaicEffect(manager *effects.EffectManager) {\n    manager.SetMosaic(1, 1, 1, 1) // Start with no mosaic\n    manager.EnableMosaic(true)\n}\n\nfunc setupWindowEffect(manager *effects.EffectManager) {\n    // Create a window in the center of the screen\n    manager.SetWindow(effects.WINDOW_0, 80, 160, 60, 100)\n    manager.SetWindowContent(effects.WINDOW_0, effects.TARGET_BG0, true, true)\n    manager.EnableWindow(effects.WINDOW_0, true)\n    \n    // Outside the window, only show BG1\n    manager.SetWindowContent(effects.WINDOW_OUT, effects.TARGET_BG1, false, false)\n}\n\nfunc setupTransitionEffect(manager *effects.EffectManager) {\n    // Prepare for transition effects\n    manager.SetBlendMode(effects.BLEND_DARKEN)\n    manager.SetBlendTargets(effects.TARGET_BG0|effects.TARGET_BG1|effects.TARGET_OBJ, 0)\n}\n\nfunc triggerEffectAction(manager *effects.EffectManager, effectType int) {\n    switch effectType {\n    case 0: // Alpha Blend - cycle alpha values\n        cycleAlphaBlend(manager)\n    case 1: // Fade Effects - trigger fade\n        triggerFade(manager)\n    case 2: // Mosaic - animate mosaic size\n        animateMosaic(manager)\n    case 3: // Window - animate window size\n        animateWindow(manager)\n    case 4: // Transitions - trigger screen wipe\n        triggerTransition(manager)\n    }\n}\n\nfunc cycleAlphaBlend(manager *effects.EffectManager) {\n    // Cycle through different alpha blend ratios\n    ratios := [][]uint8{\n        {16, 0},  // 100% first target\n        {12, 4},  // 75/25\n        {8, 8},   // 50/50\n        {4, 12},  // 25/75\n        {0, 16},  // 100% second target\n    }\n    \n    static currentRatio := 0\n    currentRatio = (currentRatio + 1) % len(ratios)\n    \n    ratio := ratios[currentRatio]\n    manager.SetAlphaBlend(ratio[0], ratio[1])\n}\n\nfunc triggerFade(manager *effects.EffectManager) {\n    // Fade to black and back\n    manager.FadeToBlack(60, func() {\n        // After 1 second of black, fade back in\n        manager.FadeIn(60, nil)\n    })\n}\n\nfunc animateMosaic(manager *effects.EffectManager) {\n    // Animate from no mosaic to heavy mosaic and back\n    manager.AnimateMosaic(8, 8, 60) // 1 second to full mosaic\n    \n    // Schedule return to normal\n    scheduler := timers.NewScheduler()\n    scheduler.ScheduleOnce(120, func() { // After 2 seconds\n        manager.AnimateMosaic(1, 1, 60) // 1 second back to normal\n    })\n}\n\nfunc animateWindow(manager *effects.EffectManager) {\n    // Animate window expanding and contracting\n    manager.AnimateWindow(effects.WINDOW_0, 20, 220, 20, 140, 60) // Expand\n    \n    scheduler := timers.NewScheduler()\n    scheduler.ScheduleOnce(120, func() {\n        manager.AnimateWindow(effects.WINDOW_0, 80, 160, 60, 100, 60) // Contract\n    })\n}\n\nfunc triggerTransition(manager *effects.EffectManager) {\n    // Cycle through different transition effects\n    transitions := []func(){\n        func() { manager.WipeLeft(60, nil) },\n        func() { manager.WipeRight(60, nil) },\n        func() { manager.WipeUp(60, nil) },\n        func() { manager.WipeDown(60, nil) },\n        func() { manager.CircleWipe(120, 80, 60, nil) },\n    }\n    \n    static currentTransition := 0\n    transitions[currentTransition]()\n    currentTransition = (currentTransition + 1) % len(transitions)\n}\n```\n\n## Advanced Features\n\n### Particle Effects\n```go\ntype ParticleSystem struct {\n    particles []Particle\n    manager   *effects.EffectManager\n}\n\ntype Particle struct {\n    x, y        float32\n    vx, vy      float32\n    life        int\n    maxLife     int\n    sprite      *sprites.Sprite\n}\n\nfunc (ps *ParticleSystem) Update() {\n    // Update particle physics and apply effects\n}\n```\n\n### Screen Shake\n```go\nfunc (em *EffectManager) Shake(intensity, duration int) {\n    // Implement screen shake using background scrolling\n}\n```\n\n## File Structure\n```\nlib/effects/\n├── hardware.go     // Hardware register interface\n├── blending.go     // Alpha blending and brightness\n├── mosaic.go       // Mosaic effects\n├── windows.go      // Window masking\n├── manager.go      // Effect coordination\n└── transitions.go  // Screen transition effects\n```\n\n## Integration Points\n- Use backgrounds from Task 07\n- Integrate with sprites from Task 04\n- Coordinate with timers from Task 09\n- Support game framework from Task 12\n\n## Resources\n- [GBATEK Special Effects](https://problemkaputt.de/gbatek.htm#lcdcolorspecialeffects)\n- [Tonc Special Effects](https://www.coranac.com/tonc/text/gfx.htm)\n- [GBA Graphics Effects](https://www.cs.rit.edu/~tjh8300/CowBite/CowBiteSpec.htm#Graphics%20Special%20Effects)\n\n## Success Criteria\n- All hardware effects work correctly\n- Smooth effect animations without frame drops\n- Proper integration with backgrounds and sprites\n- Visual effects enhance game presentation\n- Example program demonstrates all effect types\n- Performance suitable for real-time effects\n- Comprehensive test coverage (>90%)