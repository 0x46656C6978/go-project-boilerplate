package dicontainer

import (
	"fmt"

	"github.com/0x46656C6978/go-project-boilerplate/internal/config"
	"github.com/0x46656C6978/go-project-boilerplate/internal/httpapi"
	"github.com/0x46656C6978/go-project-boilerplate/internal/migrator"
	"github.com/0x46656C6978/go-project-boilerplate/internal/repository"
	"github.com/0x46656C6978/go-project-boilerplate/internal/service"
	"go.uber.org/dig"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func ProvideCore(c *dig.Container) {
	c.Provide(func() *config.Config {
		return config.New()
	})
	c.Provide(func() (*zap.Logger, error) {
		return zap.NewProduction()
	})
}

func ProvideDB(c *dig.Container) {
	c.Provide(func(cfg *config.Config) (*gorm.DB, error) {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Etc/GMT",
			cfg.DB.Host,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.DBName,
			cfg.DB.Port,
		)

		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  dsn,
			PreferSimpleProtocol: true,
		}), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		})
		if err != nil {
			return nil, err
		}
		return db, nil
	})

	c.Provide(func(db *gorm.DB) (*migrator.Migrator, error) {
		return migrator.New(db)
	})
}

func ProvideRepositories(c *dig.Container) {
	c.Provide(func(db *gorm.DB) repository.UserRepoInterface {
		return repository.NewUserRepo(db)
	})
}

func ProvideServices(c *dig.Container) {
	c.Provide(func(repo repository.UserRepoInterface) service.UserServiceInterface {
		return service.NewUserService(repo)
	})
}

func ProvideHttpApis(c *dig.Container) {
	c.Provide(func(cfg *config.Config, srv service.UserServiceInterface) *httpapi.UserHttpApi {
		return httpapi.NewUserHttpApi(cfg, srv)
	})
}
