# qrcode
二维码生成开发包

## 使用

```go
import (
	"github.com/qq51529210/qrcode"
)

func main() {
  var out bytes.Buffer
  err := qrcode.PNG(&out, "Hello World!", qrcode.LevelL, png.BestCompression)
  if err != nil {
    panic(err)
  }
  err = qrcode.JPEG(&out, "Hello World!", qrcode.LevelQ, 100)
  if err != nil {
    panic(err)
  }
}
```

## 测试

下面是与“github.com/skip2/go-qrcode”包的benchmark

```go
goos: darwin
goarch: amd64
pkg: github.com/qq51529210/qrcode
BenchmarkImageMy-4                  5775            194493 ns/op            1504 B/op          2 allocs/op
BenchmarkImageSkip2-4               2247            476112 ns/op          159376 B/op       3348 allocs/op
PASS
ok      github.com/qq51529210/qrcode    2.822s
```
