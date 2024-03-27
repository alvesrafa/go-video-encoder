package services

import (
	"errors"
	"os"
	"strconv"

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
	if err != nil {
		return j.failJob(err)
	}

	err = j.ChangeJobStatus("FRAGMENTING")
	if err != nil {
		return j.failJob(err)
	}
	err = j.VideoService.Fragment()
	if err != nil {
		return j.failJob(err)
	}

	err = j.ChangeJobStatus("ENCODING")
	if err != nil {
		return j.failJob(err)
	}
	err = j.VideoService.Encode()
	if err != nil {
		return j.failJob(err)
	}

	err = j.performUpload()
	if err != nil {
		return j.failJob(err)
	}

	err = j.ChangeJobStatus("FINISHING")
	if err != nil {
		return j.failJob(err)
	}
	err = j.VideoService.Finish()
	if err != nil {
		return j.failJob(err)
	}

	err = j.ChangeJobStatus("COMPLETED")
	if err != nil {
		return j.failJob(err)
	}

	return nil
}

func (j *JobService) performUpload() error {
	err := j.ChangeJobStatus("UPLOADING")

	if err != nil {
		return j.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + j.VideoService.Video.ID
	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	uploadResult := <-doneUpload

	if uploadResult != "upload completed" {
		return j.failJob(errors.New(uploadResult))
	}

	return err
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
