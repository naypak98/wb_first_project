package internal

import (
	"database/sql"
)

func GetOrderFromDB(db *sql.DB, orderUID string) (*Order, error) {
	var order Order

	err := db.QueryRow(`
        SELECT order_uid, track_number, entry, locale, internal_signature, 
               customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard
        FROM orders WHERE order_uid = $1
    `, orderUID).Scan(
		&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerID, &order.DeliveryService,
		&order.Shardkey, &order.SmID, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`
        SELECT name, phone, zip, city, address, region, email
        FROM delivery WHERE order_uid = $1
    `, orderUID).Scan(
		&order.Delivery.Name, &order.Delivery.Phone, &order.Delivery.Zip,
		&order.Delivery.City, &order.Delivery.Address, &order.Delivery.Region,
		&order.Delivery.Email,
	)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow(`
        SELECT transaction, request_id, currency, provider, amount, payment_dt, 
               bank, delivery_cost, goods_total, custom_fee
        FROM payment WHERE order_uid = $1
    `, orderUID).Scan(
		&order.Payment.Transaction, &order.Payment.RequestID, &order.Payment.Currency,
		&order.Payment.Provider, &order.Payment.Amount, &order.Payment.PaymentDt,
		&order.Payment.Bank, &order.Payment.DeliveryCost, &order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)
	if err != nil {
		return nil, err
	}

	rows, err := db.Query(`
        SELECT chrt_id, track_number, price, rid, name, sale, size, 
               total_price, nm_id, brand, status
        FROM items WHERE order_uid = $1
    `, orderUID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item Item
		err := rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name,
			&item.Sale, &item.Size, &item.TotalPrice, &item.NmID, &item.Brand,
			&item.Status,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return &order, nil
}
