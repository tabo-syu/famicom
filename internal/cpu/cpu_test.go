package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO: Test Addressing mode
// - AccumulatorMode
// - RelativeMode
// - IndirectMode

func Test_AND_Accumulator(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x29, 0b1001_0110, 0x00})
	cpu.Reset()
	cpu.registerA = 0b0000_1111
	cpu.Run()

	assert.Equal(t, uint8(0b0000_0110), cpu.registerA)
}

func Test_AND_SetZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x29, 0b0000_0000, 0x00})
	cpu.Reset()
	cpu.registerA = 0b0000_0000
	cpu.Run()

	assert.True(t, cpu.status.z())
	assert.False(t, cpu.status.n())
}

func Test_AND_SetNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x29, 0b0100_0000, 0x00})
	cpu.Reset()
	cpu.registerA = 0b0100_0000
	cpu.Run()

	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
}

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

func Test_ORA_Accumulator(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x09, 0b1001_0110, 0x00})
	cpu.Reset()
	cpu.registerA = 0b0000_1111
	cpu.Run()

	assert.Equal(t, uint8(0b1001_1111), cpu.registerA)
}

func Test_ORA_SetZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x09, 0b0000_0000, 0x00})
	cpu.Reset()
	cpu.registerA = 0b0000_0000
	cpu.Run()

	assert.True(t, cpu.status.z())
	assert.False(t, cpu.status.n())
}

func Test_ORA_SetNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x09, 0b0100_0000, 0x00})
	cpu.Reset()
	cpu.registerA = 0b0000_0000
	cpu.Run()

	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
}

func Test_TAX_MoveAtoX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xAA, 0x00})
	cpu.Reset()
	cpu.registerA = 0x10
	cpu.Run()

	assert.Equal(t, cpu.registerX, uint8(0x10))
}

func Test_TAY_MoveAtoY(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xA8, 0x00})
	cpu.Reset()
	cpu.registerA = 0x10
	cpu.Run()

	assert.Equal(t, cpu.registerY, uint8(0x10))
}

func Test_TSX_MoveStoX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xBA, 0x00})
	cpu.Reset()
	cpu.stackPointer = 0x10
	cpu.Run()

	assert.Equal(t, cpu.registerX, uint8(0x10))
}

func Test_TXA_MoveXtoA(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x8A, 0x00})
	cpu.Reset()
	cpu.registerX = 0x10
	cpu.Run()

	assert.Equal(t, cpu.registerA, uint8(0x10))
}

func Test_TXS_MoveXtoS(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x9A, 0x00})
	cpu.Reset()
	cpu.registerX = 0x10
	cpu.Run()

	assert.Equal(t, cpu.stackPointer, uint8(0x10))
}

func Test_TYA_MoveYtoA(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x98, 0x00})
	cpu.Reset()
	cpu.registerY = 0x10
	cpu.Run()

	assert.Equal(t, cpu.registerA, uint8(0x10))
}

func Test_CLC_UnsetCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x18, 0x00})
	cpu.Reset()
	cpu.status.setC(true)
	assert.True(t, cpu.status.c())
	cpu.Run()

	assert.False(t, cpu.status.c())
}

func Test_CLD_UnsetDecimalFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xD8, 0x00})
	cpu.Reset()
	cpu.status.setD(true)
	assert.True(t, cpu.status.d())
	cpu.Run()

	assert.False(t, cpu.status.d())
}

func Test_CLI_UnsetInterruptDisableFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x58, 0x00})
	cpu.Reset()
	cpu.status.setI(true)
	assert.True(t, cpu.status.i())
	cpu.Run()

	assert.False(t, cpu.status.i())
}

func Test_CLV_UnsetOverflowFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xB8, 0x00})
	cpu.Reset()
	cpu.status.setO(true)
	assert.True(t, cpu.status.o())
	cpu.Run()

	assert.False(t, cpu.status.o())
}

func Test_CMP_SetZeroAndCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xC9, 0x10, 0x00})
	cpu.Reset()
	cpu.registerA = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.True(t, cpu.status.z())
}

func Test_CMP_SetCarryOnly(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xC9, 0x09, 0x00})
	cpu.Reset()
	cpu.registerA = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
}

func Test_CPX_SetZeroAndCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xE0, 0x10, 0x00})
	cpu.Reset()
	cpu.registerX = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.True(t, cpu.status.z())
}

func Test_CPX_SetCarryOnly(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xE0, 0x09, 0x00})
	cpu.Reset()
	cpu.registerX = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
}

func Test_CPY_SetZeroAndCarry(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xC0, 0x10, 0x00})
	cpu.Reset()
	cpu.registerY = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.True(t, cpu.status.z())
}

func Test_CPY_SetCarryOnly(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xC0, 0x09, 0x00})
	cpu.Reset()
	cpu.registerY = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
}

func Test_DEC_SetZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xC6, 0x02, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x02, 0x01)
	cpu.Run()

	assert.Equal(t, uint8(0x00), cpu.memory.Read(0x01))
	assert.True(t, cpu.status.z())
}

func Test_DEC_SetNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xC6, 0x01, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x01, 0b0100_0001)
	cpu.Run()

	assert.Equal(t, uint8(0x40), cpu.memory.Read(0x01))
	assert.True(t, cpu.status.n())
}

func Test_DEC_Decrement(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xC6, 0x02, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x02, 0x03)
	cpu.Run()

	assert.Equal(t, uint8(0x02), cpu.memory.Read(0x02))
}

func Test_DEC_Underflow(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xC6, 0x01, 0xC6, 0x01, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x01, 0x00)
	cpu.Run()

	assert.Equal(t, uint8(0xFE), cpu.memory.Read(0x01))
}

func Test_DEX_Decrement(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xCA, 0xCA, 0x00})
	cpu.Reset()
	cpu.registerX = 0x03
	cpu.Run()

	assert.Equal(t, uint8(0x01), cpu.registerX)
}

func Test_DEX_UnderflowX(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xCA, 0xCA, 0x00})
	cpu.Reset()
	cpu.registerX = 0x00
	cpu.Run()

	assert.Equal(t, uint8(0xFE), cpu.registerX)
}

func Test_DEY_Decrement(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x88, 0x88, 0x00})
	cpu.Reset()
	cpu.registerY = 0x03
	cpu.Run()

	assert.Equal(t, uint8(0x01), cpu.registerY)
}

func Test_EOR_Accumulator(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x49, 0b1001_0110, 0x00})
	cpu.Reset()
	cpu.registerA = 0b0000_1111
	cpu.Run()

	assert.Equal(t, uint8(0b1001_1001), cpu.registerA)
}

func Test_EOR_SetZeroFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x49, 0b0000_0000, 0x00})
	cpu.Reset()
	cpu.registerA = 0b0000_0000
	cpu.Run()

	assert.True(t, cpu.status.z())
	assert.False(t, cpu.status.n())
}

func Test_EOR_SetNegativeFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x49, 0b0000_0000, 0x00})
	cpu.Reset()
	cpu.registerA = 0b0100_0000
	cpu.Run()

	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
}

func Test_DEY_UnderflowY(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x88, 0x88, 0x00})
	cpu.Reset()
	cpu.registerY = 0x00
	cpu.Run()

	assert.Equal(t, uint8(0xFE), cpu.registerY)
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

func Test_JMP_Absolute(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x4C, 0x30, 0x40, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x40_30, 0xE8)
	cpu.memory.Write(0x40_31, 0x00)
	cpu.Run()

	assert.Equal(t, uint16(0x40_32), cpu.programCounter)
	assert.Equal(t, uint8(0x01), cpu.registerX)
}

func Test_JMP_Indirect(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x6C, 0x30, 0x40, 0x00})
	cpu.Reset()
	cpu.memory.Write(0x40_30, 0x44)
	cpu.memory.Write(0x40_31, 0x55)
	cpu.memory.Write(0x55_44, 0xE8)
	cpu.memory.Write(0x55_45, 0x00)
	cpu.Run()

	assert.Equal(t, uint16(0x55_46), cpu.programCounter)
	assert.Equal(t, uint8(0x01), cpu.registerX)
}

func Test_STA_Immediate(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x85, 0x01, 0x00})
	cpu.Reset()
	cpu.registerA = 0x05
	cpu.Run()

	assert.Equal(t, uint8(0x05), cpu.memory.Read(0x01))
}

func Test_SEC_SetCarryFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x38, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.True(t, cpu.status.c())
}

func Test_SED_SetDecimalFlag(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0xF8, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.True(t, cpu.status.d())
}

func Test_SEI_SetInterruptDisable(t *testing.T) {
	cpu := NewCPU()
	cpu.Load([]uint8{0x78, 0x00})
	cpu.Reset()
	cpu.Run()

	assert.True(t, cpu.status.i())
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
