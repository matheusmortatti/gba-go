package bios

import (
	"device/arm"
)

func VBlankIntrWait() {
	arm.Asm("swi 0x50000" /* Instr_VBlankIntrWait */)
}
