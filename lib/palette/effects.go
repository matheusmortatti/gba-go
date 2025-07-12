package palette

import "errors"

// FadeTo creates a fade effect to a target palette over the specified number of steps
func (pm *PaletteManager) FadeTo(targetPalette *Palette16, steps int, paletteIndex int) error {
	if paletteIndex < 0 || paletteIndex >= MAX_BG_PALETTES {
		return errors.New("palette index out of bounds")
	}

	if pm.bg256Mode {
		return errors.New("fade effects not supported in 256-color mode")
	}

	currentPalette := pm.bgPalettes[paletteIndex]
	if currentPalette == nil {
		return errors.New("source palette not loaded")
	}

	if steps <= 0 {
		// Immediate transition
		return pm.LoadBGPalette16(paletteIndex, targetPalette)
	}

	// Create intermediate palettes for smooth fade
	for step := 1; step <= steps; step++ {
		ratio := float32(step) / float32(steps)
		fadePalette := &Palette16{}

		for colorIndex := 0; colorIndex < COLORS_PER_PALETTE_16; colorIndex++ {
			currentColor := currentPalette.GetColor(colorIndex)
			targetColor := targetPalette.GetColor(colorIndex)
			fadedColor := BlendColors(currentColor, targetColor, ratio)
			fadePalette.SetColor(colorIndex, fadedColor)
		}

		pm.LoadBGPalette16(paletteIndex, fadePalette)

		// In a real implementation, you would add a delay here
		// For now, this provides the palette data for manual timing
	}

	return nil
}

// FadeOBJTo creates a fade effect for sprite palettes
func (pm *PaletteManager) FadeOBJTo(targetPalette *Palette16, steps int, paletteIndex int) error {
	if paletteIndex < 0 || paletteIndex >= MAX_OBJ_PALETTES {
		return errors.New("palette index out of bounds")
	}

	if pm.obj256Mode {
		return errors.New("fade effects not supported in 256-color mode")
	}

	currentPalette := pm.objPalettes[paletteIndex]
	if currentPalette == nil {
		return errors.New("source palette not loaded")
	}

	if steps <= 0 {
		// Immediate transition
		return pm.LoadOBJPalette16(paletteIndex, targetPalette)
	}

	// Create intermediate palettes for smooth fade
	for step := 1; step <= steps; step++ {
		ratio := float32(step) / float32(steps)
		fadePalette := &Palette16{}

		for colorIndex := 0; colorIndex < COLORS_PER_PALETTE_16; colorIndex++ {
			currentColor := currentPalette.GetColor(colorIndex)
			targetColor := targetPalette.GetColor(colorIndex)
			fadedColor := BlendColors(currentColor, targetColor, ratio)
			fadePalette.SetColor(colorIndex, fadedColor)
		}

		pm.LoadOBJPalette16(paletteIndex, fadePalette)
	}

	return nil
}

// RotatePalette rotates colors within a palette between startColor and endColor indices
func (pm *PaletteManager) RotatePalette(paletteIndex int, startColor, endColor int) error {
	if paletteIndex < 0 || paletteIndex >= MAX_BG_PALETTES {
		return errors.New("palette index out of bounds")
	}

	if pm.bg256Mode {
		return errors.New("palette rotation not supported in 256-color mode")
	}

	palette := pm.bgPalettes[paletteIndex]
	if palette == nil {
		return errors.New("palette not loaded")
	}

	if startColor < 0 || endColor >= COLORS_PER_PALETTE_16 || startColor >= endColor {
		return errors.New("invalid color range for rotation")
	}

	// Create rotated palette
	rotatedPalette := palette.Copy()

	// Save the first color in the range
	firstColor := palette.GetColor(startColor)

	// Shift colors left
	for i := startColor; i < endColor; i++ {
		nextColor := palette.GetColor(i + 1)
		rotatedPalette.SetColor(i, nextColor)
	}

	// Place the first color at the end
	rotatedPalette.SetColor(endColor, firstColor)

	return pm.LoadBGPalette16(paletteIndex, rotatedPalette)
}

// RotateOBJPalette rotates colors within a sprite palette
func (pm *PaletteManager) RotateOBJPalette(paletteIndex int, startColor, endColor int) error {
	if paletteIndex < 0 || paletteIndex >= MAX_OBJ_PALETTES {
		return errors.New("palette index out of bounds")
	}

	if pm.obj256Mode {
		return errors.New("palette rotation not supported in 256-color mode")
	}

	palette := pm.objPalettes[paletteIndex]
	if palette == nil {
		return errors.New("palette not loaded")
	}

	if startColor < 0 || endColor >= COLORS_PER_PALETTE_16 || startColor >= endColor {
		return errors.New("invalid color range for rotation")
	}

	// Create rotated palette
	rotatedPalette := palette.Copy()

	// Save the first color in the range
	firstColor := palette.GetColor(startColor)

	// Shift colors left
	for i := startColor; i < endColor; i++ {
		nextColor := palette.GetColor(i + 1)
		rotatedPalette.SetColor(i, nextColor)
	}

	// Place the first color at the end
	rotatedPalette.SetColor(endColor, firstColor)

	return pm.LoadOBJPalette16(paletteIndex, rotatedPalette)
}

// CreateFadeSteps generates intermediate palettes for smooth fade transitions
func CreateFadeSteps(from, to *Palette16, steps int) []*Palette16 {
	if steps <= 0 {
		return []*Palette16{}
	}

	fadeSteps := make([]*Palette16, steps+1)
	fadeSteps[0] = from.Copy()
	fadeSteps[steps] = to.Copy()

	for step := 1; step < steps; step++ {
		ratio := float32(step) / float32(steps)
		fadePalette := &Palette16{}

		for colorIndex := 0; colorIndex < COLORS_PER_PALETTE_16; colorIndex++ {
			fromColor := from.GetColor(colorIndex)
			toColor := to.GetColor(colorIndex)
			blendedColor := BlendColors(fromColor, toColor, ratio)
			fadePalette.SetColor(colorIndex, blendedColor)
		}

		fadeSteps[step] = fadePalette
	}

	return fadeSteps
}

// ModulateBrightness adjusts the brightness of all colors in a palette
func ModulateBrightness(palette *Palette16, factor float32) *Palette16 {
	if factor < 0 {
		factor = 0
	}
	if factor > 2 {
		factor = 2
	}

	modulated := &Palette16{}

	for i := 0; i < COLORS_PER_PALETTE_16; i++ {
		color := palette.GetColor(i)
		r, g, b := color.R(), color.G(), color.B()

		// Apply brightness factor
		newR := uint8(float32(r) * factor)
		newG := uint8(float32(g) * factor)
		newB := uint8(float32(b) * factor)

		// Clamp to 5-bit range
		if newR > 31 {
			newR = 31
		}
		if newG > 31 {
			newG = 31
		}
		if newB > 31 {
			newB = 31
		}

		modulated.SetColor(i, RGB15(newR, newG, newB))
	}

	return modulated
}

// CreateWaveEffect creates a wave-based color effect for palette animation
func CreateWaveEffect(basePalette *Palette16, amplitude float32, phase float32) *Palette16 {
	waveEffect := &Palette16{}

	for i := 0; i < COLORS_PER_PALETTE_16; i++ {
		color := basePalette.GetColor(i)
		r, g, b := color.R(), color.G(), color.B()

		// Apply simple wave effect (could use math.Sin in real implementation)
		// This is a simplified version for demonstration
		waveOffset := int(amplitude * float32(i%4-2) * phase)

		newR := int(r) + waveOffset
		newG := int(g) + waveOffset
		newB := int(b) + waveOffset

		// Clamp values
		if newR < 0 {
			newR = 0
		}
		if newR > 31 {
			newR = 31
		}
		if newG < 0 {
			newG = 0
		}
		if newG > 31 {
			newG = 31
		}
		if newB < 0 {
			newB = 0
		}
		if newB > 31 {
			newB = 31
		}

		waveEffect.SetColor(i, RGB15(uint8(newR), uint8(newG), uint8(newB)))
	}

	return waveEffect
}
