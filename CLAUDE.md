# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Game Boy Advance (GBA) development library built with TinyGo, providing Go abstractions for GBA hardware functionality. The project enables writing GBA games in Go that compile to GBA ROM files.

## Build and Development Commands

### Building
- `tinygo build -target gameboy-advance -o main.gba examples/main.go` - Build GBA ROM from main example
- `tinygo build -target gameboy-advance -o <output>.gba <source>.go` - Build any Go file to GBA ROM

### Testing and Running
- `tinygo test -target gameboy-advance .` - Run tests for GBA target
- `tinygo flash -target gameboy-advance examples/main.go` - Flash to real GBA hardware (requires flashcart)
- `mgba <.gba file path>` - Run the built gba ROM to check if behaviour is correct

### Development
- `go mod tidy` - Update dependencies
- `go fmt ./...` - Format all Go code
- `go vet ./...` - Run Go static analysis

## Architecture

### Core Library Structure (`lib/`)

- **gba.go** - Main package entry point
- **registers/** - Hardware register abstractions for GBA components:
  - LCD control and status registers
  - Keypad input registers  
  - Interrupt management registers
  - DMA transfer channel registers
  - Timer, sound, and serial communication registers
- **input/** - Input handling with button state tracking and interrupt-based polling
- **video/** - Video synchronization utilities (VSync implementation)
- **interrupts/** - Interrupt management system
- **bios/** - BIOS function wrappers
- **drawing/** - Drawing utilities (extends tinydraw)

### Key Concepts

- **Memory-mapped I/O**: All hardware interaction uses volatile register access at specific memory addresses
- **Hardware Registers**: Low-level access to GBA hardware through memory-mapped registers
- **VSync**: Video synchronization for smooth animation (video.VSync())
- **Input Polling**: Frame-based input state management with click detection
- **TinyGo Target**: Uses `gameboy-advance` target for cross-compilation

### Dependencies

- **TinyGo**: Required for GBA cross-compilation
- **tinygo.org/x/tinydraw**: Drawing primitives and graphics utilities
- **machine**: Hardware abstraction layer for display and peripherals

### Example Structure

The `examples/` directory contains sample GBA programs demonstrating library usage. Built ROMs output to `examples/bin/` directory.

## Key Development Notes

- Always use `tinygo` instead of standard `go` for building GBA targets
- Memory access uses `runtime/volatile` and `unsafe` packages for hardware registers
- Input state is polled each frame - call `input.Poll()` in main loop
- Use `video.VSync()` to synchronize with display refresh rate
- VRAM is accessed directly at memory address 0x06000000 for pixel manipulation

## Critical GBA Programming Requirements

### Display Initialization (CRITICAL for visible output)
When creating GBA programs, the display MUST be properly initialized or you will get a blank white screen:

```go
// CORRECT: Set video mode and enable background in one operation
registers.Lcd.DISPCNT.Set(memory.MODE_3 | (1 << 10)) // Mode 3 + BG2 enabled

// INCORRECT: Using separate SetBits/ClearBits calls may not work reliably
registers.Lcd.DISPCNT.SetBits(memory.MODE_3)  // Don't do this
registers.Lcd.DISPCNT.SetBits(1 << 10)        // Don't do this
```

### Double Buffering Display Control
For double buffering in Mode 4/5, update the entire DISPCNT register:

```go
// CORRECT: Update entire register with mode, BG2, and frame selection
var displayValue uint16
if currentPage == 1 {
    displayValue = uint16(mode) | (1 << 10) | (1 << 4) // Mode + BG2 + Frame1
} else {
    displayValue = uint16(mode) | (1 << 10)             // Mode + BG2 + Frame0
}
registers.Lcd.DISPCNT.Set(displayValue)

// INCORRECT: Partial register updates may not work
registers.Lcd.DISPCNT.SetBits(1 << 4)   // Don't do this
registers.Lcd.DISPCNT.ClearBits(1 << 4) // Don't do this
```

### Volatile Register Usage
Always use the correct volatile register API:

```go
// CORRECT: Use Get() and Set() methods
(*volatile.Register16)(unsafe.Pointer(addr)).Set(value)
value := (*volatile.Register16)(unsafe.Pointer(addr)).Get()

// INCORRECT: Direct assignment with volatile.Register16()
*(*volatile.Register16)(unsafe.Pointer(addr)) = volatile.Register16(value) // Don't do this
```

### Testing Programs
Always test built programs with: `mgba examples/bin/program.gba` to verify visual output is working correctly.