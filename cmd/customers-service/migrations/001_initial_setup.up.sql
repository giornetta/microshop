CREATE TABLE IF NOT EXISTS customers (
    customer_id UUID PRIMARY KEY,
    
    email       TEXT UNIQUE NOT NULL,
    first_name  TEXT NOT NULL,
    last_name   TEXT NOT NULL,

    shipping_country    TEXT,
    shipping_city       TEXT,
    shipping_zipcode    TEXT,
    shipping_street     TEXT
);
