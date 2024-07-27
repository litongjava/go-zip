package main

import (
  "flag"
  "fmt"
  "strings"
)

func main() {
  // Define command line parameters
  zipFile := flag.String("r", "", "Archive zip file name")
  excludeFile := flag.String("x", "", "Exclude file pattern")
  flag.Parse()

  // Check if there are enough parameters
  //if *zipFile == "" && flag.NArg() < 1 {
  //  log.Fatalln("Usage: go-zip -r archive.zip [-x exclude_pattern] [directory_or_file]")
  //}

  // Get extra non-flag parameters, i.e., directory or file name
  var sourcePath string
  if flag.NArg() < 1 {
    sourcePath = "."
  } else {
    sourcePath = flag.Arg(0)
  }

  // Set zip file name and directory
  err := Zip(*zipFile, sourcePath, strings.Split(*excludeFile, ""))
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println("File zipped successfully!")
  }
}
