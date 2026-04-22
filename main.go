package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		printUsage()
		return nil
	}

	if len(args) < 3 {
		printUsage()
		return errors.New("missing required arguments")
	}

	inputPath := args[0]
	outputPath := args[1]
	targetFields := args[2:]

	if err := decodeCSVFields(inputPath, outputPath, targetFields); err != nil {
		return err
	}

	return nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  csv-decode <input.csv> <output.csv> <field1> [field2 field3 ...]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -h, --help    Show this help message")
	fmt.Println("")
	fmt.Println("Example:")
	fmt.Println("  csv-decode input.csv output.csv message title")
}

func decodeCSVFields(inputPath, outputPath string, fields []string) error {
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
			record[idx] = decodeUnicodeEscapes(record[idx])
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

func decodeUnicodeEscapes(value string) string {
	var b strings.Builder
	b.Grow(len(value))

	for i := 0; i < len(value); {
		if value[i] != '\\' || i+1 >= len(value) {
			b.WriteByte(value[i])
			i++
			continue
		}

		switch value[i+1] {
		case 'u':
			if i+6 <= len(value) {
				r, err := parseHexRune(value[i+2 : i+6])
				if err == nil {
					b.WriteRune(r)
					i += 6
					continue
				}
			}
		case 'U':
			if i+10 <= len(value) {
				r, err := parseHexRune(value[i+2 : i+10])
				if err == nil {
					b.WriteRune(r)
					i += 10
					continue
				}
			}
		}

		b.WriteByte(value[i])
		i++
	}

	return b.String()
}

func parseHexRune(hex string) (rune, error) {
	v, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		return 0, err
	}

	r := rune(v)
	if !utf8.ValidRune(r) {
		return 0, fmt.Errorf("invalid rune: %U", r)
	}

	return r, nil
}
