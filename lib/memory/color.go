package memory

// Color represents a 15-bit GBA color
type Color uint16

// RGB15 constructs a Color from 5-bit RGB components (0-31 each)
func RGB15(r, g, b uint8) Color {
	// Clamp values to 5 bits (0-31)
	if r > 31 {
		r = 31
	}
	if g > 31 {
		g = 31
	}
	if b > 31 {
		b = 31
	}
	
	// GBA color format: 0bbbbbgggggrrrrr
	return Color(uint16(r) | (uint16(g) << 5) | (uint16(b) << 10))
}

// RGB constructs a Color from 8-bit RGB components (0-255 each)
// Automatically converts to 5-bit GBA format
func RGB(r, g, b uint8) Color {
	return RGB15(r>>3, g>>3, b>>3)
}

// R extracts the red component (0-31)
func (c Color) R() uint8 {
	return uint8(c & 0x1F)
}

// G extracts the green component (0-31)  
func (c Color) G() uint8 {
	return uint8((c >> 5) & 0x1F)
}

// B extracts the blue component (0-31)
func (c Color) B() uint8 {
	return uint8((c >> 10) & 0x1F)
}

// ToRGB converts to 8-bit RGB components
func (c Color) ToRGB() (r, g, b uint8) {
	return (c.R() << 3) | (c.R() >> 2),
		   (c.G() << 3) | (c.G() >> 2), 
		   (c.B() << 3) | (c.B() >> 2)
}

// Common color constants
var (
	BLACK   = RGB15(0, 0, 0)
	WHITE   = RGB15(31, 31, 31)
	RED     = RGB15(31, 0, 0)
	GREEN   = RGB15(0, 31, 0)
	BLUE    = RGB15(0, 0, 31)
	YELLOW  = RGB15(31, 31, 0)
	MAGENTA = RGB15(31, 0, 31)
	CYAN    = RGB15(0, 31, 31)
	GRAY    = RGB15(16, 16, 16)
	SILVER  = RGB15(24, 24, 24)
)

// Palette represents a color palette
type Palette [COLORS_PER_PALETTE]Color

// SetColor sets a color at the specified index
func (p *Palette) SetColor(index uint8, color Color) {
	p[index] = color
}

// GetColor gets a color at the specified index
func (p *Palette) GetColor(index uint8) Color {
	return p[index]
}