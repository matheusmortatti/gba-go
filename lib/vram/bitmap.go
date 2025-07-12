package vram

import (
	"errors"
	"runtime/volatile"
	"unsafe"
)

// BitmapBuffer represents a framebuffer for bitmap modes
type BitmapBuffer struct {
	base   uintptr
	width  int
	height int
	bpp    int // bits per pixel (8 or 16)
}

// NewBitmapBuffer creates a new bitmap buffer with the given parameters
func NewBitmapBuffer(base uintptr, width, height, bpp int) *BitmapBuffer {
	return &BitmapBuffer{
		base:   base,
		width:  width,
		height: height,
		bpp:    bpp,
	}
}

// GetWidth returns the buffer width
func (bb *BitmapBuffer) GetWidth() int {
	return bb.width
}

// GetHeight returns the buffer height
func (bb *BitmapBuffer) GetHeight() int {
	return bb.height
}

// GetBPP returns the bits per pixel
func (bb *BitmapBuffer) GetBPP() int {
	return bb.bpp
}

// GetBase returns the base address
func (bb *BitmapBuffer) GetBase() uintptr {
	return bb.base
}

// InBounds checks if coordinates are within buffer bounds
func (bb *BitmapBuffer) InBounds(x, y int) bool {
	return x >= 0 && x < bb.width && y >= 0 && y < bb.height
}

// getPixelAddr calculates the memory address for a pixel
func (bb *BitmapBuffer) getPixelAddr(x, y int) uintptr {
	if bb.bpp == 16 {
		return bb.base + uintptr((y*bb.width+x)*2)
	} else {
		return bb.base + uintptr(y*bb.width+x)
	}
}

// PlotPixel sets a pixel at the specified coordinates
func (bb *BitmapBuffer) PlotPixel(x, y int, color uint16) error {
	if !bb.InBounds(x, y) {
		return errors.New("pixel coordinates out of bounds")
	}

	addr := bb.getPixelAddr(x, y)
	if !InVRAMBounds(addr) {
		return errors.New("pixel address out of VRAM bounds")
	}

	if bb.bpp == 16 {
		(*volatile.Register16)(unsafe.Pointer(addr)).Set(color)
	} else {
		// For 8-bit mode, we need to handle byte alignment
		offset := uintptr(y*bb.width + x)
		if offset%2 == 0 {
			// Even offset - lower byte
			current := (*volatile.Register16)(unsafe.Pointer(bb.base + (offset &^ uintptr(1)))).Get()
			(*volatile.Register16)(unsafe.Pointer(bb.base + (offset &^ uintptr(1)))).Set((current & 0xFF00) | uint16(color&0xFF))
		} else {
			// Odd offset - upper byte
			current := (*volatile.Register16)(unsafe.Pointer(bb.base + (offset &^ uintptr(1)))).Get()
			(*volatile.Register16)(unsafe.Pointer(bb.base + (offset &^ uintptr(1)))).Set((current & 0x00FF) | (uint16(color&0xFF) << 8))
		}
	}

	return nil
}

// GetPixel gets the color of a pixel at the specified coordinates
func (bb *BitmapBuffer) GetPixel(x, y int) (uint16, error) {
	if !bb.InBounds(x, y) {
		return 0, errors.New("pixel coordinates out of bounds")
	}

	addr := bb.getPixelAddr(x, y)
	if !InVRAMBounds(addr) {
		return 0, errors.New("pixel address out of VRAM bounds")
	}

	if bb.bpp == 16 {
		return uint16((*volatile.Register16)(unsafe.Pointer(addr)).Get()), nil
	} else {
		// For 8-bit mode, extract the correct byte
		offset := uintptr(y*bb.width + x)
		value := (*volatile.Register16)(unsafe.Pointer(bb.base + (offset &^ uintptr(1)))).Get()
		if offset%2 == 0 {
			return uint16(value & 0xFF), nil
		} else {
			return uint16(value >> 8), nil
		}
	}
}

// PlotPixelFast sets a pixel without bounds checking for maximum performance
func (bb *BitmapBuffer) PlotPixelFast(x, y int, color uint16) {
	addr := bb.getPixelAddr(x, y)
	if bb.bpp == 16 {
		(*volatile.Register16)(unsafe.Pointer(addr)).Set(color)
	} else {
		offset := uintptr(y*bb.width + x)
		if offset%2 == 0 {
			current := (*volatile.Register16)(unsafe.Pointer(bb.base + (offset &^ uintptr(1)))).Get()
			(*volatile.Register16)(unsafe.Pointer(bb.base + (offset &^ uintptr(1)))).Set((current & 0xFF00) | uint16(color&0xFF))
		} else {
			current := (*volatile.Register16)(unsafe.Pointer(bb.base + (offset &^ uintptr(1)))).Get()
			(*volatile.Register16)(unsafe.Pointer(bb.base + (offset &^ uintptr(1)))).Set((current & 0x00FF) | (uint16(color&0xFF) << 8))
		}
	}
}

// Clear fills the entire buffer with the specified color
func (bb *BitmapBuffer) Clear(color uint16) {
	bb.FillRect(0, 0, bb.width, bb.height, color)
}

// FillRect fills a rectangular area with the specified color
func (bb *BitmapBuffer) FillRect(x, y, width, height int, color uint16) {
	// Clamp to buffer bounds
	if x < 0 {
		width += x
		x = 0
	}
	if y < 0 {
		height += y
		y = 0
	}
	if x+width > bb.width {
		width = bb.width - x
	}
	if y+height > bb.height {
		height = bb.height - y
	}

	if width <= 0 || height <= 0 {
		return
	}

	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			bb.PlotPixelFast(x+dx, y+dy, color)
		}
	}
}

// DrawLine draws a line between two points using Bresenham's algorithm
func (bb *BitmapBuffer) DrawLine(x1, y1, x2, y2 int, color uint16) {
	dx := x2 - x1
	dy := y2 - y1

	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	sx := -1
	if x1 < x2 {
		sx = 1
	}

	sy := -1
	if y1 < y2 {
		sy = 1
	}

	err := dx - dy

	for {
		if bb.InBounds(x1, y1) {
			bb.PlotPixelFast(x1, y1, color)
		}

		if x1 == x2 && y1 == y2 {
			break
		}

		e2 := 2 * err

		if e2 > -dy {
			err -= dy
			x1 += sx
		}

		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// CopyFrom copies a rectangular region from another buffer
func (bb *BitmapBuffer) CopyFrom(src *BitmapBuffer, srcX, srcY, dstX, dstY, width, height int) {
	// Clamp source coordinates
	if srcX < 0 {
		width += srcX
		dstX -= srcX
		srcX = 0
	}
	if srcY < 0 {
		height += srcY
		dstY -= srcY
		srcY = 0
	}
	if srcX+width > src.width {
		width = src.width - srcX
	}
	if srcY+height > src.height {
		height = src.height - srcY
	}

	// Clamp destination coordinates
	if dstX < 0 {
		width += dstX
		srcX -= dstX
		dstX = 0
	}
	if dstY < 0 {
		height += dstY
		srcY -= dstY
		dstY = 0
	}
	if dstX+width > bb.width {
		width = bb.width - dstX
	}
	if dstY+height > bb.height {
		height = bb.height - dstY
	}

	if width <= 0 || height <= 0 {
		return
	}

	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			pixel, err := src.GetPixel(srcX+dx, srcY+dy)
			if err == nil {
				bb.PlotPixelFast(dstX+dx, dstY+dy, pixel)
			}
		}
	}
}

// DrawCircle draws a circle using the midpoint circle algorithm
func (bb *BitmapBuffer) DrawCircle(centerX, centerY, radius int, color uint16) {
	x := radius
	y := 0
	err := 0

	for x >= y {
		if bb.InBounds(centerX+x, centerY+y) {
			bb.PlotPixelFast(centerX+x, centerY+y, color)
		}
		if bb.InBounds(centerX+y, centerY+x) {
			bb.PlotPixelFast(centerX+y, centerY+x, color)
		}
		if bb.InBounds(centerX-y, centerY+x) {
			bb.PlotPixelFast(centerX-y, centerY+x, color)
		}
		if bb.InBounds(centerX-x, centerY+y) {
			bb.PlotPixelFast(centerX-x, centerY+y, color)
		}
		if bb.InBounds(centerX-x, centerY-y) {
			bb.PlotPixelFast(centerX-x, centerY-y, color)
		}
		if bb.InBounds(centerX-y, centerY-x) {
			bb.PlotPixelFast(centerX-y, centerY-x, color)
		}
		if bb.InBounds(centerX+y, centerY-x) {
			bb.PlotPixelFast(centerX+y, centerY-x, color)
		}
		if bb.InBounds(centerX+x, centerY-y) {
			bb.PlotPixelFast(centerX+x, centerY-y, color)
		}

		if err <= 0 {
			y++
			err += 2*y + 1
		}

		if err > 0 {
			x--
			err -= 2*x + 1
		}
	}
}

// FillCircle draws a filled circle
func (bb *BitmapBuffer) FillCircle(centerX, centerY, radius int, color uint16) {
	for y := centerY - radius; y <= centerY+radius; y++ {
		for x := centerX - radius; x <= centerX+radius; x++ {
			dx := x - centerX
			dy := y - centerY
			if dx*dx+dy*dy <= radius*radius {
				if bb.InBounds(x, y) {
					bb.PlotPixelFast(x, y, color)
				}
			}
		}
	}
}