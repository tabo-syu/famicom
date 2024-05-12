package cpu

import "github.com/tabo-syu/famicom/internal/memory"

type CPU struct {
	ProgramCounter uint16
	RegisterA      uint8
	RegisterX      uint8
	RegisterY      uint8
	Status         Status

	memory memory.Memory
}

type AddressingMode int

const (
	AddressingMode_Implied AddressingMode = iota
	AddressingMode_Accumulator
	AddressingMode_Immediate
	AddressingMode_ZeroPage
	AddressingMode_ZeroPageX
	AddressingMode_ZeroPageY
	AddressingMode_Relative
	AddressingMode_Absolute
	AddressingMode_AbsoluteX
	AddressingMode_AbsoluteY
	AddressingMode_Indirect
	AddressingMode_IndirectX
	AddressingMode_IndirectY
	AddressingMode_NoneAddressing
)

func NewCPU() CPU {
	return CPU{
		ProgramCounter: 0,
		RegisterA:      0,
		RegisterX:      0,
		Status:         NewStatus(),
		memory:         memory.NewMemory(),
	}
}

func (cpu *CPU) Load(program []uint8) {
	copy(cpu.memory[0x8000:], program)
	cpu.memory.WriteUint16(0xFF_FC, 0x80_00)
}

func (cpu *CPU) Reset() {
	cpu.RegisterA = 0
	cpu.RegisterX = 0
	cpu.Status = NewStatus()

	// 0x80_00
	cpu.ProgramCounter = cpu.memory.ReadUint16(0xFF_FC)
}

func (cpu *CPU) getOperandAddress(mode AddressingMode) uint16 {
	switch mode {
	case AddressingMode_Immediate:
		return cpu.ProgramCounter

	case AddressingMode_ZeroPage:
		return uint16(cpu.memory.Read(cpu.ProgramCounter))

	case AddressingMode_ZeroPageX:
		position := cpu.memory.Read(cpu.ProgramCounter)
		address := uint16(position + cpu.RegisterX)

		return address

	case AddressingMode_ZeroPageY:
		position := cpu.memory.Read(cpu.ProgramCounter)
		address := uint16(position + cpu.RegisterY)

		return address

	case AddressingMode_Absolute:
		return cpu.memory.ReadUint16(cpu.ProgramCounter)

	case AddressingMode_AbsoluteX:
		base := cpu.memory.ReadUint16(cpu.ProgramCounter)
		address := base + uint16(cpu.RegisterX)

		return address

	case AddressingMode_AbsoluteY:
		base := cpu.memory.ReadUint16(cpu.ProgramCounter)
		address := base + uint16(cpu.RegisterY)

		return address

	case AddressingMode_IndirectX:
		base := cpu.memory.Read(cpu.ProgramCounter)
		pointer := base + cpu.RegisterX
		low := cpu.memory.Read(uint16(pointer))
		high := cpu.memory.Read(uint16(pointer + 1))
		address := uint16(high)<<8 | uint16(low)

		return address

	case AddressingMode_IndirectY:
		base := cpu.memory.Read(cpu.ProgramCounter)
		low := cpu.memory.Read(uint16(base))
		high := cpu.memory.Read(uint16(base + 1))
		derefBase := uint16(high)<<8 | uint16(low)
		deref := derefBase + uint16(cpu.RegisterY)

		return deref

	case AddressingMode_NoneAddressing:
		panic("through `NoneAddressing`")

	default:
		return 0
	}
}

func (cpu *CPU) Run() {
	for {
		opscode := cpu.memory.Read(cpu.ProgramCounter)
		cpu.ProgramCounter++

		switch opscode {
		// STA(ZeroPage)
		case 0x85:
			cpu.sta(AddressingMode_ZeroPage)
			cpu.ProgramCounter++

		// STA(ZeroPageX)
		case 0x95:
			cpu.sta(AddressingMode_ZeroPageX)
			cpu.ProgramCounter++

		// STA(Absolute)
		case 0x8D:
			cpu.sta(AddressingMode_Absolute)
			cpu.ProgramCounter += 2

		// STA(AbsoluteX)
		case 0x9D:
			cpu.sta(AddressingMode_AbsoluteX)
			cpu.ProgramCounter += 2

		// STA(AbsoluteY)
		case 0x99:
			cpu.sta(AddressingMode_AbsoluteY)
			cpu.ProgramCounter += 2

		// STA(IndirectX)
		case 0x81:
			cpu.sta(AddressingMode_IndirectX)
			cpu.ProgramCounter++

		// STA(IndirectY)
		case 0x91:
			cpu.sta(AddressingMode_IndirectY)
			cpu.ProgramCounter++

		// LDA(Immediate)
		case 0xA9:
			cpu.lda(AddressingMode_Immediate)
			cpu.ProgramCounter++

		// LDA(ZeroPage)
		case 0xA5:
			cpu.lda(AddressingMode_ZeroPage)
			cpu.ProgramCounter++

		// LDA(ZeroPageX)
		case 0xB5:
			cpu.lda(AddressingMode_ZeroPageX)
			cpu.ProgramCounter++

		// LDA(Absolute)
		case 0xAD:
			cpu.lda(AddressingMode_Absolute)
			cpu.ProgramCounter += 2

		// LDA(AbsoluteX)
		case 0xBD:
			cpu.lda(AddressingMode_AbsoluteX)
			cpu.ProgramCounter += 2

		// LDA(AbsoluteY)
		case 0xB9:
			cpu.lda(AddressingMode_AbsoluteY)
			cpu.ProgramCounter += 2

		// LDA(IndirectX)
		case 0xA1:
			cpu.lda(AddressingMode_IndirectX)
			cpu.ProgramCounter++

		// LDA(IndirectY)
		case 0xB1:
			cpu.lda(AddressingMode_IndirectY)
			cpu.ProgramCounter++

		// TAX
		case 0xAA:
			cpu.tax()

		// INX
		case 0xE8:
			cpu.inx()

		// BRK
		case 0x00:
			return

		default:
		}
	}
}

func (cpu *CPU) loadAndRun(program []uint8) {
	cpu.Load(program)
	cpu.Reset()
	cpu.Run()
}

func (cpu *CPU) sta(mode AddressingMode) {
	address := cpu.getOperandAddress(mode)
	cpu.memory.Write(address, cpu.RegisterA)
}

func (cpu *CPU) lda(mode AddressingMode) {
	address := cpu.getOperandAddress(mode)
	value := cpu.memory.Read(address)

	cpu.RegisterA = value
	cpu.updateZeroAndNegativeFlags(cpu.RegisterA)
}

func (cpu *CPU) tax() {
	cpu.RegisterX = cpu.RegisterA
	cpu.updateZeroAndNegativeFlags(cpu.RegisterA)
}

func (cpu *CPU) inx() {
	cpu.RegisterX++
	cpu.updateZeroAndNegativeFlags(cpu.RegisterX)
}

func (cpu *CPU) updateZeroAndNegativeFlags(register uint8) {
	if register == 0 {
		cpu.Status.SetZ(true)
	} else {
		cpu.Status.SetZ(false)
	}

	if register&0b0100_0000 != 0 {
		cpu.Status.SetN(true)
	} else {
		cpu.Status.SetN(false)
	}
}
