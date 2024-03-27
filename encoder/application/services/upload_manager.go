package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (vu *VideoUpload) UploadObject(objectPath string, client *storage.Client, ctx context.Context) error {
	path := strings.Split(objectPath, os.Getenv("localStoragePath")+"/")
	f, err := os.Open(objectPath)

	if err != nil {
		return nil
	}

	defer f.Close()
	fmt.Println("vu.OutputBucket" + vu.OutputBucket)
	wc := client.Bucket(vu.OutputBucket).Object(path[1]).NewWriter(ctx)
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, err := io.Copy(wc, f); err != nil {
		return err
	}

	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}

func (vu *VideoUpload) loadPaths() error {
	err := filepath.Walk(vu.VideoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			vu.Paths = append(vu.Paths, path)
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func getClientUpload() (*storage.Client, context.Context, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return nil, nil, err
	}

	return client, ctx, nil
}

func (vu *VideoUpload) ProcessUpload(concurrency int, doneUpload chan<- string) error {
	in := make(chan int, runtime.NumCPU()) // Qual arquivo baseado na posicao do slice Paths

	returnChannel := make(chan string) // Motivos de parada, sendo erro ou conclusao daa leitura

	err := vu.loadPaths()

	fmt.Println(vu.Paths)

	if err != nil {
		return err
	}

	uploadClient, ctx, err := getClientUpload()

	if err != nil {
		return err
	}

	for process := 0; process < concurrency; process++ {

		go vu.uploadWorker(in, returnChannel, uploadClient, ctx)
	}

	go func() {
		for x := 0; x < len(vu.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	for message := range returnChannel {
		if message != "" {
			doneUpload <- message
			break
		}
	}
	fmt.Println("Donneeee")
	return nil
}

func (vu *VideoUpload) uploadWorker(in <-chan int, returnChan chan<- string, uploadClient *storage.Client, ctx context.Context) {
	for position := range in {

		err := vu.UploadObject(vu.Paths[position], uploadClient, ctx)

		if err != nil {
			vu.Errors = append(vu.Errors, vu.Paths[position])
			log.Printf("Error during the upload: %v. Error: %v", vu.Paths[position], err)
			returnChan <- err.Error()
		}

		returnChan <- ""
	}
	returnChan <- "upload completed2"
}
