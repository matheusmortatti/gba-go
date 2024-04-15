package registers

mport (
	"runtime/volatile"
	"unsafe"
)

type sound struct {
	SOUND1CNT_L *volatile.Register16, // Channel 1 Sweep register       (NR10)
	SOUND1CNT_H *volatile.Register16, // Channel 1 Duty/Length/Envelope (NR11, NR12)
	SOUND1CNT_X *volatile.Register16, // Channel 1 Frequency/Control    (NR13, NR14)
	SOUND2CNT_L *volatile.Register16, // Channel 2 Duty/Length/Envelope (NR21, NR22)
	SOUND2CNT_H *volatile.Register16, // Channel 2 Frequency/Control    (NR23, NR24)
	SOUND3CNT_L *volatile.Register16, // Channel 3 Stop/Wave RAM select (NR30)
	SOUND3CNT_H *volatile.Register16, // Channel 3 Length/Volume        (NR31, NR32)
	SOUND3CNT_X *volatile.Register16, // Channel 3 Frequency/Control    (NR33, NR34)
	SOUND4CNT_L *volatile.Register16, // Channel 4 Length/Envelope      (NR41, NR42)
	SOUND4CNT_H *volatile.Register16, // Channel 4 Frequency/Control    (NR43, NR44)
	SOUNDCNT_L *volatile.Register16, // Control Stereo/Volume/Enable   (NR50)
	SOUNDCNT_H *volatile.Register16, // Control Mixing/DMA Control     (NR51)
	SOUNDCNT_X *volatile.Register16, // Control Sound on/off           (NR52)
	SOUNDBIAS *volatile.Register16, // Sound PWM Control              (NR53)
	WAVE_RAM *volatile.Register64, // Wave Pattern RAM               (NR30-3F)
	FIFO_A *volatile.Register32, // Channel A FIFO                 (NR30-3F)
	FIFO_B *volatile.Register32, // Channel B FIFO                 (NR30-3F)
}

var Sound = &sound{
	SOUND1CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000060))),
	SOUND1CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000062))),
	SOUND1CNT_X: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000064))),
	SOUND2CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000068))),
	SOUND2CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400006C))),
	SOUND3CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000070))),
	SOUND3CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000072))),
	SOUND3CNT_X: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000074))),
	SOUND4CNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000078))),
	SOUND4CNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400007C))),
	SOUNDCNT_L: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000080))),
	SOUNDCNT_H: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000082))),
	SOUNDCNT_X: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000084))),
	SOUNDBIAS: (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000088))),
	WAVE_RAM: (*volatile.Register64)(unsafe.Pointer(uintptr(0x04000090))),
	FIFO_A: (*volatile.Register32)(unsafe.Pointer(uintptr(0x040000A0))),
	FIFO_B: (*volatile.Register32)(unsafe.Pointer(uintptr(0x040000A4))),
}