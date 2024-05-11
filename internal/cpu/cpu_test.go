package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_LDA_ImmediateloadAndRunData(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA9, 0x05, 0x00})
	cpu.Reset(0, 0)
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.RegisterA)
	assert.False(t, cpu.Status.Z())
	assert.False(t, cpu.Status.N())
}

func Test_LDA_SetZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA9, 0x00, 0x00})
	cpu.Reset(0, 0)
	cpu.Run()

	assert.True(t, cpu.Status.Z())
}

func Test_LDA_FromMemory(t *testing.T) {
	cpu := NewCPU()
	cpu.writeMemory(0x10, 0x55)

	cpu.Load([]uint8{0xA5, 0x10, 0x00})
	cpu.Reset(0, 0)
	cpu.Run()

	assert.Equal(t, uint8(0x55), cpu.RegisterA)
}

func Test_TAX_MoveAtoX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xAA, 0x00})
	cpu.Reset(0x10, 0)
	cpu.Run()

	assert.Equal(t, cpu.RegisterA, uint8(0x10))
}

func Test_INX_OverflowX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xE8, 0xE8, 0x00})
	cpu.Reset(0, 0xFF)
	cpu.Run()

	assert.Equal(t, uint8(0x01), cpu.RegisterX)
}

func Test_5OpsWorkingTogether(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00})
	cpu.Reset(0, 0)
	cpu.Run()

	assert.Equal(t, uint8(0xC1), cpu.RegisterX)
}
