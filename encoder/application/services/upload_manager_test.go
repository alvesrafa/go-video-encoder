package services_test

import (
	"log"
	"os"
	"testing"

	"github.com/alvesrafa/video-encoder/application/services"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func TestUploadManager(t *testing.T) {

	t.Run("should upload a video", func(t *testing.T) {
		video, repo := prepare()
		videoService := services.NewVideoService()
		videoService.Video = video
		videoService.VideoRepository = repo

		err := videoService.Download("video-encoder-alvesrafa")
		require.Nil(t, err)

		err = videoService.Fragment()
		require.Nil(t, err)

		err = videoService.Encode()
		require.Nil(t, err)

		videoUpload := services.NewVideoUpload()
		videoUpload.OutputBucket = "video-encoder-alvesrafa"
		videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + video.ID

		doneUpload := make(chan string)

		go videoUpload.ProcessUpload(50, doneUpload)

		result := <-doneUpload
		require.Equal(t, result, "upload completed")

		err = videoService.Finish()
		require.Nil(t, err)
	})
}
