package main

import (
	"image/color"
	"machine"
	"math/rand"
	"runtime/interrupt"
	"runtime/volatile"
	"unsafe"

	"tinygo.org/x/tinydraw"
)

type Box struct {
	x      int16
	y      int16
	width  int16
	height int16
	color  color.RGBA
}

var (
	// Key codes
	keyDown      = uint16(895)
	keyUP        = uint16(959)
	keyLEFT      = uint16(991)
	keyRIGHT     = uint16(1007)
	keyLSHOULDER = uint16(511)
	keyRSHOULDER = uint16(767)
	keyA         = uint16(1022)
	keyB         = uint16(1021)
	keySTART     = uint16(1015)
	keySELECT    = uint16(1019)

	// Register display
	regDISPSTAT = (*volatile.Register16)(unsafe.Pointer(uintptr(0x4000004)))

	// Register keyboard
	regKEYPAD = (*volatile.Register16)(unsafe.Pointer(uintptr(0x04000130)))

	// Display from machine
	display = machine.Display

	// Screen resolution
	screenWidth, screenHeight = display.Size()

	// Colors
	black  = color.RGBA{}
	white  = color.RGBA{255, 255, 255, 255}
	green  = color.RGBA{0, 255, 0, 255}
	red    = color.RGBA{255, 0, 0, 255}
	yellow = color.RGBA{255, 255, 0, 255}
	gray   = color.RGBA{149, 165, 166, 255}

	player    Box = Box{100, 100, 15, 15, yellow}
	randomBox Box = generateRandomBox()
)

func main() {
	// Set up display
	display.Configure()

	// Register display status
	regDISPSTAT.SetBits(1<<3 | 1<<4)

	drawBox(randomBox)

	interrupt.New(machine.IRQ_VBLANK, update).Enable()

	for {
		// Infinite loop for keep alive
	}
}

func update(interrupt.Interrupt) {
	keyValue := regKEYPAD.Get()

	if keyValue == keySTART {
		// TODO
	} else if keyValue == keySELECT {
		// TODO
	} else {
		moveBox(&player, keyValue)
	}
}

func clearScreen() {
	tinydraw.FilledRectangle(display, int16(0), int16(0), screenWidth, screenHeight, black)
}

func clearBox(box Box) {
	tinydraw.FilledRectangle(display, box.x, box.y, box.width, box.height, black)
}

func drawBox(box Box) {
	tinydraw.FilledRectangle(display, box.x, box.y, box.width, box.height, box.color)
}

func moveBox(box *Box, keyValue uint16) {
	clearBox(*box)

	switch keyValue {
	case keyRIGHT:
		box.x = box.x + 10

		if box.x+box.width > screenWidth {
			box.x = screenWidth - box.width
		}
	case keyLEFT:
		box.x = box.x - 10

		if box.x < 0 {
			box.x = 0
		}
	case keyUP:
		box.y = box.y - 10

		if box.y < 0 {
			box.y = 0
		}
	case keyDown:
		box.y += 10

		if box.y+box.height > screenHeight {
			box.y = screenHeight - box.height
		}
	}

	drawBox(*box)

	eatBox()
}

func generateRandomBox() Box {
	x := rand.Intn(int(screenWidth - 15))
	y := rand.Intn(int(screenHeight - 15))

	randBox := Box{int16(x), int16(y), 15, 15, gray}

	return randBox
}

func eatBox() {
	if checkPoint(player.x, player.y, randomBox.x, randomBox.y, randomBox.width, randomBox.height) {
		clearBox(randomBox)

		randomBox.x = int16(rand.Intn(int(screenWidth - 15)))
		randomBox.y = int16(rand.Intn(int(screenHeight - 15)))

		drawBox(randomBox)
	}

	if checkPoint(player.x+player.width, player.y, randomBox.x, randomBox.y, randomBox.width, randomBox.height) {
		clearBox(randomBox)

		randomBox.x = int16(rand.Intn(int(screenWidth - 15)))
		randomBox.y = int16(rand.Intn(int(screenHeight - 15)))

		drawBox(randomBox)
	}

	if checkPoint(player.x, player.y+player.height, randomBox.x, randomBox.y, randomBox.width, randomBox.height) {
		clearBox(randomBox)

		randomBox.x = int16(rand.Intn(int(screenWidth - 15)))
		randomBox.y = int16(rand.Intn(int(screenHeight - 15)))

		drawBox(randomBox)
	}

	if checkPoint(player.x+player.width, player.y+player.height, randomBox.x, randomBox.y, randomBox.width, randomBox.height) {
		clearBox(randomBox)

		randomBox.x = int16(rand.Intn(int(screenWidth - 15)))
		randomBox.y = int16(rand.Intn(int(screenHeight - 15)))

		drawBox(randomBox)
	}
}

func checkPoint(playerX int16, playerY int16, boxX int16, boxY int16, boxWidth int16, boxHeight int16) bool {
	return playerX >= boxX && playerX <= boxX+boxWidth && playerY >= boxY && playerY <= boxY+boxHeight
}
