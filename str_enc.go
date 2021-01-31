package qrcode

var (
	// 编码字符串函数
	strEncFunc = [maxMode]func(*strEncoder){
		encNumericStr,
		encAlphanumericStr,
		encByteStr,
		encKanJiStr,
	}
)

type strEncoder struct {
	str     string  // 原始字符串
	buff    *buffer // 共享缓存，在字节编码和交错会用到
	bitD    []byte  // 编码的数据
	bitN    int     // 最后一个字节的bit个数
	version         // 版本
	Level           // 纠错级别
	mode            // 选择的模式
}

// 添加bit，c是小端字节，n是bit的个数
func (e *strEncoder) appendBit(c byte, n int) {
	i := 8 - e.bitN
	if n < i {
		e.bitD[len(e.bitD)-1] |= c << (i - n)
		e.bitN += n
		return
	}
	if n == i {
		e.bitD[len(e.bitD)-1] |= c
		e.bitD = append(e.bitD, 0)
		e.bitN = 0
		return
	}
	j := n - i
	e.bitD[len(e.bitD)-1] |= c >> j
	e.bitD = append(e.bitD, 0)
	e.bitD[len(e.bitD)-1] |= c << i
	e.bitN = j
}

// 编码
func (e *strEncoder) Encode(str string, level Level) error {
	e.Level = level
	e.str = str
	// 确定编码模式
	e.mode = analysisMode(e.str)
	// 确定最小版本
	var err error
	e.version, err = analysisVersion(e.str, e.Level, e.mode)
	if err != nil {
		return err
	}
	// 准备编码
	e.bitD = e.bitD[:1]
	e.bitN = 0
	// 指示器
	e.encIndicator()
	// 字符串长度
	e.encStrLength()
	// 字符串数据
	strEncFunc[e.mode](e)
	// 填充字节
	e.appendPadBytes()
	return nil
}

// 编码指示器
func (e *strEncoder) encIndicator() {
	e.appendBit(indicatorTable[e.mode], 4)
}

// 编码字符串长度
func (e *strEncoder) encStrLength() {
	n := uint16(len(e.str))
	if e.version <= 8 {
		// v1-9，10，9，8，8
		switch e.mode {
		case numericMode:
			e.appendBit(byte(n>>8), 2)
			e.appendBit(byte(n), 8)
		case alphanumericMode:
			e.appendBit(byte(n>>8), 1)
			e.appendBit(byte(n), 8)
		case byteMode:
			e.appendBit(byte(n), 8)
		case kanJiMode:
			e.appendBit(byte(n), 8)
		}
	} else if e.version >= 9 && e.version <= 25 {
		// v10-26，12，11，16，10
		switch e.mode {
		case numericMode:
			e.appendBit(byte(n>>8), 4)
			e.appendBit(byte(n), 8)
		case alphanumericMode:
			e.appendBit(byte(n>>8), 3)
			e.appendBit(byte(n), 8)
		case byteMode:
			e.appendBit(byte(n>>8), 8)
			e.appendBit(byte(n), 8)
		case kanJiMode:
			e.appendBit(byte(n>>8), 2)
			e.appendBit(byte(n), 8)
		}
	} else {
		// v27-40，14，13，16，12
		switch e.mode {
		case numericMode:
			e.appendBit(byte(n>>8), 6)
			e.appendBit(byte(n), 8)
		case alphanumericMode:
			e.appendBit(byte(n>>8), 5)
			e.appendBit(byte(n), 8)
		case byteMode:
			e.appendBit(byte(n>>8), 8)
			e.appendBit(byte(n), 8)
		case kanJiMode:
			e.appendBit(byte(n>>8), 4)
			e.appendBit(byte(n), 8)
		}
	}
}

// 调整编码的数据大小
func (e *strEncoder) appendPadBytes() {
	if len(e.bitD) < errorCorrectionTable[e.version][e.Level].TotalBytes {
		if e.bitN > 4 {
			e.bitD = append(e.bitD, 0)
		}
	}
	e.bitN = 0
	for {
		if len(e.bitD) >= errorCorrectionTable[e.version][e.Level].TotalBytes {
			return
		}
		e.bitD = append(e.bitD, 236)
		if len(e.bitD) >= errorCorrectionTable[e.version][e.Level].TotalBytes {
			return
		}
		e.bitD = append(e.bitD, 17)
	}
}

// 数字模式编码
func encNumericStr(e *strEncoder) {
	// 将字符分组，3个（10bit），2个（7bit），1个（4bit）
	i := 0
	var n int16
	for i < len(e.str) {
		switch len(e.str[i:]) {
		case 1:
			n = int16(e.str[i] - '0')
			i++
			e.appendBit(byte(n), 4)
		case 2:
			n = int16(e.str[i]-'0') * 10
			i++
			n += int16(e.str[i] - '0')
			i++
			e.appendBit(byte(n), 7)
		default:
			n = int16(e.str[i]-'0') * 100
			i++
			n += int16(e.str[i]-'0') * 10
			i++
			n += int16(e.str[i] - '0')
			i++
			e.appendBit(byte(n>>8), 2)
			e.appendBit(byte(n), 8)
		}
	}
}

// 字母模式编码
func encAlphanumericStr(e *strEncoder) {
	// 两个字符一组，alphanumericTable[0]*45+alphanumericTable[1]，(11bit)
	i1, i2 := 0, 1
	var n uint16
	for i2 < len(e.str) {
		n = uint16(alphanumericTable[e.str[i1]])*45 + uint16(alphanumericTable[e.str[i2]])
		e.appendBit(byte(n>>8), 3)
		e.appendBit(byte(n), 8)
		i1 += 2
		i2 += 2
	}
	// 如果1个字符，6bit
	if i1 < len(e.str) {
		e.appendBit(alphanumericTable[e.str[i1]], 6)
	}
}

// 字节模式编码
func encByteStr(e *strEncoder) {
	e.buff.Reset()
	e.buff.data = append(e.buff.data, e.str...)
	for i := 0; i < len(e.buff.data); i++ {
		e.appendBit(e.buff.data[i], 8)
	}
}

// 日文模式编码
func encKanJiStr(e *strEncoder) {
	var m uint16
	for _, c := range e.str {
		if uint16(c) <= 0x9FFC {
			// 减去0x8140
			c = c - 0x8140
		} else {
			// 减去0xC140
			c = c - 0xC140
		}
		// 结果（高字节*0xC0+低字节）是一个13bit的数
		m = uint16(c>>8)*0xC0 + uint16(c)
		e.appendBit(byte(m>>8), 5)
		e.appendBit(byte(m), 8)
	}
}
