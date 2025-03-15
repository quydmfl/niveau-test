//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/quydmfl/niveau-test/internal/handler"
	"github.com/quydmfl/niveau-test/internal/job"
	"github.com/quydmfl/niveau-test/internal/repository"
	"github.com/quydmfl/niveau-test/internal/server"
	"github.com/quydmfl/niveau-test/internal/service"
	"github.com/quydmfl/niveau-test/pkg/app"
	"github.com/quydmfl/niveau-test/pkg/jwt"
	"github.com/quydmfl/niveau-test/pkg/log"
	"github.com/quydmfl/niveau-test/pkg/server/http"
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

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewProductService,
	service.NewSupplierService,
	service.NewCategoryService,
)

var handlerSet = wire.NewSet(
	handler.NewHandler,
	handler.NewUserHandler,
	handler.NewProductHandler,
	handler.NewSupplierHandler,
	handler.NewCategoryHandler,
)

var jobSet = wire.NewSet(
	job.NewJob,
	job.NewUserJob,
)
var serverSet = wire.NewSet(
	server.NewHTTPServer,
	server.NewJobServer,
)

// build App
func newApp(
	httpServer *http.Server,
	jobServer *server.JobServer,
	// task *server.Task,
) *app.App {
	return app.NewApp(
		app.WithServer(httpServer, jobServer),
		app.WithName("demo-server"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serviceSet,
		handlerSet,
		jobSet,
		serverSet,
		sid.NewSid,
		jwt.NewJwt,
		newApp,
	))
}
