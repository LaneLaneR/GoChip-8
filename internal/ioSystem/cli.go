package iosystem

import "fmt"

type CliIO struct {
	Keys    [16]bool
	Display [32][64]bool
}

func (c *CliIO) Clear() {
	for y, _ := range c.Display {
		for x, _ := range c.Display[y] {
			c.Display[y][x] = false
		}
	}
}

func (c *CliIO) Draw() {
	fmt.Print("\033[2J\033[H")
	for y, _ := range c.Display {
		for _, b := range c.Display[y] {
			if b {
				fmt.Print("█")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Print("\n")
	}
}

func (c *CliIO) DrawPixel(x int, y int, value bool) bool {
	old := c.Display[y][x]
	c.Display[y][x] = c.Display[y][x] != value // XOR
	return old && value
}
