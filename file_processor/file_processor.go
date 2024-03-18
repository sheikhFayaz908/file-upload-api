package fileprocessor

import (
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
	"strings"
)

type FileProcessor interface {
	ProcessFile(file *multipart.FileHeader) ([]string, error)
}

type CSVProcessor struct{}

func (cp *CSVProcessor) ProcessFile(file *multipart.FileHeader) ([]string, error) {
	uploadedFile, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer uploadedFile.Close()

	reader := csv.NewReader(uploadedFile)
	// Handle quoted fields with lazy quotes
	reader.LazyQuotes = true

	recordList := make([]string, 0)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}

		var finalColumnValue string
		if len(row) > 0 {
			finalColumnValue = row[len(row)-1]
		}

		// Handle cells that are quoted and split across lines
		if strings.HasPrefix(finalColumnValue, `"`) && !strings.HasSuffix(finalColumnValue, `"`) {
			for {
				nextRow, err := reader.Read()
				if err != nil {
					if err == io.EOF {
						return nil, errors.New("unterminated quoted cell across lines")
					}
					return nil, err
				}

				finalColumnValue += "\n" + nextRow[len(nextRow)-1]
				if strings.HasSuffix(nextRow[len(nextRow)-1], `"`) {
					break
				}
			}
		}

		recordList = append(recordList, finalColumnValue)
	}

	return recordList, nil
}

func GetFileProcessor(fileType string) (FileProcessor, error) {
	switch fileType {
	case "csv":
		return &CSVProcessor{}, nil
	default:
		return nil, errors.New("unsupported file type")
	}
}
