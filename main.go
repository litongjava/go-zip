package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

func main() {
	i := len(os.Args)
	if i < 2 {
		log.Fatalln("usage zip file.zip file")
	}
	sourceDir := os.Args[2]
	target := os.Args[1]
	// 设置压缩包名称和目录
	log.Println(target, sourceDir)
	err := Zip(sourceDir, target)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("File zipped successfully!")
	}
}

func Zip(sourceDir, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	base := filepath.Base(sourceDir)
	filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		// 获取相对路径
		relPath, err := filepath.Rel(sourceDir, path)
		if err != nil {
			return err
		}
		if relPath == "." {
			return nil
		} else if strings.Contains(relPath, string(os.PathSeparator)) {
			relPath = strings.Replace(relPath, string(os.PathSeparator), "/", len(relPath))
		}
		//处理中文编码
		if IsChineseChar(relPath) {
			relPath = GetChineseName(relPath)
		}

		header.Name = base + "/" + relPath

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

// 或者封装函数调用
func IsChineseChar(str string) bool {
	for _, r := range str {
		compile := regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]")
		if unicode.Is(unicode.Scripts["Han"], r) || (compile.MatchString(string(r))) {
			return true
		}
	}
	return false
}

//获取中文名称
func GetChineseName(filename string) string {
	i := bytes.NewReader([]byte(filename))
	decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewEncoder())
	content, _ := ioutil.ReadAll(decoder)
	filename = string(content)
	return filename
}
