package cpu

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tabo-syu/famicom/internal/bus"
	"github.com/tabo-syu/famicom/internal/memory"
	"github.com/tabo-syu/famicom/internal/rom"
)

// newTestCPU は 0x03_00 にプログラムを配置し、リセット済みの CPU を返す。
// プログラムは必ず BRK (0x00) で終端すること。Run() は BRK で停止する。
func newTestCPU(t *testing.T, program []byte) *CPU {
	t.Helper()

	mem := memory.NewMemory()
	r, err := rom.NewROM(validrom)
	assert.NoError(t, err)

	cpu := NewCPU(bus.NewBus(&mem, r))
	cpu.loadForTest(program)
	cpu.Reset(0x00_00)

	return &cpu
}

// ---------------------------------------------------------------------------
// 命令テーブルの整合性
// ---------------------------------------------------------------------------

// Test_Instructions_OfficialOpcodeCount は公式オペコードが 151 個揃っていることを確認する。
func Test_Instructions_OfficialOpcodeCount(t *testing.T) {
	assert.Len(t, NewInstructions(), 151)
}

// Test_Instructions_AllOpcodesAreDispatchable はテーブルに載っている全オペコードが
// instruction.Call の switch に実装されていることを確認する。
// テーブルにニーモニックを追加したのに switch に case を足し忘れると、
// 実行時に "unexpected opcode" で静かに止まってしまうため。
func Test_Instructions_AllOpcodesAreDispatchable(t *testing.T) {
	for code, i := range NewInstructions() {
		t.Run(fmt.Sprintf("%#02x_%s", code, i.opcode), func(t *testing.T) {
			cpu := newTestCPU(t, []byte{0x00})

			err := i.Call(cpu)
			if i.opcode == "BRK" {
				// BRK は Run() を止めるためのエラーを返す。
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
		})
	}
}

// Test_Instructions_BytesMatchAddressingMode はアドレッシングモードごとの
// 命令長が仕様どおりかを確認する。テーブルの打ち間違いを検出する。
func Test_Instructions_BytesMatchAddressingMode(t *testing.T) {
	want := map[addressingMode]uint16{
		ImpliedMode:     1,
		AccumulatorMode: 1,
		ImmediateMode:   2,
		ZeroPageMode:    2,
		ZeroPageXMode:   2,
		ZeroPageYMode:   2,
		RelativeMode:    2,
		IndirectXMode:   2,
		IndirectYMode:   2,
		AbsoluteMode:    3,
		AbsoluteXMode:   3,
		AbsoluteYMode:   3,
		IndirectMode:    3,
	}

	for code, i := range NewInstructions() {
		bytes, ok := want[i.mode]
		assert.True(t, ok, "%#02x %s: 未知のアドレッシングモード", code, i.opcode)
		assert.Equal(t, bytes, i.bytes, "%#02x %s", code, i.opcode)
	}
}

// ---------------------------------------------------------------------------
// BIT: N と V はメモリの値から、Z だけが A との AND 結果から決まる
// ---------------------------------------------------------------------------

func Test_BIT_FlagsComeFromMemoryNotAccumulator(t *testing.T) {
	tests := []struct {
		name                string
		a, value            byte
		wantZ, wantN, wantO bool
	}{
		// A が 0 でも N/V はメモリのビット 7・6 がそのまま出る
		{"A=0x00 M=0xC0", 0x00, 0b1100_0000, true, true, true},
		// A が全ビット 1 でもメモリが 0 なら N/V は下りる
		{"A=0xFF M=0x00", 0xFF, 0b0000_0000, true, false, false},
		// AND 結果には現れないビット 7 だけを見る
		{"A=0x01 M=0x81", 0x01, 0b1000_0001, false, true, false},
		// AND 結果には現れないビット 6 だけを見る
		{"A=0x01 M=0x41", 0x01, 0b0100_0001, false, false, true},
		// N/V 両方立つがビットは重ならない
		{"A=0x0F M=0xC0", 0x0F, 0b1100_0000, true, true, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t, []byte{0x24, 0x40, 0x00})
			cpu.registerA = tt.a
			cpu.Bus.WriteMemory(0x00_40, tt.value)
			cpu.Run()

			assert.Equal(t, tt.wantZ, cpu.status.z(), "Z")
			assert.Equal(t, tt.wantN, cpu.status.n(), "N")
			assert.Equal(t, tt.wantO, cpu.status.o(), "V")
			// BIT はアキュムレータを破壊しない
			assert.Equal(t, tt.a, cpu.registerA)
		})
	}
}

// ---------------------------------------------------------------------------
// シフト・ローテート: メモリ操作時の Z フラグはメモリの結果で決まる
// ---------------------------------------------------------------------------

func Test_ShiftAndRotate_ZeroFlagFromMemoryResult(t *testing.T) {
	tests := []struct {
		name    string
		code    byte
		carryIn bool
		value   byte
		want    byte
		wantZ   bool
		wantN   bool
		wantC   bool
	}{
		// 結果が 0 なので Z が立つ (A は非 0 のまま)
		{"ASL result zero", 0x06, false, 0b1000_0000, 0b0000_0000, true, false, true},
		{"LSR result zero", 0x46, false, 0b0000_0001, 0b0000_0000, true, false, true},
		{"ROL result zero", 0x26, false, 0b1000_0000, 0b0000_0000, true, false, true},
		{"ROR result zero", 0x66, false, 0b0000_0001, 0b0000_0000, true, false, true},
		// 結果が非 0 なので Z は下りる
		{"ASL result nonzero", 0x06, false, 0b0100_0001, 0b1000_0010, false, true, false},
		{"LSR result nonzero", 0x46, false, 0b1000_0010, 0b0100_0001, false, false, false},
		{"ROL carry in", 0x26, true, 0b0000_0000, 0b0000_0001, false, false, false},
		{"ROR carry in", 0x66, true, 0b0000_0000, 0b1000_0000, false, true, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t, []byte{tt.code, 0x40, 0x00})
			// A を非 0 にしておき、フラグが A ではなく結果から決まることを保証する
			cpu.registerA = 0xFF
			cpu.status.setC(tt.carryIn)
			cpu.Bus.WriteMemory(0x00_40, tt.value)
			cpu.Run()

			assert.Equal(t, tt.want, cpu.Bus.ReadMemory(0x00_40))
			assert.Equal(t, tt.wantZ, cpu.status.z(), "Z")
			assert.Equal(t, tt.wantN, cpu.status.n(), "N")
			assert.Equal(t, tt.wantC, cpu.status.c(), "C")
			// メモリ操作なのでアキュムレータは変化しない
			assert.Equal(t, byte(0xFF), cpu.registerA)
		})
	}
}

// ---------------------------------------------------------------------------
// JMP (Indirect) のページ跨ぎバグ
// ---------------------------------------------------------------------------

// Test_JMP_IndirectPageBoundaryBug は下位バイトが 0xFF のとき、
// 上位バイトが次ページではなく同一ページの先頭から読まれることを確認する。
func Test_JMP_IndirectPageBoundaryBug(t *testing.T) {
	cpu := newTestCPU(t, []byte{0x6C, 0xFF, 0x04, 0x00})
	// 実機は 0x04FF と 0x0400 を読む (0x0500 ではない)
	cpu.Bus.WriteMemory(0x04_FF, 0x44)
	cpu.Bus.WriteMemory(0x04_00, 0x05)
	cpu.Bus.WriteMemory(0x05_00, 0x06)
	// 0x0544 に INX を置く
	cpu.Bus.WriteMemory(0x05_44, 0xE8)
	cpu.Bus.WriteMemory(0x05_45, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x05_46), cpu.ProgramCounter)
}

// Test_JMP_IndirectNoPageBoundary はページ跨ぎでない通常ケースが
// バグ再現によって壊れていないことを確認する。
func Test_JMP_IndirectNoPageBoundary(t *testing.T) {
	cpu := newTestCPU(t, []byte{0x6C, 0xFE, 0x04, 0x00})
	cpu.Bus.WriteMemory(0x04_FE, 0x44)
	cpu.Bus.WriteMemory(0x04_FF, 0x05)
	cpu.Bus.WriteMemory(0x05_44, 0xE8)
	cpu.Bus.WriteMemory(0x05_45, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	assert.Equal(t, uint16(0x05_46), cpu.ProgramCounter)
}

// ---------------------------------------------------------------------------
// ADC / SBC のオーバーフロー真理値表
// ---------------------------------------------------------------------------

// Test_ADC_OverflowTruthTable は符号付き加算の 4 象限すべてを検証する。
// V は「同符号どうしを足して符号が変わったとき」だけ立つ。
func Test_ADC_OverflowTruthTable(t *testing.T) {
	tests := []struct {
		a, value byte
		carryIn  bool
		want     byte
		wantC    bool
		wantO    bool
	}{
		{0x50, 0x10, false, 0x60, false, false}, // 正 + 正 = 正
		{0x50, 0x50, false, 0xA0, false, true},  // 正 + 正 = 負 -> V
		{0x50, 0x90, false, 0xE0, false, false}, // 正 + 負 = 負
		{0x50, 0xD0, false, 0x20, true, false},  // 正 + 負 = 正 (桁上がり)
		{0xD0, 0x10, false, 0xE0, false, false}, // 負 + 正 = 負
		{0xD0, 0x50, false, 0x20, true, false},  // 負 + 正 = 正
		{0xD0, 0x90, false, 0x60, true, true},   // 負 + 負 = 正 -> V
		{0xD0, 0xD0, false, 0xA0, true, false},  // 負 + 負 = 負
		{0x50, 0x10, true, 0x61, false, false},  // キャリー入力が加算される
		{0xFF, 0x00, true, 0x00, true, false},   // キャリー入力だけで桁上がり
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#02x+%#02x+%t", tt.a, tt.value, tt.carryIn), func(t *testing.T) {
			cpu := newTestCPU(t, []byte{0x69, tt.value, 0x00})
			cpu.registerA = tt.a
			cpu.status.setC(tt.carryIn)
			cpu.Run()

			assert.Equal(t, tt.want, cpu.registerA)
			assert.Equal(t, tt.wantC, cpu.status.c(), "C")
			assert.Equal(t, tt.wantO, cpu.status.o(), "V")
			assert.Equal(t, tt.want == 0, cpu.status.z(), "Z")
			assert.Equal(t, tt.want&0b1000_0000 != 0, cpu.status.n(), "N")
		})
	}
}

// Test_SBC_OverflowTruthTable は符号付き減算の 4 象限すべてを検証する。
// C は「借りが発生しなかった」ことを表すので、C=true が正常な減算を意味する。
func Test_SBC_OverflowTruthTable(t *testing.T) {
	tests := []struct {
		a, value byte
		carryIn  bool
		want     byte
		wantC    bool
		wantO    bool
	}{
		{0x50, 0xF0, true, 0x60, false, false}, // 正 - 負 = 正 (借りあり)
		{0x50, 0xB0, true, 0xA0, false, true},  // 正 - 負 = 負 -> V
		{0x50, 0x70, true, 0xE0, false, false}, // 正 - 正 = 負 (借りあり)
		{0x50, 0x30, true, 0x20, true, false},  // 正 - 正 = 正
		{0xD0, 0xF0, true, 0xE0, false, false}, // 負 - 負 = 負 (借りあり)
		{0xD0, 0xB0, true, 0x20, true, false},  // 負 - 負 = 正
		{0xD0, 0x70, true, 0x60, true, true},   // 負 - 正 = 正 -> V
		{0xD0, 0x30, true, 0xA0, true, false},  // 負 - 正 = 負
		{0x50, 0x30, false, 0x1F, true, false}, // C=0 のとき追加で 1 引かれる
		{0x00, 0x00, false, 0xFF, false, false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%#02x-%#02x-%t", tt.a, tt.value, tt.carryIn), func(t *testing.T) {
			cpu := newTestCPU(t, []byte{0xE9, tt.value, 0x00})
			cpu.registerA = tt.a
			cpu.status.setC(tt.carryIn)
			cpu.Run()

			assert.Equal(t, tt.want, cpu.registerA)
			assert.Equal(t, tt.wantC, cpu.status.c(), "C")
			assert.Equal(t, tt.wantO, cpu.status.o(), "V")
			assert.Equal(t, tt.want == 0, cpu.status.z(), "Z")
			assert.Equal(t, tt.want&0b1000_0000 != 0, cpu.status.n(), "N")
		})
	}
}

// ---------------------------------------------------------------------------
// 比較命令の N フラグ
// ---------------------------------------------------------------------------

// Test_Compare_AllFlags は CMP/CPX/CPY が「レジスタ - メモリ」の 3 通り
// (未満・等しい・より大きい) すべてで N/Z/C を正しく更新することを確認する。
func Test_Compare_AllFlags(t *testing.T) {
	tests := []struct {
		name  string
		code  byte
		reg   func(*CPU, byte)
		value byte
		set   byte
		wantN bool
		wantZ bool
		wantC bool
	}{
		{"CMP less", 0xC9, func(c *CPU, v byte) { c.registerA = v }, 0x20, 0x10, true, false, false},
		{"CMP equal", 0xC9, func(c *CPU, v byte) { c.registerA = v }, 0x10, 0x10, false, true, true},
		{"CMP greater", 0xC9, func(c *CPU, v byte) { c.registerA = v }, 0x10, 0x20, false, false, true},
		{"CMP negative result", 0xC9, func(c *CPU, v byte) { c.registerA = v }, 0x01, 0x80, false, false, true},
		{"CPX less", 0xE0, func(c *CPU, v byte) { c.registerX = v }, 0x20, 0x10, true, false, false},
		{"CPX equal", 0xE0, func(c *CPU, v byte) { c.registerX = v }, 0x10, 0x10, false, true, true},
		{"CPX greater", 0xE0, func(c *CPU, v byte) { c.registerX = v }, 0x10, 0x20, false, false, true},
		{"CPY less", 0xC0, func(c *CPU, v byte) { c.registerY = v }, 0x20, 0x10, true, false, false},
		{"CPY equal", 0xC0, func(c *CPU, v byte) { c.registerY = v }, 0x10, 0x10, false, true, true},
		{"CPY greater", 0xC0, func(c *CPU, v byte) { c.registerY = v }, 0x10, 0x20, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t, []byte{tt.code, tt.value, 0x00})
			tt.reg(cpu, tt.set)
			cpu.Run()

			assert.Equal(t, tt.wantN, cpu.status.n(), "N")
			assert.Equal(t, tt.wantZ, cpu.status.z(), "Z")
			assert.Equal(t, tt.wantC, cpu.status.c(), "C")
		})
	}
}

// ---------------------------------------------------------------------------
// アドレッシングモードの折り返し・ページ跨ぎ
// ---------------------------------------------------------------------------

// Test_ZeroPageX_WrapsWithinZeroPage は 0xFF + X がゼロページ内で折り返すことを確認する。
func Test_ZeroPageX_WrapsWithinZeroPage(t *testing.T) {
	cpu := newTestCPU(t, []byte{0xB5, 0xFF, 0x00})
	cpu.registerX = 0x02
	// 0x01FF ではなく 0x0001 を読む
	cpu.Bus.WriteMemory(0x00_01, 0x42)
	cpu.Bus.WriteMemory(0x01_FF, 0x99)
	cpu.Run()

	assert.Equal(t, byte(0x42), cpu.registerA)
}

// Test_ZeroPageY_WrapsWithinZeroPage は LDX $nn,Y の折り返しを確認する。
func Test_ZeroPageY_WrapsWithinZeroPage(t *testing.T) {
	cpu := newTestCPU(t, []byte{0xB6, 0xFF, 0x00})
	cpu.registerY = 0x03
	cpu.Bus.WriteMemory(0x00_02, 0x42)
	cpu.Run()

	assert.Equal(t, byte(0x42), cpu.registerX)
}

// Test_IndirectX_WrapsWithinZeroPage はポインタ自体と、
// ポインタ +1 の両方がゼロページ内で折り返すことを確認する。
func Test_IndirectX_WrapsWithinZeroPage(t *testing.T) {
	cpu := newTestCPU(t, []byte{0xA1, 0xFE, 0x00})
	cpu.registerX = 0x01
	// ポインタは 0x00FF。上位バイトは 0x0100 ではなく 0x0000 から読む。
	cpu.Bus.WriteMemory(0x00_FF, 0x30)
	cpu.Bus.WriteMemory(0x00_00, 0x04)
	cpu.Bus.WriteMemory(0x01_00, 0x99)
	cpu.Bus.WriteMemory(0x04_30, 0x42)
	cpu.Run()

	assert.Equal(t, byte(0x42), cpu.registerA)
}

// Test_IndirectY_WrapsWithinZeroPage は基底ポインタの上位バイト読み出しが
// ゼロページ内で折り返すことを確認する。
func Test_IndirectY_WrapsWithinZeroPage(t *testing.T) {
	cpu := newTestCPU(t, []byte{0xB1, 0xFF, 0x00})
	cpu.registerY = 0x02
	cpu.Bus.WriteMemory(0x00_FF, 0x30)
	cpu.Bus.WriteMemory(0x00_00, 0x04)
	cpu.Bus.WriteMemory(0x04_32, 0x42)
	cpu.Run()

	assert.Equal(t, byte(0x42), cpu.registerA)
}

// Test_IndirectY_CrossesPage は Y の加算でページを跨ぐケースを確認する。
func Test_IndirectY_CrossesPage(t *testing.T) {
	cpu := newTestCPU(t, []byte{0xB1, 0x40, 0x00})
	cpu.registerY = 0x02
	cpu.Bus.WriteMemory(0x00_40, 0xFF)
	cpu.Bus.WriteMemory(0x00_41, 0x04)
	// 0x04FF + 0x02 = 0x0501
	cpu.Bus.WriteMemory(0x05_01, 0x42)
	cpu.Run()

	assert.Equal(t, byte(0x42), cpu.registerA)
}

// Test_AbsoluteX_CrossesPage は絶対アドレス + X がページを跨いでも
// 正しいアドレスになることを確認する。
func Test_AbsoluteX_CrossesPage(t *testing.T) {
	cpu := newTestCPU(t, []byte{0xBD, 0xFF, 0x04, 0x00})
	cpu.registerX = 0x02
	cpu.Bus.WriteMemory(0x05_01, 0x42)
	cpu.Run()

	assert.Equal(t, byte(0x42), cpu.registerA)
}

// ---------------------------------------------------------------------------
// スタック
// ---------------------------------------------------------------------------

// Test_Stack_PushWrapsAround はスタックポインタが 0x00 から 0xFF へ
// 折り返すことを確認する。
func Test_Stack_PushWrapsAround(t *testing.T) {
	cpu := newTestCPU(t, []byte{0x48, 0x48, 0x00})
	cpu.stackPointer = stackPointer(0x00)
	cpu.registerA = 0x42
	cpu.Run()

	assert.Equal(t, byte(0x42), cpu.Bus.ReadMemory(0x01_00))
	assert.Equal(t, byte(0x42), cpu.Bus.ReadMemory(0x01_FF))
	assert.Equal(t, byte(0xFE), byte(cpu.stackPointer))
}

// Test_Stack_PopWrapsAround はポップ側の折り返しを確認する。
func Test_Stack_PopWrapsAround(t *testing.T) {
	cpu := newTestCPU(t, []byte{0x68, 0x00})
	cpu.stackPointer = stackPointer(0xFF)
	cpu.Bus.WriteMemory(0x01_00, 0x42)
	cpu.Run()

	assert.Equal(t, byte(0x42), cpu.registerA)
	assert.Equal(t, byte(0x00), byte(cpu.stackPointer))
}

// Test_JSR_RTS_Nested は入れ子のサブルーチン呼び出しで
// 戻りアドレスが正しく積み降ろしされることを確認する。
func Test_JSR_RTS_Nested(t *testing.T) {
	// 0x0300: JSR $0400
	// 0x0303: BRK
	cpu := newTestCPU(t, []byte{0x20, 0x00, 0x04, 0x00})
	// 0x0400: JSR $0500 / RTS
	cpu.Bus.WriteMemory(0x04_00, 0x20)
	cpu.Bus.WriteMemory(0x04_01, 0x00)
	cpu.Bus.WriteMemory(0x04_02, 0x05)
	cpu.Bus.WriteMemory(0x04_03, 0x60)
	// 0x0500: INX / RTS
	cpu.Bus.WriteMemory(0x05_00, 0xE8)
	cpu.Bus.WriteMemory(0x05_01, 0x60)
	cpu.Run()

	assert.Equal(t, byte(0x01), cpu.registerX)
	// 2 回積んで 2 回降ろすのでリセット直後の値に戻る
	assert.Equal(t, byte(0xFF), byte(cpu.stackPointer))
	assert.Equal(t, uint16(0x03_04), cpu.ProgramCounter)
}

// Test_RTI_RestoresStatusAndProgramCounter は RTI が積まれた PC を
// そのまま (+1 せずに) 復元することを確認する。
func Test_RTI_RestoresStatusAndProgramCounter(t *testing.T) {
	cpu := newTestCPU(t, []byte{0x40, 0x00})
	// 積まれている順: PC 上位 -> PC 下位 -> ステータス
	cpu.Bus.WriteMemory(0x01_FF, 0x04)
	cpu.Bus.WriteMemory(0x01_FE, 0x00)
	cpu.Bus.WriteMemory(0x01_FD, 0b0100_0001)
	cpu.stackPointer = stackPointer(0xFC)
	// 復帰先 0x0400 に BRK を置いて即停止させる
	cpu.Bus.WriteMemory(0x04_00, 0x00)
	cpu.Run()

	assert.Equal(t, byte(0b0100_0001), byte(cpu.status))
	assert.Equal(t, uint16(0x04_01), cpu.ProgramCounter)
	assert.Equal(t, byte(0xFF), byte(cpu.stackPointer))
}

// ---------------------------------------------------------------------------
// ストア命令
// ---------------------------------------------------------------------------

// Test_Store_AddressingModes は STA/STX/STY が各アドレッシングモードで
// 正しい場所に書き込むことを確認する。
func Test_Store_AddressingModes(t *testing.T) {
	tests := []struct {
		name    string
		program []byte
		setup   func(*CPU)
		address uint16
		want    byte
	}{
		{"STA ZeroPage", []byte{0x85, 0x40, 0x00}, func(c *CPU) { c.registerA = 0x42 }, 0x00_40, 0x42},
		{"STA ZeroPageX", []byte{0x95, 0x40, 0x00}, func(c *CPU) { c.registerA, c.registerX = 0x42, 0x02 }, 0x00_42, 0x42},
		{"STA Absolute", []byte{0x8D, 0x40, 0x04, 0x00}, func(c *CPU) { c.registerA = 0x42 }, 0x04_40, 0x42},
		{"STA AbsoluteX", []byte{0x9D, 0x40, 0x04, 0x00}, func(c *CPU) { c.registerA, c.registerX = 0x42, 0x02 }, 0x04_42, 0x42},
		{"STA AbsoluteY", []byte{0x99, 0x40, 0x04, 0x00}, func(c *CPU) { c.registerA, c.registerY = 0x42, 0x03 }, 0x04_43, 0x42},
		{"STX ZeroPage", []byte{0x86, 0x40, 0x00}, func(c *CPU) { c.registerX = 0x42 }, 0x00_40, 0x42},
		{"STX ZeroPageY", []byte{0x96, 0x40, 0x00}, func(c *CPU) { c.registerX, c.registerY = 0x42, 0x02 }, 0x00_42, 0x42},
		{"STX Absolute", []byte{0x8E, 0x40, 0x04, 0x00}, func(c *CPU) { c.registerX = 0x42 }, 0x04_40, 0x42},
		{"STY ZeroPage", []byte{0x84, 0x40, 0x00}, func(c *CPU) { c.registerY = 0x42 }, 0x00_40, 0x42},
		{"STY ZeroPageX", []byte{0x94, 0x40, 0x00}, func(c *CPU) { c.registerY, c.registerX = 0x42, 0x02 }, 0x00_42, 0x42},
		{"STY Absolute", []byte{0x8C, 0x40, 0x04, 0x00}, func(c *CPU) { c.registerY = 0x42 }, 0x04_40, 0x42},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t, tt.program)
			tt.setup(cpu)
			cpu.Run()

			assert.Equal(t, tt.want, cpu.Bus.ReadMemory(tt.address))
		})
	}
}

// Test_Store_IndirectModes は STA の間接アドレッシングを確認する。
func Test_Store_IndirectModes(t *testing.T) {
	t.Run("STA IndirectX", func(t *testing.T) {
		cpu := newTestCPU(t, []byte{0x81, 0x40, 0x00})
		cpu.registerA = 0x42
		cpu.registerX = 0x02
		cpu.Bus.WriteMemory(0x00_42, 0x30)
		cpu.Bus.WriteMemory(0x00_43, 0x04)
		cpu.Run()

		assert.Equal(t, byte(0x42), cpu.Bus.ReadMemory(0x04_30))
	})

	t.Run("STA IndirectY", func(t *testing.T) {
		cpu := newTestCPU(t, []byte{0x91, 0x40, 0x00})
		cpu.registerA = 0x42
		cpu.registerY = 0x02
		cpu.Bus.WriteMemory(0x00_40, 0x30)
		cpu.Bus.WriteMemory(0x00_41, 0x04)
		cpu.Run()

		assert.Equal(t, byte(0x42), cpu.Bus.ReadMemory(0x04_32))
	})
}

// ---------------------------------------------------------------------------
// リード・モディファイ・ライト命令のアドレッシングモード
// ---------------------------------------------------------------------------

// Test_ReadModifyWrite_AddressingModes は INC/DEC/ASL/LSR/ROL/ROR が
// 各アドレッシングモードで正しいアドレスを更新することを確認する。
func Test_ReadModifyWrite_AddressingModes(t *testing.T) {
	tests := []struct {
		name    string
		program []byte
		setup   func(*CPU)
		address uint16
		initial byte
		want    byte
	}{
		{"INC ZeroPage", []byte{0xE6, 0x40, 0x00}, nil, 0x00_40, 0x41, 0x42},
		{"INC ZeroPageX", []byte{0xF6, 0x40, 0x00}, func(c *CPU) { c.registerX = 0x02 }, 0x00_42, 0x41, 0x42},
		{"INC Absolute", []byte{0xEE, 0x40, 0x04, 0x00}, nil, 0x04_40, 0x41, 0x42},
		{"INC AbsoluteX", []byte{0xFE, 0x40, 0x04, 0x00}, func(c *CPU) { c.registerX = 0x02 }, 0x04_42, 0x41, 0x42},
		{"DEC ZeroPage", []byte{0xC6, 0x40, 0x00}, nil, 0x00_40, 0x43, 0x42},
		{"DEC ZeroPageX", []byte{0xD6, 0x40, 0x00}, func(c *CPU) { c.registerX = 0x02 }, 0x00_42, 0x43, 0x42},
		{"DEC Absolute", []byte{0xCE, 0x40, 0x04, 0x00}, nil, 0x04_40, 0x43, 0x42},
		{"DEC AbsoluteX", []byte{0xDE, 0x40, 0x04, 0x00}, func(c *CPU) { c.registerX = 0x02 }, 0x04_42, 0x43, 0x42},
		{"ASL Absolute", []byte{0x0E, 0x40, 0x04, 0x00}, nil, 0x04_40, 0b0010_0001, 0b0100_0010},
		{"ASL AbsoluteX", []byte{0x1E, 0x40, 0x04, 0x00}, func(c *CPU) { c.registerX = 0x02 }, 0x04_42, 0b0010_0001, 0b0100_0010},
		{"LSR Absolute", []byte{0x4E, 0x40, 0x04, 0x00}, nil, 0x04_40, 0b0100_0010, 0b0010_0001},
		{"LSR AbsoluteX", []byte{0x5E, 0x40, 0x04, 0x00}, func(c *CPU) { c.registerX = 0x02 }, 0x04_42, 0b0100_0010, 0b0010_0001},
		{"ROL Absolute", []byte{0x2E, 0x40, 0x04, 0x00}, nil, 0x04_40, 0b0010_0001, 0b0100_0010},
		{"ROL AbsoluteX", []byte{0x3E, 0x40, 0x04, 0x00}, func(c *CPU) { c.registerX = 0x02 }, 0x04_42, 0b0010_0001, 0b0100_0010},
		{"ROR Absolute", []byte{0x6E, 0x40, 0x04, 0x00}, nil, 0x04_40, 0b0100_0010, 0b0010_0001},
		{"ROR AbsoluteX", []byte{0x7E, 0x40, 0x04, 0x00}, func(c *CPU) { c.registerX = 0x02 }, 0x04_42, 0b0100_0010, 0b0010_0001},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t, tt.program)
			if tt.setup != nil {
				tt.setup(cpu)
			}
			cpu.Bus.WriteMemory(tt.address, tt.initial)
			cpu.Run()

			assert.Equal(t, tt.want, cpu.Bus.ReadMemory(tt.address))
		})
	}
}

// ---------------------------------------------------------------------------
// 論理演算・算術演算のアドレッシングモード網羅
// ---------------------------------------------------------------------------

// Test_Alu_AddressingModes は ADC/AND/EOR/ORA/SBC/CMP が
// メモリ系アドレッシングモードで同じオペランドを読み出すことを確認する。
func Test_Alu_AddressingModes(t *testing.T) {
	// いずれも 0x0440 (またはゼロページ 0x40) の値 0x0F を読む
	tests := []struct {
		name    string
		program []byte
		setup   func(*CPU)
		address uint16
	}{
		{"ZeroPage", []byte{0x25, 0x40, 0x00}, nil, 0x00_40},
		{"ZeroPageX", []byte{0x35, 0x3E, 0x00}, func(c *CPU) { c.registerX = 0x02 }, 0x00_40},
		{"Absolute", []byte{0x2D, 0x40, 0x04, 0x00}, nil, 0x04_40},
		{"AbsoluteX", []byte{0x3D, 0x3E, 0x04, 0x00}, func(c *CPU) { c.registerX = 0x02 }, 0x04_40},
		{"AbsoluteY", []byte{0x39, 0x3E, 0x04, 0x00}, func(c *CPU) { c.registerY = 0x02 }, 0x04_40},
	}

	for _, tt := range tests {
		t.Run("AND_"+tt.name, func(t *testing.T) {
			cpu := newTestCPU(t, tt.program)
			cpu.registerA = 0b1111_1100
			if tt.setup != nil {
				tt.setup(cpu)
			}
			cpu.Bus.WriteMemory(tt.address, 0b0000_1111)
			cpu.Run()

			assert.Equal(t, byte(0b0000_1100), cpu.registerA)
		})
	}

	t.Run("AND_IndirectX", func(t *testing.T) {
		cpu := newTestCPU(t, []byte{0x21, 0x3E, 0x00})
		cpu.registerA = 0b1111_1100
		cpu.registerX = 0x02
		cpu.Bus.WriteMemory(0x00_40, 0x40)
		cpu.Bus.WriteMemory(0x00_41, 0x04)
		cpu.Bus.WriteMemory(0x04_40, 0b0000_1111)
		cpu.Run()

		assert.Equal(t, byte(0b0000_1100), cpu.registerA)
	})

	t.Run("AND_IndirectY", func(t *testing.T) {
		cpu := newTestCPU(t, []byte{0x31, 0x40, 0x00})
		cpu.registerA = 0b1111_1100
		cpu.registerY = 0x02
		cpu.Bus.WriteMemory(0x00_40, 0x3E)
		cpu.Bus.WriteMemory(0x00_41, 0x04)
		cpu.Bus.WriteMemory(0x04_40, 0b0000_1111)
		cpu.Run()

		assert.Equal(t, byte(0b0000_1100), cpu.registerA)
	})
}

// ---------------------------------------------------------------------------
// 分岐命令
// ---------------------------------------------------------------------------

// Test_Branch_TakenTargetIsRelativeToNextInstruction は分岐先が
// 「分岐命令の次のアドレス + オフセット」であることを全分岐命令で確認する。
func Test_Branch_TakenTargetIsRelativeToNextInstruction(t *testing.T) {
	tests := []struct {
		name  string
		code  byte
		setup func(*CPU)
	}{
		{"BCC", 0x90, func(c *CPU) { c.status.setC(false) }},
		{"BCS", 0xB0, func(c *CPU) { c.status.setC(true) }},
		{"BEQ", 0xF0, func(c *CPU) { c.status.setZ(true) }},
		{"BNE", 0xD0, func(c *CPU) { c.status.setZ(false) }},
		{"BMI", 0x30, func(c *CPU) { c.status.setN(true) }},
		{"BPL", 0x10, func(c *CPU) { c.status.setN(false) }},
		{"BVC", 0x50, func(c *CPU) { c.status.setO(false) }},
		{"BVS", 0x70, func(c *CPU) { c.status.setO(true) }},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_forward", func(t *testing.T) {
			// 0x0300: 分岐命令 / 0x0301: オフセット -> 次命令は 0x0302
			// 0x0302 + 0x04 = 0x0306 に INX を置く
			cpu := newTestCPU(t, []byte{tt.code, 0x04, 0x00, 0x00, 0x00, 0x00, 0xE8, 0x00})
			tt.setup(cpu)
			cpu.Run()

			assert.Equal(t, byte(0x01), cpu.registerX)
			assert.Equal(t, uint16(0x03_08), cpu.ProgramCounter)
		})

		t.Run(tt.name+"_backward", func(t *testing.T) {
			// 0x0300: JMP $0306 / 0x0303: INX / 0x0304: BRK
			// 0x0306: 分岐命令 (-5) -> 次命令 0x0308 + (-5) = 0x0303
			cpu := newTestCPU(t, []byte{
				0x4C, 0x06, 0x03,
				0xE8, 0x00, 0x00,
				tt.code, 0xFB, 0x00,
			})
			tt.setup(cpu)
			cpu.Run()

			assert.Equal(t, byte(0x01), cpu.registerX)
			assert.Equal(t, uint16(0x03_05), cpu.ProgramCounter)
		})
	}
}

// Test_Branch_NotTakenFallsThrough は条件不成立のとき
// 分岐命令の直後に進むことを確認する。
func Test_Branch_NotTakenFallsThrough(t *testing.T) {
	tests := []struct {
		name  string
		code  byte
		setup func(*CPU)
	}{
		{"BCC", 0x90, func(c *CPU) { c.status.setC(true) }},
		{"BCS", 0xB0, func(c *CPU) { c.status.setC(false) }},
		{"BEQ", 0xF0, func(c *CPU) { c.status.setZ(false) }},
		{"BNE", 0xD0, func(c *CPU) { c.status.setZ(true) }},
		{"BMI", 0x30, func(c *CPU) { c.status.setN(false) }},
		{"BPL", 0x10, func(c *CPU) { c.status.setN(true) }},
		{"BVC", 0x50, func(c *CPU) { c.status.setO(true) }},
		{"BVS", 0x70, func(c *CPU) { c.status.setO(false) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 分岐しなければ 0x0302 の INY が実行される
			cpu := newTestCPU(t, []byte{tt.code, 0x04, 0xC8, 0x00, 0x00, 0x00, 0xE8, 0x00})
			tt.setup(cpu)
			cpu.Run()

			assert.Equal(t, byte(0x00), cpu.registerX)
			assert.Equal(t, byte(0x01), cpu.registerY)
			assert.Equal(t, uint16(0x03_04), cpu.ProgramCounter)
		})
	}
}

// ---------------------------------------------------------------------------
// 転送命令とフラグ
// ---------------------------------------------------------------------------

// Test_Transfer_UpdatesFlags は転送命令の Z/N 更新有無を確認する。
// TXS だけはフラグを一切変更しない。
func Test_Transfer_UpdatesFlags(t *testing.T) {
	tests := []struct {
		name    string
		code    byte
		setup   func(*CPU)
		check   func(*testing.T, *CPU)
		wantZ   bool
		wantN   bool
		noFlags bool
	}{
		{"TAX zero", 0xAA, func(c *CPU) { c.registerA = 0x00; c.registerX = 0xFF },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0x00), c.registerX) }, true, false, false},
		{"TAX negative", 0xAA, func(c *CPU) { c.registerA = 0x80 },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0x80), c.registerX) }, false, true, false},
		{"TAY negative", 0xA8, func(c *CPU) { c.registerA = 0x80 },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0x80), c.registerY) }, false, true, false},
		{"TXA zero", 0x8A, func(c *CPU) { c.registerX = 0x00; c.registerA = 0xFF },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0x00), c.registerA) }, true, false, false},
		{"TYA negative", 0x98, func(c *CPU) { c.registerY = 0x80 },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0x80), c.registerA) }, false, true, false},
		{"TSX negative", 0xBA, func(c *CPU) { c.stackPointer = stackPointer(0x80) },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0x80), c.registerX) }, false, true, false},
		// TXS はスタックポインタを書き換えるだけでフラグに影響しない
		{"TXS no flags", 0x9A, func(c *CPU) { c.registerX = 0x00 },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0x00), byte(c.stackPointer)) }, false, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t, []byte{tt.code, 0x00})
			tt.setup(cpu)
			cpu.Run()

			tt.check(t, cpu)
			assert.Equal(t, tt.wantZ, cpu.status.z(), "Z")
			assert.Equal(t, tt.wantN, cpu.status.n(), "N")
		})
	}
}

// ---------------------------------------------------------------------------
// インクリメント / デクリメントの折り返し
// ---------------------------------------------------------------------------

// Test_IncrementDecrement_Wrap は 0xFF -> 0x00 と 0x00 -> 0xFF の
// 折り返し時のフラグを確認する。
func Test_IncrementDecrement_Wrap(t *testing.T) {
	tests := []struct {
		name    string
		program []byte
		setup   func(*CPU)
		check   func(*testing.T, *CPU)
		wantZ   bool
		wantN   bool
	}{
		{"INX 0xFF->0x00", []byte{0xE8, 0x00}, func(c *CPU) { c.registerX = 0xFF },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0x00), c.registerX) }, true, false},
		{"DEX 0x00->0xFF", []byte{0xCA, 0x00}, func(c *CPU) { c.registerX = 0x00 },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0xFF), c.registerX) }, false, true},
		{"INY 0xFF->0x00", []byte{0xC8, 0x00}, func(c *CPU) { c.registerY = 0xFF },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0x00), c.registerY) }, true, false},
		{"DEY 0x00->0xFF", []byte{0x88, 0x00}, func(c *CPU) { c.registerY = 0x00 },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0xFF), c.registerY) }, false, true},
		{"INC 0xFF->0x00", []byte{0xE6, 0x40, 0x00}, func(c *CPU) { c.Bus.WriteMemory(0x00_40, 0xFF) },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0x00), c.Bus.ReadMemory(0x00_40)) }, true, false},
		{"DEC 0x00->0xFF", []byte{0xC6, 0x40, 0x00}, func(c *CPU) { c.Bus.WriteMemory(0x00_40, 0x00) },
			func(t *testing.T, c *CPU) { assert.Equal(t, byte(0xFF), c.Bus.ReadMemory(0x00_40)) }, false, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t, tt.program)
			tt.setup(cpu)
			cpu.Run()

			tt.check(t, cpu)
			assert.Equal(t, tt.wantZ, cpu.status.z(), "Z")
			assert.Equal(t, tt.wantN, cpu.status.n(), "N")
		})
	}
}

// ---------------------------------------------------------------------------
// フラグ操作命令
// ---------------------------------------------------------------------------

// Test_FlagInstructions はセット/クリア命令が対象のフラグだけを
// 操作することを確認する。
func Test_FlagInstructions(t *testing.T) {
	tests := []struct {
		name   string
		code   byte
		before status
		after  status
	}{
		{"CLC", 0x18, 0b1111_1111, 0b1111_1110},
		{"SEC", 0x38, 0b0000_0000, 0b0000_0001},
		{"CLI", 0x58, 0b1111_1111, 0b1111_1011},
		{"SEI", 0x78, 0b0000_0000, 0b0000_0100},
		{"CLD", 0xD8, 0b1111_1111, 0b1111_0111},
		{"SED", 0xF8, 0b0000_0000, 0b0000_1000},
		{"CLV", 0xB8, 0b1111_1111, 0b1011_1111},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cpu := newTestCPU(t, []byte{tt.code, 0x00})
			cpu.status = tt.before
			cpu.Run()

			assert.Equal(t, byte(tt.after), byte(cpu.status))
		})
	}
}
