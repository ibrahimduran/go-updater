package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	serve := flag.String("serve", "", "serve update data (e.g.: :8080)")
	addr := flag.String("addr", "", "connect to server (e.g.: http://localhost:8080)")
	dataDir := flag.String("data", "data", "set data directory to download updates to")
	secretKey := flag.String("secret", "", "set secret key for authentication")

	flag.Parse()

	if *serve == "" && *addr == "" {
		flag.Usage()
	} else {
		err := os.MkdirAll("./"+*dataDir, 0777)

		if err != nil {
			log.Fatal(err)
		}
	}

	hashes, err := MD5Dir("./" + *dataDir)
	if err != nil {
		log.Fatal(err)
	}

	if *serve != "" {
		fmt.Println("Serving", *dataDir, " on", *serve)
		err := ServeStatic(":8080", *secretKey, &hashes, "./"+*dataDir)

		if err != nil {
			log.Fatal(err)
		}
	} else if *addr != "" {
		fmt.Println("Checking update data using", *addr)
		outdated, err := CheckUpdates(*addr, *secretKey, hashes)

		if err != nil {
			log.Fatal(err)
		}

		if count := len(outdated); count > 0 {
			fmt.Printf("Found %d outdated files\n", count)

			for _, file := range outdated {
				fmt.Printf("Downloading file: %s\n", file)
				size, err := Download(*addr, *secretKey, file, "./"+*dataDir)

				if err != nil {
					log.Fatal(err)
				}

				fmt.Printf("Download completed, size: %d bytes\n", size)
			}
		} else {
			fmt.Println("Outdated files not found")
		}
	}
}
