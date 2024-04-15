package registers

import (
	"runtime/volatile"
	"unsafe"
)

type keypad struct {
	KEYINPUT *volatile.Register16 // Key Input
	KEYCNT   *volatile.Register16 // Key Control
}

var Keypad = &keypad{
	KEYINPUT: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000130))),
	KEYCNT:   (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000132))),
}
