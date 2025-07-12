package vram

import (
	"runtime/volatile"
	"unsafe"

	"github.com/matheusmortatti/gba-go/lib/registers"
)

// DMA control flags
const (
	DMA_ENABLE      = 1 << 15
	DMA_IRQ_ENABLE  = 1 << 14
	DMA_TIMING_IMMD = 0 << 12 // Start immediately
	DMA_32BIT       = 1 << 10 // 32-bit transfer
	DMA_16BIT       = 0 << 10 // 16-bit transfer
	DMA_INC         = 0 << 7  // Increment destination
	DMA_DEC         = 1 << 7  // Decrement destination
	DMA_FIXED       = 2 << 7  // Fixed destination
	DMA_SRC_INC     = 0 << 9  // Increment source
	DMA_SRC_DEC     = 1 << 9  // Decrement source
	DMA_SRC_FIXED   = 2 << 9  // Fixed source
)

// FastClear fills the entire buffer using DMA for maximum speed
func (bb *BitmapBuffer) FastClear(color uint16) {
	bb.FastFill(0, 0, bb.width, bb.height, color)
}

// FastFill fills a rectangular area using DMA when beneficial
func (bb *BitmapBuffer) FastFill(x, y, width, height int, color uint16) {
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

	totalPixels := width * height
	
	// Use DMA for larger fills, CPU for smaller ones
	if totalPixels >= 64 {
		bb.dmaFill(x, y, width, height, color)
	} else {
		bb.FillRect(x, y, width, height, color)
	}
}

// FastCopy copies from another buffer using DMA when beneficial
func (bb *BitmapBuffer) FastCopy(src *BitmapBuffer) {
	if src.bpp != bb.bpp {
		// Different bit depths, fall back to regular copy
		bb.CopyFrom(src, 0, 0, 0, 0, src.width, src.height)
		return
	}

	totalPixels := bb.width * bb.height
	
	// Use DMA for larger copies
	if totalPixels >= 128 {
		bb.dmaCopy(src, 0, 0, 0, 0, bb.width, bb.height)
	} else {
		bb.CopyFrom(src, 0, 0, 0, 0, src.width, src.height)
	}
}

// dmaFill performs DMA-accelerated fill operation
func (bb *BitmapBuffer) dmaFill(x, y, width, height int, color uint16) {
	if bb.bpp == 16 {
		bb.dmaFill16(x, y, width, height, color)
	} else {
		bb.dmaFill8(x, y, width, height, uint8(color))
	}
}

// dmaFill16 fills using DMA for 16-bit modes
func (bb *BitmapBuffer) dmaFill16(x, y, width, height int, color uint16) {
	// For simple cases where we can fill entire rows
	if x == 0 && width == bb.width {
		// Fill entire rows
		startAddr := bb.base + uintptr(y*bb.width*2)
		wordCount := (width * height) / 2 // Number of 32-bit words
		
		if wordCount > 0 {
			color32 := uint32(color) | (uint32(color) << 16)
			bb.dma3Fill32(startAddr, color32, wordCount)
		}
		
		// Handle remaining pixels if odd count
		remaining := (width * height) % 2
		if remaining > 0 {
			addr := startAddr + uintptr(wordCount*4)
			(*volatile.Register16)(unsafe.Pointer(addr)).Set(color)
		}
	} else {
		// Fill row by row for partial fills
		for dy := 0; dy < height; dy++ {
			rowStart := bb.base + uintptr((y+dy)*bb.width*2 + x*2)
			wordCount := width / 2
			
			if wordCount > 0 {
				color32 := uint32(color) | (uint32(color) << 16)
				bb.dma3Fill32(rowStart, color32, wordCount)
			}
			
			// Handle remaining pixel if odd width
			if width%2 == 1 {
				addr := rowStart + uintptr(wordCount*4)
				(*volatile.Register16)(unsafe.Pointer(addr)).Set(color)
			}
		}
	}
}

// dmaFill8 fills using DMA for 8-bit modes
func (bb *BitmapBuffer) dmaFill8(x, y, width, height int, color uint8) {
	// For 8-bit mode, we need to be careful about alignment
	if x == 0 && width == bb.width {
		// Fill entire rows
		startAddr := bb.base + uintptr(y*bb.width)
		wordCount := (width * height) / 4 // Number of 32-bit words
		
		if wordCount > 0 {
			color32 := uint32(color) | (uint32(color) << 8) | (uint32(color) << 16) | (uint32(color) << 24)
			bb.dma3Fill32(startAddr, color32, wordCount)
		}
		
		// Handle remaining bytes
		remaining := (width * height) % 4
		addr := startAddr + uintptr(wordCount*4)
		for i := 0; i < remaining; i += 2 {
			if i+1 < remaining {
				color16 := uint16(color) | (uint16(color) << 8)
				(*volatile.Register16)(unsafe.Pointer(addr + uintptr(i))).Set(color16)
			} else {
				// Single byte remaining - need to preserve the other byte
				current := (*volatile.Register16)(unsafe.Pointer(addr + uintptr(i&^1))).Get()
				if (i % 2) == 0 {
					(*volatile.Register16)(unsafe.Pointer(addr + uintptr(i&^1))).Set((current & 0xFF00) | uint16(color))
				} else {
					(*volatile.Register16)(unsafe.Pointer(addr + uintptr(i&^1))).Set((current & 0x00FF) | (uint16(color) << 8))
				}
			}
		}
	} else {
		// Fall back to regular fill for partial fills in 8-bit mode
		bb.FillRect(x, y, width, height, uint16(color))
	}
}

// dmaCopy performs DMA-accelerated copy operation
func (bb *BitmapBuffer) dmaCopy(src *BitmapBuffer, srcX, srcY, dstX, dstY, width, height int) {
	// For now, implement simple full-buffer copy
	if srcX == 0 && srcY == 0 && dstX == 0 && dstY == 0 &&
		width == bb.width && height == bb.height &&
		width == src.width && height == src.height {
		
		var wordCount int
		if bb.bpp == 16 {
			wordCount = (width * height) / 2
		} else {
			wordCount = (width * height) / 4
		}
		
		if wordCount > 0 {
			bb.dma3Copy32(src.base, bb.base, wordCount)
		}
	} else {
		// Fall back to regular copy for partial copies
		bb.CopyFrom(src, srcX, srcY, dstX, dstY, width, height)
	}
}

// dma3Fill32 performs a 32-bit DMA fill using DMA channel 3
func (bb *BitmapBuffer) dma3Fill32(destAddr uintptr, value uint32, wordCount int) {
	if wordCount <= 0 {
		return
	}

	// Wait for any previous DMA to complete
	for registers.DmaTransferChannels.DMA3CNT_H.Get()&DMA_ENABLE != 0 {
		// Wait
	}

	// Set up source (fixed address pointing to our value)
	valueAddr := uintptr(unsafe.Pointer(&value))
	registers.DmaTransferChannels.DMA3SAD.Set(uint32(valueAddr))
	
	// Set destination
	registers.DmaTransferChannels.DMA3DAD.Set(uint32(destAddr))
	
	// Set transfer count
	registers.DmaTransferChannels.DMA3CNT_L.Set(uint16(wordCount))
	
	// Start DMA: Enable | 32-bit | Source Fixed | Dest Increment
	control := DMA_ENABLE | DMA_32BIT | DMA_SRC_FIXED | DMA_INC | DMA_TIMING_IMMD
	registers.DmaTransferChannels.DMA3CNT_H.Set(uint16(control))
	
	// Wait for completion
	for registers.DmaTransferChannels.DMA3CNT_H.Get()&DMA_ENABLE != 0 {
		// Wait
	}
}

// dma3Copy32 performs a 32-bit DMA copy using DMA channel 3
func (bb *BitmapBuffer) dma3Copy32(srcAddr, destAddr uintptr, wordCount int) {
	if wordCount <= 0 {
		return
	}

	// Wait for any previous DMA to complete
	for registers.DmaTransferChannels.DMA3CNT_H.Get()&DMA_ENABLE != 0 {
		// Wait
	}

	// Set source
	registers.DmaTransferChannels.DMA3SAD.Set(uint32(srcAddr))
	
	// Set destination  
	registers.DmaTransferChannels.DMA3DAD.Set(uint32(destAddr))
	
	// Set transfer count
	registers.DmaTransferChannels.DMA3CNT_L.Set(uint16(wordCount))
	
	// Start DMA: Enable | 32-bit | Source Increment | Dest Increment
	control := DMA_ENABLE | DMA_32BIT | DMA_SRC_INC | DMA_INC | DMA_TIMING_IMMD
	registers.DmaTransferChannels.DMA3CNT_H.Set(uint16(control))
	
	// Wait for completion
	for registers.DmaTransferChannels.DMA3CNT_H.Get()&DMA_ENABLE != 0 {
		// Wait
	}
}

// GetDMAStatus returns true if DMA channel 3 is currently active
func GetDMAStatus() bool {
	return registers.DmaTransferChannels.DMA3CNT_H.Get()&DMA_ENABLE != 0
}

// WaitForDMA waits for DMA channel 3 to complete
func WaitForDMA() {
	for GetDMAStatus() {
		// Wait
	}
}