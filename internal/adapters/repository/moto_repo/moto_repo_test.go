package motorepo

import (
	"context"
	"flag"
	"log"
	"os"
	"testing"
	"errors"

	"github.com/vvetta/electoral_system/internal/adapters/logger"
	"github.com/vvetta/electoral_system/internal/domain"
	"github.com/vvetta/electoral_system/internal/usecase"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	db *gorm.DB
	integration = flag.Bool("integration", false, "run integration tests")
	mtRepo usecase.MotoRepo
	lg usecase.Logger
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
		mtRepo = NewMotoRepo(db, lg)
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

func TestMotoRepo_Crud(t *testing.T) {
	if !*integration {
		t.Skip("integration tests disabled")
	}

	ctx := context.Background()

	// Create test

	testMoto := domain.Moto{
		Name: "Yamaha",
		Year: 2000,
		Mileage: 25000,
		MotoType: "Быстрый",
		Location: "ВДНХ",
		EngineSize: 2,
		Price: int64(1000000),
	}

	resultMoto, err := mtRepo.Create(ctx, testMoto)
	if err != nil {
		t.Errorf("create moto error: %v", err)
	}

	result, err := mtRepo.Read(ctx, resultMoto.ID)
	if err != nil {
		t.Errorf("read moto error: %v", err)
	}

	if resultMoto.ID != result.ID {
		t.Errorf("create != read moto")
	}

	// Update test
	
	newMoto := domain.Moto{
		ID: 2,
		Name: "Yamaha",
		Year: 2002,
		Mileage: 250,
		MotoType: "Быстрый",
		Location: "ВДНХ",
		EngineSize: 2,
		Price: int64(1000000),
	}

	updated, err := mtRepo.Update(ctx, newMoto)
	if err != nil {
		t.Errorf("update moto error: %v", err)
	}

	r, err := mtRepo.Read(ctx, updated.ID)
	if err != nil {
		t.Errorf("read updated moto error: %v", err)
	}

	if r.ID != updated.ID {
		t.Errorf("updated id != read updated id")
	}

	// Delete test

	err = mtRepo.Delete(ctx, r.ID)
	if err != nil {
		t.Errorf("delete moto error: %v", err)
	}

	_, err = mtRepo.Read(ctx, r.ID)
	if !errors.Is(err, domain.RecordNotFound) {
		t.Errorf("moto dont deleted! %v", err)
	}
}

