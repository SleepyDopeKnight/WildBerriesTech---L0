create table if not exists delivery
(
    order_uid varchar primary key,
    name varchar,
    phone varchar,
    zip varchar,
    city varchar,
    address varchar,
    region varchar,
    email varchar
);

create table if not exists payment
(
    transaction varchar primary key,
    request_id varchar,
    currency varchar,
    provider varchar,
    amount integer,
    payment_dt integer,
    bank varchar,
    delivery_cost integer,
    goods_total integer,
    custom_fee integer
);

create table if not exists orders
(
    order_uid varchar ,
    track_number varchar primary key,
    entry varchar,
    locale varchar,
    internal_signature varchar,
    customer_id varchar,
    delivery_service varchar,
    shardkey varchar,
    sm_id integer,
    date_created timestamp,
    oof_shard varchar,

    constraint fk_delivery_order_uid foreign key (order_uid) references delivery (order_uid),
    constraint fk_payment_transaction foreign key (order_uid) references payment (transaction)
    );

create table if not exists items
(
    chrt_id integer,
    track_number varchar,
    price integer,
    rid varchar,
    name varchar,
    sale integer,
    size varchar,
    total_price integer,
    nm_id integer,
    brand varchar,
    status integer,

    constraint fk_items_track_number foreign key (track_number) references orders (track_number)
    );