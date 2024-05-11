package main

import (
	readDB "L0/database"
	"database/sql"
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
		rows, err := db.Query("select * from orders join public.delivery d on orders.order_uid = d.order_uid "+
			"join public.items i on i.track_number = orders.track_number "+
			"join public.payment p on orders.order_uid = p.transaction "+
			"WHERE orders.order_uid = $1", message.Data)
		for rows.Next() {
			fmt.Println(rows)
		}
		defer rows.Close()
		natsStreamConnection.Publish("data", message.Data)
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

//func FoundData(string id) {
//
//}
