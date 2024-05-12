package memory

type Memory [0x1_00_00]uint8

func NewMemory() Memory {
	return [0x1_00_00]uint8{}
}

func (m *Memory) Read(address uint16) uint8 {
	return m[address]
}

func (m *Memory) ReadUint16(position uint16) uint16 {
	low := uint16(m.Read(position))
	high := uint16(m.Read(position + 1))

	return high<<8 | low
}

func (m *Memory) Write(address uint16, data uint8) {
	m[address] = data
}

func (m *Memory) WriteUint16(position uint16, data uint16) {
	high := uint8(data >> 8)
	low := uint8(data & 0x00_FF)

	m.Write(position, low)
	m.Write(position+1, high)
}
