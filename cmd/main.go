package main

import (
	"log"
	"net/http"
	"wb_first_project/internal"
)

func main() {
	db, err := internal.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// HTTP сервер
	http.HandleFunc("/order/", internal.GetOrderHandler(db))
	log.Println("HTTP сервер запущен на :8081")
	go http.ListenAndServe(":8081", nil)

	// Kafka consumer
	go internal.ConsumeOrders("localhost:9092", "orders", "order-group", db)

	// Блокировка main
	select {}
}
