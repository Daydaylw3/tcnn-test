package biz

import "testing"

// 1292957 ns/op
func BenchmarkPreallocate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var data [][]byte
		data = make([][]byte, 0, 10000) // 预先分配容量

		for j := 0; j < 10000; j++ {
			newBytes := make([]byte, 1024)
			data = append(data, newBytes)
		}
	}
}

// 1318550 ns/op
func BenchmarkWithoutPreallocate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var data [][]byte
		data = make([][]byte, 0) // 不预先分配

		for j := 0; j < 10000; j++ {
			newBytes := make([]byte, 1024)
			data = append(data, newBytes)
		}
	}
}
