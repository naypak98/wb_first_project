package internal

import (
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	return sql.Open("postgres", "postgres://testuser:testpass@localhost:5432/orders_db?sslmode=disable")
}

func SaveOrder(db *sql.DB, order Order) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(`
        INSERT INTO orders (order_uid, track_number, entry, locale, internal_signature, customer_id, 
                            delivery_service, shardkey, sm_id, date_created, oof_shard)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
        ON CONFLICT (order_uid) DO NOTHING
    `,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService,
		order.Shardkey, order.SmID, order.DateCreated, order.OofShard,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
    `, order.OrderUID, order.Delivery.Name, order.Delivery.Phone,
		order.Delivery.Zip, order.Delivery.City, order.Delivery.Address,
		order.Delivery.Region, order.Delivery.Email,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
        INSERT INTO payment (order_uid, transaction, request_id, currency, provider, amount,
                             payment_dt, bank, delivery_cost, goods_total, custom_fee)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
    `, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID,
		order.Payment.Currency, order.Payment.Provider, order.Payment.Amount,
		order.Payment.PaymentDt, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
	)
	if err != nil {
		return err
	}

	for _, item := range order.Items {
		_, err = tx.Exec(`
            INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, sale,
                               size, total_price, nm_id, brand, status)
            VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)
        `, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price,
			item.Rid, item.Name, item.Sale, item.Size, item.TotalPrice,
			item.NmID, item.Brand, item.Status,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
