package application

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/eliezerraj/go-core/v3/logger"
	"github.com/eliezerraj/go-core/v3/database/connector"
	"github.com/eliezerraj/go-core/v3/httpclient"

	"github.com/go-semantic-engine-v2/application/infrastructure/repository"
	"github.com/go-semantic-engine-v2/application/config"
	"github.com/go-semantic-engine-v2/application/controller"
	"github.com/go-semantic-engine-v2/application/domain/usecase"
	"github.com/go-semantic-engine-v2/application/infrastructure/module"
)

type Application struct {
	SemanticController *controller.SemanticController
}

type UseCase struct {
	SemanticUseCase usecase.ISemanticUseCase
}

type Repository struct {
	SemanticRepository repository.ISemanticRepository
}

func NewApplication(cfg *config.Config) (*Application, error) {
	logger.InfoOutCtx("initializing application SUCCESSFULLY")

	// Initialize database connector
	readerConfig := connector.ConnectorConfig{
		DSN:              "postgres://" + cfg.Database.Username + ":" + cfg.Database.Password + "@" + cfg.Database.Host + ":" + cfg.Database.Port + "/" + cfg.Database.Name,
		MaxConnIdleTime:  cfg.Database.ConnIdleTime * time.Minute,
		MaxConnLifeTime:  cfg.Database.ConnLifetime * time.Minute,
		MaxConns:         cfg.Database.MaxConnections,
		MinConns:         cfg.Database.MinConnections,
		DBConnTimeout:    cfg.Database.ConnTimeout,
		HealthCheckPeriod: cfg.Database.ConnIdleTime * time.Minute / 2,
	}

	writerConfig := connector.ConnectorConfig{
		DSN:              "postgres://" + cfg.Database.Username + ":" + cfg.Database.Password + "@" + cfg.Database.Host + ":" + cfg.Database.Port + "/" + cfg.Database.Name,
		MaxConnIdleTime:  cfg.Database.ConnIdleTime * time.Minute,
		MaxConnLifeTime:  cfg.Database.ConnLifetime * time.Minute,
		MaxConns:         cfg.Database.MaxConnections,
		MinConns:         cfg.Database.MinConnections,
		DBConnTimeout:    cfg.Database.ConnTimeout,
		HealthCheckPeriod: cfg.Database.ConnIdleTime * time.Minute / 2,
	}

	logger.InfoOutCtx("readerConfig initialized SUCCESSFULLY", zap.Any("readerConfig", readerConfig), zap.Any("writerConfig", writerConfig))

	dbConnector, err := connector.NewDatabaseConnector(cfg.App.Name, readerConfig, writerConfig)
	if err != nil {
		logger.FatalOutCtx("failed to initialize database connector")
		return nil, err
	}
	
	logger.InfoOutCtx("dbConnector initialized SUCCESSFULLY", zap.Any("dbConnector", dbConnector))
	
	pgConnection := &connector.PgConnection{}
	_, err = pgConnection.NewPool(context.Background(), readerConfig)
	if err != nil {
		logger.FatalOutCtx("failed to create database pool")
		return nil, err
	}
	err = pgConnection.Ping(context.Background())
	if err != nil {
		logger.FatalOutCtx("failed to ping pg connection")
		return nil, err
	}
	
	// Repository initialization (where the RSA keys are loaded and managed)
	semanticRepository := repository.NewSemanticRepository(dbConnector)

	// Create the forwards modules.
	httpConfig := &httpclient.HttpConfig{
		Timeout:             cfg.HTTP.Timeout * time.Second,
		KeepAlive:           cfg.HTTP.KeepAlive * time.Second,
		IdleConnTimeout:     cfg.HTTP.IdleConnTimeout * time.Second,
		MaxIdleConns:        cfg.HTTP.MaxIdleConns,
		MaxIdleConnsPerHost: cfg.HTTP.MaxIdleConnsPerHost,
		MaxConnsPerHost:     cfg.HTTP.MaxConnsPerHost,
		ServiceName:         cfg.App.Name,
	}

	modelHttpClient := httpclient.NewHttpClient(httpConfig)
	modelModule := module.NewModelModule(cfg, modelHttpClient, nil)

	// UseCase initialization	
	semanticUsecase := usecase.NewSemanticUseCase(semanticRepository, modelModule)

	// Controller initialization
	semanticController := controller.NewSemanticController(semanticUsecase)

	return &Application{
		SemanticController: semanticController,
	}, nil
}
