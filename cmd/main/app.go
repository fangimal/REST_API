package main

import (
	author2 "REST_API/internal/author"
	"REST_API/internal/author/db"
	"REST_API/internal/config"
	"REST_API/internal/user"
	"REST_API/pkg/client/postgresql"
	"REST_API/pkg/logging"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/julienschmidt/httprouter"
)

func main() {
	logger := logging.GetLogger()
	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	postgreSQLClient, err := postgresql.NewClient(context.TODO(), 3, cfg.Storage)
	if err != nil {
		logger.Fatal("%v", err)
	}
	repository := author.NewRepository(postgreSQLClient, logger)

	//one, err := repository.FindOne(context.TODO(), "e459ef69-c63a-48ad-8b77-b089ffa43c58")
	//if err != nil {
	//	return
	//}
	//logger.Info(one)

	//newAth := author2.Author{
	//	Name: "MIR",
	//}
	//err = repository.Create(context.TODO(), &newAth)
	//if err != nil {
	//	logger.Fatalf("%v", err)
	//}
	//logger.Infof("%v", newAth)

	//all, err := repository.FindAll(context.TODO())
	//if err != nil {
	//	logger.Fatalf("%v", err)
	//}
	//
	//for _, ath := range all {
	//	logger.Infof("%v", ath)
	//}

	//cfgMongo := cfg.MongoDB
	//mongoDBClient, err := mongodb.NewClient(context.Background(), cfgMongo.Host, cfgMongo.Port, cfgMongo.Username,
	//	cfgMongo.Password, cfgMongo.Database, cfgMongo.AuthDB)
	//if err != nil {
	//	panic(err)
	//}
	//storage := db.NewStorage(mongoDBClient, cfg.MongoDB.Collection, logger)

	//TODO for test
	//users, err := storage.FindAll(context.Background())
	//fmt.Println(users)

	logger.Info("register author handler")
	authorHandler := author2.NewHandler(repository, logger)
	authorHandler.Register(router)

	logger.Info("register user handler")
	handler := user.NewHandler(logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen unix tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		panic(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal((server.Serve(listener)))
}
