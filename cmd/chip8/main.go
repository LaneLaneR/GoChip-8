package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/LaneLaneR/GoChip-8/internal/chip8"
)

func main() {
	var text string

	hz := flag.Int("hz", 1000, "Hz to CPU")
	debug := flag.Bool("debug", false, "Debug chip-8")

	flag.Parse()

	if flag.Arg(0) == "" {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Enter the path to ROM: ")
		scanner.Scan()
		text = scanner.Text()
	} else {
		if os.Args[1] == "--help" || os.Args[1] == "-help" {
			return
		} else {
			text = flag.Arg(0)
		}
	}

	rand.Seed(time.Now().UnixNano())
	cpu8 := chip8.NewChip8()
	cpu8.LoadFont()

	if err := cpu8.LoadFromFile(text); err != nil {
		fmt.Println(err)
		return
	}

	if !*debug {
		cpu8.StartChip8(*hz)
	} else {
		cpu8.StartDebugChip8(*hz)
	}
}
