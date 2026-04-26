CREATE TABLE suppliers (
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    login       VARCHAR(50) NOT NULL UNIQUE,
    password    TEXT NOT NULL,
    email       VARCHAR(100) NOT NULL UNIQUE
);