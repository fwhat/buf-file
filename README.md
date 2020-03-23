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
testOsFileWrite: avg cost: 1773.500701 ms p: 22.554262M/s
testBufferFileWrite: avg cost: 68.483454 ms p: 584.082692M/s
testBuffFileReader: avg cost: 47.458180 ms p: 842.847320M/s
testOsFileReader: avg cost: 36.782501 ms p: 1087.473639M/s
testBuffFileWriteAndRead: avg cost: 80.396841 ms p: 497.531984M/s
testOsFileWriteAndRead: avg cost: 2473.529887 ms p: 16.171222M/s
```