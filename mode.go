package qrcode

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

const (
	numericMode mode = iota
	alphanumericMode
	byteMode
	kanJiMode // 这个模式在没有太大的意义，go的string存放utf-8，可以直接用byteMode
	maxMode
)

var (
	alphanumericTable = [256]byte{} // 用于快速判断模式
	modeIndicator     = [maxMode]byte{
		0b00010000,
		0b00100000,
		0b01000000,
		0b10000000,
	}
)

func init() {
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
