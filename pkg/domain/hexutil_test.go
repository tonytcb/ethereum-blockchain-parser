package domain

import "testing"

func Test_hexToDecimalString(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "should be equals to 0",
			args:    args{input: "0x0"},
			want:    0,
			wantErr: false,
		},
		{
			name:    "should be equals to 65",
			args:    args{input: "0x41"},
			want:    65,
			wantErr: false,
		},
		{
			name:    "should be equals to 1024",
			args:    args{input: "0x400"},
			want:    1024,
			wantErr: false,
		},

		{
			name:    "should be invalid due to invalid formatting",
			args:    args{input: "0x"},
			want:    0,
			wantErr: true,
		},
		{
			name:    "should be invalid due to invalid formatting",
			args:    args{input: "0x0400"},
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := hexToDecimalString(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("hexToDecimalString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("hexToDecimalString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
