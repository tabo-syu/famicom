package rom

import (
	"errors"
	"slices"
)

type Mirroring int

const (
	Vertical Mirroring = iota
	Horizontal
	FourScreen
)

const (
	PrgROMPageSize = 16_384
	ChrROMPageSize = 8_192
)

type ROM struct {
	Prg             []byte
	Chr             []byte
	Mapper          byte
	ScreenMirroring Mirroring
}

func NewROM(raw []byte) (*ROM, error) {
	if !slices.Equal(raw[0:4], []byte{'N', 'E', 'S', 0x1A}) {
		return nil, errors.New("file is not iNES file foramt")
	}

	inesVer := (raw[7] >> 2) & 0b0000_0011
	if inesVer != 0x00 {
		return nil, errors.New("NES2.0 format is not supported")
	}

	// [prg|chr]romsize
	prgROMSize := uint(raw[4]) * PrgROMPageSize
	chrROMSize := uint(raw[5]) * ChrROMPageSize

	skipTrainer := raw[6]&0b0000_0100 != 0

	var prgROMStartPos uint = 16
	if skipTrainer {
		prgROMStartPos += 512
	}

	chrROMStartPos := prgROMStartPos + prgROMSize

	// mapper
	mapper := (raw[7] & 0b1111_0000) | (raw[6] >> 4)

	// screenMirroring
	isFourScreen := raw[6]&0b0000_1000 != 0
	isVerticalMirroring := raw[6]&0b0000_0001 != 0

	var screenMirroring Mirroring
	if isFourScreen {
		screenMirroring = FourScreen
	} else if !isFourScreen && isVerticalMirroring {
		screenMirroring = Vertical
	} else if !isFourScreen && !isVerticalMirroring {
		screenMirroring = Horizontal
	}

	return &ROM{
		Prg:             raw[prgROMStartPos : prgROMStartPos+prgROMSize],
		Chr:             raw[chrROMStartPos : chrROMStartPos+chrROMSize],
		Mapper:          mapper,
		ScreenMirroring: screenMirroring,
	}, nil
}
