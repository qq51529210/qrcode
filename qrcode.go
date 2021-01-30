package qrcode

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"sync"
)

const (
	maxMark = 8
)

var (
	_pool    sync.Pool
	_palette = color.Palette{
		color.Black,
		color.White,
		color.RGBA{R: 255, A: 255},
		color.RGBA{G: 255, A: 255},
		color.RGBA{B: 255, A: 255},
	}
	//_palette      = color.Palette{color.Black, color.White}
	_paletteBlack  = uint8(_palette.Index(color.Black))
	_paletteWhite  = uint8(_palette.Index(color.White))
	_paletteR      = uint8(_palette.Index(color.RGBA{R: 255, A: 255}))
	_paletteG      = uint8(_palette.Index(color.RGBA{G: 255, A: 255}))
	_paletteB      = uint8(_palette.Index(color.RGBA{B: 255, A: 255}))
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
)

func init() {
	_pool.New = func() interface{} {
		return new(qrCode)
	}
}

func resetBytes(b []byte, n int) []byte {
	if cap(b) < n {
		b = make([]byte, n, n)
	} else {
		b = b[:n]
		for i := 0; i < n; i++ {
			b[i] = 0
		}
	}
	return b
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
	q.str.level = level
	q.str.str = str
	// 字符串编码
	err := q.str.Encode()
	if err != nil {
		_pool.Put(q)
		return nil, err
	}
	// 纠错
	q.str.ECC()
	// 交错，得到最终的数据
	q.data = q.str.Interleave()
	// 位图
	img := image.NewPaletted(image.Rect(0, 0, q.str.size, q.str.size), _palette)
	q.img = img
	// 开始画图
	q.Draw()
	q.img = nil
	// 回收缓存
	_pool.Put(q)
	// 返回
	return img, err
}

type qrCode struct {
	str  strEncoder      // 字符串编码器
	data []byte          // 交错后的数据，指针
	img  *image.Paletted // 图
	mark int             // 使用的mark图
}

// 画图
func (q *qrCode) Draw() {
	// 所有的像素设为白色
	for i := 0; i < len(q.img.Pix); i++ {
		q.img.Pix[i] = _paletteWhite
	}
	//q.test()
	q.DrawFinderPatterns()
	q.DrawTimingPatterns()
	q.DrawAlignmentPatterns()
	q.DrawBottomLeftPoint()
	q.DrawData()
	q.Mark()
	q.DrawFormat()
	q.DrawVersion()
}

func (q *qrCode) test() {
	for x := 0; x <= 8; x++ {
		q.img.SetColorIndex(x, 8, _paletteB)
	}
	for y := 0; y <= 8; y++ {
		q.img.SetColorIndex(8, y, _paletteB)
	}

	for x := 0; x <= 8; x++ {
		q.img.SetColorIndex(x+q.img.Rect.Max.X-8, 8, _paletteB)
	}
	for x := 0; x < 3; x++ {
		for y := 0; y < 6; y++ {
			q.img.SetColorIndex(x+q.img.Rect.Max.X-11, y, _paletteB)
		}
	}

	for y := 0; y < 3; y++ {
		for x := 0; x < 6; x++ {
			q.img.SetColorIndex(x, y+q.img.Rect.Max.Y-11, _paletteB)
		}
	}
	for y := 0; y <= 8; y++ {
		q.img.SetColorIndex(8, q.img.Rect.Max.Y-7+y, _paletteB)
	}
}

func (q *qrCode) DrawRectangle(x1, y1, x2, y2 int, c uint8, fill bool) {
	if fill {
		for i := y1; i <= y2; i++ {
			for j := x1; j <= x2; j++ {
				q.img.SetColorIndex(j, i, c)
			}
		}
		return
	}
	// 上下
	for i := x1; i <= x2; i++ {
		q.img.SetColorIndex(i, y1, c)
		q.img.SetColorIndex(i, y2, c)
	}
	// 左右
	for i := y1 + 1; i < y2; i++ {
		q.img.SetColorIndex(x1, i, c)
		q.img.SetColorIndex(x2, i, c)
	}
}

func (q *qrCode) DrawFinderPatterns() {
	// 左上角
	q.DrawRectangle(0, 0, 6, 6, _paletteBlack, false)
	q.DrawRectangle(2, 2, 4, 4, _paletteBlack, true)
	// 右上角
	q.DrawRectangle(q.img.Rect.Max.X-7, 0, q.img.Rect.Max.X-1, 6, _paletteBlack, false)
	q.DrawRectangle(q.img.Rect.Max.X-5, 2, q.img.Rect.Max.X-3, 4, _paletteBlack, true)
	// 左下角
	q.DrawRectangle(0, q.img.Rect.Max.Y-7, 6, q.img.Rect.Max.Y-1, _paletteBlack, false)
	q.DrawRectangle(2, q.img.Rect.Max.Y-5, 4, q.img.Rect.Max.Y-3, _paletteBlack, true)
}

func (q *qrCode) DrawTimingPatterns() {
	// 水平
	for i := 8; i < q.img.Rect.Max.X-8; {
		q.img.SetColorIndex(i, 6, _paletteBlack)
		i += 2
	}
	// 垂直
	for i := 8; i < q.img.Rect.Max.Y-8; {
		q.img.SetColorIndex(6, i, _paletteBlack)
		i += 2
	}
}

func (q *qrCode) DrawAlignmentPatterns() {
	for _, r := range versionAlignmentTable[q.str.version] {
		q.DrawRectangle(r.Min.X, r.Min.Y, r.Max.X, r.Max.Y, _paletteBlack, false)
		q.img.SetColorIndex(r.Min.X+2, r.Min.Y+2, _paletteBlack)
	}
}

func (q *qrCode) DrawBottomLeftPoint() {
	// y=version*4+4+9，
	q.img.SetColorIndex(8, int(q.str.version)*4+13, _paletteBlack)
}

func (q *qrCode) DrawFormat() {
	f := formatBitTable[q.str.level][q.mark]
	idx := 0
	// 左上角
	for x := 0; x < 6; x++ {
		if f[idx] == 1 {
			q.img.SetColorIndex(x, 8, _paletteBlack)
		}
		idx++
	}
	if f[idx] == 1 {
		q.img.SetColorIndex(7, 8, _paletteBlack)
	}
	idx++
	if f[idx] == 1 {
		q.img.SetColorIndex(8, 8, _paletteBlack)
	}
	idx++
	if f[idx] == 1 {
		q.img.SetColorIndex(8, 7, _paletteBlack)
	}
	idx++
	for y := 5; y >= 0; y-- {
		if f[idx] == 1 {
			q.img.SetColorIndex(8, y, _paletteBlack)
		}
		idx++
	}
	idx = 0
	// 左下角
	for y := q.img.Rect.Max.Y - 1; y > q.img.Rect.Max.Y-8; y-- {
		if f[idx] == 1 {
			q.img.SetColorIndex(8, y, _paletteBlack)
		}
		idx++
	}
	// 右上角
	for x := q.img.Rect.Max.X - 8; x <= q.img.Rect.Max.X-1; x++ {
		if f[idx] == 1 {
			q.img.SetColorIndex(x, 8, _paletteBlack)
		}
		idx++
	}
}

func (q *qrCode) DrawVersion() {
	if q.str.version < 6 {
		return
	}
	ver := versionBitTable[q.str.version]
	x := q.img.Rect.Max.X - 11
	y := q.img.Rect.Max.Y - 8
	idx := 0
	for i := 0; i < 6; i++ {
		for j := 0; j < 3; j++ {
			if ver[idx] == 1 {
				// 左下角
				q.img.SetColorIndex(i, j+y, _paletteBlack)
				q.img.SetColorIndex(i+x, j, _paletteBlack)
			}
			idx++
		}
	}
}

func (q *qrCode) DrawData() {
	// 从右下角开始
	x := q.img.Rect.Max.X - 1
	y := q.img.Rect.Max.Y - 1
	// finder patterns，左上0，右上1，左下2
	finderPatterns := [3]image.Point{{8, 8}, {x - 7, 8}, {8, y - 7}}
	// align patterns 矩形
	alignPatterns := versionAlignmentTable[q.str.version]
	// timing patterns
	timingPatterns := image.Point{X: 6, Y: 6}
	// version patterns，0右上x，1左下y
	versionPatterns := image.Point{X: finderPatterns[1].X - 3, Y: finderPatterns[2].Y - 3}
	idx := 0
	bit := byte(0b10000000)
	char := q.data[idx]
	up := true
	setColor := func() bool {

		if char&bit != 0 {
			q.img.SetColorIndex(x, y, _paletteBlack)
		}
		bit >>= 1
		if bit == 0 {
			bit = 0b10000000
			idx++
			if idx == len(q.data) {
				return false
			}
			char = q.data[idx]
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
			// 上移
			y--
			// finder patterns
			if y == finderPatterns[1].Y {
				// 右上
				if x >= finderPatterns[1].X {
					x--
					y++
					up = !up
					continue Loop
				}
				// 左上
				if x <= finderPatterns[0].X {
					x--
					// timing patterns，垂直
					if x == timingPatterns.X {
						x--
					}
					y++
					up = !up
					continue Loop
				}
				// 左下，不可能
			}
			// timing patterns，水平
			if y == timingPatterns.Y {
				// 版本7以上
				if q.str.version >= 6 && x >= versionPatterns.X {
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
				} else {
					// 上移
					if x > finderPatterns[0].X {
						y--
					} else {
						x--
						y++
						up = !up
						continue Loop
					}
				}
			}
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
			// 上边缘
			if y < 0 {
				// 左移
				x--
				y++
				// 向下
				up = !up
				continue Loop
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
			// 下移
			y++
			// finder patterns
			if x < timingPatterns.X {
				// 左下
				if q.str.version >= 6 {
					if y == versionPatterns.Y {
						y--
						x--
						up = !up
						continue Loop
					}
				} else {
					if y == finderPatterns[2].Y {
						y--
						x--
						up = !up
						continue Loop
					}
				}
			} else if x < finderPatterns[2].X {
				if y == q.img.Rect.Max.Y-8 {
					y++
					x++
					continue Loop
				}
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
			// timing patterns，水平
			if y == timingPatterns.Y {
				x++
				y++
				continue Loop
			}
			// 下边缘
			if y > q.img.Rect.Max.Y {
				// 左移
				x--
				if x == finderPatterns[2].X {
					// 左下角
					y = finderPatterns[2].Y - 1
				} else {
					y = q.img.Rect.Max.Y - 1
				}
				// 向下
				up = !up
				continue Loop
			}
			x++
		}
	}
}

// 检查并mark
func (q *qrCode) Mark() {

}
