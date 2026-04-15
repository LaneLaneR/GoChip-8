package chip8

import "fmt"

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

func (c *Chip8) StepOpcode() (error, uint16) {
	opcode := c.getOpcode()

	c.PC += 2

	err := c.execOpcode(opcode)

	return err, opcode
}

func (c *Chip8) getOpcode() uint16 {
	if int(c.PC)+1 >= len(c.Memory) {
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
		if int(c.SP) >= len(c.Stack) {
			err := fmt.Errorf("Стэк переполнен")
			return err
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
		c.V[x] = c.V[x] + c.V[y]
		if c.V[x] > 0x00FF {
			c.V[VF] = 1
		}
		c.V[VF] = 0
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
		} else {
			c.V[VF] = 0
		}

		c.V[x] = c.V[y] - c.V[x]
	case 0x000E:
		if c.V[x]&0x0001 == 0x0001 {
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
	}
}
