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

func JobRepositoryTest(t *testing.T) {
	t.Run("should create a new job", func(t *testing.T) {

		db := database.NewDatabaseTest()
		defer db.Close()

		video := domain.NewVideo()

		video.ID = uuid.NewV4().String()
		video.FilePath = "path"
		video.CreatedAt = time.Now()
		video.ResourceID = "resource"

		repo := repositories.NewVideoRepository(db)
		repo.Insert(video)

		job, err := domain.NewJob("outputgpath", "PENDING", video)

		require.Nil(t, err)

		repoJob := repositories.JobRepositoryDb{Db: db}
		repoJob.Insert(job)

		createdJob, err := repoJob.Find(job.ID)

		require.NotEmpty(t, createdJob.ID)
		require.Nil(t, err)
		require.Equal(t, createdJob.ID, job.ID)
		require.Equal(t, createdJob.VideoID, video.ID)
	})

	t.Run("should update a job", func(t *testing.T) {

		db := database.NewDatabaseTest()
		defer db.Close()

		video := domain.NewVideo()

		video.ID = uuid.NewV4().String()
		video.FilePath = "path"
		video.CreatedAt = time.Now()
		video.ResourceID = "resource"

		repo := repositories.NewVideoRepository(db)
		repo.Insert(video)

		job, err := domain.NewJob("outputgpath", "PENDING", video)

		require.Nil(t, err)

		repoJob := repositories.JobRepositoryDb{Db: db}
		repoJob.Insert(job)

		job.Status = "COMPLETED"
		repoJob.Update(job)

		updatedJob, err := repoJob.Find(job.ID)

		require.NotEmpty(t, updatedJob.ID)
		require.Nil(t, err)
		require.Equal(t, updatedJob.Status, job.Status)
	})
}
