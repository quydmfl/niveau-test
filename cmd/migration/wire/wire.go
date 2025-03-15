//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/quydmfl/niveau-test/internal/repository"
	"github.com/quydmfl/niveau-test/internal/server"
	"github.com/quydmfl/niveau-test/pkg/app"
	"github.com/quydmfl/niveau-test/pkg/log"
	"github.com/spf13/viper"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	//repository.NewRedis,
	repository.NewRepository,
	repository.NewUserRepository,
	repository.NewProductRepository,
	repository.NewSupplierRepository,
	repository.NewCategoryRepository,
	repository.NewDocumentsRepository,
)
var serverSet = wire.NewSet(
	server.NewMigrateServer,
)

// build App
func newApp(
	migrateServer *server.MigrateServer,
) *app.App {
	return app.NewApp(
		app.WithServer(migrateServer),
		app.WithName("niveau-test-migrate"),
	)
}

func NewWire(*viper.Viper, *log.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		repositorySet,
		serverSet,
		newApp,
	))
}
