package main

import (
	"fmt"
	"time"

	"github.com/LaneLaneR/GoChip-8/internal/chip8"
)

func main() {
	cpu8 := chip8.Chip8{}

	if err := cpu8.LoadFromFile("rom.ch8"); err != nil {
		panic(err)
	}

	for {
		if cpu8.PC < 0x200 || cpu8.PC >= 4095 {
			break
		}
		err, opcode := cpu8.StepOpcode()
		if err != nil {
			panic(err)
		}

		for i := 0; i <= 15; i++ {
			fmt.Printf("V%X = %d ", i, cpu8.V[i])
			fmt.Printf("S%X = %X\n", i, cpu8.Stack[i])
		}
		fmt.Println("")
		fmt.Printf("SP = %d ", cpu8.SP)
		fmt.Printf("PC = %d\n", cpu8.PC)
		fmt.Println("")
		fmt.Printf("opcode:%X\n", opcode)
		fmt.Println("")

		time.Sleep(1 * time.Second)
	}
}
