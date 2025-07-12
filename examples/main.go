package main

import (
	"runtime/volatile"
	"unsafe"
	
	"github.com/matheusmortatti/gba-go/lib/input"
	"github.com/matheusmortatti/gba-go/lib/memory"
	"github.com/matheusmortatti/gba-go/lib/registers"
	"github.com/matheusmortatti/gba-go/lib/video"

	"math/rand"
)

func main() {
	// Set video mode 3 (16-bit bitmap) and enable BG2
	registers.Lcd.DISPCNT.Set(memory.MODE_3 | (1 << 10))

	// Test direct VRAM access first to make sure the screen works
	vram := (*[160][240]volatile.Register16)(unsafe.Pointer(uintptr(0x06000000)))
	
	// Fill screen with red using direct access to verify hardware works
	for y := 0; y < 160; y++ {
		for x := 0; x < 240; x++ {
			vram[y][x].Set(0x001F) // Red in 15-bit color
		}
	}

	// Now test our memory system - draw 4 colored squares
	memory.DrawRectMode3(0, 0, 120, 80, memory.GREEN)   // Top-left: Green
	memory.DrawRectMode3(120, 0, 120, 80, memory.BLUE)  // Top-right: Blue  
	memory.DrawRectMode3(0, 80, 120, 80, memory.YELLOW) // Bottom-left: Yellow
	memory.DrawRectMode3(120, 80, 120, 80, memory.WHITE) // Bottom-right: White

	for {
		video.VSync()
		input.Poll()
		
		// Draw random pixels when Down is pressed
		if input.BtnClicked(input.KeyDown) {
			randomColor := memory.RGB15(
				uint8(rand.Intn(32)),
				uint8(rand.Intn(32)), 
				uint8(rand.Intn(32)),
			)
			x := rand.Intn(memory.SCREEN_WIDTH)
			y := rand.Intn(memory.SCREEN_HEIGHT)
			memory.SetPixelMode3(x, y, randomColor)
		}
		
		// Clear screen and redraw squares when A is pressed
		if input.BtnClicked(input.KeyA) {
			memory.FillScreenMode3(memory.BLACK)
			memory.DrawRectMode3(0, 0, 120, 80, memory.GREEN)   
			memory.DrawRectMode3(120, 0, 120, 80, memory.BLUE)  
			memory.DrawRectMode3(0, 80, 120, 80, memory.YELLOW) 
			memory.DrawRectMode3(120, 80, 120, 80, memory.WHITE)
		}
	}
}

