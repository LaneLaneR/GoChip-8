package iosystem

import (
	"os"
	"unsafe"

	"github.com/veandco/go-sdl2/sdl"
)

// 1 2 3 C      → 1 2 3 4
// 4 5 6 D      → Q W E R
// 7 8 9 E      → A S D F
// A 0 B F      → Z X C V

var keymap = map[sdl.Scancode]int{
	sdl.SCANCODE_X: 0,
	sdl.SCANCODE_1: 1,
	sdl.SCANCODE_2: 2,
	sdl.SCANCODE_3: 3,
	sdl.SCANCODE_Q: 4,
	sdl.SCANCODE_W: 5,
	sdl.SCANCODE_E: 6,
	sdl.SCANCODE_A: 7,
	sdl.SCANCODE_S: 8,
	sdl.SCANCODE_D: 9,
	sdl.SCANCODE_Z: 0xA,
	sdl.SCANCODE_C: 0xB,
	sdl.SCANCODE_4: 0xC,
	sdl.SCANCODE_R: 0xD,
	sdl.SCANCODE_F: 0xE,
	sdl.SCANCODE_V: 0xF,
}

type IoSDL struct {
	Window   *sdl.Window
	Rendered *sdl.Renderer
	Texture  *sdl.Texture
	//		 X   Y
	Display  [64 * 32]bool
	Pixels   []uint32
	Keys     [16]bool
	DeviceID sdl.AudioDeviceID
	Beep     []byte

	Scale    int
	DrawFlag bool
}

func (io *IoSDL) Init() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return err
	}

	io.Scale = 20
	xScale := int32(64 * io.Scale)
	yScale := int32(32 * io.Scale)
	io.Pixels = make([]uint32, 64*32)
	io.Window, err = sdl.CreateWindow("GoChip-8", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		xScale, yScale, sdl.WINDOW_SHOWN)
	if err != nil {
		return err
	}

	io.Rendered, err = sdl.CreateRenderer(io.Window, -1, sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC)
	if err != nil {
		return err
	}

	io.Texture, err = io.Rendered.CreateTexture(
		sdl.PIXELFORMAT_RGBA8888,
		sdl.TEXTUREACCESS_STREAMING,
		64, 32,
	)

	spec := sdl.AudioSpec{
		Freq:     44100,
		Format:   sdl.AUDIO_S16SYS,
		Channels: 1,
		Samples:  4096,
	}

	io.DeviceID, err = sdl.OpenAudioDevice("", false, &spec, nil, 0)
	if err != nil {
		return err
	}

	sdl.PauseAudioDevice(io.DeviceID, false)

	io.GenerateBeep()

	return nil
}

func (io *IoSDL) Draw() error {
	if !io.DrawFlag {
		return nil
	}
	io.DrawFlag = false

	for i, _ := range io.Pixels {
		if io.Display[i] {
			io.Pixels[i] = 0xFFFFFFFF
		} else {
			io.Pixels[i] = 0x000000FF
		}
	}

	io.Texture.Update(nil, unsafe.Pointer(&io.Pixels[0]), 64*4) // Нужен unsafe ибо нужен указатель на первый байт памяти в массиве
	// 3 аргумент это количество байт в одной строке, один пиксель - 4 байта
	//
	io.Rendered.Copy(io.Texture, nil, nil) // Копируем текстуру в рендер
	io.Rendered.Present()                  // Показываем
	return nil
}

func (io *IoSDL) PixelUpdate(x int, y int, pixel bool) bool {
	io.DrawFlag = true
	if x < 0 || x >= 64 || y < 0 || y >= 32 {
		return false // Если выходит за пределы
	}

	idx := x + y*64

	collision := io.Display[idx] && pixel
	io.Display[idx] = io.Display[idx] != pixel

	return collision
}

func (io *IoSDL) DrawTrue() error {
	if err := io.Draw(); err != nil {
		return err
	}

	return nil
}

func (io *IoSDL) Clear() {
	for i, _ := range io.Display {
		io.Display[i] = false
	}
}

func (io *IoSDL) UpdateKeys() {
	keys := sdl.GetKeyboardState()

	for i, v := range keymap {
		if keys[i] != 0 {
			io.Keys[v] = true
		} else {
			io.Keys[v] = false
		}
	}
}

func (io *IoSDL) WaitKeyPress() int {
	for {
		keys := sdl.GetKeyboardState()

		for i, v := range keymap {
			if keys[i] != 0 {
				return v
			}
		}

		sdl.Delay(1) // Остановить текущий поток на миллисекунду
		// Если не останавливать то хост ос будет загружена на 100%
	}
}

func (io *IoSDL) GetKey(key int) bool {
	return io.Keys[key]
}

func (io *IoSDL) PollEvents() { // Проверка на закрытие
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			os.Exit(0)
		}
	}
}

func (io *IoSDL) GenerateBeep() {
	samples := 44100 * 100 / 1000     // 4410 Сэмплов в секунду
	buffer := make([]byte, samples*2) // Создаем сам звук

	freq := 440.0                   // Частота звука
	period := float64(44100) / freq // Сколько сэмплов занимает один сэмпл волны

	high := int16(3000) // Перепады
	low := int16(-3000)

	for i := 0; i < samples; i++ {
		var value int16 // Присваиваем значение

		if int(float64(i)/period)%2 == 0 { // Показывает, в какой части периода мы сейчас
			value = high
		} else {
			value = low
		}

		buffer[i*2] = byte(value)
		buffer[i*2+1] = byte(value >> 8) // Мы делим 16 битной число на два 8 битных
	}

	io.Beep = buffer
}

func (io *IoSDL) Play() {
	sdl.QueueAudio(io.DeviceID, io.Beep) // Запускаем Beep на нашем аудиодевайсе
}
