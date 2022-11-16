package server

import (
	"balance_manager/dbmanager"
	"balance_manager/handlers"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type Server struct {
	db   *dbmanager.PostgresDB
	addr string
}

func CreateServer(db dbmanager.PostgresDB, addr string) *Server {
	return &Server{
		db:   &db,
		addr: addr,
	}
}

func (s *Server) setupRoutes() chi.Router {
	r := chi.NewRouter()
	h := handlers.CreateHandlers(*s.db)

	r.Get("/balance", h.GetUserBalance)
	r.Post("/balance/inc", h.IncUserBalance)
	r.Post("/balance/dec", h.DecUserBalance)
	r.Post("/balance/translate", h.TranslationMoney)
	r.Get("/reserve", h.GetReserveBalance)
	r.Post("/reserve/inc", h.ReserveMoney)
	r.Post("/reserve/dec", h.DecReservedMoney)
	return r
}

func (s *Server) Run() error {
	http.Handle("/", s.setupRoutes())
	log.Println("Start working")
	return http.ListenAndServe(s.addr, nil)
}
