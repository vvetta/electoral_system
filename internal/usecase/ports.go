package usecase

import (
	"context"

	"github.com/vvetta/electoral_system/internal/domain"
)

type MotoParser interface {
	GetAllMoto() ([]domain.Moto, error)
}

type MotoRepo interface {
	Create(ctx context.Context, moto domain.Moto) (domain.Moto, error)
	Read(ctx context.Context, motoID uint) (domain.Moto, error)
	Update(ctx context.Context, moto domain.Moto) (domain.Moto, error)
	Delete(ctx context.Context, motoID uint) error

	GetMotosByFilter(ctx context.Context, filter domain.MotoFilter) ([]domain.Moto, error)
}

type MotoService interface {
	GetMoto(ctx context.Context, motoID uint) (domain.Moto, error)
	GetAllMoto(ctx context.Context) (domain.Moto, error)
	ParseAndUpdateAllMoto(ctx context.Context) ([]domain.Moto, error)
	UpdateMoto(ctx context.Context, moto domain.Moto) (domain.Moto, error)
	DeleteMoto(ctx context.Context, motoID uint) error

	//TODO offset и limit делать не будут, но вообще он тут нужен!
	GetMotosByFilter(ctx context.Context, filter domain.MotoFilter) ([]domain.Moto, error)
}

type Logger interface {
	Info(msg string, kv ...any)
	Debug(msg string, kv ...any)
	Error(msg string, kv ...any)
}
