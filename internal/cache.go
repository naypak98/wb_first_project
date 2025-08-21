package internal

import (
	"database/sql"
	"log"
	"sync"
)

var cache = struct {
	sync.RWMutex
	data map[string]Order
}{data: make(map[string]Order)}

func SaveToCache(order Order) {
	cache.Lock()
	defer cache.Unlock()
	cache.data[order.OrderUID] = order
}

func GetFromCache(orderUID string) (Order, bool) {
	cache.RLock()
	defer cache.RUnlock()
	order, ok := cache.data[orderUID]
	return order, ok
}

func LoadCacheFromDB(db *sql.DB) {
	rows, err := db.Query(`SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
                                   delivery_service, shardkey, sm_id, date_created, oof_shard
                            FROM orders`)
	if err != nil {
		log.Println("Ошибка загрузки кеша из БД:", err)
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var o Order
		err := rows.Scan(&o.OrderUID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature,
			&o.CustomerID, &o.DeliveryService, &o.Shardkey, &o.SmID, &o.DateCreated, &o.OofShard)
		if err != nil {
			log.Println("Ошибка чтения строки из БД:", err)
			continue
		}

		err = db.QueryRow(`SELECT name, phone, zip, city, address, region, email
                           FROM delivery WHERE order_uid=$1`, o.OrderUID).
			Scan(&o.Delivery.Name, &o.Delivery.Phone, &o.Delivery.Zip, &o.Delivery.City,
				&o.Delivery.Address, &o.Delivery.Region, &o.Delivery.Email)
		if err != nil {
			log.Println("Ошибка загрузки Delivery для", o.OrderUID, ":", err)
		}

		err = db.QueryRow(`SELECT transaction, request_id, currency, provider, amount,
                                  payment_dt, bank, delivery_cost, goods_total, custom_fee
                           FROM payment WHERE order_uid=$1`, o.OrderUID).
			Scan(&o.Payment.Transaction, &o.Payment.RequestID, &o.Payment.Currency,
				&o.Payment.Provider, &o.Payment.Amount, &o.Payment.PaymentDt,
				&o.Payment.Bank, &o.Payment.DeliveryCost, &o.Payment.GoodsTotal,
				&o.Payment.CustomFee)
		if err != nil {
			log.Println("Ошибка загрузки Payment для", o.OrderUID, ":", err)
		}

		itemRows, err := db.Query(`SELECT chrt_id, track_number, price, rid, name, sale,
                                           size, total_price, nm_id, brand, status
                                    FROM items WHERE order_uid=$1`, o.OrderUID)
		if err != nil {
			log.Println("Ошибка загрузки Items для", o.OrderUID, ":", err)
		} else {
			for itemRows.Next() {
				var it Item
				err := itemRows.Scan(&it.ChrtID, &it.TrackNumber, &it.Price, &it.Rid, &it.Name,
					&it.Sale, &it.Size, &it.TotalPrice, &it.NmID, &it.Brand, &it.Status)
				if err != nil {
					log.Println("Ошибка сканирования Item для", o.OrderUID, ":", err)
					continue
				}
				o.Items = append(o.Items, it)
			}
			itemRows.Close()
		}

		SaveToCache(o)
		count++
	}

	log.Printf("Cache loaded: %d orders\n", count)
}

func CacheSize() int {
	cache.RLock()
	defer cache.RUnlock()
	return len(cache.data)
}
