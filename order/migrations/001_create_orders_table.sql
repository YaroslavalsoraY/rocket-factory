-- +goose Up
CREATE TABLE orders (
    id uuid PRIMARY KEY,
    user_uuid uuid NOT NULL,
    part_uuids uuid[],
    total_price numeric(10, 2),
    transactional_uuid uuid,
    payment_method text,
    status text
);