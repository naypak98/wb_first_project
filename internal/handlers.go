package internal

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

func GetOrderHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		orderID := r.URL.Path[len("/order/"):]
		log.Println("Запрос заказа:", orderID)

		// 1. Кэш
		if order, ok := GetFromCache(orderID); ok {
			log.Println("Нашёл в кэше:", orderID)
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(order)
			return
		}

		// 2. БД
		order, err := GetOrderFromDB(db, orderID)
		if err == sql.ErrNoRows {
			log.Println("Заказ не найден в БД:", orderID)
			http.Error(w, "Order not found", http.StatusNotFound) // 404
			return
		} else if err != nil {
			log.Println("Ошибка при чтении из БД:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError) // 500
			return
		}

		log.Println("Нашёл в БД:", orderID)

		// 3. Кэшируем
		SaveToCache(*order)

		// 4. Отдаём JSON
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(order); err != nil {
			log.Println("Ошибка кодирования JSON:", err)
		}
	}
}
