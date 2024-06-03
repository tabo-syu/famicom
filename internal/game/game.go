package game

import (
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/tabo-syu/famicom/internal/cpu"
)

const ScreenSize = 320

type game struct {
	cpu   *cpu.CPU
	rng   *rand.Rand
	board *Board
}

func NewGame(cpu *cpu.CPU, rng *rand.Rand) *game {
	return &game{
		cpu:   cpu,
		rng:   rng,
		board: NewBoard(cpu),
	}
}

func (g *game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		log.Fatal("Game exited by user")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		g.cpu.Memory.Write(0xFF, 0x77)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.cpu.Memory.Write(0xFF, 0x73)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		g.cpu.Memory.Write(0xFF, 0x61)
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		g.cpu.Memory.Write(0xFF, 0x64)
	}

	g.cpu.Memory.Write(0xFE, byte(g.rng.Intn(15)+1))
	g.board.Update()

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	g.board.Draw(screen)
}

func (g *game) Layout(width, height int) (int, int) {
	return ScreenSize, ScreenSize
}
