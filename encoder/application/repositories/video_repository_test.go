package repositories_test

import (
	"testing"
	"time"

	"github.com/alvesrafa/video-encoder/application/repositories"
	"github.com/alvesrafa/video-encoder/domain"
	"github.com/alvesrafa/video-encoder/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func VideoRepositoryTest(t *testing.T) {
	t.Run("should insert a new video", func(t *testing.T) {
		db := database.NewDatabaseTest()
		defer db.Close()

		video := domain.NewVideo()

		video.ID = uuid.NewV4().String()
		video.FilePath = "path"
		video.CreatedAt = time.Now()
		video.ResourceID = "resource"

		repo := repositories.NewVideoRepository(db)
		repo.Insert(video)

		createdVideo, err := repo.Find(video.ID)

		require.NotEmpty(t, createdVideo.ID)
		require.Nil(t, err)
		require.Equal(t, createdVideo.ID, video.ID)
	})
}
