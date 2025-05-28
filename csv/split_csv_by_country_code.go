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
)

func main() {
	// 输入文件
	inputFile := "ams_0501_0523_ip.csv.tar.gz"

	// 打开 tar.gz 文件
	f, err := os.Open(inputFile)
	if err != nil {
		log.Fatalf("无法打开文件: %v", err)
	}
	defer f.Close()

	// 创建 gzip 读取器
	gzr, err := gzip.NewReader(f)
	if err != nil {
		log.Fatalf("无法创建 gzip 读取器: %v", err)
	}
	defer gzr.Close()

	// 创建 tar 读取器
	tr := tar.NewReader(gzr)

	// 读取 tar 文件中的第一个文件（假设只有一个 CSV 文件）
	header, err := tr.Next()
	if err != nil {
		log.Fatalf("无法读取 tar 文件: %v", err)
	}

	if header.Typeflag != tar.TypeReg {
		log.Fatalf("不是常规文件: %v", header.Name)
	}

	// 创建 CSV 读取器
	csvr := csv.NewReader(tr)

	// 读取 CSV 头部
	headerRow, err := csvr.Read()
	if err != nil {
		log.Fatalf("无法读取 CSV 头部: %v", err)
	}

	// 查找 country_code 列的索引
	countryCodeIndex := -1
	for i, col := range headerRow {
		if col == "country_code" {
			countryCodeIndex = i
			break
		}
	}
	if countryCodeIndex == -1 {
		log.Fatal("CSV 文件中没有 country_code 列")
	}

	// 创建输出目录
	outputDir := "output"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("无法创建输出目录: %v", err)
	}

	// 用于存储不同国家的 CSV 写入器
	writers := make(map[string]*csv.Writer)
	files := make(map[string]*os.File)

	// 读取并处理 CSV 数据
	for {
		record, err := csvr.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("无法读取 CSV 记录: %v", err)
			continue
		}

		countryCode := record[countryCodeIndex]
		if countryCode == "" {
			countryCode = "unknown"
		}

		// 如果该国家的写入器不存在，则创建一个
		if _, exists := writers[countryCode]; !exists {
			outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.csv", countryCode))
			f, err := os.Create(outputFile)
			if err != nil {
				log.Printf("无法创建输出文件 %s: %v", outputFile, err)
				continue
			}
			files[countryCode] = f
			writer := csv.NewWriter(f)
			writers[countryCode] = writer

			// 写入头部
			if err := writer.Write(headerRow); err != nil {
				log.Printf("无法写入 CSV 头部到 %s: %v", outputFile, err)
			}
		}

		// 写入记录
		if err := writers[countryCode].Write(record); err != nil {
			log.Printf("无法写入 CSV 记录到 %s.csv: %v", countryCode, err)
		}
	}

	// 刷新所有写入器并关闭文件
	for countryCode, writer := range writers {
		writer.Flush()
		if err := writer.Error(); err != nil {
			log.Printf("刷新写入器 %s 时出错: %v", countryCode, err)
		}
		if err := files[countryCode].Close(); err != nil {
			log.Printf("关闭文件 %s 时出错: %v", countryCode, err)
		}
	}

	fmt.Println("CSV 文件拆分完成")
}
