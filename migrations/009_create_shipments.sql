CREATE TABLE shipments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    customer_id UUID NOT NULL,
    code VARCHAR(20) NOT NULL UNIQUE,
    date DATE NOT NULL,
    price DOUBLE PRECISION NOT NULL DEFAULT 0 CHECK (price >= 0),

    CONSTRAINT fk_shipments_customer
        FOREIGN KEY (customer_id)
        REFERENCES customers(id)
        ON DELETE CASCADE
);