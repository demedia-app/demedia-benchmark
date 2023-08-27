package main

import (
	"flag"
	"log"
	"os"

	"github.com/sithumonline/demedia-benchmark/util"

	"github.com/sithumonline/demedia-nostr/ipfs"
)

func main() {
	i := ipfs.NewIPFSClient(
		util.EnvOrDefault("IPFS_NODE", "https://ipfs.infura.io:5001"),
		util.EnvOrDefault("INFURA_PROJECT_ID", "2RnWuMLJyQznOY1k0VwGr5vPubC"),
		util.EnvOrDefault("INFURA_PROJECT_SECRET", "57d68d08848b58e2f76804a407fd1c97"),
	)

	var file_path string
	flag.StringVar(&file_path, "file", "", "file path")
	flag.Parse()

	if file_path == "" {
		log.Fatal("file path is required")
	}

	log.Printf("file path: %s", file_path)

	data, err := os.ReadFile(file_path)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	link, err := i.UploadFile(data)
	if err != nil {
		log.Fatalf("failed to upload file: %v", err)
	}

	log.Printf("link: %s", link)
}

/*
*	go run ipfs/main.go -file=/Users/sithumsandeepa/Downloads/348662e4dde23114.jpg
*/
