package main

import (
	"buf_file"
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"time"
)

const hash = "2ab6d209dbc6c58d44af9e2cec539f40"

// 40 char
var s = []byte("test buf file the writeFile performance.")
var readContent []byte
// 1024000 * 40 = 40M
func writeFile(writer io.Writer)  {
	for i := 0; i < 1024000; i ++ {
		writer.Write(s)
	}
}


func readFile (read io.ReaderAt)  {
	readContent = []byte{}
	var offset int64
	for true {
		buf := make([]byte, 1024 * 10)
		at, err := read.ReadAt(buf, offset)
		offset += int64(at)
		if at > 0 {
			readContent = append(readContent, buf[:at]...)
		}
		if len(readContent) == 40960000 {
			break
		}
		if err != nil {
			if err == io.EOF {
				continue
			}
			panic(err)
		}
	}
}

func testOsFileWrite ()  {
	// file writer avg cost: 1609.280375 ms 24.85 M/s
	start := time.Now()
	for i:=0; i < 10; i ++ {
		os.Remove("/tmp/buf_file")
		file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		writeFile(file)
		file.Close()
	}
	end := time.Now()

	fmt.Printf("avg cost: %f ms", float64(end.Sub(start).Nanoseconds()) / 1000/ 1000 / 10)
}

func testOsFileReader ()  {
	// 40M avg cost: 74.660214 ms
	start := time.Now()
	for i:=0; i < 10; i ++ {
		file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		readFile(file)
		file.Close()
	}
	end := time.Now()

	fmt.Printf("avg cost: %f ms", float64(end.Sub(start).Nanoseconds()) / 1000/ 1000 / 10)
}

func testBuffFileReader ()  {
	// 40M avg cost: 81.650695 ms
	start := time.Now()
	for i:=0; i < 10; i ++ {
		file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		bufFile := buf_file.NewBufFile(file, 1024*1024*4)
		readFile(bufFile)
		file.Close()
	}
	end := time.Now()

	fmt.Printf("avg cost: %f ms", float64(end.Sub(start).Nanoseconds()) / 1000/ 1000 / 10)
}

func testBufferFileWrite ()  {
	// buffer file writer avg cost: 10k, 38.335672 ms 1043.41 M/s
	// buffer file writer avg cost: 1M, 32.305683 ms x
	// buffer file writer avg cost: 4M, 37.540536 ms x
	// buffer file writer avg cost: 10M, 38.428869 ms x
	start := time.Now()
	for i:=0; i < 10; i ++ {
		os.Remove("/tmp/buf_file")

		file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		bufFile := buf_file.NewBufFile(file, 1024 * 1024 * 4)

		writeFile(bufFile)
		bufFile.Flush()
		file.Close()
	}
	end := time.Now()

	fmt.Printf("avg cost: %f ms", float64(end.Sub(start).Nanoseconds()) / 1000/ 1000 / 10)
}

func testOsFileWriteAndRead ()  {
	// 40M avg cost: 2263.035093 ms
	start := time.Now()
	for i:=0; i < 10; i ++ {
		os.Remove("/tmp/buf_file")
		file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		go writeFile(file)
		readFile(file)
		file.Close()
	}
	end := time.Now()

	fmt.Printf("avg cost: %f ms", float64(end.Sub(start).Nanoseconds()) / 1000/ 1000 / 10)
}

func testBuffFileWriteAndRead ()  {
	// 40M  avg cost: 65.301589 ms
	start := time.Now()
	for i:=0; i < 10; i ++ {
		os.Remove("/tmp/buf_file")
		file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			panic(err)
		}

		bufFile := buf_file.NewBufFile(file, 1024 * 1024 * 4)

		go writeFile(bufFile)
		readFile(bufFile)
		bufFile.Flush()
		file.Close()
	}
	end := time.Now()

	fmt.Printf("avg cost: %f ms", float64(end.Sub(start).Nanoseconds()) / 1000/ 1000 / 10)
}

func compareWrite ()  {
	//testOsFileWrite()
	testBufferFileWrite()

	file, err := os.OpenFile("/tmp/buf_file", os.O_CREATE|os.O_RDWR, 0644)

	if err != nil {
		panic(err)
	}
	stat, _ := file.Stat()
	if stat.Size() != 40960000 {
		panic("file size error")
	}
}

func compareRead ()  {
	//testOsFileReader()
	testBuffFileReader()

	if len(readContent) != 40960000 {
		panic("file size error")
	}
}

func compareWriteAndRead ()  {
	//testOsFileWriteAndRead()
	testBuffFileWriteAndRead()

	if hashx(readContent) != hash {
		panic("content hash error")
	}
	if len(readContent) != 40960000 {
		panic("file size error")
	}
}
// content md5 2ab6d209dbc6c58d44af9e2cec539f40

func hashx(TestString []byte) string {
	Md5Inst := md5.New()
	Md5Inst.Write(TestString)
	Result := Md5Inst.Sum([]byte(""))

	return fmt.Sprintf("%x", Result)
}

func main()  {
	compareWriteAndRead()
}