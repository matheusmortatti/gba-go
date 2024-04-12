package interrupts

import (
	"machine"
	"runtime/interrupt"

	"github.com/matheusmortatti/go/gba/registers"
)

type Interrupts struct {
	handlers map[interrupt.Interrupt]func()
}

func NewInterrupts() *Interrupts {
	i := &Interrupts{
		handlers: make(map[interrupt.Interrupt]func()),
	}
	return i
}

func (i *Interrupts) EnableVBlankInterrupt(handler func()) {
	registers.Video.DispStat.SetBits(1 << 3)
	itr := interrupt.New(machine.IRQ_VBLANK, i.handleInterrupt)
	i.enableInterrupt(itr, handler)
}

func DisableAllInterrupts() {
	interrupt.Disable()
}

func (i *Interrupts) handleInterrupt(itr interrupt.Interrupt) {
	handler, ok := i.handlers[itr]
	if ok {
		handler()
	}
}

func (i *Interrupts) enableInterrupt(itr interrupt.Interrupt, handler func()) {
	i.handlers[itr] = handler
	itr.Enable()
}
