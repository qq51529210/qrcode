package qrcode

type bit struct {
	b []byte
	n int
}

func (b *bit) Bytes() []byte {
	return b.b
}

func (b *bit) Append(c byte, n int) {
	i := 8 - b.n
	if n < i {
		b.b[len(b.b)-1] |= c << (i - n)
		b.n += n
		return
	}
	if n == i {
		b.b[len(b.b)-1] |= c
		b.b = append(b.b, 0)
		b.n = 0
		return
	}
	j := n - i
	b.b[len(b.b)-1] |= c >> j
	b.b = append(b.b, 0)
	b.b[len(b.b)-1] |= c << i
	b.n = j
}

func (b *bit) Reset() {
	b.b = b.b[:0]
	b.b = append(b.b, 0)
	b.n = 0
}
