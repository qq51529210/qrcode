package qrcode

import "unicode"

type mode int

func (m mode) String() string {
	switch m {
	case numericMode:
		return "numeric"
	case alphanumericMode:
		return "alphanumer"
	case byteMode:
		return "byte"
	default:
		panic("code bug")
	}
}

func (m mode) Indicator() byte {
	switch m {
	case numericMode:
		return 0b00010000
	case alphanumericMode:
		return 0b00100000
	case byteMode:
		return 0b01000000
	case kanJiMode:
		return 0b10000000
	default:
		panic("code bug")
	}
}

const (
	numericMode mode = iota
	alphanumericMode
	byteMode
	kanJiMode // 这个模式在没有太大的意义，go的string存放utf-8，可以直接用byteMode
	maxMode
)

var (
	numericModeTable      = [256]byte{} // 用于快速判断模式
	alphanumericModeTable = [256]byte{} // 用于快速判断模式
)

func init() {
	var i byte
	for c := '0'; c <= '9'; c++ {
		numericModeTable[c] = i
		alphanumericModeTable[c] = i
		i++
	}
	for c := 'A'; c <= 'Z'; c++ {
		alphanumericModeTable[c] = i
		i++
	}
	for _, c := range []byte{' ', '$', '%', '*', '+', '-', '.', '/', ':'} {
		alphanumericModeTable[c] = i
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
			if numericModeTable[c] != 0 {
				continue
			} else {
				if alphanumericModeTable[c] != 0 {
					mode = alphanumericMode
				} else {
					return byteMode
				}
			}
		}
	}
	return mode
}
