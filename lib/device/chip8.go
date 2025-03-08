package device

// TODO: define fontset
var fontset = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

const (
	memorySize    = 4096
	memoryFontset = 0x50
	memoryStart   = 0x200
	DisplayWidth  = 64
	DisplayHeight = 32
)

type Chip8 struct {
	memory [memorySize]byte
	stack  [16]uint16

	oc  uint16   // opcode
	pc  uint16   // program counter
	sp  uint16   // stack pointer
	V   [16]byte // registers
	I   uint16   // index register
	key [16]byte // keypad

	delayTimer uint8
	soundTimer uint8

	display    [DisplayWidth][DisplayHeight]byte
	shouldDraw bool

	beep func()

	currentPixel struct {
		x     uint8
		y     uint8
		delay uint8
	}
}

func (c *Chip8) Init() {
	c.pc = memoryStart
	c.shouldDraw = true
	c.beep = func() {}

	for i := range c.memory {
		c.memory[i] = 0
	}

	for i := range c.stack {
		c.stack[i] = 0
	}

	for i := range c.V {
		c.V[i] = 0
	}

	for i := range c.key {
		c.key[i] = 0
	}

	for i := range c.display {
		for j := range c.display[i] {
			c.display[i][j] = 0
		}
	}

	copy(c.memory[memoryFontset:], fontset[:])
}

func (c *Chip8) Load(program []byte) {
	copy(c.memory[memoryStart:], program)
	c.pc = memoryStart
}

func (c *Chip8) ShouldDraw() bool {
	tmp := c.shouldDraw
	c.shouldDraw = false
	return tmp
}

func (c *Chip8) GetBuffer() *[DisplayWidth][DisplayHeight]byte {
	return &c.display
}

func (c *Chip8) SetKey(key uint8, state bool) {
	if state {
		c.key[key] = 1
	} else {
		c.key[key] = 0
	}
}

func (c *Chip8) SetBeep(beep func()) {
	c.beep = beep
}

func (c *Chip8) Step() {
	if c.currentPixel.delay > 0 {
		c.currentPixel.delay -= 1
	} else {
		c.display[c.currentPixel.x][c.currentPixel.y] = 0
		c.currentPixel.x += 1
		if c.currentPixel.x == 64 {
			c.currentPixel.x = 0
			c.currentPixel.y += 1
			if c.currentPixel.y == 32 {
				c.currentPixel.y = 0
			}
		}
		c.currentPixel.delay = 0
		c.display[c.currentPixel.x][c.currentPixel.y] = 1
		c.shouldDraw = true
	}
}
