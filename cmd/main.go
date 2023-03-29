package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin/binding"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/companies"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/no_auth"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/service_methods"
	"gitlab.assistagro.com/back/back.auth.go/internal/api/user"
	"gitlab.assistagro.com/back/back.auth.go/internal/transport/rest"
	"gitlab.assistagro.com/back/back.auth.go/pkg/cvalidator"
	"gitlab.assistagro.com/back/back.auth.go/pkg/repository/postgres"
)

func init() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err)
	}

	postgres.InitDB(postgres.Config{
		Host:     os.Getenv("DB_IP"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME_AUTH"),
	}, viper.GetInt("db.maxOpenConnections"), viper.GetInt("db.maxIdleConnections"))
}

func main() {

	// a := struct {
	// 	FieldStr1 string
	// 	FieldInt1 int64
	// 	FieldStr2 string
	// }{
	// 	FieldStr1: "str1   ",
	// 	FieldInt1: 1,
	// 	FieldStr2: "  str2   ",
	// }
	// tools.TrimStruct(a)

	router := rest.New()

	companies.Register(router)
	service_methods.Register(router)
	no_auth.Register(router)
	user.Register(router)

	cvalidator.Register()
	binding.EnableDecoderDisallowUnknownFields = true

	srv := &http.Server{
		Addr:           "localhost:" + viper.GetString("server.port"),
		Handler:        router,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error occured while running http server: %s", err)
		}
	}()

	log.Info("Service started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Info("Shutting down the service")

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Errorf("error occured on server shutting down: %s", err)
	}

	if err := postgres.DB.Close(); err != nil {
		log.Errorf("error occured on db connection close: %s", err)
	}

	//err := router.Run("localhost:" + viper.GetString("server.port"))
	//if err != nil {
	//	panic(err)
	//}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
