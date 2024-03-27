create table if not exists items
(
    chrt_id integer,
    track_number varchar primary key,
    price integer,
    rid varchar,
    name varchar,
    sale integer,
    size varchar,
    total_price integer,
    nm_id integer,
    brand varchar,
    status integer
    );

create table if not exists orders
(
    order_uid varchar primary key,
    track_number varchar,
    entry varchar,
    locale varchar,
    internal_signature varchar,
    customer_id varchar,
    delivery_service varchar,
    shardkey varchar,
    sm_id integer,
    date_created timestamp,
    oof_shard varchar,

    constraint fk_orders_track_number foreign key (track_number) references items (track_number)
);

create table if not exists delivery
(
    order_uid varchar,
    name varchar,
    phone varchar,
    zip varchar,
    city varchar,
    address varchar,
    region varchar,
    email varchar,

    constraint fk_orders_order_uid foreign key (order_uid) references orders (order_uid)
);

create table if not exists payment
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
    custom_fee integer,

    constraint fk_orders_transaction foreign key (transaction) references orders (order_uid)
);
