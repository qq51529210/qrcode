package qrcode

import (
	"github.com/skip2/go-qrcode"
	"testing"
)

var (
	testStr = "你多1231行上sdfsd岛咖东方航空"
)

func BenchmarkImageMy(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Image(testStr, LevelL)
	}
}

func BenchmarkImageSkip2(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		code, _ := qrcode.New(testStr, qrcode.Low)
		code.Image(128)
	}
}
