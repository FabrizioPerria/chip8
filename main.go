package main

import (
	"github.com/fabrizioperria/chip8/lib/device"
	"github.com/fabrizioperria/chip8/lib/display"
	"github.com/fabrizioperria/chip8/lib/display/sdl"
)

const (
	refreshRate = 60
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

	running := true

	for running {
		currentDevice.Step()
		if currentDevice.ShouldDraw() {
			buffer := currentDevice.GetBuffer()
			for x := range 64 {
				for y := range 32 {
					currentDisplay.DrawPixel(x, y, buffer[x][y] == 1)
				}
			}
			currentDisplay.Update()
		}

		for event := currentDisplay.PollEvent(); event != nil; event = currentDisplay.PollEvent() {
			switch event.(type) {
			case *display.QuitEvent:
				running = false
			}
		}
		currentDisplay.Delay(1000 / refreshRate)
	}
}
