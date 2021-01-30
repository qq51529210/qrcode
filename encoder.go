package qrcode

import (
	"fmt"
	"unicode"
)

var (
	strEncodeFunc = [maxMode]func(*strEncoder){
		encNumeric, encAlphanumeric, encByte, encKanJi,
	} // 编码函数
)

type strEncoder struct {
	str     string // 原始字符串
	version byte   // 版本
	size    int    // 大小
	level   Level  // 纠错级别
	mode    mode   // 选择的模式
	data    bit    // 原始字符串编码后的数据
	buff    []byte // 字节编码使用的缓存
	poly    poly   // 纠错多项式
}

// 编码
func (enc *strEncoder) Encode() error {
	// 确定编码模式
	enc.Mode()
	// 确定最小版本
	err := enc.Version()
	if err != nil {
		return err
	}
	// 图像大小
	enc.size = int(enc.version)*4 + 21
	// 准备编码
	enc.data.Reset()
	// 写入指示器
	enc.Indicator()
	// 写入字符串长度
	enc.Length()
	// 写入字符串数据
	strEncodeFunc[enc.mode](enc)
	// 填充编码后的长度
	enc.AddPadBytes()
	return nil
}

// 判断编码模式
func (enc *strEncoder) Mode() {
	enc.mode = numericMode
	for _, c := range enc.str {
		if unicode.MaxLatin1 < c {
			if (c >= 0x8140 && c <= 0x9FFC) || (c >= 0xE040 && c <= 0xEBBF) {
				enc.mode = kanJiMode
			} else {
				enc.mode = byteMode
				return
			}
		} else {
			if c >= '0' && c <= '9' {
				continue
			}
			if alphanumericTable[c] != 0 {
				if enc.mode < alphanumericMode {
					enc.mode = alphanumericMode
				}
			} else {
				enc.mode = byteMode
				return
			}
		}
	}
}

// 判断编码版本
func (enc *strEncoder) Version() error {
	for i, a := range versionCapacity[enc.level][enc.mode] {
		if len(enc.str) <= a {
			enc.version = byte(i)
			return nil
		}
	}
	return fmt.Errorf("string length <%d> too lager", len(enc.str))
}

// 编码指示器
func (enc *strEncoder) Indicator() {
	enc.data.b[0] = modeIndicator[enc.mode]
	enc.data.n = 4
}

// 编码字符串长度
func (enc *strEncoder) Length() {
	n := uint16(len(enc.str))
	if enc.version <= 8 {
		// v1-9，10，9，8，8
		switch enc.mode {
		case numericMode:
			enc.data.Append(byte(n>>8), 2)
			enc.data.Append(byte(n), 8)
		case alphanumericMode:
			enc.data.Append(byte(n>>8), 1)
			enc.data.Append(byte(n), 8)
		case byteMode:
			enc.data.Append(byte(n), 8)
		case kanJiMode:
			enc.data.Append(byte(n), 8)
		}
	} else if enc.version >= 9 && enc.version <= 25 {
		// v10-26，12，11，16，10
		switch enc.mode {
		case numericMode:
			enc.data.Append(byte(n>>8), 4)
			enc.data.Append(byte(n), 8)
		case alphanumericMode:
			enc.data.Append(byte(n>>8), 3)
			enc.data.Append(byte(n), 8)
		case byteMode:
			enc.data.Append(byte(n>>8), 8)
			enc.data.Append(byte(n), 8)
		case kanJiMode:
			enc.data.Append(byte(n>>8), 2)
			enc.data.Append(byte(n), 8)
		}
	} else {
		// v27-40，14，13，16，12
		switch enc.mode {
		case numericMode:
			enc.data.Append(byte(n>>8), 6)
			enc.data.Append(byte(n), 8)
		case alphanumericMode:
			enc.data.Append(byte(n>>8), 5)
			enc.data.Append(byte(n), 8)
		case byteMode:
			enc.data.Append(byte(n>>8), 8)
			enc.data.Append(byte(n), 8)
		case kanJiMode:
			enc.data.Append(byte(n>>8), 4)
			enc.data.Append(byte(n), 8)
		}
	}
}

// 调整编码的数据大小
func (enc *strEncoder) AddPadBytes() {
	if len(enc.data.b) < versionECTable[enc.version][enc.level].TotalBytes {
		if enc.data.n > 4 {
			enc.data.b = append(enc.data.b, 0)
		}
	}
	enc.data.n = 0
	for {
		if len(enc.data.b) >= versionECTable[enc.version][enc.level].TotalBytes {
			return
		}
		enc.data.b = append(enc.data.b, 236)
		if len(enc.data.b) >= versionECTable[enc.version][enc.level].TotalBytes {
			return
		}
		enc.data.b = append(enc.data.b, 17)
	}
}

// 纠错
func (enc *strEncoder) ECC() {
	// 版本纠错表
	ect := versionECTable[enc.version][enc.level]
	// 生成多项式
	enc.poly.Gen(ect.BlockECBytes)
	// 纠错编码
	data := enc.data.Bytes()
	for i := 0; i < ect.Group1Block; i++ {
		enc.data.b = append(enc.data.b, enc.poly.Encode(data[:ect.Group1BlockBytes])...)
		data = data[ect.Group1BlockBytes:]
	}
	for i := 0; i < ect.Group2Block; i++ {
		enc.data.b = append(enc.data.b, enc.poly.Encode(data[:ect.Group2BlockBytes])...)
		data = data[ect.Group2BlockBytes:]
	}
}

// 交错
func (enc *strEncoder) Interleave() []byte {
	// 版本纠错表
	ect := versionECTable[enc.version][enc.level]
	// 交错
	if ect.Group2Block > 0 {
		enc.buff = resetBytes(enc.buff, len(enc.data.b))
		idx, col := 0, 0
		groupTotalBytes := ect.Group1Block * ect.Group1Block
		for col < ect.Group2BlockBytes {
			if col < ect.Group1BlockBytes {
				for i := 0; i < ect.Group1Block; i++ {
					enc.buff[idx] = enc.data.b[col+i*ect.Group1BlockBytes]
					idx++
				}
			}
			for i := 0; i < ect.Group2Block; i++ {
				enc.buff[idx] = enc.data.b[col+groupTotalBytes+i*ect.Group2BlockBytes]
				idx++
			}
			col++
		}
		col = 0
		for col < ect.BlockECBytes {
			for i := 0; i < ect.Group1Block; i++ {
				enc.buff[idx] = enc.data.b[col+i*ect.BlockECBytes]
				idx++
			}
			for i := 0; i < ect.Group2Block; i++ {
				enc.buff[idx] = enc.data.b[col+(ect.Group1Block+i)*ect.BlockECBytes]
				idx++
			}
			col++
		}
	}
	return enc.buff
}

// 数字模式编码
func encNumeric(enc *strEncoder) {
	// 将字符分组，3个（10bit），2个（7bit），1个（4bit）
	i := 0
	var n int16
	for i < len(enc.str) {
		switch len(enc.str[i:]) {
		case 1:
			n = int16(enc.str[i] - '0')
			i++
			enc.data.Append(byte(n), 4)
		case 2:
			n = int16(enc.str[i]-'0') * 10
			i++
			n += int16(enc.str[i] - '0')
			i++
			enc.data.Append(byte(n), 7)
		default:
			n = int16(enc.str[i]-'0') * 100
			i++
			n += int16(enc.str[i]-'0') * 10
			i++
			n += int16(enc.str[i] - '0')
			i++
			enc.data.Append(byte(n>>8), 2)
			enc.data.Append(byte(n), 8)
		}
	}
}

// 字母模式编码
func encAlphanumeric(enc *strEncoder) {
	// 两个字符一组，alphanumericTable[0]*45+alphanumericTable[1]，(11bit)
	i1, i2 := 0, 1
	var n uint16
	for i2 < len(enc.str) {
		n = uint16(alphanumericTable[enc.str[i1]])*45 + uint16(alphanumericTable[enc.str[i2]])
		enc.data.Append(byte(n>>8), 3)
		enc.data.Append(byte(n), 8)
		i1 += 2
		i2 += 2
	}
	// 如果1个字符，6bit
	if i1 < len(enc.str) {
		enc.data.Append(alphanumericTable[enc.str[i1]], 6)
	}
}

// 字节模式编码
func encByte(enc *strEncoder) {
	enc.buff = enc.buff[:0]
	enc.buff = append(enc.buff, enc.str...)
	for i := 0; i < len(enc.buff); i++ {
		enc.data.Append(enc.buff[i], 8)
	}
}

// 日文模式编码
func encKanJi(enc *strEncoder) {
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
		enc.data.Append(byte(m>>8), 5)
		enc.data.Append(byte(m), 8)
	}
}
