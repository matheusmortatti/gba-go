package vram

import (
	"github.com/matheusmortatti/gba-go/lib/registers"
	"github.com/matheusmortatti/gba-go/lib/video"
)

// DoubleBuffer manages double buffering for smooth animation
type DoubleBuffer struct {
	manager     *VRAMManager
	frontBuffer *BitmapBuffer
	backBuffer  *BitmapBuffer
	vsyncEnabled bool
}

// NewDoubleBuffer creates a new double buffer system
func NewDoubleBuffer(mode int) *DoubleBuffer {
	manager := NewVRAMManager(mode)
	
	if !manager.SupportsDoubleBuffering() {
		// For modes that don't support double buffering, return nil
		return nil
	}

	db := &DoubleBuffer{
		manager:      manager,
		vsyncEnabled: true,
	}
	
	db.frontBuffer = manager.GetCurrentBuffer()
	db.backBuffer = manager.GetBackBuffer()
	
	return db
}

// GetManager returns the underlying VRAM manager
func (db *DoubleBuffer) GetManager() *VRAMManager {
	return db.manager
}

// GetFrontBuffer returns the currently displayed buffer
func (db *DoubleBuffer) GetFrontBuffer() *BitmapBuffer {
	return db.frontBuffer
}

// GetBackBuffer returns the buffer being drawn to
func (db *DoubleBuffer) GetBackBuffer() *BitmapBuffer {
	return db.backBuffer
}

// SetVSync enables or disables VSync waiting during buffer swaps
func (db *DoubleBuffer) SetVSync(enabled bool) {
	db.vsyncEnabled = enabled
}

// IsVSyncEnabled returns whether VSync is enabled
func (db *DoubleBuffer) IsVSyncEnabled() bool {
	return db.vsyncEnabled
}

// Swap swaps the front and back buffers
func (db *DoubleBuffer) Swap() {
	if db.vsyncEnabled {
		video.VSync()
	}
	
	// Swap buffers in the manager
	db.manager.SwapBuffers()
	
	// Update buffer references
	db.frontBuffer = db.manager.GetCurrentBuffer()
	db.backBuffer = db.manager.GetBackBuffer()
	
	// Update the display control register to show the new front buffer
	db.updateDisplayControl()
}

// SwapWithoutVSync swaps buffers without waiting for VSync
func (db *DoubleBuffer) SwapWithoutVSync() {
	wasEnabled := db.vsyncEnabled
	db.vsyncEnabled = false
	db.Swap()
	db.vsyncEnabled = wasEnabled
}

// Present is an alias for Swap for clarity
func (db *DoubleBuffer) Present() {
	db.Swap()
}

// ClearBackBuffer clears the back buffer with the specified color
func (db *DoubleBuffer) ClearBackBuffer(color uint16) {
	db.backBuffer.Clear(color)
}

// FastClearBackBuffer clears the back buffer using DMA
func (db *DoubleBuffer) FastClearBackBuffer(color uint16) {
	db.backBuffer.FastClear(color)
}

// updateDisplayControl updates the LCD control register to show the correct buffer
func (db *DoubleBuffer) updateDisplayControl() {
	currentPage := db.manager.GetCurrentPage()
	mode := db.manager.GetMode()
	
	// For Mode 4 and Mode 5, bit 4 of DISPCNT selects the frame
	var displayValue uint16
	if currentPage == 1 {
		// Show frame 1 (set bit 4)
		displayValue = uint16(mode) | (1 << 10) | (1 << 4) // Mode + BG2 + Frame1
	} else {
		// Show frame 0 (clear bit 4)
		displayValue = uint16(mode) | (1 << 10) // Mode + BG2 + Frame0
	}
	
	registers.Lcd.DISPCNT.Set(displayValue)
}

// GetCurrentPage returns the current front buffer page (0 or 1)
func (db *DoubleBuffer) GetCurrentPage() int {
	return db.manager.GetCurrentPage()
}

// SyncBuffers synchronizes the back buffer with the front buffer
func (db *DoubleBuffer) SyncBuffers() {
	db.backBuffer.FastCopy(db.frontBuffer)
}

// TripleBuffer manages triple buffering for ultra-smooth animation
type TripleBuffer struct {
	manager        *VRAMManager
	frontBuffer    *BitmapBuffer
	backBuffer     *BitmapBuffer
	displayBuffer  *BitmapBuffer
	swapRequested  bool
	vsyncEnabled   bool
}

// NewTripleBuffer creates a new triple buffer system
// Note: True triple buffering requires additional memory management
// This is a simplified version using double buffering with swap requests
func NewTripleBuffer(mode int) *TripleBuffer {
	manager := NewVRAMManager(mode)
	
	if !manager.SupportsDoubleBuffering() {
		return nil
	}

	tb := &TripleBuffer{
		manager:       manager,
		swapRequested: false,
		vsyncEnabled:  true,
	}
	
	tb.frontBuffer = manager.GetCurrentBuffer()
	tb.backBuffer = manager.GetBackBuffer()
	tb.displayBuffer = tb.frontBuffer // Initially same as front
	
	return tb
}

// GetBackBuffer returns the buffer being drawn to
func (tb *TripleBuffer) GetBackBuffer() *BitmapBuffer {
	return tb.backBuffer
}

// GetDisplayBuffer returns the currently displayed buffer
func (tb *TripleBuffer) GetDisplayBuffer() *BitmapBuffer {
	return tb.displayBuffer
}

// RequestSwap requests a buffer swap (non-blocking)
func (tb *TripleBuffer) RequestSwap() {
	tb.swapRequested = true
}

// Update should be called each frame to handle pending swaps
func (tb *TripleBuffer) Update() {
	if tb.swapRequested {
		if tb.vsyncEnabled {
			video.VSync()
		}
		
		// Perform the swap
		tb.manager.SwapBuffers()
		tb.frontBuffer = tb.manager.GetCurrentBuffer()
		tb.backBuffer = tb.manager.GetBackBuffer()
		tb.displayBuffer = tb.frontBuffer
		
		// Update display control
		tb.updateDisplayControl()
		
		tb.swapRequested = false
	}
}

// SetVSync enables or disables VSync waiting
func (tb *TripleBuffer) SetVSync(enabled bool) {
	tb.vsyncEnabled = enabled
}

// updateDisplayControl updates the display control register
func (tb *TripleBuffer) updateDisplayControl() {
	currentPage := tb.manager.GetCurrentPage()
	mode := tb.manager.GetMode()
	
	var displayValue uint16
	if currentPage == 1 {
		displayValue = uint16(mode) | (1 << 10) | (1 << 4) // Mode + BG2 + Frame1
	} else {
		displayValue = uint16(mode) | (1 << 10) // Mode + BG2 + Frame0
	}
	
	registers.Lcd.DISPCNT.Set(displayValue)
}

// BufferSyncManager provides utilities for buffer synchronization
type BufferSyncManager struct {
	frameCount    int
	lastSyncFrame int
}

// NewBufferSyncManager creates a new buffer synchronization manager
func NewBufferSyncManager() *BufferSyncManager {
	return &BufferSyncManager{
		frameCount:    0,
		lastSyncFrame: 0,
	}
}

// Update increments the frame counter
func (bsm *BufferSyncManager) Update() {
	bsm.frameCount++
}

// GetFrameCount returns the current frame count
func (bsm *BufferSyncManager) GetFrameCount() int {
	return bsm.frameCount
}

// ShouldSync returns true if buffers should be synchronized
func (bsm *BufferSyncManager) ShouldSync(interval int) bool {
	if bsm.frameCount-bsm.lastSyncFrame >= interval {
		bsm.lastSyncFrame = bsm.frameCount
		return true
	}
	return false
}

// Reset resets the frame counter
func (bsm *BufferSyncManager) Reset() {
	bsm.frameCount = 0
	bsm.lastSyncFrame = 0
}