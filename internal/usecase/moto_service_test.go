package usecase_test

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"

	"github.com/vvetta/electoral_system/internal/adapters/moto_parser"
	"github.com/vvetta/electoral_system/internal/adapters/repository/moto_repo"
	"github.com/vvetta/electoral_system/internal/adapters/logger"
	"github.com/vvetta/electoral_system/internal/usecase"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
	integration = flag.Bool("integration", false, "run integration tests")
	mtRepo usecase.MotoRepo
	mtParser usecase.MotoParser
	lg usecase.Logger
	mtSVC usecase.MotoService
)

func TestMain(m *testing.M) {
	flag.Parse()

	_ = godotenv.Load(".env")

	if *integration {
		var err error

		testDSN := getTestDSN()
		db, err = gorm.Open(postgres.Open(testDSN), &gorm.Config{})
		if err != nil {
			log.Fatalf("error connetc to test db: %s", testDSN)			
		}	

		lg = logger.NewLogger()
		mtRepo = motorepo.NewMotoRepo(db, lg)
		mtParser = motoparser.NewMotoParser("https://mr-moto.ru/catalog/mototsikly/", "page-card__col", 100)

		mtSVC = usecase.NewMotoService(lg, mtRepo, mtParser)		
	}

	code := m.Run()
	os.Exit(code)
}

func getTestDSN() string {
	var dsn string

	DB_USER := os.Getenv("PG_TEST_USER")
	DB_PASS := os.Getenv("PG_TEST_PASSWORD")
	DB_HOST := os.Getenv("PG_TEST_HOST")
	DB_PORT := os.Getenv("PG_TEST_PORT")
	DB_NAME := os.Getenv("PG_TEST_DB_NAME")

	dsn = "postgres://" + DB_USER + ":" + DB_PASS + "@" + DB_HOST + ":" + DB_PORT + "/" + DB_NAME + "?sslmode=disable"	

	return dsn
}

func TestMotoService_ParseAndUpdateAllMoto(t *testing.T) {
	if !*integration {
		t.Skip("integration tests disabled")
	}

	ctx := context.Background()

	motos, err := mtSVC.ParseAndUpdateAllMoto(ctx)
	if err != nil {
		t.Errorf("parse and update moto error: %v", err)
	}

	if len(motos) <= 1 {
		t.Errorf("parser motos error. len < 1")
	}

	log.Print(motos)
}

