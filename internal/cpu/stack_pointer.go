package cpu

type stackPointer byte

func newStackPointer() stackPointer {
	return 0xFF
}

func (s *stackPointer) toAddress() uint16 {
	return 0x01_00 + uint16(*s)
}
