package main

import "testing"

func Test_strCount(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strCount(tt.args.str); got != tt.want {
				t.Errorf("strCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
