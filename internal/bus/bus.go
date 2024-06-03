package bus

import (
	"fmt"
	"log"

	"github.com/tabo-syu/famicom/internal/memory"
)

const (
	RAM                    uint16 = 0x00_00
	RAMMirrorsEnd          uint16 = 0x1F_FF
	PPURegisters           uint16 = 0x20_00
	PPURegistersMirrorsEnd uint16 = 0x3F_FF
)

type Bus interface {
	ReadMemory(address uint16) byte
	ReadMemoryUint16(address uint16) uint16
	WriteMemory(address uint16, data byte)
	WriteMemoryUint16(address uint16, data uint16)
	CopyToMemory(start int, value []byte)
}

type bus struct {
	Memory memory.Memory
}

func NewBus(memory memory.Memory) *bus {
	return &bus{
		Memory: memory,
	}
}

func (bus *bus) ReadMemory(address uint16) byte {
	masked, err := bus.mask(address)
	if err != nil {
		log.Println(err)

		return 0x00
	}

	return bus.Memory.Read(masked)
}

func (bus *bus) ReadMemoryUint16(address uint16) uint16 {
	masked, err := bus.mask(address)
	if err != nil {
		log.Println(err)

		return 0x00
	}

	return bus.Memory.ReadUint16(masked)
}

func (bus *bus) WriteMemory(address uint16, data byte) {
	masked, err := bus.mask(address)
	if err != nil {
		log.Println(err)

		return
	}

	bus.Memory.Write(masked, data)
}

func (bus *bus) WriteMemoryUint16(address uint16, data uint16) {
	masked, err := bus.mask(address)
	if err != nil {
		log.Println(err)

		return
	}

	bus.Memory.WriteUint16(masked, data)
}

func (bus *bus) CopyToMemory(start int, value []byte) {
	bus.Memory.Copy(start, value)
}

func (bus *bus) mask(address uint16) (uint16, error) {
	var (
		masked uint16
		err    error
	)

	switch {
	case RAM <= address && address <= RAMMirrorsEnd:
		masked = address & 0b0000_0111_1111_1111
	case PPURegisters <= address && address <= PPURegistersMirrorsEnd:
	default:
		err = fmt.Errorf("ignoring memory access at %#x", address)
	}

	return masked, err
}
