//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/quydmfl/niveau-test/internal/repository"
	"github.com/quydmfl/niveau-test/internal/server"
	"github.com/quydmfl/niveau-test/internal/task"
	"github.com/quydmfl/niveau-test/pkg/app"
	"github.com/quydmfl/niveau-test/pkg/log"
	"github.com/quydmfl/niveau-test/pkg/sid"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewTransaction,
	repository.NewUserRepository,
	repository.NewProductRepository,
	repository.NewSupplierRepository,
	repository.NewCategoryRepository,
	repository.NewDocumentsRepository,
)

var taskSet = wire.NewSet(
	task.NewTask,
	task.NewUserTask,
)
var serverSet = wire.NewSet(
	server.NewTaskServer,
)

// build App
func newApp(
	task *server.TaskServer,
) *app.App {
	return app.NewApp(
		app.WithServer(task),
		app.WithName("demo-task"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		taskSet,
		serverSet,
		newApp,
		sid.NewSid,
	))
}
