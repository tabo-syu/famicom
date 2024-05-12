package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_scrape(t *testing.T) {
	got, err := scrape()
	if err != nil {
		t.Errorf("scrape() error = %v", err)
		return
	}

	assert.Equal(t, 56, len(got))

	for _, op := range got {
		if op.Name == "ADC" {
			assert.Equal(t, "Set if overflow in bit 7", op.Status.C)
			assert.Equal(t, "Set if A = 0", op.Status.Z)
			assert.Equal(t, "Not affected", op.Status.I)
			assert.Equal(t, "Not affected", op.Status.D)
			assert.Equal(t, "Not affected", op.Status.B)
			assert.Equal(t, "Set if sign bit is incorrect", op.Status.V)
			assert.Equal(t, "Set if bit 7 set", op.Status.N)

			for _, mode := range op.Modes {
				switch mode.Name {
				case "Immediate":
					assert.Equal(t, "69", mode.Code)
					assert.Equal(t, "2", mode.Bytes)
					assert.Equal(t, "2", mode.Cycles)
				case "ZeroPage":
					assert.Equal(t, "65", mode.Code)
					assert.Equal(t, "2", mode.Bytes)
					assert.Equal(t, "3", mode.Cycles)
				case "ZeroPageX":
					assert.Equal(t, "75", mode.Code)
					assert.Equal(t, "2", mode.Bytes)
					assert.Equal(t, "4", mode.Cycles)
				case "Absolute":
					assert.Equal(t, "6D", mode.Code)
					assert.Equal(t, "3", mode.Bytes)
					assert.Equal(t, "4", mode.Cycles)
				case "AbsoluteX":
					assert.Equal(t, "7D", mode.Code)
					assert.Equal(t, "3", mode.Bytes)
					assert.Equal(t, "4 /*(+1 if page crossed)*/", mode.Cycles)
				case "AbsoluteY":
					assert.Equal(t, "79", mode.Code)
					assert.Equal(t, "3", mode.Bytes)
					assert.Equal(t, "4 /*(+1 if page crossed)*/", mode.Cycles)
				case "IndirectX":
					assert.Equal(t, "61", mode.Code)
					assert.Equal(t, "2", mode.Bytes)
					assert.Equal(t, "6", mode.Cycles)
				case "IndirectY":
					assert.Equal(t, "71", mode.Code)
					assert.Equal(t, "2", mode.Bytes)
					assert.Equal(t, "5 /*(+1 if page crossed)*/", mode.Cycles)
				default:
					t.Error("Error: Unexpected AddressingMode:", mode.Name)
				}
			}
		}

		if op.Name == "CLD" {
			assert.Equal(t, "Not affected", op.Status.C)
			assert.Equal(t, "Not affected", op.Status.Z)
			assert.Equal(t, "Not affected", op.Status.I)
			assert.Equal(t, "Set to 0", op.Status.D)
			assert.Equal(t, "Not affected", op.Status.B)
			assert.Equal(t, "Not affected", op.Status.V)
			assert.Equal(t, "Not affected", op.Status.N)

			for _, mode := range op.Modes {
				switch mode.Name {
				case "Implied":
					assert.Equal(t, "D8", mode.Code)
					assert.Equal(t, "1", mode.Bytes)
					assert.Equal(t, "2", mode.Cycles)
				default:
					t.Error("Error: Unexpected AddressingMode:", mode.Name)
				}
			}
		}

		if op.Name == "JMP" {
			assert.Equal(t, "Not affected", op.Status.C)
			assert.Equal(t, "Not affected", op.Status.Z)
			assert.Equal(t, "Not affected", op.Status.I)
			assert.Equal(t, "Not affected", op.Status.D)
			assert.Equal(t, "Not affected", op.Status.B)
			assert.Equal(t, "Not affected", op.Status.V)
			assert.Equal(t, "Not affected", op.Status.N)

			for _, mode := range op.Modes {
				switch mode.Name {
				case "Absolute":
					assert.Equal(t, "4C", mode.Code)
					assert.Equal(t, "3", mode.Bytes)
					assert.Equal(t, "3", mode.Cycles)
				case "Indirect":
					assert.Equal(t, "6C", mode.Code)
					assert.Equal(t, "3", mode.Bytes)
					assert.Equal(t, "5", mode.Cycles)
				default:
					t.Error("Error: Unexpected AddressingMode:", mode.Name)
				}
			}
		}

		if op.Name == "SBC" {
			assert.Equal(t, "Clear if overflow in bit 7", op.Status.C)
			assert.Equal(t, "Set if A = 0", op.Status.Z)
			assert.Equal(t, "Not affected", op.Status.I)
			assert.Equal(t, "Not affected", op.Status.D)
			assert.Equal(t, "Not affected", op.Status.B)
			assert.Equal(t, "Set if sign bit is incorrect", op.Status.V)
			assert.Equal(t, "Set if bit 7 set", op.Status.N)

			for _, mode := range op.Modes {
				switch mode.Name {
				case "Immediate":
					assert.Equal(t, "E9", mode.Code)
					assert.Equal(t, "2", mode.Bytes)
					assert.Equal(t, "2", mode.Cycles)
				case "ZeroPage":
					assert.Equal(t, "E5", mode.Code)
					assert.Equal(t, "2", mode.Bytes)
					assert.Equal(t, "3", mode.Cycles)
				case "ZeroPageX":
					assert.Equal(t, "F5", mode.Code)
					assert.Equal(t, "2", mode.Bytes)
					assert.Equal(t, "4", mode.Cycles)
				case "Absolute":
					assert.Equal(t, "ED", mode.Code)
					assert.Equal(t, "3", mode.Bytes)
					assert.Equal(t, "4", mode.Cycles)
				case "AbsoluteX":
					assert.Equal(t, "FD", mode.Code)
					assert.Equal(t, "3", mode.Bytes)
					assert.Equal(t, "4 /*(+1 if page crossed)*/", mode.Cycles)
				case "AbsoluteY":
					assert.Equal(t, "F9", mode.Code)
					assert.Equal(t, "3", mode.Bytes)
					assert.Equal(t, "4 /*(+1 if page crossed)*/", mode.Cycles)
				case "IndirectX":
					assert.Equal(t, "E1", mode.Code)
					assert.Equal(t, "2", mode.Bytes)
					assert.Equal(t, "6", mode.Cycles)
				case "IndirectY":
					assert.Equal(t, "F1", mode.Code)
					assert.Equal(t, "2", mode.Bytes)
					assert.Equal(t, "5 /*(+1 if page crossed)*/", mode.Cycles)
				default:
					t.Error("Error: Unexpected AddressingMode:", mode.Name)
				}
			}
		}

		if op.Name == "TYA" {
			assert.Equal(t, "Not affected", op.Status.C)
			assert.Equal(t, "Set if A = 0", op.Status.Z)
			assert.Equal(t, "Not affected", op.Status.I)
			assert.Equal(t, "Not affected", op.Status.D)
			assert.Equal(t, "Not affected", op.Status.B)
			assert.Equal(t, "Not affected", op.Status.V)
			assert.Equal(t, "Set if bit 7 of A is set", op.Status.N)

			for _, mode := range op.Modes {
				switch mode.Name {
				case "Implied":
					assert.Equal(t, "98", mode.Code)
					assert.Equal(t, "1", mode.Bytes)
					assert.Equal(t, "2", mode.Cycles)
				default:
					t.Error("Error: Unexpected AddressingMode:", mode.Name)
				}
			}
		}
	}
}
