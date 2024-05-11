package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LDA_ImmediateLoadData(t *testing.T) {
	cpu := NewCPU()
	cpu.Interpret([]uint8{0xA9, 0x05, 0x00})

	assert.Equal(t, cpu.RegisterA, uint8(0x05))
	assert.False(t, cpu.Status.Z())
	assert.False(t, cpu.Status.N())
}

func Test_LDA_SetZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Interpret([]uint8{0xA9, 0x00, 0x00})

	assert.True(t, cpu.Status.Z())
}

func Test_TAX_MoveAtoX(t *testing.T) {
	cpu := NewCPU()
	cpu.RegisterA = 0x10
	cpu.Interpret([]uint8{0xAA, 0x00})

	assert.Equal(t, cpu.RegisterA, uint8(0x10))
}

func Test_INX_OverflowX(t *testing.T) {
	cpu := NewCPU()
	cpu.RegisterX = 0xFF
	cpu.Interpret([]uint8{0xE8, 0xE8, 0x00})

	assert.Equal(t, cpu.RegisterX, uint8(0x01))
}

func Test_5OpsWorkingTogether(t *testing.T) {
	cpu := NewCPU()
	cpu.Interpret([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00})

	assert.Equal(t, cpu.RegisterX, uint8(0xC1))
}
