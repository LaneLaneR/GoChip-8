package chip8

import (
	"fmt"
	"time"

	iosystem "github.com/LaneLaneR/GoChip-8/core/ioSystem"
	"github.com/veandco/go-sdl2/sdl"
)

/*
 * Память в Chip-8 состоит из 4096 байт, 512 из которых зарезервированны
 *
 * Счетчик команд (PC) хранит адрес текущей инструкции в виде 16 битного целого числа
 * Каждая инструкция в CHIP-8 обновляет PC когда она завершена, чтобы перейти к следующей инструкции
 * Обращася по адресу, который записан в PC
 *
 * Память обычно используется для долгосрочного хранени, регистры кратковременная память
 * Регистов всего 16, V0 - VF
 * Индексный регистр это спец 16 битный регистр, который обращается
 * К определенной точке в памяти, I.
 * Он существует для чтения и записи в памяти
 *
 * Стэк позволяет переходить в подпрограммы, и используется для отслеживания куда возвращатся
 * Указатель стека SP это 8 битное целое число, которое указывает на место в стэке. Он должен иметь значение только от 0 до 15
 * То есть по сути это должно быть 4 битное число, но мы не можем выделить меньше 8 бит (1 байта) памяти
 *
 * Оба таймеры - 8 битные регистры: звуковой таймер для определения времени звукового сигнала
 * Таймер задержки для определения времени некоторых событий в игре
 * Все они отсчитывают время с частотой 60 ГЦ
 */

type Chip8 struct {
	Memory [4096]byte
	V      [16]byte // Регистры
	I      uint16   // Индекстный регистр
	PC     uint16   // Указатель на текущую операцию

	DelayTimer byte
	SoundTimer byte

	Stack [16]uint16 // Стэк
	SP    byte       // Показывает, сколько в стэке элементов
	RPL   [8]byte

	/* Display [32][64]bool
	 * Keys    [16]bool
	 */
	IO *iosystem.IoSDL
}

const (
	VF  = 0x000F
	FPS = 60
)

func NewChip8() *Chip8 {
	tmp := Chip8{
		IO: &iosystem.IoSDL{},
	}
	tmp.LoadFont()

	return &tmp
}

func (cpu8 *Chip8) StartChip8(hz int) error {
	cpuTicker := time.NewTicker(time.Second / time.Duration(hz))
	renderTicker := time.NewTicker(time.Second / FPS)

	defer func() {
		cpuTicker.Stop()
		renderTicker.Stop()
	}()

	if err := cpu8.IO.Init(); err != nil {
		panic(err)
	}

	for cpu8.PC >= 0x200 && cpu8.PC < 4096 {
		cpu8.IO.PollEvents()
		cpu8.IO.UpdateKeys()

		select {
		case <-cpuTicker.C:
			cpu8.StepOpcode()
		case <-renderTicker.C:
			if cpu8.DelayTimer > 0 { // Каждый такт этот таймер уменьшается
				cpu8.DelayTimer--
			}
			if cpu8.SoundTimer > 0 { // Каждый такт этот таймер уменьшается
				cpu8.IO.Play()
				cpu8.SoundTimer--
			} else {
				sdl.ClearQueuedAudio(cpu8.IO.DeviceID)
			}
			cpu8.IO.Draw()
		}
	}

	return nil
}

func (cpu8 *Chip8) StartDebugChip8(hz int) error {
	cpuTicker := time.NewTicker(time.Second / time.Duration(hz))
	renderTicker := time.NewTicker(time.Second / FPS)

	defer func() {
		cpuTicker.Stop()
		renderTicker.Stop()
	}()

	if err := cpu8.IO.Init(); err != nil {
		panic(err)
	}

	for cpu8.PC >= 0x200 && cpu8.PC < 4096 {
		cpu8.IO.PollEvents()
		cpu8.IO.UpdateKeys()

		select {
		case <-cpuTicker.C:
			opcode, _ := cpu8.StepOpcode()
			cpu8.DebugPrint(opcode)
		case <-renderTicker.C:
			if cpu8.DelayTimer > 0 { // Каждый такт этот таймер уменьшается
				cpu8.DelayTimer--
			}
			if cpu8.SoundTimer > 0 { // Каждый такт этот таймер уменьшается
				cpu8.IO.Play()
				cpu8.SoundTimer--
			} else {
				sdl.ClearQueuedAudio(cpu8.IO.DeviceID) // Функция для очистки очереди на воспрозводство звука
			}
			cpu8.IO.Draw()
		}
	}

	return nil
}

func (c *Chip8) DebugPrint(opcode uint16) {
	for i := 0; i <= 15; i++ {
		fmt.Printf("V%X = %d ", i, c.V[i])
		fmt.Printf("S%X = %X\n", i, c.Stack[i])
	}
	fmt.Println("")
	fmt.Printf("SP = %d ", c.SP)
	fmt.Printf("PC = %d ", c.PC)
	fmt.Printf("I = %X\n", c.I)
	fmt.Println("")
	fmt.Printf("Delay Timer = %d ", c.DelayTimer)
	fmt.Printf("Sound Timer = %d\n", c.SoundTimer)
	fmt.Println("")
	fmt.Printf("opcode:%X\n", opcode)
	fmt.Println("")
	fmt.Printf("HighMod:%X\n", c.IO.GetHighMode())
	fmt.Println("")
	fmt.Println("")
}
