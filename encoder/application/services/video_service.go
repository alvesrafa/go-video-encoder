package services

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/alvesrafa/video-encoder/application/repositories"
	"github.com/alvesrafa/video-encoder/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (v *VideoService) Download(bucketName string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	fmt.Printf("client2222 %v\n", client)
	fmt.Printf("err %v\n", err)
	if err != nil {
		return err
	}

	defer client.Close()

	bkt := client.Bucket(bucketName)

	obj := bkt.Object(v.Video.FilePath)

	r, err := obj.NewReader(ctx)

	if err != nil {
		return err
	}

	defer r.Close()

	body, err := io.ReadAll(r)

	if err != nil {
		return err
	}
	f, err := os.Create(os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4")

	if err != nil {
		return err
	}

	_, err = f.Write(body)

	if err != nil {
		return err
	}

	fmt.Printf("Video %v has been stored", v.Video.ID)

	return nil
}

func (v *VideoService) Fragment() error {

	err := os.Mkdir(os.Getenv("localStoragePath")+"/"+v.Video.ID, os.ModePerm)

	if err != nil {
		return err
	}

	source := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".mp4"
	target := os.Getenv("localStoragePath") + "/" + v.Video.ID + ".frag"

	cmd := exec.Command("mp4fragment", source, target)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return err
	}
	printOutput(output)

	return nil
}

func printOutput(output []byte) {
	if len(output) > 0 {
		fmt.Printf("=========> Output: %s\n", string(output))
	}
}
