package internal

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"

	"github.com/segmentio/kafka-go"
)

func ConsumeOrders(broker, topic, groupID string, db *sql.DB) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   topic,
		GroupID: groupID,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("Ошибка чтения: %v", err)
			continue
		}

		var order Order
		if err := json.Unmarshal(m.Value, &order); err != nil {
			log.Printf("Ошибка парсинга JSON: %v", err)
			continue
		}

		if err := SaveOrder(db, order); err != nil {
			log.Printf("Ошибка сохранения: %v", err)
			continue
		}

		SaveToCache(order)
		log.Printf("Заказ %s сохранён", order.OrderUID)
	}
}
