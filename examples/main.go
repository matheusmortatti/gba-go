package main

import (
	. "github.com/matheusmortatti/gba-go"

	"machine"
	"tinygo.org/x/tinydraw"
)

var (
	display    = machine.Display
	interrupts = NewInterrupts()
)

func main() {
	display.Configure()
	interrupts.EnableVBlankInterrupt(func() {})
}
