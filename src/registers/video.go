package registers

import (
	"runtime/volatile"
	"unsafe"
)

type Video struct {
	DispCnt *volatile.Register16,
	VCount *volatile.Register16,
	BG0Cnt *volatile.Register16,
	BG1Cnt *volatile.Register16,
	BG2Cnt *volatile.Register16,
	BG3Cnt *volatile.Register16,
	DispStat *volatile.Register16
}

func NewVideoRegister() *Video {
	return &Video{
		DispCnt: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000000))),
		VCount: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000006))),
		BG0Cnt: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000008))),
		BG1Cnt: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400000A))),
		BG2Cnt: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400000C))),
		BG3Cnt: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400000E))),
		DispStat: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000004))),
	}
}