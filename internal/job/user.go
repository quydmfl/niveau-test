package job

import (
	"context"
	"time"

	"github.com/quydmfl/niveau-test/internal/repository"
)

type UserJob interface {
	KafkaConsumer(ctx context.Context) error
}

func NewUserJob(
	job *Job,
	userRepo repository.UserRepository,
) UserJob {
	return &userJob{
		userRepo: userRepo,
		Job:      job,
	}
}

type userJob struct {
	userRepo repository.UserRepository
	*Job
}

func (t userJob) KafkaConsumer(ctx context.Context) error {
	// do something
	for {
		t.logger.Info("KafkaConsumer")
		time.Sleep(time.Second * 5)
	}
}
