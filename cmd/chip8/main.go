package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/LaneLaneR/GoChip-8/internal/chip8"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("Вы не передали никакого значения")
		return
	}
	rand.Seed(time.Now().UnixNano())
	cpu8 := chip8.NewChip8()
	cpu8.LoadFont()
	text := os.Args[1]
	if err := cpu8.LoadFromFile(text); err != nil {
		fmt.Println(err)
		return
	}

	cpu8.StartChip8()
}
