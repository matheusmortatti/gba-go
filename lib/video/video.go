package video

import (
	"github.com/matheusmortatti/gba-go/lib/registers"
)

func VSync() {
	for registers.Lcd.VCOUNT.Get() >= 160 {
	}
	for registers.Lcd.VCOUNT.Get() < 160 {
	}
}
