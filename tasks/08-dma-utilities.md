# Task 08: DMA Utilities System

## Objective
Implement a comprehensive Direct Memory Access (DMA) system for the GBA that provides high-performance memory operations, automated transfers, and hardware-accelerated bulk operations for graphics, audio, and general memory management.

## Background
The GBA has 4 DMA channels (DMA0-DMA3) that can perform memory transfers without CPU intervention. DMA is essential for performance-critical operations like graphics updates, audio streaming, and bulk memory operations.

### GBA DMA System
- **4 DMA Channels**: DMA0-DMA3 with different capabilities and restrictions
- **Transfer Modes**: 16-bit and 32-bit transfers
- **Timing Modes**: Immediate, VBlank, HBlank, special (audio/video)
- **Address Modes**: Increment, decrement, fixed, increment-reload
- **Transfer Sizes**: Up to 16KB (DMA0-2) or 64KB (DMA3) per operation

### DMA Channel Capabilities
- **DMA0**: Internal memory only, highest priority
- **DMA1**: Any memory, medium priority
- **DMA2**: Any memory, medium priority  
- **DMA3**: Any memory, lowest priority, largest transfers

## Requirements

### Functional Requirements
1. High-level DMA operation wrappers for common tasks
2. Safe memory transfer with bounds checking
3. Asynchronous and synchronous transfer modes
4. Integration with VBlank/HBlank timing
5. Audio DMA for sound streaming
6. Graphics DMA for fast VRAM operations
7. Queue system for multiple DMA operations
8. Performance monitoring and optimization

### Technical Requirements
- Direct hardware register manipulation
- Memory alignment handling
- Interrupt-based completion detection
- Priority management between channels
- Integration with existing memory systems
- Error handling and recovery

## API Design

### Constants
```go
package dma

const (
    // DMA channels
    DMA0 = 0
    DMA1 = 1
    DMA2 = 2
    DMA3 = 3
    
    // Transfer sizes
    TRANSFER_16BIT = 0
    TRANSFER_32BIT = 1
    
    // Timing modes
    TIMING_IMMEDIATE = 0
    TIMING_VBLANK    = 1
    TIMING_HBLANK    = 2
    TIMING_SPECIAL   = 3
    
    // Source address modes
    SRC_INCREMENT = 0
    SRC_DECREMENT = 1
    SRC_FIXED     = 2
    SRC_PROHIBITED = 3
    
    // Destination address modes
    DST_INCREMENT        = 0
    DST_DECREMENT        = 1
    DST_FIXED            = 2
    DST_INCREMENT_RELOAD = 3
    
    // Control flags
    DMA_REPEAT     = 0x0200
    DMA_WORD       = 0x0400  // 32-bit transfer
    DMA_DREQ_ON    = 0x0800  // DRQ synchronization
    DMA_TIMING_MASK = 0x3000
    DMA_IRQ        = 0x4000  // Interrupt request
    DMA_ENABLE     = 0x8000
    
    // Maximum transfer sizes
    MAX_TRANSFER_DMA0 = 0x4000  // 16KB
    MAX_TRANSFER_DMA1 = 0x4000  // 16KB
    MAX_TRANSFER_DMA2 = 0x4000  // 16KB
    MAX_TRANSFER_DMA3 = 0x10000 // 64KB
    
    // Common patterns
    PATTERN_FILL16 = 0x00000000  // Fill with 16-bit value
    PATTERN_FILL32 = 0x00000000  // Fill with 32-bit value
)
```

### Data Structures
```go
// DMATransfer represents a single DMA operation
type DMATransfer struct {
    channel     int
    source      uintptr
    dest        uintptr
    count       uint32
    control     uint16
    onComplete  func()
    active      bool
    priority    int
}

// DMAManager coordinates all DMA operations
type DMAManager struct {
    channels     [4]*DMAChannel
    queue        []*DMATransfer
    totalOps     int
    totalBytes   uint64
    activeOps    int
}

// DMAChannel represents a hardware DMA channel
type DMAChannel struct {
    id           int
    busy         bool
    currentOp    *DMATransfer
    maxTransfer  uint32
    capabilities uint32
}

// DMAConfig holds configuration for DMA operations
type DMAConfig struct {
    transferSize int  // 16 or 32 bit
    timing       int  // When to start transfer
    srcMode      int  // Source address mode
    dstMode      int  // Destination address mode
    repeat       bool // Repeat transfer
    interrupt    bool // Generate interrupt on completion
}
```

### Core Functions
```go
// DMA Manager
func NewDMAManager() *DMAManager
func (dm *DMAManager) GetChannel(channel int) *DMAChannel
func (dm *DMAManager) Update() // Check for completed operations
func (dm *DMAManager) WaitForCompletion()
func (dm *DMAManager) IsChannelBusy(channel int) bool

// High-level operations
func (dm *DMAManager) Copy(src, dst uintptr, size uint32, config DMAConfig) error
func (dm *DMAManager) Fill16(dst uintptr, value uint16, count uint32) error
func (dm *DMAManager) Fill32(dst uintptr, value uint32, count uint32) error
func (dm *DMAManager) FastCopy(src, dst uintptr, size uint32) error
func (dm *DMAManager) CopyAsync(src, dst uintptr, size uint32, onComplete func()) error

// Graphics-specific operations
func (dm *DMAManager) CopyToVRAM(src uintptr, vramOffset uintptr, size uint32) error
func (dm *DMAManager) FillVRAM(vramOffset uintptr, value uint16, count uint32) error
func (dm *DMAManager) LoadPalette(src uintptr, paletteType int, index int, colors int) error
func (dm *DMAManager) LoadTiles(src uintptr, charBlock int, tileIndex int, tileCount int) error
func (dm *DMAManager) UpdateOAM(src uintptr, oamOffset int, spriteCount int) error

// Audio-specific operations
func (dm *DMAManager) SetupAudioDMA(channel int, buffer uintptr, size uint32) error
func (dm *DMAManager) StartAudioStream(channel int) error
func (dm *DMAManager) StopAudioStream(channel int) error

// Synchronization operations
func (dm *DMAManager) CopyOnVBlank(src, dst uintptr, size uint32) error
func (dm *DMAManager) CopyOnHBlank(src, dst uintptr, size uint32) error
func (dm *DMAManager) ScheduleTransfer(transfer *DMATransfer) error

// Utility functions
func (dm *DMAManager) GetOptimalChannel(size uint32, timing int) int
func (dm *DMAManager) ValidateTransfer(src, dst uintptr, size uint32) error
func AlignAddress(addr uintptr, alignment int) uintptr
func CalculateTransferCount(size uint32, transferSize int) uint32

// Performance monitoring
func (dm *DMAManager) GetStatistics() DMAStatistics
func (dm *DMAManager) ResetStatistics()

type DMAStatistics struct {
    TotalOperations uint64
    TotalBytes      uint64
    AverageSpeed    float64
    ChannelUsage    [4]uint64
}
```

## Implementation Details

### Step 1: Hardware Interface
Create `lib/dma/hardware.go`:
- Direct DMA register access
- Channel configuration and control
- Interrupt handling for completion detection

### Step 2: DMA Manager
Create `lib/dma/manager.go`:
- Central coordination of all DMA operations
- Channel allocation and priority management
- Queue system for multiple operations

### Step 3: High-Level Operations
Create `lib/dma/operations.go`:
- Common DMA operation wrappers
- Parameter validation and optimization
- Error handling and recovery

### Step 4: Graphics Integration
Create `lib/dma/graphics.go`:
- VRAM-specific DMA operations
- Palette and tile loading utilities
- OAM update optimizations

### Step 5: Audio Integration
Create `lib/dma/audio.go`:
- Audio buffer streaming
- Sound effect transfer utilities
- Audio timing synchronization

### Step 6: Utilities and Optimization
Create `lib/dma/utils.go`:
- Performance monitoring tools
- Memory alignment utilities
- Debug and profiling helpers

## Testing Strategy

### Unit Tests
```go
func TestDMAManager(t *testing.T) {
    manager := NewDMAManager()
    
    // Test channel availability
    for i := 0; i < 4; i++ {
        channel := manager.GetChannel(i)
        assert.NotNil(t, channel)
        assert.Equal(t, i, channel.id)
        assert.False(t, channel.busy)
    }
}

func TestDMATransfer(t *testing.T) {
    manager := NewDMAManager()
    
    // Create test data
    src := make([]uint16, 100)
    dst := make([]uint16, 100)
    for i := range src {
        src[i] = uint16(i)
    }
    
    // Perform DMA transfer
    config := DMAConfig{
        transferSize: TRANSFER_16BIT,
        timing:       TIMING_IMMEDIATE,
        srcMode:      SRC_INCREMENT,
        dstMode:      DST_INCREMENT,
    }
    
    err := manager.Copy(uintptr(unsafe.Pointer(&src[0])), 
                       uintptr(unsafe.Pointer(&dst[0])), 
                       200, config) // 100 * 2 bytes
    assert.NoError(t, err)
    
    // Wait for completion
    manager.WaitForCompletion()
    
    // Verify data
    for i := range src {
        assert.Equal(t, src[i], dst[i])\n    }\n}\n\nfunc TestDMAFill(t *testing.T) {\n    manager := NewDMAManager()\n    \n    // Create destination buffer\n    dst := make([]uint16, 100)\n    \n    // Fill with pattern\n    err := manager.Fill16(uintptr(unsafe.Pointer(&dst[0])), 0x1234, 100)\n    assert.NoError(t, err)\n    \n    manager.WaitForCompletion()\n    \n    // Verify fill\n    for _, value := range dst {\n        assert.Equal(t, uint16(0x1234), value)\n    }\n}\n\nfunc TestVRAMOperations(t *testing.T) {\n    manager := NewDMAManager()\n    \n    // Test VRAM fill\n    err := manager.FillVRAM(0, 0x7FFF, 1000) // Fill with white\n    assert.NoError(t, err)\n    \n    // Test tile loading\n    tileData := make([]uint8, 32*10) // 10 tiles\n    err = manager.LoadTiles(uintptr(unsafe.Pointer(&tileData[0])), 0, 0, 10)\n    assert.NoError(t, err)\n}\n```

### Performance Tests
```go
func BenchmarkDMACopy(b *testing.B) {\n    manager := NewDMAManager()\n    src := make([]uint32, 1000)\n    dst := make([]uint32, 1000)\n    \n    config := DMAConfig{\n        transferSize: TRANSFER_32BIT,\n        timing:       TIMING_IMMEDIATE,\n        srcMode:      SRC_INCREMENT,\n        dstMode:      DST_INCREMENT,\n    }\n    \n    b.ResetTimer()\n    for i := 0; i < b.N; i++ {\n        manager.Copy(uintptr(unsafe.Pointer(&src[0])),\n                    uintptr(unsafe.Pointer(&dst[0])),\n                    4000, config)\n        manager.WaitForCompletion()\n    }\n}\n\nfunc BenchmarkCPUCopy(b *testing.B) {\n    src := make([]uint32, 1000)\n    dst := make([]uint32, 1000)\n    \n    b.ResetTimer()\n    for i := 0; i < b.N; i++ {\n        copy(dst, src)\n    }\n}\n```

## Example Program
```go
package main

import (
    \"github.com/matheusmortatti/gba-go/lib/dma\"
    \"github.com/matheusmortatti/gba-go/lib/memory\"
    \"github.com/matheusmortatti/gba-go/lib/vram\"
    \"github.com/matheusmortatti/gba-go/lib/palette\"
    \"github.com/matheusmortatti/gba-go/lib/registers\"
    \"github.com/matheusmortatti/gba-go/lib/video\"
    \"github.com/matheusmortatti/gba-go/lib/input\"
)

func main() {
    // Initialize DMA system
    dmaManager := dma.NewDMAManager()
    vramManager := vram.NewVRAMManager(memory.MODE_3)
    
    // Set video mode 3
    registers.Lcd.DISPCNT.SetBits(memory.MODE_3)
    registers.Lcd.DISPCNT.SetBits(1 << 10) // Enable BG2
    
    // Create test graphics data
    setupGraphicsData(dmaManager)\n    \n    frame := 0\n    \n    for {\n        video.VSync()\n        input.Poll()\n        \n        // Demonstrate different DMA operations\n        if input.BtnClicked(input.KeyA) {\n            demonstrateFastFill(dmaManager)\n        }\n        \n        if input.BtnClicked(input.KeyB) {\n            demonstrateImageCopy(dmaManager)\n        }\n        \n        if input.BtnClicked(input.KeyStart) {\n            demonstrateAudioDMA(dmaManager)\n        }\n        \n        // Animate using DMA\n        animateWithDMA(dmaManager, frame)\n        \n        frame++\n    }\n}\n\nfunc setupGraphicsData(manager *dma.DMAManager) {\n    // Create a test palette\n    colors := []uint16{\n        0x0000, // Black\n        0x001F, // Red\n        0x03E0, // Green\n        0x7C00, // Blue\n        0x7FFF, // White\n        0x03FF, // Yellow\n        0x7FE0, // Cyan\n        0x7C1F, // Magenta\n    }\n    \n    // Load palette using DMA\n    err := manager.LoadPalette(\n        uintptr(unsafe.Pointer(&colors[0])),\n        palette.BG_PALETTE,\n        0, // palette index\n        len(colors),\n    )\n    if err != nil {\n        // Handle error\n    }\n    \n    // Create tile data\n    tileData := createTestTiles()\n    \n    // Load tiles using DMA\n    err = manager.LoadTiles(\n        uintptr(unsafe.Pointer(&tileData[0])),\n        0, // character block\n        0, // starting tile index\n        len(tileData)/32, // tile count\n    )\n    if err != nil {\n        // Handle error\n    }\n}\n\nfunc demonstrateFastFill(manager *dma.DMAManager) {\n    // Fill screen with solid color using DMA\n    screenSize := memory.SCREEN_WIDTH * memory.SCREEN_HEIGHT\n    \n    // Generate random color\n    color := uint16(rand.Intn(32768))\n    \n    err := manager.FillVRAM(0, color, uint32(screenSize))\n    if err != nil {\n        // Fallback to CPU fill\n        fallbackFill(color)\n    }\n}\n\nfunc demonstrateImageCopy(manager *dma.DMAManager) {\n    // Create a test image in memory\n    image := createTestImage()\n    \n    // Copy to VRAM using DMA for best performance\n    err := manager.CopyToVRAM(\n        uintptr(unsafe.Pointer(&image[0])),\n        0, // VRAM offset\n        uint32(len(image)*2), // size in bytes\n    )\n    \n    if err != nil {\n        // Fallback to CPU copy\n        fallbackCopy(image)\n    }\n}\n\nfunc demonstrateAudioDMA(manager *dma.DMAManager) {\n    // Set up audio streaming using DMA\n    audioBuffer := createAudioBuffer()\n    \n    // Configure DMA for audio streaming\n    err := manager.SetupAudioDMA(\n        dma.DMA1, // Use DMA1 for audio\n        uintptr(unsafe.Pointer(&audioBuffer[0])),\n        uint32(len(audioBuffer)),\n    )\n    \n    if err != nil {\n        return\n    }\n    \n    // Start audio streaming\n    manager.StartAudioStream(dma.DMA1)\n}\n\nfunc animateWithDMA(manager *dma.DMAManager, frame int) {\n    // Create animated pattern\n    pattern := createAnimatedPattern(frame)\n    \n    // Use asynchronous DMA to update screen section\n    x := 50\n    y := 50\n    width := 100\n    height := 60\n    \n    for row := 0; row < height; row++ {\n        srcOffset := row * width\n        dstOffset := ((y + row) * memory.SCREEN_WIDTH + x) * 2\n        \n        // Copy each row asynchronously\n        manager.CopyAsync(\n            uintptr(unsafe.Pointer(&pattern[srcOffset])),\n            memory.VRAM_BASE + uintptr(dstOffset),\n            uint32(width * 2), // width in bytes\n            nil, // no completion callback\n        )\n    }\n}\n\nfunc createTestTiles() []uint8 {\n    // Create sample tile data\n    tileCount := 16\n    tileData := make([]uint8, tileCount * 32) // 32 bytes per 4-bit tile\n    \n    for t := 0; t < tileCount; t++ {\n        for y := 0; y < 8; y++ {\n            for x := 0; x < 8; x++ {\n                color := uint8((t + x + y) % 16)\n                byteIndex := t*32 + (y*8 + x)/2\n                \n                if x%2 == 0 {\n                    tileData[byteIndex] = (tileData[byteIndex] & 0xF0) | color\n                } else {\n                    tileData[byteIndex] = (tileData[byteIndex] & 0x0F) | (color << 4)\n                }\n            }\n        }\n    }\n    \n    return tileData\n}\n\nfunc createTestImage() []uint16 {\n    image := make([]uint16, memory.SCREEN_WIDTH * memory.SCREEN_HEIGHT)\n    \n    for y := 0; y < memory.SCREEN_HEIGHT; y++ {\n        for x := 0; x < memory.SCREEN_WIDTH; x++ {\n            // Create gradient pattern\n            r := uint16((x * 31) / memory.SCREEN_WIDTH)\n            g := uint16((y * 31) / memory.SCREEN_HEIGHT)\n            b := uint16(((x + y) * 31) / (memory.SCREEN_WIDTH + memory.SCREEN_HEIGHT))\n            \n            color := (b << 10) | (g << 5) | r\n            image[y*memory.SCREEN_WIDTH + x] = color\n        }\n    }\n    \n    return image\n}\n\nfunc createAnimatedPattern(frame int) []uint16 {\n    width := 100\n    height := 60\n    pattern := make([]uint16, width * height)\n    \n    for y := 0; y < height; y++ {\n        for x := 0; x < width; x++ {\n            // Create moving wave pattern\n            wave := int(math.Sin(float64(x + frame)*0.2) * 10)\n            \n            var color uint16\n            if y >= 30 + wave && y <= 32 + wave {\n                color = 0x03FF // Yellow\n            } else {\n                color = 0x001F // Red\n            }\n            \n            pattern[y*width + x] = color\n        }\n    }\n    \n    return pattern\n}\n\nfunc createAudioBuffer() []uint32 {\n    // Create test audio data (simple sine wave)\n    sampleRate := 22050\n    duration := 1 // 1 second\n    bufferSize := sampleRate * duration\n    \n    buffer := make([]uint32, bufferSize)\n    \n    for i := 0; i < bufferSize; i++ {\n        // Generate 440Hz tone\n        sample := math.Sin(2 * math.Pi * 440 * float64(i) / float64(sampleRate))\n        sample16 := int16(sample * 32767)\n        \n        // Pack as stereo 16-bit samples\n        buffer[i] = uint32(sample16) | (uint32(sample16) << 16)\n    }\n    \n    return buffer\n}\n\nfunc fallbackFill(color uint16) {\n    // CPU-based fill as fallback\n    vramManager := vram.NewVRAMManager(memory.MODE_3)\n    buffer := vramManager.GetCurrentBuffer()\n    buffer.Clear(color)\n}\n\nfunc fallbackCopy(image []uint16) {\n    // CPU-based copy as fallback\n    vramManager := vram.NewVRAMManager(memory.MODE_3)\n    buffer := vramManager.GetCurrentBuffer()\n    \n    for y := 0; y < memory.SCREEN_HEIGHT; y++ {\n        for x := 0; x < memory.SCREEN_WIDTH; x++ {\n            buffer.PlotPixel(x, y, image[y*memory.SCREEN_WIDTH + x])\n        }\n    }\n}\n```

## Advanced Features

### DMA Chaining
```go\ntype DMAChain struct {\n    operations []*DMATransfer\n    current    int\n    loop       bool\n}\n\nfunc (dm *DMAManager) ExecuteChain(chain *DMAChain) error {\n    // Execute multiple DMA operations in sequence\n}\n```

### Memory Pool DMA
```go\nfunc (dm *DMAManager) CopyFromPool(poolId int, src, dst uintptr, size uint32) error {\n    // DMA with memory pool management\n}\n```

### Compressed Data DMA
```go\nfunc (dm *DMAManager) DecompressAndCopy(src uintptr, compression int, dst uintptr) error {\n    // DMA with hardware decompression\n}\n```

## File Structure
```
lib/dma/\n├── hardware.go     // Hardware register interface\n├── manager.go      // DMA coordination and management\n├── operations.go   // High-level DMA operations\n├── graphics.go     // Graphics-specific DMA utilities\n├── audio.go        // Audio DMA streaming\n└── utils.go        // Utilities and performance monitoring\n```

## Integration Points\n- Use memory constants from Task 01\n- Integrate with VRAM management from Task 03\n- Support sprite system from Task 04\n- Enhance tile system from Task 06\n- Accelerate background operations from Task 07\n- Enable audio system from Task 10\n\n## Resources\n- [GBATEK DMA Transfers](https://problemkaputt.de/gbatek.htm#gbadmatransfers)\n- [Tonc DMA](https://www.coranac.com/tonc/text/dma.htm)\n- [GBA DMA Programming Guide](https://www.cs.rit.edu/~tjh8300/CowBite/CowBiteSpec.htm#DMA%20Transfer%20Channels)\n\n## Success Criteria\n- All DMA channels work correctly with proper priority\n- Significant performance improvement over CPU operations\n- Audio streaming works without interruption\n- Graphics operations complete without visual artifacts\n- Error handling prevents system crashes\n- Integration with all graphics systems works seamlessly\n- Example program demonstrates all DMA capabilities\n- Comprehensive test coverage (>90%)\n- Performance benchmarks show expected speedup