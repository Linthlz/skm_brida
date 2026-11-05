CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE users (
    internal_id BIGSERIAL PRIMARY KEY,
    public_id UUID NOT NULL DEFAULT gen_random_uuid (),
    name varchar(255) not NULL,
    email VARCHAR(255) not NULL,
    password TEXT NOT NULL,
    phone_number VARCHAR(50),
    institution TEXT,
    institution_address TEXT,
    position varchar(100),
    role varchar(50) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    CONSTRAINT user_public_id_unique UNIQUE (public_id),
    CONSTRAINT user_email_unique UNIQUE (email)
)