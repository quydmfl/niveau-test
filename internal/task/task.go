package task

import (
	"github.com/quydmfl/niveau-test/internal/repository"
	"github.com/quydmfl/niveau-test/pkg/jwt"
	"github.com/quydmfl/niveau-test/pkg/log"
	"github.com/quydmfl/niveau-test/pkg/sid"
)

type Task struct {
	logger *log.Logger
	sid    *sid.Sid
	jwt    *jwt.JWT
	tm     repository.Transaction
}

func NewTask(
	tm repository.Transaction,
	logger *log.Logger,
	sid *sid.Sid,
) *Task {
	return &Task{
		logger: logger,
		sid:    sid,
		tm:     tm,
	}
}
