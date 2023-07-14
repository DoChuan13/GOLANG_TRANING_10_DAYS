package demo

import (
	"math"
	"math/rand"
	"testing"
)

func TestAbs(t *testing.T) {
	got := math.Abs(-1)
	if got != 1 {
		t.Errorf("Abs(-1) = %f; want 1", got)
	}
}
func TestSum(t *testing.T) {
	got := 2 + 3
	if got != 5 {
		t.Errorf("Sum = %d; want 5", got)
	}
}

func BenchmarkBigLen(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//fmt.Println(i)
	}
}
func BenchmarkRandInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rand.Int()
	}
}
