package service

import (
	repository "a21hc3NpZ25tZW50/repository/fileRepository"
	"encoding/csv"
	"errors"
	"strings"
)

type FileService struct {
	Repo *repository.FileRepository
}

func (s *FileService) ProcessFile(fileContent string) (map[string][]string, error) {
	// return nil, nil // TODO: replace this
	if strings.TrimSpace(fileContent) == "" {
		return nil, errors.New("file content is empty")
	}

	reader := csv.NewReader(strings.NewReader(fileContent))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("failed to parse CSV file")
	}

	if len(records) < 2 {
		return nil, errors.New("invalid CSV format: insufficient data")
	}

	headers := records[0]
	result := make(map[string][]string)

	for _, header := range headers {
		result[header] = []string{}
	}

	for i := 1; i < len(records); i++ {
		row := records[i]
		if len(row) != len(headers) {
			return nil, errors.New("invalid CSV format: row does not match header length")
		}

		for j, value := range row {
			header := headers[j]
			result[header] = append(result[header], value)
		}
	}

	return result, nil
}
