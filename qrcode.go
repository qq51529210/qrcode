package qrcode

import (
	"image"
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
		return new(encoder)
	}
}

func PNG(w io.Writer, str string, level Level, pixel int) error {
	img, err := Image(str, level, pixel)
	if err != nil {
		return err
	}
	return png.Encode(w, img)
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
	img := image.NewRGBA(image.Rectangle{
		Min: image.Point{},
		Max: image.Point{X: pixel, Y: pixel},
	})
	encPool.Put(enc)
	return img, nil
}
