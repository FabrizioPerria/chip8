package display

type Display interface {
	// Init initialises the display.
	Init() error

	// Destroy cleans up the display.
	Destroy()

	// Clear clears the display.
	Clear()

	// SetScale sets the scale of the display.
	SetScale(x, y float32)

	// DrawRect draws a rectangle on the display.
	DrawRect(x, y, w, h int, r, g, b, a uint8)

	// DrawPixel draws a pixel on the display. A pixel is a rectangle that follow the scale set by SetScale.
	DrawPixel(x, y int, on bool)

	// Update updates the display.
	Update()

	// PollEvent polls for an event.
	PollEvent() Event

	// Delay delays for a number of milliseconds.
	Delay(ms uint32)
}
