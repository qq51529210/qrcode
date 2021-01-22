package qrcode

import (
	"io"
	"strconv"
)

var (
	encFunc = [maxMode]func(*encoder){
		encNumeric, encAlphanumeric, encByte, encKanJi,
	} // 编码函数
)

type encoder struct {
	str     string // 原始字符串
	buf     []byte // 数据
	bit     int    // buf最后一个字节的bit的数量
	level          // 纠错级别
	mode           // 选择的模式
	version        // 版本
}

func (enc *encoder) Encode(writer io.Writer, str string, level level) (err error) {
	// 确定编码模式
	enc.mode = analysisMode(str)
	// 确定最小版本
	enc.version, err = analysisVersion(level, enc.mode, str)
	if err != nil {
		return err
	}
	// 数据
	enc.buf = enc.buf[:0]
	enc.str = str
	enc.level = level
	// 指示器
	enc.buf = append(enc.buf, enc.mode.Indicator())
	enc.bit = 4
	// 字符串长度
	enc.encLength()
	// 编码
	encFunc[enc.mode](enc)
	enc.growBuff()
	// 纠错
	enc.ec()
	// 输出结果
	_, err = writer.Write(enc.buf)
	return err
}

// 追加多少位，n是小端模式的字节，nBit表示位数
func (enc *encoder) append(n byte, nBit int) {
	m := 8 - enc.bit
	if nBit < m {
		enc.buf[len(enc.buf)-1] |= n<<m - nBit
		enc.bit += nBit
		return
	}
	if nBit == m {
		enc.buf[len(enc.buf)-1] |= n
		enc.buf = append(enc.buf, 0)
		enc.bit = 0
		return
	}
	m = nBit - m
	enc.buf[len(enc.buf)-1] |= n >> m
	enc.buf = append(enc.buf, 0)
	enc.buf[len(enc.buf)-1] |= n << enc.bit
	enc.bit = m
}

// 编码字符串长度
func (enc *encoder) encLength() {
	//n := len(enc.str)
	// v1-9，10，9，8，8
	// v10-26，12，11，16，10
	// v27-40，14，13，16，12
	if enc.version >= 0 && enc.version <= 8 {

	} else if enc.version >= 9 && enc.version <= 25 {

	} else {

	}
	encFunc[enc.mode](enc)
}

// 调整编码的数据大小
func (enc *encoder) growBuff() {
	for {
		if len(enc.buf) >= versionECTable[enc.version][enc.mode].TotalBytes {
			return
		}
		enc.buf = append(enc.buf, 236)
		if len(enc.buf) >= versionECTable[enc.version][enc.mode].TotalBytes {
			return
		}
		enc.buf = append(enc.buf, 17)
	}
}

// 纠错
func (enc *encoder) ec() {

}

// 数字模式编码
func encNumeric(enc *encoder) {
	// 将字符分组，3个（10bit），2个（7bit），1个（4bit）
	i1, i2 := 0, 3
	var n int64
	for i2 < len(enc.str) {
		n, _ = strconv.ParseInt(enc.str[i1:i2], 10, 64)
		enc.append(byte(n>>8), 2)
		enc.append(byte(n), 8)
		i1 = i2
		i2 += 3
	}
	if i1 < len(enc.str) {
		n, _ = strconv.ParseInt(enc.str[i1:], 10, 64)
		switch len(enc.str[i1:]) {
		case 1:
			enc.append(byte(n), 4)
		case 2:
			enc.append(byte(n), 7)
		default:
			enc.append(byte(n>>8), 2)
			enc.append(byte(n), 8)
		}
	}
}

// 字母模式编码
func encAlphanumeric(enc *encoder) {
	// 两个字符一组，alphanumericTable[0]*45+alphanumericTable[1]，(11bit)
	i1, i2 := 0, 1
	var n uint16
	for i2 < len(enc.str) {
		n = uint16(alphanumericModeTable[enc.str[i1]])*45 + uint16(alphanumericModeTable[enc.str[i2]])
		enc.append(byte(n>>8), 3)
		enc.append(byte(n), 8)
		i1 += 2
		i2 += 2
	}
	// 如果1个字符，6bit
	if i1 < len(enc.str) {
		enc.append(alphanumericModeTable[enc.str[i1]], 6)
	}
}

// 字节模式编码
func encByte(enc *encoder) {
	for _, c := range enc.str {
		enc.append(byte(c), 8)
	}
}

// 日文模式编码
func encKanJi(enc *encoder) {
	var m uint16
	for _, c := range enc.str {
		if uint16(c) <= 0x9FFC {
			// 减去0x8140
			c = c - 0x8140
		} else {
			// 减去0xC140
			c = c - 0xC140
		}
		// 结果（高字节*0xC0+低字节）是一个13bit的数
		m = uint16(c>>8)*0xC0 + uint16(c)
		enc.append(byte(m>>8), 5)
		enc.append(byte(m), 8)
	}
}
