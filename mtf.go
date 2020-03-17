package mof

import (
  "os"
  "fmt"
  "path/filepath"
  "strings"
  "io/ioutil"
  "bufio"
  "strconv"
)

// MTF: Make Two Files.
//
// This does the packing of files in to two files.
// It expects a directory for its input. this would contain all the files or subdirectories
// with files you want to archive.
// It also expects another directory for the output of itself.
func MTF(dir, outDir string) (string, string, error) {
  indexFilePath := filepath.Join(outDir, filepath.Base(dir) + ".f1")
  indexFile, err := os.OpenFile(indexFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
  if err != nil {
    return "", "", err
  }
  defer indexFile.Close()

  dataFilePath := filepath.Join(outDir, filepath.Base(dir) + ".f2")
  dataFile, err := os.OpenFile(dataFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
  if err != nil {
    return "", "", err
  }
  defer dataFile.Close()

  err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }
    if ! info.IsDir() {
      if ! strings.HasSuffix(dir, "/") {
        dir += "/"
      }

      data, err := ioutil.ReadFile(path)
      if err != nil {
        return nil
      }
      dataFile.Write(data)

      pathToWrite := strings.Replace(path, dir, "", 1)
      outStr := fmt.Sprintf("%s,,,%d\n", pathToWrite, info.Size())
      indexFile.WriteString(outStr)
    }
    return nil
  })

  return indexFilePath, dataFilePath, nil
}


// UndoMTF extracts the files to their full names ( directory and file name) in the folder
// supplied as the outDir
func UndoMTF(indexFilePath, dataFilePath, outDir string) error {
  indexFile, err := os.Open(indexFilePath)
  if err != nil {
    return err
  }
  defer indexFile.Close()

  dataFile, err := os.Open(dataFilePath)
  if err != nil {
    return err
  }
  defer dataFile.Close()

  scanner := bufio.NewScanner(indexFile)
  var seekSize int64 = 0
  for scanner.Scan() {
    line := scanner.Text()
    lineParts := strings.Split(line, ",,,")
    sizeInt64, err := strconv.ParseInt(lineParts[1], 10, 64)
    if err != nil {
      return err
    }
    outData := make([]byte, sizeInt64)
    if seekSize == 0 {
      _, err := dataFile.Read(outData)
      if err != nil {
        return err
      }
    } else {
      _, err := dataFile.ReadAt(outData, seekSize)
      if err != nil {
        return err
      }
    }
    seekSize += sizeInt64

    outFilePath := filepath.Join(outDir, strings.Replace(filepath.Base(indexFilePath), ".f1", "", 1), lineParts[0])
    // make missing directories
    err = os.MkdirAll(filepath.Dir(outFilePath), 0777)
    if err != nil {
      return err
    }

    outFile, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
    if err != nil {
      return err
    }
    outFile.Write(outData)
    outFile.Close()
  }

  if err = scanner.Err(); err != nil {
    return err
  }

  return nil
}
