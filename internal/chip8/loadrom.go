package chip8

import (
	"fmt"
	"os"
)

func (c *Chip8) LoadROM(rom []byte) error {
	if len(rom) > 4096-0x200 {
		return fmt.Errorf("ROM больше границ памяти")
	}

	for i, b := range rom {
		c.Memory[0x200+i] = b
	}

	c.PC = 0x200
	return nil
}

func (c *Chip8) LoadFromFile(path string) error {
	rom, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	return c.LoadROM(rom)
}
