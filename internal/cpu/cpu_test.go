package cpu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tabo-syu/famicom/internal/bus"
	"github.com/tabo-syu/famicom/internal/memory"
)

func (cpu *CPU) loadForTest(program []byte) {
	cpu.Bus.CopyToMemory(0x03_00, program)
	cpu.Bus.WriteMemoryUint16(0x00_00, 0x03_00)
}

func Test_ADC_SetCarryFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x69, 0b1111_1111, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b0000_0001
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.n())
	assert.False(t, cpu.status.o())
	assert.True(t, cpu.status.z())
	assert.Equal(t, byte(0b0000_0000), cpu.registerA)
}

func Test_ADC_SetOverflowFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x69, 0b0111_1111, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b0000_0010
	cpu.Run()

	assert.False(t, cpu.status.c())
	assert.True(t, cpu.status.n())
	assert.True(t, cpu.status.o())
	assert.False(t, cpu.status.z())
	assert.Equal(t, byte(0b1000_0001), cpu.registerA)
}

func Test_AND_SetZeroFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x29, 0b0000_0000, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b0000_0000
	cpu.Run()

	assert.True(t, cpu.status.z())
	assert.False(t, cpu.status.n())
}

func Test_AND_SetNegativeFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x29, 0b1000_0000, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b1000_0000
	cpu.Run()

	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
}

func Test_ASL_ArithmeticShiftLeft(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x0A, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b1101_0101
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
	assert.Equal(t, byte(0b1010_1010), cpu.registerA)
}

func Test_ASL_ShiftFromMemory(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x06, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x05, 0b1101_0101)
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.True(t, cpu.status.z())
	assert.True(t, cpu.status.n())
	assert.Equal(t, byte(0b1010_1010), cpu.Bus.ReadMemory(0x05))
}

func Test_BCC_WhenSetCarry(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	// 0x8000, 0x8001, 0x8002, 0x8003
	cpu.loadForTest([]byte{0x90, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(true)
	cpu.Bus.WriteMemory(0x80_10, 0xE8)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.registerX)
	assert.Equal(t, uint16(0x80_03), cpu.ProgramCounter)
}

func Test_BCC_WhenUnsetCarry(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	// 0x8000, 0x8001, 0x8012, 0x8013
	cpu.loadForTest([]byte{0x90, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(false)
	cpu.Bus.WriteMemory(0x80_12, 0xE8)
	cpu.Bus.WriteMemory(0x80_13, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x80_14), cpu.ProgramCounter)
}

func Test_BCC_WhenMinusOperand(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	// 0x8000, 0x8001, 0x8012, 0x8013
	cpu.loadForTest([]byte{0x90, 0xF6, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(false)
	cpu.Bus.WriteMemory(0x7F_F8, 0xE8)
	cpu.Bus.WriteMemory(0x7F_F9, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x7F_FA), cpu.ProgramCounter)
}

func Test_BCS_WhenSetCarry(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	// 0x8000, 0x8001, 0x8012, 0x8013
	cpu.loadForTest([]byte{0xB0, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(true)
	cpu.Bus.WriteMemory(0x80_12, 0xE8)
	cpu.Bus.WriteMemory(0x80_13, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x80_14), cpu.ProgramCounter)
}

func Test_BCS_WhenUnsetCarry(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	// 0x8000, 0x8001, 0x8002, 0x8003
	cpu.loadForTest([]byte{0xB0, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(false)
	cpu.Bus.WriteMemory(0x80_10, 0xE8)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.registerX)
	assert.Equal(t, uint16(0x80_03), cpu.ProgramCounter)
}

func Test_BCS_WhenMinusOperand(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	// 0x8000, 0x8001, 0x8012, 0x8013
	cpu.loadForTest([]byte{0xB0, 0xF6, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(true)
	cpu.Bus.WriteMemory(0x7F_F8, 0xE8)
	cpu.Bus.WriteMemory(0x7F_F9, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x7F_FA), cpu.ProgramCounter)
}

func Test_BEQ_WhenSetZero(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xF0, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setZ(true)
	cpu.Bus.WriteMemory(0x80_12, 0xE8)
	cpu.Bus.WriteMemory(0x80_13, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x80_14), cpu.ProgramCounter)
}

func Test_BEQ_WhenUnsetZero(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xF0, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setZ(false)
	cpu.Bus.WriteMemory(0x80_10, 0xE8)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.registerX)
	assert.Equal(t, uint16(0x80_03), cpu.ProgramCounter)
}

func Test_BEQ_WhenMinusOperand(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xF0, 0xF6, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setZ(true)
	cpu.Bus.WriteMemory(0x7F_F8, 0xE8)
	cpu.Bus.WriteMemory(0x7F_F9, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x7F_FA), cpu.ProgramCounter)
}

func Test_BIT_SetNegativeFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x24, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b1010_1010
	cpu.Bus.WriteMemory(0x05, 0b1111_0000)
	cpu.Run()

	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
	assert.False(t, cpu.status.o())
}

func Test_BIT_SetOverflowFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x24, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b1100_1111
	cpu.Bus.WriteMemory(0x05, 0b1111_0000)
	cpu.Run()

	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
	assert.True(t, cpu.status.o())
}

func Test_BIT_SetZeroFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x24, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b0000_0000
	cpu.Bus.WriteMemory(0x05, 0b1111_0000)
	cpu.Run()

	assert.True(t, cpu.status.z())
	assert.False(t, cpu.status.n())
	assert.False(t, cpu.status.o())
}

func Test_BIT_Absolute(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x2C, 0x05, 0x33, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b1100_1111
	cpu.Bus.WriteMemory(0x33_05, 0b1111_0000)
	cpu.Run()

	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
	assert.True(t, cpu.status.o())
}

func Test_BMI_WhenSetNegative(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x30, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setN(true)
	cpu.Bus.WriteMemory(0x80_12, 0xE8)
	cpu.Bus.WriteMemory(0x80_13, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x80_14), cpu.ProgramCounter)
}

func Test_BMI_WhenUnsetNegative(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x30, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setN(false)
	cpu.Bus.WriteMemory(0x80_10, 0xE8)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.registerX)
	assert.Equal(t, uint16(0x80_03), cpu.ProgramCounter)
}

func Test_BMI_WhenMinusOperand(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x30, 0xF6, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setN(true)
	cpu.Bus.WriteMemory(0x7F_F8, 0xE8)
	cpu.Bus.WriteMemory(0x7F_F9, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x7F_FA), cpu.ProgramCounter)
}

func Test_BNE_WhenSetZero(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xD0, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setZ(true)
	cpu.Bus.WriteMemory(0x80_10, 0xE8)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.registerX)
	assert.Equal(t, uint16(0x80_03), cpu.ProgramCounter)
}

func Test_BNE_WhenUnsetZero(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xD0, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setZ(false)
	cpu.Bus.WriteMemory(0x80_12, 0xE8)
	cpu.Bus.WriteMemory(0x80_13, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x80_14), cpu.ProgramCounter)
}

func Test_BNE_WhenMinusOperand(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xD0, 0xF6, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setZ(false)
	cpu.Bus.WriteMemory(0x7F_F8, 0xE8)
	cpu.Bus.WriteMemory(0x7F_F9, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x7F_FA), cpu.ProgramCounter)
}

func Test_BPL_WhenSetNegative(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x10, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setN(true)
	cpu.Bus.WriteMemory(0x80_10, 0xE8)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.registerX)
	assert.Equal(t, uint16(0x80_03), cpu.ProgramCounter)
}

func Test_BPL_WhenUnsetNegative(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x10, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setN(false)
	cpu.Bus.WriteMemory(0x80_12, 0xE8)
	cpu.Bus.WriteMemory(0x80_13, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x80_14), cpu.ProgramCounter)
}

func Test_BPL_WhenMinusOperand(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x10, 0xF6, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setN(false)
	cpu.Bus.WriteMemory(0x7F_F8, 0xE8)
	cpu.Bus.WriteMemory(0x7F_F9, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x7F_FA), cpu.ProgramCounter)
}

func Test_BVC_WhenSetOverflow(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x50, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setO(true)
	cpu.Bus.WriteMemory(0x80_10, 0xE8)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.registerX)
	assert.Equal(t, uint16(0x80_03), cpu.ProgramCounter)
}

func Test_BVC_WhenUnsetOverflow(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x50, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setO(false)
	cpu.Bus.WriteMemory(0x80_12, 0xE8)
	cpu.Bus.WriteMemory(0x80_13, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x80_14), cpu.ProgramCounter)
}

func Test_BVC_WhenMinusOperand(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x50, 0xF6, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setO(false)
	cpu.Bus.WriteMemory(0x7F_F8, 0xE8)
	cpu.Bus.WriteMemory(0x7F_F9, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x7F_FA), cpu.ProgramCounter)
}

func Test_BVS_WhenSetOverflow(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x70, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setO(true)
	cpu.Bus.WriteMemory(0x80_12, 0xE8)
	cpu.Bus.WriteMemory(0x80_13, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x80_14), cpu.ProgramCounter)
}

func Test_BVS_WhenUnsetOverflow(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x70, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setO(false)
	cpu.Bus.WriteMemory(0x80_10, 0xE8)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.registerX)
	assert.Equal(t, uint16(0x80_03), cpu.ProgramCounter)
}

func Test_BVS_WhenMinusOperand(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x70, 0xF6, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setO(true)
	cpu.Bus.WriteMemory(0x7F_F8, 0xE8)
	cpu.Bus.WriteMemory(0x7F_F9, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x7F_FA), cpu.ProgramCounter)
}

func Test_LDA_SetZeroFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA9, 0x00, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.True(t, cpu.status.z())
}

func Test_LDA_SetNegativeFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA9, 0b1000_0000, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.True(t, cpu.status.n())
}

func Test_LDA_Immediate(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	// cpu.pc: 8000, 8001, 8002
	cpu.loadForTest([]byte{0xA9, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.Equal(t, byte(0x05), cpu.registerA)
}

func Test_LDA_ZeroPage(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA5, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x05, 0x11)
	cpu.Run()

	assert.Equal(t, byte(0x11), cpu.registerA)
}

func Test_LDA_ZeroPageX(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xB5, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x01
	cpu.Bus.WriteMemory(0x06, 0x11)
	cpu.Run()

	assert.Equal(t, byte(0x11), cpu.registerA)
}

func Test_LDA_Absolute(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xAD, 0x11, 0x12, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemoryUint16(0x12_11, 0x13)
	cpu.Run()

	assert.Equal(t, byte(0x13), cpu.registerA)
}

func Test_LDA_AbsoluteX(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xBD, 0x11, 0x12, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x01
	cpu.Bus.WriteMemoryUint16(0x12_12, 0x13)
	cpu.Run()

	assert.Equal(t, byte(0x13), cpu.registerA)
}

func Test_LDA_AbsoluteY(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xB9, 0x11, 0x12, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0x01
	cpu.Bus.WriteMemoryUint16(0x12_12, 0x13)
	cpu.Run()

	assert.Equal(t, byte(0x13), cpu.registerA)
}

func Test_LDA_IndirectX(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA1, 0x11, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x01
	cpu.Bus.WriteMemoryUint16(0x12, 0x13_14)
	cpu.Bus.WriteMemoryUint16(0x13_14, 0x05)
	cpu.Run()

	assert.Equal(t, byte(0x05), cpu.registerA)
}

func Test_LDA_IndirectY(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xB1, 0x11, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0x01
	cpu.Bus.WriteMemory(0x11, 0x31)
	cpu.Bus.WriteMemory(0x12, 0x32)
	cpu.Bus.WriteMemory(0x32_32, 0x05)
	cpu.Run()

	assert.Equal(t, byte(0x05), cpu.registerA)
}

func Test_LDX_SetZeroFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA2, 0x00, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.True(t, cpu.status.z())
}

func Test_LDX_SetNegativeFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA2, 0b1000_0000, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.True(t, cpu.status.n())
}

func Test_LDX_Immediate(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA2, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.Equal(t, byte(0x05), cpu.registerX)
}

func Test_LDX_ZeroPage(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA6, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x05, 0x11)
	cpu.Run()

	assert.Equal(t, byte(0x11), cpu.registerX)
}

func Test_LDX_ZeroPageY(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xB6, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0x01
	cpu.Bus.WriteMemory(0x06, 0x11)
	cpu.Run()

	assert.Equal(t, byte(0x11), cpu.registerX)
}

func Test_LDX_Absolute(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xAE, 0x11, 0x12, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemoryUint16(0x12_11, 0x13)
	cpu.Run()

	assert.Equal(t, byte(0x13), cpu.registerX)
}

func Test_LDX_AbsoluteY(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xBE, 0x11, 0x12, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0x01
	cpu.Bus.WriteMemoryUint16(0x12_12, 0x13)
	cpu.Run()

	assert.Equal(t, byte(0x13), cpu.registerX)
}

func Test_LDY_SetZeroFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA0, 0x00, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.True(t, cpu.status.z())
}

func Test_LDY_SetNegativeFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA0, 0b1000_0000, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.True(t, cpu.status.n())
}
func Test_LDY_Immediate(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA0, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.Equal(t, byte(0x05), cpu.registerY)
}

func Test_LDY_ZeroPage(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA4, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x05, 0x11)
	cpu.Run()

	assert.Equal(t, byte(0x11), cpu.registerY)
}

func Test_LDY_ZeroPageX(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xB4, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x01
	cpu.Bus.WriteMemory(0x06, 0x11)
	cpu.Run()

	assert.Equal(t, byte(0x11), cpu.registerY)
}

func Test_LDY_Absolute(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xAC, 0x11, 0x12, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemoryUint16(0x12_11, 0x13)
	cpu.Run()

	assert.Equal(t, byte(0x13), cpu.registerY)
}

func Test_LDY_AbsoluteX(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xBC, 0x11, 0x12, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x01
	cpu.Bus.WriteMemoryUint16(0x12_12, 0x13)
	cpu.Run()

	assert.Equal(t, byte(0x13), cpu.registerY)
}

func Test_LSR_LogicalShiftRight(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x4A, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b1001_0101
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
	assert.False(t, cpu.status.n())
	assert.Equal(t, byte(0b0100_1010), cpu.registerA)
}

func Test_LSR_ShiftFromMemory(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x46, 0x05, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x05, 0b1001_0101)
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
	assert.False(t, cpu.status.n())
	assert.Equal(t, byte(0b0100_1010), cpu.Bus.ReadMemory(0x05))
}

func Test_NOP_NoOperation(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xEA, 0xE8, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.Equal(t, uint16(0x80_03), cpu.ProgramCounter)
	assert.Equal(t, byte(0x01), cpu.registerX)
}

func Test_ORA_Accumulator(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x09, 0b1001_0110, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b0000_1111
	cpu.Run()

	assert.Equal(t, byte(0b1001_1111), cpu.registerA)
}

func Test_ORA_SetZeroFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x09, 0b0000_0000, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b0000_0000
	cpu.Run()

	assert.True(t, cpu.status.z())
	assert.False(t, cpu.status.n())
}

func Test_ORA_SetNegativeFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x09, 0b1000_0000, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b1000_0000
	cpu.Run()

	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
}

func Test_PHA_PushAccumulator(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x48, 0x00})
	cpu.Reset(0x00_00)
	cpu.stackPointer = stackPointer(0x05)
	cpu.registerA = 0x22
	cpu.Run()

	assert.Equal(t, byte(0x04), byte(cpu.stackPointer))
	assert.Equal(t, byte(0x22), cpu.Bus.ReadMemory(0x01_05))
}

func Test_PHP_PushStatus(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x08, 0x00})
	cpu.Reset(0x00_00)
	cpu.stackPointer = stackPointer(0x05)
	cpu.status = 0b1010_0110
	cpu.Run()

	assert.Equal(t, byte(0x04), byte(cpu.stackPointer))
	assert.Equal(t, byte(0b1010_0110), cpu.Bus.ReadMemory(0x01_05))
}

func Test_PLA_PopAccumulator(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x68, 0x00})
	cpu.Reset(0x00_00)
	cpu.stackPointer = stackPointer(0x05)
	cpu.Bus.WriteMemory(0x01_06, 0x22)
	cpu.Run()

	assert.Equal(t, byte(0x06), byte(cpu.stackPointer))
	assert.Equal(t, byte(0x22), cpu.registerA)
}

func Test_PLA_SetNegativeFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x68, 0x00})
	cpu.Reset(0x00_00)
	cpu.stackPointer = stackPointer(0x05)
	cpu.Bus.WriteMemory(0x01_06, 0b1000_0000)
	cpu.Run()

	assert.Equal(t, byte(0x06), byte(cpu.stackPointer))
	assert.Equal(t, byte(0b1000_0000), cpu.registerA)
	assert.True(t, cpu.status.n())
}

func Test_PLP_PopStatus(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x28, 0x00})
	cpu.Reset(0x00_00)
	cpu.stackPointer = stackPointer(0x05)
	cpu.Bus.WriteMemory(0x01_06, 0b1010_0110)
	cpu.Run()

	assert.Equal(t, byte(0x06), byte(cpu.stackPointer))
	assert.Equal(t, byte(0b1010_0110), byte(cpu.status))
}

func Test_ROL_Set0Bit(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x2A, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(true)
	cpu.registerA = 0b1001_0101
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
	assert.False(t, cpu.status.n())
	assert.Equal(t, byte(0b0010_1011), cpu.registerA)
}

func Test_ROL_Unset0Bit(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x2A, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(false)
	cpu.registerA = 0b0101_1001
	cpu.Run()

	assert.False(t, cpu.status.c())
	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
	assert.Equal(t, byte(0b1011_0010), cpu.registerA)
}

func Test_ROL_RotateFromMemory(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x2E, 0x30, 0x40, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(true)
	cpu.Bus.WriteMemory(0x40_30, 0b1001_0101)
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
	assert.False(t, cpu.status.n())
	assert.Equal(t, byte(0b0010_1011), cpu.Bus.ReadMemory(0x40_30))
}

func Test_ROR_Set7Bit(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x6A, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(true)
	cpu.registerA = 0b1001_0101
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
	assert.Equal(t, byte(0b1100_1010), cpu.registerA)
}

func Test_ROR_Unset7Bit(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x6A, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(false)
	cpu.registerA = 0b0000_0001
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.True(t, cpu.status.z())
	assert.False(t, cpu.status.n())
	assert.Equal(t, byte(0b0000_0000), cpu.registerA)
}

func Test_ROR_RotateFromMemory(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x76, 0x30, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x02
	cpu.status.setC(true)
	cpu.Bus.WriteMemory(0x32, 0b1001_0101)
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.True(t, cpu.status.z())
	assert.True(t, cpu.status.n())
	assert.Equal(t, byte(0b1100_1010), cpu.Bus.ReadMemory(0x32))
}

func Test_RTI_ReturnFromInterrupt(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x40, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x01_FF, 0x05)
	cpu.Bus.WriteMemory(0x01_FE, 0x06)
	cpu.Bus.WriteMemory(0x01_FD, 0b1001_0110)
	// SEC
	cpu.Bus.WriteMemory(0x05_07, 0x38)
	cpu.Bus.WriteMemory(0x05_08, 0x00)
	cpu.stackPointer = stackPointer(0xFC)
	cpu.Run()

	// SEC affected
	assert.Equal(t, byte(0b1001_0111), byte(cpu.status))
	// assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x05_09), cpu.ProgramCounter)
	assert.Equal(t, byte(0xFF), byte(cpu.stackPointer))
}

func Test_RTS_PopStack(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x60, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x01_FF, 0x05)
	cpu.Bus.WriteMemory(0x01_FE, 0x06)
	cpu.Bus.WriteMemory(0x05_07, 0xE8)
	cpu.Bus.WriteMemory(0x05_08, 0x00)
	cpu.stackPointer = stackPointer(0xFD)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x05_09), cpu.ProgramCounter)
	assert.Equal(t, byte(0xFF), byte(cpu.stackPointer))
}

func Test_SBC_Immediate(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xE9, 0b0111_1111, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b0111_1110
	cpu.Run()

	// 0b0111_1110 - 0b0111_1111
	// 0b0111_1110 + 0b1000_0000 + 1
	//
	//   0b0111_1110
	// + 0b1000_0001
	// -------------
	//   0b1111_1111 + (carry-1)

	assert.False(t, cpu.status.c())
	assert.True(t, cpu.status.n())
	assert.False(t, cpu.status.o())
	assert.False(t, cpu.status.z())
	assert.Equal(t, byte(0b1111_1110), cpu.registerA)
}

func Test_SBC_SetCarryFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xE9, 0b0000_0010, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b0000_0011
	cpu.Run()

	// 0b0000_0011 - 0b0000_0010
	// 0b0000_0011 + 0b1111_1101 + 1
	//
	//   0b0000_0011
	// + 0b1111_1110
	// -------------
	// 1 0b0000_0001 + (carry-1)

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.n())
	assert.False(t, cpu.status.o())
	assert.True(t, cpu.status.z())
	assert.Equal(t, byte(0b0000_0000), cpu.registerA)
}

func Test_SBC_SetCarryAndOverflowFlags(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xE9, 0b0111_1111, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b1011_0000
	cpu.Run()

	// 0b1011_0000 - 0b0111_1111
	// 0b1011_0000 + 0b1000_0000 + 1
	//
	//   0b1011_0000
	// + 0b1000_0001
	// -------------
	// 1 0b0011_0001 + (carry-1)

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.n())
	assert.True(t, cpu.status.o())
	assert.False(t, cpu.status.z())
	assert.Equal(t, byte(0b0011_0000), cpu.registerA)
}

func Test_JSRandRTS(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x20, 0x30, 0x40, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x40_30, 0xE8)
	cpu.Bus.WriteMemory(0x40_31, 0x60)
	cpu.Bus.WriteMemory(0x40_32, 0x00)
	cpu.Run()

	assert.Equal(t, uint16(0x80_04), cpu.ProgramCounter)
	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, byte(0xFF), byte(cpu.stackPointer))
}

func Test_TAX_MoveAtoX(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xAA, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0x10
	cpu.Run()

	assert.Equal(t, cpu.registerX, byte(0x10))
}

func Test_TAY_MoveAtoY(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xA8, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0x10
	cpu.Run()

	assert.Equal(t, cpu.registerY, byte(0x10))
}

func Test_TSX_MoveStoX(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xBA, 0x00})
	cpu.Reset(0x00_00)
	cpu.stackPointer = 0x10
	cpu.Run()

	assert.Equal(t, cpu.registerX, byte(0x10))
}

func Test_TXA_MoveXtoA(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x8A, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x10
	cpu.Run()

	assert.Equal(t, cpu.registerA, byte(0x10))
}

func Test_TXS_MoveXtoS(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x9A, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x10
	cpu.Run()

	assert.Equal(t, cpu.stackPointer, stackPointer(0x10))
}

func Test_TYA_MoveYtoA(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x98, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0x10
	cpu.Run()

	assert.Equal(t, cpu.registerA, byte(0x10))
}

func Test_CLC_UnsetCarryFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x18, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setC(true)
	assert.True(t, cpu.status.c())
	cpu.Run()

	assert.False(t, cpu.status.c())
}

func Test_CLD_UnsetDecimalFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xD8, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setD(true)
	assert.True(t, cpu.status.d())
	cpu.Run()

	assert.False(t, cpu.status.d())
}

func Test_CLI_UnsetInterruptDisableFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x58, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setI(true)
	assert.True(t, cpu.status.i())
	cpu.Run()

	assert.False(t, cpu.status.i())
}

func Test_CLV_UnsetOverflowFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xB8, 0x00})
	cpu.Reset(0x00_00)
	cpu.status.setO(true)
	assert.True(t, cpu.status.o())
	cpu.Run()

	assert.False(t, cpu.status.o())
}

func Test_CMP_SetZeroAndCarry(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xC9, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.True(t, cpu.status.z())
}

func Test_CMP_SetCarryOnly(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xC9, 0x09, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
}

func Test_CPX_SetZeroAndCarry(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xE0, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.True(t, cpu.status.z())
}

func Test_CPX_SetCarryOnly(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xE0, 0x09, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
}

func Test_CPY_SetZeroAndCarry(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xC0, 0x10, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.True(t, cpu.status.z())
}

func Test_CPY_SetCarryOnly(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xC0, 0x09, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0x10
	cpu.Run()

	assert.True(t, cpu.status.c())
	assert.False(t, cpu.status.z())
}

func Test_DEC_SetZeroFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xC6, 0x02, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x02, 0x01)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.Bus.ReadMemory(0x01))
	assert.True(t, cpu.status.z())
}

func Test_DEC_SetNegativeFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xC6, 0x01, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x01, 0b1000_0001)
	cpu.Run()

	assert.Equal(t, byte(0b1000_0000), cpu.Bus.ReadMemory(0x01))
	assert.True(t, cpu.status.n())
}

func Test_DEC_Decrement(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xC6, 0x02, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x02, 0x03)
	cpu.Run()

	assert.Equal(t, byte(0x02), cpu.Bus.ReadMemory(0x02))
}

func Test_DEC_Underflow(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xC6, 0x01, 0xC6, 0x01, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x01, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0xFE), cpu.Bus.ReadMemory(0x01))
}

func Test_DEX_Decrement(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xCA, 0xCA, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x03
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
}

func Test_DEX_UnderflowX(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xCA, 0xCA, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x00
	cpu.Run()

	assert.Equal(t, byte(0xFE), cpu.registerX)
}

func Test_DEY_Decrement(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x88, 0x88, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0x03
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerY)
}

func Test_EOR_Accumulator(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x49, 0b1001_0110, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b0000_1111
	cpu.Run()

	assert.Equal(t, byte(0b1001_1001), cpu.registerA)
}

func Test_EOR_SetZeroFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x49, 0b0000_0000, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b0000_0000
	cpu.Run()

	assert.True(t, cpu.status.z())
	assert.False(t, cpu.status.n())
}

func Test_EOR_SetNegativeFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x49, 0b0000_0000, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0b1000_0000
	cpu.Run()

	assert.False(t, cpu.status.z())
	assert.True(t, cpu.status.n())
}

func Test_DEY_UnderflowY(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x88, 0x88, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0x00
	cpu.Run()

	assert.Equal(t, byte(0xFE), cpu.registerY)
}

func Test_INC_SetZeroFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xE6, 0x01, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x01, 0xFF)
	cpu.Run()

	assert.Equal(t, byte(0x00), cpu.Bus.ReadMemory(0x01))
	assert.True(t, cpu.status.z())
}

func Test_INC_SetNegativeFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xE6, 0x01, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x01, 0b0111_1111)
	cpu.Run()

	assert.Equal(t, byte(0x80), cpu.Bus.ReadMemory(0x01))
	assert.True(t, cpu.status.n())
}

func Test_INC_Increment(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xE6, 0x01, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x01, 0x02)
	cpu.Run()

	assert.Equal(t, byte(0x03), cpu.Bus.ReadMemory(0x01))
}

func Test_INC_Overflow(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xE6, 0x01, 0xE6, 0x01, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x01, 0xFF)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.Bus.ReadMemory(0x01))
}

func Test_INX_Increment(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xE8, 0xE8, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x01
	cpu.Run()

	assert.Equal(t, byte(0x03), cpu.registerX)
}

func Test_INX_OverflowX(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xE8, 0xE8, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0xFF
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
}

func Test_INY_Increment(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xC8, 0xC8, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0x01
	cpu.Run()

	assert.Equal(t, byte(0x03), cpu.registerY)
}

func Test_INY_OverflowY(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xC8, 0xC8, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0xFF
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerY)
}

func Test_JMP_Absolute(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x4C, 0x30, 0x40, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x40_30, 0xE8)
	cpu.Bus.WriteMemory(0x40_31, 0x00)
	cpu.Run()

	assert.Equal(t, uint16(0x40_32), cpu.ProgramCounter)
	assert.Equal(t, byte(0x01), cpu.registerX)
}

func Test_JMP_Indirect(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x6C, 0x30, 0x40, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x40_30, 0x44)
	cpu.Bus.WriteMemory(0x40_31, 0x55)
	cpu.Bus.WriteMemory(0x55_44, 0xE8)
	cpu.Bus.WriteMemory(0x55_45, 0x00)
	cpu.Run()

	assert.Equal(t, uint16(0x55_46), cpu.ProgramCounter)
	assert.Equal(t, byte(0x01), cpu.registerX)
}

func Test_JSR_PushStack(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x20, 0x30, 0x40, 0x00})
	cpu.Reset(0x00_00)
	cpu.Bus.WriteMemory(0x40_30, 0xE8)
	cpu.Bus.WriteMemory(0x40_31, 0x00)
	cpu.Run()

	assert.Equal(t, uint16(0x40_32), cpu.ProgramCounter)
	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, byte(0x80), cpu.Bus.ReadMemory(0x01_FF))
	assert.Equal(t, byte(0x02), cpu.Bus.ReadMemory(0x01_FE))
}

func Test_STA_Immediate(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x85, 0x01, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerA = 0x05
	cpu.Run()

	assert.Equal(t, byte(0x05), cpu.Bus.ReadMemory(0x01))
}

func Test_SEC_SetCarryFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x38, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.True(t, cpu.status.c())
}

func Test_SED_SetDecimalFlag(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xF8, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.True(t, cpu.status.d())
}

func Test_SEI_SetInterruptDisable(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x78, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.True(t, cpu.status.i())
}

func Test_STX_Immediate(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x86, 0x01, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerX = 0x05
	cpu.Run()

	assert.Equal(t, byte(0x05), cpu.Bus.ReadMemory(0x01))
}

func Test_STY_Immediate(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0x84, 0x01, 0x00})
	cpu.Reset(0x00_00)
	cpu.registerY = 0x05
	cpu.Run()

	assert.Equal(t, byte(0x05), cpu.Bus.ReadMemory(0x01))
}

func Test_5OpsWorkingTogether(t *testing.T) {
	memory := memory.NewMemory()
	cpu := NewCPU(bus.NewBus(&memory))
	cpu.loadForTest([]byte{0xa9, 0xc0, 0xaa, 0xe8, 0x00})
	cpu.Reset(0x00_00)
	cpu.Run()

	assert.Equal(t, byte(0xC1), cpu.registerX)
}
