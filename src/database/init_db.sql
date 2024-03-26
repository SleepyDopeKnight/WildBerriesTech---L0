create table if not exist orders
(
    order_uid varchar primary key,
    track_number varchar primary key,
    delivery_name varchar primary key
    entry varchar,
    locale varchar,
    internal_signature varchar,
    customer_id varchar,
    delivery_service varchar,
    shardkey varchar,
    sm_id integer varchar,
    date_created timestamp,
    oof_shard varchar
);

create table if not exist delivery
(
    name varchar,
    phone varchar,
    zip varchar,
    city varchar,
    address varchar,
    region varchar,
    email varchar

    constraint fk_orders_name foreign key (name) references orders (delivery_name)
);

create table if not exist payment
(
    transaction varchar,
    request_id varchar,
    currency varchar,
    provider varchar,
    amount integer,
    payment_dt integer,
    bank varchar,
    delivery_cost integer,
    goods_total integer,
    custom_fee integer

    constraint fk_orders_transaction foreign key (transaction) references orders (order_uid)
);

create table if not exist items
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
    status integer

    constraint fk_orders_track_number foreign key (track_number) references orders (track_number)
);
