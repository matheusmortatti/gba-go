# Task 09: Timer Utilities System

## Objective
Implement a comprehensive timer system for the GBA that provides precise timing, frame rate management, delays, scheduling, and performance measurement utilities using the hardware timers and software timing mechanisms.

## Background
The GBA has 4 hardware timers (TM0-TM3) that can be used for precise timing, frame rate control, audio generation, and general timing utilities. Proper timer management is essential for game timing, animation, and performance optimization.

### GBA Timer System
- **4 Hardware Timers**: TM0-TM3 with cascade support
- **Timer Frequencies**: 16.78MHz, 262.21kHz, 65.536kHz, 16.384kHz
- **Timer Modes**: Free-running, interrupt-generating, cascaded
- **16-bit Counters**: 0-65535 count range with reload values

## Requirements

### Functional Requirements
1. High-precision timing and delays
2. Frame rate measurement and control
3. Timer-based scheduling system
4. Performance profiling utilities
5. Audio timing support
6. Game timing abstractions (delta time, fixed timestep)
7. Timer interrupt handling
8. Stopwatch and timeout functionality

### Technical Requirements
- Hardware timer register access
- Interrupt-based timing
- Microsecond precision timing
- Integration with VSync and game loops
- Memory-efficient timer management

## API Design

### Constants
```go
package timers

const (
    // Hardware timers
    TIMER0 = 0
    TIMER1 = 1
    TIMER2 = 2
    TIMER3 = 3
    
    // Timer frequencies (cycles per second)
    FREQ_1024_CYCLE  = 0  // 16.384 KHz
    FREQ_64_CYCLE    = 1  // 262.21 KHz  
    FREQ_256_CYCLE   = 2  // 65.536 KHz
    FREQ_1_CYCLE     = 3  // 16.78 MHz
    
    // Timer control flags
    TIMER_IRQ        = 0x40   // Interrupt enable
    TIMER_CASCADE    = 0x04   // Cascade mode
    TIMER_ENABLE     = 0x80
    
    // Common timing values
    CYCLES_PER_SECOND = 16777216  // 16.78 MHz
    CYCLES_PER_FRAME  = 280896    // 59.73 FPS
    TICKS_PER_SECOND  = 1024      // Using 1024-cycle timer
    MICROSECONDS_PER_SECOND = 1000000
)
```

### Data Structures
```go
// Timer represents a hardware timer
type Timer struct {
    id        int
    frequency int
    active    bool
    reload    uint16
    callback  func()
}

// GameTimer provides high-level timing functionality
type GameTimer struct {
    startTime    uint64
    lastTime     uint64
    deltaTime    float32
    totalTime    float32
    frameCount   uint32
    fps          float32
    paused       bool
}

// Scheduler manages timed events
type Scheduler struct {
    events     []*ScheduledEvent
    currentTime uint64
}

// ScheduledEvent represents a timed callback
type ScheduledEvent struct {
    id          int
    executeTime uint64
    interval    uint64
    callback    func()
    repeat      bool
    active      bool
}

// Profiler measures performance
type Profiler struct {
    samples     map[string][]uint64
    startTimes  map[string]uint64
    enabled     bool
}

// DelayManager handles various delay operations
type DelayManager struct {
    delays      []*DelayOperation
    frameDelays []*FrameDelay
}

// DelayOperation represents a time-based delay
type DelayOperation struct {
    duration    uint64
    startTime   uint64
    callback    func()
    completed   bool
}

// FrameDelay represents a frame-based delay
type FrameDelay struct {
    frames      int
    remaining   int
    callback    func()
    completed   bool
}
```

### Core Functions
```go
// Timer system initialization
func InitTimers()
func GetHardwareTimer(id int) *Timer
func SetTimerFrequency(id int, freq int)
func StartTimer(id int, reload uint16, callback func())
func StopTimer(id int)

// Game timing
func NewGameTimer() *GameTimer
func (gt *GameTimer) Update()
func (gt *GameTimer) GetDeltaTime() float32
func (gt *GameTimer) GetTotalTime() float32
func (gt *GameTimer) GetFPS() float32
func (gt *GameTimer) Pause()
func (gt *GameTimer) Resume()
func (gt *GameTimer) Reset()

// Scheduling system
func NewScheduler() *Scheduler
func (s *Scheduler) Update()
func (s *Scheduler) ScheduleOnce(delay uint64, callback func()) int
func (s *Scheduler) ScheduleRepeating(interval uint64, callback func()) int
func (s *Scheduler) Cancel(eventId int)
func (s *Scheduler) GetCurrentTime() uint64

// Delay operations
func NewDelayManager() *DelayManager
func (dm *DelayManager) Update()
func (dm *DelayManager) DelayMicroseconds(us uint64, callback func())
func (dm *DelayManager) DelayMilliseconds(ms uint64, callback func())
func (dm *DelayManager) DelayFrames(frames int, callback func())
func (dm *DelayManager) DelaySeconds(seconds float32, callback func())

// Profiling
func NewProfiler() *Profiler
func (p *Profiler) Begin(name string)
func (p *Profiler) End(name string)
func (p *Profiler) GetAverageTime(name string) float32
func (p *Profiler) GetReport() string
func (p *Profiler) Reset()

// Utility functions
func GetSystemTime() uint64
func MicrosecondsToTicks(us uint64) uint64
func TicksToMicroseconds(ticks uint64) uint64
func WaitMicroseconds(us uint64)
func WaitFrames(frames int)
func MeasureFrameRate() float32

// High-precision timing
func GetHighPrecisionTime() uint64
func SetupPrecisionTimer() 
func TimeFunctionCall(fn func()) uint64
```

## Implementation Details

### Step 1: Hardware Timer Interface
Create `lib/timers/hardware.go`:
- Direct timer register access
- Timer configuration and control
- Interrupt handling setup

### Step 2: Game Timer System
Create `lib/timers/gametimer.go`:
- Delta time calculation
- Frame rate measurement
- Game timing abstractions

### Step 3: Scheduler Implementation
Create `lib/timers/scheduler.go`:
- Event scheduling and management
- Callback execution system
- Event prioritization

### Step 4: Delay Manager
Create `lib/timers/delays.go`:
- Various delay mechanisms
- Frame-based and time-based delays
- Asynchronous delay operations

### Step 5: Performance Profiler
Create `lib/timers/profiler.go`:
- Code execution timing
- Performance measurement tools
- Statistical analysis

### Step 6: Utility Functions
Create `lib/timers/utils.go`:
- Common timing utilities
- Conversion functions
- Debug helpers

## Testing Strategy

### Unit Tests
```go
func TestGameTimer(t *testing.T) {
    timer := NewGameTimer()
    
    // Simulate frame updates
    for i := 0; i < 60; i++ {
        timer.Update()
        time.Sleep(16 * time.Millisecond) // ~60 FPS
    }
    
    // Check timing accuracy
    assert.InDelta(t, 1.0, timer.GetTotalTime(), 0.1) // ~1 second
    assert.InDelta(t, 60.0, timer.GetFPS(), 5.0)      // ~60 FPS
}

func TestScheduler(t *testing.T) {
    scheduler := NewScheduler()
    executed := false
    
    // Schedule event
    eventId := scheduler.ScheduleOnce(1000, func() {
        executed = true
    })
    
    // Update scheduler
    for i := 0; i < 1100; i++ {
        scheduler.Update()
    }
    
    assert.True(t, executed)
}

func TestDelayManager(t *testing.T) {
    manager := NewDelayManager()
    completed := false
    
    manager.DelayMilliseconds(100, func() {
        completed = true
    })
    
    start := time.Now()
    for !completed && time.Since(start) < 200*time.Millisecond {
        manager.Update()
        time.Sleep(1 * time.Millisecond)
    }
    
    assert.True(t, completed)
}
```

## Example Program
```go
package main

import (
    \"github.com/matheusmortatti/gba-go/lib/timers\"
    \"github.com/matheusmortatti/gba-go/lib/video\"
    \"github.com/matheusmortatti/gba-go/lib/input\"
    \"github.com/matheusmortatti/gba-go/lib/registers\"
)

func main() {
    // Initialize timer systems
    timers.InitTimers()
    gameTimer := timers.NewGameTimer()
    scheduler := timers.NewScheduler()
    delayManager := timers.NewDelayManager()
    profiler := timers.NewProfiler()
    
    // Schedule recurring events
    scheduler.ScheduleRepeating(60, func() {
        // Print FPS every second (60 frames)
        println(\"FPS:\", gameTimer.GetFPS())
    })
    
    // Schedule one-time events
    scheduler.ScheduleOnce(300, func() {
        println(\"5 seconds elapsed!\")
    })
    
    frame := 0
    
    for {
        profiler.Begin(\"frame\")
        
        video.VSync()
        input.Poll()
        
        // Update timing systems
        gameTimer.Update()
        scheduler.Update()
        delayManager.Update()
        
        // Handle input with timing
        handleTimedInput(delayManager, gameTimer)
        
        // Demonstrate timing utilities
        if frame%300 == 0 { // Every 5 seconds
            demonstrateTimingFeatures(scheduler, delayManager, profiler)
        }
        
        profiler.End(\"frame\")
        
        // Print performance stats periodically
        if frame%3600 == 0 { // Every minute
            println(profiler.GetReport())
            profiler.Reset()
        }
        
        frame++
    }
}

func handleTimedInput(dm *timers.DelayManager, gt *timers.GameTimer) {
    if input.BtnClicked(input.KeyA) {
        println(\"A pressed at time:\", gt.GetTotalTime())
        
        // Delay action by 1 second
        dm.DelaySeconds(1.0, func() {
            println(\"Delayed action executed!\")
        })
    }
    
    if input.BtnClicked(input.KeyB) {
        // Measure button response time
        startTime := timers.GetHighPrecisionTime()
        
        dm.DelayFrames(1, func() {
            endTime := timers.GetHighPrecisionTime()
            responseTime := timers.TicksToMicroseconds(endTime - startTime)
            println(\"Button response time:\", responseTime, \"microseconds\")
        })
    }
}

func demonstrateTimingFeatures(scheduler *timers.Scheduler, dm *timers.DelayManager, profiler *timers.Profiler) {
    profiler.Begin(\"demo\")
    
    // Demonstrate cascading delays
    dm.DelayMilliseconds(100, func() {
        println(\"First delay completed\")
        
        dm.DelayMilliseconds(200, func() {
            println(\"Second delay completed\")
            
            dm.DelayFrames(30, func() {
                println(\"Frame delay completed\")\n            })\n        })\n    })\n    \n    // Schedule event chain\n    scheduler.ScheduleOnce(60, func() {\n        println(\"Scheduled event 1\")\n        \n        scheduler.ScheduleOnce(120, func() {\n            println(\"Scheduled event 2\")\n        })\n    })\n    \n    profiler.End(\"demo\")\n}\n```\n\n## File Structure\n```\nlib/timers/\n├── hardware.go     // Hardware timer interface\n├── gametimer.go    // Game timing abstractions\n├── scheduler.go    // Event scheduling system\n├── delays.go       // Delay management\n├── profiler.go     // Performance profiling\n└── utils.go        // Utility functions\n```\n\n## Integration Points\n- Support animation system from Task 05\n- Enable DMA timing from Task 08\n- Integrate with sound system from Task 10\n- Coordinate with game framework from Task 12\n\n## Resources\n- [GBATEK Timer Control](https://problemkaputt.de/gbatek.htm#gbatimers)\n- [Tonc Timers](https://www.coranac.com/tonc/text/timers.htm)\n- [GBA Timer Programming](https://www.cs.rit.edu/~tjh8300/CowBite/CowBiteSpec.htm#Timers)\n\n## Success Criteria\n- Accurate timing measurements and delays\n- Stable frame rate measurement\n- Precise event scheduling\n- Performance profiling provides useful metrics\n- Integration with game systems works smoothly\n- Example program demonstrates all timing features\n- Timer interrupts work correctly\n- Comprehensive test coverage (>90%)