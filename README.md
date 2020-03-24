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
testOsFileWrite: avg cost: 1670.818854 ms p: 23.940357M/s
testBufferFileWrite: avg cost: 67.421164 ms p: 593.285511M/s
testBufioFileWrite: avg cost: 35.548261 ms p: 1125.230858M/s
testBuffFileReader: avg cost: 58.648416 ms p: 682.030358M/s
testOsFileReader: avg cost: 36.335624 ms p: 1100.848018M/s
testBuffFileWriteAndRead: avg cost: 81.141164 ms p: 492.968035M/s
testOsFileWriteAndRead: avg cost: 2300.550556 ms p: 17.387142M/s

```