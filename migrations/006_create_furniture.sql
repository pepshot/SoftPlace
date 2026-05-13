CREATE TABLE furniture (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        VARCHAR(100) NOT NULL,
    code        VARCHAR(20) NOT NULL UNIQUE,
    price       DOUBLE PRECISION NOT NULL DEFAULT 0 CHECK (price >= 0),
    stock_count INTEGER NOT NULL DEFAULT 0 CHECK (stock_count >= 0)
);