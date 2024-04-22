package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/alvesrafa/video-encoder/domain"
	"github.com/alvesrafa/video-encoder/framework/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

// pega a msg do rabbitmq e executa o sstart.

type JobWorkerResult struct {
	Job     domain.Job
	Message *amqp.Delivery
	Error   error
}

func JobWorker(messageChannel <-chan amqp.Delivery, returnChannel chan<- JobWorkerResult, jobService JobService, workerID int, job domain.Job) {
	// Formato do JSON
	// {
	// 	resource_id: "id do video da pessoa que enviou para a fila.",
	// 	file_path: "dluffy.mp4"
	// }

	for message := range messageChannel {
		// pega message body do json
		// validar se o json é um json valido
		// validar o video e inserir o video no banco de dados
		// startar
		err := utils.IsJson(string(message.Body))
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = json.Unmarshal(message.Body, &jobService.VideoService.Video) // Injetando valores do body diretamente no video do videservice
		jobService.VideoService.Video.ID = uuid.NewV4().String()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.Video.Validate()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		err = jobService.VideoService.InsertVideo()
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		job.Video = jobService.VideoService.Video
		job.OutputBucketPath = os.Getenv("outputBucketName")
		job.ID = uuid.NewV4().String()
		job.Status = "STARTING"
		job.CreatedAt = time.Now()

		_, err = jobService.JobRepository.Insert(&job)
		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}
		fmt.Println("=========> =========> &job")
		fmt.Println(&job)

		jobService.Job = &job
		err = jobService.Start()

		if err != nil {
			returnChannel <- returnJobResult(domain.Job{}, message, err)
			continue
		}

		returnChannel <- returnJobResult(job, message, nil)
	}
}

func returnJobResult(job domain.Job, message amqp.Delivery, err error) JobWorkerResult {
	result := JobWorkerResult{
		Job:     job,
		Message: &message,
		Error:   err,
	}

	return result
}
