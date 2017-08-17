package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func CheckUpdates(addr string, hashes map[string]string) ([]string, error) {
	checksums, err := readChecksum(addr)

	if err != nil {
		return nil, err
	}

	outdated := []string{}

	for file, hash := range checksums {
		if hashes[file] != hash {
			outdated = append(outdated, file)
		}
	}

	return outdated, nil
}

func Download(addr string, file string, dir string) (int64, error) {
	resp, err := http.Get(addr + "/" + file)

	if err != nil {
		return 0, err
	}

	path := strings.Replace(file, "/", string(filepath.Separator), -1)
	err = os.MkdirAll(filepath.Join(dir, filepath.Dir(path)), 0777)

	if err != nil {
		return 0, err
	}

	f, err := os.Create(filepath.Join(dir, path))

	if err != nil {
		return 0, err
	}

	defer f.Close()

	return io.Copy(f, resp.Body)
}

func readChecksum(addr string) (map[string]string, error) {
	resp, err := http.Get(addr + "/CHECKSUM.txt")

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(body), "\n")
	hashes := make(map[string]string)

	for _, line := range lines {
		if line == "" {
			continue
		}

		words := strings.Split(line, " ")
		hashes[words[0]] = words[1]
	}

	return hashes, nil
}
