package motorepo

import (
	"fmt"
	"context"
	"errors"

	"github.com/vvetta/electoral_system/internal/domain"
	"github.com/vvetta/electoral_system/internal/usecase"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type motoRepo struct {
	db *gorm.DB
	log usecase.Logger
}

func NewMotoRepo(db *gorm.DB, log usecase.Logger) usecase.MotoRepo {
	return &motoRepo{
		db: db,
		log: log,
	}
}

func (r *motoRepo) Create(ctx context.Context, moto domain.Moto) (domain.Moto, error) {
	r.log.Debug("MotoRepo_Create: Start!")	

	var gormMoto GormMoto
	gormMoto = toGormMoto(moto)

	result := r.db.WithContext(ctx).Create(&gormMoto)
	if result.Error != nil {
		r.log.Error("MotoRepo_Create: create new moto error!", "err", result.Error)
		return domain.Moto{}, fmt.Errorf("%w: create moto error: %v", domain.InternalError, result.Error)
	}

	if result.RowsAffected != 1 {
		r.log.Debug("MotoRepo_Create: moto already exists!", "id", gormMoto.ID)
		return domain.Moto{}, fmt.Errorf("%w: moto already exists!", domain.RecordAlreadyExists)
	}
	r.log.Debug("MotoRepo_Create: create new moto success!", "id", gormMoto.ID)

	domainMoto := toDomainMoto(gormMoto)
	
	r.log.Debug("MotoRepo_Create: End!")
	return domainMoto, nil
}

func (r *motoRepo) Read(ctx context.Context, motoID uint) (domain.Moto, error) {
	r.log.Debug("MotoRepo_Read: Start!")	

	var gormMoto GormMoto
	result := r.db.WithContext(ctx).Where("id = ?", motoID).First(&gormMoto)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.log.Debug("MotoRepo_Read: record not found", "id", motoID)
			return domain.Moto{}, domain.RecordNotFound
		}
		r.log.Error("MotoRepo_Read: internal error", "id", motoID, "err", result.Error)
		return domain.Moto{}, domain.InternalError
	}
	r.log.Debug("MotoRepo_Read: read moto from db success!", "id", gormMoto.ID)

	domainMoto := toDomainMoto(gormMoto)
	
	r.log.Debug("MotoRepo_Read: End!")	
	return domainMoto, nil
}

func (r *motoRepo) Update(ctx context.Context, moto domain.Moto) (domain.Moto, error) {
	r.log.Debug("MotoRepo_Update: Start!")	

	var gormMoto GormMoto
	gormMoto = toGormMoto(moto)

	result := r.db.WithContext(ctx).Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"year",
				"name",
				"mileage",
				"engine_size",
				"moto_type",
				"location",
				"price",
				"updated_at",
			}),
		},
	).Create(&gormMoto)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			r.log.Debug("MotoRepo_Update: record not found", "err", result.Error, "id", gormMoto.ID)	
			return domain.Moto{}, domain.RecordNotFound
		}
		r.log.Error("MotoRepo_Update: internal error", "err", result.Error)
		return domain.Moto{}, fmt.Errorf("%w: internal error: %v", domain.InternalError, result.Error)
	}

	if result.RowsAffected == 0 {
		r.log.Error("MotoRepo_Update: RowsAffected == 0, internal error", "id", gormMoto.ID)
		return domain.Moto{}, domain.InternalError
	}
	r.log.Debug("MotoRepo_Update: record update success!", "id", gormMoto.ID)

	var update GormMoto
	if err := r.db.WithContext(ctx).Where("id = ?", gormMoto.ID).First(&update).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.log.Error("MotoRepo_Update: updated record not found!", "id", gormMoto.ID)
			return domain.Moto{}, domain.RecordNotFound
		}
		r.log.Error("MotoRepo_Update: internal error befor get", "id", gormMoto.ID, "err", err)
		return domain.Moto{}, fmt.Errorf("%w: internal error: %v", domain.InternalError, err)
	}

	domainMoto := toDomainMoto(update)

	r.log.Debug("MotoRepo_Update: End!")
	return domainMoto, nil
}

func (r *motoRepo) Delete(ctx context.Context, motoID uint) error {
	r.log.Debug("MotoRepo_Delete: Start!")	

	var gormMoto GormMoto

	err := r.db.WithContext(ctx).Where("id = ?", motoID).Delete(&gormMoto).Error
	if err != nil {
		r.log.Error("MotoRepo_Delete: delete record error", "err", err)
		return fmt.Errorf("%w, delete record error: %v", domain.InternalError, err)
	}

	r.log.Debug("MotoRepo_Delete: End!")
	return nil
}

func (r *motoRepo) GetMotosByFilter(
	ctx context.Context, 
	filter domain.MotoFilter,
) ([]domain.Moto, error) {
	r.log.Debug("MotoRepo_GetMotosByFilter: Start!")

	var gormMotos []GormMoto
	err := r.db.WithContext(ctx).
		Scopes(MotoFilterScope(filter)).
		Find(&gormMotos).Error
	
	if err != nil {
		return nil, fmt.Errorf("list motos error: %w", err)
	}

	r.log.Debug("MotoRepo_GetMotosByFilter: End!")

	var domainMotos []domain.Moto

	for _, gormMoto := range gormMotos {
		domainMoto := toDomainMoto(gormMoto)
		domainMotos = append(domainMotos, domainMoto)	
	}

	return domainMotos, nil
}

func MotoFilterScope(f domain.MotoFilter) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if f.EngineSizeMin != nil {
			db = db.Where("engine_size >= ?", *f.EngineSizeMin)
		}
		if f.EngineSizeMax != nil {
			db = db.Where("engine_size < ?", *f.EngineSizeMax)
		}

		if f.YearMin != nil {
			db = db.Where("year >= ?", *f.YearMin)
		}
		if f.YearMax != nil {
			db = db.Where("year < ?", *f.YearMax)
		}

		if f.MileageMin != nil {
			db = db.Where("mileage >= ?", *f.MileageMin)
		}
		if f.MileageMax != nil {
			db = db.Where("mileage < ?", *f.MileageMax)
		}

		if f.PriceMin != nil {
			db = db.Where("price >= ?", *f.PriceMin)
		}
		if f.PriceMax != nil {
			db = db.Where("price <= ?", *f.PriceMax)
		}

		if f.MotoType != "-" {
			db = db.Where("moto_type = ?", f.MotoType)
		}

		return db
	}
}
