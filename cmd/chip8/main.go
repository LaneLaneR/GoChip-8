package main

import (
	"math/rand"
	"time"

	"github.com/LaneLaneR/GoChip-8/internal/chip8"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cpu8 := chip8.NewChip8()
	cpu8.LoadFont()
	cpu8.LoadFromFile("rom.ch8")

	cpu8.StartChip8(false)
}
