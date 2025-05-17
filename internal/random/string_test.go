package random

import "testing"

func TestString(t *testing.T) {
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name    string
		args    args
		wantMin int
		wantMax int
	}{
		{
			name: "Positive_test_one_letter",
			args: args{
				min: 1,
				max: 1,
			},
			wantMin: 1,
			wantMax: 1,
		},
		{
			name: "Positive_test_several_letters",
			args: args{
				min: 2,
				max: 10,
			},
			wantMin: 2,
			wantMax: 10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := len(RandString(tt.args.min, tt.args.max))
			if got > tt.wantMax || got < tt.wantMin {
				t.Errorf("RandString() got len = %v, want len %v", got, tt.wantMin)
			}
		})
	}
}
