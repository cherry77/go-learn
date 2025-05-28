package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run script.go <input.csv.tar.gz>")
	}

	inputFile := os.Args[1]
	outputDir := strings.TrimSuffix(filepath.Base(inputFile), ".tar.gz") + "_split"

	// 创建输出目录
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// 打开输入文件
	inFile, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer inFile.Close()

	// 处理 tar.gz 文件
	if err := processTarGz(inFile, outputDir); err != nil {
		log.Fatalf("Error processing file: %v", err)
	}

	log.Printf("Successfully split files into directory: %s", outputDir)
}

func processTarGz(inFile io.Reader, outputDir string) error {
	// 创建 gzip 读取器
	gzReader, err := gzip.NewReader(inFile)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	// 创建 tar 读取器
	tarReader := tar.NewReader(gzReader)

	// 处理 tar 文件中的每个文件
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("tar reader error: %v", err)
		}

		// 只处理普通文件
		if header.Typeflag != tar.TypeReg {
			continue
		}

		// 只处理 CSV 文件
		if !strings.HasSuffix(strings.ToLower(header.Name), ".csv") {
			continue
		}

		// 处理 CSV 文件
		if err := processCSV(tarReader, outputDir, filepath.Base(header.Name)); err != nil {
			return fmt.Errorf("error processing CSV: %v", err)
		}
	}

	return nil
}

func processCSV(reader io.Reader, outputDir, originalFilename string) error {
	csvReader := csv.NewReader(reader)

	// 读取标题行
	headers, err := csvReader.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers: %v", err)
	}

	// 查找列索引
	advertisingIDIndex, countryCodeIndex := -1, -1
	for i, header := range headers {
		switch strings.ToLower(header) {
		case "advertising_id":
			advertisingIDIndex = i
		case "country_code":
			countryCodeIndex = i
		}
	}

	if advertisingIDIndex == -1 || countryCodeIndex == -1 {
		return fmt.Errorf("CSV must contain 'advertising_id' and 'country_code' columns")
	}

	// 创建按国家分组的写入器
	countryWriters := make(map[string]*csv.Writer)
	defer func() {
		// 确保所有写入器都被刷新
		for _, writer := range countryWriters {
			writer.Flush()
		}
	}()

	// 处理每一行数据
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Warning: error reading record: %v", err)
			continue
		}

		countryCode := record[countryCodeIndex]
		if countryCode == "" {
			continue
		}

		writer, exists := countryWriters[countryCode]
		if !exists {
			// 创建新文件
			outputFilename := filepath.Join(outputDir, fmt.Sprintf("%s.csv", countryCode))
			outputFile, err := os.Create(outputFilename)
			if err != nil {
				return fmt.Errorf("failed to create output file: %v", err)
			}

			// 创建 CSV 写入器并写入标题行
			writer = csv.NewWriter(outputFile)
			if err := writer.Write(headers); err != nil {
				outputFile.Close()
				return fmt.Errorf("failed to write headers: %v", err)
			}

			countryWriters[countryCode] = writer
		}

		// 写入记录
		if err := writer.Write(record); err != nil {
			log.Printf("Warning: failed to write record: %v", err)
		}
	}

	return nil
}
