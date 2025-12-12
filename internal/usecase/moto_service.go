package usecase

import (
	"context"

	"github.com/vvetta/electoral_system/internal/domain"
)

type motoService struct {
	log Logger
	motoRepo MotoRepo
	motoParser MotoParser
}

func NewMotoService(
	log Logger, 
	motoRepo MotoRepo, 
	motoParser MotoParser,
) MotoService {
	return &motoService{
		log: log,
		motoRepo: motoRepo,
		motoParser: motoParser,
	}
}

func (s *motoService) GetMoto(
	ctx context.Context, 
	motoID uint,
) (domain.Moto, error) {
	return s.motoRepo.Read(ctx, motoID)
}

func (s *motoService)	GetAllMoto(
	ctx context.Context,
) (domain.Moto, error) {
	return domain.Moto{}, nil
}

func (s *motoService) ParseAndUpdateAllMoto(
	ctx context.Context,
) ([]domain.Moto, error) {
	s.log.Debug("MotoService_ParseAndUpdateAllMoto: Start!")

	motos, err := s.motoParser.GetAllMoto()
	if err != nil {
		s.log.Error("MotoService_ParseAndUpdateAllMoto: parsing moto error", "err", err)
		return nil, err
	}

	//TODO тут можно использовать канал с мотоциклами и сделать несколько воркеров.
	//TODO функция хрень. Мотоциклы приходят без id, поэтому при каждом запуске будет происходить запись
	//можно попробовать считать хеш от всех полей и сделать это поле уникальным.
	var updatedMotos []domain.Moto
	for i, moto := range motos {
		moto.ID = uint(i)
		updatedMoto, err := s.motoRepo.Update(ctx, moto)
		if err != nil {
			s.log.Error("MotoService_ParseAndUpdateAllMoto: update moto error", "id", moto.ID, "err", err)
			continue
		}

		updatedMotos = append(updatedMotos, updatedMoto)
	}

	s.log.Debug("MotoService_ParseAndUpdateAlLMoto: End!")
	return updatedMotos, nil
}

func (s *motoService) UpdateMoto(
	ctx context.Context, 
	moto domain.Moto,
) (domain.Moto, error) {
	return s.motoRepo.Update(ctx, moto)
}

func (s *motoService) DeleteMoto(
	ctx context.Context, 
	motoID uint,
) error {
	return s.motoRepo.Delete(ctx, motoID)
}

func (s *motoService) GetMotosByFilter(
	ctx context.Context,
	filter domain.MotoFilter,
) ([]domain.Moto, error) {
	return s.motoRepo.GetMotosByFilter(ctx, filter)
}
