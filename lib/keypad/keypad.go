package keypad

import "github.com/veandco/go-sdl2/sdl"

var KeyMap = map[sdl.Scancode]byte{
	sdl.SCANCODE_1: 0x1, // 1
	sdl.SCANCODE_2: 0x2, // 2
	sdl.SCANCODE_3: 0x3, // 3
	sdl.SCANCODE_4: 0xC, // C

	sdl.SCANCODE_Q: 0x4, // 4
	sdl.SCANCODE_W: 0x5, // 5
	sdl.SCANCODE_E: 0x6, // 6
	sdl.SCANCODE_R: 0xD, // D

	sdl.SCANCODE_A: 0x7, // 7
	sdl.SCANCODE_S: 0x8, // 8
	sdl.SCANCODE_D: 0x9, // 9
	sdl.SCANCODE_F: 0xE, // E

	sdl.SCANCODE_Z: 0xA, // A
	sdl.SCANCODE_X: 0x0, // 0
	sdl.SCANCODE_C: 0xB, // B
	sdl.SCANCODE_V: 0xF, // F
}

type Keypad struct {
	key [16]byte
}

func (k *Keypad) UpdateKeys() *[16]byte {
	keyboard := sdl.GetKeyboardState()

	k.reset()

	for scancode, chip8key := range KeyMap {
		if keyboard[scancode] == 1 {
			k.key[chip8key] = 1
		}
	}
	return &k.key
}

func (k *Keypad) GetPressedKey() byte {
	for i := range len(k.key) {
		if k.key[i] == 1 {
			return byte(i)
		}
	}
	return 255
}

func (k *Keypad) reset() {
	for i := range len(k.key) {
		k.key[i] = 0
	}
}
