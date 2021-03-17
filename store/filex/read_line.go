package filex

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ReadFileByLine(fileName string) ([]string, error) {
	fd, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer fd.Close()
	return read(fd)
}

func read(fd io.Reader) ([]string, error) {

	reader := bufio.NewReader(fd)

	conts := make([]string, 0, 10)

	for {

		line, err := reader.ReadString('\n')
		if nil != err && err != io.EOF {
			return conts, err
		}

		line = strings.TrimSpace(line)
		if io.EOF == err {
			if len(line) <= 0 {
				return conts, nil
			}
		}

		if len(line) <= 0 {
			continue
		}

		conts = append(conts, line)
	}
}
