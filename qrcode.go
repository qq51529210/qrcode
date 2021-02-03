package qrcode

/*
参考文档：https://www.thonky.com/qr-code-tutorial/
*/

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"sync"
)

const (
	maxMark       = 8
	timingPattern = 6
)

var (
	_pool    sync.Pool
	_palette = color.Palette{
		color.White,
		color.Black,
	}
	_paletteBlack  = uint8(_palette.Index(color.Black))
	_paletteWhite  = uint8(_palette.Index(color.White))
	formatBitTable = [maxLevel][maxMark][]byte{
		{
			{1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 1, 0, 0},
			{1, 1, 1, 0, 0, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1},
			{1, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 0, 1, 0},
			{1, 1, 1, 1, 0, 0, 0, 1, 0, 0, 1, 1, 1, 0, 1},
			{1, 1, 0, 0, 1, 1, 0, 0, 0, 1, 0, 1, 1, 1, 1},
			{1, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 0, 0},
			{1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0, 0, 0, 0, 1},
			{1, 1, 0, 1, 0, 0, 1, 0, 1, 1, 1, 0, 1, 1, 0},
		},
		{
			{1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0},
			{1, 0, 1, 0, 0, 0, 1, 0, 0, 1, 0, 0, 1, 0, 1},
			{1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0},
			{1, 0, 1, 1, 0, 1, 1, 0, 1, 0, 0, 1, 0, 1, 1},
			{1, 0, 0, 0, 1, 0, 1, 1, 1, 1, 1, 1, 0, 0, 1},
			{1, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0},
			{1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 1, 1, 1},
			{1, 0, 0, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0},
		},
		{
			{0, 1, 1, 0, 1, 0, 1, 0, 1, 0, 1, 1, 1, 1, 1},
			{0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0, 0, 0},
			{0, 1, 1, 1, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 1},
			{0, 1, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 1, 1, 0},
			{0, 1, 0, 0, 1, 0, 0, 1, 0, 1, 1, 0, 1, 0, 0},
			{0, 1, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1},
			{0, 1, 0, 1, 1, 1, 0, 1, 1, 0, 1, 1, 0, 1, 0},
			{0, 1, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 0, 1},
		},
		{
			{0, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 0, 0, 1},
			{0, 0, 1, 0, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0},
			{0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1},
			{0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0},
			{0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0},
			{0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 1},
			{0, 0, 0, 1, 1, 0, 1, 0, 0, 0, 0, 1, 1, 0, 0},
			{0, 0, 0, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1},
		},
	}
	versionBitTable = [maxVersion][]byte{
		{},
		{},
		{},
		{},
		{},
		{},
		{0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0},
		{0, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0},
		{0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1},
		{0, 0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1},
		{0, 0, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0},
		{0, 0, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0},
		{0, 0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1},
		{0, 0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 1},
		{0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0, 0},
		{0, 1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1},
		{0, 1, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 1},
		{0, 1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0},
		{0, 1, 0, 1, 0, 0, 1, 0, 0, 1, 1, 0, 1, 0, 0, 1, 1, 0},
		{0, 1, 0, 1, 0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 1, 1},
		{0, 1, 0, 1, 1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 0, 0, 1},
		{0, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 0, 1, 1, 0, 0},
		{0, 1, 1, 0, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 0, 1, 0, 0},
		{0, 1, 1, 0, 0, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 1},
		{0, 1, 1, 0, 1, 0, 1, 1, 1, 1, 1, 0, 1, 0, 1, 0, 1, 1},
		{0, 1, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 0, 0, 1, 1, 1, 0},
		{0, 1, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 1, 1, 0, 1, 0},
		{0, 1, 1, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 1, 1, 1},
		{0, 1, 1, 1, 1, 0, 1, 1, 0, 1, 0, 1, 1, 1, 0, 1, 0, 1},
		{0, 1, 1, 1, 1, 1, 0, 0, 1, 0, 0, 1, 0, 1, 0, 0, 0, 0},
		{1, 0, 0, 0, 0, 0, 1, 0, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1},
		{1, 0, 0, 0, 0, 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0, 0, 0},
		{1, 0, 0, 0, 1, 0, 1, 0, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0},
		{1, 0, 0, 0, 1, 1, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1},
		{1, 0, 0, 1, 0, 0, 1, 0, 1, 1, 0, 0, 0, 0, 1, 0, 1, 1},
		{1, 0, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 1, 0, 1, 1, 1, 0},
		{1, 0, 0, 1, 1, 0, 1, 0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 0},
		{1, 0, 0, 1, 1, 1, 0, 1, 0, 1, 0, 1, 0, 0, 0, 0, 0, 1},
		{1, 0, 1, 0, 0, 0, 1, 1, 0, 0, 0, 1, 1, 0, 1, 0, 0, 1},
	}
	alignmentPatternTable = [maxVersion][]*image.Rectangle{
		{},
		{
			{Min: image.Point{X: 16, Y: 16}, Max: image.Point{X: 20, Y: 20}},
		},
		{
			{Min: image.Point{X: 20, Y: 20}, Max: image.Point{X: 24, Y: 24}},
		},
		{
			{Min: image.Point{X: 24, Y: 24}, Max: image.Point{X: 28, Y: 28}},
		},
		{
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
		},
		{
			{Min: image.Point{X: 32, Y: 32}, Max: image.Point{X: 36, Y: 36}},
		},
		{
			{Min: image.Point{X: 4, Y: 20}, Max: image.Point{X: 8, Y: 24}},
			{Min: image.Point{X: 20, Y: 4}, Max: image.Point{X: 24, Y: 8}},
			{Min: image.Point{X: 20, Y: 20}, Max: image.Point{X: 24, Y: 24}},
			{Min: image.Point{X: 20, Y: 36}, Max: image.Point{X: 24, Y: 40}},
			{Min: image.Point{X: 36, Y: 20}, Max: image.Point{X: 40, Y: 24}},
			{Min: image.Point{X: 36, Y: 36}, Max: image.Point{X: 40, Y: 40}},
		},
		{
			{Min: image.Point{X: 4, Y: 22}, Max: image.Point{X: 8, Y: 26}},
			{Min: image.Point{X: 22, Y: 4}, Max: image.Point{X: 26, Y: 8}},
			{Min: image.Point{X: 22, Y: 22}, Max: image.Point{X: 26, Y: 26}},
			{Min: image.Point{X: 22, Y: 40}, Max: image.Point{X: 26, Y: 44}},
			{Min: image.Point{X: 40, Y: 22}, Max: image.Point{X: 44, Y: 26}},
			{Min: image.Point{X: 40, Y: 40}, Max: image.Point{X: 44, Y: 44}},
		},
		{
			{Min: image.Point{X: 4, Y: 24}, Max: image.Point{X: 8, Y: 28}},
			{Min: image.Point{X: 24, Y: 4}, Max: image.Point{X: 28, Y: 8}},
			{Min: image.Point{X: 24, Y: 24}, Max: image.Point{X: 28, Y: 28}},
			{Min: image.Point{X: 24, Y: 44}, Max: image.Point{X: 28, Y: 48}},
			{Min: image.Point{X: 44, Y: 24}, Max: image.Point{X: 48, Y: 28}},
			{Min: image.Point{X: 44, Y: 44}, Max: image.Point{X: 48, Y: 48}},
		},
		{
			{Min: image.Point{X: 4, Y: 26}, Max: image.Point{X: 8, Y: 30}},
			{Min: image.Point{X: 26, Y: 4}, Max: image.Point{X: 30, Y: 8}},
			{Min: image.Point{X: 26, Y: 26}, Max: image.Point{X: 30, Y: 30}},
			{Min: image.Point{X: 26, Y: 48}, Max: image.Point{X: 30, Y: 52}},
			{Min: image.Point{X: 48, Y: 26}, Max: image.Point{X: 52, Y: 30}},
			{Min: image.Point{X: 48, Y: 48}, Max: image.Point{X: 52, Y: 52}},
		},
		{
			{Min: image.Point{X: 4, Y: 28}, Max: image.Point{X: 8, Y: 32}},
			{Min: image.Point{X: 28, Y: 4}, Max: image.Point{X: 32, Y: 8}},
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
			{Min: image.Point{X: 28, Y: 52}, Max: image.Point{X: 32, Y: 56}},
			{Min: image.Point{X: 52, Y: 28}, Max: image.Point{X: 56, Y: 32}},
			{Min: image.Point{X: 52, Y: 52}, Max: image.Point{X: 56, Y: 56}},
		},
		{
			{Min: image.Point{X: 4, Y: 30}, Max: image.Point{X: 8, Y: 34}},
			{Min: image.Point{X: 30, Y: 4}, Max: image.Point{X: 34, Y: 8}},
			{Min: image.Point{X: 30, Y: 30}, Max: image.Point{X: 34, Y: 34}},
			{Min: image.Point{X: 30, Y: 56}, Max: image.Point{X: 34, Y: 60}},
			{Min: image.Point{X: 56, Y: 30}, Max: image.Point{X: 60, Y: 34}},
			{Min: image.Point{X: 56, Y: 56}, Max: image.Point{X: 60, Y: 60}},
		},
		{
			{Min: image.Point{X: 4, Y: 32}, Max: image.Point{X: 8, Y: 36}},
			{Min: image.Point{X: 32, Y: 4}, Max: image.Point{X: 36, Y: 8}},
			{Min: image.Point{X: 32, Y: 32}, Max: image.Point{X: 36, Y: 36}},
			{Min: image.Point{X: 32, Y: 60}, Max: image.Point{X: 36, Y: 64}},
			{Min: image.Point{X: 60, Y: 32}, Max: image.Point{X: 64, Y: 36}},
			{Min: image.Point{X: 60, Y: 60}, Max: image.Point{X: 64, Y: 64}},
		},
		{
			{Min: image.Point{X: 4, Y: 24}, Max: image.Point{X: 8, Y: 28}},
			{Min: image.Point{X: 4, Y: 44}, Max: image.Point{X: 8, Y: 48}},
			{Min: image.Point{X: 24, Y: 4}, Max: image.Point{X: 28, Y: 8}},
			{Min: image.Point{X: 24, Y: 24}, Max: image.Point{X: 28, Y: 28}},
			{Min: image.Point{X: 24, Y: 44}, Max: image.Point{X: 28, Y: 48}},
			{Min: image.Point{X: 24, Y: 64}, Max: image.Point{X: 28, Y: 68}},
			{Min: image.Point{X: 44, Y: 4}, Max: image.Point{X: 48, Y: 8}},
			{Min: image.Point{X: 44, Y: 24}, Max: image.Point{X: 48, Y: 28}},
			{Min: image.Point{X: 44, Y: 44}, Max: image.Point{X: 48, Y: 48}},
			{Min: image.Point{X: 44, Y: 64}, Max: image.Point{X: 48, Y: 68}},
			{Min: image.Point{X: 64, Y: 24}, Max: image.Point{X: 68, Y: 28}},
			{Min: image.Point{X: 64, Y: 44}, Max: image.Point{X: 68, Y: 48}},
			{Min: image.Point{X: 64, Y: 64}, Max: image.Point{X: 68, Y: 68}},
		},
		{
			{Min: image.Point{X: 4, Y: 24}, Max: image.Point{X: 8, Y: 28}},
			{Min: image.Point{X: 4, Y: 46}, Max: image.Point{X: 8, Y: 50}},
			{Min: image.Point{X: 24, Y: 4}, Max: image.Point{X: 28, Y: 8}},
			{Min: image.Point{X: 24, Y: 24}, Max: image.Point{X: 28, Y: 28}},
			{Min: image.Point{X: 24, Y: 46}, Max: image.Point{X: 28, Y: 50}},
			{Min: image.Point{X: 24, Y: 68}, Max: image.Point{X: 28, Y: 72}},
			{Min: image.Point{X: 46, Y: 4}, Max: image.Point{X: 50, Y: 8}},
			{Min: image.Point{X: 46, Y: 24}, Max: image.Point{X: 50, Y: 28}},
			{Min: image.Point{X: 46, Y: 46}, Max: image.Point{X: 50, Y: 50}},
			{Min: image.Point{X: 46, Y: 68}, Max: image.Point{X: 50, Y: 72}},
			{Min: image.Point{X: 68, Y: 24}, Max: image.Point{X: 72, Y: 28}},
			{Min: image.Point{X: 68, Y: 46}, Max: image.Point{X: 72, Y: 50}},
			{Min: image.Point{X: 68, Y: 68}, Max: image.Point{X: 72, Y: 72}},
		},
		{
			{Min: image.Point{X: 4, Y: 24}, Max: image.Point{X: 8, Y: 28}},
			{Min: image.Point{X: 4, Y: 48}, Max: image.Point{X: 8, Y: 52}},
			{Min: image.Point{X: 24, Y: 4}, Max: image.Point{X: 28, Y: 8}},
			{Min: image.Point{X: 24, Y: 24}, Max: image.Point{X: 28, Y: 28}},
			{Min: image.Point{X: 24, Y: 48}, Max: image.Point{X: 28, Y: 52}},
			{Min: image.Point{X: 24, Y: 72}, Max: image.Point{X: 28, Y: 76}},
			{Min: image.Point{X: 48, Y: 4}, Max: image.Point{X: 52, Y: 8}},
			{Min: image.Point{X: 48, Y: 24}, Max: image.Point{X: 52, Y: 28}},
			{Min: image.Point{X: 48, Y: 48}, Max: image.Point{X: 52, Y: 52}},
			{Min: image.Point{X: 48, Y: 72}, Max: image.Point{X: 52, Y: 76}},
			{Min: image.Point{X: 72, Y: 24}, Max: image.Point{X: 76, Y: 28}},
			{Min: image.Point{X: 72, Y: 48}, Max: image.Point{X: 76, Y: 52}},
			{Min: image.Point{X: 72, Y: 72}, Max: image.Point{X: 76, Y: 76}},
		},
		{
			{Min: image.Point{X: 4, Y: 28}, Max: image.Point{X: 8, Y: 32}},
			{Min: image.Point{X: 4, Y: 52}, Max: image.Point{X: 8, Y: 56}},
			{Min: image.Point{X: 28, Y: 4}, Max: image.Point{X: 32, Y: 8}},
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
			{Min: image.Point{X: 28, Y: 52}, Max: image.Point{X: 32, Y: 56}},
			{Min: image.Point{X: 28, Y: 76}, Max: image.Point{X: 32, Y: 80}},
			{Min: image.Point{X: 52, Y: 4}, Max: image.Point{X: 56, Y: 8}},
			{Min: image.Point{X: 52, Y: 28}, Max: image.Point{X: 56, Y: 32}},
			{Min: image.Point{X: 52, Y: 52}, Max: image.Point{X: 56, Y: 56}},
			{Min: image.Point{X: 52, Y: 76}, Max: image.Point{X: 56, Y: 80}},
			{Min: image.Point{X: 76, Y: 28}, Max: image.Point{X: 80, Y: 32}},
			{Min: image.Point{X: 76, Y: 52}, Max: image.Point{X: 80, Y: 56}},
			{Min: image.Point{X: 76, Y: 76}, Max: image.Point{X: 80, Y: 80}},
		},
		{
			{Min: image.Point{X: 4, Y: 28}, Max: image.Point{X: 8, Y: 32}},
			{Min: image.Point{X: 4, Y: 54}, Max: image.Point{X: 8, Y: 58}},
			{Min: image.Point{X: 28, Y: 4}, Max: image.Point{X: 32, Y: 8}},
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
			{Min: image.Point{X: 28, Y: 54}, Max: image.Point{X: 32, Y: 58}},
			{Min: image.Point{X: 28, Y: 80}, Max: image.Point{X: 32, Y: 84}},
			{Min: image.Point{X: 54, Y: 4}, Max: image.Point{X: 58, Y: 8}},
			{Min: image.Point{X: 54, Y: 28}, Max: image.Point{X: 58, Y: 32}},
			{Min: image.Point{X: 54, Y: 54}, Max: image.Point{X: 58, Y: 58}},
			{Min: image.Point{X: 54, Y: 80}, Max: image.Point{X: 58, Y: 84}},
			{Min: image.Point{X: 80, Y: 28}, Max: image.Point{X: 84, Y: 32}},
			{Min: image.Point{X: 80, Y: 54}, Max: image.Point{X: 84, Y: 58}},
			{Min: image.Point{X: 80, Y: 80}, Max: image.Point{X: 84, Y: 84}},
		},
		{
			{Min: image.Point{X: 4, Y: 28}, Max: image.Point{X: 8, Y: 32}},
			{Min: image.Point{X: 4, Y: 56}, Max: image.Point{X: 8, Y: 60}},
			{Min: image.Point{X: 28, Y: 4}, Max: image.Point{X: 32, Y: 8}},
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
			{Min: image.Point{X: 28, Y: 56}, Max: image.Point{X: 32, Y: 60}},
			{Min: image.Point{X: 28, Y: 84}, Max: image.Point{X: 32, Y: 88}},
			{Min: image.Point{X: 56, Y: 4}, Max: image.Point{X: 60, Y: 8}},
			{Min: image.Point{X: 56, Y: 28}, Max: image.Point{X: 60, Y: 32}},
			{Min: image.Point{X: 56, Y: 56}, Max: image.Point{X: 60, Y: 60}},
			{Min: image.Point{X: 56, Y: 84}, Max: image.Point{X: 60, Y: 88}},
			{Min: image.Point{X: 84, Y: 28}, Max: image.Point{X: 88, Y: 32}},
			{Min: image.Point{X: 84, Y: 56}, Max: image.Point{X: 88, Y: 60}},
			{Min: image.Point{X: 84, Y: 84}, Max: image.Point{X: 88, Y: 88}},
		},
		{
			{Min: image.Point{X: 4, Y: 32}, Max: image.Point{X: 8, Y: 36}},
			{Min: image.Point{X: 4, Y: 60}, Max: image.Point{X: 8, Y: 64}},
			{Min: image.Point{X: 32, Y: 4}, Max: image.Point{X: 36, Y: 8}},
			{Min: image.Point{X: 32, Y: 32}, Max: image.Point{X: 36, Y: 36}},
			{Min: image.Point{X: 32, Y: 60}, Max: image.Point{X: 36, Y: 64}},
			{Min: image.Point{X: 32, Y: 88}, Max: image.Point{X: 36, Y: 92}},
			{Min: image.Point{X: 60, Y: 4}, Max: image.Point{X: 64, Y: 8}},
			{Min: image.Point{X: 60, Y: 32}, Max: image.Point{X: 64, Y: 36}},
			{Min: image.Point{X: 60, Y: 60}, Max: image.Point{X: 64, Y: 64}},
			{Min: image.Point{X: 60, Y: 88}, Max: image.Point{X: 64, Y: 92}},
			{Min: image.Point{X: 88, Y: 32}, Max: image.Point{X: 92, Y: 36}},
			{Min: image.Point{X: 88, Y: 60}, Max: image.Point{X: 92, Y: 64}},
			{Min: image.Point{X: 88, Y: 88}, Max: image.Point{X: 92, Y: 92}},
		},
		{
			{Min: image.Point{X: 4, Y: 26}, Max: image.Point{X: 8, Y: 30}},
			{Min: image.Point{X: 4, Y: 48}, Max: image.Point{X: 8, Y: 52}},
			{Min: image.Point{X: 4, Y: 70}, Max: image.Point{X: 8, Y: 74}},
			{Min: image.Point{X: 26, Y: 4}, Max: image.Point{X: 30, Y: 8}},
			{Min: image.Point{X: 26, Y: 26}, Max: image.Point{X: 30, Y: 30}},
			{Min: image.Point{X: 26, Y: 48}, Max: image.Point{X: 30, Y: 52}},
			{Min: image.Point{X: 26, Y: 70}, Max: image.Point{X: 30, Y: 74}},
			{Min: image.Point{X: 26, Y: 92}, Max: image.Point{X: 30, Y: 96}},
			{Min: image.Point{X: 48, Y: 4}, Max: image.Point{X: 52, Y: 8}},
			{Min: image.Point{X: 48, Y: 26}, Max: image.Point{X: 52, Y: 30}},
			{Min: image.Point{X: 48, Y: 48}, Max: image.Point{X: 52, Y: 52}},
			{Min: image.Point{X: 48, Y: 70}, Max: image.Point{X: 52, Y: 74}},
			{Min: image.Point{X: 48, Y: 92}, Max: image.Point{X: 52, Y: 96}},
			{Min: image.Point{X: 70, Y: 4}, Max: image.Point{X: 74, Y: 8}},
			{Min: image.Point{X: 70, Y: 26}, Max: image.Point{X: 74, Y: 30}},
			{Min: image.Point{X: 70, Y: 48}, Max: image.Point{X: 74, Y: 52}},
			{Min: image.Point{X: 70, Y: 70}, Max: image.Point{X: 74, Y: 74}},
			{Min: image.Point{X: 70, Y: 92}, Max: image.Point{X: 74, Y: 96}},
			{Min: image.Point{X: 92, Y: 26}, Max: image.Point{X: 96, Y: 30}},
			{Min: image.Point{X: 92, Y: 48}, Max: image.Point{X: 96, Y: 52}},
			{Min: image.Point{X: 92, Y: 70}, Max: image.Point{X: 96, Y: 74}},
			{Min: image.Point{X: 92, Y: 92}, Max: image.Point{X: 96, Y: 96}},
		},
		{
			{Min: image.Point{X: 4, Y: 24}, Max: image.Point{X: 8, Y: 28}},
			{Min: image.Point{X: 4, Y: 48}, Max: image.Point{X: 8, Y: 52}},
			{Min: image.Point{X: 4, Y: 72}, Max: image.Point{X: 8, Y: 76}},
			{Min: image.Point{X: 24, Y: 4}, Max: image.Point{X: 28, Y: 8}},
			{Min: image.Point{X: 24, Y: 24}, Max: image.Point{X: 28, Y: 28}},
			{Min: image.Point{X: 24, Y: 48}, Max: image.Point{X: 28, Y: 52}},
			{Min: image.Point{X: 24, Y: 72}, Max: image.Point{X: 28, Y: 76}},
			{Min: image.Point{X: 24, Y: 96}, Max: image.Point{X: 28, Y: 100}},
			{Min: image.Point{X: 48, Y: 4}, Max: image.Point{X: 52, Y: 8}},
			{Min: image.Point{X: 48, Y: 24}, Max: image.Point{X: 52, Y: 28}},
			{Min: image.Point{X: 48, Y: 48}, Max: image.Point{X: 52, Y: 52}},
			{Min: image.Point{X: 48, Y: 72}, Max: image.Point{X: 52, Y: 76}},
			{Min: image.Point{X: 48, Y: 96}, Max: image.Point{X: 52, Y: 100}},
			{Min: image.Point{X: 72, Y: 4}, Max: image.Point{X: 76, Y: 8}},
			{Min: image.Point{X: 72, Y: 24}, Max: image.Point{X: 76, Y: 28}},
			{Min: image.Point{X: 72, Y: 48}, Max: image.Point{X: 76, Y: 52}},
			{Min: image.Point{X: 72, Y: 72}, Max: image.Point{X: 76, Y: 76}},
			{Min: image.Point{X: 72, Y: 96}, Max: image.Point{X: 76, Y: 100}},
			{Min: image.Point{X: 96, Y: 24}, Max: image.Point{X: 100, Y: 28}},
			{Min: image.Point{X: 96, Y: 48}, Max: image.Point{X: 100, Y: 52}},
			{Min: image.Point{X: 96, Y: 72}, Max: image.Point{X: 100, Y: 76}},
			{Min: image.Point{X: 96, Y: 96}, Max: image.Point{X: 100, Y: 100}},
		},
		{
			{Min: image.Point{X: 4, Y: 28}, Max: image.Point{X: 8, Y: 32}},
			{Min: image.Point{X: 4, Y: 52}, Max: image.Point{X: 8, Y: 56}},
			{Min: image.Point{X: 4, Y: 76}, Max: image.Point{X: 8, Y: 80}},
			{Min: image.Point{X: 28, Y: 4}, Max: image.Point{X: 32, Y: 8}},
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
			{Min: image.Point{X: 28, Y: 52}, Max: image.Point{X: 32, Y: 56}},
			{Min: image.Point{X: 28, Y: 76}, Max: image.Point{X: 32, Y: 80}},
			{Min: image.Point{X: 28, Y: 100}, Max: image.Point{X: 32, Y: 104}},
			{Min: image.Point{X: 52, Y: 4}, Max: image.Point{X: 56, Y: 8}},
			{Min: image.Point{X: 52, Y: 28}, Max: image.Point{X: 56, Y: 32}},
			{Min: image.Point{X: 52, Y: 52}, Max: image.Point{X: 56, Y: 56}},
			{Min: image.Point{X: 52, Y: 76}, Max: image.Point{X: 56, Y: 80}},
			{Min: image.Point{X: 52, Y: 100}, Max: image.Point{X: 56, Y: 104}},
			{Min: image.Point{X: 76, Y: 4}, Max: image.Point{X: 80, Y: 8}},
			{Min: image.Point{X: 76, Y: 28}, Max: image.Point{X: 80, Y: 32}},
			{Min: image.Point{X: 76, Y: 52}, Max: image.Point{X: 80, Y: 56}},
			{Min: image.Point{X: 76, Y: 76}, Max: image.Point{X: 80, Y: 80}},
			{Min: image.Point{X: 76, Y: 100}, Max: image.Point{X: 80, Y: 104}},
			{Min: image.Point{X: 100, Y: 28}, Max: image.Point{X: 104, Y: 32}},
			{Min: image.Point{X: 100, Y: 52}, Max: image.Point{X: 104, Y: 56}},
			{Min: image.Point{X: 100, Y: 76}, Max: image.Point{X: 104, Y: 80}},
			{Min: image.Point{X: 100, Y: 100}, Max: image.Point{X: 104, Y: 104}},
		},
		{
			{Min: image.Point{X: 4, Y: 26}, Max: image.Point{X: 8, Y: 30}},
			{Min: image.Point{X: 4, Y: 52}, Max: image.Point{X: 8, Y: 56}},
			{Min: image.Point{X: 4, Y: 78}, Max: image.Point{X: 8, Y: 82}},
			{Min: image.Point{X: 26, Y: 4}, Max: image.Point{X: 30, Y: 8}},
			{Min: image.Point{X: 26, Y: 26}, Max: image.Point{X: 30, Y: 30}},
			{Min: image.Point{X: 26, Y: 52}, Max: image.Point{X: 30, Y: 56}},
			{Min: image.Point{X: 26, Y: 78}, Max: image.Point{X: 30, Y: 82}},
			{Min: image.Point{X: 26, Y: 104}, Max: image.Point{X: 30, Y: 108}},
			{Min: image.Point{X: 52, Y: 4}, Max: image.Point{X: 56, Y: 8}},
			{Min: image.Point{X: 52, Y: 26}, Max: image.Point{X: 56, Y: 30}},
			{Min: image.Point{X: 52, Y: 52}, Max: image.Point{X: 56, Y: 56}},
			{Min: image.Point{X: 52, Y: 78}, Max: image.Point{X: 56, Y: 82}},
			{Min: image.Point{X: 52, Y: 104}, Max: image.Point{X: 56, Y: 108}},
			{Min: image.Point{X: 78, Y: 4}, Max: image.Point{X: 82, Y: 8}},
			{Min: image.Point{X: 78, Y: 26}, Max: image.Point{X: 82, Y: 30}},
			{Min: image.Point{X: 78, Y: 52}, Max: image.Point{X: 82, Y: 56}},
			{Min: image.Point{X: 78, Y: 78}, Max: image.Point{X: 82, Y: 82}},
			{Min: image.Point{X: 78, Y: 104}, Max: image.Point{X: 82, Y: 108}},
			{Min: image.Point{X: 104, Y: 26}, Max: image.Point{X: 108, Y: 30}},
			{Min: image.Point{X: 104, Y: 52}, Max: image.Point{X: 108, Y: 56}},
			{Min: image.Point{X: 104, Y: 78}, Max: image.Point{X: 108, Y: 82}},
			{Min: image.Point{X: 104, Y: 104}, Max: image.Point{X: 108, Y: 108}},
		},
		{
			{Min: image.Point{X: 4, Y: 30}, Max: image.Point{X: 8, Y: 34}},
			{Min: image.Point{X: 4, Y: 56}, Max: image.Point{X: 8, Y: 60}},
			{Min: image.Point{X: 4, Y: 82}, Max: image.Point{X: 8, Y: 86}},
			{Min: image.Point{X: 30, Y: 4}, Max: image.Point{X: 34, Y: 8}},
			{Min: image.Point{X: 30, Y: 30}, Max: image.Point{X: 34, Y: 34}},
			{Min: image.Point{X: 30, Y: 56}, Max: image.Point{X: 34, Y: 60}},
			{Min: image.Point{X: 30, Y: 82}, Max: image.Point{X: 34, Y: 86}},
			{Min: image.Point{X: 30, Y: 108}, Max: image.Point{X: 34, Y: 112}},
			{Min: image.Point{X: 56, Y: 4}, Max: image.Point{X: 60, Y: 8}},
			{Min: image.Point{X: 56, Y: 30}, Max: image.Point{X: 60, Y: 34}},
			{Min: image.Point{X: 56, Y: 56}, Max: image.Point{X: 60, Y: 60}},
			{Min: image.Point{X: 56, Y: 82}, Max: image.Point{X: 60, Y: 86}},
			{Min: image.Point{X: 56, Y: 108}, Max: image.Point{X: 60, Y: 112}},
			{Min: image.Point{X: 82, Y: 4}, Max: image.Point{X: 86, Y: 8}},
			{Min: image.Point{X: 82, Y: 30}, Max: image.Point{X: 86, Y: 34}},
			{Min: image.Point{X: 82, Y: 56}, Max: image.Point{X: 86, Y: 60}},
			{Min: image.Point{X: 82, Y: 82}, Max: image.Point{X: 86, Y: 86}},
			{Min: image.Point{X: 82, Y: 108}, Max: image.Point{X: 86, Y: 112}},
			{Min: image.Point{X: 108, Y: 30}, Max: image.Point{X: 112, Y: 34}},
			{Min: image.Point{X: 108, Y: 56}, Max: image.Point{X: 112, Y: 60}},
			{Min: image.Point{X: 108, Y: 82}, Max: image.Point{X: 112, Y: 86}},
			{Min: image.Point{X: 108, Y: 108}, Max: image.Point{X: 112, Y: 112}},
		},
		{
			{Min: image.Point{X: 4, Y: 28}, Max: image.Point{X: 8, Y: 32}},
			{Min: image.Point{X: 4, Y: 56}, Max: image.Point{X: 8, Y: 60}},
			{Min: image.Point{X: 4, Y: 84}, Max: image.Point{X: 8, Y: 88}},
			{Min: image.Point{X: 28, Y: 4}, Max: image.Point{X: 32, Y: 8}},
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
			{Min: image.Point{X: 28, Y: 56}, Max: image.Point{X: 32, Y: 60}},
			{Min: image.Point{X: 28, Y: 84}, Max: image.Point{X: 32, Y: 88}},
			{Min: image.Point{X: 28, Y: 112}, Max: image.Point{X: 32, Y: 116}},
			{Min: image.Point{X: 56, Y: 4}, Max: image.Point{X: 60, Y: 8}},
			{Min: image.Point{X: 56, Y: 28}, Max: image.Point{X: 60, Y: 32}},
			{Min: image.Point{X: 56, Y: 56}, Max: image.Point{X: 60, Y: 60}},
			{Min: image.Point{X: 56, Y: 84}, Max: image.Point{X: 60, Y: 88}},
			{Min: image.Point{X: 56, Y: 112}, Max: image.Point{X: 60, Y: 116}},
			{Min: image.Point{X: 84, Y: 4}, Max: image.Point{X: 88, Y: 8}},
			{Min: image.Point{X: 84, Y: 28}, Max: image.Point{X: 88, Y: 32}},
			{Min: image.Point{X: 84, Y: 56}, Max: image.Point{X: 88, Y: 60}},
			{Min: image.Point{X: 84, Y: 84}, Max: image.Point{X: 88, Y: 88}},
			{Min: image.Point{X: 84, Y: 112}, Max: image.Point{X: 88, Y: 116}},
			{Min: image.Point{X: 112, Y: 28}, Max: image.Point{X: 116, Y: 32}},
			{Min: image.Point{X: 112, Y: 56}, Max: image.Point{X: 116, Y: 60}},
			{Min: image.Point{X: 112, Y: 84}, Max: image.Point{X: 116, Y: 88}},
			{Min: image.Point{X: 112, Y: 112}, Max: image.Point{X: 116, Y: 116}},
		},
		{
			{Min: image.Point{X: 4, Y: 32}, Max: image.Point{X: 8, Y: 36}},
			{Min: image.Point{X: 4, Y: 60}, Max: image.Point{X: 8, Y: 64}},
			{Min: image.Point{X: 4, Y: 88}, Max: image.Point{X: 8, Y: 92}},
			{Min: image.Point{X: 32, Y: 4}, Max: image.Point{X: 36, Y: 8}},
			{Min: image.Point{X: 32, Y: 32}, Max: image.Point{X: 36, Y: 36}},
			{Min: image.Point{X: 32, Y: 60}, Max: image.Point{X: 36, Y: 64}},
			{Min: image.Point{X: 32, Y: 88}, Max: image.Point{X: 36, Y: 92}},
			{Min: image.Point{X: 32, Y: 116}, Max: image.Point{X: 36, Y: 120}},
			{Min: image.Point{X: 60, Y: 4}, Max: image.Point{X: 64, Y: 8}},
			{Min: image.Point{X: 60, Y: 32}, Max: image.Point{X: 64, Y: 36}},
			{Min: image.Point{X: 60, Y: 60}, Max: image.Point{X: 64, Y: 64}},
			{Min: image.Point{X: 60, Y: 88}, Max: image.Point{X: 64, Y: 92}},
			{Min: image.Point{X: 60, Y: 116}, Max: image.Point{X: 64, Y: 120}},
			{Min: image.Point{X: 88, Y: 4}, Max: image.Point{X: 92, Y: 8}},
			{Min: image.Point{X: 88, Y: 32}, Max: image.Point{X: 92, Y: 36}},
			{Min: image.Point{X: 88, Y: 60}, Max: image.Point{X: 92, Y: 64}},
			{Min: image.Point{X: 88, Y: 88}, Max: image.Point{X: 92, Y: 92}},
			{Min: image.Point{X: 88, Y: 116}, Max: image.Point{X: 92, Y: 120}},
			{Min: image.Point{X: 116, Y: 32}, Max: image.Point{X: 120, Y: 36}},
			{Min: image.Point{X: 116, Y: 60}, Max: image.Point{X: 120, Y: 64}},
			{Min: image.Point{X: 116, Y: 88}, Max: image.Point{X: 120, Y: 92}},
			{Min: image.Point{X: 116, Y: 116}, Max: image.Point{X: 120, Y: 120}},
		},
		{
			{Min: image.Point{X: 4, Y: 24}, Max: image.Point{X: 8, Y: 28}},
			{Min: image.Point{X: 4, Y: 48}, Max: image.Point{X: 8, Y: 52}},
			{Min: image.Point{X: 4, Y: 72}, Max: image.Point{X: 8, Y: 76}},
			{Min: image.Point{X: 4, Y: 96}, Max: image.Point{X: 8, Y: 100}},
			{Min: image.Point{X: 24, Y: 4}, Max: image.Point{X: 28, Y: 8}},
			{Min: image.Point{X: 24, Y: 24}, Max: image.Point{X: 28, Y: 28}},
			{Min: image.Point{X: 24, Y: 48}, Max: image.Point{X: 28, Y: 52}},
			{Min: image.Point{X: 24, Y: 72}, Max: image.Point{X: 28, Y: 76}},
			{Min: image.Point{X: 24, Y: 96}, Max: image.Point{X: 28, Y: 100}},
			{Min: image.Point{X: 24, Y: 120}, Max: image.Point{X: 28, Y: 124}},
			{Min: image.Point{X: 48, Y: 4}, Max: image.Point{X: 52, Y: 8}},
			{Min: image.Point{X: 48, Y: 24}, Max: image.Point{X: 52, Y: 28}},
			{Min: image.Point{X: 48, Y: 48}, Max: image.Point{X: 52, Y: 52}},
			{Min: image.Point{X: 48, Y: 72}, Max: image.Point{X: 52, Y: 76}},
			{Min: image.Point{X: 48, Y: 96}, Max: image.Point{X: 52, Y: 100}},
			{Min: image.Point{X: 48, Y: 120}, Max: image.Point{X: 52, Y: 124}},
			{Min: image.Point{X: 72, Y: 4}, Max: image.Point{X: 76, Y: 8}},
			{Min: image.Point{X: 72, Y: 24}, Max: image.Point{X: 76, Y: 28}},
			{Min: image.Point{X: 72, Y: 48}, Max: image.Point{X: 76, Y: 52}},
			{Min: image.Point{X: 72, Y: 72}, Max: image.Point{X: 76, Y: 76}},
			{Min: image.Point{X: 72, Y: 96}, Max: image.Point{X: 76, Y: 100}},
			{Min: image.Point{X: 72, Y: 120}, Max: image.Point{X: 76, Y: 124}},
			{Min: image.Point{X: 96, Y: 4}, Max: image.Point{X: 100, Y: 8}},
			{Min: image.Point{X: 96, Y: 24}, Max: image.Point{X: 100, Y: 28}},
			{Min: image.Point{X: 96, Y: 48}, Max: image.Point{X: 100, Y: 52}},
			{Min: image.Point{X: 96, Y: 72}, Max: image.Point{X: 100, Y: 76}},
			{Min: image.Point{X: 96, Y: 96}, Max: image.Point{X: 100, Y: 100}},
			{Min: image.Point{X: 96, Y: 120}, Max: image.Point{X: 100, Y: 124}},
			{Min: image.Point{X: 120, Y: 24}, Max: image.Point{X: 124, Y: 28}},
			{Min: image.Point{X: 120, Y: 48}, Max: image.Point{X: 124, Y: 52}},
			{Min: image.Point{X: 120, Y: 72}, Max: image.Point{X: 124, Y: 76}},
			{Min: image.Point{X: 120, Y: 96}, Max: image.Point{X: 124, Y: 100}},
			{Min: image.Point{X: 120, Y: 120}, Max: image.Point{X: 124, Y: 124}},
		},
		{
			{Min: image.Point{X: 4, Y: 28}, Max: image.Point{X: 8, Y: 32}},
			{Min: image.Point{X: 4, Y: 52}, Max: image.Point{X: 8, Y: 56}},
			{Min: image.Point{X: 4, Y: 76}, Max: image.Point{X: 8, Y: 80}},
			{Min: image.Point{X: 4, Y: 100}, Max: image.Point{X: 8, Y: 104}},
			{Min: image.Point{X: 28, Y: 4}, Max: image.Point{X: 32, Y: 8}},
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
			{Min: image.Point{X: 28, Y: 52}, Max: image.Point{X: 32, Y: 56}},
			{Min: image.Point{X: 28, Y: 76}, Max: image.Point{X: 32, Y: 80}},
			{Min: image.Point{X: 28, Y: 100}, Max: image.Point{X: 32, Y: 104}},
			{Min: image.Point{X: 28, Y: 124}, Max: image.Point{X: 32, Y: 128}},
			{Min: image.Point{X: 52, Y: 4}, Max: image.Point{X: 56, Y: 8}},
			{Min: image.Point{X: 52, Y: 28}, Max: image.Point{X: 56, Y: 32}},
			{Min: image.Point{X: 52, Y: 52}, Max: image.Point{X: 56, Y: 56}},
			{Min: image.Point{X: 52, Y: 76}, Max: image.Point{X: 56, Y: 80}},
			{Min: image.Point{X: 52, Y: 100}, Max: image.Point{X: 56, Y: 104}},
			{Min: image.Point{X: 52, Y: 124}, Max: image.Point{X: 56, Y: 128}},
			{Min: image.Point{X: 76, Y: 4}, Max: image.Point{X: 80, Y: 8}},
			{Min: image.Point{X: 76, Y: 28}, Max: image.Point{X: 80, Y: 32}},
			{Min: image.Point{X: 76, Y: 52}, Max: image.Point{X: 80, Y: 56}},
			{Min: image.Point{X: 76, Y: 76}, Max: image.Point{X: 80, Y: 80}},
			{Min: image.Point{X: 76, Y: 100}, Max: image.Point{X: 80, Y: 104}},
			{Min: image.Point{X: 76, Y: 124}, Max: image.Point{X: 80, Y: 128}},
			{Min: image.Point{X: 100, Y: 4}, Max: image.Point{X: 104, Y: 8}},
			{Min: image.Point{X: 100, Y: 28}, Max: image.Point{X: 104, Y: 32}},
			{Min: image.Point{X: 100, Y: 52}, Max: image.Point{X: 104, Y: 56}},
			{Min: image.Point{X: 100, Y: 76}, Max: image.Point{X: 104, Y: 80}},
			{Min: image.Point{X: 100, Y: 100}, Max: image.Point{X: 104, Y: 104}},
			{Min: image.Point{X: 100, Y: 124}, Max: image.Point{X: 104, Y: 128}},
			{Min: image.Point{X: 124, Y: 28}, Max: image.Point{X: 128, Y: 32}},
			{Min: image.Point{X: 124, Y: 52}, Max: image.Point{X: 128, Y: 56}},
			{Min: image.Point{X: 124, Y: 76}, Max: image.Point{X: 128, Y: 80}},
			{Min: image.Point{X: 124, Y: 100}, Max: image.Point{X: 128, Y: 104}},
			{Min: image.Point{X: 124, Y: 124}, Max: image.Point{X: 128, Y: 128}},
		},
		{
			{Min: image.Point{X: 4, Y: 24}, Max: image.Point{X: 8, Y: 28}},
			{Min: image.Point{X: 4, Y: 50}, Max: image.Point{X: 8, Y: 54}},
			{Min: image.Point{X: 4, Y: 76}, Max: image.Point{X: 8, Y: 80}},
			{Min: image.Point{X: 4, Y: 102}, Max: image.Point{X: 8, Y: 106}},
			{Min: image.Point{X: 24, Y: 4}, Max: image.Point{X: 28, Y: 8}},
			{Min: image.Point{X: 24, Y: 24}, Max: image.Point{X: 28, Y: 28}},
			{Min: image.Point{X: 24, Y: 50}, Max: image.Point{X: 28, Y: 54}},
			{Min: image.Point{X: 24, Y: 76}, Max: image.Point{X: 28, Y: 80}},
			{Min: image.Point{X: 24, Y: 102}, Max: image.Point{X: 28, Y: 106}},
			{Min: image.Point{X: 24, Y: 128}, Max: image.Point{X: 28, Y: 132}},
			{Min: image.Point{X: 50, Y: 4}, Max: image.Point{X: 54, Y: 8}},
			{Min: image.Point{X: 50, Y: 24}, Max: image.Point{X: 54, Y: 28}},
			{Min: image.Point{X: 50, Y: 50}, Max: image.Point{X: 54, Y: 54}},
			{Min: image.Point{X: 50, Y: 76}, Max: image.Point{X: 54, Y: 80}},
			{Min: image.Point{X: 50, Y: 102}, Max: image.Point{X: 54, Y: 106}},
			{Min: image.Point{X: 50, Y: 128}, Max: image.Point{X: 54, Y: 132}},
			{Min: image.Point{X: 76, Y: 4}, Max: image.Point{X: 80, Y: 8}},
			{Min: image.Point{X: 76, Y: 24}, Max: image.Point{X: 80, Y: 28}},
			{Min: image.Point{X: 76, Y: 50}, Max: image.Point{X: 80, Y: 54}},
			{Min: image.Point{X: 76, Y: 76}, Max: image.Point{X: 80, Y: 80}},
			{Min: image.Point{X: 76, Y: 102}, Max: image.Point{X: 80, Y: 106}},
			{Min: image.Point{X: 76, Y: 128}, Max: image.Point{X: 80, Y: 132}},
			{Min: image.Point{X: 102, Y: 4}, Max: image.Point{X: 106, Y: 8}},
			{Min: image.Point{X: 102, Y: 24}, Max: image.Point{X: 106, Y: 28}},
			{Min: image.Point{X: 102, Y: 50}, Max: image.Point{X: 106, Y: 54}},
			{Min: image.Point{X: 102, Y: 76}, Max: image.Point{X: 106, Y: 80}},
			{Min: image.Point{X: 102, Y: 102}, Max: image.Point{X: 106, Y: 106}},
			{Min: image.Point{X: 102, Y: 128}, Max: image.Point{X: 106, Y: 132}},
			{Min: image.Point{X: 128, Y: 24}, Max: image.Point{X: 132, Y: 28}},
			{Min: image.Point{X: 128, Y: 50}, Max: image.Point{X: 132, Y: 54}},
			{Min: image.Point{X: 128, Y: 76}, Max: image.Point{X: 132, Y: 80}},
			{Min: image.Point{X: 128, Y: 102}, Max: image.Point{X: 132, Y: 106}},
			{Min: image.Point{X: 128, Y: 128}, Max: image.Point{X: 132, Y: 132}},
		},
		{
			{Min: image.Point{X: 4, Y: 28}, Max: image.Point{X: 8, Y: 32}},
			{Min: image.Point{X: 4, Y: 54}, Max: image.Point{X: 8, Y: 58}},
			{Min: image.Point{X: 4, Y: 80}, Max: image.Point{X: 8, Y: 84}},
			{Min: image.Point{X: 4, Y: 106}, Max: image.Point{X: 8, Y: 110}},
			{Min: image.Point{X: 28, Y: 4}, Max: image.Point{X: 32, Y: 8}},
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
			{Min: image.Point{X: 28, Y: 54}, Max: image.Point{X: 32, Y: 58}},
			{Min: image.Point{X: 28, Y: 80}, Max: image.Point{X: 32, Y: 84}},
			{Min: image.Point{X: 28, Y: 106}, Max: image.Point{X: 32, Y: 110}},
			{Min: image.Point{X: 28, Y: 132}, Max: image.Point{X: 32, Y: 136}},
			{Min: image.Point{X: 54, Y: 4}, Max: image.Point{X: 58, Y: 8}},
			{Min: image.Point{X: 54, Y: 28}, Max: image.Point{X: 58, Y: 32}},
			{Min: image.Point{X: 54, Y: 54}, Max: image.Point{X: 58, Y: 58}},
			{Min: image.Point{X: 54, Y: 80}, Max: image.Point{X: 58, Y: 84}},
			{Min: image.Point{X: 54, Y: 106}, Max: image.Point{X: 58, Y: 110}},
			{Min: image.Point{X: 54, Y: 132}, Max: image.Point{X: 58, Y: 136}},
			{Min: image.Point{X: 80, Y: 4}, Max: image.Point{X: 84, Y: 8}},
			{Min: image.Point{X: 80, Y: 28}, Max: image.Point{X: 84, Y: 32}},
			{Min: image.Point{X: 80, Y: 54}, Max: image.Point{X: 84, Y: 58}},
			{Min: image.Point{X: 80, Y: 80}, Max: image.Point{X: 84, Y: 84}},
			{Min: image.Point{X: 80, Y: 106}, Max: image.Point{X: 84, Y: 110}},
			{Min: image.Point{X: 80, Y: 132}, Max: image.Point{X: 84, Y: 136}},
			{Min: image.Point{X: 106, Y: 4}, Max: image.Point{X: 110, Y: 8}},
			{Min: image.Point{X: 106, Y: 28}, Max: image.Point{X: 110, Y: 32}},
			{Min: image.Point{X: 106, Y: 54}, Max: image.Point{X: 110, Y: 58}},
			{Min: image.Point{X: 106, Y: 80}, Max: image.Point{X: 110, Y: 84}},
			{Min: image.Point{X: 106, Y: 106}, Max: image.Point{X: 110, Y: 110}},
			{Min: image.Point{X: 106, Y: 132}, Max: image.Point{X: 110, Y: 136}},
			{Min: image.Point{X: 132, Y: 28}, Max: image.Point{X: 136, Y: 32}},
			{Min: image.Point{X: 132, Y: 54}, Max: image.Point{X: 136, Y: 58}},
			{Min: image.Point{X: 132, Y: 80}, Max: image.Point{X: 136, Y: 84}},
			{Min: image.Point{X: 132, Y: 106}, Max: image.Point{X: 136, Y: 110}},
			{Min: image.Point{X: 132, Y: 132}, Max: image.Point{X: 136, Y: 136}},
		},
		{
			{Min: image.Point{X: 4, Y: 32}, Max: image.Point{X: 8, Y: 36}},
			{Min: image.Point{X: 4, Y: 58}, Max: image.Point{X: 8, Y: 62}},
			{Min: image.Point{X: 4, Y: 84}, Max: image.Point{X: 8, Y: 88}},
			{Min: image.Point{X: 4, Y: 110}, Max: image.Point{X: 8, Y: 114}},
			{Min: image.Point{X: 32, Y: 4}, Max: image.Point{X: 36, Y: 8}},
			{Min: image.Point{X: 32, Y: 32}, Max: image.Point{X: 36, Y: 36}},
			{Min: image.Point{X: 32, Y: 58}, Max: image.Point{X: 36, Y: 62}},
			{Min: image.Point{X: 32, Y: 84}, Max: image.Point{X: 36, Y: 88}},
			{Min: image.Point{X: 32, Y: 110}, Max: image.Point{X: 36, Y: 114}},
			{Min: image.Point{X: 32, Y: 136}, Max: image.Point{X: 36, Y: 140}},
			{Min: image.Point{X: 58, Y: 4}, Max: image.Point{X: 62, Y: 8}},
			{Min: image.Point{X: 58, Y: 32}, Max: image.Point{X: 62, Y: 36}},
			{Min: image.Point{X: 58, Y: 58}, Max: image.Point{X: 62, Y: 62}},
			{Min: image.Point{X: 58, Y: 84}, Max: image.Point{X: 62, Y: 88}},
			{Min: image.Point{X: 58, Y: 110}, Max: image.Point{X: 62, Y: 114}},
			{Min: image.Point{X: 58, Y: 136}, Max: image.Point{X: 62, Y: 140}},
			{Min: image.Point{X: 84, Y: 4}, Max: image.Point{X: 88, Y: 8}},
			{Min: image.Point{X: 84, Y: 32}, Max: image.Point{X: 88, Y: 36}},
			{Min: image.Point{X: 84, Y: 58}, Max: image.Point{X: 88, Y: 62}},
			{Min: image.Point{X: 84, Y: 84}, Max: image.Point{X: 88, Y: 88}},
			{Min: image.Point{X: 84, Y: 110}, Max: image.Point{X: 88, Y: 114}},
			{Min: image.Point{X: 84, Y: 136}, Max: image.Point{X: 88, Y: 140}},
			{Min: image.Point{X: 110, Y: 4}, Max: image.Point{X: 114, Y: 8}},
			{Min: image.Point{X: 110, Y: 32}, Max: image.Point{X: 114, Y: 36}},
			{Min: image.Point{X: 110, Y: 58}, Max: image.Point{X: 114, Y: 62}},
			{Min: image.Point{X: 110, Y: 84}, Max: image.Point{X: 114, Y: 88}},
			{Min: image.Point{X: 110, Y: 110}, Max: image.Point{X: 114, Y: 114}},
			{Min: image.Point{X: 110, Y: 136}, Max: image.Point{X: 114, Y: 140}},
			{Min: image.Point{X: 136, Y: 32}, Max: image.Point{X: 140, Y: 36}},
			{Min: image.Point{X: 136, Y: 58}, Max: image.Point{X: 140, Y: 62}},
			{Min: image.Point{X: 136, Y: 84}, Max: image.Point{X: 140, Y: 88}},
			{Min: image.Point{X: 136, Y: 110}, Max: image.Point{X: 140, Y: 114}},
			{Min: image.Point{X: 136, Y: 136}, Max: image.Point{X: 140, Y: 140}},
		},
		{
			{Min: image.Point{X: 4, Y: 28}, Max: image.Point{X: 8, Y: 32}},
			{Min: image.Point{X: 4, Y: 56}, Max: image.Point{X: 8, Y: 60}},
			{Min: image.Point{X: 4, Y: 84}, Max: image.Point{X: 8, Y: 88}},
			{Min: image.Point{X: 4, Y: 112}, Max: image.Point{X: 8, Y: 116}},
			{Min: image.Point{X: 28, Y: 4}, Max: image.Point{X: 32, Y: 8}},
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
			{Min: image.Point{X: 28, Y: 56}, Max: image.Point{X: 32, Y: 60}},
			{Min: image.Point{X: 28, Y: 84}, Max: image.Point{X: 32, Y: 88}},
			{Min: image.Point{X: 28, Y: 112}, Max: image.Point{X: 32, Y: 116}},
			{Min: image.Point{X: 28, Y: 140}, Max: image.Point{X: 32, Y: 144}},
			{Min: image.Point{X: 56, Y: 4}, Max: image.Point{X: 60, Y: 8}},
			{Min: image.Point{X: 56, Y: 28}, Max: image.Point{X: 60, Y: 32}},
			{Min: image.Point{X: 56, Y: 56}, Max: image.Point{X: 60, Y: 60}},
			{Min: image.Point{X: 56, Y: 84}, Max: image.Point{X: 60, Y: 88}},
			{Min: image.Point{X: 56, Y: 112}, Max: image.Point{X: 60, Y: 116}},
			{Min: image.Point{X: 56, Y: 140}, Max: image.Point{X: 60, Y: 144}},
			{Min: image.Point{X: 84, Y: 4}, Max: image.Point{X: 88, Y: 8}},
			{Min: image.Point{X: 84, Y: 28}, Max: image.Point{X: 88, Y: 32}},
			{Min: image.Point{X: 84, Y: 56}, Max: image.Point{X: 88, Y: 60}},
			{Min: image.Point{X: 84, Y: 84}, Max: image.Point{X: 88, Y: 88}},
			{Min: image.Point{X: 84, Y: 112}, Max: image.Point{X: 88, Y: 116}},
			{Min: image.Point{X: 84, Y: 140}, Max: image.Point{X: 88, Y: 144}},
			{Min: image.Point{X: 112, Y: 4}, Max: image.Point{X: 116, Y: 8}},
			{Min: image.Point{X: 112, Y: 28}, Max: image.Point{X: 116, Y: 32}},
			{Min: image.Point{X: 112, Y: 56}, Max: image.Point{X: 116, Y: 60}},
			{Min: image.Point{X: 112, Y: 84}, Max: image.Point{X: 116, Y: 88}},
			{Min: image.Point{X: 112, Y: 112}, Max: image.Point{X: 116, Y: 116}},
			{Min: image.Point{X: 112, Y: 140}, Max: image.Point{X: 116, Y: 144}},
			{Min: image.Point{X: 140, Y: 28}, Max: image.Point{X: 144, Y: 32}},
			{Min: image.Point{X: 140, Y: 56}, Max: image.Point{X: 144, Y: 60}},
			{Min: image.Point{X: 140, Y: 84}, Max: image.Point{X: 144, Y: 88}},
			{Min: image.Point{X: 140, Y: 112}, Max: image.Point{X: 144, Y: 116}},
			{Min: image.Point{X: 140, Y: 140}, Max: image.Point{X: 144, Y: 144}},
		},
		{
			{Min: image.Point{X: 4, Y: 32}, Max: image.Point{X: 8, Y: 36}},
			{Min: image.Point{X: 4, Y: 60}, Max: image.Point{X: 8, Y: 64}},
			{Min: image.Point{X: 4, Y: 88}, Max: image.Point{X: 8, Y: 92}},
			{Min: image.Point{X: 4, Y: 116}, Max: image.Point{X: 8, Y: 120}},
			{Min: image.Point{X: 32, Y: 4}, Max: image.Point{X: 36, Y: 8}},
			{Min: image.Point{X: 32, Y: 32}, Max: image.Point{X: 36, Y: 36}},
			{Min: image.Point{X: 32, Y: 60}, Max: image.Point{X: 36, Y: 64}},
			{Min: image.Point{X: 32, Y: 88}, Max: image.Point{X: 36, Y: 92}},
			{Min: image.Point{X: 32, Y: 116}, Max: image.Point{X: 36, Y: 120}},
			{Min: image.Point{X: 32, Y: 144}, Max: image.Point{X: 36, Y: 148}},
			{Min: image.Point{X: 60, Y: 4}, Max: image.Point{X: 64, Y: 8}},
			{Min: image.Point{X: 60, Y: 32}, Max: image.Point{X: 64, Y: 36}},
			{Min: image.Point{X: 60, Y: 60}, Max: image.Point{X: 64, Y: 64}},
			{Min: image.Point{X: 60, Y: 88}, Max: image.Point{X: 64, Y: 92}},
			{Min: image.Point{X: 60, Y: 116}, Max: image.Point{X: 64, Y: 120}},
			{Min: image.Point{X: 60, Y: 144}, Max: image.Point{X: 64, Y: 148}},
			{Min: image.Point{X: 88, Y: 4}, Max: image.Point{X: 92, Y: 8}},
			{Min: image.Point{X: 88, Y: 32}, Max: image.Point{X: 92, Y: 36}},
			{Min: image.Point{X: 88, Y: 60}, Max: image.Point{X: 92, Y: 64}},
			{Min: image.Point{X: 88, Y: 88}, Max: image.Point{X: 92, Y: 92}},
			{Min: image.Point{X: 88, Y: 116}, Max: image.Point{X: 92, Y: 120}},
			{Min: image.Point{X: 88, Y: 144}, Max: image.Point{X: 92, Y: 148}},
			{Min: image.Point{X: 116, Y: 4}, Max: image.Point{X: 120, Y: 8}},
			{Min: image.Point{X: 116, Y: 32}, Max: image.Point{X: 120, Y: 36}},
			{Min: image.Point{X: 116, Y: 60}, Max: image.Point{X: 120, Y: 64}},
			{Min: image.Point{X: 116, Y: 88}, Max: image.Point{X: 120, Y: 92}},
			{Min: image.Point{X: 116, Y: 116}, Max: image.Point{X: 120, Y: 120}},
			{Min: image.Point{X: 116, Y: 144}, Max: image.Point{X: 120, Y: 148}},
			{Min: image.Point{X: 144, Y: 32}, Max: image.Point{X: 148, Y: 36}},
			{Min: image.Point{X: 144, Y: 60}, Max: image.Point{X: 148, Y: 64}},
			{Min: image.Point{X: 144, Y: 88}, Max: image.Point{X: 148, Y: 92}},
			{Min: image.Point{X: 144, Y: 116}, Max: image.Point{X: 148, Y: 120}},
			{Min: image.Point{X: 144, Y: 144}, Max: image.Point{X: 148, Y: 148}},
		},
		{
			{Min: image.Point{X: 4, Y: 28}, Max: image.Point{X: 8, Y: 32}},
			{Min: image.Point{X: 4, Y: 52}, Max: image.Point{X: 8, Y: 56}},
			{Min: image.Point{X: 4, Y: 76}, Max: image.Point{X: 8, Y: 80}},
			{Min: image.Point{X: 4, Y: 100}, Max: image.Point{X: 8, Y: 104}},
			{Min: image.Point{X: 4, Y: 124}, Max: image.Point{X: 8, Y: 128}},
			{Min: image.Point{X: 28, Y: 4}, Max: image.Point{X: 32, Y: 8}},
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
			{Min: image.Point{X: 28, Y: 52}, Max: image.Point{X: 32, Y: 56}},
			{Min: image.Point{X: 28, Y: 76}, Max: image.Point{X: 32, Y: 80}},
			{Min: image.Point{X: 28, Y: 100}, Max: image.Point{X: 32, Y: 104}},
			{Min: image.Point{X: 28, Y: 124}, Max: image.Point{X: 32, Y: 128}},
			{Min: image.Point{X: 28, Y: 148}, Max: image.Point{X: 32, Y: 152}},
			{Min: image.Point{X: 52, Y: 4}, Max: image.Point{X: 56, Y: 8}},
			{Min: image.Point{X: 52, Y: 28}, Max: image.Point{X: 56, Y: 32}},
			{Min: image.Point{X: 52, Y: 52}, Max: image.Point{X: 56, Y: 56}},
			{Min: image.Point{X: 52, Y: 76}, Max: image.Point{X: 56, Y: 80}},
			{Min: image.Point{X: 52, Y: 100}, Max: image.Point{X: 56, Y: 104}},
			{Min: image.Point{X: 52, Y: 124}, Max: image.Point{X: 56, Y: 128}},
			{Min: image.Point{X: 52, Y: 148}, Max: image.Point{X: 56, Y: 152}},
			{Min: image.Point{X: 76, Y: 4}, Max: image.Point{X: 80, Y: 8}},
			{Min: image.Point{X: 76, Y: 28}, Max: image.Point{X: 80, Y: 32}},
			{Min: image.Point{X: 76, Y: 52}, Max: image.Point{X: 80, Y: 56}},
			{Min: image.Point{X: 76, Y: 76}, Max: image.Point{X: 80, Y: 80}},
			{Min: image.Point{X: 76, Y: 100}, Max: image.Point{X: 80, Y: 104}},
			{Min: image.Point{X: 76, Y: 124}, Max: image.Point{X: 80, Y: 128}},
			{Min: image.Point{X: 76, Y: 148}, Max: image.Point{X: 80, Y: 152}},
			{Min: image.Point{X: 100, Y: 4}, Max: image.Point{X: 104, Y: 8}},
			{Min: image.Point{X: 100, Y: 28}, Max: image.Point{X: 104, Y: 32}},
			{Min: image.Point{X: 100, Y: 52}, Max: image.Point{X: 104, Y: 56}},
			{Min: image.Point{X: 100, Y: 76}, Max: image.Point{X: 104, Y: 80}},
			{Min: image.Point{X: 100, Y: 100}, Max: image.Point{X: 104, Y: 104}},
			{Min: image.Point{X: 100, Y: 124}, Max: image.Point{X: 104, Y: 128}},
			{Min: image.Point{X: 100, Y: 148}, Max: image.Point{X: 104, Y: 152}},
			{Min: image.Point{X: 124, Y: 4}, Max: image.Point{X: 128, Y: 8}},
			{Min: image.Point{X: 124, Y: 28}, Max: image.Point{X: 128, Y: 32}},
			{Min: image.Point{X: 124, Y: 52}, Max: image.Point{X: 128, Y: 56}},
			{Min: image.Point{X: 124, Y: 76}, Max: image.Point{X: 128, Y: 80}},
			{Min: image.Point{X: 124, Y: 100}, Max: image.Point{X: 128, Y: 104}},
			{Min: image.Point{X: 124, Y: 124}, Max: image.Point{X: 128, Y: 128}},
			{Min: image.Point{X: 124, Y: 148}, Max: image.Point{X: 128, Y: 152}},
			{Min: image.Point{X: 148, Y: 28}, Max: image.Point{X: 152, Y: 32}},
			{Min: image.Point{X: 148, Y: 52}, Max: image.Point{X: 152, Y: 56}},
			{Min: image.Point{X: 148, Y: 76}, Max: image.Point{X: 152, Y: 80}},
			{Min: image.Point{X: 148, Y: 100}, Max: image.Point{X: 152, Y: 104}},
			{Min: image.Point{X: 148, Y: 124}, Max: image.Point{X: 152, Y: 128}},
			{Min: image.Point{X: 148, Y: 148}, Max: image.Point{X: 152, Y: 152}},
		},
		{
			{Min: image.Point{X: 4, Y: 22}, Max: image.Point{X: 8, Y: 26}},
			{Min: image.Point{X: 4, Y: 48}, Max: image.Point{X: 8, Y: 52}},
			{Min: image.Point{X: 4, Y: 74}, Max: image.Point{X: 8, Y: 78}},
			{Min: image.Point{X: 4, Y: 100}, Max: image.Point{X: 8, Y: 104}},
			{Min: image.Point{X: 4, Y: 126}, Max: image.Point{X: 8, Y: 130}},
			{Min: image.Point{X: 22, Y: 4}, Max: image.Point{X: 26, Y: 8}},
			{Min: image.Point{X: 22, Y: 22}, Max: image.Point{X: 26, Y: 26}},
			{Min: image.Point{X: 22, Y: 48}, Max: image.Point{X: 26, Y: 52}},
			{Min: image.Point{X: 22, Y: 74}, Max: image.Point{X: 26, Y: 78}},
			{Min: image.Point{X: 22, Y: 100}, Max: image.Point{X: 26, Y: 104}},
			{Min: image.Point{X: 22, Y: 126}, Max: image.Point{X: 26, Y: 130}},
			{Min: image.Point{X: 22, Y: 152}, Max: image.Point{X: 26, Y: 156}},
			{Min: image.Point{X: 48, Y: 4}, Max: image.Point{X: 52, Y: 8}},
			{Min: image.Point{X: 48, Y: 22}, Max: image.Point{X: 52, Y: 26}},
			{Min: image.Point{X: 48, Y: 48}, Max: image.Point{X: 52, Y: 52}},
			{Min: image.Point{X: 48, Y: 74}, Max: image.Point{X: 52, Y: 78}},
			{Min: image.Point{X: 48, Y: 100}, Max: image.Point{X: 52, Y: 104}},
			{Min: image.Point{X: 48, Y: 126}, Max: image.Point{X: 52, Y: 130}},
			{Min: image.Point{X: 48, Y: 152}, Max: image.Point{X: 52, Y: 156}},
			{Min: image.Point{X: 74, Y: 4}, Max: image.Point{X: 78, Y: 8}},
			{Min: image.Point{X: 74, Y: 22}, Max: image.Point{X: 78, Y: 26}},
			{Min: image.Point{X: 74, Y: 48}, Max: image.Point{X: 78, Y: 52}},
			{Min: image.Point{X: 74, Y: 74}, Max: image.Point{X: 78, Y: 78}},
			{Min: image.Point{X: 74, Y: 100}, Max: image.Point{X: 78, Y: 104}},
			{Min: image.Point{X: 74, Y: 126}, Max: image.Point{X: 78, Y: 130}},
			{Min: image.Point{X: 74, Y: 152}, Max: image.Point{X: 78, Y: 156}},
			{Min: image.Point{X: 100, Y: 4}, Max: image.Point{X: 104, Y: 8}},
			{Min: image.Point{X: 100, Y: 22}, Max: image.Point{X: 104, Y: 26}},
			{Min: image.Point{X: 100, Y: 48}, Max: image.Point{X: 104, Y: 52}},
			{Min: image.Point{X: 100, Y: 74}, Max: image.Point{X: 104, Y: 78}},
			{Min: image.Point{X: 100, Y: 100}, Max: image.Point{X: 104, Y: 104}},
			{Min: image.Point{X: 100, Y: 126}, Max: image.Point{X: 104, Y: 130}},
			{Min: image.Point{X: 100, Y: 152}, Max: image.Point{X: 104, Y: 156}},
			{Min: image.Point{X: 126, Y: 4}, Max: image.Point{X: 130, Y: 8}},
			{Min: image.Point{X: 126, Y: 22}, Max: image.Point{X: 130, Y: 26}},
			{Min: image.Point{X: 126, Y: 48}, Max: image.Point{X: 130, Y: 52}},
			{Min: image.Point{X: 126, Y: 74}, Max: image.Point{X: 130, Y: 78}},
			{Min: image.Point{X: 126, Y: 100}, Max: image.Point{X: 130, Y: 104}},
			{Min: image.Point{X: 126, Y: 126}, Max: image.Point{X: 130, Y: 130}},
			{Min: image.Point{X: 126, Y: 152}, Max: image.Point{X: 130, Y: 156}},
			{Min: image.Point{X: 152, Y: 22}, Max: image.Point{X: 156, Y: 26}},
			{Min: image.Point{X: 152, Y: 48}, Max: image.Point{X: 156, Y: 52}},
			{Min: image.Point{X: 152, Y: 74}, Max: image.Point{X: 156, Y: 78}},
			{Min: image.Point{X: 152, Y: 100}, Max: image.Point{X: 156, Y: 104}},
			{Min: image.Point{X: 152, Y: 126}, Max: image.Point{X: 156, Y: 130}},
			{Min: image.Point{X: 152, Y: 152}, Max: image.Point{X: 156, Y: 156}},
		},
		{
			{Min: image.Point{X: 4, Y: 26}, Max: image.Point{X: 8, Y: 30}},
			{Min: image.Point{X: 4, Y: 52}, Max: image.Point{X: 8, Y: 56}},
			{Min: image.Point{X: 4, Y: 78}, Max: image.Point{X: 8, Y: 82}},
			{Min: image.Point{X: 4, Y: 104}, Max: image.Point{X: 8, Y: 108}},
			{Min: image.Point{X: 4, Y: 130}, Max: image.Point{X: 8, Y: 134}},
			{Min: image.Point{X: 26, Y: 4}, Max: image.Point{X: 30, Y: 8}},
			{Min: image.Point{X: 26, Y: 26}, Max: image.Point{X: 30, Y: 30}},
			{Min: image.Point{X: 26, Y: 52}, Max: image.Point{X: 30, Y: 56}},
			{Min: image.Point{X: 26, Y: 78}, Max: image.Point{X: 30, Y: 82}},
			{Min: image.Point{X: 26, Y: 104}, Max: image.Point{X: 30, Y: 108}},
			{Min: image.Point{X: 26, Y: 130}, Max: image.Point{X: 30, Y: 134}},
			{Min: image.Point{X: 26, Y: 156}, Max: image.Point{X: 30, Y: 160}},
			{Min: image.Point{X: 52, Y: 4}, Max: image.Point{X: 56, Y: 8}},
			{Min: image.Point{X: 52, Y: 26}, Max: image.Point{X: 56, Y: 30}},
			{Min: image.Point{X: 52, Y: 52}, Max: image.Point{X: 56, Y: 56}},
			{Min: image.Point{X: 52, Y: 78}, Max: image.Point{X: 56, Y: 82}},
			{Min: image.Point{X: 52, Y: 104}, Max: image.Point{X: 56, Y: 108}},
			{Min: image.Point{X: 52, Y: 130}, Max: image.Point{X: 56, Y: 134}},
			{Min: image.Point{X: 52, Y: 156}, Max: image.Point{X: 56, Y: 160}},
			{Min: image.Point{X: 78, Y: 4}, Max: image.Point{X: 82, Y: 8}},
			{Min: image.Point{X: 78, Y: 26}, Max: image.Point{X: 82, Y: 30}},
			{Min: image.Point{X: 78, Y: 52}, Max: image.Point{X: 82, Y: 56}},
			{Min: image.Point{X: 78, Y: 78}, Max: image.Point{X: 82, Y: 82}},
			{Min: image.Point{X: 78, Y: 104}, Max: image.Point{X: 82, Y: 108}},
			{Min: image.Point{X: 78, Y: 130}, Max: image.Point{X: 82, Y: 134}},
			{Min: image.Point{X: 78, Y: 156}, Max: image.Point{X: 82, Y: 160}},
			{Min: image.Point{X: 104, Y: 4}, Max: image.Point{X: 108, Y: 8}},
			{Min: image.Point{X: 104, Y: 26}, Max: image.Point{X: 108, Y: 30}},
			{Min: image.Point{X: 104, Y: 52}, Max: image.Point{X: 108, Y: 56}},
			{Min: image.Point{X: 104, Y: 78}, Max: image.Point{X: 108, Y: 82}},
			{Min: image.Point{X: 104, Y: 104}, Max: image.Point{X: 108, Y: 108}},
			{Min: image.Point{X: 104, Y: 130}, Max: image.Point{X: 108, Y: 134}},
			{Min: image.Point{X: 104, Y: 156}, Max: image.Point{X: 108, Y: 160}},
			{Min: image.Point{X: 130, Y: 4}, Max: image.Point{X: 134, Y: 8}},
			{Min: image.Point{X: 130, Y: 26}, Max: image.Point{X: 134, Y: 30}},
			{Min: image.Point{X: 130, Y: 52}, Max: image.Point{X: 134, Y: 56}},
			{Min: image.Point{X: 130, Y: 78}, Max: image.Point{X: 134, Y: 82}},
			{Min: image.Point{X: 130, Y: 104}, Max: image.Point{X: 134, Y: 108}},
			{Min: image.Point{X: 130, Y: 130}, Max: image.Point{X: 134, Y: 134}},
			{Min: image.Point{X: 130, Y: 156}, Max: image.Point{X: 134, Y: 160}},
			{Min: image.Point{X: 156, Y: 26}, Max: image.Point{X: 160, Y: 30}},
			{Min: image.Point{X: 156, Y: 52}, Max: image.Point{X: 160, Y: 56}},
			{Min: image.Point{X: 156, Y: 78}, Max: image.Point{X: 160, Y: 82}},
			{Min: image.Point{X: 156, Y: 104}, Max: image.Point{X: 160, Y: 108}},
			{Min: image.Point{X: 156, Y: 130}, Max: image.Point{X: 160, Y: 134}},
			{Min: image.Point{X: 156, Y: 156}, Max: image.Point{X: 160, Y: 160}},
		},
		{
			{Min: image.Point{X: 4, Y: 30}, Max: image.Point{X: 8, Y: 34}},
			{Min: image.Point{X: 4, Y: 56}, Max: image.Point{X: 8, Y: 60}},
			{Min: image.Point{X: 4, Y: 82}, Max: image.Point{X: 8, Y: 86}},
			{Min: image.Point{X: 4, Y: 108}, Max: image.Point{X: 8, Y: 112}},
			{Min: image.Point{X: 4, Y: 134}, Max: image.Point{X: 8, Y: 138}},
			{Min: image.Point{X: 30, Y: 4}, Max: image.Point{X: 34, Y: 8}},
			{Min: image.Point{X: 30, Y: 30}, Max: image.Point{X: 34, Y: 34}},
			{Min: image.Point{X: 30, Y: 56}, Max: image.Point{X: 34, Y: 60}},
			{Min: image.Point{X: 30, Y: 82}, Max: image.Point{X: 34, Y: 86}},
			{Min: image.Point{X: 30, Y: 108}, Max: image.Point{X: 34, Y: 112}},
			{Min: image.Point{X: 30, Y: 134}, Max: image.Point{X: 34, Y: 138}},
			{Min: image.Point{X: 30, Y: 160}, Max: image.Point{X: 34, Y: 164}},
			{Min: image.Point{X: 56, Y: 4}, Max: image.Point{X: 60, Y: 8}},
			{Min: image.Point{X: 56, Y: 30}, Max: image.Point{X: 60, Y: 34}},
			{Min: image.Point{X: 56, Y: 56}, Max: image.Point{X: 60, Y: 60}},
			{Min: image.Point{X: 56, Y: 82}, Max: image.Point{X: 60, Y: 86}},
			{Min: image.Point{X: 56, Y: 108}, Max: image.Point{X: 60, Y: 112}},
			{Min: image.Point{X: 56, Y: 134}, Max: image.Point{X: 60, Y: 138}},
			{Min: image.Point{X: 56, Y: 160}, Max: image.Point{X: 60, Y: 164}},
			{Min: image.Point{X: 82, Y: 4}, Max: image.Point{X: 86, Y: 8}},
			{Min: image.Point{X: 82, Y: 30}, Max: image.Point{X: 86, Y: 34}},
			{Min: image.Point{X: 82, Y: 56}, Max: image.Point{X: 86, Y: 60}},
			{Min: image.Point{X: 82, Y: 82}, Max: image.Point{X: 86, Y: 86}},
			{Min: image.Point{X: 82, Y: 108}, Max: image.Point{X: 86, Y: 112}},
			{Min: image.Point{X: 82, Y: 134}, Max: image.Point{X: 86, Y: 138}},
			{Min: image.Point{X: 82, Y: 160}, Max: image.Point{X: 86, Y: 164}},
			{Min: image.Point{X: 108, Y: 4}, Max: image.Point{X: 112, Y: 8}},
			{Min: image.Point{X: 108, Y: 30}, Max: image.Point{X: 112, Y: 34}},
			{Min: image.Point{X: 108, Y: 56}, Max: image.Point{X: 112, Y: 60}},
			{Min: image.Point{X: 108, Y: 82}, Max: image.Point{X: 112, Y: 86}},
			{Min: image.Point{X: 108, Y: 108}, Max: image.Point{X: 112, Y: 112}},
			{Min: image.Point{X: 108, Y: 134}, Max: image.Point{X: 112, Y: 138}},
			{Min: image.Point{X: 108, Y: 160}, Max: image.Point{X: 112, Y: 164}},
			{Min: image.Point{X: 134, Y: 4}, Max: image.Point{X: 138, Y: 8}},
			{Min: image.Point{X: 134, Y: 30}, Max: image.Point{X: 138, Y: 34}},
			{Min: image.Point{X: 134, Y: 56}, Max: image.Point{X: 138, Y: 60}},
			{Min: image.Point{X: 134, Y: 82}, Max: image.Point{X: 138, Y: 86}},
			{Min: image.Point{X: 134, Y: 108}, Max: image.Point{X: 138, Y: 112}},
			{Min: image.Point{X: 134, Y: 134}, Max: image.Point{X: 138, Y: 138}},
			{Min: image.Point{X: 134, Y: 160}, Max: image.Point{X: 138, Y: 164}},
			{Min: image.Point{X: 160, Y: 30}, Max: image.Point{X: 164, Y: 34}},
			{Min: image.Point{X: 160, Y: 56}, Max: image.Point{X: 164, Y: 60}},
			{Min: image.Point{X: 160, Y: 82}, Max: image.Point{X: 164, Y: 86}},
			{Min: image.Point{X: 160, Y: 108}, Max: image.Point{X: 164, Y: 112}},
			{Min: image.Point{X: 160, Y: 134}, Max: image.Point{X: 164, Y: 138}},
			{Min: image.Point{X: 160, Y: 160}, Max: image.Point{X: 164, Y: 164}},
		},
		{
			{Min: image.Point{X: 4, Y: 24}, Max: image.Point{X: 8, Y: 28}},
			{Min: image.Point{X: 4, Y: 52}, Max: image.Point{X: 8, Y: 56}},
			{Min: image.Point{X: 4, Y: 80}, Max: image.Point{X: 8, Y: 84}},
			{Min: image.Point{X: 4, Y: 108}, Max: image.Point{X: 8, Y: 112}},
			{Min: image.Point{X: 4, Y: 136}, Max: image.Point{X: 8, Y: 140}},
			{Min: image.Point{X: 24, Y: 4}, Max: image.Point{X: 28, Y: 8}},
			{Min: image.Point{X: 24, Y: 24}, Max: image.Point{X: 28, Y: 28}},
			{Min: image.Point{X: 24, Y: 52}, Max: image.Point{X: 28, Y: 56}},
			{Min: image.Point{X: 24, Y: 80}, Max: image.Point{X: 28, Y: 84}},
			{Min: image.Point{X: 24, Y: 108}, Max: image.Point{X: 28, Y: 112}},
			{Min: image.Point{X: 24, Y: 136}, Max: image.Point{X: 28, Y: 140}},
			{Min: image.Point{X: 24, Y: 164}, Max: image.Point{X: 28, Y: 168}},
			{Min: image.Point{X: 52, Y: 4}, Max: image.Point{X: 56, Y: 8}},
			{Min: image.Point{X: 52, Y: 24}, Max: image.Point{X: 56, Y: 28}},
			{Min: image.Point{X: 52, Y: 52}, Max: image.Point{X: 56, Y: 56}},
			{Min: image.Point{X: 52, Y: 80}, Max: image.Point{X: 56, Y: 84}},
			{Min: image.Point{X: 52, Y: 108}, Max: image.Point{X: 56, Y: 112}},
			{Min: image.Point{X: 52, Y: 136}, Max: image.Point{X: 56, Y: 140}},
			{Min: image.Point{X: 52, Y: 164}, Max: image.Point{X: 56, Y: 168}},
			{Min: image.Point{X: 80, Y: 4}, Max: image.Point{X: 84, Y: 8}},
			{Min: image.Point{X: 80, Y: 24}, Max: image.Point{X: 84, Y: 28}},
			{Min: image.Point{X: 80, Y: 52}, Max: image.Point{X: 84, Y: 56}},
			{Min: image.Point{X: 80, Y: 80}, Max: image.Point{X: 84, Y: 84}},
			{Min: image.Point{X: 80, Y: 108}, Max: image.Point{X: 84, Y: 112}},
			{Min: image.Point{X: 80, Y: 136}, Max: image.Point{X: 84, Y: 140}},
			{Min: image.Point{X: 80, Y: 164}, Max: image.Point{X: 84, Y: 168}},
			{Min: image.Point{X: 108, Y: 4}, Max: image.Point{X: 112, Y: 8}},
			{Min: image.Point{X: 108, Y: 24}, Max: image.Point{X: 112, Y: 28}},
			{Min: image.Point{X: 108, Y: 52}, Max: image.Point{X: 112, Y: 56}},
			{Min: image.Point{X: 108, Y: 80}, Max: image.Point{X: 112, Y: 84}},
			{Min: image.Point{X: 108, Y: 108}, Max: image.Point{X: 112, Y: 112}},
			{Min: image.Point{X: 108, Y: 136}, Max: image.Point{X: 112, Y: 140}},
			{Min: image.Point{X: 108, Y: 164}, Max: image.Point{X: 112, Y: 168}},
			{Min: image.Point{X: 136, Y: 4}, Max: image.Point{X: 140, Y: 8}},
			{Min: image.Point{X: 136, Y: 24}, Max: image.Point{X: 140, Y: 28}},
			{Min: image.Point{X: 136, Y: 52}, Max: image.Point{X: 140, Y: 56}},
			{Min: image.Point{X: 136, Y: 80}, Max: image.Point{X: 140, Y: 84}},
			{Min: image.Point{X: 136, Y: 108}, Max: image.Point{X: 140, Y: 112}},
			{Min: image.Point{X: 136, Y: 136}, Max: image.Point{X: 140, Y: 140}},
			{Min: image.Point{X: 136, Y: 164}, Max: image.Point{X: 140, Y: 168}},
			{Min: image.Point{X: 164, Y: 24}, Max: image.Point{X: 168, Y: 28}},
			{Min: image.Point{X: 164, Y: 52}, Max: image.Point{X: 168, Y: 56}},
			{Min: image.Point{X: 164, Y: 80}, Max: image.Point{X: 168, Y: 84}},
			{Min: image.Point{X: 164, Y: 108}, Max: image.Point{X: 168, Y: 112}},
			{Min: image.Point{X: 164, Y: 136}, Max: image.Point{X: 168, Y: 140}},
			{Min: image.Point{X: 164, Y: 164}, Max: image.Point{X: 168, Y: 168}},
		},
		{
			{Min: image.Point{X: 4, Y: 28}, Max: image.Point{X: 8, Y: 32}},
			{Min: image.Point{X: 4, Y: 56}, Max: image.Point{X: 8, Y: 60}},
			{Min: image.Point{X: 4, Y: 84}, Max: image.Point{X: 8, Y: 88}},
			{Min: image.Point{X: 4, Y: 112}, Max: image.Point{X: 8, Y: 116}},
			{Min: image.Point{X: 4, Y: 140}, Max: image.Point{X: 8, Y: 144}},
			{Min: image.Point{X: 28, Y: 4}, Max: image.Point{X: 32, Y: 8}},
			{Min: image.Point{X: 28, Y: 28}, Max: image.Point{X: 32, Y: 32}},
			{Min: image.Point{X: 28, Y: 56}, Max: image.Point{X: 32, Y: 60}},
			{Min: image.Point{X: 28, Y: 84}, Max: image.Point{X: 32, Y: 88}},
			{Min: image.Point{X: 28, Y: 112}, Max: image.Point{X: 32, Y: 116}},
			{Min: image.Point{X: 28, Y: 140}, Max: image.Point{X: 32, Y: 144}},
			{Min: image.Point{X: 28, Y: 168}, Max: image.Point{X: 32, Y: 172}},
			{Min: image.Point{X: 56, Y: 4}, Max: image.Point{X: 60, Y: 8}},
			{Min: image.Point{X: 56, Y: 28}, Max: image.Point{X: 60, Y: 32}},
			{Min: image.Point{X: 56, Y: 56}, Max: image.Point{X: 60, Y: 60}},
			{Min: image.Point{X: 56, Y: 84}, Max: image.Point{X: 60, Y: 88}},
			{Min: image.Point{X: 56, Y: 112}, Max: image.Point{X: 60, Y: 116}},
			{Min: image.Point{X: 56, Y: 140}, Max: image.Point{X: 60, Y: 144}},
			{Min: image.Point{X: 56, Y: 168}, Max: image.Point{X: 60, Y: 172}},
			{Min: image.Point{X: 84, Y: 4}, Max: image.Point{X: 88, Y: 8}},
			{Min: image.Point{X: 84, Y: 28}, Max: image.Point{X: 88, Y: 32}},
			{Min: image.Point{X: 84, Y: 56}, Max: image.Point{X: 88, Y: 60}},
			{Min: image.Point{X: 84, Y: 84}, Max: image.Point{X: 88, Y: 88}},
			{Min: image.Point{X: 84, Y: 112}, Max: image.Point{X: 88, Y: 116}},
			{Min: image.Point{X: 84, Y: 140}, Max: image.Point{X: 88, Y: 144}},
			{Min: image.Point{X: 84, Y: 168}, Max: image.Point{X: 88, Y: 172}},
			{Min: image.Point{X: 112, Y: 4}, Max: image.Point{X: 116, Y: 8}},
			{Min: image.Point{X: 112, Y: 28}, Max: image.Point{X: 116, Y: 32}},
			{Min: image.Point{X: 112, Y: 56}, Max: image.Point{X: 116, Y: 60}},
			{Min: image.Point{X: 112, Y: 84}, Max: image.Point{X: 116, Y: 88}},
			{Min: image.Point{X: 112, Y: 112}, Max: image.Point{X: 116, Y: 116}},
			{Min: image.Point{X: 112, Y: 140}, Max: image.Point{X: 116, Y: 144}},
			{Min: image.Point{X: 112, Y: 168}, Max: image.Point{X: 116, Y: 172}},
			{Min: image.Point{X: 140, Y: 4}, Max: image.Point{X: 144, Y: 8}},
			{Min: image.Point{X: 140, Y: 28}, Max: image.Point{X: 144, Y: 32}},
			{Min: image.Point{X: 140, Y: 56}, Max: image.Point{X: 144, Y: 60}},
			{Min: image.Point{X: 140, Y: 84}, Max: image.Point{X: 144, Y: 88}},
			{Min: image.Point{X: 140, Y: 112}, Max: image.Point{X: 144, Y: 116}},
			{Min: image.Point{X: 140, Y: 140}, Max: image.Point{X: 144, Y: 144}},
			{Min: image.Point{X: 140, Y: 168}, Max: image.Point{X: 144, Y: 172}},
			{Min: image.Point{X: 168, Y: 28}, Max: image.Point{X: 172, Y: 32}},
			{Min: image.Point{X: 168, Y: 56}, Max: image.Point{X: 172, Y: 60}},
			{Min: image.Point{X: 168, Y: 84}, Max: image.Point{X: 172, Y: 88}},
			{Min: image.Point{X: 168, Y: 112}, Max: image.Point{X: 172, Y: 116}},
			{Min: image.Point{X: 168, Y: 140}, Max: image.Point{X: 172, Y: 144}},
			{Min: image.Point{X: 168, Y: 168}, Max: image.Point{X: 172, Y: 172}},
		},
	} // 版本的对齐图形
	markFunc = [maxMark]func(x, y int) bool{
		func(x, y int) bool {
			return (x+y)%2 == 0
		},
		func(x, y int) bool {
			return y%2 == 0
		},
		func(x, y int) bool {
			return x%3 == 0
		},
		func(x, y int) bool {
			return (x+y)%3 == 0
		},
		func(x, y int) bool {
			return (int(float64(y)/2)+int(float64(x)/3))%2 == 0
		},
		func(x, y int) bool {
			return ((x*y)%2 + (x*y)%3) == 0
		},
		func(x, y int) bool {
			return ((x*y)%2+(x*y)%3)%2 == 0
		},
		func(x, y int) bool {
			return ((x+y)%2+(x*y)%3)%2 == 0
		},
	}
	evaluation3Bytes = [2][]uint8{
		{
			_paletteBlack, _paletteWhite, _paletteBlack, _paletteBlack, _paletteBlack,
			_paletteWhite, _paletteBlack, _paletteWhite, _paletteWhite, _paletteWhite, _paletteWhite,
		},
		{
			_paletteWhite, _paletteWhite, _paletteWhite, _paletteWhite, _paletteBlack,
			_paletteWhite, _paletteBlack, _paletteBlack, _paletteBlack, _paletteWhite, _paletteBlack,
		},
	} // 0水平，1垂直
)

func init() {
	_pool.New = func() interface{} {
		q := new(qrCode)
		q.buffer.data = make([]byte, 1)
		q.strEnc.bitD = make([]byte, 1)
		q.strEnc.buff = &q.buffer
		q.eccEnc.poly = make([]byte, 1)
		q.eccEnc.buff = &q.buffer
		return q
	}
}

func PNG(w io.Writer, str string, level Level, compress png.CompressionLevel) error {
	img, err := Image(str, level)
	if err != nil {
		return err
	}
	enc := png.Encoder{
		CompressionLevel: compress,
		BufferPool:       nil,
	}
	return enc.Encode(w, img)
}

func JPEG(w io.Writer, str string, level Level, quality int) error {
	img, err := Image(str, level)
	if err != nil {
		return err
	}
	return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
}

func Image(str string, level Level) (image.Image, error) {
	q := _pool.Get().(*qrCode)
	// 字符串编码
	err := q.strEnc.Encode(str, level)
	if err != nil {
		_pool.Put(q)
		return nil, err
	}
	// 纠错编码
	q.eccEnc.Encode(q.strEnc.bitD, q.strEnc.version, level)
	// 位图
	img := new(image.Paletted)
	img.Stride = qrCodeSizeTable[q.strEnc.version] + 8
	img.Rect.Max.X = img.Stride
	img.Rect.Max.Y = img.Stride
	img.Palette = _palette
	img.Pix = make([]uint8, img.Stride*img.Stride)
	q.Draw(img)
	// 回收缓存
	_pool.Put(q)
	// 返回
	return img, err
}

// 缓存
type buffer struct {
	data []byte
}

func (b *buffer) Reset() {
	b.data = b.data[:0]
}

func (b *buffer) Resize(n, o int) {
	if cap(b.data) < n {
		b.data = make([]byte, n)
	} else {
		b.data = b.data[:n]
		if o >= 0 {
			for i := o; i < n; i++ {
				b.data[i] = 0
			}
		}
	}
}

type qrCode struct {
	buffer     buffer     // 共享缓存
	strEnc     strEncoder // 字符串编码
	eccEnc     eccEncoder // 纠错编码
	pixXY      [][]uint8  // 位图二维数组指针
	markNum    int        // 使用的mark图编号
	markData   buffer     // mark后的最终数据
	markBuffXY [][]uint8  // mark后的缓存数组的二维指针
	markDataXY [][]uint8  // mark后的缓存数组的二维指针
}

// 画图
func (q *qrCode) Draw(img *image.Paletted) {
	// 图像数据，包括两边的4个空白
	q.buffer.Resize((img.Stride-8)*(img.Stride-8), -1)
	q.markData.Resize(len(q.buffer.data), -1)
	// 二维表，便于操作
	pix1 := img.Pix[4*(img.Stride)+4:]
	pix2 := q.buffer.data
	pix3 := q.markData.data
	q.pixXY = q.pixXY[:0]
	q.markBuffXY = q.markBuffXY[:0]
	q.markDataXY = q.markDataXY[:0]
	for i := 0; i < qrCodeSizeTable[q.strEnc.version]; i++ {
		q.pixXY = append(q.pixXY, pix1[:qrCodeSizeTable[q.strEnc.version]])
		pix1 = pix1[qrCodeSizeTable[q.strEnc.version]+8:]
		q.markBuffXY = append(q.markBuffXY, pix2[:qrCodeSizeTable[q.strEnc.version]])
		pix2 = pix2[qrCodeSizeTable[q.strEnc.version]:]
		q.markDataXY = append(q.markDataXY, pix3[:qrCodeSizeTable[q.strEnc.version]])
		pix3 = pix3[qrCodeSizeTable[q.strEnc.version]:]
	}
	// 开始画图
	q.drawFinderPatterns()
	q.drawTimingPatterns()
	q.drawAlignmentPatterns()
	q.drawBottomLeftPoint()
	q.drawData()
	q.mark()
	q.drawFormatInformation()
	q.drawVersionInformation()
}

// 画点
func (q *qrCode) drawPoint(x, y int, c uint8) {
	q.pixXY[y][x] = c
}

// 画矩形
func (q *qrCode) drawRectangle(x1, y1, x2, y2 int, c uint8) {
	// 上下
	for x := x1; x <= x2; x++ {
		q.drawPoint(x, y1, c)
		q.drawPoint(x, y2, c)
	}
	// 左右
	for y := y1 + 1; y < y2; y++ {
		q.drawPoint(x1, y, c)
		q.drawPoint(x2, y, c)
	}
}

// 画矩形
func (q *qrCode) drawSolidRectangle(x1, y1, x2, y2 int, c uint8) {
	for y := y1; y <= y2; y++ {
		for x := x1; x <= x2; x++ {
			q.drawPoint(x, y, c)
		}
	}
}

// finder patterns
func (q *qrCode) drawFinderPatterns() {
	// 左上角
	q.drawRectangle(0, 0, 6, 6, _paletteBlack)
	q.drawSolidRectangle(2, 2, 4, 4, _paletteBlack)
	// 右上角
	q.drawRectangle(qrCodeSizeTable[q.strEnc.version]-7, 0,
		qrCodeSizeTable[q.strEnc.version]-1, 6, _paletteBlack)
	q.drawSolidRectangle(qrCodeSizeTable[q.strEnc.version]-5, 2,
		qrCodeSizeTable[q.strEnc.version]-3, 4, _paletteBlack)
	// 左下角
	q.drawRectangle(0, qrCodeSizeTable[q.strEnc.version]-7, 6,
		qrCodeSizeTable[q.strEnc.version]-1, _paletteBlack)
	q.drawSolidRectangle(2, qrCodeSizeTable[q.strEnc.version]-5, 4,
		qrCodeSizeTable[q.strEnc.version]-3, _paletteBlack)
}

// timing patterns
func (q *qrCode) drawTimingPatterns() {
	// 水平
	for i := 8; i < qrCodeSizeTable[q.strEnc.version]-8; {
		q.drawPoint(i, 6, _paletteBlack)
		i += 2
	}
	// 垂直
	for i := 8; i < qrCodeSizeTable[q.strEnc.version]-8; {
		q.drawPoint(6, i, _paletteBlack)
		i += 2
	}
}

// alignment patterns，
func (q *qrCode) drawAlignmentPatterns() {
	for _, r := range alignmentPatternTable[q.strEnc.version] {
		q.drawRectangle(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, _paletteBlack)
		q.drawPoint(r.Min.X+2, r.Min.Y+2, _paletteBlack)
	}
}

// 左下角，格式信息上的一个黑点
func (q *qrCode) drawBottomLeftPoint() {
	// y=version*4+4+9，
	q.drawPoint(8, int(q.strEnc.version)*4+13, _paletteBlack)
}

// 格式信息
func (q *qrCode) drawFormatInformation() {
	f := formatBitTable[q.strEnc.Level][q.markNum]
	idx := 0
	// 左上角
	for x := 0; x < 6; x++ {
		if f[idx] == 1 {
			q.drawPoint(x, 8, _paletteBlack)
		}
		idx++
	}
	if f[idx] == 1 {
		q.drawPoint(7, 8, _paletteBlack)
	}
	idx++
	if f[idx] == 1 {
		q.drawPoint(8, 8, _paletteBlack)
	}
	idx++
	if f[idx] == 1 {
		q.drawPoint(8, 7, _paletteBlack)
	}
	idx++
	for y := 5; y >= 0; y-- {
		if f[idx] == 1 {
			q.drawPoint(8, y, _paletteBlack)
		}
		idx++
	}
	idx = 0
	// 左下角
	for y := qrCodeSizeTable[q.strEnc.version] - 1; y > qrCodeSizeTable[q.strEnc.version]-8; y-- {
		if f[idx] == 1 {
			q.drawPoint(8, y, _paletteBlack)
		}
		idx++
	}
	// 右上角
	for x := qrCodeSizeTable[q.strEnc.version] - 8; x <= qrCodeSizeTable[q.strEnc.version]-1; x++ {
		if f[idx] == 1 {
			q.drawPoint(x, 8, _paletteBlack)
		}
		idx++
	}
}

// 版本信息
func (q *qrCode) drawVersionInformation() {
	if q.strEnc.version < version7 {
		return
	}
	ver := versionBitTable[q.strEnc.version]
	x := qrCodeSizeTable[q.strEnc.version] - 11
	y := qrCodeSizeTable[q.strEnc.version] - 8
	idx := 0
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			if ver[idx] == 1 {
				// 左下角
				q.drawPoint(i, j+y, _paletteBlack)
				// 右上角
				q.drawPoint(i+x, j, _paletteBlack)
			}
			idx++
		}
	}
}

// 数据
func (q *qrCode) drawData() {
	// 从右下角开始
	x := qrCodeSizeTable[q.strEnc.version] - 1
	y := qrCodeSizeTable[q.strEnc.version] - 1
	// finder patterns，包括format区域
	topCornerMaxY := 9
	leftCornerMaxX := 9
	leftBottomCornerMinY := y - 8
	rightTopCornerMinX := x - 8
	// version information area
	topVersionMinX := rightTopCornerMinX - 3
	bottomVersionMinY := leftBottomCornerMinY - 3
	// 画
	idx := 0
	bit := byte(0b10000000)
	char := q.eccEnc.data[idx]
	up := true
	drawPoint := func() bool {
		if char&bit != 0 {
			q.drawPoint(x, y, _paletteBlack)
		}
		bit >>= 1
		if bit == 0 {
			bit = 0b10000000
			idx++
			if idx == len(q.eccEnc.data) {
				return false
			}
			char = q.eccEnc.data[idx]
		}
		return true
	}
Loop:
	for {
		if up {
			// 右点
			if !drawPoint() {
				return
			}
			// 左点
			x--
			if !drawPoint() {
				return
			}
			// 顶部
			if y == 0 {
				x--
				up = !up
				continue Loop
			}
			// finder patterns
			if y == topCornerMaxY {
				// 右上
				if x > rightTopCornerMinX {
					x--
					up = !up
					continue Loop
				}
				// 左上
				if x < leftCornerMaxX {
					x--
					// timing patterns，垂直
					if x == timingPattern {
						x--
					}
					up = !up
					continue Loop
				}
				// 左下，不可能
			}
			// 上移
			y--
			// 检查alignment patterns
			for _, r := range alignmentPatternTable[q.strEnc.version] {
				if y == r.Max.Y {
					// 右边向上
					if x == r.Max.X {
						x++
						for i := 0; i < 5; i++ {
							if !drawPoint() {
								break Loop
							}
							y--
							// timing patterns，水平
							if y == timingPattern {
								i++
								y--
							}
						}
						continue Loop
					}
					// 跳过矩形
					if x >= r.Min.X && x < r.Max.X {
						x++
						y = r.Min.Y - 1
						continue Loop
					}
					// 左边向上
					if x == r.Min.X-1 {
						for i := 0; i < 5; i++ {
							if !drawPoint() {
								break Loop
							}
							y--
							// timing patterns，水平
							if y == timingPattern {
								i++
								y--
							}
						}
						x++
						continue Loop
					}
				}
			}
			// timing patterns，水平
			if y == timingPattern {
				// 版本7以上
				if q.strEnc.version >= version7 && x > topVersionMinX {
					// 右上版本区左边向下
					x -= 2
					y = 0
					for i := 0; i < 6; i++ {
						if !drawPoint() {
							break Loop
						}
						y++
					}
					x++
					y++
					up = !up
					continue Loop
				}
				// 上移
				y--
			}
			x++
		} else {
			// 右点
			if !drawPoint() {
				break Loop
			}
			// 左点
			x--
			if !drawPoint() {
				break Loop
			}
			// 最左边
			if x < timingPattern {
				// 左下
				if q.strEnc.version >= version7 {
					if y == bottomVersionMinY {
						x--
						up = !up
						continue Loop
					}
				} else {
					if y == leftBottomCornerMinY {
						x--
						up = !up
						continue Loop
					}
				}
			}
			// 下移
			y++
			// timing patterns，水平
			if y == timingPattern {
				x++
				y++
				continue Loop
			}
			// 检查alignment patterns
			for _, r := range alignmentPatternTable[q.strEnc.version] {
				if y == r.Min.Y {
					// 右边向下
					if x == r.Max.X {
						x++
						for i := 0; i < 5; i++ {
							if !drawPoint() {
								break Loop
							}
							y++
							// timing patterns，水平
							if y == timingPattern {
								i++
								y++
							}
						}
						continue Loop
					}
					// 跳过矩形
					if x >= r.Min.X && x < r.Max.X {
						x++
						y = r.Max.Y + 1
						continue Loop
					}
					// 左边向下
					if x == r.Min.X-1 {
						for i := 0; i < 5; i++ {
							if !drawPoint() {
								break Loop
							}
							y++
							// timing patterns，水平
							if y == timingPattern {
								i++
								y++
							}
						}
						x++
						continue Loop
					}
				}
			}
			// 底部
			if y == qrCodeSizeTable[q.strEnc.version] {
				// 左移
				x--
				if x < leftCornerMaxX {
					// 左下角
					y = leftBottomCornerMinY
				} else {
					y--
				}
				// 向下
				up = !up
				continue Loop
			}
			x++
		}
	}
}

// 对原始位图数据pix分别进行8种mark，最小评分的mark将作为最终的输出数据。
func (q *qrCode) mark() {
	// 以下是不能mark的区域
	finderPatterns := [3]image.Rectangle{
		{
			Max: image.Point{X: 8, Y: 8},
		}, // 左上
		{
			Min: image.Point{X: qrCodeSizeTable[q.strEnc.version] - 8},
			Max: image.Point{X: qrCodeSizeTable[q.strEnc.version] - 1, Y: 8},
		}, // 右上
		{
			Min: image.Point{Y: qrCodeSizeTable[q.strEnc.version] - 8},
			Max: image.Point{X: 8, Y: qrCodeSizeTable[q.strEnc.version] - 1},
		}, // 左下
	}
	alignmentPatterns := alignmentPatternTable[q.strEnc.version]
	versionArea := [2]image.Rectangle{
		{
			Min: image.Point{X: finderPatterns[1].Min.X - 3},
			Max: image.Point{X: finderPatterns[1].Min.X - 1, Y: 5},
		}, // 右上
		{
			Min: image.Point{Y: finderPatterns[2].Min.Y - 3},
			Max: image.Point{X: 5, Y: finderPatterns[1].Min.X - 1},
		}, // 左下
	}
	// 得分
	score, minScore := 0, 0xffffffff
	// 生成mark图
	for i := 0; i < maxMark; i++ {
		for y := 0; y < len(q.pixXY); y++ {
		Next:
			for x := 0; x < len(q.pixXY[y]); x++ {
				// finder patterns
				for j := 0; j < len(finderPatterns); j++ {
					if x >= finderPatterns[j].Min.X && x <= finderPatterns[j].Max.X &&
						y >= finderPatterns[j].Min.Y && y <= finderPatterns[j].Max.Y {
						q.markBuffXY[y][x] = q.pixXY[y][x]
						continue Next
					}
				}
				// timing patterns
				if x == timingPattern || y == timingPattern {
					q.markBuffXY[y][x] = q.pixXY[y][x]
					continue Next
				}
				// alignment patterns
				for j := 0; j < len(alignmentPatterns); j++ {
					if x >= alignmentPatterns[j].Min.X && x <= alignmentPatterns[j].Max.X &&
						y >= alignmentPatterns[j].Min.Y && y <= alignmentPatterns[j].Max.Y {
						q.markBuffXY[y][x] = q.pixXY[y][x]
						continue Next
					}
				}
				// version information
				if q.strEnc.version >= version7 {
					for j := 0; j < len(versionArea); j++ {
						if x >= versionArea[j].Min.X && x <= versionArea[j].Max.X &&
							y >= versionArea[j].Min.Y && y <= versionArea[j].Max.Y {
							q.markBuffXY[y][x] = q.pixXY[y][x]
							continue Next
						}
					}
				}
				if markFunc[i](x, y) {
					q.markBuffXY[y][x] = _paletteBlack ^ q.pixXY[y][x]
					//q.markBuffXY[y][x] = _paletteBlack
				} else {
					q.markBuffXY[y][x] = _paletteWhite ^ q.pixXY[y][x]
					//q.markBuffXY[y][x] = _paletteWhite
				}
			}
		}
		// 评估
		score = q.evaluation1() + q.evaluation2() + q.evaluation3() + q.evaluation4()
		// 最小得分
		if score < minScore {
			minScore = score
			q.markNum = i
			t1 := q.markBuffXY
			q.markBuffXY = q.markDataXY
			q.markDataXY = t1
			t2 := q.buffer.data
			q.buffer.data = q.markData.data
			q.markData.data = t2
		}
	}
	// 最终的数据
	for y := 0; y < len(q.markDataXY); y++ {
		copy(q.pixXY[y], q.markDataXY[y])
	}
}

// 找到5个连续颜色的点，+3分
// 5个连续颜色的点之后，每多1个点+1分
func (q *qrCode) evaluation1() int {
	score := 0
	var lastBlock uint8
	var consecutive int
	x, y := 0, 0
	// 行
	for ; y < len(q.markBuffXY); y++ {
		consecutive = 0
		lastBlock = q.markBuffXY[y][0]
		for x = 1; x < len(q.markBuffXY[y]); x++ {
			if q.markBuffXY[y][x] == lastBlock {
				consecutive++
				if consecutive == 5 {
					score += 3
				} else if consecutive > 5 {
					score++
				}
			} else {
				lastBlock = q.markBuffXY[y][x]
				consecutive = 0
			}
		}
	}
	// 列
	x = 0
	for ; x < len(q.markBuffXY[0]); x++ {
		consecutive = 0
		lastBlock = q.markBuffXY[0][x]
		for y = 1; y < len(q.markBuffXY); y++ {
			if q.markBuffXY[y][x] == lastBlock {
				consecutive++
				if consecutive == 5 {
					score += 3
				} else if consecutive > 5 {
					score++
				}
			} else {
				lastBlock = q.markBuffXY[y][x]
				consecutive = 0
			}
		}
	}
	return score
}

// 找到相同颜色的最小矩形（2*2），+3分
func (q *qrCode) evaluation2() int {
	score := 0
	x, y := 0, 0
	for ; y < len(q.markBuffXY)-1; y++ {
		for ; x < len(q.markBuffXY[y])-1; x++ {
			if q.markBuffXY[y][x] == q.markBuffXY[y][x+1] &&
				q.markBuffXY[y][x] == q.markBuffXY[y+1][x] &&
				q.markBuffXY[y][x] == q.markBuffXY[y+1][x+1] {
				score += 3
			}
		}
	}
	return score
}

// 找到[10111010000]或者[00001011101]，+40分
func (q *qrCode) evaluation3() int {
	score := 0
	x, y, n, m := 0, 0, 0, 0
	o := true
	// 行
	m = len(q.markBuffXY[0]) - len(evaluation3Bytes[0]) - 1
	for ; y < len(q.markBuffXY); y++ {
		for x < m {
			o = true
			n = x
			for _, c := range evaluation3Bytes[0] {
				if c != q.markBuffXY[y][n] {
					o = false
					break
				}
				n++
			}
			if o {
				score += 40
				x += len(evaluation3Bytes[0])
			} else {
				x++
			}
		}
		x = 0
	}
	// 列
	m = len(q.markBuffXY) - len(evaluation3Bytes[1]) - 1
	x = 0
	for ; x < len(q.markBuffXY[0]); x++ {
		y = 0
		for y < m {
			o = true
			n = y
			for _, c := range evaluation3Bytes[1] {
				if c != q.markBuffXY[n][x] {
					o = false
					break
				}
				n++
			}
			if o {
				score += 40
				y += len(evaluation3Bytes[1])
			} else {
				y++
			}
		}
	}
	return score
}

// 计算黑点的百分比，（黑点总数/总点）*100，然后得到前后两个5的倍数。
// 这两个数-50，再求绝对值，然后除以5。
// 得到的最小商*10，就是最终得分
func (q *qrCode) evaluation4() int {
	// 黑点的总数
	n := 0
	for i := 0; i < len(q.buffer.data); i++ {
		if q.buffer.data[i] == _paletteBlack {
			n++
		}
	}
	// 百分比
	m := int(math.Ceil(float64(n) / float64(len(q.buffer.data)) * 100))
	// 得到前后两个5的倍数
	n1 := m / 5 * 5
	n2 := (m/5 + 1) * 5
	// 这两个数-50，再求绝对值
	if n1 >= 50 {
		n1 -= 50
	} else {
		n1 = 50 - n1
	}
	if n2 >= 50 {
		n2 -= 50
	} else {
		n2 = 50 - n2
	}
	// 除以5
	n1 /= 5
	n2 /= 5
	// 最小商*10
	if n1 > n2 {
		return n2 * 10
	} else {
		return n1 * 10
	}
}
