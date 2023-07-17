CREATE TABLE IF NOT EXISTS products (
    product_id  UUID    PRIMARY KEY,
    name        TEXT    UNIQUE NOT NULL,
    description TEXT    NOT NULL,
    price       REAL    NOT NULL CHECK (price > 0),
    amount      INT     NOT NULL CHECK (amount >= 0)
);
