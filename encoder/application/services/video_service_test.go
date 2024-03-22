package services_test

import (
	"fmt"
	"log"
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

	t.Run("should download a video", func(t *testing.T) {
		video, repo := prepare()
		videoService := services.NewVideoService()
		videoService.Video = video
		videoService.VideoRepository = repo

		err := videoService.Download("videoencoder")
		fmt.Println(err)
		require.Nil(t, err)
	})
}
