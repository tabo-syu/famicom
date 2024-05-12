package cpu

import "github.com/tabo-syu/famicom/internal/memory"

type CPU struct {
	programCounter uint16
	registerA      uint8
	registerX      uint8
	registerY      uint8
	status         status

	memory memory.Memory
}

type addressingMode int

const (
	addressingMode_Implied addressingMode = iota
	addressingMode_Accumulator
	addressingMode_Immediate
	addressingMode_ZeroPage
	addressingMode_ZeroPageX
	addressingMode_ZeroPageY
	addressingMode_Relative
	addressingMode_Absolute
	addressingMode_AbsoluteX
	addressingMode_AbsoluteY
	addressingMode_Indirect
	addressingMode_IndirectX
	addressingMode_IndirectY
	addressingMode_NoneAddressing
)

func NewCPU() CPU {
	return CPU{
		programCounter: 0,
		registerA:      0,
		registerX:      0,
		status:         newStatus(),

		memory: memory.NewMemory(),
	}
}

func (cpu *CPU) Load(program []uint8) {
	copy(cpu.memory[0x8000:], program)
	cpu.memory.WriteUint16(0xFF_FC, 0x80_00)
}

func (cpu *CPU) Reset() {
	cpu.registerA = 0
	cpu.registerX = 0
	cpu.status = newStatus()

	// 0x80_00
	cpu.programCounter = cpu.memory.ReadUint16(0xFF_FC)
}

func (cpu *CPU) getOperandAddress(mode addressingMode) uint16 {
	switch mode {
	case addressingMode_Immediate:
		return cpu.programCounter

	case addressingMode_ZeroPage:
		return uint16(cpu.memory.Read(cpu.programCounter))

	case addressingMode_ZeroPageX:
		position := cpu.memory.Read(cpu.programCounter)
		address := uint16(position + cpu.registerX)

		return address

	case addressingMode_ZeroPageY:
		position := cpu.memory.Read(cpu.programCounter)
		address := uint16(position + cpu.registerY)

		return address

	case addressingMode_Absolute:
		return cpu.memory.ReadUint16(cpu.programCounter)

	case addressingMode_AbsoluteX:
		base := cpu.memory.ReadUint16(cpu.programCounter)
		address := base + uint16(cpu.registerX)

		return address

	case addressingMode_AbsoluteY:
		base := cpu.memory.ReadUint16(cpu.programCounter)
		address := base + uint16(cpu.registerY)

		return address

	case addressingMode_IndirectX:
		base := cpu.memory.Read(cpu.programCounter)
		pointer := base + cpu.registerX
		low := cpu.memory.Read(uint16(pointer))
		high := cpu.memory.Read(uint16(pointer + 1))
		address := uint16(high)<<8 | uint16(low)

		return address

	case addressingMode_IndirectY:
		base := cpu.memory.Read(cpu.programCounter)
		low := cpu.memory.Read(uint16(base))
		high := cpu.memory.Read(uint16(base + 1))
		derefBase := uint16(high)<<8 | uint16(low)
		deref := derefBase + uint16(cpu.registerY)

		return deref

	case addressingMode_NoneAddressing:
		panic("through `NoneAddressing`")

	default:
		return 0
	}
}

func (cpu *CPU) Run() {
	for {
		opscode := cpu.memory.Read(cpu.programCounter)
		cpu.programCounter++

		switch opscode {
		// STA(ZeroPage)
		case 0x85:
			cpu.sta(addressingMode_ZeroPage)
			cpu.programCounter++

		// STA(ZeroPageX)
		case 0x95:
			cpu.sta(addressingMode_ZeroPageX)
			cpu.programCounter++

		// STA(Absolute)
		case 0x8D:
			cpu.sta(addressingMode_Absolute)
			cpu.programCounter += 2

		// STA(AbsoluteX)
		case 0x9D:
			cpu.sta(addressingMode_AbsoluteX)
			cpu.programCounter += 2

		// STA(AbsoluteY)
		case 0x99:
			cpu.sta(addressingMode_AbsoluteY)
			cpu.programCounter += 2

		// STA(IndirectX)
		case 0x81:
			cpu.sta(addressingMode_IndirectX)
			cpu.programCounter++

		// STA(IndirectY)
		case 0x91:
			cpu.sta(addressingMode_IndirectY)
			cpu.programCounter++

		// LDA(Immediate)
		case 0xA9:
			cpu.lda(addressingMode_Immediate)
			cpu.programCounter++

		// LDA(ZeroPage)
		case 0xA5:
			cpu.lda(addressingMode_ZeroPage)
			cpu.programCounter++

		// LDA(ZeroPageX)
		case 0xB5:
			cpu.lda(addressingMode_ZeroPageX)
			cpu.programCounter++

		// LDA(Absolute)
		case 0xAD:
			cpu.lda(addressingMode_Absolute)
			cpu.programCounter += 2

		// LDA(AbsoluteX)
		case 0xBD:
			cpu.lda(addressingMode_AbsoluteX)
			cpu.programCounter += 2

		// LDA(AbsoluteY)
		case 0xB9:
			cpu.lda(addressingMode_AbsoluteY)
			cpu.programCounter += 2

		// LDA(IndirectX)
		case 0xA1:
			cpu.lda(addressingMode_IndirectX)
			cpu.programCounter++

		// LDA(IndirectY)
		case 0xB1:
			cpu.lda(addressingMode_IndirectY)
			cpu.programCounter++

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

func (cpu *CPU) sta(mode addressingMode) {
	address := cpu.getOperandAddress(mode)
	cpu.memory.Write(address, cpu.registerA)
}

func (cpu *CPU) lda(mode addressingMode) {
	address := cpu.getOperandAddress(mode)
	value := cpu.memory.Read(address)

	cpu.registerA = value
	cpu.updateZeroAndNegativeFlags(cpu.registerA)
}

func (cpu *CPU) tax() {
	cpu.registerX = cpu.registerA
	cpu.updateZeroAndNegativeFlags(cpu.registerA)
}

func (cpu *CPU) inx() {
	cpu.registerX++
	cpu.updateZeroAndNegativeFlags(cpu.registerX)
}

func (cpu *CPU) updateZeroAndNegativeFlags(register uint8) {
	if register == 0 {
		cpu.status.setZ(true)
	} else {
		cpu.status.setZ(false)
	}

	if register&0b0100_0000 != 0 {
		cpu.status.setN(true)
	} else {
		cpu.status.setN(false)
	}
}
