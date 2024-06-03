package cpu

import (
	"log"
	"time"

	"github.com/tabo-syu/famicom/internal/bus"
)

type CPU struct {
	ProgramCounter uint16
	registerA      byte
	registerX      byte
	registerY      byte
	stackPointer   stackPointer
	status         status

	Bus          bus.Bus
	Instructions map[byte]instruction
}

func NewCPU(bus bus.Bus) CPU {
	return CPU{
		ProgramCounter: 0,
		registerA:      0,
		registerX:      0,
		registerY:      0,
		stackPointer:   0,
		status:         0,

		Bus:          bus,
		Instructions: NewInstructions(),
	}
}

func (cpu *CPU) Load(program []byte) {
	cpu.Bus.CopyToMemory(0x06_00, program)
	cpu.Bus.WriteMemoryUint16(0xFF_FC, 0x06_00)
}

func (cpu *CPU) Reset(position uint16) {
	// 0x80_00
	cpu.ProgramCounter = cpu.Bus.ReadMemoryUint16(position)
	cpu.registerA = 0
	cpu.registerX = 0
	cpu.registerY = 0
	cpu.stackPointer = newStackPointer()
	cpu.status = newStatus()
}

func (cpu *CPU) Run() {
	for {
		code := cpu.Bus.ReadMemory(cpu.ProgramCounter)
		cpu.ProgramCounter++

		if err := cpu.Instructions[code].Call(cpu); err != nil {
			log.Println(err)

			break
		}

		time.Sleep(10 * time.Microsecond)
	}
}

func (cpu *CPU) LoadAndRun(program []byte) {
	cpu.Load(program)
	cpu.Reset(0xFF_FC)
	cpu.Run()
}
