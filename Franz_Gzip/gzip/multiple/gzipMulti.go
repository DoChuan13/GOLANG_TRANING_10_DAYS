package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
)

var path = "gzip/multiple/test.tar.gz"

func main() {
	//writeFiles()
	readFiles()
}

func writeFiles() {
	files := []string{"gzip/multiple/file1.txt", "gzip/multiple/file2.txt", "gzip/multiple/file3.txt"}
	//files := []string{"gzip/multiple/file1.txt"}

	//output, err := os.Create(path)
	output, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer output.Close()

	gzipWriter := gzip.NewWriter(output)
	defer gzipWriter.Close()

	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()

	for _, file := range files {
		inputFile, err := os.Open(file)
		if err != nil {
			panic(err)
		}
		defer inputFile.Close()

		fileInfo, err := inputFile.Stat()
		if err != nil {
			panic(err)
		}

		header := &tar.Header{
			Name: file,
			Size: fileInfo.Size(),
		}

		if err := tarWriter.WriteHeader(header); err != nil {
			panic(err)
		}

		if _, err := io.Copy(tarWriter, inputFile); err != nil {
			panic(err)
		}
	}
}

func readFiles() {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	gzipReader, err := gzip.NewReader(file)
	if err != nil {
		panic(err)
	}
	defer gzipReader.Close()

	tarReader := tar.NewReader(gzipReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		fmt.Println(header.Name)

		if header.Typeflag == tar.TypeDir {
			continue
		}

		buf := make([]byte, header.Size)
		_, err = io.ReadFull(tarReader, buf)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(buf))
	}
}
