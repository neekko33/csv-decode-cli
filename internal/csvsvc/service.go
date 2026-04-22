package csvsvc

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	decoder "csv-decode-cli/internal/unicode"
)

var ErrDestinationExists = errors.New("destination file already exists")

func ReadHeaders(inputPath string) ([]string, error) {
	f, err := os.Open(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open input CSV: %w", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	headers, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}
	return headers, nil
}

func DefaultOutputPath(inputPath string) string {
	dir := filepath.Dir(inputPath)
	base := filepath.Base(inputPath)
	ext := filepath.Ext(base)
	name := strings.TrimSuffix(base, ext)
	if ext == "" {
		return filepath.Join(dir, name+"-decoded")
	}
	return filepath.Join(dir, name+"-decoded"+ext)
}

func ValidateDestination(path string, allowOverwrite bool) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("output path cannot be empty")
	}
	exists, err := FileExists(path)
	if err != nil {
		return err
	}
	if exists && !allowOverwrite {
		return ErrDestinationExists
	}
	return nil
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check file %q: %w", path, err)
}

func DecodeCSVFields(inputPath, outputPath string, fields []string, allowOverwrite bool) error {
	if filepath.Clean(inputPath) == filepath.Clean(outputPath) {
		return errors.New("input and output paths cannot be the same")
	}
	if err := ValidateDestination(outputPath, allowOverwrite); err != nil {
		return err
	}

	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input CSV: %w", err)
	}
	defer inFile.Close()

	reader := csv.NewReader(inFile)
	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %w", err)
	}

	fieldIndexes, err := resolveFieldIndexes(headers, fields)
	if err != nil {
		return err
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output CSV: %w", err)
	}
	defer outFile.Close()

	writer := csv.NewWriter(outFile)
	defer writer.Flush()

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("failed to write CSV header: %w", err)
	}

	for {
		record, readErr := reader.Read()
		if errors.Is(readErr, io.EOF) {
			break
		}
		if readErr != nil {
			return fmt.Errorf("failed to read CSV record: %w", readErr)
		}

		for _, idx := range fieldIndexes {
			if idx >= len(record) {
				continue
			}
			record[idx] = decoder.DecodeEscapes(record[idx])
		}

		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to flush CSV writer: %w", err)
	}
	return nil
}

func resolveFieldIndexes(headers []string, fields []string) ([]int, error) {
	headerToIndex := make(map[string]int, len(headers))
	for i, header := range headers {
		headerToIndex[header] = i
	}

	indexes := make([]int, 0, len(fields))
	for _, field := range fields {
		idx, ok := headerToIndex[field]
		if !ok {
			return nil, fmt.Errorf("field %q not found in CSV header", field)
		}
		indexes = append(indexes, idx)
	}
	return indexes, nil
}
