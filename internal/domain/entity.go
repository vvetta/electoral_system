package domain

import (
	"time"
)

type Moto struct {
	ID uint
	Name string
	Year int
	Mileage int
	EngineSize int
	MotoType string
	Location string
	Price int64
	CreatedAt *time.Time
	UpdatedAt *time.Time
}

type MotoFilter struct {
	EngineSizeMin *int   // cc
	EngineSizeMax *int   // cc

	PriceMin      *int64 // если понадобится
	PriceMax      *int64

	YearMin       *int
	YearMax       *int

	MileageMin    *int   // км
	MileageMax    *int   // км

	MotoType      string
}

func NewMotoFilter(
	engineSizeOption int,
	yearOption int,
	mileageOption int,
	priceMaxOption *int64,
	motoType string,
) MotoFilter {
	engineMin, engineMax := engineSizeRange(engineSizeOption)
	yearMin, yearMax := yearRange(yearOption)
	mileageMin, mileageMax := mileageRange(mileageOption)
	priceMin, priceMax := priceRange(priceMaxOption)

	return MotoFilter{
		EngineSizeMin: engineMin,
		EngineSizeMax: engineMax,
		YearMin:       yearMin,
		YearMax:       yearMax,
		MileageMin:    mileageMin,
		MileageMax:    mileageMax,
		PriceMin:      priceMin,
		PriceMax:      priceMax,
		MotoType:      motoType,
	}
}

func priceRange(maxOption *int64) (min, max *int64) {
	// min сейчас не нужен, оставим на будущее
	if maxOption == nil || *maxOption <= 0 {
		return nil, nil // любой
	}
	return nil, maxOption
}

func yearRange(option int) (min, max *int) {
	switch option {
	case 1: // 2020+
		minV := 2020
		return &minV, nil
	case 2: // 2015–2020
		minV, maxV := 2015, 2020
		return &minV, &maxV
	case 3: // 2010–2015
		minV, maxV := 2010, 2015
		return &minV, &maxV
	case 4: // 2000–2010
		minV, maxV := 2000, 2010
		return &minV, &maxV
	case 5: // любой
		return nil, nil
	default:
		return nil, nil
	}
}

func mileageRange(option int) (min, max *int) {
	switch option {
	case 1: // до 10k
		minV, maxV := 0, 10_000
		return &minV, &maxV
	case 2: // 10–30k
		minV, maxV := 10_000, 30_000
		return &minV, &maxV
	case 3: // 30–50k
		minV, maxV := 30_000, 50_000
		return &minV, &maxV
	case 4: // 50–100k
		minV, maxV := 50_000, 100_000
		return &minV, &maxV
	case 5: // любой
		return nil, nil
	default:
		return nil, nil
	}
}

func engineSizeRange(option int) (min, max *int) {
	switch option {
	case 1: // до 250
		minV, maxV := 0, 250
		return &minV, &maxV
	case 2: // 250–500
		minV, maxV := 250, 500
		return &minV, &maxV
	case 3: // 500–750
		minV, maxV := 500, 750
		return &minV, &maxV
	case 4: // 750–1000
		minV, maxV := 750, 1000
		return &minV, &maxV
	case 5: // 1000+
		minV := 1000
		return &minV, nil
	default: // 0 или что-то левое → любой
		return nil, nil
	}
}
