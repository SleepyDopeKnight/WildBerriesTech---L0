package main

import (
	readDB "L0/database"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
)

const (
	db_host     = "localhost"
	db_port     = "5432"
	db_user     = "jojo"
	db_password = "123"
	db_name     = "order_db"
)

func main() {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	natsStreamConnection, err := stan.Connect("test-cluster", "consumer", stan.NatsURL(stan.DefaultNatsURL))
	if err != nil {
		log.Fatal(err)
	}
	_, err = natsStreamConnection.Subscribe("orders", func(message *stan.Msg) {
		//log.Printf("Received a message: %s\n", string(message.Data))
		orders, _ := readDB.FileDeserialize(message.Data)
		FillDatabase(orders, db)
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = natsStreamConnection.Subscribe("id", func(message *stan.Msg) {
		log.Printf("Received a message: %s\n", string(message.Data))
		wantedOrder := readDB.Orders{}
		var item readDB.Items
		orderPaymentDeliveryRows, err := db.Query("select * from orders join delivery d on orders.order_uid = d.order_uid "+
			"join payment p on orders.order_uid = p.transaction "+
			"WHERE orders.order_uid = $1 and orders.order_uid = p.transaction;", message.Data)
		for orderPaymentDeliveryRows.Next() {
			err = orderPaymentDeliveryRows.Scan(&wantedOrder.OrderUid,
				&wantedOrder.TrackNumber, &wantedOrder.Locale, &wantedOrder.InternalSignature, &wantedOrder.Entry,
				&wantedOrder.CustomerId, &wantedOrder.DeliveryService, &wantedOrder.Shardkey, &wantedOrder.SmId,
				&wantedOrder.DateCreated, &wantedOrder.OofShard, &wantedOrder.Delivery.OrderUid, &wantedOrder.Delivery.Name,
				&wantedOrder.Delivery.Phone, &wantedOrder.Delivery.Zip, &wantedOrder.Delivery.City, &wantedOrder.Delivery.Address,
				&wantedOrder.Delivery.Region, &wantedOrder.Delivery.Email, &wantedOrder.Payment.Transaction,
				&wantedOrder.Payment.RequestId, &wantedOrder.Payment.Currency, &wantedOrder.Payment.Provider,
				&wantedOrder.Payment.Amount, &wantedOrder.Payment.PaymentDt, &wantedOrder.Payment.Bank,
				&wantedOrder.Payment.DeliveryCost, &wantedOrder.Payment.GoodsTotal, &wantedOrder.Payment.CustomFee)
		}
		defer orderPaymentDeliveryRows.Close()
		itemsRows, err := db.Query("select * from items where track_number = (select track_number from orders where order_uid = $1);", message.Data)
		for itemsRows.Next() {
			err = itemsRows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size,
				&item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
		}
		wantedOrder.Items = append(wantedOrder.Items, item)
		defer itemsRows.Close()
		log.Println(wantedOrder)
		outgoingOrder, err := json.Marshal(wantedOrder)
		natsStreamConnection.Publish("data", []byte(outgoingOrder))
		if err != nil {
			log.Fatal(err)
		}

	})
	//err = natsStreamConnection.Close()
	//if err != nil {
	//	log.Fatal(err)
	//}

	select {}
}

func FillDatabase(orders *readDB.Orders, db *sql.DB) {
	for i := 0; i < len(orders.Items); i++ {
		_, _ = db.Exec("INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			orders.Items[i].ChrtId, orders.Items[i].TrackNumber, orders.Items[i].Price, orders.Items[i].Rid, orders.Items[i].Name, orders.Items[i].Sale,
			orders.Items[i].Size, orders.Items[i].TotalPrice, orders.Items[i].NmId, orders.Items[i].Brand, orders.Items[i].Status)
	}
	_, _ = db.Exec("INSERT INTO orders (order_uid, track_number, entry ,locale, internal_signature, customer_id, delivery_service, shardkey ,sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		orders.OrderUid, orders.TrackNumber, orders.Entry, orders.Locale, orders.InternalSignature, orders.CustomerId,
		orders.DeliveryService, orders.Shardkey, orders.SmId, orders.DateCreated, orders.OofShard)
	_, _ = db.Exec("INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		orders.Delivery.OrderUid, orders.Delivery.Name, orders.Delivery.Phone, orders.Delivery.Zip, orders.Delivery.City, orders.Delivery.Address,
		orders.Delivery.Region, orders.Delivery.Email)
	_, _ = db.Exec("INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		orders.Payment.Transaction, orders.Payment.RequestId, orders.Payment.Currency, orders.Payment.Provider, orders.Payment.Amount, orders.Payment.PaymentDt,
		orders.Payment.Bank, orders.Payment.DeliveryCost, orders.Payment.GoodsTotal, orders.Payment.CustomFee)
}

//func FindData(string id) {
//
//}
