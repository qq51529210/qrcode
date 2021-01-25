package qrcode

import (
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"sync"
)

var (
	encPool sync.Pool
)

func init() {
	encPool.New = func() interface{} {
		enc := new(encoder)
		return enc
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

func PNG(w io.Writer, str string, level Level, pixel int, compress png.CompressionLevel) error {
	img, err := Image(str, level, pixel)
	if err != nil {
		return err
	}
	enc := png.Encoder{
		CompressionLevel: compress,
		BufferPool:       nil,
	}
	return enc.Encode(w, img)
}

func JPEG(w io.Writer, str string, level Level, pixel, quality int) error {
	img, err := Image(str, level, pixel)
	if err != nil {
		return err
	}
	return jpeg.Encode(w, img, &jpeg.Options{Quality: quality})
}

func Image(str string, level Level, pixel int) (image.Image, error) {
	enc := encPool.Get().(*encoder)
	enc.Level = level
	enc.str = str
	err := enc.Encode()
	if err != nil {
		encPool.Put(enc)
		return nil, err
	}
	img := image.NewPaletted(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: pixel, Y: pixel},
	}, []color.Color{
		color.Black,
		color.White,
	})
	//img.Pix
	encPool.Put(enc)
	return img, nil
}
