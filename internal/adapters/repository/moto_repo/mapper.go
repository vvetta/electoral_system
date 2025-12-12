package motorepo

import (
	"github.com/vvetta/electoral_system/internal/domain"
)

func toDomainMoto(moto GormMoto) domain.Moto {
	return domain.Moto{
		ID: moto.ID,
		Name: moto.Name,
		Year: moto.Year,
		Mileage: moto.Mileage,
		EngineSize: moto.EngineSize,
		MotoType: moto.MotoType,
		Location: moto.Location,
		Price: moto.Price,
		CreatedAt: moto.CreatedAt,
		UpdatedAt: moto.UpdatedAt,
	}
}

func toGormMoto(moto domain.Moto) GormMoto {
	return GormMoto{
		ID: moto.ID,
		Name: moto.Name,
		Year: moto.Year,
		Mileage: moto.Mileage,
		EngineSize: moto.EngineSize,
		MotoType: moto.MotoType,
		Location: moto.Location,
		Price: moto.Price,
		CreatedAt: moto.CreatedAt,
		UpdatedAt: moto.UpdatedAt,
	}
}

