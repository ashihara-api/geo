package usecase

import "testing"

func Test_generateCheckDigit(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput int
		wantErr    bool
	}{
		{
			name: "include_zero",
			args: args{
				s: "05201",
			},
			wantOutput: 9,
			wantErr:    false,
		},
		{
			name: "not_include_zero",
			args: args{
				s: "43511",
			},
			wantOutput: 2,
			wantErr:    false,
		},
		{
			name: "4_length",
			args: args{
				s: "1234",
			},
			wantOutput: 0,
			wantErr:    true,
		},
		{
			name: "4_length",
			args: args{
				s: "123456",
			},
			wantOutput: 0,
			wantErr:    true,
		},
		{
			name: "non_numeric",
			args: args{
				s: "abcde",
			},
			wantOutput: 0,
			wantErr:    true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotOutput, err := generateCheckDigit(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateCheckDigit() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput != tt.wantOutput {
				t.Errorf("generateCheckDigit() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}
