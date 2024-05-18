package cpu

import (
	"log"

	"github.com/tabo-syu/famicom/internal/memory"
)

type CPU struct {
	programCounter uint16
	registerA      uint8
	registerX      uint8
	registerY      uint8
	stackPointer   stackPointer
	status         status

	memory       memory.Memory
	instructions map[uint8]instruction
}

func NewCPU() CPU {
	return CPU{
		programCounter: 0,
		registerA:      0,
		registerX:      0,
		registerY:      0,
		stackPointer:   0,
		status:         0,

		memory:       memory.NewMemory(),
		instructions: NewInstructions(),
	}
}

func (cpu *CPU) Load(program []uint8) {
	copy(cpu.memory[0x8000:], program)
	cpu.memory.WriteUint16(0xFF_FC, 0x80_00)
}

func (cpu *CPU) Reset() {
	// 0x80_00
	cpu.programCounter = cpu.memory.ReadUint16(0xFF_FC)
	cpu.registerA = 0
	cpu.registerX = 0
	cpu.registerY = 0
	cpu.stackPointer = newStackPointer()
	cpu.status = newStatus()
}

func (cpu *CPU) Run() {
	for {
		code := cpu.memory.Read(cpu.programCounter)
		cpu.programCounter++

		if err := cpu.instructions[code].Call(cpu); err != nil {
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
