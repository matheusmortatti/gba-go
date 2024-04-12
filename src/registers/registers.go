package registers

import (
	"runtime/volatile"
	"unsafe"
)

type Interrupt {
	IE *volatile.Register16,
	IF *volatile.Register16,
	IFBios *volatile.Register16
}

func NewInterruptRegister() *Interrupt {
	return &Interrupt{
		IE: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000200))),
		IF: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000202))),
		IFBios: (*volatile.Register16)(unsafe.Pointer(uintptr(0x03007FF8))),
	}
}