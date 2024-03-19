package domain_test

import (
	"testing"
	"time"

	"github.com/alvesrafa/video-encoder/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideo(t *testing.T) {
	t.Run("should create a video", func(t *testing.T) {
		video := domain.NewVideo()

		video.ID = uuid.NewV4().String()
		video.ResourceID = "a"
		video.FilePath = "path"
		video.CreatedAt = time.Now()

		err := video.Validate()

		require.Nil(t, err)
	})
	t.Run("should return an error when video is empty", func(t *testing.T) {
		video := domain.NewVideo()

		err := video.Validate()

		require.Error(t, err)
	})
	t.Run("should return an error when id is not an uuid", func(t *testing.T) {
		video := domain.NewVideo()

		video.ID = "not an uuid"
		video.ResourceID = "a"
		video.FilePath = "path"
		video.CreatedAt = time.Now()

		err := video.Validate()

		require.Error(t, err)
	})

}
