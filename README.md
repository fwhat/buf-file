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
>>>> write 819200 times each 512byte writeBufSize 31457280byte total 400M

testOsFileWrite: avg cost: 339.525128 ms p: 117.811604M/s
testOsWriteAll: avg cost: 34.729239 ms p: 1151.767244M/s
testBufferFileWrite: avg cost: 61.794869 ms p: 647.302931M/s
testBufioFileWrite: avg cost: 60.539214 ms p: 660.728762M/s
testBuffFileWriteAndRead: avg cost: 142.100109 ms p: 281.491690M/s

>>>> write 102400 times each 4096byte writeBufSize 31457280byte total 400M

testOsFileWrite: avg cost: 107.945281 ms p: 370.558115M/s
testOsWriteAll: avg cost: 65.054610 ms p: 614.868029M/s
testBufferFileWrite: avg cost: 60.757820 ms p: 658.351465M/s
testBufioFileWrite: avg cost: 60.297148 ms p: 663.381294M/s
testBuffFileWriteAndRead: avg cost: 100.220669 ms p: 399.119267M/s

>>>> write 10240 times each 40960byte writeBufSize 31457280byte total 400M

testOsFileWrite: avg cost: 37.746534 ms p: 1059.699936M/s
testOsWriteAll: avg cost: 52.466312 ms p: 762.393969M/s
testBufferFileWrite: avg cost: 61.121876 ms p: 654.430177M/s
testBufioFileWrite: avg cost: 60.452312 ms p: 661.678582M/s
testBuffFileWriteAndRead: avg cost: 96.610790 ms p: 414.032426M/s

```