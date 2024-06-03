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
type status byte

func newStatus() status {
	return 0b0000_0000
}

// Get Carry flag
func (s *status) c() bool {
	return (*s & 0b0000_0001) == 0b0000_0001
}

// Set Carry flag
func (s *status) setC(c bool) {
	if c {
		// Set Carry flag
		*s = *s | 0b0000_0001
	} else {
		// Unset Carry flag
		*s = *s & 0b1111_1110
	}
}

// Get Zero flag
func (s *status) z() bool {
	return (*s & 0b0000_0010) == 0b0000_0010
}

// Set Zero flag
func (s *status) setZ(z bool) {
	if z {
		// Set Zero flag
		*s = *s | 0b0000_0010
	} else {
		// Unset Zero flag
		*s = *s & 0b1111_1101
	}
}

// Get Interrupt Disable flag
func (s *status) i() bool {
	return (*s & 0b0000_0100) == 0b0000_0100
}

// Set Interrupt Disable flag
func (s *status) setI(i bool) {
	if i {
		// Set Interrupt Disable flag
		*s = *s | 0b0000_0100
	} else {
		// Unset Interrupt Disable flag
		*s = *s & 0b1111_1011
	}
}

// Get Decimal flag
func (s *status) d() bool {
	return (*s & 0b0000_1000) == 0b0000_1000
}

// Set Decimal flag
func (s *status) setD(d bool) {
	if d {
		// Set Decimal flag
		*s = *s | 0b0000_1000
	} else {
		// Unset Decimal flag
		*s = *s & 0b1111_0111
	}
}

// Get b flag
func (s *status) b() bool {
	return (*s & 0b0001_0000) == 0b0001_0000
}

// Set B flag
func (s *status) setB() {}

// Get Overflow flag
func (s *status) o() bool {
	return (*s & 0b0100_0000) == 0b0100_0000
}

// Set Overflow flag
func (s *status) setO(o bool) {
	if o {
		// Set Decimal flag
		*s = *s | 0b0100_0000
	} else {
		// Unset Decimal flag
		*s = *s & 0b1011_1111
	}
}

// Get Negative flag
func (s *status) n() bool {
	return (*s & 0b1000_0000) == 0b1000_0000
}

// Set Negative flag
func (s *status) setN(n bool) {
	if n {
		// Set Negative flag
		*s = *s | 0b1000_0000
	} else {
		// Unset Negative flag
		*s = *s & 0b0111_1111
	}
}
