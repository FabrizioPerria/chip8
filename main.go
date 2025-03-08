package main

import (
	"github.com/fabrizioperria/chip8/lib/device"
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

	for !currentDisplay.ShouldQuit() {
		currentDevice.Step()
		if currentDevice.ShouldDraw() {
			buffer := currentDevice.GetBuffer()
			currentDisplay.DrawBuffer(buffer)
		}
	}
}
