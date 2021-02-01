package qrcode

// 纠错等级
type Level int

const (
	LevelL Level = iota // 7%
	LevelM              // 15%
	LevelQ              // 25％
	LevelH              // 30％
	maxLevel
)

var (
	levelString = [maxMode]string{
		"L", "M", "Q", "H",
	}
)
