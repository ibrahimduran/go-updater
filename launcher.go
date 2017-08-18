package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var publicDir = "./public"
var localDataDir = "./public/data"
var remoteDataDir = "/data/"

func writeHashes(hashes map[string]string) error {
	f, err := os.Create(publicDir + "/CHECKSUM.txt")

	if err != nil {
		return err
	}

	for path, hash := range hashes {
		f.Write([]byte(path + " " + hash + "\n"))
	}

	defer f.Close()

	return nil
}

func main() {
	serve := flag.String("serve", "", "serve update data (e.g.: :8080)")
	addr := flag.String("addr", "", "connect to server (e.g.: http://localhost:8080)")
	flag.Parse()

	hashes, err := MD5Dir(localDataDir)
	if err != nil {
		log.Fatal(err)
	}

	if *serve != "" {
		fmt.Println("Serving on", *serve)
		writeHashes(hashes)
		ServeStatic(":8080", publicDir)
	} else if *addr != "" {
		fmt.Println("Checking update data using", *addr)
		outdated, err := CheckUpdates(*addr, hashes)

		if err != nil {
			log.Fatal(err)
		}

		if count := len(outdated); count > 0 {
			fmt.Printf("Found %d outdated files\n", count)

			for _, file := range outdated {
				fmt.Printf("Downloading file: %s\n", file)
				size, err := Download(*addr+remoteDataDir, file, localDataDir)

				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Download completed, size: %d bytes\n", size)
			}
		} else {
			fmt.Println("Outdated files not found")
		}
	} else {
		flag.Usage()
	}
}
