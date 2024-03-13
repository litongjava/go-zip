package main

import (
  "archive/zip"
  "bytes"
  "flag"
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
  // 定义命令行参数
  zipFile := flag.String("r", "", "Archive zip file name")
  excludeFile := flag.String("x", "", "Exclude file pattern")
  flag.Parse()

  // 检查是否有足够的参数
  if *zipFile == "" || flag.NArg() < 1 {
    log.Fatalln("Usage: go-zip -r archive.zip directory")
  }

  // 获取额外的非flag参数，即目录名称
  sourceDir := flag.Arg(0)

  // 输出解析得到的参数值，用于验证
  log.Printf("Zip File: %s, Exclude Pattern: %s, Source Directory: %s\n", *zipFile, *excludeFile, sourceDir)

  // 设置压缩包名称和目录
  err := Zip(sourceDir, zipFile, excludeFile)
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("File zipped successfully!")
  }
}

func Zip(sourceDir string, target *string, excludeFile *string) error {
  zipfile, err := os.Create(*target)
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

    // 获取相对路径
    relPath, err := filepath.Rel(sourceDir, path)

    if relPath == "." || (excludeFile != nil && matchExcludeFile(relPath, *excludeFile)) {
      return nil // 跳过根目录或匹配排除模式的文件
    }

    header, err := zip.FileInfoHeader(info)
    if err != nil {
      return err
    }
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

// matchExcludeFile 检查文件名是否匹配排除模式
func matchExcludeFile(filename string, pattern string) bool {
  matched, err := filepath.Match(pattern, filename)
  if err != nil {
    log.Printf("Error matching file with pattern: %v", err)
    return false
  }
  return matched
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

//对中文文件进行编码
func GetChineseName(filename string) string {
  reader := bytes.NewReader([]byte(filename))
  encoder := transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder())
  content, _ := ioutil.ReadAll(encoder)
  return string(content)
}
