package bus

import (
	"fmt"
	"log"

	"github.com/tabo-syu/famicom/internal/memory"
	"github.com/tabo-syu/famicom/internal/rom"
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
	ReadPrgROM(address uint16) byte
}

type bus struct {
	Memory memory.Memory
	ROM    *rom.ROM
}

func NewBus(memory memory.Memory, rom *rom.ROM) *bus {
	return &bus{
		Memory: memory,
		ROM:    rom,
	}
}

func (bus *bus) ReadMemory(address uint16) byte {
	if 0x8000 <= address && address <= 0xFFFF {
		return bus.ReadPrgROM(address)
	}

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
	if 0x8000 <= address && address <= 0xFFFF {
		panic("Attempt to write to Cartridge ROM space")
	}

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

func (bus *bus) ReadPrgROM(address uint16) byte {
	address -= 0x8000
	if len(bus.ROM.Prg) == 0x4000 && address >= 0x4000 {
		address = address % 0x4000
	}

	return bus.ROM.Prg[address]
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
