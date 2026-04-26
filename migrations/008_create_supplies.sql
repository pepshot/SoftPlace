CREATE TABLE supplies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    supplier_id UUID NOT NULL,
    code VARCHAR(20) NOT NULL UNIQUE,
    date DATE NOT NULL,
    price DOUBLE PRECISION NOT NULL DEFAULT 0 CHECK (price >= 0),

    CONSTRAINT fk_supplies_supplier
        FOREIGN KEY (supplier_id)
        REFERENCES suppliers(id)
        ON DELETE CASCADE
);