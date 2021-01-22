package qrcode

// 纠错等级
type level int

func (l level) String() string {
	switch l {
	case levelL:
		return "L"
	case levelM:
		return "M"
	case levelQ:
		return "Q"
	case levelH:
		return "H"
	default:
		panic("code bug")
	}
}

const (
	levelL level = iota // 7%
	levelM              // 15%
	levelQ              // 25％
	levelH              // 30％
	maxLevel
)
