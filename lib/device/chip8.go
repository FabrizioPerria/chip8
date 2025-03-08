package device

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

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
	memorySize     = 4096
	memoryFontset  = 0x50
	memoryStart    = 0x200
	DisplayWidth   = 64
	DisplayHeight  = 32
	cycleFrequency = 500 // commands are executed at 500Hz
	timerFrequency = 60  // timers are updated at 60Hz
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

	clock struct {
		lastCycle     time.Time
		lastTimerTick time.Time
		cycleDelay    time.Duration
		timerDelay    time.Duration
	}

	display    [DisplayWidth][DisplayHeight]byte
	shouldDraw bool

	beep func()

	currentPixel struct {
		x uint8
		y uint8
	}
}

func (c *Chip8) Init() {
	c.pc = memoryStart
	c.beep = func() {}

	c.memory = [memorySize]byte{}
	c.stack = [16]uint16{}
	c.V = [16]byte{}
	c.key = [16]byte{}

	c.clearDisplay()

	copy(c.memory[memoryFontset:], fontset[:])

	c.clock.lastCycle = time.Now()
	c.clock.lastTimerTick = time.Now()
	c.clock.cycleDelay = time.Second / cycleFrequency
	c.clock.timerDelay = time.Second / timerFrequency
}

func (c *Chip8) clearDisplay() {
	c.display = [DisplayWidth][DisplayHeight]byte{}
	c.shouldDraw = true
}

func (c *Chip8) LoadFile(filename string) {
	program, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	c.LoadBytes(program)
}

func (c *Chip8) LoadBytes(program []byte) {
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
	now := time.Now()
	if now.Sub(c.clock.lastCycle) >= c.clock.cycleDelay {
		c.fetch()
		c.decode()
		// c.display[c.currentPixel.x][c.currentPixel.y] = 0
		// c.currentPixel.x += 1
		// if c.currentPixel.x == 64 {
		// 	c.currentPixel.x = 0
		// 	c.currentPixel.y += 1
		// 	if c.currentPixel.y == 32 {
		// 		c.currentPixel.y = 0
		// 	}
		// }
		// c.display[c.currentPixel.x][c.currentPixel.y] = 1
		// c.shouldDraw = true
		c.clock.lastCycle = now
	}
	c.updateTimers(&now)
}

func (c *Chip8) updateTimers(now *time.Time) {
	if now.Sub(c.clock.lastTimerTick) >= c.clock.timerDelay {
		if c.delayTimer > 0 {
			c.delayTimer--
		}
		if c.soundTimer > 0 {
			c.soundTimer--
		}
		c.clock.lastTimerTick = *now
	}
}

// 0000 0000 1110 0000
func (c *Chip8) fetch() {
	c.oc = uint16(c.memory[c.pc])<<8 | uint16(c.memory[c.pc+1])
	c.pc += 2
}

func (c *Chip8) decode() {
	opcode_type := (c.oc & 0xF000) >> 12
	x := (c.oc & 0x0F00) >> 8
	y := (c.oc & 0x00F0) >> 4
	n := c.oc & 0x000F
	nn := c.oc & 0x00FF
	nnn := c.oc & 0x0FFF
	fmt.Println(hex2str(c.oc))
	switch opcode_type {
	case 0x0:
		switch nn {
		case 0xE0:
			// 00E0 - CLS
			c.clearDisplay()
		case 0xEE:
			// 00EE - RET
			c.pc = c.stack[c.sp]
			c.sp--
		default:
			// 0nnn - SYS addr
			c.pc = nnn
		}

	case 0x1:
		// 1nnn - JP addr
		c.pc = nnn

	case 0x2:
		// 2nnn - CALL addr
		c.sp++
		c.stack[c.sp] = c.pc
		c.pc = nnn

	case 0x3:
		// 3xkk - SE Vx, byte
		if c.V[x] == byte(nn) {
			c.pc += 2
		}

	case 0x4:
		// 4xkk - SNE Vx, byte
		if c.V[x] != byte(nn) {
			c.pc += 2
		}

	case 0x5:
		// 5xy0 - SE Vx, Vy
		if c.V[x] == c.V[y] {
			c.pc += 2
		}

	case 0x6:
		// 6xkk - LD Vx, byte
		c.V[x] = byte(nn)

	case 0x7:
		// 7xkk - ADD Vx, byte
		c.V[x] += byte(nn)

	case 0x8:
		switch n {
		case 0x0:
			// 8xy0 - LD Vx, Vy
			c.V[x] = c.V[y]
		case 0x1:
			// 8xy1 - OR Vx, Vy
			c.V[x] |= c.V[y]
		case 0x2:
			// 8xy2 - AND Vx, Vy
			c.V[x] &= c.V[y]
		case 0x3:
			// 8xy3 - XOR Vx, Vy
			c.V[x] ^= c.V[y]
		case 0x4:
			// 8xy4 - ADD Vx, Vy
			sum := uint16(c.V[x]) + uint16(c.V[y])
			c.V[0xF] = 0
			if sum > 0xFF {
				c.V[0xF] = 1
			}
			c.V[x] = byte(sum)
		case 0x5:
			// 8xy5 - SUB Vx, Vy
			c.V[0xF] = 0
			if c.V[x] > c.V[y] {
				c.V[0xF] = 1
			}
			c.V[x] -= c.V[y]
		case 0x6:
			// 8xy6 - SHR Vx {, Vy}
			c.V[0xF] = c.V[x] & 0x1
			c.V[x] >>= 1
		case 0x7:
			// 8xy7 - SUBN Vx, Vy
			c.V[0xF] = 0
			if c.V[y] > c.V[x] {
				c.V[0xF] = 1
			}
			c.V[x] = c.V[y] - c.V[x]
		case 0xE:
			// 8xyE - SHL Vx {, Vy}
			c.V[0xF] = c.V[x] >> 7
			c.V[x] <<= 1

		default:
			panic("Unknown opcode - " + hex2str(opcode_type))

		}

	case 0x9:
		// 9xy0 - SNE Vx, Vy
		if c.V[x] != c.V[y] {
			c.pc += 2
		}

	case 0xA:
		// Annn - LD I, addr
		c.I = nnn

	case 0xB:
		// Bnnn - JP V0, addr
		c.pc = nnn + uint16(c.V[0])

	case 0xC:
		// Cxkk - RND Vx, byte
		c.V[x] = byte(nn) & byte(rand.Intn(256))

	case 0xD:
		// Dxyn - DRW Vx, Vy, nibble
		sprite := c.memory[c.I : c.I+n]
		x := int(c.V[x])
		y := int(c.V[y])
		c.drawSprite(x, y, sprite)

	case 0xE:
		switch nn {
		case 0x9E:
			// Ex9E - SKP Vx
			if c.key[c.V[x]] == 1 {
				c.pc += 2
			}
		case 0xA1:
			// ExA1 - SKNP Vx
			if c.key[c.V[x]] == 0 {
				c.pc += 2
			}
		default:
			panic("Unknown opcode - " + hex2str(c.oc))
		}

	case 0xF:
		switch nn {
		case 0x07:
			// Fx07 - LD Vx, DT
			c.V[x] = c.delayTimer
		case 0x0A:
			// Fx0A - LD Vx, K
			c.V[x] = c.getKeyPress()
		case 0x15:
			// Fx15 - LD DT, Vx
			c.delayTimer = c.V[x]
		case 0x18:
			// Fx18 - LD ST, Vx
			c.soundTimer = c.V[x]
		case 0x1E:
			// Fx1E - ADD I, Vx
			c.I += uint16(c.V[x])
		case 0x29:
			// Fx29 - LD F, Vx
			c.I = uint16(c.V[x]) * 5
		case 0x33:
			// Fx33 - LD B, Vx
			c.memory[c.I] = c.V[x] / 100
			c.memory[c.I+1] = (c.V[x] / 10) % 10
			c.memory[c.I+2] = c.V[x] % 10
		case 0x55:
			// Fx55 - LD [I], Vx
			copy(c.memory[c.I:], c.V[:x+1])
		case 0x65:
			// Fx65 - LD Vx, [I]
			copy(c.V[:x+1], c.memory[c.I:])
		}
	}
}

func (c *Chip8) getKeyPress() byte {
	for {
		for i := range len(c.key) {
			if c.key[i] == 1 {
				return byte(i)
			}
		}

		time.Sleep(10 * time.Millisecond)
	}
}

func (c *Chip8) SetKeyPress(key uint8) {
	c.key[key] = 1
}

func (c *Chip8) drawSprite(x, y int, sprite []byte) {
	c.V[0xF] = 0
	for j := range sprite {
		for i := range 8 {
			spritePixel := (sprite[j] & (0x80 >> uint(i)))
			if spritePixel != 0 {
				xCoord := (x + i) % DisplayWidth
				yCoord := (y + j) % DisplayHeight
				if c.display[xCoord][yCoord] == 1 {
					c.V[0xF] = 1
				}
				c.display[xCoord][yCoord] ^= 1
				c.shouldDraw = true
			}
		}
	}
}

func hex2str(h uint16) string {
	return fmt.Sprintf("%04X", h)
}
