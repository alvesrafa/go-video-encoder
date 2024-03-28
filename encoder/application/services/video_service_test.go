package services_test

import (
	"log"
	"os"
	"testing"
	"time"

	"github.com/alvesrafa/video-encoder/application/repositories"
	"github.com/alvesrafa/video-encoder/application/services"
	"github.com/alvesrafa/video-encoder/domain"
	"github.com/alvesrafa/video-encoder/framework/database"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDatabaseTest()
	defer db.Close()

	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.FilePath = "dluffy.mp4"
	video.CreatedAt = time.Now()
	video.ResourceID = "resource"

	repo := repositories.VideoRepositoryDb{Db: db}

	return video, repo
}
func TestVideoService(t *testing.T) {
	video, repo := prepare()
	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	t.Run("should download a video", func(t *testing.T) {
		err := videoService.Download(os.Getenv("inputBucketName"))
		require.Nil(t, err)
	})
	t.Run("should fragment a video and fragment", func(t *testing.T) {
		err := videoService.Fragment()
		require.Nil(t, err)
	})
	t.Run("should download, fragment and encode", func(t *testing.T) {
		err := videoService.Encode()
		require.Nil(t, err)
	})

	t.Run("should download, fragment, encode and remove temp a video", func(t *testing.T) {
		video, repo := prepare()
		videoService := services.NewVideoService()
		videoService.Video = video
		videoService.VideoRepository = repo

		err := videoService.Download(os.Getenv("inputBucketName"))
		require.Nil(t, err)

		err = videoService.Fragment()
		require.Nil(t, err)

		err = videoService.Encode()
		require.Nil(t, err)

		err = videoService.Finish()
		require.Nil(t, err)
	})
}
