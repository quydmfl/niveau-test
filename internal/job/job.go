package job

import (
	"github.com/quydmfl/niveau-test/internal/repository"
	"github.com/quydmfl/niveau-test/pkg/jwt"
	"github.com/quydmfl/niveau-test/pkg/log"
	"github.com/quydmfl/niveau-test/pkg/sid"
)

type Job struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewJob(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
) *Job {
	return &Job{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
