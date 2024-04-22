package main

import (
	"github.com/matheusmortatti/gba-go/lib/drawing"
	"github.com/matheusmortatti/gba-go/lib/input"
	"github.com/matheusmortatti/gba-go/lib/interrupts"
	"github.com/matheusmortatti/gba-go/lib/registers"

	"image/color"
	"machine"
	"tinygo.org/x/tinydraw"
)

var (
	display = machine.Display
	// Screen resolution
	screenWidth, screenHeight = display.Size()
	green                     = color.RGBA{0, 255, 0, 255}
	black                     = color.RGBA{}
)

func main() {
	display.Configure()
	interrupts.EnableVBlankInterrupt(func() {
		registers.Interrupt.IFBios.SetBits(1)
		update()
		drawing.VSync()
		drawing.Display()
	})
	input.EnablePolling()

	for {
	}
}

func update() {
	clearScreen()
	if input.BtnDown(input.KeyA) {
		tinydraw.FilledRectangle(
			&display,
			int16(0), int16(0),
			screenWidth, screenHeight,
			green,
		)
	}
}

func clearScreen() {
	tinydraw.FilledRectangle(
		&display,
		int16(0), int16(0),
		screenWidth, screenHeight,
		black,
	)
}
