package scan

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

// devNull is a ByteReader that provides 0 bytes of input
type devNull struct{}

func (*devNull) ReadByte() (byte, error) {
	return 0, io.EOF
}

type filesReader struct {
	filenames []string
	nextFile  int
	currFile  io.ByteReader
}

func NewFilesReader(filenames []string) io.ByteReader {
	reader := &filesReader{
		filenames: filenames,
		nextFile:  0,
		currFile:  &devNull{},
	}
	return reader
}

func (r *filesReader) openNextFile() {
	filename := r.filenames[r.nextFile]
	r.nextFile++
	if f, err := os.Open(filename); err != nil {
		fmt.Fprintf(os.Stderr, "go-lox: %s\n", err)
		os.Exit(66) // see "sysexits.h"
	} else {
		r.currFile = bufio.NewReader(f)
	}
}

func (r *filesReader) ReadByte() (b byte, err error) {
	for {
		if b, err = r.currFile.ReadByte(); err == nil {
			return
		}
		if err == io.EOF && r.nextFile < len(r.filenames) {
			return
		}
		r.openNextFile()
	}
}
