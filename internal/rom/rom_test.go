package rom

import (
	"reflect"
	"testing"
)

func TestNewRom(t *testing.T) {
	type args struct {
		raw []byte
	}
	tests := []struct {
		name    string
		args    args
		want    *ROM
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				raw: []byte{
					'N', 'E', 'S', 0x1A,
					0x00, 0x00,
					0b0110_0001, 0b1001_0000,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				},
			},
			want: &ROM{
				Prg:             []byte{},
				Chr:             []byte{},
				Mapper:          0b1001_0110,
				ScreenMirroring: Vertical,
			},
			wantErr: false,
		},
		{
			name: "Failure/Validate 'NES^Z'",
			args: args{
				raw: []byte{
					'N', 'E', 'S', 'A',
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Failure/Validate iNES Version",
			args: args{
				raw: []byte{
					'N', 'E', 'S', 0x1A,
					0x00, 0x00,
					0b0000_0000, 0b0000_1000,
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRom(tt.args.raw)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRom() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRom() = %v, want %v", got, tt.want)
			}
		})
	}
}
