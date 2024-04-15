package registers

import (
	"runtime/volatile"
	"unsafe"
)

type lcd struct {
	DISPCNT *volatile.Register16, // LCD Control
	DISPSTAT *volatile.Register16, // General LCD Status (STAT,LYC)
	VCOUNT *volatile.Register16, // Vertical Counter (LY)
	BG0CNT *volatile.Register16, // BG0 Control
	BG1CNT *volatile.Register16, // BG1 Control
	BG2CNT *volatile.Register16, // BG2 Control
	BG3CNT *volatile.Register16, // BG3 Control
	BG0HOFS *volatile.Register16, // BG0 X-Offset
	BG0VOFS *volatile.Register16, // BG0 Y-Offset
	BG1HOFS *volatile.Register16, // BG1 X-Offset
	BG1VOFS *volatile.Register16, // BG1 Y-Offset
	BG2HOFS *volatile.Register16, // BG2 X-Offset
	BG2VOFS *volatile.Register16, // BG2 Y-Offset
	BG3HOFS *volatile.Register16, // BG3 X-Offset
	BG3VOFS *volatile.Register16, // BG3 Y-Offset
	BG2PA *volatile.Register16, // BG2 Rotation/Scaling Parameter A (dx)
	BG2PB *volatile.Register16, // BG2 Rotation/Scaling Parameter B (dmx)
	BG2PC *volatile.Register16, // BG2 Rotation/Scaling Parameter C (dy)
	BG2PD *volatile.Register16, // BG2 Rotation/Scaling Parameter D (dmy)
	BG2X *volatile.Register32, // BG2 Reference Point X-Coordinate
	BG2Y *volatile.Register32, // BG2 Reference Point Y-Coordinate
	BG3PA *volatile.Register16, // BG3 Rotation/Scaling Parameter A (dx)
	BG3PB *volatile.Register16, // BG3 Rotation/Scaling Parameter B (dmx)
	BG3PC *volatile.Register16, // BG3 Rotation/Scaling Parameter C (dy)
	BG3PD *volatile.Register16, // BG3 Rotation/Scaling Parameter D (dmy)
	BG3X *volatile.Register32, // BG3 Reference Point X-Coordinate
	BG3Y *volatile.Register32, // BG3 Reference Point Y-Coordinate
	WIN0H *volatile.Register16, // Window 0 Horizontal Dimensions
	WIN1H *volatile.Register16, // Window 1 Horizontal Dimensions
	WIN0V *volatile.Register16, // Window 0 Vertical Dimensions
	WIN1V *volatile.Register16, // Window 1 Vertical Dimensions
	WININ *volatile.Register16, // Inside of Window 0 and 1
	WINOUT *volatile.Register16, // Inside of OBJ Window & Outside of Windows
	MOSAIC *volatile.Register16, // Mosaic Size
	BLDCNT *volatile.Register16, // Color Special Effects Selection
	BLDALPHA *volatile.Register16, // Alpha Blending Coefficients
	BLDY *volatile.Register16 // Brightness (Fade-In/Out) Coefficient
}

var Lcd = &lcd{
	DISPCNT: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000000))),
	DISPSTAT: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000004))),
	VCOUNT: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000006))),
	BG0CNT: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000008))),
	BG1CNT: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400000A))),
	BG2CNT: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400000C))),
	BG3CNT: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400000E))),
	BG0HOFS: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000010))),
	BG0VOFS: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000012))),
	BG1HOFS: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000014))),
	BG1VOFS: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000016))),
	BG2HOFS: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000018))),
	BG2VOFS: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400001A))),
	BG3HOFS: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400001C))),
	BG3VOFS: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400001E))),
	BG2PA: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000020))),
	BG2PB: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000022))),
	BG2PC: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000024))),
	BG2PD: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000026))),
	BG2X: (*volatile.Register32)(unsafe.Pointer(uintptr(0x04000028))),
	BG2Y: (*volatile.Register32)(unsafe.Pointer(uintptr(0x0400002C))),
	BG3PA: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000030))),
	BG3PB: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000032))),
	BG3PC: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000034))),
	BG3PD: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000036))),
	BG3X: (*volatile.Register32)(unsafe.Pointer(uintptr(0x04000038))),
	BG3Y: (*volatile.Register32)(unsafe.Pointer(uintptr(0x0400003C))),
	WIN0H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000040))),
	WIN1H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000042))),
	WIN0V: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000044))),
	WIN1V: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000046))),
	WININ: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000048))),
	WINOUT: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400004A))),
	MOSAIC: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400004C))),
	BLDCNT: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000050))),
	BLDALPHA: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000052))),
	BLDY: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000054))),
}