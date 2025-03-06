package main

import (
	"github.com/fabrizioperria/chip8/lib/display"
	"github.com/fabrizioperria/chip8/lib/display/sdl"
)

func main() {
	var currentDisplay display.Display = sdl.New()
	if err := currentDisplay.Init(); err != nil {
		panic(err)
	}
	defer currentDisplay.Destroy()

	currentDisplay.Clear()
	currentDisplay.DrawRect(0, 0, 200, 200, 255, 0, 255, 255)

	currentDisplay.Update()

	running := true

	for running {
		for event := currentDisplay.PollEvent(); event != nil; event = currentDisplay.PollEvent() {
			switch event.(type) {
			case *display.QuitEvent:
				running = false
			}
		}
		currentDisplay.Delay(33)
	}
}
