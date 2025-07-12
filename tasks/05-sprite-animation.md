# Task 05: Sprite Animation System

## Objective
Create a comprehensive sprite animation system that extends the OAM sprite system with time-based frame animation, tweening, state machines, and performance-optimized animation playback for smooth character and object animation.

## Background
Game sprites require complex animation systems to bring characters and objects to life. The GBA's limited resources require efficient animation management that minimizes memory usage and CPU overhead while providing flexible animation capabilities.

### Animation Requirements
- **Frame-based Animation**: Sequence of sprite tiles played over time
- **Timing Control**: Variable frame durations and playback speeds
- **Animation States**: Multiple animations per sprite (idle, walk, jump, etc.)
- **Looping Control**: One-shot, looping, and ping-pong animations
- **Tweening**: Smooth interpolation of position, scale, rotation
- **Event System**: Callbacks for animation events and completion
- **Memory Efficiency**: Shared animation data between sprite instances

## Requirements

### Functional Requirements
1. Frame-based sprite animation with variable timing
2. Animation state management (idle, walk, attack, etc.)
3. Smooth tweening for position, scale, rotation, and color
4. Animation events and completion callbacks
5. Performance optimization for many animated sprites
6. Animation blending and transitions
7. Resource management for animation data
8. Integration with sprite system from Task 04

### Technical Requirements
- Minimal memory overhead per animated sprite
- 60fps animation playback capability
- Shared animation definitions between sprites
- Efficient update loop for many sprites
- Integration with timer system for precise timing

## API Design

### Constants
```go
package animation

const (
    // Animation playback modes
    ANIM_ONCE       = 0  // Play once and stop
    ANIM_LOOP       = 1  // Loop continuously
    ANIM_PINGPONG   = 2  // Play forward then backward
    ANIM_REVERSE    = 3  // Play in reverse once
    
    // Animation states
    ANIM_STOPPED    = 0
    ANIM_PLAYING    = 1
    ANIM_PAUSED     = 2
    ANIM_FINISHED   = 3
    
    // Easing types for tweening
    EASE_LINEAR     = 0
    EASE_IN_QUAD    = 1
    EASE_OUT_QUAD   = 2
    EASE_IN_OUT_QUAD = 3
    EASE_IN_CUBIC   = 4
    EASE_OUT_CUBIC  = 5
    EASE_BOUNCE     = 6
    
    // Maximum animation events per animation
    MAX_ANIM_EVENTS = 8
)
```

### Data Structures
```go
// AnimationFrame represents a single frame of animation
type AnimationFrame struct {
    TileIndex   int     // Sprite tile to display
    Duration    int     // Frame duration in game ticks (1/60 second)
    OffsetX     int8    // X offset from sprite position
    OffsetY     int8    // Y offset from sprite position
    FlipH       bool    // Horizontal flip
    FlipV       bool    // Vertical flip
}

// AnimationEvent represents an event during animation playback
type AnimationEvent struct {
    Frame    int                    // Frame number to trigger on
    Callback func(sprite *Sprite)  // Function to call
}

// AnimationDef defines a complete animation sequence
type AnimationDef struct {
    Name        string
    Frames      []AnimationFrame
    Events      []AnimationEvent
    Mode        int     // Playback mode (once, loop, etc.)
    Speed       float32 // Playback speed multiplier
    Priority    int     // Animation priority for blending
}

// SpriteAnimator manages animation state for a single sprite
type SpriteAnimator struct {
    sprite          *sprites.Sprite
    currentAnim     *AnimationDef
    nextAnim        *AnimationDef
    currentFrame    int
    frameTimer      int
    state           int
    speed           float32
    blendTimer      int
    blendDuration   int
    eventFlags      [MAX_ANIM_EVENTS]bool
}

// Tween represents a smooth interpolation animation
type Tween struct {
    target      interface{}  // Target object (sprite, etc.)
    property    string       // Property to animate ("x", "y", "scaleX", etc.)
    startValue  float32
    endValue    float32
    duration    int
    elapsed     int
    easeType    int
    onComplete  func()
    active      bool
}

// AnimationManager coordinates all sprite animations
type AnimationManager struct {
    animators   map[*sprites.Sprite]*SpriteAnimator
    tweens      []*Tween
    animations  map[string]*AnimationDef
    globalSpeed float32
    paused      bool
}
```

### Core Functions
```go
// Animation manager
func NewAnimationManager() *AnimationManager
func (am *AnimationManager) Update()
func (am *AnimationManager) SetGlobalSpeed(speed float32)
func (am *AnimationManager) SetPaused(paused bool)

// Animation definitions
func (am *AnimationManager) LoadAnimation(name string, anim *AnimationDef)
func (am *AnimationManager) GetAnimation(name string) *AnimationDef
func CreateAnimationDef(name string, frames []AnimationFrame, mode int) *AnimationDef

// Sprite animation control
func (am *AnimationManager) GetAnimator(sprite *sprites.Sprite) *SpriteAnimator
func (sa *SpriteAnimator) PlayAnimation(animName string) error
func (sa *SpriteAnimator) SetAnimation(anim *AnimationDef)
func (sa *SpriteAnimator) Stop()
func (sa *SpriteAnimator) Pause()
func (sa *SpriteAnimator) Resume()
func (sa *SpriteAnimator) SetSpeed(speed float32)

// Animation state queries
func (sa *SpriteAnimator) IsPlaying() bool
func (sa *SpriteAnimator) IsFinished() bool
func (sa *SpriteAnimator) GetCurrentFrame() int
func (sa *SpriteAnimator) GetCurrentAnimation() string
func (sa *SpriteAnimator) GetProgress() float32

// Animation blending and transitions
func (sa *SpriteAnimator) BlendToAnimation(animName string, duration int) error
func (sa *SpriteAnimator) QueueAnimation(animName string)
func (sa *SpriteAnimator) SetTransitionSpeed(speed float32)

// Tweening system
func (am *AnimationManager) TweenTo(target interface{}, property string, endValue float32, duration int, easeType int) *Tween
func (am *AnimationManager) MoveTo(sprite *sprites.Sprite, x, y int, duration int, easeType int) *Tween
func (am *AnimationManager) ScaleTo(sprite *sprites.Sprite, scaleX, scaleY float32, duration int, easeType int) *Tween
func (am *AnimationManager) RotateTo(sprite *sprites.Sprite, angle float32, duration int, easeType int) *Tween
func (am *AnimationManager) FadeTo(sprite *sprites.Sprite, alpha float32, duration int, easeType int) *Tween

// Tween control
func (t *Tween) Stop()
func (t *Tween) SetOnComplete(callback func())
func (t *Tween) IsActive() bool

// Utility functions
func (am *AnimationManager) StopAllAnimations()
func (am *AnimationManager) PauseAllAnimations()
func (am *AnimationManager) GetActiveAnimationCount() int
func LoadAnimationsFromData(data []byte) (map[string]*AnimationDef, error)
```

## Implementation Details

### Step 1: Animation Data Structures
Create `lib/animation/definitions.go`:
- AnimationFrame and AnimationDef structures
- Animation loading and storage
- Animation validation and optimization

### Step 2: Sprite Animator
Create `lib/animation/animator.go`:
- SpriteAnimator for individual sprite animation state
- Frame timing and playback logic
- Animation events and callbacks

### Step 3: Animation Manager
Create `lib/animation/manager.go`:
- Central coordinator for all animations
- Efficient batch updates
- Resource management for animation data

### Step 4: Tweening System
Create `lib/animation/tweening.go`:
- Smooth interpolation animations
- Multiple easing functions
- Property-based animation system

### Step 5: Animation Blending
Create `lib/animation/blending.go`:
- Smooth transitions between animations
- Animation priority system
- Cross-fade effects

### Step 6: Utilities and Helpers
Create `lib/animation/utils.go`:
- Animation creation helpers
- Debug visualization
- Performance profiling tools

## Testing Strategy

### Unit Tests
```go
func TestAnimationCreation(t *testing.T) {
    frames := []AnimationFrame{
        {TileIndex: 0, Duration: 10},
        {TileIndex: 1, Duration: 10},
        {TileIndex: 2, Duration: 10},
    }
    
    anim := CreateAnimationDef("test", frames, ANIM_LOOP)
    assert.Equal(t, "test", anim.Name)
    assert.Equal(t, 3, len(anim.Frames))
    assert.Equal(t, ANIM_LOOP, anim.Mode)
}

func TestSpriteAnimation(t *testing.T) {
    manager := NewAnimationManager()
    spriteManager := sprites.NewSpriteManager()
    sprite := spriteManager.CreateSprite(100, 100, 0, sprites.SIZE_16x16)
    
    // Create test animation
    frames := []AnimationFrame{
        {TileIndex: 0, Duration: 5},
        {TileIndex: 1, Duration: 5},
    }
    anim := CreateAnimationDef("walk", frames, ANIM_LOOP)
    manager.LoadAnimation("walk", anim)
    
    // Play animation
    animator := manager.GetAnimator(sprite)
    err := animator.PlayAnimation("walk")
    assert.NoError(t, err)
    assert.True(t, animator.IsPlaying())
    
    // Update animation
    for i := 0; i < 5; i++ {
        manager.Update()
    }
    
    // Should be on second frame now
    assert.Equal(t, 1, animator.GetCurrentFrame())
}

func TestTweening(t *testing.T) {
    manager := NewAnimationManager()
    spriteManager := sprites.NewSpriteManager()
    sprite := spriteManager.CreateSprite(0, 0, 0, sprites.SIZE_16x16)
    
    // Start tween
    tween := manager.MoveTo(sprite, 100, 50, 60, EASE_LINEAR)
    assert.True(t, tween.IsActive())
    
    // Update half duration
    for i := 0; i < 30; i++ {
        manager.Update()
    }
    
    // Should be halfway
    x, y := sprite.GetPosition()
    assert.InDelta(t, 50, x, 2) // Allow small floating point error
    assert.InDelta(t, 25, y, 2)
    
    // Complete tween
    for i := 0; i < 30; i++ {
        manager.Update()
    }
    
    x, y = sprite.GetPosition()
    assert.Equal(t, 100, x)
    assert.Equal(t, 50, y)
    assert.False(t, tween.IsActive())
}
```

### Performance Tests
```go
func BenchmarkAnimationUpdate(b *testing.B) {
    manager := NewAnimationManager()
    spriteManager := sprites.NewSpriteManager()
    
    // Create many animated sprites
    sprites := make([]*sprites.Sprite, 100)
    for i := 0; i < 100; i++ {
        sprites[i] = spriteManager.CreateSprite(i*2, i*2, 0, sprites.SIZE_16x16)
        animator := manager.GetAnimator(sprites[i])
        animator.PlayAnimation("walk")
    }
    
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        manager.Update()
    }
}
```

## Example Program
```go
package main

import (
    "github.com/matheusmortatti/gba-go/lib/animation"
    "github.com/matheusmortatti/gba-go/lib/sprites"
    "github.com/matheusmortatti/gba-go/lib/memory"
    "github.com/matheusmortatti/gba-go/lib/registers"
    "github.com/matheusmortatti/gba-go/lib/video"
    "github.com/matheusmortatti/gba-go/lib/input"
)

func main() {
    // Initialize systems
    spriteManager := sprites.NewSpriteManager()
    animManager := animation.NewAnimationManager()
    
    // Set up video mode
    registers.Lcd.DISPCNT.SetBits(memory.MODE_0)
    registers.Lcd.DISPCNT.SetBits(1 << 12) // Enable sprites
    
    // Load sprite graphics and palettes
    setupGraphics()
    
    // Create character animations
    setupAnimations(animManager)
    
    // Create player character
    player := spriteManager.CreateSprite(120, 80, 0, sprites.SIZE_16x16)
    playerAnimator := animManager.GetAnimator(player)
    playerAnimator.PlayAnimation("idle")
    
    // Create NPCs with different animations
    npcs := createNPCs(spriteManager, animManager)
    
    frame := 0
    playerState := "idle"
    
    for {
        video.VSync()
        input.Poll()
        
        // Handle player input and animation
        newState := handlePlayerInput(player, playerAnimator)
        if newState != playerState {
            playerState = newState
            // Smooth transition between animations
            playerAnimator.BlendToAnimation(playerState, 10)
        }
        
        // Update NPC behaviors
        updateNPCs(npcs, animManager, frame)
        
        // Demonstrate tweening
        if input.BtnClicked(input.KeyA) {
            demonstrateTweening(animManager, player)
        }
        
        // Update all animations
        animManager.Update()
        spriteManager.Update()
        
        frame++
    }
}

func setupAnimations(manager *animation.AnimationManager) {
    // Player idle animation
    idleFrames := []animation.AnimationFrame{
        {TileIndex: 0, Duration: 30},
        {TileIndex: 1, Duration: 30},
    }
    idleAnim := animation.CreateAnimationDef("idle", idleFrames, animation.ANIM_LOOP)
    manager.LoadAnimation("idle", idleAnim)
    
    // Player walk animation
    walkFrames := []animation.AnimationFrame{
        {TileIndex: 2, Duration: 8},
        {TileIndex: 3, Duration: 8},
        {TileIndex: 4, Duration: 8},
        {TileIndex: 5, Duration: 8},
    }
    walkAnim := animation.CreateAnimationDef("walk", walkFrames, animation.ANIM_LOOP)
    manager.LoadAnimation("walk", walkAnim)
    
    // Player jump animation
    jumpFrames := []animation.AnimationFrame{
        {TileIndex: 6, Duration: 5, OffsetY: -2},
        {TileIndex: 7, Duration: 10, OffsetY: -4},
        {TileIndex: 8, Duration: 15, OffsetY: -2},
        {TileIndex: 6, Duration: 5},
    }
    jumpAnim := animation.CreateAnimationDef("jump", jumpFrames, animation.ANIM_ONCE)
    // Add landing event
    jumpAnim.Events = []animation.AnimationEvent{
        {Frame: 3, Callback: func(sprite *sprites.Sprite) {
            // Play landing sound effect here
        }},
    }
    manager.LoadAnimation("jump", jumpAnim)
    
    // Enemy patrol animation
    patrolFrames := []animation.AnimationFrame{
        {TileIndex: 10, Duration: 15},
        {TileIndex: 11, Duration: 15},
        {TileIndex: 12, Duration: 15},
        {TileIndex: 11, Duration: 15},
    }
    patrolAnim := animation.CreateAnimationDef("patrol", patrolFrames, animation.ANIM_LOOP)
    manager.LoadAnimation("patrol", patrolAnim)
    
    // Spinning coin animation
    coinFrames := []animation.AnimationFrame{
        {TileIndex: 20, Duration: 6},
        {TileIndex: 21, Duration: 6},
        {TileIndex: 22, Duration: 6},
        {TileIndex: 23, Duration: 6},
    }
    coinAnim := animation.CreateAnimationDef("spin", coinFrames, animation.ANIM_LOOP)
    coinAnim.Speed = 1.5 // Faster spinning
    manager.LoadAnimation("spin", coinAnim)
}

func handlePlayerInput(player *sprites.Sprite, animator *animation.SpriteAnimator) string {
    x, y := player.GetPosition()
    moved := false
    
    if input.BtnDown(input.KeyLeft) {
        player.SetPosition(x-2, y)
        player.SetFlip(true, false)
        moved = true
    } else if input.BtnDown(input.KeyRight) {
        player.SetPosition(x+2, y)
        player.SetFlip(false, false)
        moved = true
    }
    
    if input.BtnDown(input.KeyUp) {
        player.SetPosition(x, y-2)
        moved = true
    } else if input.BtnDown(input.KeyDown) {
        player.SetPosition(x, y+2)
        moved = true
    }
    
    if input.BtnClicked(input.KeyB) {
        return "jump"
    }
    
    if moved {
        return "walk"
    } else {
        return "idle"
    }
}

func createNPCs(spriteManager *sprites.SpriteManager, animManager *animation.AnimationManager) []*NPCData {
    npcs := make([]*NPCData, 3)
    
    // Patrolling enemy
    enemy := spriteManager.CreateSprite(50, 50, 10, sprites.SIZE_16x16)
    enemyAnimator := animManager.GetAnimator(enemy)
    enemyAnimator.PlayAnimation("patrol")
    npcs[0] = &NPCData{
        sprite: enemy,
        animator: enemyAnimator,
        aiType: "patrol",
        startX: 50,
        patrolRange: 80,
    }
    
    // Spinning coins
    for i := 1; i < 3; i++ {
        coin := spriteManager.CreateSprite(100+i*50, 30, 20, sprites.SIZE_8x8)
        coinAnimator := animManager.GetAnimator(coin)
        coinAnimator.PlayAnimation("spin")
        npcs[i] = &NPCData{
            sprite: coin,
            animator: coinAnimator,
            aiType: "static",
        }
    }
    
    return npcs
}

type NPCData struct {
    sprite       *sprites.Sprite
    animator     *animation.SpriteAnimator
    aiType       string
    startX       int
    patrolRange  int
    direction    int
}

func updateNPCs(npcs []*NPCData, animManager *animation.AnimationManager, frame int) {
    for _, npc := range npcs {
        switch npc.aiType {
        case "patrol":
            updatePatrolAI(npc, frame)
        case "static":
            // Coins just animate in place
            // Add slight bobbing motion
            x, y := npc.sprite.GetPosition()
            newY := y + int(math.Sin(float64(frame)*0.1))
            npc.sprite.SetPosition(x, newY)
        }
    }
}

func updatePatrolAI(npc *NPCData, frame int) {
    x, y := npc.sprite.GetPosition()
    
    // Simple patrol movement
    if npc.direction == 0 {
        x++
        if x >= npc.startX + npc.patrolRange {
            npc.direction = 1
            npc.sprite.SetFlip(true, false)
        }
    } else {
        x--
        if x <= npc.startX {
            npc.direction = 0
            npc.sprite.SetFlip(false, false)
        }
    }
    
    npc.sprite.SetPosition(x, y)
}

func demonstrateTweening(manager *animation.AnimationManager, player *sprites.Sprite) {
    // Create a smooth movement tween with bounce easing
    x, y := player.GetPosition()
    targetX := (x + 50) % memory.SCREEN_WIDTH
    targetY := (y + 30) % memory.SCREEN_HEIGHT
    
    tween := manager.MoveTo(player, targetX, targetY, 30, animation.EASE_BOUNCE)
    tween.SetOnComplete(func() {
        // Could play a sound or effect when movement completes
    })
}

func setupGraphics() {
    // Set up sprite palettes and load tile graphics
    // This would integrate with palette and tile systems
}
```

## Advanced Features

### State Machine Integration
```go
type AnimationStateMachine struct {
    states      map[string]*AnimationState
    current     string
    transitions map[string]map[string]bool
}

type AnimationState struct {
    animation   string
    enterAction func()
    exitAction  func()
}

func (asm *AnimationStateMachine) TransitionTo(state string) bool {
    // Validate transition and change state
}
```

### Animation Compression
```go
func CompressAnimation(anim *AnimationDef) *CompressedAnimation {
    // Compress repetitive frame data
    // Delta compression for similar frames
}
```

## File Structure
```
lib/animation/
├── definitions.go   // Animation data structures
├── animator.go      // Individual sprite animation
├── manager.go       // Animation coordination
├── tweening.go      // Smooth interpolation
├── blending.go      // Animation transitions
└── utils.go         // Helper functions
```

## Integration Points
- Extend sprite system from Task 04
- Use timer utilities from Task 09 for precise timing
- Integrate with sound system from Task 10 for audio events
- Coordinate with game framework from Task 12

## Resources
- [Game Programming Patterns - Animation](https://gameprogrammingpatterns.com/state.html)
- [Real-Time Rendering Animation Techniques](https://www.realtimerendering.com/)
- [GBA Animation Optimization Techniques](https://www.coranac.com/tonc/text/gba.htm)

## Success Criteria
- Smooth 60fps animation playback with many sprites
- Animation blending works without visual artifacts
- Tweening system provides smooth interpolations
- Memory usage remains efficient with many animations
- Animation events trigger correctly and on time
- Example program demonstrates all animation features
- Performance suitable for complex games
- Comprehensive test coverage (>90%)