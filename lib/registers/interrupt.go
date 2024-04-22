package registers

import (
	"runtime/volatile"
	"unsafe"
)

type interrupt struct {
	IE     *volatile.Register16 // Interrupt Enable
	IF     *volatile.Register16 // Interrupt Request Flags
	IME    *volatile.Register16 // Interrupt Master Enable
	IFBios *volatile.Register16 // Interrupt Request Flags (BIOS)
}

var Interrupt = &interrupt{
	IE:     (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000200))),
	IF:     (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000202))),
	IME:    (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000208))),
	IFBios: (*volatile.Register16)(unsafe.Pointer(uintptr(0x03007FF8))),
}
