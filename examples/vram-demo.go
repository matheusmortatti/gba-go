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
	
	// Set up video mode - Mode 4 (8-bit bitmap) with BG2 enabled
	registers.Lcd.DISPCNT.Set(memory.MODE_4 | (1 << 10))
	
	// Set up palette
	paletteManager := palette.NewPaletteManager()
	rainbowPalette := createSimplePalette()
	paletteManager.LoadBGPalette256(rainbowPalette)
	
	// Create double buffer for smooth animation
	doubleBuffer := vram.NewDoubleBuffer(memory.MODE_4)
	if doubleBuffer == nil {
		// Fallback if double buffering not supported
		runSingleBufferDemo(vramManager)
		return
	}
	
	runDoubleBufferDemo(doubleBuffer)
}

func runDoubleBufferDemo(db *vram.DoubleBuffer) {
	frame := 0
	
	for {
		// Get back buffer for drawing
		backBuffer := db.GetBackBuffer()
		
		// Clear back buffer with black
		backBuffer.FastClear(0)
		
		// Draw animated scene
		drawAnimatedScene(backBuffer, frame)
		
		// Handle input for interactive elements
		input.Poll()
		if input.BtnDown(input.KeyA) {
			drawPlayerSprite(backBuffer, frame)
		}
		
		if input.BtnDown(input.KeyB) {
			// Draw test patterns
			drawTestPatterns(backBuffer, frame)
		}
		
		if input.BtnDown(input.KeySelect) {
			// Show performance information
			drawPerformanceInfo(backBuffer, frame)
		}
		
		// Present the buffer (swap and wait for VSync)
		db.Present()
		
		frame++
	}
}

func runSingleBufferDemo(vm *vram.VRAMManager) {
	// Switch to Mode 3 for single buffer demo
	vm.SetMode(memory.MODE_3)
	registers.Lcd.DISPCNT.Set(memory.MODE_3 | (1 << 10))
	
	buffer := vm.GetCurrentBuffer()
	frame := 0
	
	for {
		video.VSync()
		
		// Clear screen
		buffer.FastClear(0)
		
		// Draw animated content directly to screen
		drawAnimatedScene(buffer, frame)
		
		input.Poll()
		if input.BtnDown(input.KeyA) {
			drawPlayerSprite(buffer, frame)
		}
		
		frame++
	}
}

func drawAnimatedScene(buffer *vram.BitmapBuffer, frame int) {
	width := buffer.GetWidth()
	height := buffer.GetHeight()
	
	// Draw moving circles
	for i := 0; i < 5; i++ {
		x := (frame*2 + i*40) % width
		y := height/2 + int(30*sin(float64(frame+i*20)*0.1))
		color := uint8(i + 1) // Palette color index
		
		drawCircle(buffer, x, y, 15, uint16(color))
	}
	
	// Draw scrolling background pattern
	drawScrollingPattern(buffer, frame)
	
	// Draw bouncing rectangles
	for i := 0; i < 3; i++ {
		x := (frame*3 + i*60) % width
		y := int(20 + 15*sin(float64(frame+i*30)*0.15))
		color := uint8(10 + i)
		
		buffer.FillRect(x, y, 20, 15, uint16(color))
	}
}

func drawCircle(buffer *vram.BitmapBuffer, centerX, centerY, radius int, color uint16) {
	for y := centerY - radius; y <= centerY + radius; y++ {
		for x := centerX - radius; x <= centerX + radius; x++ {
			dx := x - centerX
			dy := y - centerY
			if dx*dx + dy*dy <= radius*radius {
				if buffer.InBounds(x, y) {
					buffer.PlotPixelFast(x, y, color)
				}
			}
		}
	}
}

func drawScrollingPattern(buffer *vram.BitmapBuffer, frame int) {
	width := buffer.GetWidth()
	height := buffer.GetHeight()
	
	for y := 0; y < height; y += 8 {
		for x := 0; x < width; x += 8 {
			patternX := (x + frame) % 32
			patternY := (y + frame/2) % 32
			color := uint8(((patternX + patternY) / 4) % 16)
			buffer.FillRect(x, y, 8, 8, uint16(color))
		}
	}
}

func drawPlayerSprite(buffer *vram.BitmapBuffer, frame int) {
	width := buffer.GetWidth()
	height := buffer.GetHeight()
	
	// Simple animated player sprite
	spriteX := width / 2
	spriteY := height / 2
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
	
	// Animate the sprite by changing its vertical position
	animOffset := int(5 * sin(float64(frame)*0.2))
	
	for y := 0; y < 8; y++ {
		for x := 0; x < 8; x++ {
			if spriteData[y][x] > 0 {
				color := spriteColor + spriteData[y][x] - 1
				pixelX := spriteX + x - 4
				pixelY := spriteY + y - 4 + animOffset
				if buffer.InBounds(pixelX, pixelY) {
					buffer.PlotPixelFast(pixelX, pixelY, uint16(color))
				}
			}
		}
	}
}

func drawTestPatterns(buffer *vram.BitmapBuffer, frame int) {
	// Cycle through different test patterns
	patternIndex := (frame / 60) % 8
	
	switch patternIndex {
	case 0:
		vram.FillPattern(buffer, vram.PATTERN_CHECKERBOARD, 1, 2)
	case 1:
		vram.FillPattern(buffer, vram.PATTERN_GRADIENT_H, 0, 0)
	case 2:
		vram.FillPattern(buffer, vram.PATTERN_GRADIENT_V, 0, 0)
	case 3:
		vram.FillPattern(buffer, vram.PATTERN_STRIPES_H, 3, 4)
	case 4:
		vram.FillPattern(buffer, vram.PATTERN_STRIPES_V, 5, 6)
	case 5:
		vram.FillPattern(buffer, vram.PATTERN_RAINBOW, 0, 0)
	case 6:
		vram.FillPattern(buffer, vram.PATTERN_NOISE, 0, 0)
	default:
		vram.FillPattern(buffer, vram.PATTERN_SOLID, 7, 0)
	}
}

func drawPerformanceInfo(buffer *vram.BitmapBuffer, frame int) {
	// Draw some simple performance indicators
	
	// Frame counter display (simple digit representation)
	frameDigits := []int{
		(frame / 1000) % 10,
		(frame / 100) % 10,
		(frame / 10) % 10,
		frame % 10,
	}
	
	for i, digit := range frameDigits {
		drawDigit(buffer, 10+i*12, 10, digit, 15)
	}
	
	// Draw memory usage indicator
	usage := vram.GetMemoryUsage(memory.MODE_4)
	totalSize := vram.VRAM_SIZE
	usagePercent := (usage.VRAMUsed * 100) / totalSize
	
	// Draw usage bar
	barWidth := (usagePercent * 100) / 100
	buffer.FillRect(10, 30, barWidth, 8, 14) // Green for used
	buffer.FillRect(10+barWidth, 30, 100-barWidth, 8, 1) // Red for free
}

func drawDigit(buffer *vram.BitmapBuffer, x, y, digit int, color uint16) {
	// Simple 7-segment style digit patterns (5x7 pixels)
	patterns := [][]uint8{
		// 0
		{1,1,1,1,1,
		 1,0,0,0,1,
		 1,0,0,0,1,
		 1,0,0,0,1,
		 1,0,0,0,1,
		 1,0,0,0,1,
		 1,1,1,1,1},
		// 1
		{0,0,1,0,0,
		 0,1,1,0,0,
		 0,0,1,0,0,
		 0,0,1,0,0,
		 0,0,1,0,0,
		 0,0,1,0,0,
		 1,1,1,1,1},
		// 2
		{1,1,1,1,1,
		 0,0,0,0,1,
		 0,0,0,0,1,
		 1,1,1,1,1,
		 1,0,0,0,0,
		 1,0,0,0,0,
		 1,1,1,1,1},
		// 3
		{1,1,1,1,1,
		 0,0,0,0,1,
		 0,0,0,0,1,
		 1,1,1,1,1,
		 0,0,0,0,1,
		 0,0,0,0,1,
		 1,1,1,1,1},
		// 4
		{1,0,0,0,1,
		 1,0,0,0,1,
		 1,0,0,0,1,
		 1,1,1,1,1,
		 0,0,0,0,1,
		 0,0,0,0,1,
		 0,0,0,0,1},
		// 5
		{1,1,1,1,1,
		 1,0,0,0,0,
		 1,0,0,0,0,
		 1,1,1,1,1,
		 0,0,0,0,1,
		 0,0,0,0,1,
		 1,1,1,1,1},
		// 6
		{1,1,1,1,1,
		 1,0,0,0,0,
		 1,0,0,0,0,
		 1,1,1,1,1,
		 1,0,0,0,1,
		 1,0,0,0,1,
		 1,1,1,1,1},
		// 7
		{1,1,1,1,1,
		 0,0,0,0,1,
		 0,0,0,0,1,
		 0,0,0,0,1,
		 0,0,0,0,1,
		 0,0,0,0,1,
		 0,0,0,0,1},
		// 8
		{1,1,1,1,1,
		 1,0,0,0,1,
		 1,0,0,0,1,
		 1,1,1,1,1,
		 1,0,0,0,1,
		 1,0,0,0,1,
		 1,1,1,1,1},
		// 9
		{1,1,1,1,1,
		 1,0,0,0,1,
		 1,0,0,0,1,
		 1,1,1,1,1,
		 0,0,0,0,1,
		 0,0,0,0,1,
		 1,1,1,1,1},
	}
	
	if digit < 0 || digit > 9 {
		return
	}
	
	pattern := patterns[digit]
	for py := 0; py < 7; py++ {
		for px := 0; px < 5; px++ {
			if pattern[py*5+px] == 1 {
				if buffer.InBounds(x+px, y+py) {
					buffer.PlotPixelFast(x+px, y+py, color)
				}
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

// Simple sine approximation for animation
func sin(x float64) float64 {
	// Simple sine approximation using Taylor series (first few terms)
	// This is a very basic implementation for demo purposes
	for x > 6.28318530718 {
		x -= 6.28318530718
	}
	for x < 0 {
		x += 6.28318530718
	}
	
	x2 := x * x
	x3 := x2 * x
	x5 := x3 * x2
	
	return x - x3/6.0 + x5/120.0
}