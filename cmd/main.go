package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: gitx <command> [<args>...]\n")
		os.Exit(1)
	}

	switch command := os.Args[1]; command {
	case "init":
		for _, dir := range []string{".gitx", ".gitx/objects", ".gitx/refs"} {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory: %s\n", err)
			}
		}

		headFileContents := []byte("ref: refs/heads/main\n")
		if err := os.WriteFile(".gitx/HEAD", headFileContents, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing file: %s\n", err)
		}

		fmt.Println("Initialized gitx directory")
	case "cat-file":
		filePath := fmt.Sprintf(".gitx/objects/%s/%s", os.Args[3][:2], os.Args[3][2:])
		file, err := os.Open(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %s\n", err)
			os.Exit(1)
		}
		fileReader := io.Reader(file)
		r, _ := zlib.NewReader(fileReader)
		w, _ := io.ReadAll(r)
		parts := strings.Split(string(w), "\x00")
		if os.Args[2] == "-p" {
			fmt.Print(parts[1])
		}
		r.Close()
	case "hash-object":
		file, _ := os.ReadFile(os.Args[3])
		stats, _ := os.Stat(os.Args[3])
		content := string(file)
		contentHeader := fmt.Sprintf("blob %d\x00%s", stats.Size(), content)
		sha := sha1.Sum([]byte(contentHeader))
		hash := fmt.Sprintf("%x", sha)
		blobName := []rune(hash)
		blobPath := ".gitx/objects" + string(blobName[:2]) + string(blobName[:2])
		var buffer bytes.Buffer
		z := zlib.NewWriter(&buffer)
		z.Write([]byte(contentHeader))
		z.Close()
		os.Mkdir(filepath.Dir(blobPath), os.ModePerm)
		f, _ := os.Create(blobPath)
		defer f.Close()
		f.Write(buffer.Bytes())
		fmt.Print(hash)

	default:
		fmt.Fprintf(os.Stderr, "Unknown command %s\n", command)
		os.Exit(1)
	}
}
