package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/Dowte/buf-file"
	"io"
	"os"
	"time"
)

const hash = "6273962f98e8f3591daa492de531a209"
const buffSize = 1024 * 1024 * 10

// 64 char
var s = []byte("test buf file the readWrite performance when write 400 * 1048576")
var readContent []byte
var repeatStr []byte // 4096 char

func init() {
	for i := 0; i < 64; i++ {
		repeatStr = append(repeatStr, s...)
	}
}

// 4096 *1024 * 100 = 400M 419430400
func writeFile(writer io.Writer) error {
	// 1024 * 100
	for i := 0; i < 102400; i++ {
		_, err := writer.Write(repeatStr)
		if err != nil {
			return err
		}
	}

	return nil
}

func readFileWithReadAt(read io.ReaderAt) {
	readContent = []byte{}
	var offset int64
	for true {
		buf := make([]byte, 1024*10)
		at, err := read.ReadAt(buf, offset)
		offset += int64(at)
		if at > 0 {
			readContent = append(readContent, buf[:at]...)
		}
		if len(readContent) == 419430400 {
			break
		}
		if err != nil {
			panic(err)
		}
	}
}

func testOsFileWrite() {
	start := time.Now()
	os.Remove("/tmp/buf_file")
	file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	writeFile(file)
	file.Close()
	end := time.Now()

	echo("testOsFileWrite", end.Sub(start), 10)
}

func testBufioFileWrite() {
	os.Remove("/tmp/buf_file")
	start := time.Now()
	file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	size := bufio.NewWriterSize(file, buffSize)

	writeFile(size)
	size.Flush()
	file.Close()
	end := time.Now()

	echo("testBufioFileWrite", end.Sub(start), 10)
}

func testOsFileReader() {
	// 40M avg cost: 74.660214 ms
	start := time.Now()
	for i := 0; i < 10; i++ {
		file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		readFileWithReadAt(file)
		file.Close()
	}
	end := time.Now()

	echo("testOsFileReader", end.Sub(start), 10)
}

func testBuffFileReader() {
	start := time.Now()
	for i := 0; i < 10; i++ {
		file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE, 0644)
		if err != nil {
			panic(err)
		}

		bufFile, err := buf_file.NewBufFile("/tmp/buf_file", buffSize)
		if err != nil {
			panic(err)
		}

		reader, err := bufFile.GetReader()

		if err != nil {
			panic(err)
		}

		readFileWithReadAt(reader)
		// reader need close
		reader.Close()
		file.Close()
	}
	end := time.Now()

	echo("testBuffFileReader", end.Sub(start), 10)

	if len(readContent) != 419430400 {
		panic("file size error")
	}
	if hashx(readContent) != hash {
		panic("content hash error")
	}
}

func testBufferFileWrite() {
	os.Remove("/tmp/buf_file")
	start := time.Now()

	file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	bufFile, err := buf_file.NewBufFile("/tmp/buf_file", buffSize)

	if err != nil {
		panic(err)
	}

	file.Close()
	writer, err := bufFile.GetWriter()

	if err != nil {
		panic(err)
	}
	writeFileAndClose(writer)
	end := time.Now()

	echo("testBufferFileWrite", end.Sub(start), 10)
}

func writeFileAndClose(writer *buf_file.BufFileWriter) {
	writeFile(writer)
	writer.Close()
}

func testBuffFileWriteAndRead() {
	os.Remove("/tmp/buf_file")
	start := time.Now()
	file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		panic(err)
	}

	bufFile, err := buf_file.NewBufFile("/tmp/buf_file", buffSize)

	if err != nil {
		panic(err)
	}

	writer, err := bufFile.GetWriter()

	if err != nil {
		panic(err)
	}

	file.Close()

	reader, err := bufFile.GetReader()

	if err != nil {
		panic(err)
	}

	go writeFileAndClose(writer)
	readFileWithReadAt(reader)
	end := time.Now()

	echo("testBuffFileWriteAndRead", end.Sub(start), 10)

	if hashx(readContent) != hash {
		panic("content hash error")
	}
	if len(readContent) != 419430400 {
		panic("file size error")
	}
}

func echo(key string, time time.Duration, count int) {
	ms := time.Seconds() * 1000 / float64(count)

	fmt.Printf("%s: avg cost: %f ms p: %fM/s\n", key, ms, 40/ms*1000)
}

func compareWrite() {
	testOsFileWrite()

	file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		panic(err)
	}
	stat, _ := file.Stat()
	if stat.Size() != 419430400 {
		panic("file size error")
	}

	testBufferFileWrite()

	file2, err2 := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)

	if err2 != nil {
		panic(err2)
	}
	stat2, _ := file2.Stat()
	if stat2.Size() != 419430400 {
		panic("file size error")
	}

	testBufioFileWrite()
	file3, err3 := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)

	if err3 != nil {
		panic(err3)
	}
	stat3, _ := file3.Stat()
	if stat3.Size() != 419430400 {
		panic("file size error")
	}
}

// content md5 27c8c5fc3875c1ab804885a2cccc8bcb

func hashx(TestString []byte) string {
	Md5Inst := md5.New()
	Md5Inst.Write(TestString)
	Result := Md5Inst.Sum(nil)

	return fmt.Sprintf("%x", Result)
}

func main() {
	compareWrite()
	testOsFileReader()
	testBuffFileReader()
	testBuffFileWriteAndRead()
}
