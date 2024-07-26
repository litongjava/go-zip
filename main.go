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
  // Define command line parameters
  zipFile := flag.String("r", "", "Archive zip file name")
  excludeFile := flag.String("x", "", "Exclude file pattern")
  flag.Parse()

  // Check if there are enough parameters
  if *zipFile == "" || flag.NArg() < 1 {
    log.Fatalln("Usage: go-zip -r archive.zip directory_or_file")
  }

  // Get extra non-flag parameters, i.e., directory or file name
  sourcePath := flag.Arg(0)

  // Output the parsed parameter values for validation
  log.Printf("Zip File: %s, Exclude Pattern: %s, Source Path: %s\n", *zipFile, *excludeFile, sourcePath)

  // Set zip file name and directory
  err := Zip(zipFile, sourcePath, excludeFile)
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("File zipped successfully!")
  }
}

func Zip(target *string, sourcePath string, excludeFile *string) error {
  zipfile, err := os.Create(*target)
  if err != nil {
    return err
  }
  defer zipfile.Close()

  archive := zip.NewWriter(zipfile)
  defer archive.Close()

  info, err := os.Stat(sourcePath)
  if err != nil {
    return err
  }

  var base string
  if info.IsDir() {
    base = filepath.Base(sourcePath)
  } else {
    base = filepath.Base(filepath.Dir(sourcePath))
  }

  if info.IsDir() {
    err = filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
      return addFileToZip(archive, path, sourcePath, base, excludeFile, info)
    })
  } else {
    err = addFileToZip(archive, sourcePath, filepath.Dir(sourcePath), base, excludeFile, info)
  }

  return err
}

func addFileToZip(archive *zip.Writer, path, sourcePath, base string, excludeFile *string, info os.FileInfo) error {
  if info.IsDir() && path == sourcePath {
    return nil
  }

  // Get relative path
  relPath, err := filepath.Rel(sourcePath, path)
  if err != nil {
    return err
  }

  if excludeFile != nil && matchExcludeFile(relPath, *excludeFile) {
    return nil // Skip files matching the exclude pattern
  }

  header, err := zip.FileInfoHeader(info)
  if err != nil {
    return err
  }

  if strings.Contains(relPath, string(os.PathSeparator)) {
    relPath = strings.Replace(relPath, string(os.PathSeparator), "/", -1)
  }

  // Handle Chinese encoding
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
}

// matchExcludeFile checks if the filename matches the exclude pattern
func matchExcludeFile(filename string, pattern string) bool {
  matched, err := filepath.Match(pattern, filename)
  if err != nil {
    log.Printf("Error matching file with pattern: %v", err)
    return false
  }
  return matched
}

// IsChineseChar checks if the string contains Chinese characters
func IsChineseChar(str string) bool {
  for _, r := range str {
    compile := regexp.MustCompile("[\u3002\uff1b\uff0c\uff1a\u201c\u201d\uff08\uff09\u3001\uff1f\u300a\u300b]")
    if unicode.Is(unicode.Scripts["Han"], r) || compile.MatchString(string(r)) {
      return true
    }
  }
  return false
}

// GetChineseName encodes the filename in GB18030
func GetChineseName(filename string) string {
  reader := bytes.NewReader([]byte(filename))
  encoder := transform.NewReader(reader, simplifiedchinese.GB18030.NewEncoder())
  content, _ := ioutil.ReadAll(encoder)
  return string(content)
}
