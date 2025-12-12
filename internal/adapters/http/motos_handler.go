package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vvetta/electoral_system/internal/adapters/http/dto"
	"github.com/vvetta/electoral_system/internal/domain"
	"github.com/vvetta/electoral_system/internal/usecase"
)

type MotosHandler struct {
	svc usecase.MotoService
	lg usecase.Logger
}

func NewMotosHandler(
	svc usecase.MotoService, 
	lg usecase.Logger,
) *MotosHandler {
	return &MotosHandler{
		svc: svc,
		lg: lg,
	}
}

func (h *MotosHandler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/motos/getByFilter", h.handleGetMotos)
	mux.HandleFunc("POST /api/v1/motos/parseAndUpdate", h.handleParseAndUpdate)
}

func (h *MotosHandler) handleParseAndUpdate(
	w http.ResponseWriter,
	r *http.Request,
) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errorResponse{})	
		return
	}

	_, err := h.svc.ParseAndUpdateAllMoto(r.Context())
	if err != nil {
		writeError(w, http.StatusBadRequest, errorResponse{})
	}

	writeJSON(w, http.StatusOK, nil)
}

func (h *MotosHandler) handleGetMotos(
	w http.ResponseWriter, 
	r *http.Request,
) {
	h.lg.Debug("MotosHandler_GetMotos: Start!")

	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, errorResponse{})	
		return
	}

	var request dto.RequestGetMotos
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		writeError(w, http.StatusBadRequest, errorResponse{})
		return	
	}

	filter := domain.NewMotoFilter(
		request.EngineSizeOption,
		request.YearOption,
		request.MileageOption,
		request.PriceMax,
		request.MotoType,
	)

	motos, err := h.svc.GetMotosByFilter(r.Context(), filter)
	if err != nil {
		if errors.Is(err, domain.RecordNotFound) {
			writeError(w, http.StatusNotFound, errorResponse{})
			return
		} else if errors.Is(err, domain.InternalError) {
			writeError(w, http.StatusInternalServerError, errorResponse{})
			return
		}
	}

	response := dto.ResponseGetMotos{
		Motos: motos,
	}

	h.lg.Debug("MotosHandler_GetMotos: End!")
	writeJSON(w, http.StatusOK, response)
}
