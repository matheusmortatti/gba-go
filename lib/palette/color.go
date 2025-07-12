package palette

const (
	// Common colors (RGB15 format)
	BLACK      = 0x0000
	WHITE      = 0x7FFF
	RED        = 0x001F
	GREEN      = 0x03E0
	BLUE       = 0x7C00
	YELLOW     = 0x03FF
	CYAN       = 0x7FE0
	MAGENTA    = 0x7C1F
	GRAY       = 0x39CE
	DARK_GRAY  = 0x1CE7
	LIGHT_GRAY = 0x5AD6
)

// Color represents a 15-bit GBA color in RGB15 format: 0bGGGGGRRRRRBBBBB
type Color uint16

// RGB15 creates a color from 5-bit RGB components (0-31)
func RGB15(r, g, b uint8) Color {
	return Color(uint16(r&0x1F) | (uint16(g&0x1F) << 5) | (uint16(b&0x1F) << 10))
}

// RGB24ToRGB15 converts 24-bit RGB values (0-255) to 15-bit GBA format
func RGB24ToRGB15(r, g, b uint8) Color {
	r5 := (r >> 3) & 0x1F
	g5 := (g >> 3) & 0x1F
	b5 := (b >> 3) & 0x1F
	return RGB15(r5, g5, b5)
}

// ToRGB24 converts the 15-bit color back to 24-bit RGB values
func (c Color) ToRGB24() (r, g, b uint8) {
	r5 := uint8(c & 0x1F)
	g5 := uint8((c >> 5) & 0x1F)
	b5 := uint8((c >> 10) & 0x1F)

	// Scale 5-bit values to 8-bit
	r = (r5 << 3) | (r5 >> 2)
	g = (g5 << 3) | (g5 >> 2)
	b = (b5 << 3) | (b5 >> 2)

	return r, g, b
}

// R returns the red component (0-31)
func (c Color) R() uint8 {
	return uint8(c & 0x1F)
}

// G returns the green component (0-31)
func (c Color) G() uint8 {
	return uint8((c >> 5) & 0x1F)
}

// B returns the blue component (0-31)
func (c Color) B() uint8 {
	return uint8((c >> 10) & 0x1F)
}

// BlendColors blends two colors with the given ratio (0.0 = color1, 1.0 = color2)
func BlendColors(color1, color2 Color, ratio float32) Color {
	if ratio <= 0 {
		return color1
	}
	if ratio >= 1 {
		return color2
	}

	r1, g1, b1 := color1.R(), color1.G(), color1.B()
	r2, g2, b2 := color2.R(), color2.G(), color2.B()

	r := uint8(float32(r1)*(1-ratio) + float32(r2)*ratio)
	g := uint8(float32(g1)*(1-ratio) + float32(g2)*ratio)
	b := uint8(float32(b1)*(1-ratio) + float32(b2)*ratio)

	return RGB15(r, g, b)
}

// CreateGradient creates a gradient between two colors with the specified number of steps
func CreateGradient(start, end Color, steps int) []Color {
	if steps <= 0 {
		return []Color{}
	}
	if steps == 1 {
		return []Color{start}
	}

	gradient := make([]Color, steps)
	gradient[0] = start
	gradient[steps-1] = end

	for i := 1; i < steps-1; i++ {
		ratio := float32(i) / float32(steps-1)
		gradient[i] = BlendColors(start, end, ratio)
	}

	return gradient
}
