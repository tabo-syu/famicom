package cpu

type CPU struct {
	ProgramCounter uint16
	RegisterA      uint8
	RegisterX      uint8
	RegisterY      uint8
	Status         Status

	memory [0x1_00_00]uint8
}

type AddressingMode int

const (
	AddressingMode_Immediate AddressingMode = iota
	AddressingMode_ZeroPage
	AddressingMode_ZeroPageX
	AddressingMode_ZeroPageY
	AddressingMode_Absolute
	AddressingMode_AbsoluteX
	AddressingMode_AbsoluteY
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
	}
}

func (cpu *CPU) readMemory(address uint16) uint8 {
	return cpu.memory[address]
}

func (cpu *CPU) writeMemory(address uint16, data uint8) {
	cpu.memory[address] = data
}

func (cpu *CPU) readMemoryUint16(position uint16) uint16 {
	low := uint16(cpu.readMemory(position))
	high := uint16(cpu.readMemory(position + 1))

	return high<<8 | low
}

func (cpu *CPU) writeMemoryUint16(position uint16, data uint16) {
	high := uint8(data >> 8)
	low := uint8(data & 0x00_FF)

	cpu.writeMemory(position, low)
	cpu.writeMemory(position+1, high)
}

func (cpu *CPU) Load(program []uint8) {
	copy(cpu.memory[0x8000:], program)
	cpu.writeMemoryUint16(0xFF_FC, 0x80_00)
}

func (cpu *CPU) Reset() {
	cpu.RegisterA = 0
	cpu.RegisterX = 0
	cpu.Status = NewStatus()

	// 0x80_00
	cpu.ProgramCounter = cpu.readMemoryUint16(0xFF_FC)
}

func (cpu *CPU) getOperandAddress(mode AddressingMode) uint16 {
	switch mode {
	case AddressingMode_Immediate:
		return cpu.ProgramCounter

	case AddressingMode_ZeroPage:
		return uint16(cpu.readMemory(cpu.ProgramCounter))

	case AddressingMode_ZeroPageX:
		position := cpu.readMemory(cpu.ProgramCounter)
		address := uint16(position + cpu.RegisterX)

		return address

	case AddressingMode_ZeroPageY:
		position := cpu.readMemory(cpu.ProgramCounter)
		address := uint16(position + cpu.RegisterY)

		return address

	case AddressingMode_Absolute:
		return cpu.readMemoryUint16(cpu.ProgramCounter)

	case AddressingMode_AbsoluteX:
		base := cpu.readMemoryUint16(cpu.ProgramCounter)
		address := base + uint16(cpu.RegisterX)

		return address

	case AddressingMode_AbsoluteY:
		base := cpu.readMemoryUint16(cpu.ProgramCounter)
		address := base + uint16(cpu.RegisterY)

		return address

	case AddressingMode_IndirectX:
		base := cpu.readMemory(cpu.ProgramCounter)
		pointer := base + cpu.RegisterX
		low := cpu.readMemory(uint16(pointer))
		high := cpu.readMemory(uint16(pointer + 1))
		address := uint16(high)<<8 | uint16(low)

		return address

	case AddressingMode_IndirectY:
		base := cpu.readMemory(cpu.ProgramCounter)
		low := cpu.readMemory(uint16(base))
		high := cpu.readMemory(uint16(base + 1))
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
		opscode := cpu.readMemory(cpu.ProgramCounter)
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
	cpu.writeMemory(address, cpu.RegisterA)
}

func (cpu *CPU) lda(mode AddressingMode) {
	address := cpu.getOperandAddress(mode)
	value := cpu.readMemory(address)

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
