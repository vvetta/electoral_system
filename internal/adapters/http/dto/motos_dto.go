package dto

import (
	"github.com/vvetta/electoral_system/internal/domain"
)

type RequestGetMotos struct {
    EngineSizeOption int    `json:"engine_size_option"`
    YearOption       int    `json:"year_option"`
    MileageOption    int    `json:"mileage_option"`
    PriceMax *int64  `json:"price_max"`
    MotoType string  `json:"moto_type"`
}

type ResponseGetMotos struct {
	Motos []domain.Moto `json:"motos"`
}
