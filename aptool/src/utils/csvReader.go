package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
)

type CsvReader struct {
	CsvFilePath string
}

func (c *CsvReader) ReadContent() chan []string {
	var (
		err      error
		f        *os.File
		contents [][]string
		data     chan []string
	)
	f, err = os.Open(c.CsvFilePath)
	if err != nil {
		panic(fmt.Sprintf("csv:[%s] open fail!!", filepath.Base(c.CsvFilePath)))
	}
	reader := csv.NewReader(f)
	contents, err = reader.ReadAll()
	if err != nil {
		panic("Error reading csv file content!")
	}
	data = make(chan []string)
	go func() {
		for i, v := range contents {
			if i == 0 {
				continue
			}
			data <- v
		}

		close(data)
	}()
	return data
}
