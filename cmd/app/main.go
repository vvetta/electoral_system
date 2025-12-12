package main

import (
	"log"
	"os"
	"net/http"


	"github.com/vvetta/electoral_system/internal/adapters/http"
	"github.com/vvetta/electoral_system/internal/adapters/logger"
	motoparser "github.com/vvetta/electoral_system/internal/adapters/moto_parser"
	"github.com/vvetta/electoral_system/internal/adapters/repository/moto_repo"
	"github.com/vvetta/electoral_system/internal/usecase"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	url = "https://mr-moto.ru/catalog/mototsikly/"
)

func main() {

	_ = godotenv.Load(".env")

	dsn := getDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	lg := logger.NewLogger()

	motoRepo := motorepo.NewMotoRepo(db, lg)
	motoParser := motoparser.NewMotoParser(url, "page-card__col", 100)

	motoSVC := usecase.NewMotoService(lg, motoRepo, motoParser)
	
	srv := httpserver.NewServer(motoSVC, lg)
	if err := http.ListenAndServe(":8080", srv); err != nil {
		log.Fatal(err)
	}
}

func getDSN() string {
	DB_USER := os.Getenv("DB_USER")
	DB_PASS := os.Getenv("DB_PASS")
	DB_HOST := os.Getenv("DB_HOST")
	DB_PORT := os.Getenv("DB_PORT")
	DB_NAME := os.Getenv("DB_NAME")

	dsn := "postgres://" + DB_USER + ":" + DB_PASS + "@" + DB_HOST + ":" + DB_PORT + "/" + DB_NAME + "?sslmode=disable"
	return dsn
}
