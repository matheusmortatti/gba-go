package main

import (
	"github.com/matheusmortatti/gba-go/lib/interrupts"

	"machine"
	"tinygo.org/x/tinydraw"
)

var (
	display    = machine.Display
	interrupts = interrupts.NewInterrupts()
)

func main() {
	display.Configure()
	interrupts.EnableVBlankInterrupt(func() {})
}
