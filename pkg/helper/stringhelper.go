package helper

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
)

// LineCounter count number of line in string/file
func LineCounter(r io.Reader) (int, error) {
    buf := make([]byte, 32*1024)
    count := 0
    lineSep := []byte{'\n'}

    for {
        c, err := r.Read(buf)
        count = count + bytes.Count(buf[:c], lineSep)
        switch {
            case err == io.EOF:
            return count, nil
            case err != nil:
            return count, err
        }
    }
}

// ReadLines read batch number of lines in bufio
func ReadLines(scanner *bufio.Scanner, numOfLines int) (lines []string, err error) {
    i := 0
    for scanner.Scan() {
        i++
        if i <= numOfLines {
            lines = append(lines, scanner.Text())
            if i == numOfLines {
                return lines, scanner.Err()
            }
        }
    }

    return lines, io.EOF
}

// WriteText write string to txt file
func WriteText(text string, path string) error {
    fileHandle, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND, 0666)

    if err != nil {
        return err
    }
    defer fileHandle.Close()

    writer := bufio.NewWriter(fileHandle)
    fmt.Fprintln(writer, text)
    
    return writer.Flush()
}
