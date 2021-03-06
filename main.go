package main

import (
	"context"
	"fmt"
	"log"
	"os"

	common "github.com/open-Q/common/golang"
	commonLog "github.com/open-Q/common/golang/log"
	proto "github.com/open-Q/common/golang/proto/user"
	commonService "github.com/open-Q/common/golang/service"
	"github.com/open-Q/user/controller"
	"github.com/open-Q/user/storage"
)

const (
	envMongoConn = "mongo:conn"
	envMongoDB   = "mongo:db"
)

// This variable is assigned during build time using build flags.
var version string

func main() {
	ctx := context.Background()

	// initialize logger.
	logger, err := commonLog.NewFileLogger("./log", fmt.Sprintf("log_%s.json", version), os.ModePerm)
	if err != nil {
		log.Fatalf("could not initialize logger: %v", err)
	}

	// setup interrupt hook.
	go common.InterruptHook(func() {
		logger.Error("interrupted")
		os.Exit(1)
	})

	// create service instance.
	microService, flagsMap, err := commonService.New("./.contract/contract.json")
	if err != nil {
		logger.Fatalf("could not create service: %v", err)
	}

	// initialize storage.
	mongoConnString := flagsMap[envMongoConn]
	mongoDBName := flagsMap[envMongoDB]
	userStorage, err := storage.NewMongoStorage(ctx, mongoConnString.Value().(string), mongoDBName.Value().(string))
	if err != nil {
		logger.Fatalf("could not create connection to storage: %v", err)
	}
	defer func() {
		if err := userStorage.Disconnect(ctx); err != nil {
			logger.Errorf("could not close storage connection: %v", err)
		}
	}()

	// register service controller.
	service := controller.New(controller.Config{
		Logger:      logger,
		UserStorage: userStorage,
	})
	if err := proto.RegisterUserHandler(microService.Server(), service); err != nil {
		logger.Fatalf("could not register service controller: %v", err)
	}

	// run service.
	logger.Infof("service started, version: %s", version)
	if err := microService.Run(); err != nil {
		logger.Fatalf("could not run service: %v", err)
	}
	logger.Info("service stopped")
}
