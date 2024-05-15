package cpu

import (
	"errors"
)

type addressingMode int

const (
	ImpliedMode addressingMode = iota
	AccumulatorMode
	ImmediateMode
	ZeroPageMode
	ZeroPageXMode
	ZeroPageYMode
	RelativeMode
	AbsoluteMode
	AbsoluteXMode
	AbsoluteYMode
	IndirectMode
	IndirectXMode
	IndirectYMode
	NoneAddressingMode
)

func (cpu *CPU) BRK(mode addressingMode) error {
	return errors.New("BRK called")
}

func (cpu *CPU) AND(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)
	value := cpu.memory.Read(address)

	result := cpu.registerA & value
	cpu.registerA = result
	cpu.updateZeroAndNegativeFlags(result)

	return nil
}

func (cpu *CPU) SEC(mode addressingMode) error {
	cpu.status.setC(true)

	return nil
}

func (cpu *CPU) SED(mode addressingMode) error {
	cpu.status.setD(true)

	return nil
}

func (cpu *CPU) SEI(mode addressingMode) error {
	cpu.status.setI(true)

	return nil
}

func (cpu *CPU) STA(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)
	cpu.memory.Write(address, cpu.registerA)

	return nil
}

func (cpu *CPU) STX(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)
	cpu.memory.Write(address, cpu.registerX)

	return nil
}

func (cpu *CPU) STY(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)
	cpu.memory.Write(address, cpu.registerY)

	return nil
}

func (cpu *CPU) LDA(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)
	value := cpu.memory.Read(address)

	cpu.registerA = value
	cpu.updateZeroAndNegativeFlags(cpu.registerA)

	return nil
}

func (cpu *CPU) LDX(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)
	value := cpu.memory.Read(address)

	cpu.registerX = value
	cpu.updateZeroAndNegativeFlags(cpu.registerX)

	return nil
}

func (cpu *CPU) LDY(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)
	value := cpu.memory.Read(address)

	cpu.registerY = value
	cpu.updateZeroAndNegativeFlags(cpu.registerY)

	return nil
}

func (cpu *CPU) ORA(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)
	value := cpu.memory.Read(address)

	result := cpu.registerA | value
	cpu.registerA = result
	cpu.updateZeroAndNegativeFlags(result)

	return nil
}

func (cpu *CPU) TAX(mode addressingMode) error {
	cpu.registerX = cpu.registerA
	cpu.updateZeroAndNegativeFlags(cpu.registerX)

	return nil
}

func (cpu *CPU) TAY(mode addressingMode) error {
	cpu.registerY = cpu.registerA
	cpu.updateZeroAndNegativeFlags(cpu.registerY)

	return nil
}

func (cpu *CPU) TSX(mode addressingMode) error {
	cpu.registerX = cpu.stackPointer
	cpu.updateZeroAndNegativeFlags(cpu.registerX)

	return nil
}

func (cpu *CPU) TXA(mode addressingMode) error {
	cpu.registerA = cpu.registerX
	cpu.updateZeroAndNegativeFlags(cpu.registerA)

	return nil
}

func (cpu *CPU) TXS(mode addressingMode) error {
	cpu.stackPointer = cpu.registerX

	return nil
}

func (cpu *CPU) TYA(mode addressingMode) error {
	cpu.registerA = cpu.registerY
	cpu.updateZeroAndNegativeFlags(cpu.registerA)

	return nil
}

func (cpu *CPU) CLC(mode addressingMode) error {
	cpu.status.setC(false)

	return nil
}

func (cpu *CPU) CLD(mode addressingMode) error {
	cpu.status.setD(false)

	return nil
}

func (cpu *CPU) CLI(mode addressingMode) error {
	cpu.status.setI(false)

	return nil
}

func (cpu *CPU) CLV(mode addressingMode) error {
	cpu.status.setO(false)

	return nil
}

func (cpu *CPU) CMP(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)

	value := cpu.memory.Read(address)
	if cpu.registerA == value {
		cpu.status.setZ(true)
	}
	if cpu.registerA >= value {
		cpu.status.setC(true)
	}
	if value&0b0100_0000 != 0 {
		cpu.status.setN(true)
	} else {
		cpu.status.setN(false)
	}

	return nil
}

func (cpu *CPU) CPX(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)

	value := cpu.memory.Read(address)
	if cpu.registerX == value {
		cpu.status.setZ(true)
	}
	if cpu.registerX >= value {
		cpu.status.setC(true)
	}
	if value&0b0100_0000 != 0 {
		cpu.status.setN(true)
	} else {
		cpu.status.setN(false)
	}

	return nil
}

func (cpu *CPU) CPY(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)

	value := cpu.memory.Read(address)
	if cpu.registerY == value {
		cpu.status.setZ(true)
	}
	if cpu.registerY >= value {
		cpu.status.setC(true)
	}
	if value&0b0100_0000 != 0 {
		cpu.status.setN(true)
	} else {
		cpu.status.setN(false)
	}

	return nil
}

func (cpu *CPU) DEC(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)

	value := cpu.memory.Read(address)
	value--
	cpu.memory.Write(address, value)

	cpu.updateZeroAndNegativeFlags(value)

	return nil
}

func (cpu *CPU) DEX(mode addressingMode) error {
	cpu.registerX--
	cpu.updateZeroAndNegativeFlags(cpu.registerX)

	return nil
}

func (cpu *CPU) DEY(mode addressingMode) error {
	cpu.registerY--
	cpu.updateZeroAndNegativeFlags(cpu.registerY)

	return nil
}

func (cpu *CPU) EOR(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)
	value := cpu.memory.Read(address)

	result := cpu.registerA ^ value
	cpu.registerA = result
	cpu.updateZeroAndNegativeFlags(result)

	return nil
}

func (cpu *CPU) INC(mode addressingMode) error {
	address := cpu.getOperandAddress(mode)

	value := cpu.memory.Read(address)
	value++
	cpu.memory.Write(address, value)

	cpu.updateZeroAndNegativeFlags(value)

	return nil
}

func (cpu *CPU) INX(mode addressingMode) error {
	cpu.registerX++
	cpu.updateZeroAndNegativeFlags(cpu.registerX)

	return nil
}

func (cpu *CPU) INY(mode addressingMode) error {
	cpu.registerY++
	cpu.updateZeroAndNegativeFlags(cpu.registerY)

	return nil
}

func (cpu *CPU) getOperandAddress(mode addressingMode) uint16 {
	switch mode {
	case ImmediateMode:
		return cpu.programCounter

	case ZeroPageMode:
		return uint16(cpu.memory.Read(cpu.programCounter))

	case ZeroPageXMode:
		position := cpu.memory.Read(cpu.programCounter)
		address := uint16(position + cpu.registerX)

		return address

	case ZeroPageYMode:
		position := cpu.memory.Read(cpu.programCounter)
		address := uint16(position + cpu.registerY)

		return address

	case AbsoluteMode:
		return cpu.memory.ReadUint16(cpu.programCounter)

	case AbsoluteXMode:
		base := cpu.memory.ReadUint16(cpu.programCounter)
		address := base + uint16(cpu.registerX)

		return address

	case AbsoluteYMode:
		base := cpu.memory.ReadUint16(cpu.programCounter)
		address := base + uint16(cpu.registerY)

		return address

	case IndirectXMode:
		base := cpu.memory.Read(cpu.programCounter)
		pointer := base + cpu.registerX
		low := cpu.memory.Read(uint16(pointer))
		high := cpu.memory.Read(uint16(pointer + 1))
		address := uint16(high)<<8 | uint16(low)

		return address

	case IndirectYMode:
		base := cpu.memory.Read(cpu.programCounter)
		low := cpu.memory.Read(uint16(base))
		high := cpu.memory.Read(uint16(base + 1))
		derefBase := uint16(high)<<8 | uint16(low)
		deref := derefBase + uint16(cpu.registerY)

		return deref

	case NoneAddressingMode:
		panic("through `NoneAddressing`")

	default:
		return 0
	}
}

func (cpu *CPU) updateZeroFlag(value uint8) {
	if value == 0 {
		cpu.status.setZ(true)
	} else {
		cpu.status.setZ(false)
	}
}

func (cpu *CPU) updateNegativeFlag(value uint8) {
	if value&0b0100_0000 != 0 {
		cpu.status.setN(true)
	} else {
		cpu.status.setN(false)
	}
}

func (cpu *CPU) updateZeroAndNegativeFlags(value uint8) {
	cpu.updateZeroFlag(value)
	cpu.updateNegativeFlag(value)
}
