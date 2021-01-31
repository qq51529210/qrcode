package qrcode

import "unicode"

type mode int

const (
	numericMode mode = iota
	alphanumericMode
	byteMode
	kanJiMode
	maxMode
)

var (
	// 用于快速判断模式
	alphanumericTable = [256]byte{}
	// 指示器表
	indicatorTable = [maxMode]byte{
		0b00010000,
		0b00100000,
		0b01000000,
		0b10000000,
	}
)

func init() {
	initAlphanumericTable()
}

// 初始化字母表
func initAlphanumericTable() {
	var i byte
	for c := '0'; c <= '9'; c++ {
		alphanumericTable[c] = i
		i++
	}
	for c := 'A'; c <= 'Z'; c++ {
		alphanumericTable[c] = i
		i++
	}
	for _, c := range []byte{' ', '$', '%', '*', '+', '-', '.', '/', ':'} {
		alphanumericTable[c] = i
		i++
	}
}

// 判断编码模式
func analysisMode(str string) mode {
	mode := numericMode
	for _, c := range str {
		if unicode.MaxLatin1 < c {
			if (c >= 0x8140 && c <= 0x9FFC) || (c >= 0xE040 && c <= 0xEBBF) {
				mode = kanJiMode
			} else {
				return byteMode
			}
		} else {
			if c >= '0' && c <= '9' {
				continue
			}
			if alphanumericTable[c] != 0 {
				if mode < alphanumericMode {
					mode = alphanumericMode
				}
			} else {
				return byteMode
			}
		}
	}
	return mode
}
