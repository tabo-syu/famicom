package cpu

type CPU struct {
	ProgramCounter uint16
	RegisterA      uint8
	RegisterX      uint8
	Status         Status
}

func NewCPU() CPU {
	return CPU{
		ProgramCounter: 0b0000_0000,
		RegisterA:      0b0000_0000,
		RegisterX:      0b0000_0000,
		Status:         NewStatus(),
	}
}

func (cpu *CPU) Interpret(program []uint8) {
	cpu.ProgramCounter = 0

	for {
		opscode := program[cpu.ProgramCounter]
		cpu.ProgramCounter++

		switch opscode {

		// LDA
		case 0xA9:
			arg := program[cpu.ProgramCounter]
			cpu.ProgramCounter++

			cpu.lda(arg)

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

func (cpu *CPU) lda(param uint8) {
	cpu.RegisterA = param
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

	if register&0b1000_000 != 0 {
		cpu.Status.SetN(true)
	} else {
		cpu.Status.SetN(false)
	}
}
