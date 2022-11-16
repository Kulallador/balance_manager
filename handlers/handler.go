package handlers

import (
	"balance_manager/dbmanager"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Handlers struct {
	db *dbmanager.PostgresDB
}

func CreateHandlers(db dbmanager.PostgresDB) *Handlers {
	return &Handlers{db: &db}
}

func (h *Handlers) GetUserBalance(w http.ResponseWriter, r *http.Request) {
	var transaction dbmanager.BalanceTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Printf("Decode json: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("Decode json: %v", err.Error()), http.StatusBadRequest)
		return
	}

	balance, err := h.db.GetUserBalance(r.Context(), transaction.UserID)
	if err != nil {
		log.Printf("GetUserBalance: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("GetUserBalance: %v", err.Error()), http.StatusBadRequest)
		return
	}

	transaction.Money = balance

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(transaction)
}

func (h *Handlers) IncUserBalance(w http.ResponseWriter, r *http.Request) {
	var transaction dbmanager.BalanceTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Printf("Decode json: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("Decode json: %v", err.Error()), http.StatusBadRequest)
		return
	}

	err := h.db.IncrementMoney(r.Context(), transaction)
	if err != nil {
		log.Printf("IncUserBalance: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("IncUserBalance: %v", err.Error()), http.StatusBadRequest)
		return
	}
}

func (h *Handlers) DecUserBalance(w http.ResponseWriter, r *http.Request) {
	var transaction dbmanager.BalanceTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Printf("Decode json: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("Decode json: %v", err.Error()), http.StatusBadRequest)
		return
	}

	err := h.db.DecrementMoney(r.Context(), transaction)
	if err != nil {
		log.Printf("DecUserBalance: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("DecUserBalance: %v", err.Error()), http.StatusBadRequest)
		return
	}
}

func (h *Handlers) TranslationMoney(w http.ResponseWriter, r *http.Request) {
	var transaction dbmanager.TranslateTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Printf("Decode json: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("Decode json: %v", err.Error()), http.StatusBadRequest)
		return
	}

	err := h.db.TranslationMoney(r.Context(), transaction)
	if err != nil {
		log.Printf("TranslationMoney: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("TranslationMoney: %v", err.Error()), http.StatusBadRequest)
		return
	}
}

func (h *Handlers) ReserveMoney(w http.ResponseWriter, r *http.Request) {
	var transaction dbmanager.ReserveTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Printf("Decode json: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("Decode json: %v", err.Error()), http.StatusBadRequest)
		return
	}

	err := h.db.ReserveMoney(r.Context(), transaction)
	if err != nil {
		log.Printf("ReserveMoney: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("ReserveMoney: %v", err.Error()), http.StatusBadRequest)
		return
	}
}

func (h *Handlers) GetReserveBalance(w http.ResponseWriter, r *http.Request) {
	var transaction dbmanager.ReserveTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Printf("Decode json: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("Decode json: %v", err.Error()), http.StatusBadRequest)
		return
	}

	balance, err := h.db.GetReserveBalance(r.Context(), transaction)
	if err != nil {
		log.Printf("GetReserveBalance: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("GetReserveBalance: %v", err.Error()), http.StatusBadRequest)
		return
	}

	transaction.Money = balance

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(transaction)
}

func (h *Handlers) DecReservedMoney(w http.ResponseWriter, r *http.Request) {
	var transaction dbmanager.ReserveTransaction
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		log.Printf("Decode json: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("Decode json: %v", err.Error()), http.StatusBadRequest)
		return
	}

	err := h.db.DecReservedMoney(r.Context(), transaction)
	if err != nil {
		log.Printf("DecReservedMoney: %v\n", err.Error())
		http.Error(w, fmt.Sprintf("DecReservedMoney: %v", err.Error()), http.StatusBadRequest)
		return
	}
}