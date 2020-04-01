# A File With A Buffer Space

### BufferWriter

```go
package main

import (
    "github.com/Dowte/buf-file"
    "errors"
    "io"
    "os"
)
var s = []byte("test buf file the readWrite performance.")

func main () {
   // file not exists will create
    bufFile, err := buf_file.NewBufFile("/tmp/buf_file", 1024*1024*4)
    
    if err != nil {
        panic(err)
    }
    
    writer, err := bufFile.GetWriter()
    
    if err != nil {
        panic(err)
    }
    // do not forget Close writer
    defer writer.Close()
    
    for i := 0; i < 1024000; i++ {
        writer.Write(s)
    }
}
```

### BufferReader

```go
package main

import (
    "github.com/Dowte/buf-file"
    "errors"
    "io"
    "os"
)
var s = []byte("test buf file the readWrite performance.")

func main () {
    // file not exists will create
    bufFile, err := buf_file.NewBufFile("/tmp/buf_file", 1024*1024*4)
    
    if err != nil {
        panic(err)
    }
    
    writer, err := bufFile.GetWriter()
    
    if err != nil {
        panic(err)
    }
    
    go func() {
        // do not forget to close writer
    	defer writer.Close()
        
        for i := 0; i < 1024000; i++ {
            writer.Write(s)
        }
    }()
    
    reader, err := bufFile.GetReader()

    if err != nil {
        panic(err)
    }
    // do not forget to close reader
    defer reader.Close()

    var readContent []byte
    var offset int64
    for true {
        buf := make([]byte, 1024*10)
        at, err := reader.ReadAt(buf, offset)
        offset += int64(at)
        if at > 0 {
            readContent = append(readContent, buf[:at]...)
        }
        if len(readContent) == 40960000 {
            break
        }
        if err != nil {
            panic(err)
        }
    }
}
```

#### test

```
>>>> write 102400 times each 4096byte writeBufSize 31457280byte total 400M

testOsFileWrite: avg cost: 109.787377 ms p: 364.340612M/s
testOsWriteAll: avg cost: 56.966013 ms p: 702.173062M/s
testBufferFileWrite: avg cost: 60.949642 ms p: 656.279491M/s
testBufioFileWrite: avg cost: 60.766215 ms p: 658.260515M/s
testOsFileReader: avg cost: 561.259386 ms p: 71.268296M/s
testBuffFileReader: avg cost: 448.621274 ms p: 89.162067M/s
testBuffFileWriteAndRead: avg cost: 72.040656 ms p: 555.242026M/s

>>>> write 10240 times each 40960byte writeBufSize 31457280byte total 400M

testOsFileWrite: avg cost: 62.626169 ms p: 638.710637M/s
testOsWriteAll: avg cost: 62.516247 ms p: 639.833671M/s
testBufferFileWrite: avg cost: 58.992151 ms p: 678.056305M/s
testBufioFileWrite: avg cost: 59.194844 ms p: 675.734527M/s
testOsFileReader: avg cost: 433.153303 ms p: 92.346058M/s
testBuffFileReader: avg cost: 414.256288 ms p: 96.558583M/s
testBuffFileWriteAndRead: avg cost: 74.541580 ms p: 536.613255M/s
```