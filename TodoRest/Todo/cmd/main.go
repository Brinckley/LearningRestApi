package main

import (
	todo "Todo"
	handler "Todo/pkg/handler"
	"Todo/pkg/repository"
	"Todo/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter)) // setting format for logrus

	if err := initConfig(); err != nil { // initing config for viper
		logrus.Fatalf("Error getting config info: %s\n", err.Error())
	}

	if err := godotenv.Load(); err != nil { // loading environment
		logrus.Fatalln("Error getting environment :", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		logrus.Fatalln("Error occurred while initializing db :", err)
	}

	// Dependency injection
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		logrus.Fatalln("Error occurred on starting server on port 8000")
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
