package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/LaneLaneR/GoChip-8/internal/chip8"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cpu8 := chip8.NewChip8()
	cpu8.LoadFont()

	if err := cpu8.LoadFromFile("mySnake.ch8"); err != nil {
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
		fmt.Printf("PC = %d ", cpu8.PC)
		fmt.Printf("I = %X\n", cpu8.I)
		fmt.Println("")
		fmt.Printf("opcode:%X\n", opcode)
		fmt.Println("")
		fmt.Println("")
		fmt.Println("")
		cpu8.IO.Draw()

		time.Sleep(10 * time.Millisecond)
	}
}
