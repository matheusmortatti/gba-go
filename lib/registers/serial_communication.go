package registers

import (
	"runtime/volatile"
	"unsafe"
)

type serialCommunication struct {
	SIODATA32   *volatile.Register32 // Data
	SIOMULTI0   *volatile.Register16 // Data 0
	SIOMULTI1   *volatile.Register16 // Data 1
	SIOMULTI2   *volatile.Register16 // Data 2
	SIOMULTI3   *volatile.Register16 // Data 3
	SIOCNT      *volatile.Register16 // Control
	SIOMLT_SEND *volatile.Register16 // Data Send
	SIODATA8    *volatile.Register16 // Data

	RCNT      *volatile.Register16 // SIO Mode Select/General Purpose Data
	IR        *volatile.Register16 // Ancient - Infrared Register (Prototypes only)
	JOYCNT    *volatile.Register16 // SIO JOY Bus Control
	JOY_RECV  *volatile.Register32 // SIO JOY Bus Receive Data
	JOY_TRANS *volatile.Register32 // SIO JOY Bus Transmit Data
	JOYSTAT   *volatile.Register16 // SIO JOY Bus Receive Status
}

var SerialCommunication = &serialCommunication{
	SIODATA32:   (*volatile.Register32)(unsafe.Pointer(uintptr(0x04000120))),
	SIOMULTI0:   (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000120))),
	SIOMULTI1:   (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000122))),
	SIOMULTI2:   (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000124))),
	SIOMULTI3:   (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000126))),
	SIOCNT:      (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000128))),
	SIOMLT_SEND: (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400012A))),
	SIODATA8:    (*volatile.Register16)(unsafe.Pointer(uintptr(0x0400012A))),

	RCNT:      (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000134))),
	IR:        (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000136))),
	JOYCNT:    (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000140))),
	JOY_RECV:  (*volatile.Register32)(unsafe.Pointer(uintptr(0x04000150))),
	JOY_TRANS: (*volatile.Register32)(unsafe.Pointer(uintptr(0x04000154))),
	JOYSTAT:   (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000158))),
}
