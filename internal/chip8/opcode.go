package chip8

import (
	"fmt"
	"math/rand"
)

/*
 * Инструкция 16 бит хранится в 2 байтах памяти
 * Поэтому читаются 2 байта памяти
 * | - объединение 2 байтов в 1
 * <<8 перейди на старшие байты
 * Побитовое складование OR (|)
 * 0x1200
 * OR (|)
 * 0x0034
 * =
 * 0x1234
 */

func (c *Chip8) StepOpcode() (uint16, error) {
	opcode := c.getOpcode()

	c.PC += 2

	err := c.execOpcode(opcode)

	return opcode, err
}

func (c *Chip8) getOpcode() uint16 {
	if int(c.PC) >= len(c.Memory) {
		panic("PC выходит за границу")
	}

	// Сначала берется первый байт, и в опкоде он сдвигается налево, становясь старшим байтом
	// Потом берется второй байт [c.PC+1] и через | эти два байта объединяются
	// Старший байт берет левее, а младший записывается ему в конец
	// Из чего образуется двух байтовая инструкция, опкод
	opcode := uint16(c.Memory[c.PC])<<8 | uint16(c.Memory[c.PC+1])

	return opcode
}

func (c *Chip8) execOpcode(opcode uint16) error {
	/*
			 * & - оператор побитового сравнения
			 * 0xF000 это число в двочиной системе 1111 0000 0000 0000
			 * Наша задача - откинуть 12 младших бит (те, которые в 0xF000 равны нулю)
			 * Оставив только 4 бита, для этого мы воспользуемся оператором И, который применяет его к каждому биту
			 * Представим
			 * opcode: 1010 1011 1100 1101 (0xABCD)
		     * mask:   1111 0000 0000 0000 (0xF000)
			 * result: 1010 0000 0000 0000 (0xA000)
			 * Именно в этих 4 битах записана инструкция
			 * Остальные 12 бит это адрес (NNN)
	*/
	switch opcode & 0xF000 {
	case 0x0000:
		c.zeroOpcode(opcode)
	case 0x1000: // JMP NNN безусловный переход
		c.PC = opcode & 0x0FFF // Отбрасываем старшие 4 бита, получая 12 битный адрес и отдаем его PC
	case 0x2000: // CALL NNN условный переход
		if c.SP >= byte(len(c.Stack)) {
			return fmt.Errorf("Стэк переполнен")
		}
		c.Stack[c.SP] = c.PC
		c.SP++
		c.PC = opcode & 0x0FFF
	case 0x3000:
		x := ((opcode & 0x0F00) >> 8)
		value := byte(opcode & 0x00FF)
		if c.V[x] == value {
			c.PC += 2 // Сдвигаемся на 1 инструкцию, то есть на 2 байта
		}
	case 0x4000:
		x := ((opcode & 0x0F00) >> 8)
		value := byte(opcode & 0x00FF)
		if c.V[x] != value {
			c.PC += 2
		}
	case 0x5000:
		x := ((opcode & 0x0F00) >> 8) // Так как два нуля после F сдвиг на 8 бит
		y := ((opcode & 0x00F0) >> 4) // Так как один ноль после F сдвиг на 4 бит
		if c.V[x] == c.V[y] {
			c.PC += 2
		}
	case 0x6000: // LD XNN запись в регистр X значение NN
		x := ((opcode & 0x0F00) >> 8) // 0xF00 - 0000 1111 0000 0000 сдвигаются правее на 8 бит, получая адекватное число
		value := (opcode & 0x00FF)    // 0x00FF - 0000 0000 1111 1111
		c.V[x] = byte(value)
	case 0x7000: // ADD XNN сложение
		x := ((opcode & 0x0F00) >> 8)
		value := (opcode & 0x00FF)
		c.V[x] = c.V[x] + byte(value)
	case 0x8000:
		c.eightOpcode(opcode)
	case 0x9000:
		x := ((opcode & 0x0F00) >> 8)
		y := ((opcode & 0x00F0) >> 4)
		if c.V[x] != c.V[y] {
			c.PC += 2
		}
	case 0xA000:
		c.I = opcode & 0x0FFF
	case 0xB000:
		c.PC = (opcode & 0x0FFF) + uint16(c.V[0])
	case 0xC000:
		x := ((opcode & 0x0F00) >> 8)
		k := byte(opcode & 0x00FF)
		c.V[x] = byte(rand.Intn(256)) & k
	case 0xD000:
		x := int(c.V[(opcode&0x0F00)>>8])
		y := int(c.V[(opcode&0x00F0)>>4])
		n := int(opcode & 0x000F)

		c.V[VF] = 0 // Сброс флага коллизии

		for row := 0; row < n; row++ { // Идем по строкам спрайта
			sprite := c.Memory[int(c.I)+row] // Берем строку из памяти, она хранится в формате 0 и 1

			for col := 0; col < 8; col++ { // Цикл по горизонтали восьми пикселей
				pixel := (sprite >> (7 - col)) & 1 // Сдвигаем нужный бит в конец берем единичку и получаем пиксель

				if pixel == 1 { // Если пиксель равен единице
					xx := (x + col) % 64 // x + col > движение по ширине
					yy := (y + row) % 32 // y + col > движение по высоте
					// На проценты можно не смотреть, это на случай, что если символ выйдет за границы
					// То он вернется назад ибо поделится с остатком на 64 или 32

					if c.IO.PixelUpdate(xx, yy, true) {
						c.V[VF] = 1
					}
				}
			}
		}
	case 0xE000:
		c.eOpcode(opcode)
	case 0xF000:
		c.fOpcode(opcode)
	}
	return nil
}

func (c *Chip8) eightOpcode(opcode uint16) {
	x := ((opcode & 0x0F00) >> 8)
	y := ((opcode & 0x00F0) >> 4)

	switch opcode & 0x000F {
	case 0x0000: // =
		c.V[x] = c.V[y]
	case 0x0001: // OR Или
		c.V[x] = c.V[x] | c.V[y]
	case 0x0002: // AND И
		c.V[x] = c.V[x] & c.V[y]
	case 0x0003: // XOR Исключающее ИЛИ
		c.V[x] = c.V[x] ^ c.V[y]
	case 0x0004:
		sum := uint16(c.V[x] + c.V[y])

		if sum > 0x00FF {
			c.V[VF] = 1
		}

		c.V[x] = byte(sum)
	case 0x0005:
		if c.V[x] >= c.V[y] {
			c.V[VF] = 1
		} else {
			c.V[VF] = 0
		}
		c.V[x] = c.V[x] - c.V[y]
	case 0x0006:
		if c.V[x]&0x0001 == 0x0001 {
			c.V[VF] = 1
		} else {
			c.V[VF] = 0
		}

		// Modern style
		c.V[x] = c.V[x] >> 1
	case 0x0007:
		if c.V[x] >= c.V[y] {
			c.V[VF] = 1
		}

		c.V[x] = c.V[y] - c.V[x]
	case 0x000E:
		if c.V[x]&0x80 != 0 {
			c.V[VF] = 1
		} else {
			c.V[VF] = 0
		}

		c.V[x] = c.V[x] << 1
	}
}

func (c *Chip8) zeroOpcode(opcode uint16) {
	switch opcode {
	case 0x00EE:
		c.SP--
		c.PC = c.Stack[c.SP]
	case 0x00E0:
		c.IO.Clear()
	}
}

func (c *Chip8) eOpcode(opcode uint16) {
	x := int((opcode & 0x0F00) >> 8)
	vx := int(c.V[x])

	switch opcode & 0x00FF {
	case 0x009E:
		if c.IO.GetKey(vx) {
			c.PC += 2
		}

	case 0x00A1:
		if !c.IO.GetKey(vx) {
			c.PC += 2
		}
	}
}

func (c *Chip8) fOpcode(opcode uint16) error {
	x := ((opcode & 0x0F00) >> 8)

	switch opcode & 0x00FF {
	case 0x0007:
		c.V[x] = c.DelayTimer
	case 0x000A:
		key := c.IO.WaitKeyPress()
		c.V[x] = byte(key)
	case 0x0015:
		c.DelayTimer = c.V[x]
	case 0x0018:
		c.SoundTimer = c.V[x]
	case 0x001E:
		sum := c.I + uint16(c.V[x])
		c.I = sum
	case 0x0029:
		c.I = 0x050 + (5 * uint16(c.V[x]))
	case 0x0033:
		integer := c.V[x]
		hundrets := integer / 100
		teens := (integer - (hundrets * 100)) / 10
		ones := integer - (hundrets * 100) - (teens * 10)
		c.Memory[c.I] = hundrets
		c.Memory[c.I+1] = teens
		c.Memory[c.I+2] = ones
	case 0x0055:
		// Modern Style
		for i, _ := range c.V {
			if i <= int(x) {
				c.Memory[c.I+uint16(i)] = c.V[i]
			} else {
				break
			}
		}

		// c.I += uint16(x) + 1
	case 0x0065:
		for i, _ := range c.V {
			if i <= int(x) {
				c.V[i] = c.Memory[c.I+uint16(i)]
			} else {
				break
			}
		}
	}

	return nil
}
