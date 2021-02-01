package qrcode

/*
参考文档：https://www.thonky.com/qr-code-tutorial/
*/

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"strings"
	"sync"
)

const (
	maxMark = 8
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
			return int(float64(y)/2+float64(x)/3)%2 == 0
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
		q.strEncoder.bitD = make([]byte, 1)
		q.strEncoder.buff = &q.buffer
		q.eccEncoder.poly = make([]byte, 1)
		q.eccEncoder.buff = &q.buffer
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
	err := q.strEncoder.Encode(str, level)
	if err != nil {
		_pool.Put(q)
		return nil, err
	}
	// 纠错编码
	q.eccEncoder.Encode(q.strEncoder.bitD, q.strEncoder.version, level)
	// 位图
	img := new(image.Paletted)
	img.Stride = qrCodeSizeTable[q.strEncoder.version]
	img.Rect.Max.X = img.Stride
	img.Rect.Max.Y = img.Stride
	img.Palette = _palette
	img.Pix = q.Draw()
	// 回收缓存
	_pool.Put(q)
	// 返回
	return img, err
}

func printBits(b []byte) {
	var str strings.Builder
	for i := 0; i < len(b); i++ {
		for j := 7; j >= 0; j-- {
			str.WriteString(fmt.Sprint((b[i] >> j) & 0b00000001))
		}
	}
	fmt.Println(str.String())
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
	strEncoder strEncoder // 字符串编码
	eccEncoder eccEncoder // 纠错编码
	buffer     buffer     // 共享缓存
	pix        []uint8    // 原始位图数据
	pixXY      [][]uint8  // 位图二维数组指针
	markNum    int        // 使用的mark图编号
	markData   buffer     // mark后的最终数据
	markXY     [][]uint8  // mark后的缓存数组的二维指针
}

// 根据原始图像的一维数组，生成二维表，便于操作
func (q *qrCode) initPixXY() {
	q.pixXY = q.pixXY[:0]
	b := q.pix
	for i := 0; i < qrCodeSizeTable[q.strEncoder.version]; i++ {
		q.pixXY = append(q.pixXY, b[:qrCodeSizeTable[q.strEncoder.version]])
		b = b[qrCodeSizeTable[q.strEncoder.version]:]
	}
}

// 根据mark后的图像的一维数组缓存，生成二维表，便于操作
func (q *qrCode) initMarkXY() {
	q.markXY = q.markXY[:0]
	b := q.buffer.data
	for i := 0; i < qrCodeSizeTable[q.strEncoder.version]; i++ {
		q.markXY = append(q.markXY, b[:qrCodeSizeTable[q.strEncoder.version]])
		b = b[qrCodeSizeTable[q.strEncoder.version]:]
	}
}

// 画图
func (q *qrCode) Draw() []uint8 {
	q.pix = make([]uint8, qrCodeSizeTable[q.strEncoder.version]*qrCodeSizeTable[q.strEncoder.version])
	q.buffer.Resize(len(q.pix), 0)
	q.markData.Resize(len(q.pix), -1)
	q.initPixXY()
	q.drawFinderPatterns()
	q.drawTimingPatterns()
	q.drawAlignmentPatterns()
	q.drawBottomLeftPoint()
	q.drawData()
	q.mark()
	q.drawFormat()
	q.drawVersion()
	return q.pix
}

// 画点
func (q *qrCode) drawPoint(x, y int, c uint8) {
	q.pixXY[y][x] = c
}

// 画矩形
func (q *qrCode) drawRectangle(x1, y1, x2, y2 int, c uint8, fill bool) {
	if fill {
		for i := y1; i <= y2; i++ {
			for j := x1; j <= x2; j++ {
				q.drawPoint(j, i, c)
			}
		}
		return
	}
	// 上下
	for i := x1; i <= x2; i++ {
		q.drawPoint(i, y1, c)
		q.drawPoint(i, y2, c)
	}
	// 左右
	for i := y1 + 1; i < y2; i++ {
		q.drawPoint(x1, i, c)
		q.drawPoint(x2, i, c)
	}
}

// finder patterns
func (q *qrCode) drawFinderPatterns() {
	// 左上角
	q.drawRectangle(0, 0, 6, 6, _paletteBlack, false)
	q.drawRectangle(2, 2, 4, 4, _paletteBlack, true)
	// 右上角
	q.drawRectangle(qrCodeSizeTable[q.strEncoder.version]-7, 0, qrCodeSizeTable[q.strEncoder.version]-1, 6, _paletteBlack, false)
	q.drawRectangle(qrCodeSizeTable[q.strEncoder.version]-5, 2, qrCodeSizeTable[q.strEncoder.version]-3, 4, _paletteBlack, true)
	// 左下角
	q.drawRectangle(0, qrCodeSizeTable[q.strEncoder.version]-7, 6, qrCodeSizeTable[q.strEncoder.version]-1, _paletteBlack, false)
	q.drawRectangle(2, qrCodeSizeTable[q.strEncoder.version]-5, 4, qrCodeSizeTable[q.strEncoder.version]-3, _paletteBlack, true)
}

// timing patterns
func (q *qrCode) drawTimingPatterns() {
	// 水平
	for i := 8; i < qrCodeSizeTable[q.strEncoder.version]-8; {
		q.drawPoint(i, 6, _paletteBlack)
		i += 2
	}
	// 垂直
	for i := 8; i < qrCodeSizeTable[q.strEncoder.version]-8; {
		q.drawPoint(6, i, _paletteBlack)
		i += 2
	}
}

// alignment patterns，
func (q *qrCode) drawAlignmentPatterns() {
	for _, r := range alignmentPatternTable[q.strEncoder.version] {
		q.drawRectangle(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, _paletteBlack, false)
		q.drawPoint(r.Min.X+2, r.Min.Y+2, _paletteBlack)
	}
}

// 左下角，格式信息上的一个黑点
func (q *qrCode) drawBottomLeftPoint() {
	// y=version*4+4+9，
	q.drawPoint(8, int(q.strEncoder.version)*4+13, _paletteBlack)
}

// 格式信息
func (q *qrCode) drawFormat() {
	f := formatBitTable[q.strEncoder.Level][q.markNum]
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
	for y := qrCodeSizeTable[q.strEncoder.version] - 1; y > qrCodeSizeTable[q.strEncoder.version]-8; y-- {
		if f[idx] == 1 {
			q.drawPoint(8, y, _paletteBlack)
		}
		idx++
	}
	// 右上角
	for x := qrCodeSizeTable[q.strEncoder.version] - 8; x <= qrCodeSizeTable[q.strEncoder.version]-1; x++ {
		if f[idx] == 1 {
			q.drawPoint(x, 8, _paletteBlack)
		}
		idx++
	}
}

// 版本信息
func (q *qrCode) drawVersion() {
	if q.strEncoder.version < version7 {
		return
	}
	ver := versionBitTable[q.strEncoder.version]
	x := qrCodeSizeTable[q.strEncoder.version] - 11
	y := qrCodeSizeTable[q.strEncoder.version] - 8
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
	x := qrCodeSizeTable[q.strEncoder.version] - 1
	y := qrCodeSizeTable[q.strEncoder.version] - 1
	// finder patterns，左上0，右上1，左下2
	finderPatterns := [3]image.Point{{9, 9}, {x - 8, 9}, {9, y - 8}}
	// align patterns 矩形
	alignPatterns := alignmentPatternTable[q.strEncoder.version]
	// timing patterns
	timingPatterns := image.Point{X: 6, Y: 6}
	// version patterns，0右上x，1左下y
	versionPatterns := image.Point{X: finderPatterns[1].X - 3, Y: finderPatterns[2].Y - 3}
	idx := 0
	bit := byte(0b10000000)
	char := q.eccEncoder.data[idx]
	up := true
	setColor := func() bool {
		if char&bit != 0 {
			q.drawPoint(x, y, _paletteBlack)
		}
		bit >>= 1
		if bit == 0 {
			bit = 0b10000000
			idx++
			if idx == len(q.eccEncoder.data) {
				return false
			}
			char = q.eccEncoder.data[idx]
		}
		return true
	}
Loop:
	for {
		if up {
			// 右点
			if !setColor() {
				break Loop
			}
			// 左点
			x--
			if !setColor() {
				break Loop
			}
			// finder patterns
			if y == finderPatterns[1].Y {
				// 右上
				if x > finderPatterns[1].X {
					x--
					up = !up
					continue Loop
				}
				// 左上
				if x < finderPatterns[0].X {
					x--
					// timing patterns，垂直
					if x == timingPatterns.X {
						x--
					}
					up = !up
					continue Loop
				}
				// 左下，不可能
			}
			// timing patterns，水平
			if y == timingPatterns.Y {
				// 版本7以上
				if q.strEncoder.version >= 6 && x > versionPatterns.X {
					// 右上版本区左边向下
					x -= 2
					y = 0
					for i := 0; i < 6; i++ {
						if !setColor() {
							break Loop
						}
						y++
					}
					x++
					y++
					up = !up
					continue Loop
				}
			}
			// 上边缘
			if y == 0 {
				x--
				if x < finderPatterns[0].X {
					y = finderPatterns[0].Y
				}
				up = !up
				continue Loop
			}
			// 上移
			y--
			// 检查alignment patterns
			for _, r := range alignPatterns {
				if y == r.Max.Y {
					// 右边向上
					if x == r.Max.X {
						x++
						for i := 0; i < 5; i++ {
							if !setColor() {
								break Loop
							}
							y--
							// timing patterns，水平
							if y == 6 {
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
							if !setColor() {
								break Loop
							}
							y--
							// timing patterns，水平
							if y == 6 {
								i++
								y--
							}
						}
						x++
						continue Loop
					}
				}
			}
			x++
		} else {
			// 右点
			if !setColor() {
				break Loop
			}
			// 左点
			x--
			if !setColor() {
				break Loop
			}
			// finder patterns
			if x < timingPatterns.X {
				// 左下
				if q.strEncoder.version >= 6 {
					if y == versionPatterns.Y {
						x--
						up = !up
						continue Loop
					}
				} else {
					if y == finderPatterns[2].Y {
						x--
						up = !up
						continue Loop
					}
				}
			} else if x < finderPatterns[2].X {
				if y == finderPatterns[2].Y {
					x--
					if x == timingPatterns.X {
						x--
					}
					up = !up
					continue Loop
				}
			}
			// 下移
			y++
			// timing patterns，水平
			if y == timingPatterns.Y {
				x++
				y++
				continue Loop
			}
			// alignment patterns
			for _, r := range alignPatterns {
				if y == r.Min.Y {
					// 右边向下
					if x == r.Max.X {
						x++
						for i := 0; i < 5; i++ {
							if !setColor() {
								break Loop
							}
							y++
							// timing patterns，水平
							if y == 6 {
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
							if !setColor() {
								break Loop
							}
							y++
							// timing patterns，水平
							if y == 6 {
								i++
								y++
							}
						}
						x++
						continue Loop
					}
				}
			}
			// 下边缘
			if y == qrCodeSizeTable[q.strEncoder.version] {
				// 左移
				x--
				if x < finderPatterns[2].X {
					// 左下角
					y = finderPatterns[2].Y - 1
				} else {
					y = qrCodeSizeTable[q.strEncoder.version] - 1
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
		{Max: image.Point{X: 8, Y: 8}}, // 左上
		{Min: image.Point{X: qrCodeSizeTable[q.strEncoder.version] - 8}, Max: image.Point{X: qrCodeSizeTable[q.strEncoder.version] - 1, Y: 8}}, // 右上
		{Min: image.Point{Y: qrCodeSizeTable[q.strEncoder.version] - 8}, Max: image.Point{X: 8, Y: qrCodeSizeTable[q.strEncoder.version] - 1}}, // 左下
	}
	timingPatterns := image.Point{X: 6, Y: 6}
	alignPatterns := alignmentPatternTable[q.strEncoder.version]
	versionArea := image.Rectangle{}
	// 版本大于7才有区域
	if q.strEncoder.version >= version7 {
		versionArea.Min.X = finderPatterns[1].Min.X - 3
		versionArea.Min.Y = 5
		versionArea.Max.Y = 5
	}
	// 得分
	score, minScore := 0, 0xffffffff
	idx := -1
	// 8种mark的方法
	for i := 0; i < maxMark; i++ {
		// 原始数据
		copy(q.buffer.data, q.pix)
		// 逐行
		for y := 0; y < len(q.pixXY); y++ {
		Next:
			for x := 0; x < len(q.pixXY[y]); x++ {
				idx++
				// finder patterns
				for i := 0; i < len(finderPatterns); i++ {
					if x >= finderPatterns[i].Min.X && y <= finderPatterns[i].Max.Y {
						continue Next
					}
				}
				// timing patterns
				if x == timingPatterns.X || y == timingPatterns.Y {
					continue Next
				}
				// alignment patterns
				for i := 0; i < len(alignPatterns); i++ {
					if x >= alignPatterns[i].Min.X && y <= alignPatterns[i].Max.Y {
						continue Next
					}
				}
				// 每种mark都有相应的计算公式
				if markFunc[i](x, y) {
					// 反色
					if q.buffer.data[idx] == _paletteBlack {
						q.buffer.data[idx] = _paletteWhite
					} else {
						q.buffer.data[idx] = _paletteBlack
					}
				}
			}
		}
		// 评估
		q.initMarkXY()
		score = q.evaluation1() + q.evaluation2() + q.evaluation3() + q.evaluation4()
		// 最小得分
		if score < minScore {
			minScore = score
			// 保存mark后的数据和mark编码
			copy(q.markData.data, q.buffer.data)
			q.markNum = i
		}
		idx = -1
	}
	// 最终的数据
	copy(q.pix, q.markData.data)
}

// 找到5个连续颜色的点，+3分
// 5个连续颜色的点之后，每多1个点+1分
func (q *qrCode) evaluation1() int {
	data := q.buffer.data
	score := 0
	var lastBlock uint8
	var consecutive int
	// 行
	for y := 0; y < qrCodeSizeTable[q.strEncoder.version]; y++ {
		consecutive = 0
		lastBlock = data[0]
		for x := 1; x < qrCodeSizeTable[q.strEncoder.version]; x++ {
			if data[x] == lastBlock {
				consecutive++
				if consecutive == 5 {
					score += 3
				} else if consecutive > 5 {
					score++
				}
			} else {
				lastBlock = data[x]
				consecutive = 0
			}
		}
		data = data[qrCodeSizeTable[q.strEncoder.version]:]
	}
	// 列
	data = q.buffer.data
	idx := 0
	for i := 0; i < qrCodeSizeTable[q.strEncoder.version]; i++ {
		consecutive = 0
		lastBlock = data[i]
		idx = i
		for j := 0; j < qrCodeSizeTable[q.strEncoder.version]; j++ {
			if data[idx] == lastBlock {
				consecutive++
				if consecutive == 5 {
					score += 3
				} else if consecutive > 5 {
					score++
				}
			} else {
				lastBlock = data[idx]
				consecutive = 0
			}
			idx += qrCodeSizeTable[q.strEncoder.version]
		}
	}
	return score
}

// 找到相同颜色的最小矩形（2*2），+3分
func (q *qrCode) evaluation2() int {
	score := 0
	for y := 0; y < qrCodeSizeTable[q.strEncoder.version]-1; y++ {
		i := 0
		for x := 0; x < qrCodeSizeTable[q.strEncoder.version]-1; x++ {
			i1 := i + x
			i2 := i1 + 1
			i3 := i1 + qrCodeSizeTable[q.strEncoder.version]
			i4 := i2 + qrCodeSizeTable[q.strEncoder.version]
			if q.buffer.data[i1] == q.buffer.data[i2] &&
				q.buffer.data[i1] == q.buffer.data[i3] &&
				q.buffer.data[i1] == q.buffer.data[i4] {
				score += 3
			}
		}
		i += qrCodeSizeTable[q.strEncoder.version]
	}
	return score
}

// 找到[10111010000]或者[00001011101]，+40分
func (q *qrCode) evaluation3() int {
	score := 0
	b := evaluation3Bytes[0]
	d := q.buffer.data
	o := true
	for y := 0; y < qrCodeSizeTable[q.strEncoder.version]; y++ {
		for x := 0; x < qrCodeSizeTable[q.strEncoder.version]-len(b); {
			o = true
			for i := 0; i < len(b); i++ {
				if b[i] != d[i] {
					o = false
					break
				}
			}
			if o {
				score += 40
				x += len(b)
			} else {
				x++
			}
		}
	}
	b = evaluation3Bytes[1]
	for x := 0; x < qrCodeSizeTable[q.strEncoder.version]; x++ {
		for y := 0; y < qrCodeSizeTable[q.strEncoder.version]-len(b); y++ {
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
