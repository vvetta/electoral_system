package httpserver

import (
	"net/http"

	"github.com/vvetta/electoral_system/internal/usecase"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer(
	motosSVC usecase.MotoService,
	lg usecase.Logger,
) *Server {
	mux := http.NewServeMux()

	motosHandler := NewMotosHandler(motosSVC, lg)
	motosHandler.Register(mux)

	mux.Handle("/", http.FileServer(http.Dir("web/")))

	return &Server{
		mux: mux,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	origin := r.Header.Get("Origin")

	// Разрешаем только фронт, но на dev это можно считать "все, кто нам нужен"
	if origin == "http://localhost:5173" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	} else {
		// на время разработки можно и всем:
		w.Header().Set("Access-Control-Allow-Origin", origin)
		w.Header().Set("Access-Control-Allow-Credentials", "true")
	}

	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	s.mux.ServeHTTP(w, r)
}

