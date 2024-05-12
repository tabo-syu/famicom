package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_STA_Immediate(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x85, 0x01, 0x00})
	cpu.Reset()
	cpu.RegisterA = 0x05
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.memory.Read(0x01))
}

func Test_LDA_Immediate(t *testing.T) {
	cpu := NewCPU()
	// cpu.pc: 8000, 8001, 8002
	cpu.Load([]uint8{0xA9, 0x05, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.RegisterA)
}

func Test_LDA_ZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA5, 0x05, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x05, 0x11)
	cpu.Run()

	assert.Equal(t, uint8(0x11), cpu.RegisterA)
}

func Test_LDA_ZeroPageX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xB5, 0x05, 0x00})
	cpu.Reset()
	cpu.RegisterX = 0x01
	cpu.memory.Write(0x06, 0x11)
	cpu.Run()

	assert.Equal(t, uint8(0x11), cpu.RegisterA)
}

func Test_LDA_Absolute(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xAD, 0x11, 0x12, 0x00})
	cpu.Reset()
	cpu.memory.WriteUint16(0x12_11, 0x13)
	cpu.Run()

	assert.Equal(t, uint8(0x13), cpu.RegisterA)
}

func Test_LDA_AbsoluteX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xBD, 0x11, 0x12, 0x00})
	cpu.Reset()
	cpu.RegisterX = 0x01
	cpu.memory.WriteUint16(0x12_12, 0x13)
	cpu.Run()

	assert.Equal(t, uint8(0x13), cpu.RegisterA)
}

func Test_LDA_AbsoluteY(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xB9, 0x11, 0x12, 0x00})
	cpu.Reset()
	cpu.RegisterY = 0x01
	cpu.memory.WriteUint16(0x12_12, 0x13)
	cpu.Run()

	assert.Equal(t, uint8(0x13), cpu.RegisterA)
}

func Test_LDA_IndirectX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA1, 0x11, 0x00})
	cpu.Reset()
	cpu.RegisterX = 0x01
	cpu.memory.WriteUint16(0x12, 0x13_14)
	cpu.memory.WriteUint16(0x13_14, 0x05)
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.RegisterA)
}

func Test_LDA_IndirectY(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xB1, 0x11, 0x00})
	cpu.Reset()
	cpu.RegisterY = 0x01
	cpu.memory.Write(0x11, 0x31)
	cpu.memory.Write(0x12, 0x32)
	cpu.memory.Write(0x32_32, 0x05)
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.RegisterA)
}

func Test_LDA_SetZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA9, 0x00, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.True(t, cpu.Status.Z())
}

func Test_LDA_SetNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA9, 0b0100_0000, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.True(t, cpu.Status.N())
}

func Test_TAX_MoveAtoX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xAA, 0x00})
	cpu.Reset()
	cpu.RegisterA = 0x10
	cpu.Run()

	assert.Equal(t, cpu.RegisterA, uint8(0x10))
}

func Test_INX_OverflowX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xE8, 0xE8, 0x00})
	cpu.Reset()
	cpu.RegisterX = 0xFF
	cpu.Run()

	assert.Equal(t, uint8(0x01), cpu.RegisterX)
}

func Test_5OpsWorkingTogether(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.Equal(t, uint8(0xC1), cpu.RegisterX)
}
