package palette

import "errors"

// PaletteManager handles background and sprite palettes
type PaletteManager struct {
	bgPalettes  [MAX_BG_PALETTES]*Palette16
	objPalettes [MAX_OBJ_PALETTES]*Palette16
	bg256       *Palette256
	obj256      *Palette256
	bank        *PaletteBank
	bg256Mode   bool
	obj256Mode  bool
}

// NewPaletteManager creates a new palette manager instance
func NewPaletteManager() *PaletteManager {
	return &PaletteManager{
		bank: GetPaletteBank(),
	}
}

// LoadBGPalette16 loads a 16-color palette to the specified background palette slot
func (pm *PaletteManager) LoadBGPalette16(index int, palette *Palette16) error {
	if index < 0 || index >= MAX_BG_PALETTES {
		return errors.New("background palette index out of bounds")
	}

	if pm.bg256Mode {
		return errors.New("background palette is in 256-color mode")
	}

	// Store a copy of the palette
	pm.bgPalettes[index] = palette.Copy()

	// Load to hardware
	pm.bank.LoadBGPalette16(index, palette)

	return nil
}

// LoadBGPalette256 loads a 256-color palette to background palette memory
func (pm *PaletteManager) LoadBGPalette256(palette *Palette256) error {
	// Store a copy of the palette
	pm.bg256 = palette.Copy()
	pm.bg256Mode = true

	// Clear 16-color palettes when switching to 256-color mode
	for i := range pm.bgPalettes {
		pm.bgPalettes[i] = nil
	}

	// Load to hardware
	pm.bank.LoadBGPalette256(palette)

	return nil
}

// LoadOBJPalette16 loads a 16-color palette to the specified sprite palette slot
func (pm *PaletteManager) LoadOBJPalette16(index int, palette *Palette16) error {
	if index < 0 || index >= MAX_OBJ_PALETTES {
		return errors.New("sprite palette index out of bounds")
	}

	if pm.obj256Mode {
		return errors.New("sprite palette is in 256-color mode")
	}

	// Store a copy of the palette
	pm.objPalettes[index] = palette.Copy()

	// Load to hardware
	pm.bank.LoadOBJPalette16(index, palette)

	return nil
}

// LoadOBJPalette256 loads a 256-color palette to sprite palette memory
func (pm *PaletteManager) LoadOBJPalette256(palette *Palette256) error {
	// Store a copy of the palette
	pm.obj256 = palette.Copy()
	pm.obj256Mode = true

	// Clear 16-color palettes when switching to 256-color mode
	for i := range pm.objPalettes {
		pm.objPalettes[i] = nil
	}

	// Load to hardware
	pm.bank.LoadOBJPalette256(palette)

	return nil
}

// GetBGPalette16 gets a copy of the background palette at the specified index
func (pm *PaletteManager) GetBGPalette16(index int) (*Palette16, error) {
	if index < 0 || index >= MAX_BG_PALETTES {
		return nil, errors.New("background palette index out of bounds")
	}

	if pm.bg256Mode {
		return nil, errors.New("background palette is in 256-color mode")
	}

	if pm.bgPalettes[index] == nil {
		return nil, errors.New("palette not loaded")
	}

	return pm.bgPalettes[index].Copy(), nil
}

// GetBGPalette256 gets a copy of the 256-color background palette
func (pm *PaletteManager) GetBGPalette256() (*Palette256, error) {
	if !pm.bg256Mode {
		return nil, errors.New("background palette is not in 256-color mode")
	}

	if pm.bg256 == nil {
		return nil, errors.New("256-color palette not loaded")
	}

	return pm.bg256.Copy(), nil
}

// GetOBJPalette16 gets a copy of the sprite palette at the specified index
func (pm *PaletteManager) GetOBJPalette16(index int) (*Palette16, error) {
	if index < 0 || index >= MAX_OBJ_PALETTES {
		return nil, errors.New("sprite palette index out of bounds")
	}

	if pm.obj256Mode {
		return nil, errors.New("sprite palette is in 256-color mode")
	}

	if pm.objPalettes[index] == nil {
		return nil, errors.New("palette not loaded")
	}

	return pm.objPalettes[index].Copy(), nil
}

// GetOBJPalette256 gets a copy of the 256-color sprite palette
func (pm *PaletteManager) GetOBJPalette256() (*Palette256, error) {
	if !pm.obj256Mode {
		return nil, errors.New("sprite palette is not in 256-color mode")
	}

	if pm.obj256 == nil {
		return nil, errors.New("256-color palette not loaded")
	}

	return pm.obj256.Copy(), nil
}

// SetBGMode256 switches background palette to 256-color mode
func (pm *PaletteManager) SetBGMode256(enable bool) {
	if enable {
		pm.bg256Mode = true
		// Clear 16-color palettes
		for i := range pm.bgPalettes {
			pm.bgPalettes[i] = nil
		}
	} else {
		pm.bg256Mode = false
		pm.bg256 = nil
	}
}

// SetOBJMode256 switches sprite palette to 256-color mode
func (pm *PaletteManager) SetOBJMode256(enable bool) {
	if enable {
		pm.obj256Mode = true
		// Clear 16-color palettes
		for i := range pm.objPalettes {
			pm.objPalettes[i] = nil
		}
	} else {
		pm.obj256Mode = false
		pm.obj256 = nil
	}
}

// IsBGMode256 returns true if background palette is in 256-color mode
func (pm *PaletteManager) IsBGMode256() bool {
	return pm.bg256Mode
}

// IsOBJMode256 returns true if sprite palette is in 256-color mode
func (pm *PaletteManager) IsOBJMode256() bool {
	return pm.obj256Mode
}

// ClearBGPalettes clears all background palettes
func (pm *PaletteManager) ClearBGPalettes() {
	for i := range pm.bgPalettes {
		pm.bgPalettes[i] = nil
	}
	pm.bg256 = nil
	pm.bg256Mode = false

	// Clear hardware palettes
	blackPalette := &Palette256{}
	pm.bank.LoadBGPalette256(blackPalette)
}

// ClearOBJPalettes clears all sprite palettes
func (pm *PaletteManager) ClearOBJPalettes() {
	for i := range pm.objPalettes {
		pm.objPalettes[i] = nil
	}
	pm.obj256 = nil
	pm.obj256Mode = false

	// Clear hardware palettes
	blackPalette := &Palette256{}
	pm.bank.LoadOBJPalette256(blackPalette)
}
