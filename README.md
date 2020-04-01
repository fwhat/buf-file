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
# write 102400 times each 4k total 400M

testOsFileWrite: avg cost: 111.123039 ms p: 359.961358M/s
testBufferFileWrite: avg cost: 59.361948 ms p: 673.832336M/s
testBufioFileWrite: avg cost: 70.071467 ms p: 570.845760M/s
testOsFileReader: avg cost: 610.278888 ms p: 65.543804M/s
testBuffFileReader: avg cost: 483.647542 ms p: 82.704855M/s
testBuffFileWriteAndRead: avg cost: 95.381160 ms p: 419.370030M/s

```