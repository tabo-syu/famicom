package memory

type Memory interface {
	Read(address uint16) byte
	ReadUint16(position uint16) uint16
	Write(address uint16, data byte)
	WriteUint16(position uint16, data uint16)
	Copy(start int, program []byte)
}

type memory [0x1_00_00]byte

func NewMemory() memory {
	return [0x1_00_00]byte{}
}

func (m *memory) Read(address uint16) byte {
	return m[address]
}

func (m *memory) ReadUint16(position uint16) uint16 {
	low := uint16(m.Read(position))
	high := uint16(m.Read(position + 1))

	return high<<8 | low
}

func (m *memory) Write(address uint16, data byte) {
	m[address] = data
}

func (m *memory) WriteUint16(position uint16, data uint16) {
	high := byte(data >> 8)
	low := byte(data & 0x00_FF)

	m.Write(position, low)
	m.Write(position+1, high)
}

func (m *memory) Copy(start int, program []byte) {
	end := start + len(program)

	copy(m[start:end], program)
}
