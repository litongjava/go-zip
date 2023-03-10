package main

import (
  "archive/zip"
  "fmt"
  "io"
  "log"
  "os"
  "path/filepath"
)

func main() {
  i := len(os.Args)
  if i < 2 {
    log.Fatalln("usage zip file.zip file")
  }
  source := os.Args[2]
  target := os.Args[1]
  log.Println(target, source)
  err := Zip(source, target)
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("File zipped successfully!")
  }
}

func Zip(source, target string) error {
  zipfile, err := os.Create(target)
  if err != nil {
    return err
  }
  defer zipfile.Close()

  archive := zip.NewWriter(zipfile)
  defer archive.Close()

  filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }

    header, err := zip.FileInfoHeader(info)
    if err != nil {
      return err
    }
    header.Name = path

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
