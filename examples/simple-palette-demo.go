package main

import (
	"github.com/matheusmortatti/gba-go/lib/palette"
	"github.com/matheusmortatti/gba-go/lib/registers"
	"github.com/matheusmortatti/gba-go/lib/video"
	"github.com/matheusmortatti/gba-go/lib/input"
	"runtime/volatile"
	"unsafe"
)

const (
	MODE_3 = 3
	BG2_ENABLE = 1 << 10
	VRAM_BASE = 0x06000000
	SCREEN_WIDTH = 240
	SCREEN_HEIGHT = 160
)

func main() {
	// Initialize palette manager
	paletteManager := palette.NewPaletteManager()
	
	// Set video mode 3 for direct bitmap mode (15-bit colors)
	registers.Lcd.DISPCNT.Set(MODE_3 | BG2_ENABLE)
	
	// Create different themed palettes
	rainbowPalette := palette.CreateRainbowPalette()
	firePalette := palette.CreateFirePalette()
	waterPalette := palette.CreateWaterPalette()
	earthPalette := palette.CreateEarthPalette()
	grayPalette := palette.CreateGrayscalePalette()
	
	// Load initial palette (this loads to hardware palette RAM)
	paletteManager.LoadBGPalette16(0, rainbowPalette)
	
	// Get VRAM pointer for direct pixel access
	vram := (*[SCREEN_WIDTH * SCREEN_HEIGHT]volatile.Register16)(unsafe.Pointer(uintptr(VRAM_BASE)))
	
	currentPaletteIndex := 0
	palettes := []*palette.Palette16{rainbowPalette, firePalette, waterPalette, earthPalette, grayPalette}
	
	colorIndex := 1 // Start from 1 (skip transparent color)
	frameCounter := 0
	
	for {
		video.VSync()
		input.Poll()
		
		// Switch palettes with A button
		if input.BtnClicked(input.KeyA) {
			currentPaletteIndex = (currentPaletteIndex + 1) % len(palettes)
			paletteManager.LoadBGPalette16(0, palettes[currentPaletteIndex])
		}
		
		// Rotate colors with B button
		if input.BtnClicked(input.KeyB) {
			paletteManager.RotatePalette(0, 1, 15) // Rotate colors 1-15
		}
		
		// Auto-cycle through colors every 60 frames
		if frameCounter%60 == 0 {
			colorIndex = (colorIndex % 15) + 1 // Keep in range 1-15
		}
		
		// In Mode 3, we need to write the actual RGB15 color values to VRAM
		// Since Mode 3 doesn't use palette indices, we get the actual color value
		currentColor := palettes[currentPaletteIndex].GetColor(colorIndex)
		
		// Simple demo: fill screen with current color
		fillScreen(vram, currentColor)
		
		frameCounter++
	}
}

// fillScreen fills the entire screen with a single color
func fillScreen(vram *[SCREEN_WIDTH * SCREEN_HEIGHT]volatile.Register16, color palette.Color) {
	colorValue := uint16(color)
	for i := 0; i < SCREEN_WIDTH*SCREEN_HEIGHT; i++ {
		vram[i].Set(colorValue)
	}
}

// More complex demo showing multiple colors
func drawColorBars(vram *[SCREEN_WIDTH * SCREEN_HEIGHT]volatile.Register16, currentPalette *palette.Palette16) {
	barWidth := SCREEN_WIDTH / 16
	
	for i := 0; i < 16; i++ {
		color := currentPalette.GetColor(i)
		colorValue := uint16(color)
		startX := i * barWidth
		endX := startX + barWidth
		
		for y := 0; y < SCREEN_HEIGHT; y++ {
			for x := startX; x < endX && x < SCREEN_WIDTH; x++ {
				offset := y*SCREEN_WIDTH + x
				vram[offset].Set(colorValue)
			}
		}
	}
}