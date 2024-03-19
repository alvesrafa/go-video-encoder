package domain_test

import (
	"testing"
	"time"

	"github.com/alvesrafa/video-encoder/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestJob(t *testing.T) {
	t.Run("should create a job", func(t *testing.T) {
		video := domain.NewVideo()
		video.ID = uuid.NewV4().String()
		video.ResourceID = uuid.NewV4().String()
		video.FilePath = "path"
		video.CreatedAt = time.Now()

		job, err := domain.NewJob("path", "CONVERTED", video)

		require.Nil(t, err)
		require.NotNil(t, job)

	})

}
