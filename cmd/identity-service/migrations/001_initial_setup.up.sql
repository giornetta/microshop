CREATE TABLE IF NOT EXISTS identities (
    identity_id UUID PRIMARY KEY,
    email       TEXT UNIQUE NOT NULL,
    roles TEXT[],
    password_hash TEXT NOT NULL
);
