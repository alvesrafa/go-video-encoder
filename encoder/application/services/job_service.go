package services

import (
	"os"

	"github.com/alvesrafa/video-encoder/application/repositories"
	"github.com/alvesrafa/video-encoder/domain"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (j *JobService) Start() error {
	err := j.ChangeJobStatus("DOWNLOADING")

	if err != nil {
		return j.failJob(err)
	}
	err = j.VideoService.Download(os.Getenv("inputBucketName"))

	return nil
}
func (j *JobService) ChangeJobStatus(status string) error {
	var err error
	j.Job.Status = status
	j.Job, err = j.JobRepository.Update(j.Job)

	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) failJob(error error) error {

	j.Job.Status = "FAILED"
	j.Job.Error = error.Error()

	_, err := j.JobRepository.Update(j.Job)

	if err != nil {
		return err
	}

	return error
}
