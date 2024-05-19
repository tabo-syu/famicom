package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

type game struct{}

func NewGame() *game {
	return &game{}
}

func (g *game) Update(image *ebiten.Image) error {
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})
	ebitenutil.DebugPrint(screen, "Snake Game")
}

func (g *game) Layout(width, height int) (int, int) {
	return width, height
}
