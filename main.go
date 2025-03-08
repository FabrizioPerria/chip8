package main

import (
	"flag"
	"log"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/fabrizioperria/chip8/lib/device"
	"github.com/fabrizioperria/chip8/lib/display/sdl"
)

var (
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to `file`")
	memprofile = flag.String("memprofile", "", "write memory profile to `file`")
)

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close() // error handling omitted for example
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}
	if *memprofile != "" {
		defer func() {
			f, err := os.Create(*memprofile)
			if err != nil {
				log.Fatal("could not create memory profile: ", err)
			}
			defer f.Close() // error handling omitted for example
			runtime.GC()    // get up-to-date statistics
			// Lookup("allocs") creates a profile similar to go test -memprofile.
			// Alternatively, use Lookup("heap") for a profile
			// that has inuse_space as the default index.
			if err := pprof.Lookup("allocs").WriteTo(f, 0); err != nil {
				log.Fatal("could not write memory profile: ", err)
			}
		}()
	}

	currentDisplay := sdl.New()
	if err := currentDisplay.Init("CHIP-8", 800, 600); err != nil {
		panic(err)
	}
	defer currentDisplay.Destroy()

	var currentDevice device.Chip8
	currentDevice.Init()
	currentDevice.LoadFile("roms/1-chip8-logo.ch8")

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
