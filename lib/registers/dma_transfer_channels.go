package registers

import (
	"runtime/volatile"
	"unsafe"
)

type dmaTransferChannels struct {
	DMA0SAD   *volatile.Register32 // DMA 0 Source Address
	DMA0DAD   *volatile.Register32 // DMA 0 Destination Address
	DMA0CNT_L *volatile.Register16 // Word Count
	DMA0CNT_H *volatile.Register16 // DMA 0 Control
	DMA1SAD   *volatile.Register32 // DMA 1 Source Address
	DMA1DAD   *volatile.Register32 // DMA 1 Destination Address
	DMA1CNT_L *volatile.Register16 // Word Count
	DMA1CNT_H *volatile.Register16 // DMA 1 Control
	DMA2SAD   *volatile.Register32 // DMA 2 Source Address
	DMA2DAD   *volatile.Register32 // DMA 2 Destination Address
	DMA2CNT_L *volatile.Register16 // Word Count
	DMA2CNT_H *volatile.Register16 // DMA 2 Control
	DMA3SAD   *volatile.Register32 // DMA 3 Source Address
	DMA3DAD   *volatile.Register32 // DMA 3 Destination Address
	DMA3CNT_L *volatile.Register16 // Word Count
	DMA3CNT_H *volatile.Register16 // DMA 3 Control
}

var DmaTransferChannels = &dmaTransferChannels{
	DMA0SAD:   (*volatile.Register32)(unsafe.Pointer(uintptr(0x040000B0))),
	DMA0DAD:   (*volatile.Register32)(unsafe.Pointer(uintptr(0x040000B4))),
	DMA0CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x040000B8))),
	DMA0CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x040000BA))),
	DMA1SAD:   (*volatile.Register32)(unsafe.Pointer(uintptr(0x040000BC))),
	DMA1DAD:   (*volatile.Register32)(unsafe.Pointer(uintptr(0x040000C0))),
	DMA1CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x040000C4))),
	DMA1CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x040000C6))),
	DMA2SAD:   (*volatile.Register32)(unsafe.Pointer(uintptr(0x040000C8))),
	DMA2DAD:   (*volatile.Register32)(unsafe.Pointer(uintptr(0x040000CC))),
	DMA2CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x040000D0))),
	DMA2CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x040000D2))),
	DMA3SAD:   (*volatile.Register32)(unsafe.Pointer(uintptr(0x040000D4))),
	DMA3DAD:   (*volatile.Register32)(unsafe.Pointer(uintptr(0x040000D8))),
	DMA3CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x040000DC))),
	DMA3CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x040000DE))),
}
