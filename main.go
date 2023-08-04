package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// 定义命令行参数
	inputPath := flag.String("i", "", "Input directory path")
	outputPath := flag.String("o", "", "Output directory path")
	copyMode := flag.Bool("c", false, "Copy files instead of move")
	flag.Parse()

	// 检查参数是否有效
	if *inputPath == "" || *outputPath == "" {
		fmt.Println("Usage: -i <input_directory> -o <output_directory> [-c]")
		return
	}

	err := filepath.Walk(*inputPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 判断是否为文件，并且是 .mp4 格式
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".MP4") {
			// 获取文件的父目录路径
			parentDirPath := filepath.Dir(path)
			// 获取父目录的最后一层名称
			parentDirName := filepath.Base(parentDirPath)

			// 解析年和月
			year := info.Name()[0:4]
			month := info.Name()[4:6]

			// 构建输出路径
			outputYearPath := filepath.Join(*outputPath, year)
			outputMonthPath := filepath.Join(outputYearPath, month)
			outputVideoPath := filepath.Join(outputMonthPath, parentDirName)
			outputFilePath := filepath.Join(outputVideoPath, info.Name())
			// 创建目录
			err := os.MkdirAll(outputVideoPath, os.ModePerm)
			if err != nil {
				return err
			}

			// 复制或移动文件内容
			if *copyMode {
				err = copyFile(path, outputFilePath)
			} else {
				err = moveFile(path, outputFilePath)
			}
			if err != nil {
				return err
			}

			// 保留原始修改时间和创建时间
			err = os.Chtimes(outputFilePath, info.ModTime(), info.ModTime())
			if err != nil {
				return err
			}

			// 输出处理信息
			action := "moved"
			if *copyMode {
				action = "copied"
			}
			fmt.Printf("Processed videos - %s %s to %s\n", action, path, outputFilePath)
		}

		return nil
	})

	if err != nil {
		fmt.Println("Error:", err)
	}
}

func copyFile(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	return err
}

func moveFile(src, dest string) error {
	err := os.Rename(src, dest)
	return err
}
