package device

// TODO: define fontset
var fontset = []uint8{}

const (
	memorySize    = 4096
	memoryStart   = 0x200
	displayWidth  = 64
	displayHeight = 32
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

	display    [displayWidth][displayHeight]byte
	shouldDraw bool

	beep func()
}
