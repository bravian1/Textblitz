package internals

import "testing"

func Test_hammingDistance(t *testing.T) {
	type args struct {
		a uint64
		b uint64
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "identical hashes",
		args: args{a:0b101010, b: 0b101010},
		want:0,
		},
		{
			name: "different hashes",
            args: args{a: 0b101010, b: 0b010101},
            want: 6,
		},
		{
			name: "alternating bits",
			args: args{a: 0b1010101010101010, b: 0b0101010101010101},
			want: 16,
		},
	}


	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := hammingDistance(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("hammingDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}
