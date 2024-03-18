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
	recordList := make([]string, 0)

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// Handle varying numbers of columns in rows
		if len(row) > 0 {
			finalColumnValue := row[len(row)-1]
			recordList = append(recordList, finalColumnValue)
		}
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
