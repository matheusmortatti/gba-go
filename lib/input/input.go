package input

import (
	"github.com/matheusmortatti/gba-go/lib/interrupts"
	"github.com/matheusmortatti/gba-go/lib/registers"
)

const (
	KeyA = 1 << iota
	KeyB
	KeySelect
	KeyStart
	KeyRight
	KeyLeft
	KeyUp
	KeyDown
	KeyR
	KeyL
)

var (
	lastState    uint16 = 0x3FF
	currentState uint16 = 0x3FF
)

// WasBtnDown returns true if the key was down in the last frame.
func WasBtnDown(key uint16) bool {
	return lastState&key != 0
}

// BtnDown returns true if the key is currently down.
func BtnDown(key uint16) bool {
	return currentState&key == 0
}

// BtnUp returns true if the key is currently up.
func BtnUp(key uint16) bool {
	return !BtnDown(key)
}

// BtnClicked returns true if the key was pressed in the current frame.
func BtnClicked(key uint16) bool {
	return BtnDown(key) && !WasBtnDown(key)
}

// Poll updates the current and last key states.
func Poll() {
	lastState = currentState
	currentState = registers.Keypad.KEYINPUT.Get()
}

// EnablePolling enables the keypad polling interrupt.
func EnablePolling() {
	registers.Keypad.KEYCNT.SetBits(1 << 0xE)
	registers.Keypad.KEYCNT.SetBits(0b1111111111)
	interrupts.EnableKeypadPollingInterrupt(keyInterruptHandler)
}

func keyInterruptHandler() {
	Poll()
}
