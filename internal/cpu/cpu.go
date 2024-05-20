package cpu

import (
	"log"

	"github.com/tabo-syu/famicom/internal/memory"
)

type CPU struct {
	ProgramCounter uint16
	registerA      uint8
	registerX      uint8
	registerY      uint8
	stackPointer   stackPointer
	status         status

	Memory       memory.Memory
	Instructions map[uint8]instruction
}

func NewCPU() CPU {
	return CPU{
		ProgramCounter: 0,
		registerA:      0,
		registerX:      0,
		registerY:      0,
		stackPointer:   0,
		status:         0,

		Memory:       memory.NewMemory(),
		Instructions: NewInstructions(),
	}
}

func (cpu *CPU) Load(program []uint8) {
	copy(cpu.Memory[0x06_00:], program)
	cpu.Memory.WriteUint16(0xFF_FC, 0x06_00)
}

func (cpu *CPU) Reset() {
	// 0x80_00
	cpu.ProgramCounter = cpu.Memory.ReadUint16(0xFF_FC)
	cpu.registerA = 0
	cpu.registerX = 0
	cpu.registerY = 0
	cpu.stackPointer = newStackPointer()
	cpu.status = newStatus()
}

func (cpu *CPU) Run() {
	cpu.RunWithCallback(func(cpu *CPU) {})
}

func (cpu *CPU) RunWithCallback(callback func(*CPU)) {
	for {
		callback(cpu)

		code := cpu.Memory.Read(cpu.ProgramCounter)
		cpu.ProgramCounter++

		if err := cpu.Instructions[code].Call(cpu); err != nil {
			log.Println(err)

			break
		}
	}
}

func (cpu *CPU) loadAndRun(program []uint8) {
	cpu.Load(program)
	cpu.Reset()
	cpu.Run()
}
