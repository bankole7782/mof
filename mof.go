package mof

import (
  "os"
  "path/filepath"
  "fmt"
  "io/ioutil"
  "strconv"
  "strings"
  "math/rand"
  "time"
)


func randomString(length int) string {
  var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
  const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}


// MOF : Make One File
// This first calls MTF to get two files and then combine them into one file.
//
// It expects a directory for its input. this would contain all the files or subdirectories
// with files you want to archive.
//
// It expects a path which ends with the name of the archive.
func MOF(dir, archivePath string) error {
  outDir := filepath.Dir(archivePath)
  indexFilePath, dataFilePath, err := MTF(dir, outDir)
  if err != nil {
    return err
  }

  outFilePath := archivePath
  indexFileFileInfo, err := os.Stat(indexFilePath)
  if err != nil {
    return err
  }

  outFile, err := os.OpenFile(outFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
  if err != nil {
    return err
  }
  defer outFile.Close()

  outFile.WriteString(fmt.Sprintf("%d\n", indexFileFileInfo.Size()))
  mtfData, err := ioutil.ReadFile(indexFilePath)
  if err != nil {
    return err
  }
  outFile.Write(mtfData)

  mtfData1, err := ioutil.ReadFile(dataFilePath)
  if err != nil {
    return err
  }
  outFile.Write(mtfData1)

  // delete the mtf files
  err = os.Remove(indexFilePath)
  if err != nil {
    return err
  }
  err = os.Remove(dataFilePath)
  if err != nil {
    return err
  }
  return nil
}


// UndoMOF undoes archives created with the MOF function
//
// It first undoes the archive into two files and then passes the two files to the UndoMTF function.
// It expects the path of the mof file created with the MOF function and the outDir as the output directory.
func UndoMOF(mofFilePath, outDir string) error {
  sizeOfIndexFileBytes := make([]byte, 0)
  mofFile, err := os.Open(mofFilePath)
  if err != nil {
    return err
  }

  // read out the size of the indexFile
  var i int64
  for {
    charBytes := make([]byte, 1)
    if i == 0 {
      _, err = mofFile.Read(charBytes)
      if err != nil {
        return err
      }
    } else {
      _, err = mofFile.ReadAt(charBytes, i)
      if err != nil {
        return err
      }
    }

    sizeOfIndexFileBytes = append(sizeOfIndexFileBytes, charBytes...)
    if string(charBytes) == "\n" {
      break
    }
  }

  trueSizeOfIndexFileBytes := sizeOfIndexFileBytes[0: len(sizeOfIndexFileBytes) - 1] // remove the new line
  sizeOfIndexFile, err := strconv.Atoi(string(trueSizeOfIndexFileBytes))
  if err != nil {
    return err
  }

  var startReadingIndex int64 = int64(len(sizeOfIndexFileBytes))
  indexFileData := make([]byte, sizeOfIndexFile)
  _, err = mofFile.ReadAt(indexFileData, startReadingIndex)
  if err != nil {
    return err
  }

  tmpFolder := filepath.Join("/tmp", fmt.Sprintf("mof-%s", randomString(10)))
  os.MkdirAll(tmpFolder, 0777)
  tmpBasePath := filepath.Join(tmpFolder, strings.Replace(filepath.Base(mofFilePath), filepath.Ext(mofFilePath), "", 1))
  tmpIndexFilePath := tmpBasePath + ".f1"
  tmpDataFilePath := tmpBasePath + ".f2"
  ioutil.WriteFile(tmpIndexFilePath, indexFileData, 0777)

  mofFileFileInfo, err := os.Stat(mofFilePath)
  if err != nil {
    return err
  }
  dataFileSize := mofFileFileInfo.Size() - int64(sizeOfIndexFile) - startReadingIndex
  dataFileData := make([]byte, dataFileSize)

  _, err = mofFile.ReadAt(dataFileData, (int64(sizeOfIndexFile) + startReadingIndex))
  if err != nil {
    return err
  }
  ioutil.WriteFile(tmpDataFilePath, dataFileData, 0777)

  UndoMTF(tmpIndexFilePath, tmpDataFilePath, outDir)
  err = os.RemoveAll(tmpFolder)
  if err != nil {
    return err
  }
  return nil
}
