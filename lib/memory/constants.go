package memory

const (
	// Memory Base Addresses
	BIOS_BASE    = 0x00000000
	EWRAM_BASE   = 0x02000000
	IWRAM_BASE   = 0x03000000
	IO_BASE      = 0x04000000
	PALETTE_BASE = 0x05000000
	VRAM_BASE    = 0x06000000
	OAM_BASE     = 0x07000000
	ROM_BASE     = 0x08000000

	// Memory Sizes
	BIOS_SIZE    = 0x4000  // 16KB
	EWRAM_SIZE   = 0x40000 // 256KB
	IWRAM_SIZE   = 0x8000  // 32KB
	PALETTE_SIZE = 0x400   // 1KB
	VRAM_SIZE    = 0x18000 // 96KB
	OAM_SIZE     = 0x400   // 1KB

	// Screen Constants
	SCREEN_WIDTH  = 240
	SCREEN_HEIGHT = 160
	SCREEN_PIXELS = SCREEN_WIDTH * SCREEN_HEIGHT

	// Video Modes
	MODE_0 = 0 // Text mode, 4 backgrounds
	MODE_1 = 1 // Text + Affine, 3 backgrounds
	MODE_2 = 2 // Affine mode, 2 backgrounds
	MODE_3 = 3 // Bitmap 16-bit, 1 background
	MODE_4 = 4 // Bitmap 8-bit, 1 background
	MODE_5 = 5 // Bitmap 16-bit small, 1 background

	// Tile Constants
	TILE_SIZE       = 8  // 8x8 pixels
	TILE_SIZE_BYTES = 32 // 4 bits per pixel * 64 pixels / 2
	TILES_PER_ROW   = 32
	TILES_PER_COL   = 32

	// Color Constants
	COLOR_DEPTH        = 15  // 15-bit color (32768 colors)
	COLORS_PER_PALETTE = 256
	BG_PALETTES        = 16
	OBJ_PALETTES       = 16
)