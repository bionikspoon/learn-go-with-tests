package iteration

import "testing"

func TestRepeat(t *testing.T) {
	type args struct {
		character string
		times     int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"it repeats the lettter", args{"a", 5}, "aaaaa"},
		{"it repeats the lettter", args{"a", 6}, "aaaaaa"},
		{"it repeats the lettter", args{"a", 0}, ""},
		{"it repeats the lettter", args{"b", 1}, "b"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Repeat(tt.args.character, tt.args.times); got != tt.want {
				t.Errorf("Repeat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkRepeat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Repeat("a", 5)
	}
}
