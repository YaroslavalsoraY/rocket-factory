-- +goose Up
CREATE TABLE users (
    id UUID PRIMARY KEY,
    login TEXT UNIQUE,
    password TEXT,
    email TEXT UNIQUE,
    notification_methods JSONB
);