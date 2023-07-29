CREATE TABLE IF NOT EXISTS identities (
    identity_id UUID PRIMARY KEY,
    email       TEXT UNIQUE NOT NULL,
    roles TEXT[],
    password_hash TEXT NOT NULL
);

-- Probably not the most secure way, whatever...
INSERT INTO identities(identity_id, email, roles, password_hash)
VALUES (gen_random_uuid(), 'michelegiornetta@gmail.com', ARRAY ['Customer', 'Admin'], '$2a$10$Mm24cU/V7mLnrVehTF73AOjCPd8BXl083tywtW3ZUH0JaLjylgkZS');