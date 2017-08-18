package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func MD5File(file string) []byte {
	f, err := os.Open(file)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return h.Sum(nil)
}

func MD5Dir(dir string) (map[string]string, error) {
	path, err := filepath.Abs(dir)

	if err != nil {
		log.Fatal(err)
	}

	hashes := make(map[string]string)

	e := filepath.Walk(path, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fi.IsDir() {
			base, err := os.Getwd()

			if err != nil {
				return err
			}

			file, err := filepath.Rel(filepath.Join(base, dir), path)

			if err != nil {
				return err
			}

			file = strings.Replace(file, string(filepath.Separator), "/", -1)

			hash := hex.EncodeToString(MD5File(path))
			hashes[file] = hash
		}

		return nil
	})

	if e != nil {
		return nil, e
	}

	return hashes, nil
}
