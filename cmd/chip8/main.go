package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/LaneLaneR/GoChip-8/internal/chip8"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	cpu8 := chip8.NewChip8()
	cpu8.LoadFont()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Введите путь к рому CHIP-8: ")
	scanner.Scan()
	text := scanner.Text()
	if err := cpu8.LoadFromFile(text); err != nil {
		panic(err)
	}

	cpu8.StartChip8()
}
