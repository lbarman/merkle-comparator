package main

import (
	"fmt"
	"os"
	"math"
	"flag"
	"io"
	"encoding/hex"
	"crypto/md5"
    "path/filepath"
)

const FILECHUNK = 8192    // we settle for 8KB
var out string = ""

func main() {
	flag.Parse()
    root := flag.Arg(0)
	build(root)
}

func build(folder string){
	fmt.Println("Hashes for folder", folder)
    _ = filepath.Walk(folder, visit)
    fmt.Println(out)
}

func visit(path string, f os.FileInfo, err error) error {
  out = out + path +" -> "+hashFile(path)+"\n"
  return nil
} 

func hashFile(filePath string) string {

    file, err := os.Open(filePath)
    if err != nil {
    	fmt.Println("[warn]", err)
        return "NaN"
    }
    defer file.Close()

    info, err2 := file.Stat()
    if err2 != nil {
    	fmt.Println("[warn]", err)
        return "NaN"
    }

    filesize := info.Size()  
    blocks := uint64(math.Ceil(float64(filesize) / float64(FILECHUNK)))
    hash := md5.New()

    for i := uint64(0); i < blocks; i++ {
        blocksize := int(math.Min(FILECHUNK, float64(filesize-int64(i*FILECHUNK))))
        buf := make([] byte, blocksize)

        file.Read(buf)
        io.WriteString(hash, string(buf))   // append into the hash
    }

    finalHash := hash.Sum(nil)
    return hex.EncodeToString(finalHash[:])
 }