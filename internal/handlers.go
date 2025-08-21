package internal

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func GetOrderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Path[len("/order/"):]
		order, ok := GetFromCache(orderID)
		if !ok {
			// TODO: достать из БД при необходимости
			http.Error(w, "Order not found", http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(order)
	}
}
