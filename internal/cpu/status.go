package cpu

/*
	`0bNV1B_DIZC`
	- N: Negative
	- V: Overflow
	- 1: (No CPU effect; always pushed as 1)
	- B: (No CPU effect; the B flag)
	- D: Decimal
	- I: Interrupt Disable
	- Z: Zero
	- C: Carry

https://www.nesdev.org/wiki/Status_flags
*/
type Status uint8

func NewStatus() Status {
	return 0b0000_0000
}

// Get Carry flag
func (s *Status) C() bool {
	return (*s & 0b0000_0001) == 0b0000_0001
}

// Set Carry flag
func (s *Status) SetC(c bool) {}

// Get Zero flag
func (s *Status) Z() bool {
	return (*s & 0b0000_0010) == 0b0000_0010
}

// Set Zero flag
func (s *Status) SetZ(z bool) {
	if z {
		// Set Zero flag
		*s = *s | 0b0000_0010
	} else {
		// Unset Zero flag
		*s = *s & 0b1111_1101
	}
}

// Get Interrupt Disable flag
func (s *Status) I() bool {
	return (*s & 0b0000_0100) == 0b0000_0100
}

// Set Interrupt Disable flag
func (s *Status) SetI() {}

// Get Decimal flag
func (s *Status) D() bool {
	return (*s & 0b0000_1000) == 0b0000_1000
}

// Set Decimal flag
func (s *Status) SetD() {}

// Get B flag
func (s *Status) B() bool {
	return (*s & 0b0001_0000) == 0b0001_0000
}

// Set B flag
func (s *Status) SetB() {}

// Get Overflow flag
func (s *Status) O() bool {
	return (*s & 0b0100_0000) == 0b0100_0000
}

// Set Overflow flag
func (s *Status) SetO() {}

// Get Negative flag
func (s *Status) N() bool {
	return (*s & 0b1000_0000) == 0b1000_0000
}

// Set Negative flag
func (s *Status) SetN(n bool) {
	if n {
		// Set Negative flag
		*s = *s | 0b1000_0000
	} else {
		// Unset Negative flag
		*s = *s & 0b0111_1111
	}
}
