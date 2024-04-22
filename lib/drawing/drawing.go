package drawing

import (
	"github.com/matheusmortatti/gba-go/lib/bios"
	"github.com/matheusmortatti/gba-go/lib/registers"
)

func VCount() uint16 {
	return registers.Lcd.VCOUNT.Get()
}

func VSync() {
	bios.VBlankIntrWait()
}

var drawPage = 1

func Display() error {
	old := registers.Lcd.DISPCNT.Get()
	registers.Lcd.DISPCNT.Set(old ^ (uint16(drawPage) << 4)) // flip display
	drawPage ^= 1                                            // switch drawPage
	return nil
}
