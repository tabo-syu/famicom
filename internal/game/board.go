package game

import (
	clr "image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tabo-syu/famicom/internal/cpu"
)

const size = 32

type Board struct {
	cpu  *cpu.CPU
	info [size][size]clr.Color
}

func NewBoard(cpu *cpu.CPU) *Board {
	return &Board{cpu, [size][size]clr.Color{}}
}

func (b *Board) Update() {
	for i := 0x0200; i < 0x0600; i++ {
		row := (i - 0x0200) / size
		col := (i - 0x0200) % size

		data := b.cpu.Memory.Read(uint16(i))
		b.info[row][col] = c(data)
	}
}

var dot = ebiten.NewImage(1, 1)

func (b *Board) Draw(board *ebiten.Image) {
	for y := range size {
		for x := range size {
			dot.Fill(b.info[y][x])
			ops := &ebiten.DrawImageOptions{}
			ops.GeoM.Translate(float64(x), float64(y))
			ops.GeoM.Scale(float64(10), float64(10))

			board.DrawImage(dot, ops)
			dot.Clear()
		}
	}
}

func c(data byte) clr.Color {
	switch data {
	case 0:
		// Black
		return clr.RGBA{0, 0, 0, 255}
	case 1:
		// White
		return clr.RGBA{255, 255, 255, 255}
	case 2, 9:
		// Grey
		return clr.RGBA{128, 128, 128, 255}
	case 3, 10:
		// RED
		return clr.RGBA{255, 0, 0, 255}
	case 4, 11:
		// Green
		return clr.RGBA{0, 255, 0, 255}
	case 5, 12:
		// Blue
		return clr.RGBA{0, 0, 255, 255}
	case 6, 13:
		// Magenta
		return clr.RGBA{255, 0, 255, 255}
	case 7, 14:
		// Yellow
		return clr.RGBA{255, 255, 0, 255}
	case 15:
		// Cyan
		return clr.RGBA{0, 255, 255, 255}
	default:
		return clr.RGBA{255, 255, 255, 255}
	}
}
