package main

import (
	"fmt"

	"github.com/LaneLaneR/GoChip-8/internal/chip8"
)

func main() {
	cpu8 := chip8.Chip8{}

	if err := cpu8.LoadFromFile("1-chip8-logo.ch8"); err != nil {
		panic(err)
	}

	for {
		if cpu8.PC < 0x200 || cpu8.PC >= 4095 {
			break
		}
		opcode := cpu8.GetOpcode()
		fmt.Println(opcode)
		fmt.Println(cpu8.PC)
	}
}
