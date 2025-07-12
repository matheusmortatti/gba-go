package palette

import "errors"

const (
	// Palette configurations
	COLORS_PER_PALETTE_16  = 16
	COLORS_PER_PALETTE_256 = 256
	SUB_PALETTES_16        = 16
	MAX_BG_PALETTES        = 16
	MAX_OBJ_PALETTES       = 16

	// Special indices
	TRANSPARENT_COLOR_INDEX = 0
)

// Palette16 represents a 16-color palette
type Palette16 [COLORS_PER_PALETTE_16]Color

// Palette256 represents a 256-color palette
type Palette256 [COLORS_PER_PALETTE_256]Color

// SetColor sets a color at the specified index in a 16-color palette
func (p *Palette16) SetColor(index int, color Color) error {
	if index < 0 || index >= COLORS_PER_PALETTE_16 {
		return errors.New("palette index out of bounds")
	}
	p[index] = color
	return nil
}

// GetColor gets a color at the specified index from a 16-color palette
func (p *Palette16) GetColor(index int) Color {
	if index < 0 || index >= COLORS_PER_PALETTE_16 {
		return BLACK // Return black for invalid indices
	}
	return p[index]
}

// SetTransparent sets the transparent color (index 0) for the palette
func (p *Palette16) SetTransparent(color Color) {
	p[TRANSPARENT_COLOR_INDEX] = color
}

// IsTransparent returns true if the given index is the transparent color index
func (p *Palette16) IsTransparent(index int) bool {
	return index == TRANSPARENT_COLOR_INDEX
}

// Copy creates a copy of the palette
func (p *Palette16) Copy() *Palette16 {
	copy := &Palette16{}
	*copy = *p
	return copy
}

// SetColor sets a color at the specified index in a 256-color palette
func (p *Palette256) SetColor(index int, color Color) error {
	if index < 0 || index >= COLORS_PER_PALETTE_256 {
		return errors.New("palette index out of bounds")
	}
	p[index] = color
	return nil
}

// GetColor gets a color at the specified index from a 256-color palette
func (p *Palette256) GetColor(index int) Color {
	if index < 0 || index >= COLORS_PER_PALETTE_256 {
		return BLACK // Return black for invalid indices
	}
	return p[index]
}

// Copy creates a copy of the palette
func (p *Palette256) Copy() *Palette256 {
	copy := &Palette256{}
	*copy = *p
	return copy
}

// GetSubPalette extracts a 16-color sub-palette from a 256-color palette
func (p *Palette256) GetSubPalette(index int) (*Palette16, error) {
	if index < 0 || index >= SUB_PALETTES_16 {
		return nil, errors.New("sub-palette index out of bounds")
	}

	subPalette := &Palette16{}
	startIndex := index * COLORS_PER_PALETTE_16

	for i := 0; i < COLORS_PER_PALETTE_16; i++ {
		subPalette[i] = p[startIndex+i]
	}

	return subPalette, nil
}

// SetSubPalette sets a 16-color sub-palette within a 256-color palette
func (p *Palette256) SetSubPalette(index int, subPalette *Palette16) error {
	if index < 0 || index >= SUB_PALETTES_16 {
		return errors.New("sub-palette index out of bounds")
	}

	startIndex := index * COLORS_PER_PALETTE_16

	for i := 0; i < COLORS_PER_PALETTE_16; i++ {
		p[startIndex+i] = subPalette[i]
	}

	return nil
}
