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
		{"name1", args{"杜zhenxun"}, 8},
		{"name2", args{"ad1112341234123412341243o"}, 13131},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strCount(tt.args.str); got != tt.want {
				t.Errorf("strCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

//基准测试
//go test -bench=. -benchtime=3s -run=none
//go test -bench=. -benchmem -memprofile memprofile.out -cpuprofile profile.out
//go tool pprof profile.out
func BenchmarkStrCount(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		strCount("杜振训杜振训杜振训")
	}
}
