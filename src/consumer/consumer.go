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
	db := DBConnection()
	defer db.Close()
	natsStreamConnection, err := stan.Connect("test-cluster", "consumer", stan.NatsURL(stan.DefaultNatsURL))
	defer natsStreamConnection.Close()
	if err != nil {
		log.Fatal(err)
	}
	ChannelForGetJSON(natsStreamConnection, db)
	ChannelsForHandleIdDRequest(natsStreamConnection, db)

	select {}
}

func DBConnection() *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func ChannelForGetJSON(natsStreamConnection stan.Conn, db *sql.DB) {
	_, err := natsStreamConnection.Subscribe("orders", func(message *stan.Msg) {
		orders := readDB.FileDeserialize(message.Data)
		FillDatabase(orders, db)
	})
	if err != nil {
		log.Fatal(err)
	}
}

func ChannelsForHandleIdDRequest(natsStreamConnection stan.Conn, db *sql.DB) {
	_, err := natsStreamConnection.Subscribe("id", func(message *stan.Msg) {
		err := db.Ping()
		if err == nil {
			publicationWithConnectedDB(message, natsStreamConnection, db)
		} else {
			publicationWithDisconnectedDB(natsStreamConnection)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
}

func publicationWithConnectedDB(message *stan.Msg, natsStreamConnection stan.Conn, db *sql.DB) {
	wantedOrder := FindOrder(message, db)
	outgoingOrder, err := json.Marshal(wantedOrder)
	if err != nil {
		log.Fatal(err)
	}
	err = natsStreamConnection.Publish("data", []byte(outgoingOrder))
	if err != nil {
		log.Fatal(err)
	}
}

func publicationWithDisconnectedDB(natsStreamConnection stan.Conn) {
	var emptyOrder readDB.Orders
	outgoingOrder, err := json.Marshal(emptyOrder)
	err = natsStreamConnection.Publish("data", []byte(outgoingOrder))
	if err != nil {
		log.Fatal(err)
	}
}

func FillDatabase(orders *readDB.Orders, db *sql.DB) {
	for i := 0; i < len(orders.Items); i++ {
		_, err := db.Exec("INSERT INTO items (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
			orders.Items[i].ChrtId, orders.Items[i].TrackNumber, orders.Items[i].Price, orders.Items[i].Rid, orders.Items[i].Name, orders.Items[i].Sale,
			orders.Items[i].Size, orders.Items[i].TotalPrice, orders.Items[i].NmId, orders.Items[i].Brand, orders.Items[i].Status)
		if err != nil {
			log.Fatal(err)
		}
	}
	_, err := db.Exec("INSERT INTO orders (order_uid, track_number, entry ,locale, internal_signature, customer_id, delivery_service, shardkey ,sm_id, date_created, oof_shard) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)",
		orders.OrderUid, orders.TrackNumber, orders.Entry, orders.Locale, orders.InternalSignature, orders.CustomerId,
		orders.DeliveryService, orders.Shardkey, orders.SmId, orders.DateCreated, orders.OofShard)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO delivery (order_uid, name, phone, zip, city, address, region, email) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		orders.Delivery.OrderUid, orders.Delivery.Name, orders.Delivery.Phone, orders.Delivery.Zip, orders.Delivery.City, orders.Delivery.Address,
		orders.Delivery.Region, orders.Delivery.Email)
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO payment (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)",
		orders.Payment.Transaction, orders.Payment.RequestId, orders.Payment.Currency, orders.Payment.Provider, orders.Payment.Amount, orders.Payment.PaymentDt,
		orders.Payment.Bank, orders.Payment.DeliveryCost, orders.Payment.GoodsTotal, orders.Payment.CustomFee)
	if err != nil {
		log.Fatal(err)
	}
}

func FindOrder(message *stan.Msg, db *sql.DB) readDB.Orders {
	wantedOrder := readDB.Orders{}
	var item readDB.Items
	orderRows := RowsFromDB(db, message, "orders", "order_uid")
	for orderRows.Next() {
		err := orderRows.Scan(&wantedOrder.OrderUid, &wantedOrder.TrackNumber, &wantedOrder.Entry, &wantedOrder.Locale,
			&wantedOrder.InternalSignature, &wantedOrder.CustomerId, &wantedOrder.DeliveryService, &wantedOrder.Shardkey,
			&wantedOrder.SmId, &wantedOrder.DateCreated, &wantedOrder.OofShard)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer orderRows.Close()

	paymentRows := RowsFromDB(db, message, "payment", "transaction")
	for paymentRows.Next() {
		err := paymentRows.Scan(&wantedOrder.Payment.Transaction, &wantedOrder.Payment.RequestId, &wantedOrder.Payment.Currency,
			&wantedOrder.Payment.Provider, &wantedOrder.Payment.Amount, &wantedOrder.Payment.PaymentDt, &wantedOrder.Payment.Bank,
			&wantedOrder.Payment.DeliveryCost, &wantedOrder.Payment.GoodsTotal, &wantedOrder.Payment.CustomFee)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer paymentRows.Close()

	deliveryRows := RowsFromDB(db, message, "delivery", "order_uid")
	for deliveryRows.Next() {
		err := deliveryRows.Scan(&wantedOrder.Delivery.OrderUid, &wantedOrder.Delivery.Name, &wantedOrder.Delivery.Phone,
			&wantedOrder.Delivery.Zip, &wantedOrder.Delivery.City, &wantedOrder.Delivery.Address,
			&wantedOrder.Delivery.Region, &wantedOrder.Delivery.Email)
		if err != nil {
			log.Fatal(err)
		}
	}
	defer deliveryRows.Close()

	itemsRows, err := db.Query("select * from items where track_number = (select track_number from orders where order_uid = $1);", message.Data)
	if err != nil {
		log.Fatal(err)
	}
	defer itemsRows.Close()
	for itemsRows.Next() {
		err = itemsRows.Scan(&item.ChrtId, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale, &item.Size,
			&item.TotalPrice, &item.NmId, &item.Brand, &item.Status)
	}
	wantedOrder.Items = append(wantedOrder.Items, item)
	return wantedOrder
}

func RowsFromDB(db *sql.DB, message *stan.Msg, tableName, rowName string) *sql.Rows {
	query := fmt.Sprintf("select * from %s where %s = '%s'", tableName, rowName, message.Data)
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	return rows
}
