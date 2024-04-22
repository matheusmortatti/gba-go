package interrupts

import (
	"machine"
	"runtime/interrupt"

	"github.com/matheusmortatti/gba-go/lib/registers"
)

var handlers = make(map[interrupt.Interrupt]func())

func EnableVBlankInterrupt(handler func()) {
	registers.Lcd.DISPSTAT.Set(1<<3 | 1<<4 | 1<<0xA)
	itr := interrupt.New(machine.IRQ_VBLANK, handleInterrupt)
	enableInterrupt(itr, handler)
}

func EnableKeypadPollingInterrupt(handler func()) {
	itr := interrupt.New(machine.IRQ_KEYPAD, handleInterrupt)
	enableInterrupt(itr, handler)
}

func DisableAllInterrupts() {
	interrupt.Disable()
}

func handleInterrupt(itr interrupt.Interrupt) {
	handler, ok := handlers[itr]
	if ok {
		handler()
	}
}

func enableInterrupt(itr interrupt.Interrupt, handler func()) {
	handlers[itr] = handler
	itr.Enable()
}
