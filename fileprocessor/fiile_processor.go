package fileprocessor

import (
	"encoding/csv"
	"errors"
	"io"
	"mime/multipart"
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
	finalColumnIndex := -1

	// Read the CSV headers
	headers, err := reader.Read()
	if err != nil {
		return nil, err
	}

	// Find the index of the final column
	if len(headers) > 0 {
		finalColumnIndex = len(headers) - 1
	}

	recordList := []string{headers[finalColumnIndex]}

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		finalColumnValue := ""
		if finalColumnIndex >= 0 && finalColumnIndex < len(row) {
			finalColumnValue = row[finalColumnIndex]
		}

		recordList = append(recordList, finalColumnValue)

	}

	return recordList, err
}

/*
Get file processor depending on filetype
*/
func GetFileProcessor(fileType string) (FileProcessor, error) {
	switch fileType {
	case "csv":
		return &CSVProcessor{}, nil
	default:
		return nil, errors.New("unsupported file type")
	}
}
