package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Test Addressing mode
// - AccumulatorMode
// - RelativeMode
// - IndirectMode

func Test_LDA_SetZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA9, 0x00, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.True(t, cpu.status.z())
}

func Test_LDA_SetNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA9, 0b0100_0000, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.True(t, cpu.status.n())
}

func Test_LDA_Immediate(t *testing.T) {
	cpu := NewCPU()
	// cpu.pc: 8000, 8001, 8002
	cpu.Load([]uint8{0xA9, 0x05, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.registerA)
}

func Test_LDA_ZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA5, 0x05, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x05, 0x11)
	cpu.Run()

	assert.Equal(t, uint8(0x11), cpu.registerA)
}

func Test_LDA_ZeroPageX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xB5, 0x05, 0x00})
	cpu.Reset()
	cpu.registerX = 0x01
	cpu.memory.Write(0x06, 0x11)
	cpu.Run()

	assert.Equal(t, uint8(0x11), cpu.registerA)
}

func Test_LDA_Absolute(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xAD, 0x11, 0x12, 0x00})
	cpu.Reset()
	cpu.memory.WriteUint16(0x12_11, 0x13)
	cpu.Run()

	assert.Equal(t, uint8(0x13), cpu.registerA)
}

func Test_LDA_AbsoluteX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xBD, 0x11, 0x12, 0x00})
	cpu.Reset()
	cpu.registerX = 0x01
	cpu.memory.WriteUint16(0x12_12, 0x13)
	cpu.Run()

	assert.Equal(t, uint8(0x13), cpu.registerA)
}

func Test_LDA_AbsoluteY(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xB9, 0x11, 0x12, 0x00})
	cpu.Reset()
	cpu.registerY = 0x01
	cpu.memory.WriteUint16(0x12_12, 0x13)
	cpu.Run()

	assert.Equal(t, uint8(0x13), cpu.registerA)
}

func Test_LDA_IndirectX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA1, 0x11, 0x00})
	cpu.Reset()
	cpu.registerX = 0x01
	cpu.memory.WriteUint16(0x12, 0x13_14)
	cpu.memory.WriteUint16(0x13_14, 0x05)
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.registerA)
}

func Test_LDA_IndirectY(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xB1, 0x11, 0x00})
	cpu.Reset()
	cpu.registerY = 0x01
	cpu.memory.Write(0x11, 0x31)
	cpu.memory.Write(0x12, 0x32)
	cpu.memory.Write(0x32_32, 0x05)
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.registerA)
}

func Test_LDX_SetZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA2, 0x00, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.True(t, cpu.status.z())
}

func Test_LDX_SetNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA2, 0b0100_0000, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.True(t, cpu.status.n())
}

func Test_LDX_Immediate(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA2, 0x05, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.registerX)
}

func Test_LDX_ZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA6, 0x05, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x05, 0x11)
	cpu.Run()

	assert.Equal(t, uint8(0x11), cpu.registerX)
}

func Test_LDX_ZeroPageY(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xB6, 0x05, 0x00})
	cpu.Reset()
	cpu.registerY = 0x01
	cpu.memory.Write(0x06, 0x11)
	cpu.Run()

	assert.Equal(t, uint8(0x11), cpu.registerX)
}

func Test_LDX_Absolute(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xAE, 0x11, 0x12, 0x00})
	cpu.Reset()
	cpu.memory.WriteUint16(0x12_11, 0x13)
	cpu.Run()

	assert.Equal(t, uint8(0x13), cpu.registerX)
}

func Test_LDX_AbsoluteY(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xBE, 0x11, 0x12, 0x00})
	cpu.Reset()
	cpu.registerY = 0x01
	cpu.memory.WriteUint16(0x12_12, 0x13)
	cpu.Run()

	assert.Equal(t, uint8(0x13), cpu.registerX)
}

func Test_LDY_SetZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA0, 0x00, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.True(t, cpu.status.z())
}

func Test_LDY_SetNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA0, 0b0100_0000, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.True(t, cpu.status.n())
}
func Test_LDY_Immediate(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA0, 0x05, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.registerY)
}

func Test_LDY_ZeroPage(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA4, 0x05, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x05, 0x11)
	cpu.Run()

	assert.Equal(t, uint8(0x11), cpu.registerY)
}

func Test_LDY_ZeroPageX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xB4, 0x05, 0x00})
	cpu.Reset()
	cpu.registerX = 0x01
	cpu.memory.Write(0x06, 0x11)
	cpu.Run()

	assert.Equal(t, uint8(0x11), cpu.registerY)
}

func Test_LDY_Absolute(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xAC, 0x11, 0x12, 0x00})
	cpu.Reset()
	cpu.memory.WriteUint16(0x12_11, 0x13)
	cpu.Run()

	assert.Equal(t, uint8(0x13), cpu.registerY)
}

func Test_LDY_AbsoluteX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xBC, 0x11, 0x12, 0x00})
	cpu.Reset()
	cpu.registerX = 0x01
	cpu.memory.WriteUint16(0x12_12, 0x13)
	cpu.Run()

	assert.Equal(t, uint8(0x13), cpu.registerY)
}

func Test_TAX_MoveAtoX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xAA, 0x00})
	cpu.Reset()
	cpu.registerA = 0x10
	cpu.Run()

	assert.Equal(t, cpu.registerA, uint8(0x10))
}

func Test_INC_SetZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xE6, 0x01, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x01, 0xFF)
	cpu.Run()

	assert.Equal(t, uint8(0x00), cpu.memory.Read(0x01))
	assert.True(t, cpu.status.z())
}

func Test_INC_SetNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xE6, 0x01, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x01, 0b0011_1111)
	cpu.Run()

	assert.Equal(t, uint8(0x40), cpu.memory.Read(0x01))
	assert.True(t, cpu.status.n())
}

func Test_INC_Increment(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xE6, 0x01, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x01, 0x02)
	cpu.Run()

	assert.Equal(t, uint8(0x03), cpu.memory.Read(0x01))
}

func Test_INC_Overflow(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xE6, 0x01, 0xE6, 0x01, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x01, 0xFF)
	cpu.Run()

	assert.Equal(t, uint8(0x01), cpu.memory.Read(0x01))
}

func Test_INX_Increment(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xE8, 0xE8, 0x00})
	cpu.Reset()
	cpu.registerX = 0x01
	cpu.Run()

	assert.Equal(t, uint8(0x03), cpu.registerX)
}

func Test_INX_OverflowX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xE8, 0xE8, 0x00})
	cpu.Reset()
	cpu.registerX = 0xFF
	cpu.Run()

	assert.Equal(t, uint8(0x01), cpu.registerX)
}

func Test_INY_Increment(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xC8, 0xC8, 0x00})
	cpu.Reset()
	cpu.registerY = 0x01
	cpu.Run()

	assert.Equal(t, uint8(0x03), cpu.registerY)
}

func Test_INY_OverflowY(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xC8, 0xC8, 0x00})
	cpu.Reset()
	cpu.registerY = 0xFF
	cpu.Run()

	assert.Equal(t, uint8(0x01), cpu.registerY)
}

func Test_STA_Immediate(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x85, 0x01, 0x00})
	cpu.Reset()
	cpu.registerA = 0x05
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.memory.Read(0x01))
}

func Test_STX_Immediate(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x86, 0x01, 0x00})
	cpu.Reset()
	cpu.registerX = 0x05
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.memory.Read(0x01))
}

func Test_STY_Immediate(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x84, 0x01, 0x00})
	cpu.Reset()
	cpu.registerY = 0x05
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.memory.Read(0x01))
}

func Test_5OpsWorkingTogether(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xa9, 0xc0, 0xaa, 0xe8, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.Equal(t, uint8(0xC1), cpu.registerX)
}
