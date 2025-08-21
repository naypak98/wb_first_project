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
	internal.LoadCacheFromDB(db)

	http.Handle("/order/", internal.WithCORS(internal.GetOrderHandler(db)))
	log.Println("HTTP сервер запущен на :8081")
	go http.ListenAndServe(":8081", nil)

	go internal.ConsumeOrders("localhost:9092", "orders", "order-group", db)

	select {}
}
