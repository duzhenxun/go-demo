package main

import "testing"

func test1_chanDemo(t *testing.T) {
	type args struct {
		a int
		b int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{"name", args{1, 2}, 3},
		{"name2", args{11, 2}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := chanDemo(tt.args.a, tt.args.b); got != tt.want {
				t.Errorf("chanDemo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Benchmark_chanDemo(b *testing.B) {
	args1 := 1
	args2 := 2
	want := 3
	for i := 0; i < b.N; i++ {
		got := chanDemo(args1, args2)
		if got != want {
			b.Errorf("chanDemo()=%v,want %v", got, want)
		}
	}

}
