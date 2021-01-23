package qrcode

// 纠错等级
type Level int

func (l Level) String() string {
	switch l {
	case LevelL:
		return "L"
	case LevelM:
		return "M"
	case LevelQ:
		return "Q"
	case LevelH:
		return "H"
	default:
		panic("code bug")
	}
}

const (
	LevelL Level = iota // 7%
	LevelM              // 15%
	LevelQ              // 25％
	LevelH              // 30％
	maxLevel
)
