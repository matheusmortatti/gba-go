package palette

// CreateGrayscalePalette creates a 16-color grayscale palette
func CreateGrayscalePalette() *Palette16 {
	palette := &Palette16{}

	// Create grayscale gradient from black to white
	for i := 0; i < COLORS_PER_PALETTE_16; i++ {
		// Scale from 0-31 (5-bit)
		gray := uint8(i * 31 / (COLORS_PER_PALETTE_16 - 1))
		palette.SetColor(i, RGB15(gray, gray, gray))
	}

	return palette
}

// CreateRainbowPalette creates a 16-color rainbow palette
func CreateRainbowPalette() *Palette16 {
	palette := &Palette16{}

	colors := []Color{
		BLACK,            // 0 - transparent
		RED,              // 1
		RGB15(31, 15, 0), // Orange
		YELLOW,           // 3
		RGB15(15, 31, 0), // Yellow-green
		GREEN,            // 5
		RGB15(0, 31, 15), // Cyan-green
		CYAN,             // 7
		RGB15(0, 15, 31), // Light blue
		BLUE,             // 9
		RGB15(15, 0, 31), // Purple
		MAGENTA,          // 11
		RGB15(31, 0, 15), // Pink
		WHITE,            // 13
		LIGHT_GRAY,       // 14
		DARK_GRAY,        // 15
	}

	for i, color := range colors {
		palette.SetColor(i, color)
	}

	return palette
}

// CreateFirePalette creates a fire-themed 16-color palette
func CreateFirePalette() *Palette16 {
	palette := &Palette16{}

	colors := []Color{
		BLACK,             // 0 - transparent
		RGB15(8, 0, 0),    // Dark red
		RGB15(16, 0, 0),   // Medium dark red
		RGB15(24, 0, 0),   // Red
		RGB15(31, 0, 0),   // Bright red
		RGB15(31, 8, 0),   // Red-orange
		RGB15(31, 16, 0),  // Orange
		RGB15(31, 24, 0),  // Bright orange
		RGB15(31, 31, 0),  // Yellow
		RGB15(31, 31, 8),  // Light yellow
		RGB15(31, 31, 16), // Pale yellow
		RGB15(31, 31, 24), // Very pale yellow
		WHITE,             // 12 - white hot
		LIGHT_GRAY,        // 13
		GRAY,              // 14
		DARK_GRAY,         // 15
	}

	for i, color := range colors {
		palette.SetColor(i, color)
	}

	return palette
}

// CreateWaterPalette creates a water-themed 16-color palette
func CreateWaterPalette() *Palette16 {
	palette := &Palette16{}

	colors := []Color{
		BLACK,             // 0 - transparent
		RGB15(0, 0, 8),    // Dark blue
		RGB15(0, 0, 16),   // Medium dark blue
		RGB15(0, 0, 24),   // Blue
		RGB15(0, 0, 31),   // Bright blue
		RGB15(0, 8, 31),   // Blue-cyan
		RGB15(0, 16, 31),  // Light blue
		RGB15(0, 24, 31),  // Cyan-blue
		RGB15(0, 31, 31),  // Cyan
		RGB15(8, 31, 31),  // Light cyan
		RGB15(16, 31, 31), // Pale cyan
		RGB15(24, 31, 31), // Very pale cyan
		WHITE,             // 12 - white foam
		LIGHT_GRAY,        // 13
		CYAN,              // 14
		BLUE,              // 15
	}

	for i, color := range colors {
		palette.SetColor(i, color)
	}

	return palette
}

// CreateEarthPalette creates an earth-themed 16-color palette
func CreateEarthPalette() *Palette16 {
	palette := &Palette16{}

	colors := []Color{
		BLACK,             // 0 - transparent
		RGB15(8, 4, 0),    // Dark brown
		RGB15(16, 8, 0),   // Brown
		RGB15(24, 12, 0),  // Medium brown
		RGB15(31, 16, 0),  // Light brown
		RGB15(31, 24, 8),  // Tan
		RGB15(24, 31, 8),  // Yellow-green
		RGB15(16, 31, 0),  // Green
		RGB15(8, 24, 0),   // Dark green
		RGB15(0, 16, 0),   // Forest green
		RGB15(16, 16, 16), // Stone gray
		RGB15(24, 24, 24), // Light stone
		RGB15(31, 31, 24), // Sand
		GRAY,              // 13
		LIGHT_GRAY,        // 14
		WHITE,             // 15
	}

	for i, color := range colors {
		palette.SetColor(i, color)
	}

	return palette
}

// CreateMetalPalette creates a metallic 16-color palette
func CreateMetalPalette() *Palette16 {
	palette := &Palette16{}

	colors := []Color{
		BLACK,             // 0 - transparent
		RGB15(4, 4, 4),    // Very dark gray
		RGB15(8, 8, 8),    // Dark gray
		RGB15(12, 12, 12), // Medium dark gray
		RGB15(16, 16, 16), // Gray
		RGB15(20, 20, 20), // Medium gray
		RGB15(24, 24, 24), // Light gray
		RGB15(28, 28, 28), // Very light gray
		RGB15(31, 31, 31), // White
		RGB15(24, 24, 31), // Blue tint
		RGB15(31, 31, 24), // Yellow tint
		RGB15(31, 24, 24), // Red tint
		RGB15(24, 31, 24), // Green tint
		LIGHT_GRAY,        // 13
		GRAY,              // 14
		DARK_GRAY,         // 15
	}

	for i, color := range colors {
		palette.SetColor(i, color)
	}

	return palette
}

// ValidatePalette checks if a palette has valid color values
func ValidatePalette(palette *Palette16) bool {
	for i := 0; i < COLORS_PER_PALETTE_16; i++ {
		color := palette.GetColor(i)
		// Check if color is within valid 15-bit range
		if uint16(color) > 0x7FFF {
			return false
		}
	}
	return true
}

// ValidatePalette256 checks if a 256-color palette has valid color values
func ValidatePalette256(palette *Palette256) bool {
	for i := 0; i < COLORS_PER_PALETTE_256; i++ {
		color := palette.GetColor(i)
		// Check if color is within valid 15-bit range
		if uint16(color) > 0x7FFF {
			return false
		}
	}
	return true
}

// ComparePalettes compares two 16-color palettes for equality
func ComparePalettes(palette1, palette2 *Palette16) bool {
	for i := 0; i < COLORS_PER_PALETTE_16; i++ {
		if palette1.GetColor(i) != palette2.GetColor(i) {
			return false
		}
	}
	return true
}

// PaletteToArray converts a Palette16 to a slice of Color values
func PaletteToArray(palette *Palette16) []Color {
	colors := make([]Color, COLORS_PER_PALETTE_16)
	for i := 0; i < COLORS_PER_PALETTE_16; i++ {
		colors[i] = palette.GetColor(i)
	}
	return colors
}

// ArrayToPalette converts a slice of Color values to a Palette16
func ArrayToPalette(colors []Color) *Palette16 {
	palette := &Palette16{}

	for i, color := range colors {
		if i >= COLORS_PER_PALETTE_16 {
			break
		}
		palette.SetColor(i, color)
	}

	return palette
}

// QuantizeColor reduces color precision for palette fitting
func QuantizeColor(color Color, levels int) Color {
	if levels <= 1 {
		return BLACK
	}
	if levels > 32 {
		levels = 32
	}

	r, g, b := color.R(), color.G(), color.B()

	// Quantize each component
	step := uint8(31 / (levels - 1))
	qr := (r / step) * step
	qg := (g / step) * step
	qb := (b / step) * step

	return RGB15(qr, qg, qb)
}

// FindClosestColor finds the closest color in a palette to the given color
func FindClosestColor(palette *Palette16, targetColor Color) (int, Color) {
	closestIndex := 0
	closestColor := palette.GetColor(0)
	minDistance := colorDistance(targetColor, closestColor)

	for i := 1; i < COLORS_PER_PALETTE_16; i++ {
		paletteColor := palette.GetColor(i)
		distance := colorDistance(targetColor, paletteColor)

		if distance < minDistance {
			minDistance = distance
			closestIndex = i
			closestColor = paletteColor
		}
	}

	return closestIndex, closestColor
}

// colorDistance calculates the Euclidean distance between two colors
func colorDistance(color1, color2 Color) int {
	r1, g1, b1 := color1.R(), color1.G(), color1.B()
	r2, g2, b2 := color2.R(), color2.G(), color2.B()

	dr := int(r1) - int(r2)
	dg := int(g1) - int(g2)
	db := int(b1) - int(b2)

	return dr*dr + dg*dg + db*db
}
