package main

import (
	"flag"
	"github.com/sebastianm/npmirror/server"
	"github.com/sebastianm/npmirror/storage"
	"log"
	"os"
)

func assert(mustTrue bool, errorMsg string, msgParams ...interface{}) {
	if !mustTrue {
		if len(msgParams) > 0 {
			log.Fatalf(errorMsg, msgParams)
		} else {
			log.Fatal(errorMsg)
		}

		os.Exit(1)
	}
}

func main() {
	storageType := flag.String("storage-type", "file", "Storage Type (only 'file' available in this version)")
	fileDir := flag.String("storage-file-dir", "./npmirror", "Storage directory for all cached NPM files")
	httpAddr := flag.String("http-addr", "127.0.0.1:8023", "httpd bind address")
	externalAddr := flag.String("external-addr", "", "Required! External HTTP address (registry address)")
	flag.Parse()

	// check flag values first
	assert(*storageType == "file", "Storage type '%s' not supported", *storageType)
	assert(*fileDir != "", "Please provide a valid file directory path")
	assert(*externalAddr != "", "Please provide an external http address (external-addr option)")

	s, err := storage.NewFileStorage(*fileDir)
	if err != nil {
		log.Fatalf("error creating fileStorage: %s", err.Error())
		os.Exit(1)
	}

	config := &server.StorageServerConfig{
		Storage:          s,
		HttpAddr:         *httpAddr,
		ExternalHttpAddr: *externalAddr,
	}
	server := server.NewStorageServer(config)
	log.Fatal(server.Start())
}
