CREATE TABLE IF NOT EXISTS deliveries (
    id UUID PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    phone VARCHAR(32) NOT NULL,
    zip VARCHAR(32) NOT NULL,
    city VARCHAR(64) NOT NULL,
    address VARCHAR(256) NOT NULL,
    region VARCHAR(64) NOT NULL,
    email VARCHAR(256)
);

CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY,
    request_id VARCHAR(128),
    currency VARCHAR(32) NOT NULL,
    provider VARCHAR(128) NOT NULL,
    amount BIGINT NOT NULL,
    payment_dt TIMESTAMP NOT NULL,
    bank VARCHAR(128) NOT NULL,
    delivery_cost BIGINT NOT NULL,
    goods_total BIGINT NOT NULL,
    custom_fee BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS items (
    chrt_id BIGINT PRIMARY KEY,
    track_number VARCHAR(128),
    price BIGINT NOT NULL,
    rid UUID NOT NULL,
    name VARCHAR(256) NOT NULL,
    sale BIGINT NOT NULL,
    size VARCHAR(32) NOT NULL,
    total_price BIGINT NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(128) NOT NULL,
    status BIGINT NOT NULL
);


CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(32) NOT NULL,
    delivery_id UUID NOT NULL,
    payment_id UUID NOT NULL,
    locale VARCHAR(32) NOT NULL,
    internal_signature VARCHAR(128),
    customer_id UUID NOT NULL,
    delivery_service VARCHAR(128),
    shardkey VARCHAR(32) NOT NULL,
    sm_id BIGINT NOT NULL,
    date_created TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    oof_shard VARCHAR(32) NOT NULL
);