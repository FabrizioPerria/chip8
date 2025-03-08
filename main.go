package main

import (
	"github.com/fabrizioperria/chip8/lib/device"
	"github.com/fabrizioperria/chip8/lib/display"
	"github.com/fabrizioperria/chip8/lib/display/sdl"
)

func main() {
	currentDisplay := sdl.New()
	if err := currentDisplay.Init("CHIP-8", 800, 600); err != nil {
		panic(err)
	}
	defer currentDisplay.Destroy()

	var currentDevice device.Chip8
	currentDevice.Init()

	currentDisplay.Clear()
	currentDisplay.SetScale(device.DisplayWidth, device.DisplayHeight)
	currentDisplay.DrawPixel(0, 0, true)
	currentDisplay.DrawPixel(31, 15, true)
	currentDisplay.DrawPixel(63, 31, true)

	currentDisplay.Update()

	running := true

	for running {
		// currentDevice.Step()
		// if currentDevice.ShouldDraw() {
		// 	buffer := currentDevice.GetBuffer()
		// 	for x := 0; x < display.DisplayWidth; x++ {
		// 		for y := 0; y < display.DisplayHeight; y++ {
		// 			currentDisplay.DrawPixel(x, y, buffer[x][y] == 1)
		// 		}
		// 	}
		// }

		for event := currentDisplay.PollEvent(); event != nil; event = currentDisplay.PollEvent() {
			switch event.(type) {
			case *display.QuitEvent:
				running = false
			}
		}
		currentDisplay.Delay(33)
	}
}
