package sdl

import (
	"fmt"

	"github.com/fabrizioperria/chip8/lib/display"
	"github.com/veandco/go-sdl2/sdl"
)

type SDLDisplay struct {
	window  *sdl.Window
	surface *sdl.Surface
	xScale  float32
	yScale  float32
	width   int
	height  int
}

func New() display.Display {
	return &SDLDisplay{}
}

func (d *SDLDisplay) Init(title string, width int32, height int32) error {
	if width == 0 {
		width = 640
	}
	if height == 0 {
		height = 320
	}
	d.width = int(width)
	d.height = int(height)

	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		return err
	}

	window, err := sdl.CreateWindow(title, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}
	d.window = window

	surface, err := window.GetSurface()
	if err != nil {
		return err
	}
	d.surface = surface

	return nil
}

func (d *SDLDisplay) Destroy() {
	d.window.Destroy()
	sdl.Quit()
}

func (d *SDLDisplay) Clear() {
	d.surface.FillRect(nil, 0)
}

func (d *SDLDisplay) DrawRect(x, y, w, h int, r, g, b, a uint8) {
	rect := sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)}
	colour := sdl.Color{R: r, G: g, B: b, A: a}
	pixel := sdl.MapRGBA(d.surface.Format, colour.R, colour.G, colour.B, colour.A)
	d.surface.FillRect(&rect, pixel)
}

func (d *SDLDisplay) Update() {
	d.window.UpdateSurface()
}

func (d *SDLDisplay) PollEvent() display.Event {
	event := sdl.PollEvent()
	if event == nil {
		return nil
	}

	// Convert SDL events to our display events
	switch event.(type) {
	case *sdl.QuitEvent:
		return &display.QuitEvent{}
	// Add more event type conversions as needed
	default:
		return nil
	}
}

func (d *SDLDisplay) Delay(ms uint32) {
	sdl.Delay(ms)
}

func (d *SDLDisplay) DrawPixel(x, y int, on bool) {
	// Scale the coordinates
	scaledX := int(float32(x) * d.xScale)
	scaledY := int(float32(y) * d.yScale)

	// Scale the width and height of the pixel
	pixelWidth := int(d.xScale)
	fmt.Println("pixelWidth", pixelWidth)
	pixelHeight := int(d.yScale)
	fmt.Println("pixelHeight", pixelHeight)

	if on {
		d.DrawRect(scaledX, scaledY, pixelWidth, pixelHeight, 255, 255, 255, 255)
	} else {
		d.DrawRect(scaledX, scaledY, pixelWidth, pixelHeight, 0, 0, 0, 255)
	}
}

func (d *SDLDisplay) SetScale(x, y uint) {
	if x == 0 {
		x = 1
	}
	if y == 0 {
		y = 1
	}
	d.xScale = float32(d.width) / float32(x)
	d.yScale = float32(d.height) / float32(y)
	fmt.Println("d.xScale", d.xScale)
	fmt.Println("d.yScale", d.yScale)
}
