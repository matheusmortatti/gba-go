package registers

import (
	"runtime/volatile"
	"unsafe"
)

type timer struct {
	TM0CNT_L *volatile.Register16 // Timer 0 Counter/Reload
	TM0CNT_H *volatile.Register16 // Timer 0 Control
	TM1CNT_L *volatile.Register16 // Timer 1 Counter/Reload
	TM1CNT_H *volatile.Register16 // Timer 1 Control
	TM2CNT_L *volatile.Register16 // Timer 2 Counter/Reload
	TM2CNT_H *volatile.Register16 // Timer 2 Control
	TM3CNT_L *volatile.Register16 // Timer 3 Counter/Reload
	TM3CNT_H *volatile.Register16 // Timer 3 Control
}

var Timer = &timer{
	TM0CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000100))),
	TM0CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000102))),
	TM1CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000104))),
	TM1CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000106))),
	TM2CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000108))),
	TM2CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400010A))),
	TM3CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400010C))),
	TM3CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400010E))),
}
